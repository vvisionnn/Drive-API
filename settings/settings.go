package settings

import (
	"encoding/json"
	"io/ioutil"
)

type configuration struct {
	Port             int      `json:"port"`
	AppId            string   `json:"app_id"`
	AppSecret        string   `json:"app_secret"`
	Redirect         string   `json:"redirect"`
	Scopes           []string `json:"scopes"`
	OauthEndpoint    string   `json:"oauth_endpoint"`
	Authority        string   `json:"authority"`
	OnedriveEndpoint string   `json:"onedrive_endpoint"`
}

var CONF = configuration{}

func init() {
	confPath := "./conf.json"
	confContent, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic("configuration file error")
	}
	if err := json.Unmarshal(confContent, &CONF); err != nil {
		panic("configuration file error")
	}
}
