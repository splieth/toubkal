package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/h2non/gock"
)

const (
	CREATE_ENDPOINT = "/api/atlas/v1.0/groups/some-group-id/clusters"
)

func createClient() Client {
	return Client{
		BaseURL:  "https://cloud.mongodb.com",
		UserName: "some-user@dev.null",
		APIKey:   "some-api-key",
	}
}

func TestClient_CreateClusterCreatesNewCluster(t *testing.T) {
	var clusterToBuild ClusterInput
	defer gock.Off()

	client := createClient()

	// Input
	dataIn, _ := loadFixture("cluster_create_input.json")
	json.Unmarshal(dataIn, &clusterToBuild)

	// Output
	dataOut, _ := loadFixture("cluster_create_response.json")
	gock.New(client.BaseURL).
		Post(CREATE_ENDPOINT).
		Reply(http.StatusCreated).
		JSON(dataOut)

	response, _ := client.CreateCluster(&clusterToBuild)

	assert.Equal(t, response.Name, clusterToBuild.Name)
	assert.Equal(t, response.NumShards, clusterToBuild.NumShards)
}

func TestClient_CreateClusterFailsWhenClusterExists(t *testing.T) {
	defer gock.Off()

	client := createClient()

	gock.New(client.BaseURL).
		Post(CREATE_ENDPOINT).
		Reply(http.StatusBadRequest)

	_, err := client.CreateCluster(nil)

	assert.NotNil(t, err)
}
