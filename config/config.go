package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BotToken     string `json:"bot_token"`
	GeminiApiKey string `json:"gemini_api_key"`
	ProxyUrl     string `json:"proxy_url"`
	MongoUrl     string `json:"mongo_url"`
}

var SetConfig *Config

func InitConfig() {
	//read config.json
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	//decode json
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&SetConfig)
	if err != nil {
		panic(err)
	}
}
