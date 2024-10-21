package zosmf

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Host       string
	HttpClient *http.Client
	Username   string
	Password   string
}

func NewClient(host string, username string, password string, insecure bool) (*Client, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	c := Client{
		HttpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
		},
		Host:     host,
		Username: username,
		Password: password,
	}

	return &c, nil

}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {

	req.Header.Set("X-CSRF-ZOSMF-HEADER", "dummy")

	req.SetBasicAuth(c.Username, c.Password)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, nil
}

func (c *Client) GetDataset(name string) (*Dataset, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zosmf/restfiles/ds/%s", c.Host, name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Content-Type", "text/plain")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	content := Dataset{
		Content: string(body),
	}

	return &content, nil
}

func (c *Client) CreateDataset(name string, dataset CreateDataset, content Dataset) error {
	rb, err := json.Marshal(dataset)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/zosmf/restfiles/ds/%s", c.Host, name), bytes.NewBuffer(rb))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	err = c.AddData(name, content)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddData(name string, content Dataset) error {

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/zosmf/restfiles/ds/%s", c.Host, name), bytes.NewBuffer([]byte(content.Content)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Content-Type", "text/plain")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}
