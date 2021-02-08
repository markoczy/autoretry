package formatter

import (
	"fmt"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
)

type defaultFormatter struct {
	log logger.Logger
}

func NewDefault(log logger.Logger) def.Formatter {
	return &defaultFormatter{log}
}

func (f *defaultFormatter) Format(val interface{}) (ret string, err error) {
	return fmt.Sprintf("%v", val), nil
}
