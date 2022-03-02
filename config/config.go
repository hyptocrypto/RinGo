package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port            string `yaml:"port"`
	BufferSize      uint16 `yaml:"buffer_size"`
	OverwriteBuffer bool   `yaml:"overwrite_buffer"`
	ReadSize        uint16 `yaml:"read_size"`
	ReadInterval    uint16 `yaml:"read_interval"`
}

func LoadConfig(filename string) *Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	return config
}
