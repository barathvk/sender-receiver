package main

import (
	"flag"

	"github.com/barathvk/sender-receiver/receiver"
	"github.com/barathvk/sender-receiver/sender"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppId string `envconfig:"APP_ID"`
	Port  int    `envconfig:"PORT"`
	Redis string `envconfig:"REDIS_ADDRESS"`
}

func loadConfig() Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}
	return config
}

func main() {
	godotenv.Load()
	isSender := flag.Bool("sender", false, "is sender")
	flag.Parse()
	config := loadConfig()
	if *isSender {
		sender.Start(config.AppId, config.Redis)
	} else {
		receiver.Start(config.AppId, config.Port, config.Redis)
	}
}
