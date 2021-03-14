package main

import (
	"flag"
	"os"

	"github.com/barathvk/sender-receiver/receiver"
	"github.com/barathvk/sender-receiver/sender"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	AppId string `yaml:"appId"`
	Port  int    `yaml:"port"`
}

func loadConfig() Config {
	configFile, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	decoder := yaml.NewDecoder(configFile)
	var config Config
	decoder.Decode(&config)
	return config
}

func main() {
	isSender := flag.Bool("sender", false, "is sender")
	initialCount := flag.Int("initial-count", 0, "sender initial count")
	flag.Parse()
	config := loadConfig()
	if *isSender {
		sender.Start(config.AppId, *initialCount)
	} else {
		receiver.Start(config.AppId, config.Port)
	}
}
