package config

import ce "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Env              string `yaml:"env" env-required:"true"`
	ConnectionString string `yaml:"connection_string" env-required:"true"`
	Jwt              Jwt    `yaml:"jwt" env-required:"true"`
}

type Jwt struct {
	Secret  string `yaml:"secret" env-required:"true"`
	ExpTime uint   `yaml:"exp_time" env-required:"true"`
}

func GetConfig() *Config {
	var cfg Config

	err := ce.ReadConfig("../../config/local.yaml", &cfg)

	if err != nil {
		panic(err.Error())
	}

	return &cfg
}
