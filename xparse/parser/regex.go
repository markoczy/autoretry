package parser

import (
	"regexp"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
)

type regexParser struct {
	log logger.Logger
}

func NewRegex(log logger.Logger) def.Parser {
	return &regexParser{log}
}

func (p *regexParser) Parse(input string, cfg def.Config) (ret interface{}, err error) {
	filterRx := regexp.MustCompile(cfg.Path)
	if len(filterRx.SubexpNames()) > 0 {
		return p.parseRegexWithCaptureGroups(filterRx, input, cfg), nil
	}
	return p.parseRegexWithoutCaptureGroups(filterRx, input, cfg), nil
}

func (p *regexParser) parseRegexWithoutCaptureGroups(rx *regexp.Regexp, input string, cfg def.Config) []string {
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

func (p *regexParser) parseRegexWithCaptureGroups(rx *regexp.Regexp, input string, cfg def.Config) []map[string]string {
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
