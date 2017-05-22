package core

import "golang.org/x/net/context"
import "gopkg.in/olivere/elastic.v5"

type El struct {
	config *Config

	client  *elastic.Client
	search  *elastic.SearchService
	context context.Context
}

func NewEl(config *Config) *El {
	return &El{config: config}
}

func (o *El) Init() (err error) {
	if client, err := elastic.NewClient(); err == nil {
		if o.config == nil {
			o.config = NewElConfig()
		}
		o.search = client.Search(o.config.Indexes...)
		o.context = context.Background()
	}
	return
}

func (o *El) Close() (err error) {
	if o.client != nil {
		o.client.Stop()
		o.client = nil
		o.search = nil
		o.context.Done()
		o.context = nil
	}
	return
}

func (o *El) UpdateConfig(config *Config) {
	o.Close()
	o.config = config
}

func (o *El) Export() (err error) {
	if err = o.Init(); err == nil {
		var result *elastic.SearchResult
		o.search.Source(`{
  "query": {
    "match_all": {}
  }
}`)
		if result, err = o.search.Do(o.context); err == nil {
			println(result)
		}
	}
	return
}
