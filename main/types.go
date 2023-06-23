package main

import (
	"net/http"

	"github.com/ybakhan/tax_calculator/taxclient"
)

type taxServer struct {
	ListenAddress string
	TaxClient     taxclient.TaxClient
}

type taxServerError struct {
	Error string `json:"error"`
}

type requestHandler func(http.ResponseWriter, *http.Request) (int, error)
