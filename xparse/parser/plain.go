package parser

import (
	"regexp"
	"strings"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
)

type plainParser struct {
	log logger.Logger
}

func NewPlain(log logger.Logger) def.Parser {
	return &plainParser{log}
}

func (p *plainParser) Parse(input string, cfg def.Config) (ret interface{}, err error) {
	r := regexp.MustCompile(`\r?\n`)
	split := r.Split(input, -1)
	x := []interface{}{}
	for _, v := range split {
		if strings.Contains(v, cfg.Path) {
			x = append(x, v)
		}
	}
	return
}
