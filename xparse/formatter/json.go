package formatter

import (
	"encoding/json"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
)

type jsonFormatter struct {
	log logger.Logger
}

func NewJson(log logger.Logger) def.Formatter {
	return &jsonFormatter{log}
}

func (f *jsonFormatter) Format(val interface{}) (ret string, err error) {
	var b []byte
	if b, err = json.MarshalIndent(val, "", "  "); err != nil {
		return
	}
	ret = string(b)
	return
}
