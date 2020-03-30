package solr

import (
	"encoding/json"
	"path/filepath"
	"solr-snapshot-service/backup/blob"
)

// sits on the collection zk node
type CollectionData struct {
	ConfigName string `json:"configName"`
}

func getCollectionData(zkNode *blob.ZKNode, solrPath string, collectionName string) (*CollectionData, error) {
	collectionZkNode, err := zkNode.GetNode(filepath.Join(solrPath, "collections", collectionName))
	if err != nil {
		return nil, err
	}
	collectionData := &CollectionData{}
	err = json.Unmarshal([]byte(collectionZkNode.Data), collectionData)
	if err != nil {
		return nil, err
	}
	return collectionData, nil
}
