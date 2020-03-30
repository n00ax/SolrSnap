package blob

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/vmihailenco/msgpack/v4"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const Magic = "zkBLOBv1Image"
const Version = "1.0.0"

type Source string

const (
	FusionRestExport Source = "FusionRestExport"
	ZKDirectExport   Source = "ZKDirectExport"
)

type SourceApplication string

const (
	SolrSnapshotService SourceApplication = "solr-snapshot-service"
)

type Header struct {
	Magic string
	// Base64 encoded
	SHA256Checksum    string
	Version           string
	Source            Source
	SourceApplication SourceApplication
	GenerateUser      string
	Description       string
	TimeGenerated     time.Time
	SolrZKPath        string
}
type ZKNode struct {
	Children       []*ZKNode
	Data           string
	Path           string
	IsEmpty        bool
	SHA256Checksum string
}
type IndexDataRoot struct {
	SHA256Checksum string
	Files          []IndexFile
	SolrVersion    string
}
type IndexFile struct {
	RelativePath   string
	Permissions    os.FileMode
	Data           []byte
	SHA256Checksum []string
	IsDir          bool
}
type Image struct {
	Header    Header
	ZKRoot    *ZKNode
	IndexData *IndexDataRoot
}

// Validates image as Merkle tree
func (img *Image) ValidateImage() error {
	return img.ZKRoot.validateNode()
}
func (img *Image) Save(filepath string, perm os.FileMode) error {
	bytes, err := msgpack.Marshal(img)
	if err != nil {
		return fmt.Errorf("failed to marshall msgpack %s", err.Error())
	}
	err = ioutil.WriteFile(filepath, bytes, perm)
	if err != nil {
		return fmt.Errorf("failed to save msgpack %s", err.Error())
	}
	return nil
}
func (img *Image) DeleteZK(path string) error {
	parentNode, err := img.ZKRoot.GetNode(filepath.Dir(path))
	if err != nil {
		return err
	}
	for i, ch := range parentNode.Children {
		if ch.Path == path {
			parentNode.Children = append(parentNode.Children[:i], parentNode.Children[i+1:]...)
		}
	}
	img.regenChecksum()
	return nil
}
func (img *Image) ModifyZK(path string, newData string) error {
	node, err := img.ZKRoot.GetNode(path)
	if err != nil {
		return err
	}
	node.Data = newData
	img.regenChecksum()
	return nil
}
func (img *Image) regenChecksum() { // Inefficient, but i'm lazy
	newSum := img.ZKRoot.regenChecksum()
	img.Header.SHA256Checksum = newSum
}

type Readable interface {
	GetChildren() []Readable
	GetIsEmpty() bool
	GetData() string
	GetPath() string
}

// ZK tree is validated by Merkle hashes with Data:Path:{children hashes concatenated with :}
func (zkNode *ZKNode) GenerateNodeHash() string {
	content := fmt.Sprintf("%s:%s", zkNode.Data, zkNode.Path)
	for _, child := range zkNode.Children {
		content = fmt.Sprintf("%s:%s", content, child.SHA256Checksum)
	}
	hash := sha256.Sum256([]byte(content))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func (zkNode *ZKNode) validateNode() error {
	for _, child := range zkNode.Children {
		err := child.validateNode()
		if err != nil {
			return err
		}
	}
	if zkNode.SHA256Checksum == zkNode.GenerateNodeHash() {
		//Valid leaf
		return nil
	} else {
		return fmt.Errorf("bad image tree, Merkle tree discrepency at %s, computed, %s stored %s",
			zkNode.Path, zkNode.GenerateNodeHash(), zkNode.SHA256Checksum)
	}
}
func (zkNode *ZKNode) GetChildren() [] Readable {
	s := make([]Readable, len(zkNode.Children))
	for i, v := range zkNode.Children {
		s[i] = v
	}
	return s
}
func (zkNode *ZKNode) GetData() string {
	return zkNode.Data
}
func (zkNode *ZKNode) GetIsEmpty() bool {
	return zkNode.IsEmpty
}
func (zkNode *ZKNode) GetPath() string {
	return zkNode.Path
}
func (zkNode *ZKNode) GetObjCount() int {
	objcount := 0
	for _, child := range zkNode.Children {
		objcount += child.GetObjCount()
	}
	return objcount + 1
}
func (zkNode *ZKNode) GetNode(path string) (*ZKNode, error) {
	if zkNode.Path == path {
		return zkNode, nil
	}
	var match *ZKNode
	for _, child := range zkNode.Children {
		if child.Path == path {
			return child, nil
		} else {
			ch, err := child.GetNode(path)
			if err == nil {
				match = ch
			}
		}
	}
	if match != nil {
		return match, nil
	} else {
		return zkNode, fmt.Errorf("could not find node with path %s", path)
	}
}
func (zkNode *ZKNode) regenChecksum() string {
	zkNode.SHA256Checksum = zkNode.GenerateNodeHash()
	for _, node := range zkNode.Children {
		node.regenChecksum()
	}
	return zkNode.SHA256Checksum
}
