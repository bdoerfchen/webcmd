package configloader

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/bdoerfchen/webcmd/src/model/config"
)

type configUnmarshaller interface {
	Name() string
	Unmarshal(content []byte) (*config.AppConfig, error)
}

var unmarshallers = []configUnmarshaller{
	&jsonUnmarshaller{},
}

var (
	errTypeNotCorrect = errors.New("type not correct")
	errFormatBroken   = errors.New("broken format")
	errFormatUnknown  = errors.New("unknown format")
)

type jsonUnmarshaller struct{}

func (u *jsonUnmarshaller) Name() string { return "json" }

func (u *jsonUnmarshaller) Unmarshal(content []byte) (*config.AppConfig, error) {
	var result *config.AppConfig

	decoder := json.NewDecoder(bytes.NewReader(content))
	err := decoder.Decode(&result)

	return result, err
}
