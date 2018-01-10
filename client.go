package main

import (
	"net/http"
)

type Client struct {
	BaseURL  string
	UserName string
	APIKey   string
	GroupId  string

	httpClient *http.Client
}

func NewClient(url, username, password string) *Client {
	return &Client{
		BaseURL:  url,
		UserName: username,
		APIKey:   password,
	}
}
