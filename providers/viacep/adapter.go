package viacep

import "api-cep/models"

func (v *ViaCep) ConvertData(dataViaCep DataViaCep) models.DataCep {
	dataCep := models.DataCep{
		Cep:          dataViaCep.Cep,
		Address:      dataViaCep.Address,
		Neighborhood: dataViaCep.Neighborhood,
		City:         dataViaCep.City,
		State:        dataViaCep.State,
	}
	return dataCep
}
