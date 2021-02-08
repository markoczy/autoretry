package parser

import (
	"fmt"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

type ymlParser struct {
	log logger.Logger
}

func NewYml(log logger.Logger) def.Parser {
	return &ymlParser{log}
}

func (p *ymlParser) Parse(input string, cfg def.Config) (ret interface{}, err error) {
	var nodes []*yaml.Node
	var yamlData yaml.Node
	if err = yaml.Unmarshal([]byte(input), &yamlData); err != nil {
		p.log.Error("Failed to umarshall YML")
		return
	}
	expr, err := yamlpath.NewPath(cfg.Path)
	if nodes, err = expr.Find(&yamlData); err != nil {
		p.log.Error("Failed to apply Yamlpath expression")
		return
	}
	if len(nodes) == 1 {
		ret, err = p.decode(nodes[0])
		return
	}
	ret = nodes
	return
}

func (p *ymlParser) decode(n *yaml.Node) (ret interface{}, err error) {
	// TODO could probably distinct numeric and optimize using n.Kind
	s := ""
	if err = n.Decode(&s); err == nil {
		ret = s
		return
	}
	arr := []interface{}{}
	if err = n.Decode(&arr); err == nil {
		ret = arr
		return
	}
	mp := map[string]interface{}{}
	if err = n.Decode(&mp); err == nil {
		ret = mp
		return
	}
	err = fmt.Errorf("Could not decode yaml value")
	return
}
