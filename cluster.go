package main

import (
	"encoding/json"
	"net/http"
)

type ProviderSettings struct {
	ProviderName     string `json:"providerName"`
	RegionName       string `json:"regionName"`
	InstanceSizeName string `json:"instanceSizeName"`
	DiskIOPS         int    `json:"diskIOPS,omitempty"`
	EncryptEBSVolume bool   `json:"encryptEBSVolume,omitempty"`
}

type ClusterInput struct {
	Name                string           `json:"name"`
	MongoDBMajorVersion string           `json:"mongoDBMajorVersion"`
	NumShards           int              `json:"numShards,omitempty"`
	ProviderSettings    ProviderSettings `json:"providerSettings"`
	BackupEnabled       bool             `json:"backupEnabled"`
	Autoscaling         Autoscaling      `json:"autoScaling,omitempty"`
	Paused              bool             `json:"paused,omitempty"`
	ReplicationFactor   int8             `json:"replicationFactor,omitempty"`
	DiskSizeGB          float64          `json:"diskSizeGB,omitempty"`
	BiConnector         BiConnector      `json:"biConnector,omitempty"`
}

type ClusterOutput struct {
	Name                string           `json:"name,omitempty"`
	MongoDBMajorVersion string           `json:"mongoDBMajorVersion,omitempty"`
	NumShards           int              `json:"numShards,omitempty,omitempty"`
	ProviderSettings    ProviderSettings `json:"providerSettings,omitempty"`
	BackupEnabled       bool             `json:"backupEnabled,omitempty"`
	Autoscaling         Autoscaling      `json:"autoScaling,omitempty"`
	Paused              bool             `json:"paused,omitempty"`
	ReplicationFactor   int8             `json:"replicationFactor,omitempty"`
	DiskSizeGB          float64          `json:"diskSizeGB,omitempty"`
	BiConnector         BiConnector      `json:"biConnector,omitempty"`
	GroupId             string           `json:"groupId,omitempty"`
	Id                  string           `json:"id,omitempty"`
	MongoURI            string           `json:"mongoURI,omitempty"`
	MongoURIUpdated     string           `json:"mongoURIUpdated,omitempty"`
	MongoURIWithOptions string           `json:"mongoURIWithOptions,omitempty"`
	StateName           string           `json:"stateName,omitempty"`
}

type Autoscaling struct {
	DiskGBEnabled bool `json:"diskGBEnabled,omitempty"`
}

type BiConnector struct {
	Enabled        bool   `json:"enabled,omitempty"`
	ReadPreference string `json:"replicationFactor,omitempty"`
}

func (c *Client) CreateCluster(cluster *ClusterInput) (*ClusterOutput, error) {
	var cl ClusterOutput

	data, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}

	response, err := c.request("POST", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/clusters", string(data[:]), http.StatusCreated)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &cl)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

func (c *Client) GetCluster(clusterName string) (*ClusterOutput, error) {
	var cluster ClusterOutput

	response, err := c.request("GET", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/clusters/"+clusterName, "", http.StatusOK)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(response, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (c *Client) GetClusters() ([]ClusterOutput, error) {
	var clusters []ClusterOutput
	var x map[string]*json.RawMessage

	response, err := c.request("GET", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/clusters", "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &x)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*x["results"], &clusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func (c *Client) UpdateCluster(cluster *ClusterInput) (*ClusterOutput, error) {
	var cl ClusterOutput
	data, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}

	response, err := c.request("PATCH", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/clusters/"+cluster.Name, string(data[:]), http.StatusOK)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &cl)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

func (c *Client) DeleteCluster(clusterName string) error {
	_, err := c.request("DELETE", c.BaseURL+"/api/atlas/v1.0/groups/"+c.GroupId+"/clusters/"+clusterName, "", http.StatusAccepted)
	if err != nil {
		return err
	}

	return nil
}
