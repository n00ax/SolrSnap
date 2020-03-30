package automated

import (
	"fmt"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/fatih/color"
	"os"
	"solr-snapshot-service/backup/blob"
	"solr-snapshot-service/backup/solr"
	"solr-snapshot-service/docker"
)

const localSolrNodeName = "solr"
const localSolrPort = 8983

func CreateLocal(args []string) {
	blobPath := args[0]
	// load blob
	fmt.Printf("Loading blob %s\n", blobPath)
	img, err := blob.Read(blobPath)
	if err != nil{
		color.Red("Error reading blob!, error %s", err.Error())
		os.Exit(-1)
	}
	// customize blob
	fmt.Printf("Customizing blob %s to %s\n", blobPath, localSolrNodeName)
	err = solr.CustomizeBlob(img, localSolrNodeName, localSolrPort)
	if err != nil {
		color.Red("Error customizing blob error, %s", err.Error())
		os.Exit(-1)
	}
	// create /data files
	fmt.Printf("Creating instance local and core data in temporary directory..\n")
	dir, err := ioutils.TempDir("/tmp", "solrsnap-imgd")
	if err != nil {
		color.Red("Error creating temporary directory for customized blob error, %s", err.Error())
		os.Exit(-1)
	}
	err = solr.CreateCores(img, dir)
	if err != nil {
		color.Red("Error creating core data, error %s", err.Error())
		os.Exit(-1)
	}
	// create patched configs
	fmt.Printf("Creating patched solr configs in data directory..\n")
	err = solr.WriteDefaultBaseConfig(dir)
	if err != nil{
		color.Red("Error creating patched solr configs, error %s\n", err.Error())
		os.Exit(-1)
	}
	// startup Solr and Zookeeper
	fmt.Printf("Starting Solr instance with temporary data directory\n")
	err = docker.StartLocalSolr(localSolrNodeName ,dir, img.Header.SolrZKPath)
	if err != nil{
		color.Red("Error creating Solr instance, error %s", err.Error())
		os.Exit(-1)
	}
}
