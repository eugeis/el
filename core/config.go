package core

import (
	"gee/cfg"
	"encoding/json"
	"gee/lg"
)

var Log = lg.NewLogger("EL ")

type Config struct {
	Debug bool   `default:true`

	File           string
	ClusterName    string
	Hosts          []string
	Port           int
	Settings       map[string]string
	Fields         []string
	Separator      string
	Indexes        []string
	SearchSource   string

	ConfigFiles    []string
	ConfigSuffixes []string
}

func NewElConfig() *Config {
	return &Config{
		File:         "C:\\temp\\export.log",
		ClusterName:  "", Hosts: []string{"localhost"},
		Port:         9300,
		Settings:     make(map[string]string),
		Fields:       []string{"_all"},
		Separator:    " ",
		Indexes:      []string{"logstash-2017.05.18"},
		SearchSource: "",
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
