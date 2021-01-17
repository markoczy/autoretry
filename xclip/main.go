package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	input := os.Stdin
	stat, err := input.Stat()
	check(err)

	// Read mode
	if stat.Mode()&os.ModeNamedPipe == 0 {
		str, err := clipboard.ReadAll()
		check(err)
		fmt.Print(str)
		return
	}

	// Write Mode
	data, err := ioutil.ReadAll(input)
	check(err)
	clipboard.WriteAll(string(data))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
