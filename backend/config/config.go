package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	SQLite struct {
		Path string `yaml:"path"`
	} `yaml:"sqlite"`
	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"admin"`
	Devices []DeviceConfig `yaml:"devices"`
}

type DeviceConfig struct {
	Name               string `yaml:"name"`
	MAC                string `yaml:"mac"`
	ServiceUUID        string `yaml:"service_uuid"`
	CharacteristicUUID string `yaml:"characteristic_uuid"`
	Enabled            bool   `yaml:"enabled"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}
