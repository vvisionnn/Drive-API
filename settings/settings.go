package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	ClientId         string   `json:"client_id"`
	ClientSecret     string   `json:"client_secret"`
	RedirectUrl      string   `json:"redirect_url"`
	Scopes           []string `json:"scopes"`
	OauthEndpoint    string   `json:"oauth_endpoint"`
	Authority        string   `json:"authority"`
	OnedriveEndpoint string   `json:"onedrive_endpoint"`
}

var CONF *Configuration

func init() {
	confPath := "./configuration.json"
	confContent, err := ioutil.ReadFile(confPath)
	if err != nil {
		//log.Panic("Configuration file error")
		log.Println("no configuration file found")
		return
	}
	if err := json.Unmarshal(confContent, &CONF); err != nil {
		panic("configuration file error")
	}
}
