package main

import (
	"fmt"
	"os"
	"github.com/ArthurSpeziali/worgo/pkg/optparser"
	"github.com/ArthurSpeziali/worgo/pkg/grammar"
)

var ARGS = os.Args[1:]

func presetOpts() optparser.OptionList {
	return optparser.OptionList{
	}
}

func main() {
	data := `Hello 
			Darkness
			World`

	fmt.Println(
		grammar.Input(data),
	)
}
