package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/ybakhan/tax_calculator/taxclient"
)

func main() {
	config := readConfig()
	fmt.Printf("Tax calculator config: %+v\n", config)

	httpClient := initializeHTTPClient(config)
	taxClient := taxclient.InitializeTaxClient(config.InterviewServer.BaseURL, httpClient)

	listenAddress := fmt.Sprintf(":%d", config.Port)
	server := &taxServer{listenAddress, taxClient}
	server.Start()
}

func initializeHTTPClient(config *Config) *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = time.Duration(config.HTTPClient.Retry.Wait.MinSeconds) * time.Second
	httpClient.RetryWaitMax = time.Duration(config.HTTPClient.Retry.Wait.MaxSeconds) * time.Second
	httpClient.RetryMax = config.HTTPClient.Retry.Max
	//client.CheckRetry = retryablehttp.DefaultRetryPolicy

	httpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp.StatusCode == http.StatusInternalServerError {
			return true, nil
		}
		return false, nil
	}
	return httpClient
}
