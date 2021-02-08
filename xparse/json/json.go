package json

import (
	"encoding/json"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
	"github.com/oliveagle/jsonpath"
)

type parser struct {
	log logger.Logger
}
type formatter struct {
	log logger.Logger
}

func NewParser(log logger.Logger) def.Parser {
	return &parser{log}
}

func NewFormatter(log logger.Logger) def.Formatter {
	return &formatter{log}
}

func (p *parser) Parse(input string, cfg def.Config) (ret interface{}, err error) {
	var jsonData interface{}
	if err = json.Unmarshal([]byte(input), &jsonData); err != nil {
		return
	}
	ret, err = jsonpath.JsonPathLookup(jsonData, cfg.Path)
	return
}

func (f *formatter) Format(val interface{}) (ret string, err error) {
	var b []byte
	if b, err = json.MarshalIndent(val, "", "  "); err != nil {
		return
	}
	ret = string(b)
	return
}
