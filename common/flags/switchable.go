package flags

// Switchable a hybrid flag that acts as a boolean flag as well as a string input flag. Initial value will be used when function 'Set()' sets to 'true', any Call to 'Set()' sets 'defined' to 'true'.
type Switchable struct {
	defined bool
	val     string
}

func (f *Switchable) Set(val string) error {
	f.defined = true
	if val != "true" {
		f.val = val
	}
	return nil
}

func (f *Switchable) IsBoolFlag() bool {
	return true
}

func (f *Switchable) Defined() bool {
	return f.defined
}

func (f *Switchable) String() string {
	return f.val
}

func NewSwitchable(initial string) *Switchable {
	return &Switchable{
		val: initial,
	}
}
