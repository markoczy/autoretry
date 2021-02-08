package formatter

import (
	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
	"gopkg.in/yaml.v3"
)

type ymlFormatter struct {
	log logger.Logger
}

func NewYml(log logger.Logger) def.Formatter {
	return &ymlFormatter{log}
}

func (f *ymlFormatter) Format(val interface{}) (ret string, err error) {
	var b []byte
	if b, err = yaml.Marshal(val); err != nil {
		return
	}
	ret = string(b)
	return
}
