package core

import (
	"gee/cfg"
	"encoding/json"
	"gee/lg"
)

var Log = lg.NewLogger("EL ")

type Config struct {
	Name  string   `default:El`
	Port  int   `default:5000`
	Debug bool   `default:true`

	Folder   string `default:export`
	Elastics []*Elastic

	ConfigFiles    []string
	ConfigSuffixes []string
}

type Elastic struct {
	Name      string
	Cluster   string`default:`
	Hosts     []string
	Port      int`default:9300`
	Settings  map[string]string
	Fields    []string
	Separator string`default:;`
	Indexes   []string
	Queries   []*Query
}

type Query struct {
	Name   string
	Source string
}

func NewElastic(name string) *Elastic {
	return &Elastic{
		Name:      name,
		Hosts:     []string{"localhost"},
		Settings:  make(map[string]string),
		Fields:    []string{"@logdate", "thread", "level", "logger", "dur", "kind", "message"},
		Separator: ";",
		Indexes:   []string{"logstash*"},

	}
}

func NewConfig(name string) *Config {
	return &Config{
		Name:     name,
		Elastics: []*Elastic{NewElastic(name)},
	}
}

func LoadConfig(files []string, suffixes []string) (ret *Config, err error) {
	ret = &Config{ConfigFiles: files}
	err = cfg.Unmarshal(ret, files, suffixes)

	if err == nil {
		ret.Print()
	}

	return
}

func (o *Config) Reload() (ret *Config, err error) {
	return LoadConfig(o.ConfigFiles, o.ConfigSuffixes)
}

func (o *Config) Print() {
	json, err := json.MarshalIndent(o, "", "\t")
	if err == nil {
		Log.Info(string(json))
	}
}
