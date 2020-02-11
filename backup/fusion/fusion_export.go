package fusion

import (
	"bsb-pr-solr-snapshot-service/backup/blob"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const FusionExportPath = "api/v1/zk/export/?encodeValues=utf-8"

type ExportParamsHeader struct{
	ZkHost string `json:"zkHost"`
	Path string `json:"path"`
	EncodeValues string `json:"encodeValues"`
	ExcludePaths []string`json:"excludePaths"`
	Recursive bool `json:"recursive"`
	Ephemeral bool `json:"ephemeral"`
}

type ExportRequestHeader struct {
	Timestamp string           `json:"timestamp"`
	Params *ExportParamsHeader `json:"params"`
}

type ExportNode struct{
	Path string            `json:"path"`
	Children []*ExportNode `json:"children"`
	Data string            `json:"data"`
}

type ExportResponse struct{
	Request *ExportRequestHeader `json:"request"`
	Response *ExportNode         `json:"response"`
}
func (exportNode *ExportNode) GetChildren()  []blob.Readable{
	s := make([]blob.Readable, len(exportNode.Children))
	for i, v := range exportNode.Children {
		s[i] = v
	}
	return s
}
func (exportNode *ExportNode) GetData() string{
	return exportNode.Data
}
func (exportNode *ExportNode) GetIsEmpty() bool{
	return exportNode.Data == ""
}
func (exportNode *ExportNode) GetPath() string{
	return exportNode.Path
}
func GetZKExport(baseUrl string) (*ExportResponse, error) {
	fusionExportUrl := fmt.Sprintf("%s/%s", baseUrl, FusionExportPath)
	exportResponse := ExportResponse{}
	client := resty.New()
	resp, err := client.R().EnableTrace().Get(fusionExportUrl)
	if err != nil {
		log.Error("Failed to make GET call for ", fusionExportUrl, " error ", err.Error())
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &exportResponse);
	if err != nil {
		log.Error("Failed to unmarshal JSON error ", err.Error())
	}
	return &exportResponse, nil
}
