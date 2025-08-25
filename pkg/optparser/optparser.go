// Package optparser provides parsing por Argv shell commands
package optparser

import (
	"strings"
	"fmt"
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
	Name  string
	Alias string
	Type  string
	Value string
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


func Parser(args []string, preset OptionList) (OptionList, []string) {
	var params []string
	var value bool
	var opts OptionList
	var option Option

	for _, v := range args {

		if strings.Contains(v, "--") {
			sufix := v[2:]

			res, err := preset.existsName(sufix)	
			if err == nil {
				option = res
			} else {
				continue
			}
		} else if value {
			option.Set(v)
			value = false

			opts = append(opts, option)
			continue
		} else {
			params = append(params, v)
			continue
		}

		if option.Type != "boolean" {
			value = true
		} else {
			value = false
			opts = append(opts, option)
		}
	}

	return opts, params
}

