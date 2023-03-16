package service

import (
	"encoding/json"
	"log"
	"net/http"

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
	if cep == "" {
		log.Println("cep not informed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cep, err := s.CepNumbers(cep)
	if err != nil {
		log.Printf("invalid cep %s: %v\n", cep, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := models.DataCep{}

	for _, p := range s.providers {
		resp, err = p.GetDataCep(cep)
		if err != nil {
			log.Printf("error GetDataCep by %s: %v", p.GetName(), err)
			continue
		}

		if resp.Cep != "" {
			resp.Cep, err = s.CepFormat(resp.Cep)
			if err != nil {
				log.Printf("error CepFormat by %s: %v", p.GetName(), err)
				continue
			}
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Cep == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
