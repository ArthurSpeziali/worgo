// Package optparser provides parsing por Argv shell commands
package optparser

import (
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

func (l OptionList) parseAlias(sufix string) (OptionList, error) {
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

    if len(sufix) > len(returnOpts) {
		return returnOpts, OptionError{Msg: "some aliases do not exist in sulfix"}
	} else {
		return returnOpts, nil
	}
}


func Parser(args []string, preset OptionList) (OptionList, []string) {
	var params []string
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

			res, _ := preset.parseAlias(sufix)
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

