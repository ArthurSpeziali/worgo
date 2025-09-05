package main

import (
	"fmt"
	"os"
	"github.com/ArthurSpeziali/worgo/pkg/optparser"
)

var ARGS = os.Args[1:]

func presetOpts() optparser.OptionList {
	return optparser.OptionList{
		{Name: "export", Type: "boolean"},
		{Name: "counter", Type: "integer"},
		{Name: "print", Type: "boolean"},
		{Name: "list", Type: "boolean", Alias: 'l'},
		{Name: "wright", Type: "boolean", Alias: 'w'},
		{Name: "olsen", Type: "boolean", Alias: 'o'},
		{Name: "name", Type: "string", Alias: 'n'},
	}
}

func main() {
	preset := presetOpts()
	ARGS = []string{
		"--no-print",
		"--export",
		"--counter",
		"1",
		"--unknow",
		"download",
		"-lawon",
		"hello world",
		"cmd",
	}

	opts, params, unknows := optparser.Parser(ARGS, preset)
	fmt.Printf("P: %q\nO: %q\nU: %q\n", params, opts, unknows)

	opts.TyperAll()
	fmt.Printf("\nF:%q", opts)
}
