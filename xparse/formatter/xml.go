package formatter

import (
	"encoding/xml"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
)

type xmlFormatter struct {
	log logger.Logger
}

func NewXml(log logger.Logger) def.Formatter {
	return &xmlFormatter{log}
}

func (f *xmlFormatter) Format(val interface{}) (ret string, err error) {
	var b []byte

	if b, err = xml.MarshalIndent(val, "", "  "); err != nil {
		return
	}
	ret = string(b)
	return
}
