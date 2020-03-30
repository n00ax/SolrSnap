package blob

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v4"
	"io/ioutil"
)

func Read(filepath string) (*Image, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read img blob %s", err.Error())
	}
	pre := Image{}
	err = msgpack.Unmarshal(bytes, &pre)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall image error %s", err.Error())
	}
	return &pre, nil
}
