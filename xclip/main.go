package main

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/markoczy/xtools/common/helpers"
)

func main() {
	input, err := helpers.ReadStdin()

	// Read mode
	if helpers.IsNoValue(err) {
		str, err := clipboard.ReadAll()
		check(err)
		fmt.Print(str)
		return
	}
	check(err)

	// Write Mode
	clipboard.WriteAll(input)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
