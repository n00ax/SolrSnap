package solr

import (
	"encoding/json"
	"path/filepath"
	"solr-snapshot-service/backup/blob"
)

type StateJson map[string]Collection

type Collection struct {
	PullReplicas      string            `json:"pullReplicas"`
	ReplicationFactor string            `json:"replicationFactor"`
	Shards            map[string]*Shard `json:"shards"`
	MaxShardsPerNode  string            `json:"maxShardsPerNode"`
	AutoAddReplicas   string            `json:"autoAddReplicas"`
	NrtReplicas       string            `json:"nrtReplicas"`
	TlogReplicas      string            `json:"tlogReplicas"`
	Router            *Router           `json:"router"`
}
type Shard struct {
	Range    string           `json:"range"`
	State    string           `json:"state"`
	Replicas map[string]*Core `json:"replicas"`
}
type Core struct {
	Core          string `json:"core"`
	BaseUrl       string `json:"base_url"`
	NodeName      string `json:"node_name"`
	State         string `json:"state"`
	Type          string `json:"type"`
	ForceSetState string `json:"force_set_state"`
	Leader        string `json:"leader"`
}
type Router struct {
	Name  string `json:"name"`
	Field string `json:"field"`
}

func stateJsonsFromBlob(solrPath string, blob *blob.Image) ([]StateJson, error) {
	collections, err := blob.ZKRoot.GetNode(filepath.Join(solrPath, "collections"))
	if err != nil {
		return nil, err
	}
	stateJsons := make([]StateJson, len(collections.Children))
	for i, child := range collections.Children {
		stateNode, err := child.GetNode(filepath.Join(child.Path, "state.json"))
		if err != nil {
			return nil, err
		}
		var stateJson StateJson
		err = json.Unmarshal([]byte(stateNode.Data), &stateJson)
		if err != nil {
			return nil, err
		}
		stateJsons[i] = stateJson
	}
	return stateJsons, nil
}
