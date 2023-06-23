package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ybakhan/tax_calculator/taxcalculator"
	"github.com/ybakhan/tax_calculator/taxclient"
)

func (s *taxServer) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/tax/{year}", makeHTTPHandleFunc(s.handleTaxes))

	log.Printf("Tax server running on port %s", s.ListenAddress)
	http.ListenAndServe(s.ListenAddress, router)
}

func (s *taxServer) handleTaxes(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "GET" {
		return s.handleGetTaxes(w, r)
	}

	return http.StatusBadRequest, fmt.Errorf("method not supported %s", r.Method)
}

func (s *taxServer) handleGetTaxes(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	year := vars["year"]
	if _, err := strconv.Atoi(year); err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid tax year %s", year)
	}

	salaryStr := r.FormValue("s")
	if salaryStr == "" {
		return http.StatusBadRequest, errors.New("salary missing in request")
	}

	salaryF, err := strconv.ParseFloat(salaryStr, 32)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid salary %s", salaryStr)
	}

	ctx := context.Background()
	brackets, response, err := s.TaxClient.GetBrackets(ctx, year)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if response == taxclient.Failed {
		return http.StatusInternalServerError, fmt.Errorf("get taxes failed year %s", year)
	}

	if response == taxclient.NotFound {
		return http.StatusNotFound, fmt.Errorf("tax year not found %s", year)
	}

	taxes := taxcalculator.Calculate(brackets, float32(salaryF))
	WriteJSON(w, http.StatusOK, taxes)
	return http.StatusOK, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f requestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status, err := f(w, r); err != nil {
			// handle error
			WriteJSON(w, status, &taxServerError{Error: err.Error()})
		}
	}
}
