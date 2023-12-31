package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// DB connection
type UserDataBase struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type ApplicationConfig struct {
	Service  string       `json:"service"`
	DataBase UserDataBase `json:"database"`
}

var Config *ApplicationConfig
var configFile *string

func LoadConfiguration() error {
	if configFile == nil {
		return fmt.Errorf("config not initialized")
	}

	config := new(ApplicationConfig)
	file, err := os.Open(*configFile)

	if err != nil {
		return err
	}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return err
	}
	Config = config
	return nil
}

func GetConfig() ApplicationConfig {
	if Config != nil {
		return *Config
	}
	err := Init(".")
	if err != nil {
		panic(err)
	}
	return *Config
}

func setupReload() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for range c {
			log.Printf("Reloading configurations....\n")
			if configFile == nil {
				panic("Config file path not set!")
			}
			if err := LoadConfiguration(); err != nil {
				log.Printf("Error on reloading configurations, using old configuration : Error : %s\n", err.Error())
			}
		}
	}()
}

func Init(config string) error {
	if config == "" {
		config = "."
	}
	configFilepath := strings.TrimRight(config, "/") + "/config.json"
	configFile = &configFilepath
	setupReload()
	return LoadConfiguration()
}
