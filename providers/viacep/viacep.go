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
		name:   "viacep",
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
		return dataCep, fmt.Errorf("error preparing request ViaCep: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return dataCep, fmt.Errorf("error making request ViaCep: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dataCep, fmt.Errorf("error read body response ViaCep: %v", err)
	}

	defer resp.Body.Close()

	var data DataViaCep
	err = json.Unmarshal(body, &data)
	if err != nil {
		return dataCep, fmt.Errorf("error parse response ViaCep: %v", err)
	}

	dataCep = v.ConvertData(data)
	if err != nil {
		return dataCep, fmt.Errorf("error convert data ViaCep: %v", err)
	}

	return dataCep, nil
}
