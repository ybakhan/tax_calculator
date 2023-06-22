//go:build integration
// +build integration

package taxclient

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestGetBracketsIT(t *testing.T) {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = 1 * time.Second
    httpClient.RetryWaitMax = 5 * time.Second
    httpClient.RetryMax = 3
	//client.CheckRetry = retryablehttp.DefaultRetryPolicy
	httpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp.StatusCode == http.StatusInternalServerError {
			return true, nil
		}
    	return false, nil
	}

	taxClient := InitializeTaxClient("http://localhost:5000", httpClient)

	tests := map[string]struct {
		Year string
		Response GetTaxBracketsResponse
		TaxBrackets int
	}{
		"tax bracket not found": {
			Year: "2018",
			Response: NotFound,
		},
		"get brackets": {
			Year: "2022",
		 	Response: Found,
		 	TaxBrackets: 5,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			taxBrackets, response, err := taxClient.GetBrackets(context.Background(), test.Year)
			assert.Equal(t, test.TaxBrackets, len(taxBrackets))
			assert.Equal(t, test.Response, response)
			assert.Nil(t, err)
		})
	}
}