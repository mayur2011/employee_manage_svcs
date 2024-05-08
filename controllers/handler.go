package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// response used to send HTTP responses
type response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

// Generic handler for writing response header and body for all handler functions
func ResponseHandler(h func(http.ResponseWriter, *http.Request) (interface{}, int, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r) // execute application handler
		var errorMsg string
		if err != nil {
			errorMsg = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if data != nil {
			// Send JSON response back to the client application
			err = json.NewEncoder(w).Encode(response{Data: data, Error: errorMsg})
			if err != nil {
				log.Printf("Error from Handler: %s\n", err.Error())
			}
		} else {
			err = json.NewEncoder(w).Encode(response{Data: nil, Error: errorMsg})
			if err != nil {
				fmt.Printf("Error from Handler: %s\n", err.Error())
			}
		}

	})
}
