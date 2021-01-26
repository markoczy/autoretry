package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/markoczy/xtools/common/flags"
	"github.com/markoczy/xtools/common/helpers"
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
)

type config struct {
	ParseMode   mode
	FormatMode  mode
	Path        string
	XmlTextName string
}

func parseFlags() (config, error) {
	ret := config{}

	regEx := flags.NewSwitchable(".*")
	jpathEx := flags.NewSwitchable("$")
	xpathEx := flags.NewSwitchable("/")
	ypathEx := flags.NewSwitchable("$")
	format := flags.NewEnum([]string{"go", "json", "xml", "yml", "plain"}, "plain")

	flag.Var(regEx, "regex", "Switch to Regular Expression mode and optionally input a regex expression")
	flag.Var(jpathEx, "json", "Switch to JSON mode and optionally input a JsonPath expression")
	flag.Var(xpathEx, "xml", "Switch to XML mode and optionally input a XPath expression")
	flag.Var(ypathEx, "yaml", "Switch to YAML mode and optionally input a XPath expression")
	flag.Var(format, "format", "Format parsed value to another output, possible values are 'go', 'json', 'xml', 'yml', 'plain'")

	// parseRegex := flag.String("regex", ".*", "Switch to Regular Expression mode and input a regex expression")
	// parseJson := flag.String("json", "$", "Switch to json mode and input a JsonPath expression")
	// parseXml := flag.String("xml", "//", "Switch to xml mode and input a XPath expression")
	// format := flag.String("format", "default", "Format parsed value to another output, possible values are 'go', 'json', 'xml', 'yml', 'default'")

	// formatGo := flag.Bool("to-go", false, "Format parsed value to Golang output instead of default")
	flag.Parse()
	// fmt.Println(err)

	// flag.Usage = func() {
	// 	fmt.Println("Usage of xparse: xparse [-json|-regex|-xml] [-to-json|-to-go|-to-yml] path")
	// 	flag.PrintDefaults()
	// }

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
	}

	// if *modeRegex {
	// 	ret.ParseMode = parseRegex
	// }
	// if *modeJson {
	// 	ret.ParseMode = parseJson
	// }
	// if *modeXml {
	// 	ret.ParseMode = parseXml
	// }

	args := flag.Args()
	if len(args) == 1 {
		ret.Path = args[0]
	} else if len(args) > 1 {
		flag.Usage()
		fmt.Println("\n> Be sure to define flags before the path. Trailing arguments:", args[1:])
		os.Exit(1)
	}
	// ret.FormatJson = *formatJson

	return ret, nil
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		flag.Usage()
		fmt.Println("\n> ", err)
	}

	fmt.Println(cfg)

	d, err := helpers.ReadStdin()
	check(err)

	switch cfg.ParseMode {
	case modePlain:
		r := parsePlain(d, cfg)
		fmt.Println(r)
	case modeRegex:
		r := parseRegex(d, cfg)
		fmt.Println(r)
	case modeJson:
		r := parseJson(d, cfg)
		fmt.Println(r)
	case modeXml:
		r := parseXml(d, cfg)
		fmt.Println(r)
	case modeYml:
		r := parseYaml(d, cfg)
		fmt.Println(r)
	}

	// res, err := jsonpath.JsonPathLookup(jsonData, cfg.Path)
	// check(err)

	// if cfg.FormatJson {
	// 	data, err := json.Marshal(res)
	// 	check(err)
	// 	fmt.Println(string(data))
	// 	return
	// }
	// fmt.Println(res)
}

// func parseJson() interface{} {

// }

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parsePlain(input string, cfg config) []string {
	r := regexp.MustCompile(`\r?\n`)
	split := r.Split(input, -1)
	ret := []string{}
	for _, v := range split {
		if strings.Contains(v, cfg.Path) {
			ret = append(ret, v)
		}
	}
	return ret
}

func parseRegex(input string, cfg config) []string {
	filterRx := regexp.MustCompile(cfg.Path)
	r := regexp.MustCompile(`\r?\n`)

	split := r.Split(input, -1)
	ret := []string{}
	for _, v := range split {
		if filterRx.MatchString(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func parseJson(input string, cfg config) interface{} {
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	check(err)
	res, err := jsonpath.JsonPathLookup(jsonData, cfg.Path)
	check(err)
	return res
}

func parseXml(input string, cfg config) interface{} {
	// xpath.
	fmt.Println(input)
	doc, err := xmlquery.Parse(strings.NewReader(input))
	check(err)
	expr, err := xpath.Compile(cfg.Path)
	check(err)
	// nodes := xmlquery.Find(doc, path)
	// nodes[0].
	// expr.Select()

	ret := expr.Evaluate(xmlquery.CreateXPathNavigator(doc))
	switch v := ret.(type) {
	case bool:
		fmt.Println("bool")
		return v
	case float64:
		fmt.Println("float")
		return v
	case string:
		fmt.Println("string")
		return v
	case *xpath.NodeIterator:
		fmt.Println("nodeiterator")

		return xmlDecode(v, cfg)
		// v.MoveNext()

		// v.MoveNext()
		// v.Current().NodeType()
		// fmt.Println(v.Current().Value())
		// fmt.Println(v.Current())
	default:
		fmt.Println(v)

		// v.Current()
		// v.Current().

	}
	return nil

}

func parseYaml(input string, cfg config) interface{} {
	// fmt.Println("Inside parseYaml")
	var yamlData yaml.Node
	err := yaml.Unmarshal([]byte(input), &yamlData)
	check(err)
	expr, err := yamlpath.NewPath(cfg.Path)
	check(err)
	nodes, err := expr.Find(&yamlData)
	check(err)
	// fmt.Println("Nodes:", nodes)
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

func xmlDecode(it *xpath.NodeIterator, cfg config) interface{} {
	ret := []interface{}{}
	hasMore := true
	for hasMore {
		_, val := parseXmlNode(it.Current(), cfg)
		if val != nil {
			ret = append(ret, val)
		}
		hasMore = it.MoveNext()
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

var textNodeName = "@text"

func parseXmlNode(node xpath.NodeNavigator, cfg config) (string, interface{}) {
	fmt.Println("parseXmlNode, type =", node.NodeType())
	switch node.NodeType() {
	case xpath.TextNode:
		fmt.Println("parseXmlNode found TextNode, value:", node.Value())
		if isNormalizedEmpty(node.Value()) {
			return "", nil
		}
		return textNodeName, node.Value()
	case xpath.ElementNode:
		name := node.LocalName()
		k := []map[string]interface{}{}
		fmt.Println("parseXmlNode found ElementNode, name:", name)
		k = append(k, parseXmlElementNode(node, cfg))
		fmt.Println("parseXmlNode ElementNode value", k[0])
		for node.MoveToNext() {
			k = append(k, parseXmlElementNode(node, cfg))
		}
		if len(k) == 1 {
			return name, k[0]
		}
		return name, k
	default:
		return "", nil
	}
}

func parseXmlElementNode(node xpath.NodeNavigator, cfg config) map[string]interface{} {
	ret := map[string]interface{}{}
	for node.MoveToNextAttribute() {
		fmt.Println("parseXmlElementNode attribute", node.LocalName(), node.Value())
		ret[node.LocalName()] = node.Value()
	}
	if node.NodeType() == xpath.AttributeNode {
		node.MoveToParent()
	}
	if node.MoveToChild() {
		name, cur := parseXmlNode(node, cfg)
		fmt.Println("parseXmlElementNode child", name, cur)
		if cur != nil {
			ret[name] = cur
		}
		for node.MoveToNext() {
			fmt.Println("parseXmlElementNode found next")
			name, cur = parseXmlNode(node, cfg)
			if cur != nil {
				ret[name] = cur
			}
		}
		node.MoveToParent()
	}
	return ret
}

var trimable = regexp.MustCompile("\\s|\n|\r|\t")

func normalize(s string) string {
	return trimable.ReplaceAllString(s, "")
}

func isNormalizedEmpty(s string) bool {
	return normalize(s) == ""
}

// func later() {
// 	pathPtr := flag.String("path", "$", "The JsonPath")
// 	flag.Parse()
// 	fmt.Println(flag.CommandLine.Args())

// 	d, err := ioutil.ReadAll(os.Stdin)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var jsonData interface{}
// 	json.Unmarshal(d, &jsonData)

// 	res, err := jsonpath.JsonPathLookup(jsonData, *pathPtr)
// 	fmt.Println(res)

// 	data, err := json.Marshal(res)
// 	fmt.Println(string(data))
// }
