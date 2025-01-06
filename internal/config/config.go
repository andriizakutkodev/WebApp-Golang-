package config

import ce "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Env              string `yaml:"env" env-required:"true"`
	ConnectionString string `yaml:"connection_string" env-required:"true"`
}

func GetConfig() *Config {
	var cnf Config

	err := ce.ReadConfig("../../config/local.yaml", &cnf)

	if err != nil {
		panic(err.Error())
	}

	return &cnf
}
