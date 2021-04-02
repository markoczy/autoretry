package formatter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
)

type jpropFormatter struct {
	log logger.Logger
}

func NewJprop(log logger.Logger) def.Formatter {
	return &jpropFormatter{log}
}

func (f *jpropFormatter) Format(val interface{}) (ret string, err error) {
	x := f.formatVal(reflect.ValueOf(val), "")
	ret = strings.Join(x, "\r\n")
	return
}

func (f *jpropFormatter) formatVal(val reflect.Value, prefix string) []string {
	f.log.Debug("*** Cur: %v", val)
	f.log.Debug("*** Kind: %v", val.Kind())

	switch val.Kind() {
	case reflect.Array:
		val.InterfaceData()
		return []string{fmt.Sprintf("%v", val.Interface())} //? works???
	case reflect.Slice:
		return f.formatSlice(val, prefix)
	case reflect.Map:
		return f.formatMap(val, prefix)
	case reflect.Interface:
		k := reflect.ValueOf(val.Interface())
		f.log.Debug("*** New Kind: %v", k.Kind())
		if k.Kind() == reflect.Interface {
			return []string{fmt.Sprintf("%v", val.Interface())}
		}
		return f.formatVal(k, prefix)
	case reflect.Invalid:
		return []string{}
	default:
		f.log.Debug("*** Leaf Value: %s", fmt.Sprintf("%v", val.Interface()))
		return []string{fmt.Sprintf("%s = %v", prefix, val.Interface())}
	}
}

func (f *jpropFormatter) formatSlice(slice reflect.Value, prefix string) []string {
	ret := []string{}
	formatArray := len(prefix) != 0
	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		f.log.Debug("Found next value in slice %v", v)
		if formatArray {
			if i > 0 {
				prefix = prefix[:strings.LastIndex(prefix, "[")]
			}
			prefix = prefix + "[" + strconv.Itoa(i) + "]"
		}
		ret = append(ret, f.formatVal(v, prefix)...)
	}
	return ret
}

func (f *jpropFormatter) formatMap(m reflect.Value, prefix string) []string {
	ret := []string{}
	it := m.MapRange()
	for it.Next() {
		f.log.Debug("*** Found next in map %v", it.Value())
		// we currently ignore the keys when formatting a map
		if prefix == "" {
			prefix = it.Key().String()
		} else {
			prefix = prefix + "." + it.Key().String()
		}
		ret = append(ret, f.formatVal(it.Value(), prefix)...)
	}
	return ret
}
