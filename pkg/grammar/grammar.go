// Package grammar
package grammar


import (
	"os"
	"fmt"
)

type GrammarError struct {
	Msg   string
	Code  int
}

func (e GrammarError) Error() string {
	return fmt.Sprint(e.Msg)
}


func File(path string) (string, error) {
	var returnError error
	data, err := os.ReadFile(path)
	
	if err != nil {
		returnError = GrammarError{Msg: "file not found or incorrect permissions", Code: 1}
	}
	
	res, err := Input(
		string(data),
	)

	if returnError != nil {
		return res, err
	} else {
		return res, returnError
	}
}

func Input(data string) (string, error) {
	fmt.Printf("Your string: %s", data)
	return "a", nil
}
