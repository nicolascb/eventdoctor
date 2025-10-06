package client

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) UploadSpec(filePath string) error {
	// Implement the logic to upload the OpenAPI spec file to the server
	return nil
}
