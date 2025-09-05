// Package optparser provides parsing por Argv shell commands
package optparser

import (
	"fmt"
	"strings"
	"slices"
	"strconv"
)

type OptionError struct {
	Msg    string
	Option string
	Code   int
}
func (e OptionError) Error() string {
	return fmt.Sprintf("error in option: %v.\nmsg: %v.", e.Msg, e.Option)
}


type Option struct {
	Name    string
	Alias   rune
	Type    string
	Value   any
}
func (o *Option) Set(value any) {
	o.Value = value
}

func (o *Option) Typer() error {
	var value string = fmt.Sprintf("%v", o.Value)

	switch o.Type {
	case "string":
		// Nothing, continues the same, only converts in string
		o.Set(value)
		return nil

	case "boolean":
		Fvalue := strings.ToLower(value)	

		if Fvalue == "true" {
			o.Set(true)
			return nil

		} else if Fvalue == "false" {
			o.Set(false)
			return nil

		} else {
			return OptionError{Msg: "type does not match in boolean", Option: o.Name, Code: 5}
		}

	case "integer":
		res, err := strconv.Atoi(value)
		if err == nil {
			o.Set(res)
			return nil
		} else {
			return OptionError{Msg: "type does not match in integer", Option: o.Name, Code: 6}
		}

	default:
		return OptionError{Msg: "type does not exists", Option: o.Name, Code: 7}
	}

}


type OptionList []Option
func (l OptionList) existsName(sufix string) (Option, error) {
	for i, v := range l {
		if v.Name == sufix {
			return l[i], nil
		}
	}

	return l[0], OptionError{Msg: "option does not exists", Option: sufix, Code: 1}
}

func (l OptionList) GetAliases() ([]rune) {
	var runeList []rune

	for _, option := range l {
		if option.Alias != 0 {
			runeList = append(runeList, option.Alias)
		}
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
		return OptionError{Msg: "alias does not exists in your list", Option: string(last), Code: 3}
	}
}

func (l OptionList) ParseAlias(sufix string) (OptionList, error) {
	var returnOpts OptionList

	for i, v := range l {
		if v.Alias == 0 {
			continue
		}

		for _, letter := range(sufix) {

			if v.Alias == letter {
				option := l[i]
				option.Set("true")
				returnOpts = append(returnOpts, option)

				break
			}

		}
	}
	
	if len(returnOpts) < len(sufix) {
		var missingLetter rune
		aliasses := returnOpts.GetAliases()

		for _, v := range []rune(sufix) {
			if !(slices.Contains(aliasses, v)) {
				missingLetter = v
				break
			}
		}

		return returnOpts, OptionError{Msg: "alias does not exists in your list", Option: string(missingLetter), Code: 3}

	} else if len(returnOpts) > 0 {
		return returnOpts, l.DiffAlias(returnOpts)
	} else {
		return returnOpts, OptionError{Msg: "there is no alias in your list", Option: sufix, Code: 4}
	}
}

func (l *OptionList) UniqueSlice() {
	var duplicate OptionList
	slices.Reverse(*l)

	for _, v := range *l {
		fun := func(o Option) bool {
			return o.Name == v.Name
		}

		if !(slices.ContainsFunc(duplicate, fun)) {
			duplicate = append(duplicate, v)
		}
	}	

	slices.Reverse(duplicate)
	*l = duplicate
}

func (l *OptionList) TyperAll() error {

	for i := range *l {
		err := (*l)[i].Typer() 
		if err !=  nil {
			return err
		}
	}

	return nil
}



func Parser(args []string, preset OptionList) (OptionList, []string, []string) {
	var params []string
	var unknows []string
	var value bool
	var opts OptionList
	var option Option
	splitEqual(&args)


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
				unknows = append(unknows, v)
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
				unknows = append(unknows, v)
				continue
			}
		} else if len(v) > 1 && v[:1] == "-" {
			var semiOpts OptionList
			var semiUnkws []string
			var setValue bool

			sufix := v[1: len(v) - 1]
			last := []string{
				"-" + string(v[len(v) - 1]),
			}

			if sufix != "" {
				semiOpts, _, semiUnkws = Parser(last, preset)

				if len(semiUnkws) == 0 {
					option = semiOpts[0]
				}
			} else {
				sufix = last[0][1:]
				setValue = true
			}
			unknows = append(unknows, semiUnkws...)


			res, err := preset.ParseAlias(sufix)
			if len(res) > 0 {
				if setValue {
					option = res[0]
					value = true
				}
	
				opts = append(opts, 
					append(semiOpts, res...)...
				)

			}

			if e, ok := err.(OptionError); ok {
				unknows = append(unknows, "-" + e.Option)
			}

			if setValue {
				continue
			}
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

	opts.UniqueSlice()
	return opts, params, unknows
}

func splitEqual(slice *[]string) {
	for i, str := range *slice {
		occ := strings.Index(str, "=")

		if occ != -1 {
			head := str[:occ]
			tail := str[occ+1:]

			(*slice)[i] = head
			*slice = slices.Insert(*slice, i+1, tail)
		}
	}
}
