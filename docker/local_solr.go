package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"time"
)

const solrImage = "docker.io/library/solr:8.4.1"
const zookeeperImage = "docker.io/library/zookeeper:3.5.6"
const zookeeperUIImage = "docker.io/juris/zkui:latest"
const networkName = "solrsnap-net"

const solrHostname = "solr"
const zookeeperHostname = "zookeeper"
const zookeeperUIHostname =  "zookeeper-ui"

var zookeeperContainerName = fmt.Sprintf("%s-solrsnap-container", zookeeperHostname)
var zookeeperUIContainerName = fmt.Sprintf("%s-solrsnap-container", zookeeperUIHostname)
var solrContainerName = fmt.Sprintf("%s-solrsnap-container", solrHostname)

func StartLocalSolr(hostname string, solrDataPath string, baseSolrZKPath string) error {
	cli, ctx, err := CreateInstance()
	fmt.Printf("Pulling solr, zookeeper and zk-ui images\n")
	if err != nil {
		return nil
	}
	_, err = cli.ImageList(ctx, types.ImageListOptions{})
	_, err = cli.ImagePull(ctx, solrImage, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = cli.ImagePull(ctx, zookeeperImage, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = cli.ImagePull(ctx, zookeeperUIImage, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Creating %s network\n", networkName)
	err = cli.NetworkRemove(ctx, networkName)
	if err != nil && client.IsErrNetworkNotFound(err) {
		return err
	}
	_, err = cli.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	fmt.Print("Creating and starting containers\n")
	// Create Containers
	// (1.) Solr
	err = cli.ContainerRemove(ctx, solrContainerName, types.ContainerRemoveOptions{})
	_, err = cli.ContainerCreate(ctx, &container.Config{
			Hostname: solrHostname,
			Image: solrImage,
			ExposedPorts: nat.PortSet{
				"2181" : {},
			},
			Cmd: []string{""},
		},nil, nil, solrContainerName)
	err = cli.ContainerStart(ctx, solrContainerName, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	time.Sleep(10000)
	//zkSetupCommand := fmt.Sprintf("run -p 2181:2181 -d --network=%s --name=zookeeper -it zookeeper", networkName)
	//zkUiSetupCommand := fmt.Sprintf("run --rm -p 9090:9090 -d --network=%s -e ZK_SERVER=zookeeper:2181 juris/zkui", networkName)
	//solrStartCommand := fmt.Sprintf("run --hostname=%s -v (pwd):%s -d --network=%s -p 8986:8983 solr bash -c 'solr start -f -z zookeeper:2181%s", hostname, solrDataPath, networkName, baseSolrZKPath)
	//err = exec.Command("docker", zkSetupCommand).Run()
	//if err != nil {
	//	return fmt.Errorf("error creating zookeeper instance error, %s", err.Error())
	//}
	//err = exec.Command("docker", zkUiSetupCommand).Run()
	//if err != nil {
	//	return fmt.Errorf("error creating zkui instance error, %s", err.Error())
	//}
	//err = exec.Command("docker", solrStartCommand).Run()
	//if err != nil {
	//	return fmt.Errorf("error creating solr instance error, %s", err.Error())
	//}
	return nil
}
