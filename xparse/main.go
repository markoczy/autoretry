package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/markoczy/xtools/common/flags"
	"github.com/markoczy/xtools/common/helpers"
	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
	"github.com/markoczy/xtools/xparse/json"
	"github.com/markoczy/xtools/xparse/plain"
	"github.com/markoczy/xtools/xparse/regex"
	"github.com/markoczy/xtools/xparse/xml"
	"github.com/markoczy/xtools/xparse/yml"
)

var log logger.Logger

func parseFlags() (def.Config, error) {
	ret := def.Config{}
	logFactory := logger.NewAutoFlagFactory()

	regEx := flags.NewSwitchable(".*")
	jpathEx := flags.NewSwitchable("$")
	xpathEx := flags.NewSwitchable("/")
	ypathEx := flags.NewSwitchable("$")
	format := flags.NewEnum([]string{"go", "json", "xml", "yml", "yaml", "plain"}, "plain")

	xmlTextFieldPtr := flag.String("xml-text-field", "@text", "(Only applies when parsing xml) Field name for inner text when parsing XML to map-like structure")
	xmlChildPrefixPtr := flag.String("xml-child-prefix", "", "(Only applies when parsing xml) Prefix for children when XML to map-like structure")
	xmlAttrPrefixPtr := flag.String("xml-attr-prefix", "", "(Only applies when parsing xml) Prefix for attribute fields when paring XML to map-like structure")

	flag.Var(regEx, "regex", "Switch to Regular Expression mode and input a regex expression")
	flag.Var(jpathEx, "json", "Switch to JSON mode and input a JsonPath expression")
	flag.Var(xpathEx, "xml", "Switch to XML mode and input a XPath expression")
	flag.Var(ypathEx, "yaml", "Switch to YAML mode and input a XPath expression")
	flag.Var(format, "format", "Format parsed value to another output, possible values are 'go', 'json', 'xml', 'yml', 'plain'")

	logFactory.InitFlags()
	flag.Parse()

	ret.XmlTextField = *xmlTextFieldPtr
	ret.XmlChildPrefix = *xmlChildPrefixPtr
	ret.XmlAttrPrefix = *xmlAttrPrefixPtr
	log = logFactory.Create()

	exCount := 0
	if regEx.Defined() {
		ret.ParseMode = def.ModeRegex
		ret.Path = regEx.String()
		exCount++
	}
	if jpathEx.Defined() {
		ret.ParseMode = def.ModeJson
		ret.Path = jpathEx.String()
		exCount++
	}
	if xpathEx.Defined() {
		ret.ParseMode = def.ModeXml
		ret.Path = xpathEx.String()
		exCount++
	}
	if ypathEx.Defined() {
		ret.ParseMode = def.ModeYml
		ret.Path = ypathEx.String()
		exCount++
	}
	if exCount > 1 {
		return def.Config{}, fmt.Errorf("Too many expressions, only use one of -regex, -json, -xml")
	}

	switch format.String() {
	case "plain":
		ret.FormatMode = def.ModePlain
	case "json":
		ret.FormatMode = def.ModeJson
	case "xml":
		ret.FormatMode = def.ModeXml
	case "yml":
		ret.FormatMode = def.ModeYml
	case "yaml":
		ret.FormatMode = def.ModeYml
	case "go":
		ret.FormatMode = def.ModeGo
	default:
		ret.FormatMode = def.ModePlain
	}

	args := flag.Args()
	if ret.ParseMode == def.ModePlain {
		switch len(args) {
		case 0:
			ret.Path = ""
		case 1:
			ret.Path = args[0]
		default:
			fmt.Println("Be sure to define flags before the path. Only one trailing argument is allowed for mode plain:", args[1:])
			flag.Usage()
			os.Exit(1)
		}
	} else if len(args) != 0 {
		fmt.Println("No trailing arguments allowed for when mode is not plain:", args[1:])
		flag.Usage()
		os.Exit(1)
	}

	return ret, nil
}

func main() {
	var parser def.Parser
	var formatter def.Formatter

	cfg, err := parseFlags()
	if err != nil {
		flag.Usage()
		fmt.Println("\n> ", err)
	}
	log.Debug("Config: %v", cfg)

	d, err := helpers.ReadStdin()
	check(err)

	var val interface{}
	switch cfg.ParseMode {
	case def.ModePlain:
		parser = plain.NewParser(log)
	case def.ModeRegex:
		parser = regex.NewParser(log)
	case def.ModeJson:
		parser = json.NewParser(log)
	case def.ModeXml:
		parser = xml.NewParser(log)
	case def.ModeYml:
		parser = yml.NewParser(log)
	}
	val, err = parser.Parse(d, cfg)
	check(err)

	s := ""
	switch cfg.FormatMode {
	case def.ModeGo:
		s = fmt.Sprintf("%v", val)
	case def.ModePlain:
		formatter = plain.NewFormatter(log)
	case def.ModeJson:
		formatter = json.NewFormatter(log)
	case def.ModeXml:
		formatter = xml.NewFormatter(log)
	case def.ModeYml:
		formatter = yml.NewFormatter(log)
	}
	s, err = formatter.Format(val)
	check(err)

	fmt.Println(s)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
