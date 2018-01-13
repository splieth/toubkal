package main

import (
	"encoding/json"
	"net/http"
)

type ProjectOutput struct {
	ClusterCount int    `json:"clusterCount,omitempty"`
	Created      string `json:"created,omitempty"`
	Id           string `json:"id,omitempty"`
	Links		 []Link `json:"links,omitempty"`
	Name         string `json:"name,omitempty"`
	OrgId        string `json:"orgId,omitempty"`
}

type ProjectInput struct {
	Name         string `json:"name,omitempty"`
	OrgId        string `json:"orgId,omitempty"`
}

type Link struct {
	Href string `json:"href,omitempty"`
	Rel  string `json:"rel,omitempty"`
}

func (c *Client) GetProjects() ([]ProjectOutput, error) {
	var projects []ProjectOutput
	var x map[string]*json.RawMessage

	response, err := c.request("GET", c.BaseURL+"/api/atlas/v1.0/groups", "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &x)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*x["results"], &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *Client) GetProject(groupId string) (*ProjectOutput, error) {
	var project ProjectOutput

	response, err := c.request("GET", c.BaseURL+"/api/atlas/v1.0/groups/"+groupId, "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) GetProjectByName(name string) (*ProjectOutput, error) {
	var project ProjectOutput

	response, err := c.request("GET", c.BaseURL+"/api/atlas/v1.0/groups/byName/"+name, "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) CreateProject(project *ProjectInput) (*ProjectOutput, error) {
	var proj ProjectOutput

	data, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	response, err := c.request("POST", c.BaseURL+"/api/atlas/v1.0/groups", string(data[:]), http.StatusCreated)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &proj)
	if err != nil {
		return nil, err
	}

	return &proj, nil
}
