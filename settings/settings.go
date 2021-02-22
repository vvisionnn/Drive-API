package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
var confFolder = "./data"
var confPath = confFolder + "/configuration.json"

func init() {
	if _, err := os.Stat(confFolder); os.IsNotExist(err) {
		if err = os.Mkdir(confFolder, 0755); err != nil {
			log.Panic("create config folder error: ", err)
		}
	}
	confContent, err := ioutil.ReadFile(confPath)
	if err != nil {
		//log.Panic("Configuration file error")
		log.Println("no configuration file found")
		return
	}
	log.Println("found previous file, initial settings.")
	if err := json.Unmarshal(confContent, &CONF); err != nil {
		panic("configuration file error")
	}
}

func (conf *Configuration) Save() {
	confStr, err := json.Marshal(conf)
	if err != nil {
		log.Println("marshal config file error: ", err)
	}
	if err := ioutil.WriteFile(confPath, confStr, 0644); err != nil {
		log.Println("save config file error: ", err)
	}
}
