package config

import (
	"fmt"
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
		return nil, fmt.Errorf("error loading .env file: %w", errLoad)
	}

	viaCepEnabled, err := strconv.ParseBool(GetEnvDefault("VIACEP_ENABLE", "false"))
	if err != nil {
		return nil, fmt.Errorf("error parse bool config - VIACEP_ENABLE: %w", err)
	}

	viaCep := ViaCepCfg{
		Enabled: viaCepEnabled,
	}

	cepAbertoEnabled, err := strconv.ParseBool(GetEnvDefault("CEPABERTO_ENABLE", "false"))
	if err != nil {
		return nil, fmt.Errorf("error parse bool config - CEPABERTO_ENABLE: %w", err)
	}

	cepAberto := CepAbertoCfg{
		Enabled: cepAbertoEnabled,
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
