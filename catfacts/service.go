package catfacts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Service interface {
	GetCatFact(context.Context) (*CatFact, error)
}

type CatFactService struct {
	url string
}

func NewCatFactService(url string) Service {
	return &CatFactService{
		url: url,
	}
}

func RunCatFactsServer() {
	service := NewCatFactService("https://catfact.ninja/fact")
	service = NewLoggingService(service)

	apiServer := NewApiServer(service)

	fmt.Println("Starting Cat fact service at port 3001")
	log.Fatal(apiServer.Start(":3001"))
}

func (s *CatFactService) GetCatFact(ctx context.Context) (*CatFact, error) {
	resp, err := http.Get(s.url)

	if err != nil {
		return nil, err
	}

	fact := &CatFact{}
	if err := json.NewDecoder(resp.Body).Decode(fact); err != nil {
		return nil, err
	}

	return fact, nil
}
