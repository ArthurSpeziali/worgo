package main

import (
	"fmt"
	"os"
	"github.com/ArthurSpeziali/worgo/pkg/optparser"
)

var ARGS = os.Args[1:]
func main() {
	preset := optparser.OptionList{
		{Name: "export", Type: "boolean"},
		{Name: "counter", Type: "integer"},
	}

	// ARGS := []string{
	// 	"--export",
	// 	"--counter",
	// 	"1",
	// 	"--unknow",
	// 	"download",
	// }

	opts, params := optparser.Parser(ARGS, preset)
	fmt.Printf("P: %q\nO: %q\n", params, opts)
}
