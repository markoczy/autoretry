package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/oliveagle/jsonpath"
)

func main() {
	pathPtr := flag.String("path", "$", "The JsonPath")
	flag.Parse()

	d, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var jsonData interface{}
	json.Unmarshal(d, &jsonData)

	res, err := jsonpath.JsonPathLookup(jsonData, *pathPtr)
	fmt.Println(res)
}
