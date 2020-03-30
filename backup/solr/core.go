package solr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"solr-snapshot-service/backup/blob"
)

const DefailtFilePerm = 0777

type CoreProperties struct {
	NumShards            int    `properties:"numShards"`
	CollectionConfigName string `properties:"collection.configName"`
	Name                 string `properties:"name"`
	ReplicaType          string `properties:"replicaType"`
	Shard                string `properties:"shard"`
	Collection           string `properties:"collection"`
	CoreNodeName         string `properties:"coreNodeName"`
}

func (coreProperties *CoreProperties) toProperties() string {
	return fmt.Sprintf("numShards=%d\ncollection.configName=%s\nname=%s\nreplicaType=%s\nshard=%s\ncollection="+
		"%s\ncoreNodeName=%s\n", coreProperties.NumShards, coreProperties.CollectionConfigName, coreProperties.Name,
		coreProperties.ReplicaType, coreProperties.Shard, coreProperties.Collection, coreProperties.CoreNodeName)
}
func CreateCores(blob *blob.Image, solrDataPath string) error {
	cores, err := getCores(blob)
	if err != nil {
		return err
	}
	for coreName, core := range cores {
		err := os.Mkdir(filepath.Join(solrDataPath, coreName), DefailtFilePerm)
		if err == os.ErrExist {
			err = os.RemoveAll(filepath.Join(solrDataPath, coreName))
			if err != nil {
				return err
			}
			err = os.Mkdir(filepath.Join(solrDataPath, coreName), DefailtFilePerm)
			if err != nil {
				return err
			}
		}
		err = ioutil.WriteFile(filepath.Join(solrDataPath, coreName, "core.properties"), []byte(core.toProperties()), DefailtFilePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
func getCores(blob *blob.Image) (map[string]CoreProperties, error) {
	stateJsons, err := stateJsonsFromBlob(blob.Header.SolrZKPath, blob)
	if err != nil {
		return nil, err
	}
	cores := make(map[string]CoreProperties)
	for _, stateJson := range stateJsons {
		for collectionName, collection := range stateJson {
			for shardName, shard := range collection.Shards {
				for replicaCore, replica := range shard.Replicas {
					collectionData, err := getCollectionData(blob.ZKRoot, blob.Header.SolrZKPath, collectionName)
					if err != nil {
						return nil, err
					}
					coreProperties := CoreProperties{
						NumShards:            len(collection.Shards),
						CollectionConfigName: collectionData.ConfigName,
						Name:                 replica.Core,
						ReplicaType:          replica.Type,
						Shard:                shardName,
						Collection:           collectionName,
						CoreNodeName:         replicaCore,
					}
					cores[coreProperties.Name] = coreProperties
				}
			}
		}
	}
	return cores, nil
}
