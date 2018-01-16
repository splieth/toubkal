package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/h2non/gock"
)

func TestClient_GetContainersListsContainers(t *testing.T) {
	defer gock.Off()

	client := givenClient()

	data, _ := loadFixture("get_containers_response.json")
	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups/some-group/containers").
		Reply(http.StatusOK).
		JSON(data)

	containers, _ := client.GetContainers()

	assert.Equal(t, 1, len(containers))
	assert.Equal(t, "US_EAST_1", containers[0].RegionName)
}

func TestClient_GetContainerFetchesContainer(t *testing.T) {
	defer gock.Off()

	client := givenClient()

	data, _ := loadFixture("get_container_response.json")
	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups/some-group/containers/some-container").
		Reply(http.StatusOK).
		JSON(data)

	container, _ := client.GetContainer("some-container")

	assert.Equal(t, "awesome-vpc", container.VpcId)
}

func TestClient_CreateContainerCreatesNewContainer(t *testing.T) {
	defer gock.Off()

	var containerToBuild ContainerInput
	client := givenClient()

	dataIn, _ := loadFixture("create_container_input.json")
	json.Unmarshal(dataIn, &containerToBuild)

	dataOut, _ := loadFixture("get_container_response.json")
	gock.New(client.BaseURL).
		Post("/api/atlas/v1.0/groups/some-group/containers").
		Reply(http.StatusCreated).
		JSON(dataOut)

	response, _ := client.CreateContainer(&containerToBuild)

	assert.Equal(t, response.AtlasCidrBlock, containerToBuild.AtlasCidrBlock)
}

func TestClient_CreateContainerFailsToCreatesWhenCidrIsInvalid(t *testing.T) {
	defer gock.Off()

	client := givenClient()

	gock.New(client.BaseURL).
		Post("/api/atlas/v1.0/groups/some-group/containers").
		Reply(http.StatusBadRequest)

	_, err := client.CreateContainer(nil)

	assert.NotNil(t, err)
}

func TestClient_UpdateContainerUpdatesExisting(t *testing.T) {
	defer gock.Off()

	var containerUpdate ContainerInput
	client := givenClient()

	originalContainer := givenContainer(&client)

	dataIn, _ := loadFixture("update_container_input.json")
	json.Unmarshal(dataIn, &containerUpdate)

	dataOut, _ := loadFixture("update_container_response.json")
	gock.New(client.BaseURL).
		Patch("/api/atlas/v1.0/groups/some-group/containers/some-container").
		Reply(http.StatusOK).
		JSON(dataOut)

	result, _ := client.UpdateContainer("some-container", &containerUpdate)

	assert.NotEqual(t, originalContainer.AtlasCidrBlock, result.AtlasCidrBlock)
	assert.Equal(t, result.AtlasCidrBlock, containerUpdate.AtlasCidrBlock)
}

func TestClient_UpdateContainerChangesNoting(t *testing.T) {
	defer gock.Off()

	var containerUpdate ContainerInput
	client := givenClient()

	originalContainer := givenContainer(&client)

	dataIn, _ := loadFixture("update_container_input.json")
	json.Unmarshal(dataIn, &containerUpdate)

	dataOut, _ := loadFixture("update_container_response_dummy.json")
	gock.New(client.BaseURL).
		Patch("/api/atlas/v1.0/groups/some-group/containers/some-container").
		Reply(http.StatusOK).
		JSON(dataOut)

	result, _ := client.UpdateContainer("some-container", &containerUpdate)

	assert.Equal(t, originalContainer.AtlasCidrBlock, result.AtlasCidrBlock)
}

func TestClient_UpdateContainerFailsOnIvalidCidr(t *testing.T) {
	defer gock.Off()

	client := givenClient()

	gock.New(client.BaseURL).
		Patch("/api/atlas/v1.0/groups/some-group/containers/some-container").
		Reply(http.StatusBadRequest)

	_, err := client.UpdateContainer("some-container", nil)

	assert.NotNil(t, err)
}

func givenClient() Client {
	return Client{
		BaseURL:  "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey:   "some-key",
		GroupId:  "some-group",
	}
}

func givenContainer(client *Client) (*ContainerOutput) {
	var containerCreate ContainerInput
	dataInCreate, _ := loadFixture("create_container_input.json")
	json.Unmarshal(dataInCreate, &containerCreate)
	dataOutCreate, _ := loadFixture("create_container_response.json")
	gock.New(client.BaseURL).
		Post("/api/atlas/v1.0/groups/some-group/containers").
		Reply(http.StatusCreated).
		JSON(dataOutCreate)
	originalContainer, _ := client.CreateContainer(&containerCreate)
	return originalContainer
}
