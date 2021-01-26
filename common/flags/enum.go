package flags

import "fmt"

// Enum a flag that only accepts one of the provided values.
type Enum struct {
	values []string
	val    string
}

func (f *Enum) Set(val string) error {
	for _, v := range f.values {
		if v == val {
			f.val = val
			return nil
		}
	}
	return fmt.Errorf("Only one of %v is allowed", f.values)
}

func (f *Enum) String() string {
	return f.val
}

func NewEnum(values []string, initial string) *Enum {
	return &Enum{
		values: values,
		val:    initial,
	}
}
