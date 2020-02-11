package blob

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Generate(dirPath string, source Source, generateUser string, description string) (*Image, error){
	root, err := generateNode(dirPath, "/", true)
	if err != nil {
		return nil, fmt.Errorf("failed to generate IMG image error %s", err.Error())
	}
	img := Image{
		Header:    Header{
			Magic:             Magic,
			SHA256Checksum:    root.SHA256Checksum,
			Version:           Version,
			Source:            source,
			SourceApplication: SolrSnapshotService,
			GenerateUser:      generateUser,
			Description:       description,
			TimeGenerated:     time.Now(),
		},
		ZKRoot:    root,
		IndexData: &IndexDataRoot{},
	}
	validate := img.ZKRoot.validateNode()
	if validate != nil{
		return nil, fmt.Errorf("node Merkle tree validation failed error %s", validate.Error())
	}
	return &img, err
}

func generateNode(dirPath string, relativePath string, isRoot bool) (*ZKNode, error){
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	curName := filepath.Base(dirPath)
	var data []uint8
	data, err = ioutil.ReadFile(filepath.Join(dirPath, curName))
	if err == err.(*os.PathError){
		data = []uint8("")
	} else if err != nil {
		return nil, err
	}
	children := make([]*ZKNode, 0)
	for _, file := range files {
		if file.IsDir() && !strings.HasPrefix(file.Name(), "."){
			node, err := generateNode(filepath.Join(dirPath, file.Name()), filepath.Join(relativePath, file.Name()), false)
			if err != nil {
				return nil, err
			}
			children = append(children, node)
		}
	}
	thisNode := ZKNode{
		Children: children,
		Data: string(data),
		Path: relativePath,
		IsEmpty: len(data) == 0,
	}
	thisNode.SHA256Checksum = thisNode.GenerateNodeHash()
	return &thisNode, nil
}