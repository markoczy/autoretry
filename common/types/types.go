package types

import "fmt"

type Enum interface {
	ValueOf(string) (EnumValue, error)
}

type EnumValue interface {
	String() string
}

type enumValue string

func (v *enumValue) String() string {
	return string(*v)
}

type enum struct {
	values map[string]*enumValue
}

func (e *enum) ValueOf(s string) (EnumValue, error) {
	ret := e.values[s]
	if ret == nil {
		return nil, fmt.Errorf("Enum value %s does not exist", s)
	}
	return ret, nil
}

func NewEnum(values []string) Enum {
	enumValues := map[string]*enumValue{}
	for _, v := range values {
		k := enumValue(v)
		enumValues[v] = &k
	}
	return &enum{values: enumValues}
}
