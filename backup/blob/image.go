package blob

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/vmihailenco/msgpack/v4"
	"io/ioutil"
	"os"
	"time"
)

const Magic = "zkBLOBv1Image"
const Version = "1.0.0"
type Source string
const(
	FusionRestExport Source = "FusionRestExport"
	ZKDirectExport Source = "ZKDirectExport"
)
type SourceApplication string
const(
	SolrSnapshotService SourceApplication = "solr-snapshot-service"
)
type Header struct{
	Magic string
	// Base64 encoded
	SHA256Checksum string
	Version string
	Source Source
	SourceApplication SourceApplication
	GenerateUser string
	Description string
	TimeGenerated time.Time
}
type ZKNode struct{
	Children []*ZKNode
	Data string
	Path string
	IsEmpty bool
	SHA256Checksum string
}
type IndexDataRoot struct{
	SHA256Checksum string
	Files []IndexFile
	SolrVersion string
}
type IndexFile struct{
	RelativePath string
	Permissions os.FileMode
	Data []byte
	SHA256Checksum []string
	IsDir bool
}
type Image struct {
	Header Header
	ZKRoot * ZKNode
	IndexData * IndexDataRoot
}
// Validates image as Merkle tree
func (img *Image) ValidateImage() error {
	return img.ValidateImage()
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
	if zkNode.SHA256Checksum == zkNode.GenerateNodeHash(){
		//Valid leaf
		return nil
	} else {
		return fmt.Errorf("bad image tree, Merkle tree discrepency at %s, computed, %s stored %s",
			zkNode.Path, zkNode.GenerateNodeHash(), zkNode.SHA256Checksum)
	}
}
func (img *Image) Save(filepath string, perm os.FileMode) error{
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
type Readable interface {
	GetChildren() []Readable
	GetIsEmpty() bool
	GetData() string
	GetPath() string
}
func (zkNode *ZKNode) GetChildren()  [] Readable{
	s := make([]Readable, len(zkNode.Children))
	for i, v := range zkNode.Children {
		s[i] = v
	}
	return s
}
func (zkNode *ZKNode) GetData() string{
	return zkNode.Data
}
func (zkNode *ZKNode) GetIsEmpty() bool{
	return zkNode.IsEmpty
}
func (zkNode *ZKNode) GetPath() string{
	return zkNode.Path
}
