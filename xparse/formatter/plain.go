package formatter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/markoczy/xtools/common/logger"
	"github.com/markoczy/xtools/xparse/def"
)

type plainFormatter struct {
	log logger.Logger
}

func NewPlain(log logger.Logger) def.Formatter {
	return &plainFormatter{log}
}

func (f *plainFormatter) Format(val interface{}) (ret string, err error) {
	x := f.formatVal(reflect.ValueOf(val))
	ret = strings.Join(x, "\r\n")
	return
}

func (f *plainFormatter) formatVal(val reflect.Value) []string {
	f.log.Debug("*** Cur: %v", val)
	f.log.Debug("*** Kind: %v", val.Kind())

	switch val.Kind() {
	case reflect.Array:
		val.InterfaceData()
		return []string{fmt.Sprintf("%v", val.Interface())} //? works???
	case reflect.Slice:
		return f.formatSlice(val)
	case reflect.Map:
		return f.formatMap(val)
	case reflect.Interface:
		k := reflect.ValueOf(val.Interface())
		f.log.Debug("*** New Kind: %v", k.Kind())
		if k.Kind() == reflect.Interface {
			return []string{fmt.Sprintf("%v", val.Interface())}
		}
		return f.formatVal(k)
	default:
		f.log.Debug("*** Leaf Value: %s", fmt.Sprintf("%v", val.Interface()))
		return []string{fmt.Sprintf("%v", val.Interface())} //? works???
	}
}

func (f *plainFormatter) formatSlice(slice reflect.Value) []string {
	ret := []string{}
	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		f.log.Debug("Found next value in slice %v", v)
		ret = append(ret, f.formatVal(v)...)
	}
	return ret
}

func (f *plainFormatter) formatMap(m reflect.Value) []string {
	ret := []string{}
	it := m.MapRange()
	for it.Next() {
		f.log.Debug("*** Found next in map %v", it.Value())
		// we currently ignore the keys when formatting a map
		ret = append(ret, f.formatVal(it.Value())...)
	}
	return ret
}
