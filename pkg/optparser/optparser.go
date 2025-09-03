// Package optparser provides parsing por Argv shell commands
package optparser

import (
	"fmt"
	"slices"
)

type OptionError struct {
	Msg    string
	Option string
	Type   string
}
func (e OptionError) Error() string {
	return fmt.Sprintf("error in option: %v.\nmsg: %v.", e.Msg, e.Option)
}


type Option struct {
	Name    string
	Alias   rune
	Type    string
	Value   string
}
func (o *Option) Set(value string) {
	o.Value = value
}

type OptionList []Option
func (l OptionList) existsName(sufix string) (Option, error) {
	for i, v := range l {
		if v.Name == sufix {
			return l[i], nil
		}
	}

	return l[0], OptionError{Msg: "sulfix does not exists", Option: sufix}
}

func (l OptionList) GetAliases() ([]rune) {
	var runeList []rune

	for _, option := range l {
		runeList = append(runeList, option.Alias)
	}

	return runeList
}

func (l OptionList) DiffAlias(opts OptionList) error {
	var last rune
	allAliases := l.GetAliases()
	givenAliases := opts.GetAliases()
	
	for _, v := range givenAliases {
		if !(slices.Contains(allAliases, v)) {
			last = v
			break
		}
	}

	if last == 0 {
		return nil
	} else {
		return OptionError{Msg: "option does not exists", Option: string(last)}
	}
}

func (l OptionList) ParseAlias(sufix string) (OptionList, error) {
	var returnOpts OptionList

	for i, v := range l {
		for _, letter := range(sufix) {

			if v.Alias == letter {
				option := l[i]
				option.Set("true")
				returnOpts = append(returnOpts, option)
			}

		}
	}
	
	return returnOpts, nil
}


func Parser(args []string, preset OptionList) (OptionList, []string) {
	var params []string
	// var unknows []string
	var value bool
	var opts OptionList
	var option Option

	for _, v := range args {

		if value {
			option.Set(v)
			value = false

			opts = append(opts, option)
			continue
		} else if len(v) > 5 && v[:5] == "--no-" {
			sufix := v[5:]

			res, err := preset.existsName(sufix)
			if err == nil {
				option = res
			} else {
				continue
			}

			if option.Type == "boolean" {
				option.Set("false")
			}
			
			opts = append(opts, option)
			continue
		} else if len(v) > 2 && v[:2] == "--" {
			sufix := v[2:]

			res, err := preset.existsName(sufix)	
			if err == nil {
				option = res
			} else {
				continue
			}
		} else if len(v) > 1 && v[:1] == "-" {
			sufix := v[1:]

			res, _ := preset.ParseAlias(sufix)
			opts = append(opts, res...)
			continue
		} else {
			params = append(params, v)
			continue
		}

		if option.Type != "boolean" {
			value = true
		} else {
			value = false
			option.Set("true")
			opts = append(opts, option)
		}
	}

	return opts, params
}

// func substractSlices[T comparable](first []T, second []T) []T {
// 	var found bool
// 	var more []T

// 	for _, v := range first {
// 		found = false

// 		for j, w := range second {
// 			if v == w {
// 				second = append(second[:j], second[j+1:]...)
// 				found = true
// 				break
// 			} 
// 		}

// 		if !found {
// 			more = append(more, v)
// 		}
// 	}

// 	return append(second, more...)
// }
