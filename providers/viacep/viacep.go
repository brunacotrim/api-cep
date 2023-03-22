package viacep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"api-cep/config"
	"api-cep/models"
)

type ViaCep struct {
	config *config.Config
	name   string
}

type DataViaCep struct {
	Cep          string `json:"cep"`
	Address      string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
}

func New(cfg *config.Config) *ViaCep {
	return &ViaCep{
		config: cfg,
		name:   "ViaCep",
	}
}

func (v *ViaCep) GetName() string {
	return v.name
}

func (v *ViaCep) GetDataCep(cep string) (dataCep models.DataCep, err error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/" + cep + "/json")

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dataCep, fmt.Errorf("error preparing request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return dataCep, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != 200 {
		return dataCep, fmt.Errorf("request return status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dataCep, fmt.Errorf("error read body response: %v", err)
	}

	defer resp.Body.Close()

	var data DataViaCep
	err = json.Unmarshal(body, &data)
	if err != nil {
		return dataCep, fmt.Errorf("error parse response: %v", err)
	}

	dataCep = v.ConvertData(data)
	if err != nil {
		return dataCep, fmt.Errorf("error convert data: %v", err)
	}

	return dataCep, nil
}
