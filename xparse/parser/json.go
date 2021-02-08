package parser

import (
	"encoding/json"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
	"github.com/oliveagle/jsonpath"
)

type jsonParser struct {
	log logger.Logger
}

func NewJson(log logger.Logger) def.Parser {
	return &jsonParser{log}
}

func (p *jsonParser) Parse(input string, cfg def.Config) (ret interface{}, err error) {
	var jsonData interface{}
	if err = json.Unmarshal([]byte(input), &jsonData); err != nil {
		return
	}
	ret, err = jsonpath.JsonPathLookup(jsonData, cfg.Path)
	return
}
