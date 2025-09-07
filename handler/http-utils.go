package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeHttpResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error while writing response: %v", err)
	}
}
