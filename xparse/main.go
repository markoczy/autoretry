package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/markoczy/xtools/common/flags"
	"github.com/markoczy/xtools/common/helpers"
	"github.com/markoczy/xtools/common/logger"
	"github.com/oliveagle/jsonpath"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

type mode int

const (
	modePlain = mode(iota)
	modeRegex
	modeJson
	modeXml
	modeYml
	modeGo
)

var log logger.Logger

type config struct {
	ParseMode  mode
	FormatMode mode
	Path       string
	// XML specific
	XmlTextField   string
	XmlChildPrefix string
	XmlAttrPrefix  string
}

func (c config) String() string {
	return fmt.Sprintf("config: [ParseMode: %v, FormatMode: %v, Path: %v, XmlTextField: %v, XmlChildPrefix: %v, XmlAttrPrefix: %v]", c.ParseMode, c.FormatMode, c.Path, c.XmlTextField, c.XmlChildPrefix, c.XmlAttrPrefix)
}

func parseFlags() (config, error) {
	ret := config{}
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
		ret.ParseMode = modeRegex
		ret.Path = regEx.String()
		exCount++
	}
	if jpathEx.Defined() {
		ret.ParseMode = modeJson
		ret.Path = jpathEx.String()
		exCount++
	}
	if xpathEx.Defined() {
		ret.ParseMode = modeXml
		ret.Path = xpathEx.String()
		exCount++
	}
	if ypathEx.Defined() {
		ret.ParseMode = modeYml
		ret.Path = ypathEx.String()
		exCount++
	}
	if exCount > 1 {
		return config{}, fmt.Errorf("Too many expressions, only use one of -regex, -json, -xml")
	}

	switch format.String() {
	case "plain":
		ret.FormatMode = modePlain
	case "json":
		ret.FormatMode = modeJson
	case "xml":
		ret.FormatMode = modeXml
	case "yml":
		ret.FormatMode = modeYml
	case "yaml":
		ret.FormatMode = modeYml
	case "go":
		ret.FormatMode = modeGo
	default:
		ret.FormatMode = modePlain
	}

	args := flag.Args()
	if ret.ParseMode == modePlain {
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
	case modePlain:
		val = parsePlain(d, cfg)
	case modeRegex:
		val = parseRegex(d, cfg)
	case modeJson:
		val = parseJson(d, cfg)
	case modeXml:
		val = parseXml(d, cfg)
	case modeYml:
		val = parseYaml(d, cfg)
	}

	s := ""
	switch cfg.FormatMode {
	case modeGo:
		s = fmt.Sprintf("%v", val)
	case modePlain:
		s = formatPlain(val)
	case modeJson:
		s = formatJson(val)
	case modeYml:
		s = formatYaml(val)
	}

	fmt.Println(s)
}

//==============================================================================
// Plain
//==============================================================================

func parsePlain(input string, cfg config) []interface{} {
	r := regexp.MustCompile(`\r?\n`)
	split := r.Split(input, -1)
	ret := []interface{}{}
	for _, v := range split {
		if strings.Contains(v, cfg.Path) {
			ret = append(ret, v)
		}
	}
	return ret
}

func formatPlain(val interface{}) string {
	ret := formatPlainVal(reflect.ValueOf(val))
	return strings.Join(ret, "\r\n")
}

func formatPlainVal(val reflect.Value) []string {
	log.Debug("*** Cur: %v", val)
	log.Debug("*** Kind: %v", val.Kind())

	switch val.Kind() {
	case reflect.Array:
		val.InterfaceData()
		return []string{fmt.Sprintf("%v", val.Interface())} //? works???
	case reflect.Slice:
		return formatPlainSlice(val)
	case reflect.Map:
		return formatPlainMap(val)
	case reflect.Interface:
		k := reflect.ValueOf(val.Interface())
		log.Debug("*** New Kind: %v", k.Kind())
		if k.Kind() == reflect.Interface {
			return []string{fmt.Sprintf("%v", val.Interface())}
		}
		return formatPlainVal(k)
	default:
		log.Debug("*** Leaf Value: %s", fmt.Sprintf("%v", val.Interface()))
		return []string{fmt.Sprintf("%v", val.Interface())} //? works???
	}
}

func formatPlainSlice(slice reflect.Value) []string {
	ret := []string{}
	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		log.Debug("Found next value in slice %v", v)
		ret = append(ret, formatPlainVal(v)...)
	}
	return ret
}

func formatPlainMap(m reflect.Value) []string {
	ret := []string{}
	it := m.MapRange()
	for it.Next() {
		log.Debug("*** Found next in map %v", it.Value())
		// we currently ignore the keys when formatting a map
		ret = append(ret, formatPlainVal(it.Value())...)
	}
	return ret
}

//==============================================================================
// Regex
//==============================================================================

func parseRegex(input string, cfg config) interface{} {
	filterRx := regexp.MustCompile(cfg.Path)
	if len(filterRx.SubexpNames()) > 0 {
		return parseRegexWithCaptureGroups(filterRx, input, cfg)
	}
	return parseRegexWithoutCaptureGroups(filterRx, input, cfg)
}

func parseRegexWithoutCaptureGroups(rx *regexp.Regexp, input string, cfg config) []string {
	r := regexp.MustCompile(`\r?\n`)
	split := r.Split(input, -1)
	ret := []string{}
	for _, v := range split {
		if rx.MatchString(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func parseRegexWithCaptureGroups(rx *regexp.Regexp, input string, cfg config) []map[string]string {
	r := regexp.MustCompile(`\r?\n`)
	split := r.Split(input, -1)
	ret := []map[string]string{}
	for _, v := range split {
		if !rx.MatchString(v) {
			continue
		}
		cur := map[string]string{}
		match := rx.FindStringSubmatch(v)
		for i, name := range rx.SubexpNames() {
			if i != 0 && name != "" {
				cur[name] = match[i]
			}
		}
		ret = append(ret, cur)
	}
	return ret
}

//==============================================================================
// JSON
//==============================================================================

func parseJson(input string, cfg config) interface{} {
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	check(err)
	res, err := jsonpath.JsonPathLookup(jsonData, cfg.Path)
	check(err)
	return res
}

func formatJson(val interface{}) string {
	ret, err := json.MarshalIndent(val, "", "  ")
	check(err)
	return string(ret)
}

//==============================================================================
// XML
//==============================================================================

func parseXml(input string, cfg config) interface{} {
	log.Debug("Input: %s", input)
	doc, err := xmlquery.Parse(strings.NewReader(input))
	check(err)
	expr, err := xpath.Compile(cfg.Path)
	check(err)

	ret := expr.Evaluate(xmlquery.CreateXPathNavigator(doc))
	switch v := ret.(type) {
	case bool:
		return v
	case float64:
		return v
	case string:
		return v
	case *xpath.NodeIterator:
		log.Info("Found nodeIterator, decoding")
		return xmlDecode(v, cfg)
	default:
		log.Error("Unhandled node type: %", v)
	}
	return nil
}

func xmlDecode(it *xpath.NodeIterator, cfg config) interface{} {
	ret := []interface{}{}
	for it.MoveNext() {
		_, val := parseXmlNode(it.Current(), cfg)
		log.Debug("xmlDecode Val is: %v", val)
		if val != nil {
			ret = append(ret, val)
		}
	}
	switch len(ret) {
	case 0:
		return nil
	case 1:
		return ret[0]
	default:
		return ret
	}
}

func parseXmlNode(node xpath.NodeNavigator, cfg config) (string, interface{}) {
	log.Debug("parseXmlNode, type = %v", node.NodeType())
	switch node.NodeType() {
	case xpath.TextNode:
		log.Debug("parseXmlNode found TextNode, value: %s", node.Value())
		if isNormalizedEmpty(node.Value()) {
			// ignore text between nodes
			return "", nil
		}
		return cfg.XmlTextField, node.Value()
	case xpath.RootNode:
		// TODO root node should not be named
		log.Debug("parseXmlNode found RootNode")
		name := "@root"
		k := []interface{}{}
		k = append(k, parseXmlElementNode(node, cfg))
		for node.MoveToNext() {
			k = append(k, parseXmlElementNode(node, cfg))
		}
		if len(k) == 1 {
			return name, k[0]
		}
		return name, k
	case xpath.ElementNode:
		name := node.LocalName()
		log.Debug("parseXmlNode found ElementNode, name: %s", name)
		k := []map[string]interface{}{}
		k = append(k, parseXmlElementNode(node, cfg))
		log.Debug("parseXmlNode ElementNode value %v", k[0])
		return cfg.XmlChildPrefix + name, k[0]
	case xpath.AttributeNode:
		log.Debug("Attribute node Name: %s, Value: %s", node.LocalName(), node.Value())
		return node.LocalName(), node.Value()
	default:
		return "", nil
	}
}

func parseXmlElementNode(node xpath.NodeNavigator, cfg config) map[string]interface{} {
	ret := map[string]interface{}{}
	for node.MoveToNextAttribute() {
		log.Debug("parseXmlElementNode attribute %s %s", node.LocalName(), node.Value())
		ret[cfg.XmlAttrPrefix+node.LocalName()] = node.Value()
	}
	if node.NodeType() == xpath.AttributeNode {
		node.MoveToParent()
	}
	if node.MoveToChild() {
		name, cur := parseXmlNode(node, cfg)
		log.Debug("parseXmlElementNode child %s %v", name, cur)
		if cur != nil {
			ret[name] = cur
		}
		for node.MoveToNext() {
			log.Debug("parseXmlElementNode found next")
			name, cur = parseXmlNode(node, cfg)
			if cur != nil {
				defined := ret[name]
				switch {
				case defined == nil:
					log.Debug("* parseXmlElementNode Found nil")
					ret[name] = cur
				case reflect.ValueOf(defined).Kind() == reflect.Slice:
					log.Debug("* parseXmlElementNode Found array")
					ret[name] = append(defined.([]interface{}), cur)
				default:
					log.Debug("* parseXmlElementNode Found other %v", reflect.ValueOf(defined).Kind())
					arr := []interface{}{}
					arr = append(arr, defined, cur)
					ret[name] = arr
				}
			}
		}
		node.MoveToParent()
	}
	return ret
}

//==============================================================================
// YML
//==============================================================================

func parseYaml(input string, cfg config) interface{} {
	var yamlData yaml.Node
	err := yaml.Unmarshal([]byte(input), &yamlData)
	check(err)
	expr, err := yamlpath.NewPath(cfg.Path)
	check(err)
	nodes, err := expr.Find(&yamlData)
	check(err)
	if len(nodes) == 1 {
		return yamlDecode(nodes[0])
	}

	return nodes
}

func yamlDecode(n *yaml.Node) interface{} {
	// TODO could probably distinct numeric and optimize using n.Kind
	s := ""
	if err := n.Decode(&s); err == nil {
		return s
	}
	arr := []interface{}{}
	if err := n.Decode(&arr); err == nil {
		return arr
	}
	mp := map[string]interface{}{}
	if err := n.Decode(&mp); err == nil {
		return mp
	}
	panic("Could not decode yaml value")
}

func formatYaml(val interface{}) string {
	ret, err := yaml.Marshal(val)
	check(err)
	return string(ret)
}

//==============================================================================
// Common Helpers
//==============================================================================

var trimable = regexp.MustCompile("\\s|\n|\r|\t")

func normalize(s string) string {
	return trimable.ReplaceAllString(s, "")
}

func isNormalizedEmpty(s string) bool {
	return normalize(s) == ""
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
