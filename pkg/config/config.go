package config

import (
	"os"
	"path/filepath"
	"raspberry_chrome_controller/pkg/logger"

	"go.uber.org/zap"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	HostURL  string `yaml:"host_url"`
	Path     string `yaml:"path"`
	KeyToken string `yaml:"key_token"`
}

func InitConf(confFilePath string) *AppConfig {
	path, err := os.Getwd()
	if err != nil {
		logger.Fatal("Can't locate current dir", err)
	}

	confFile := path + confFilePath
	confFile = filepath.Clean(confFile)
	logger.Info("Try read config file", zap.String("path", confFile))

	file, errP := os.Open(confFile)
	if errP != nil {
		logger.Fatal("Can't open config file", errP)
	}
	defer file.Close()
	var cfg AppConfig
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Fatal("Invalid config file", err)
	}

	return &cfg
}
