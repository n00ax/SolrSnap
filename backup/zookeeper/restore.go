package zookeeper

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"solr-snapshot-service/backup/blob"
)

func Restore(readable blob.Readable, connnect [] string, source blob.Source) error {
	conn, _, err := zk.Connect(connnect, connectDuration)
	if err != nil {
		conn.Close()
		return fmt.Errorf("zookeeper connection failed error %s", err.Error())
	}
	err = restoreNode(conn, readable)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to restore node error %s", err.Error())
	}
	conn.Close()
	return nil
}

func restoreNode(conn *zk.Conn, readable blob.Readable) error {
	_, err := conn.Create(readable.GetPath(), []byte(readable.GetData()), 0, zk.WorldACL(zk.PermAll))
	if err == zk.ErrNodeExists {
		conn.Delete(readable.GetPath(), 0)
		_, err = conn.Create(readable.GetPath(), []byte(readable.GetData()), 0, zk.WorldACL(zk.PermAll))
	} else if err != nil {
		return fmt.Errorf("ZK node creation failed error %s", err.Error())
	}
	for _, child := range readable.GetChildren() {
		restoreNode(conn, child)
	}
	return nil
}
