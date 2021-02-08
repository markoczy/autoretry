package def

import "fmt"

type Mode int

const (
	ModePlain = Mode(iota)
	ModeRegex
	ModeJson
	ModeXml
	ModeYml
	ModeGo
)

type Parser interface {
	Parse(input string, cfg Config) (interface{}, error)
}

type Formatter interface {
	Format(val interface{}) (string, error)
}

type Config struct {
	ParseMode  Mode
	FormatMode Mode
	Path       string
	// XML specific
	XmlTextField   string
	XmlChildPrefix string
	XmlAttrPrefix  string
}

func (c Config) String() string {
	return fmt.Sprintf("Config: [ParseMode: %v, FormatMode: %v, Path: %v, XmlTextField: %v, XmlChildPrefix: %v, XmlAttrPrefix: %v]", c.ParseMode, c.FormatMode, c.Path, c.XmlTextField, c.XmlChildPrefix, c.XmlAttrPrefix)
}
