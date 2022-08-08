package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	TelegramBot struct {
		Token    string `yaml:"token" envconfig:"BOT_TOKEN"`
		AdminIDs string `yaml:"admin_ids" envconfig:"BOT_ADMINIDS"`
	} `yaml:"telegram-bot"`
	Rclone struct {
		Host     string `yaml:"host" envconfig:"RCLONE_HOST"`
		User     string `yaml:"user" envconfig:"RCLONE_USER"`
		Password string `yaml:"password" envconfig:"RCLONE_PASSWORD"`
	} `yaml:"rclone"`
	TempPath struct {
		Download string `yaml:"download" envconfig:"TEMP_DOWNLOAD_DIR_ROOT"`
	} `yaml:"temp-path"`
}

func processConfigError(err error) {
	fmt.Println(err)
	// os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processConfigError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processConfigError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processConfigError(err)
	}
}

var config_yml = `telegram-bot:
  token: 
  admin_ids: 
rclone: 
  host: 
  user: 
  password: 
temp-path:
  download: `
