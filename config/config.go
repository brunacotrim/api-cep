package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ViaCepCfg    ViaCepCfg
	CepAbertoCfg CepAbertoCfg
}

type ViaCepCfg struct {
	Enabled bool
}

type CepAbertoCfg struct {
	Enabled bool
	Token   string
}

func New() (*Config, error) {
	errLoad := godotenv.Load()
	if errLoad != nil {
		log.Fatal("error loading .env file")
	}

	viaCepEnable, err := strconv.ParseBool(GetEnvDefault("VIACEP_ENABLE", "false"))
	if err != nil {
		return nil, fmt.Errorf("error parse bool config enabled ViaCep - VIACEP_ENABLE: %w", err)
	}

	viaCep := ViaCepCfg{
		Enabled: viaCepEnable,
	}

	cepAbertoEnable, err := strconv.ParseBool(GetEnvDefault("CEPABERTO_ENABLE", "false"))
	if err != nil {
		return nil, fmt.Errorf("error parse bool config enabled CepAberto - CEPABERTO_ENABLE: %w", err)
	}

	cepAberto := CepAbertoCfg{
		Enabled: cepAbertoEnable,
	}

	if cepAberto.Enabled {
		cepAberto.Token = GetEnvDefault("CEPABERTO_TOKEN", "")
		if cepAberto.Token == "" {
			return nil, fmt.Errorf("token not configured for CepAberto: CEPABERTO_TOKEN")
		}
	}

	config := Config{
		ViaCepCfg:    viaCep,
		CepAbertoCfg: cepAberto,
	}

	return &config, nil
}

func GetEnvDefault(name, defaultValue string) string {
	cfg := os.Getenv(name)
	if cfg != "" {
		return cfg
	}
	return defaultValue
}
