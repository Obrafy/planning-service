package apiclient

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/obrafy/planning/infrastructure/config"
)

const (
	INVALID_ENDPOINT_ERROR_MESSAGE = "endpoints must start with a leading slash (/)"
	POST_CONTENT_TYPE              = "application/json"
)

func NewAPIClient(config *config.APIClient) *APIClient {
	return &APIClient{
		Client: http.Client{
			Timeout: time.Second * time.Duration(config.TimeoutInSeconds),
		},
		BaseURI:          config.BaseURI,
		TimeoutInSeconds: uint16(config.TimeoutInSeconds),
	}
}

func (api *APIClient) GET(endpoint string) (*http.Response, error) {
	if endpoint[0:1] != "/" {
		return nil, fmt.Errorf(INVALID_ENDPOINT_ERROR_MESSAGE)
	}

	return api.Client.Get(api.BaseURI + endpoint)

}

func (api *APIClient) POST(endpoint string, body io.Reader) (*http.Response, error) {
	if endpoint[0:1] != "/" {
		return nil, fmt.Errorf(INVALID_ENDPOINT_ERROR_MESSAGE)
	}

	return api.Client.Post(api.BaseURI+endpoint, POST_CONTENT_TYPE, body)

}
