//go:build integration
package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/ybakhan/tax_calculator/taxclient"
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

	taxClient := taxclient.InitializeTaxClient(os.Getenv("INTERVIEW_SERVER"), httpClient)

	tests := map[string]struct {
		Year string
		Response taxclient.GetTaxBracketsResponse
		TaxBrackets int
	}{
		"tax bracket not found": {
			Year: "2018",
			Response: taxclient.NotFound,
		},
		"get brackets": {
			Year: "2022",
		 	Response: taxclient.Found,
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