package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	SteamcmdPath  string          `json:"steamcmd_path"`
	AdminPassword string          `json:"admin_password"`
	Accounts      map[int]Account `json:"accounts"`
	Games         []Game          `json:"games"`
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Game struct {
	AppID    string `json:"app_id"`
	Accounts []int  `json:"accounts"`
}

var config Config

func LoadConfig(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func SaveConfig(filename string) {
	data, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func GetConfig() *Config {
	return &config
}
