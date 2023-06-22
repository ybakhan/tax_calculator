package taxclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitializeTaxClient(t *testing.T) {
	httpClient := retryablehttp.NewClient()
	taxClient := InitializeTaxClient("http://interview-test-server:5000", httpClient)
	assert.NotNil(t, taxClient)
}

func TestGetBrackets(t *testing.T) {
	tests := map[string]struct {
		HTTPResponse *http.Response
		HTTPError error
		Response GetTaxBracketsResponse
		TaxBrackets []*TaxBracket
		ReturnsError bool
	}{
		"get brackets failed": {
			HTTPResponse: &http.Response{},
			HTTPError: errors.New("some error"),
			Response: Failed,
			ReturnsError: true,
		},
		"tax bracket not found": {
			HTTPResponse: &http.Response{
				StatusCode:http.StatusNotFound, 
				Body: io.NopCloser(strings.NewReader("")),
			},
			Response: NotFound,
		},
		"get brackets failed - server response not ok": {
			HTTPResponse: &http.Response{
				StatusCode:http.StatusInternalServerError, 
				Body: io.NopCloser(strings.NewReader("some error")),
			},
			Response: Failed,
		},
		"get brackets failed - invalid json response": {
			HTTPResponse: &http.Response{
				StatusCode:http.StatusOK, 
				Body: io.NopCloser(strings.NewReader("invalid json format")),
			},
			Response: Failed,
			ReturnsError: true,
		},
		"get brackets": {
			HTTPResponse: &http.Response{
				StatusCode:http.StatusOK, 
				Body: io.NopCloser(strings.NewReader("{\"tax_brackets\":[{\"max\":100392,\"min\":50197,\"rate\":0.205}]}")),
			},
			Response: Found,
			TaxBrackets: []*TaxBracket{ {Min:50197, Max:100392, Rate:0.205}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockHTTPClient := &mockRetryableHTTPClient{};
			mockHTTPClient.
				On("Do", mock.AnythingOfType("*retryablehttp.Request")).
				Return(test.HTTPResponse, test.HTTPError)

			client := InitializeTaxClient("http://interview-test-server:5000", mockHTTPClient)
			taxBrackets, response, err := client.GetBrackets(context.Background(), "2022")
			assert.Equal(t, test.TaxBrackets, taxBrackets)
			assert.Equal(t, test.Response, response)
			if (test.ReturnsError) {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			mockHTTPClient.AssertExpectations(t)
		})
	}
}

type mockRetryableHTTPClient struct {
	mock.Mock
}

func (m *mockRetryableHTTPClient) Do(req *retryablehttp.Request) (*http.Response, error) {
    args := m.Called(req)
    return args.Get(0).(*http.Response), args.Error(1)
}