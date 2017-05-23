package core

import "golang.org/x/net/context"
import (
	"gopkg.in/olivere/elastic.v5"
	"encoding/json"
	"os"
	"fmt"
	"strings"
	"io"
	"github.com/pkg/errors"
)

type El struct {
	config *Config

	client  *elastic.Client
	scroll  *elastic.ScrollService
	context context.Context
}

func NewEl(config *Config) *El {
	return &El{config: config}
}

func (o *El) Init() (err error) {
	if client, err := elastic.NewClient(); err == nil {
		if o.config == nil {
			o.config = NewConfig("EL")
		}
		o.scroll = client.Scroll(o.config.Elastics[0].Indexes...).Size(6000)
		o.context = context.Background()
	}
	return
}

func (o *El) Close() (err error) {
	if o.client != nil {
		o.client.Stop()
		o.client = nil
		o.scroll = nil
		o.context.Done()
		o.context = nil
	}
	return
}

func (o *El) UpdateConfig(config *Config) {
	o.Close()
	o.config = config
}

func (o *El) Export(serviceName string, query Query) (err error) {
	if err = o.Init(); err != nil {
		return
	}

	var file *os.File
	if file, err = os.Create(o.config.Folder); err != nil {
		return
	}
	defer file.Close()

	var res *elastic.SearchResult
	o.scroll.Body("")

	for {
		if res, err = o.search(); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		for _, hit := range res.Hits.Hits {
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err == nil {
				for _, field := range o.config.Fields {
					if val, ok := item[field]; ok {
						if str, ok := val.(string); ok {
							fmt.Fprintf(file, "%v", strings.Trim(str, "\r\n"))
						} else {
							fmt.Fprintf(file, "%v", val)
						}
					} else {
						file.WriteString(" ")
					}
					file.WriteString(o.config.Separator)
				}
				file.WriteString("\n")
			} else {
				println(err)
			}
		}
	}
	o.scroll.Clear(o.context)
	return
}

func (o *El) search() (ret *elastic.SearchResult, err error) {
	ret, err = o.scroll.Do(o.context)
	if err == nil {
		if ret == nil {
			err = errors.New("expected results != nil; got nil")
		}
		if ret.Hits == nil {
			err = errors.New("expected results != nil; got nil")

		}
	}
	return
}
