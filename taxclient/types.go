package taxclient

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type GetTaxBracketsResponse int

const (
    Found GetTaxBracketsResponse = -(iota)
    NotFound
    Failed
)

type TaxClient interface {
	GetBrackets(context.Context, string) ([]*TaxBracket, GetTaxBracketsResponse, error)
}

type retryableHTTPClient interface {
    Do(req *retryablehttp.Request) (*http.Response, error)
}

type taxClient struct {
	taxBracketsURL string
	client retryableHTTPClient
}

type TaxBracket struct {
	Min  uint    `json:"min"`
	Max  uint    `json:"max"`
	Rate float32 `json:"rate"`
}

type taxBrackets struct {
	Data []*TaxBracket `json:"tax_brackets"`
}