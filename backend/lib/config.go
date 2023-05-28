package lib

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Application `yaml:"application"`
}

type Application struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

func GetConfiguration(path string) (*Configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Can't close the config file")
		}
	}(f)

	var cfg Configuration
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
