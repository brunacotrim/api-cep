package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"api-cep/config"
	"api-cep/models"
	"api-cep/providers/cepaberto"
	"api-cep/providers/viacep"
)

type ServiceCep interface {
	GetDataCep(cep string) (models.DataCep, error)
	GetName() string
}

type Service struct {
	config    *config.Config
	providers map[string]ServiceCep
}

func New(cfg *config.Config) (*Service, error) {
	service := Service{
		config:    cfg,
		providers: make(map[string]ServiceCep),
	}

	if cfg.ViaCepCfg.Enabled {
		service.providers["viacep"] = viacep.New(cfg)
	}

	if cfg.CepAbertoCfg.Enabled {
		service.providers["cepaberto"] = cepaberto.New(cfg)
	}

	return &service, nil
}

func (s *Service) StartServer() {
	http.HandleFunc("/cep/", s.GetDataCepHandler)
	http.ListenAndServe(":8080", nil)
}

func (s *Service) GetDataCepHandler(w http.ResponseWriter, r *http.Request) {

	cep := r.URL.Path[len("/cep/"):]
	var statusCode int
	if cep == "" {
		log.Println("cep not informed")
		statusCode = http.StatusBadRequest
		return
	}
	cep, err := s.CepNumbers(cep)
	if err != nil {
		log.Printf("invalid cep %s: %v\n", cep, err)
		statusCode = http.StatusBadRequest
		return
	}

	resp := s.GetCep(cep)

	w.Header().Set("Content-Type", "application/json")
	if resp.Cep == "" {
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = http.StatusOK
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) GetCep(cep string) (resp models.DataCep) {

	ch := make(chan bool)
	var err error

	for _, p := range s.providers {
		go func(p ServiceCep) {
			resp, err = p.GetDataCep(cep)
			if err != nil {
				log.Printf("error GetDataCep by %s: %v", p.GetName(), err)
			}

			if resp.Cep != "" {
				resp.Cep, err = s.CepFormat(cep)
				if err != nil {
					log.Printf("error CepFormat by %s: %v", p.GetName(), err)
				} else {
					ch <- true
				}
			}
		}(p)
	}

	select {
	case <-ch:
		return resp
	case <-time.After(time.Second * 5):
		return resp
	}
}
