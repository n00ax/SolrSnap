package solr

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"solr-snapshot-service/backup/blob"
)

func CustomizeBlob(img *blob.Image, newNodeName string, newNodePort int) error {
	solrPath := img.Header.SolrZKPath
	// First delete custom paths, overseer_elect and live_nodes
	err := img.DeleteZK(filepath.Join(solrPath, "overseer_elect"))
	if err != nil {
		return err
	}
	err = img.DeleteZK(filepath.Join(solrPath, "live_nodes"))
	if err != nil {
		return err
	}
	collections, err := img.ZKRoot.GetNode(filepath.Join(solrPath, "collections"))
	if err != nil {
		return err
	}
	for _, child := range collections.Children {
		// modify state.json
		stateNode, err := child.GetNode(filepath.Join(child.Path, "state.json"))
		if err != nil {
			return err
		}
		var stateJson StateJson
		err = json.Unmarshal([]byte(stateNode.Data), &stateJson)
		if err != nil {
			return err
		}
		for k := range stateJson {
			for s := range stateJson[k].Shards {
				for r := range stateJson[k].Shards[s].Replicas {
					replica := stateJson[k].Shards[s].Replicas[r]
					replica.BaseUrl = fmt.Sprintf("http://%s:%d/solr", newNodeName, newNodePort)
					replica.NodeName = fmt.Sprintf("%s:%d_solr", newNodeName, newNodePort)
				}
			}
		}
		data, err := json.Marshal(stateJson)
		if err != nil {
			return err
		}
		err = img.ModifyZK(stateNode.Path, string(data))
		if err != nil {
			return err
		}
		// delete leaders
		leadersNode, err := child.GetNode(filepath.Join(child.Path, "leaders"))
		if err != nil {
			return err
		}
		for _, shardLeader := range leadersNode.Children {
			err := img.DeleteZK(filepath.Join(shardLeader.Path, "leader"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
