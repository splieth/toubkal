package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/h2non/gock"
)

func TestClient_GetProjectsListsAllProjects(t *testing.T) {
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

	data, _ := loadFixture("get_projects_response.json")
	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups").
		Reply(http.StatusOK).
		JSON(data)

	projects, _ := client.GetProjects()

	assert.Equal(t, 2, len(projects))
}

func TestClient_GetProjectGetsProjectByGroupId(t *testing.T) {
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

	data, _ := loadFixture("get_project_response.json")
	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups/some-group").
		Reply(http.StatusOK).
		JSON(data)

	project, _ := client.GetProject("some-group")

	assert.Equal(t, "5a0a1e7e0f2912c554080adc", project.OrgId)
}

func TestClient_GetProjectFailsWhenGroupIdDoesNotExist(t *testing.T) {
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-invalid-group",
	}

	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups/some-invalid-group").
		Reply(http.StatusNotFound)

	_, err := client.GetProject("some-invalid-group")

	assert.NotNil(t, err)
}

func TestClient_GetProjectByNameGetsProjectByGroupId(t *testing.T) {
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

	data, _ := loadFixture("get_project_response.json")
	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups/byName/some-group").
		Reply(http.StatusOK).
		JSON(data)

	project, _ := client.GetProjectByName("some-group")

	assert.Equal(t, "5a0a1e7e0f2912c554080adc", project.OrgId)
}

func TestClient_GetProjectByNameFailsWhenGroupIdDoesNotExist(t *testing.T) {
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-invalid-group",
	}

	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups/byName/some-invalid-group").
		Reply(http.StatusNotFound)

	_, err := client.GetProjectByName("some-invalid-group")

	assert.NotNil(t, err)
}

func TestClient_CreateProjectCreatesNewProject(t *testing.T) {
	var project ProjectInput
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

	dataIn, _ := loadFixture("cluster_create_input.json")
	json.Unmarshal(dataIn, &project)

	dataOut, _ := loadFixture("cluster_create_response.json")
	gock.New(client.BaseURL).
		Post("/api/atlas/v1.0/groups").
		Reply(http.StatusCreated).
		JSON(dataOut)

	response, _ := client.CreateProject(&project)

	assert.Equal(t, response.Name, project.Name)
	assert.Equal(t, response.OrgId, project.OrgId)
}

func TestClient_CreateProjectFailsWhenProjectExists(t *testing.T) {
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

	gock.New(client.BaseURL).
		Post("/api/atlas/v1.0/groups").
		Reply(http.StatusCreated)

	_, err := client.CreateProject(nil)

	assert.NotNil(t, err)
}
