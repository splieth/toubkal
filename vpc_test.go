package main
import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/h2non/gock"
	"encoding/json"
)

func TestClient_GetContainersListsContainers(t *testing.T) {
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

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

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

	data, _ := loadFixture("get_container_response.json")
	gock.New(client.BaseURL).
		Get("/api/atlas/v1.0/groups/some-group/containers/some-container").
		Reply(http.StatusOK).
		JSON(data)

	container, _ := client.GetContainer("some-container")

	assert.Equal(t, "awesome-vpc", container.VpcId)
}

func TestClient_CreateContainerCreatesNewContainer(t *testing.T) {
	var containerToBuild ContainerInput
	defer gock.Off()

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

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

	client := Client{
		BaseURL: "https://cloud.mongodb.com",
		UserName: "some-user",
		APIKey: "some-key",
		GroupId: "some-group",
	}

	gock.New(client.BaseURL).
		Post("/api/atlas/v1.0/groups/some-group/containers").
		Reply(http.StatusBadRequest)

	_, err := client.CreateContainer(nil)

	assert.NotNil(t, err)
}
