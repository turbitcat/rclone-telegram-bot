package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config is the configuration of the program.
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

// processConfigError is a helper function to process errors in config.
func processConfigError(err error) {
	fmt.Println(err)
	// os.Exit(2)
}

// ReadFile reads the config file.
func (cfg *Config) ReadFile() {
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

// ReadEnv reads the environment variables.
func (cfg *Config) ReadEnv() {
	err := envconfig.Process("", cfg)
	if err != nil {
		processConfigError(err)
	}
}

// WriteFile writes the config file.
func (cfg *Config) WriteFile() {
	f, err := os.Create("config.yml")
	if err != nil {
		processConfigError(err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cfg)
	if err != nil {
		processConfigError(err)
	}
}

// ReadAll reads the config file and the environment variables.
func (cfg *Config) ReadAll() {
	cfg.ReadFile()
	cfg.ReadEnv()
}
