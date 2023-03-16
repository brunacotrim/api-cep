package cepaberto

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"api-cep/config"
	"api-cep/models"
)

type CepAberto struct {
	config *config.Config
	name   string
	token  string
}

type DataCepAberto struct {
	Cep          string `json:"cep"`
	Address      string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         City   `json:"cidade"`
	State        State  `json:"estado"`
}

type City struct {
	City string `json:"nome"`
}

type State struct {
	State string `json:"sigla"`
}

func New(cfg *config.Config) *CepAberto {
	return &CepAberto{
		config: cfg,
		name:   "cepaberto",
		token:  cfg.CepAbertoCfg.Token,
	}
}

func (c *CepAberto) GetName() string {
	return c.name
}

func (c *CepAberto) GetDataCep(cep string) (dataCep models.DataCep, err error) {
	url := fmt.Sprintf("https://www.cepaberto.com/api/v3/cep?cep=" + cep)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dataCep, fmt.Errorf("error preparing request CepAberto: %v", err)
	}
	req.Header.Add("Authorization", "Token token="+c.token)

	resp, err := client.Do(req)
	if err != nil {
		return dataCep, fmt.Errorf("error making request CepAberto: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dataCep, fmt.Errorf("error read body response CepAberto: %v", err)
	}

	defer resp.Body.Close()

	var data DataCepAberto
	err = json.Unmarshal(body, &data)
	if err != nil {
		return dataCep, fmt.Errorf("error parse response CepAberto: %v", err)
	}

	dataCep = c.ConvertData(data)
	if err != nil {
		return dataCep, fmt.Errorf("error convert data CepAberto: %v", err)
	}

	return dataCep, nil
}
