package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP     HTTPConfig `yaml:"http"`
	Postgres PSQLConfig `yaml:"postgres"`
}

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type PSQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

func ParseFlag() string {
	cfgPathPtr := flag.String("config", "", "Path to config")
	flag.Parse()

	return *cfgPathPtr
}

func MustLoad(cfgPath string) (*Config, error) {
	if cfgPath == "" {
		return nil, errors.New("config path is empty")
	}

	if _, err := os.Stat(cfgPath); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("config does not exists: %v", err)
	}

	var config Config
	if err := cleanenv.ReadConfig(cfgPath, &config); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	return &config, nil
}
