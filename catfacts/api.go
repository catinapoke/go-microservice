package catfacts

import (
	"context"
	"encoding/json"
	"net/http"
)

type ApiServer struct {
	service Service
}

func NewApiServer(service Service) *ApiServer {
	return &ApiServer{
		service: service,
	}
}

func (s *ApiServer) Start(listenAddr string) error {
	http.HandleFunc("/", s.handleGetCatFact)
	return http.ListenAndServe(listenAddr, nil)
}

func (s *ApiServer) handleGetCatFact(w http.ResponseWriter, r *http.Request) {
	fact, err := s.service.GetCatFact(context.Background())
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, fact)
}

func writeJson(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
