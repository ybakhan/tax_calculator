package taxclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

const taxBracketsResource = "/tax-calculator/tax-year/"

func InitializeTaxClient(baseURL string, client retryableHTTPClient) TaxClient {
	taxBracketsURL, err := url.JoinPath(baseURL, taxBracketsResource)
	if err != nil {
		err = fmt.Errorf("error intializing tax client: %w", err)
		panic(err)
	}

	return &taxClient{
		taxBracketsURL,
		client,
	}
}

func (tc *taxClient) GetBrackets(ctx context.Context, year string) ([]*TaxBracket, GetTaxBracketsResponse, error) {
	taxBracketsURL, err := url.JoinPath(tc.taxBracketsURL, year)
	if err != nil {
		return nil, Failed, err
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, taxBracketsURL, nil)
	if err != nil {
		return nil, Failed, err
	}

	resp, err := tc.client.Do(req)
	if err != nil {
		return nil, Failed, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFound, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, Failed, nil
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, Failed, err
	}

	var taxbrackets TaxBrackets
	err = json.Unmarshal(respBytes, &taxbrackets)
	if err != nil {
		return nil, Failed, err
	}

	return taxbrackets.Data, Found, nil
}
