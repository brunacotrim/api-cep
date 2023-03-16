package cepaberto

import "api-cep/models"

func (c *CepAberto) ConvertData(dataCepAberto DataCepAberto) models.DataCep {
	dataCep := models.DataCep{
		Cep:          dataCepAberto.Cep,
		Address:      dataCepAberto.Address,
		Neighborhood: dataCepAberto.Neighborhood,
		City:         dataCepAberto.City.City,
		State:        dataCepAberto.State.State,
	}
	return dataCep
}
