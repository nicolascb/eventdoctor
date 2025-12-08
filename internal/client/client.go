package client

import (
	"bytes"
	"fmt"
	"net/http"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

const (
	configEndpoint = "/api/v1/config"
)

type Client struct {
	http *http.Client
}

func newHTTPClient() *http.Client {
	// TODO: customize the HTTP client
	return &http.Client{}
}

func NewClient() *Client {
	return &Client{
		http: newHTTPClient(),
	}
}

// UploadSpec envia o spec já validado para o servidor
func (c *Client) UploadSpec(serverURL string, spec eventdoctor.EventDoctorSpec) error {
	body, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal spec: %w", err)
	}

	endpoint := serverURL + path.Join(configEndpoint)
	res, err := c.http.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload spec: status %s", res.Status)
	}

	return nil
}
