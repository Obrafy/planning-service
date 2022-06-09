package apiclient

import (
	"net/http"
)

type APIClient struct {
	Client           http.Client
	BaseURI          string
	TimeoutInSeconds uint16
}
