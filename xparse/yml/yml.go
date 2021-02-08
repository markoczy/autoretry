package yml

import (
	"fmt"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
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

func (p *parser) decode(n *yaml.Node) (ret interface{}, err error) {
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

func (f *formatter) Format(val interface{}) (ret string, err error) {
	var b []byte
	if b, err = yaml.Marshal(val); err != nil {
		return
	}
	ret = string(b)
	return
}
