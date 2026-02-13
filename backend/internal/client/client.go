package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

const (
	configEndpoint   = "/v1/config"
	topicsEndpoint   = "/v1/topics"
	servicesEndpoint = "/v1/services"
)

type Client struct {
	http *http.Client
}

func newHTTPClient() *http.Client {
	return &http.Client{}
}

func NewClient() *Client {
	return &Client{
		http: newHTTPClient(),
	}
}

// UploadSpec sends the validated spec to the server.
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

// GetTopicView fetches the topic detail view from the server.
func (c *Client) GetTopicView(serverURL, topicName string) (*response.TopicView, error) {
	endpoint := serverURL + topicsEndpoint + "/" + topicName
	res, err := c.http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %s", res.Status)
	}

	var view response.TopicView
	if err := json.NewDecoder(res.Body).Decode(&view); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &view, nil
}

// GetServiceView fetches the service detail view from the server.
func (c *Client) GetServiceView(serverURL, serviceName string) (*response.ServiceView, error) {
	endpoint := serverURL + servicesEndpoint + "/" + serviceName
	res, err := c.http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %s", res.Status)
	}

	var view response.ServiceView
	if err := json.NewDecoder(res.Body).Decode(&view); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &view, nil
}
