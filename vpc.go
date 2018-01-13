package main

import (
	"encoding/json"
	"net/http"
)

type ContainerOutput struct {
	Id             string `json"id,omitempty"`
	ProviderName   string `json"providerName,omitempty"`
	AtlasCidrBlock string `json"atlasCidrBlock,omitempty"`
	RegionName     string `json"regionName,omitempty"`
	VpcId          string `json"vpcId,omitempty"`
	IsProvisioned  bool   `json"isProvisioned,omitempty"`
}

type ContainerInput struct {
	AtlasCidrBlock string `json"atlasCidrBlock,omitempty"`
	ProviderName   string `json"providerName,omitempty"`
	RegionName     string `json"regionName,omitempty"`
}

func (c *Client) GetContainers() ([]ContainerOutput, error) {
	var containers []ContainerOutput
	var x map[string]*json.RawMessage

	response, err := c.request("GET", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/containers", "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &x)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*x["results"], &containers)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (c *Client) GetContainer(containerId string) (*ContainerOutput, error) {
	var container ContainerOutput

	response, err := c.request("GET", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/containers/"+containerId, "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &container)
	if err != nil {
		return nil, err
	}

	return &container, nil
}

func (c *Client) CreateContainer(container *ContainerInput) (*ContainerOutput, error) {
	var cnt ContainerOutput

	data, err := json.Marshal(container)
	if err != nil {
		return nil, err
	}

	response, err := c.request("POST", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/containers", string(data[:]), http.StatusCreated)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &cnt)
	if err != nil {
		return nil, err
	}

	return &cnt, nil
}
