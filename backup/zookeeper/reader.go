package zookeeper

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"path"
	"solr-snapshot-service/backup/blob"
)

const connectDuration = 90000000

func Read(connect [] string, source blob.Source, generateUser string, description string, solrZKPath string) (*blob.Image, error) {
	conn, _, err := zk.Connect(connect, connectDuration)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("zookeeper connection failed %s", err.Error())
	}
	zk.WithLogger(log.New())
	root, err := readNode("/", conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	conn.Close()
	return blob.GenerateFromNode(root, source, generateUser, description, solrZKPath)
}

func readNode(nodePath string, conn *zk.Conn) (*blob.ZKNode, error) {
	childZKNodes, _, err := conn.Children(nodePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get node children at nodePath %s error %s", nodePath, err.Error())
	}
	children := make([] *blob.ZKNode, len(childZKNodes))
	for i, zkChild := range childZKNodes {
		child, err := readNode(path.Join(nodePath, zkChild), conn)
		children[i] = child
		if err != nil {
			return nil, err
		}
	}
	data, _, err := conn.Get(nodePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get node data at nodePath %s error %s", nodePath, err.Error())
	}
	thisNode := &blob.ZKNode{
		Children:       children,
		Data:           string(data),
		Path:           nodePath,
		IsEmpty:        string(data) == "",
		SHA256Checksum: "",
	}
	thisNode.SHA256Checksum = thisNode.GenerateNodeHash()
	fmt.Printf("%s\n", nodePath)
	return thisNode, nil
}
