package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
)

type errNoValue error

var noValue = errNoValue(fmt.Errorf("No Value"))

func IsNoValue(err error) bool {
	switch err.(type) {
	case errNoValue:
		return true
	}
	return false
}

func ReadStdin() (string, error) {
	var err error
	var stat os.FileInfo
	var data []byte
	input := os.Stdin
	if stat, err = input.Stat(); err != nil {
		return "", err
	}
	if stat.Mode()&os.ModeNamedPipe == 0 {
		return "", errNoValue(fmt.Errorf("No Value"))
	}
	if data, err = ioutil.ReadAll(input); err != nil {
		return "", err
	}
	return string(data), nil
}
