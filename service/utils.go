package service

import (
	"fmt"
	"regexp"
	"strings"
)

func (s *Service) CepNumbers(cep string) (string, error) {
	regexNumber := regexp.MustCompile(`\d`)
	cepNumberSlc := regexNumber.FindAllString(cep, -1)
	cepNumber := strings.Join(cepNumberSlc, "")
	if cepNumber == "" || len(cepNumber) != 8 {
		return "", fmt.Errorf("error invalid CEP: %s", cep)
	}
	return cepNumber, nil
}

func (s *Service) CepFormat(cep string) (string, error) {

	regexFormat := regexp.MustCompile(`(\d{2})(\d{3})(\d{3})`)
	cepFormat := regexFormat.ReplaceAllString(cep, "$1.$2-$3")

	return cepFormat, nil
}
