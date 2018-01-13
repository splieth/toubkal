package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	auth "github.com/splieth/go-http-digest-auth-client"
)

func (c *Client) request(method, url, body string, expectedStatusCode int) ([]byte, error) {
	req := auth.NewRequest(c.UserName, c.APIKey, method, url, body)

	headers := map[string]string{"Content-Type": "application/json"}
	req.Header = headers

	resp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf("Got status %v instead of %v: %s", resp.StatusCode, expectedStatusCode, string(result))
	}

	return result, nil
}

func loadFixture(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join("fixtures", name))
}
