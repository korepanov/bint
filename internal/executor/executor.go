package executor

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"bint.com/internal/const/options"
	. "bint.com/internal/transpiler"
	. "bint.com/pkg/serviceTools"
)

func execute(systemStack []interface{}, OP string, LO []interface{}, RO []interface{}) ([]interface{}, []interface{}, error) {
	if "print" == OP {
		return []interface{}{"print", LO}, systemStack, nil
	} else if "reg_find" == OP {
		if len(fmt.Sprintf("%v", LO[0])) >= 2 && `"` == string(fmt.Sprintf("%v", LO[0])[0]) &&
			`"` == string(fmt.Sprintf("%v", LO[0])[len(fmt.Sprintf("%v", LO[0]))-1]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if len(fmt.Sprintf("%v", RO[0])) >= 2 && `"` == string(fmt.Sprintf("%v", RO[0])[0]) &&
			`"` == string(fmt.Sprintf("%v", RO[0])[len(fmt.Sprintf("%v", RO[0]))-1]) {
			RO[0] = RO[0].(string)[1 : len(RO[0].(string))-1]
		}

		pattern, err := regexp.Compile(fmt.Sprintf("%v", LO[0]))

		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		intListList := pattern.FindAllIndex([]byte(fmt.Sprintf("%v", RO[0])), -1)
		var res []interface{}
		res = append(res, []interface{}{"end"})

		for i := len(intListList) - 1; i >= 0; i-- {
			res = append(res, []interface{}{[]interface{}{"end"}})
			for j := len(intListList[i]) - 1; j >= 0; j-- {
				res[len(res)-1] = append(res[len(res)-1].([]interface{}), []interface{}{intListList[i][j]})
			}
		}

		return []interface{}{res}, systemStack, nil

	} else if "index" == OP {

		if len(fmt.Sprintf("%v", LO[0])) >= 2 && `"` == string(fmt.Sprintf("%v", LO[0])[0]) &&
			`"` == string(fmt.Sprintf("%v", LO[0])[len(fmt.Sprintf("%v", LO[0]))-1]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if len(fmt.Sprintf("%v", RO[0])) >= 2 && `"` == string(fmt.Sprintf("%v", RO[0])[0]) &&
			`"` == string(fmt.Sprintf("%v", RO[0])[len(fmt.Sprintf("%v", RO[0]))-1]) {
			RO[0] = RO[0].(string)[1 : len(RO[0].(string))-1]
		}

		return []interface{}{strings.Index(fmt.Sprintf("%v", LO[0]), fmt.Sprintf("%v", RO[0]))}, systemStack, nil
	} else if "len" == OP {
		if 0 == len(fmt.Sprintf("%v", LO[0])) {
			return []interface{}{0}, systemStack, nil
		}
		if `"` == string(fmt.Sprintf("%v", LO[0])[0]) && `"` == string(fmt.Sprintf("%v", LO[0])[len(fmt.Sprintf("%v", LO[0]))-1]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}

		return []interface{}{len([]rune(fmt.Sprintf("%v", LO[0])))}, systemStack, nil
	} else if "exists" == OP {
		if `"` == string(fmt.Sprintf("%v", LO[0])[0]) && `"` == string(fmt.Sprintf("%v", LO[0])[len(fmt.Sprintf("%v", LO[0]))-1]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		return []interface{}{Exists(fmt.Sprintf("%v", LO[0]))}, systemStack, nil
	} else if "AND" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) && "bool" == WhatsType(fmt.Sprintf("%v", RO[0])) {
			return []interface{}{StrToBool(fmt.Sprintf("%v", LO[0])) && StrToBool(fmt.Sprintf("%v", RO[0]))}, systemStack, nil
		}
		err := errors.New("executor: AND: error: data type mismatch")
		return LO, systemStack, err
	} else if "OR" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) && "bool" == WhatsType(fmt.Sprintf("%v", RO[0])) {
			return []interface{}{StrToBool(fmt.Sprintf("%v", LO[0])) || StrToBool(fmt.Sprintf("%v", RO[0]))}, systemStack, nil
		}
		err := errors.New("executor: OR: error: data type mismatch")
		return LO, systemStack, err
	} else if "XOR" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) && "bool" == WhatsType(fmt.Sprintf("%v", RO[0])) {
			return []interface{}{StrToBool(fmt.Sprintf("%v", LO[0])) != StrToBool(fmt.Sprintf("%v", RO[0]))}, systemStack, nil
		}
		err := errors.New("executor: XOR: error: data type mismatch")
		return LO, systemStack, err
	} else if "NOT" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			return []interface{}{!StrToBool(fmt.Sprintf("%v", LO[0]))}, systemStack, nil
		}
		err := errors.New("executor: NOT: error: data type mismatch")
		return LO, systemStack, err
	} else if "L: True" == OP || "L: true" == OP {
		return LO, systemStack, nil
	} else if "L: False" == OP || "L: false" == OP {
		return RO, systemStack, nil
	} else if "<" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: < : error: data type mismatch")
			return LO, systemStack, err
		}

		floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

		if nil != err {
			return LO, systemStack, err
		}

		floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

		if nil != err {
			return RO, systemStack, err
		}

		return []interface{}{floatLO < floatRO}, systemStack, nil

	} else if "<=" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: <= : error: data type mismatch")
			return LO, systemStack, err
		}

		floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

		if nil != err {
			return LO, systemStack, err
		}

		floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

		if nil != err {
			return RO, systemStack, err
		}

		return []interface{}{floatLO <= floatRO}, systemStack, nil

	} else if "==" == OP {
		LOType := WhatsType(fmt.Sprintf("%v", LO[0]))
		ROType := WhatsType(fmt.Sprintf("%v", RO[0]))

		if "string" == LOType {
			ROType = "string"
		} else if "string" == ROType {
			LOType = "string"
		}

		if LOType != ROType {
			err := errors.New("executor: == : ERROR: data type mismatch")
			return LO, systemStack, err
		}

		if "int" == LOType {
			intLO, err := strconv.Atoi(fmt.Sprintf("%v", LO[0]))
			if nil != err {
				return LO, systemStack, err
			}

			intRO, err := strconv.Atoi(fmt.Sprintf("%v", RO[0]))
			if nil != err {
				return LO, systemStack, err
			}

			return []interface{}{BoolToStr(intLO == intRO)}, systemStack, nil
		}

		if "float" == LOType {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)
			if nil != err {
				return LO, systemStack, err
			}

			floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
			if nil != err {
				return LO, systemStack, err
			}

			return []interface{}{BoolToStr(floatLO == floatRO)}, systemStack, nil
		}

		if "string" == LOType {
			if len(fmt.Sprintf("%v", LO[0])) > 0 && "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
				LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
			}
			if len(fmt.Sprintf("%v", RO[0])) > 0 && "\"" == string(fmt.Sprintf("%v", RO[0])[0]) {
				RO[0] = RO[0].(string)[1 : len(RO[0].(string))-1]
			}

			if "\\n" == fmt.Sprintf("%v", LO[0]) {
				LO[0] = "\n"
			}

			if "\\n" == fmt.Sprintf("%v", RO[0]) {
				RO[0] = "\n"
			}

			if "\\t" == fmt.Sprintf("%v", LO[0]) {
				LO[0] = "\t"
			}

			if "\\t" == fmt.Sprintf("%v", RO[0]) {
				RO[0] = "\t"
			}

			return []interface{}{BoolToStr(fmt.Sprintf("%v", LO[0]) == fmt.Sprintf("%v", RO[0]))}, systemStack, nil
		}

		err := errors.New("executor: == : ERROR: data type mismatch")
		return LO, systemStack, err
	} else if ">" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: > : error: data type mismatch")
			return LO, systemStack, err
		}

		floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

		if nil != err {
			return LO, systemStack, err
		}

		floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

		if nil != err {
			return LO, systemStack, err
		}

		return []interface{}{BoolToStr(floatLO > floatRO)}, systemStack, nil
	} else if ">=" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: >= : error: data type mismatch")
			return LO, systemStack, err
		}

		floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

		if nil != err {
			return LO, systemStack, err
		}

		floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

		if nil != err {
			return LO, systemStack, err
		}

		return []interface{}{BoolToStr(floatLO >= floatRO)}, systemStack, nil
	} else if "+" == OP {

		if "int" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			intLO, err := strconv.Atoi(fmt.Sprintf("%v", LO[0]))
			if nil != err {
				return LO, systemStack, err
			}

			if "int" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				intRO, err := strconv.Atoi(fmt.Sprintf("%v", RO[0]))
				if nil != err {
					return RO, systemStack, err
				}

				return []interface{}{intLO + intRO}, systemStack, nil
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return RO, systemStack, err
				}
				floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)
				if nil != err {
					return RO, systemStack, err
				}

				return []interface{}{fmt.Sprintf("%.6f", floatLO+floatRO)}, systemStack, nil
			}

		}

		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) || "int" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{fmt.Sprintf("%.6f", floatLO+floatRO)}, systemStack, nil
			}
		}

		if len(fmt.Sprintf("%v", LO[0])) >= 2 && "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if len(fmt.Sprintf("%v", RO[0])) >= 2 && "\"" == string(fmt.Sprintf("%v", RO[0])[0]) {
			RO[0] = RO[0].(string)[1 : len(RO[0].(string))-1]
		}

		return []interface{}{"\"" + fmt.Sprintf("%v", LO[0]) + fmt.Sprintf("%v", RO[0]) + "\""}, systemStack, nil

	} else if "-" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: - : error: data type mismatch")
			return LO, systemStack, err
		}

		if "int" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			intLO, err := strconv.Atoi(fmt.Sprintf("%v", LO[0]))
			if nil != err {
				return LO, systemStack, err
			}

			if "int" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				intRO, err := strconv.Atoi(fmt.Sprintf("%v", RO[0]))
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{intLO - intRO}, systemStack, nil
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}
				floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{fmt.Sprintf("%.6f", floatLO-floatRO)}, systemStack, nil
			}
			err = errors.New("executor: - : error: data type mismatch")
			return LO, systemStack, err

		}

		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) || "int" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{fmt.Sprintf("%.6f", floatLO-floatRO)}, systemStack, nil
			}

			err = errors.New("executor: - : error: data type mismatch")
			return LO, systemStack, err
		}
	} else if "*" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: * : error: data type mismatch")
			return LO, systemStack, err
		}

		if "int" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			intLO, err := strconv.Atoi(fmt.Sprintf("%v", LO[0]))
			if nil != err {
				return LO, systemStack, err
			}

			if "int" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				intRO, err := strconv.Atoi(fmt.Sprintf("%v", RO[0]))
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{intLO * intRO}, systemStack, nil
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}
				floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{fmt.Sprintf("%.6f", floatLO*floatRO)}, systemStack, nil
			}

			err = errors.New("executor: * : error: data type mismatch")
			return LO, systemStack, err

		}

		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) || "int" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{fmt.Sprintf("%.6f", floatLO*floatRO)}, systemStack, nil
			}

			err = errors.New("executor: * : error: data type mismatch")
			return LO, systemStack, err
		}
	} else if "/" == OP {
		if ("float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: / : error: data type mismatch")
			return LO, systemStack, err
		}

		floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)
		if nil != err {
			return LO, systemStack, err
		}
		floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
		if nil != err {
			return RO, systemStack, err
		}
		res := fmt.Sprintf("%.6f", floatLO/floatRO)

		return []interface{}{res}, systemStack, nil

	} else if "@" == OP {
		if !("int" == WhatsType(fmt.Sprintf("%v", LO[0])) && "int" == WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: @ : error: data type mismatch")
			return LO, systemStack, err
		}
		intLO, err := strconv.Atoi(fmt.Sprintf("%v", LO[0]))
		if nil != err {
			return LO, systemStack, err
		}
		intRO, err := strconv.Atoi(fmt.Sprintf("%v", RO[0]))
		if nil != err {
			return LO, systemStack, err
		}
		return []interface{}{intLO / intRO}, systemStack, nil
	} else if "^" == OP {
		if ("float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: ^ : error: data type mismatch")
			return LO, systemStack, err
		}

		floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)
		if nil != err {
			return LO, systemStack, err
		}

		floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
		if nil != err {
			return LO, systemStack, err
		}

		res := fmt.Sprintf("%.6f", math.Pow(floatLO, floatRO))

		return []interface{}{res}, systemStack, nil
	} else if "str" == OP {
		return []interface{}{"\"" + fmt.Sprintf("%v", LO[0]) + "\""}, systemStack, nil
	} else if "=" == OP {
		return []interface{}{0}, systemStack, nil // успех
	} else if "." == OP {
		return []interface{}{0}, systemStack, nil
	} else if "UNDEFINE" == OP {
		return []interface{}{0}, systemStack, nil
	} else if "int" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)
			if nil != err {
				return LO, systemStack, err
			}

			return []interface{}{int(floatLO)}, systemStack, nil
		}

		intLO, err := strconv.Atoi(fmt.Sprintf("%v", LO[0]))

		if nil != err {
			return LO, systemStack, err
		}

		return []interface{}{intLO}, systemStack, nil
	} else if "float" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}

		floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

		return []interface{}{floatLO}, systemStack, err
	} else if "bool" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}

		return LO, systemStack, nil
	} else if "input" == OP {
		var s string
		_, err := fmt.Scan(&s)
		return []interface{}{s}, systemStack, err
	} else if "goto" == OP {
		if "string" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			err := errors.New("executor: goto : error: data type mismatch")
			return LO, systemStack, err
		}
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if "#" != string(fmt.Sprintf("%v", LO[0])[0]) {
			err := errors.New("executor: goto: ERROR: mark must start with \"#\", mark: " + fmt.Sprintf("%v", LO[0]))
			return LO, systemStack, err
		}

		return []interface{}{"goto", LO[0]}, systemStack, nil
	} else if "exit" == OP {
		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			err := errors.New("executor: exit: error: data type mismatch")
			return LO, systemStack, err
		}
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		code, err := strconv.Atoi(fmt.Sprintf("%v", LO[0]))
		if nil != err {
			panic(err)
		}

		os.Exit(code)
	} else if "is_letter" == OP {
		if "string" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			err := errors.New("executor: is_letter : error: data type mismatch")
			return LO, systemStack, err
		}
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if 1 != len(fmt.Sprintf("%v", LO[0])) {
			err := errors.New("executor: is_letter : error: length of the argument is more than 1, argument: " +
				fmt.Sprintf("%v", LO[0]))

			return LO, systemStack, err
		}
		return []interface{}{unicode.IsLetter([]rune(fmt.Sprintf("%v", LO[0]))[0])}, systemStack, nil
	} else if "is_digit" == OP {
		if "string" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			err := errors.New("executor: is_digit : error: data type mismatch")
			return LO, systemStack, err
		}
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if 1 != len(fmt.Sprintf("%v", LO[0])) {
			err := errors.New("executor: is_digit : error: length of the argument is more than 1, argument: " +
				fmt.Sprintf("%v", LO[0]))

			return LO, systemStack, err
		}
		return []interface{}{unicode.IsDigit([]rune(fmt.Sprintf("%v", LO[0]))[0])}, systemStack, nil
	} else if "SET_SOURCE" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		f, err := os.Open(fmt.Sprintf("%v", LO[0]))
		return []interface{}{"SET_SOURCE", f}, systemStack, err
	} else if "SET_DEST" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		f, err := os.Create(fmt.Sprintf("%v", LO[0]))
		return []interface{}{"SET_DEST", f}, systemStack, err
	} else if "SEND_DEST" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		fin, err := os.Open("benv/program.basm")
		if nil != err {
			return []interface{}{1}, systemStack, err
		}

		fout, err := os.Create(fmt.Sprintf("%v", LO[0]))

		if nil != err {
			return []interface{}{1}, systemStack, err
		}
		_, err = io.Copy(fout, fin)

		if nil != err {
			return []interface{}{1}, systemStack, err
		}

		err = fin.Close()
		if nil != err {
			return []interface{}{1}, systemStack, err
		}
		err = fout.Close()
		if nil != err {
			return []interface{}{1}, systemStack, err
		}

		err = os.Remove("benv/program.basm")
		if nil != err {
			return []interface{}{1}, systemStack, err
		}

		return []interface{}{0}, systemStack, nil
	} else if "DEL_DEST" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if Exists(fmt.Sprintf("%v", LO[0])) {
			//err := LO[0].(*os.File).Close()
			//if nil != err{
			//	return []interface{}{1}, systemStack, err
			//}
			err := os.Remove(fmt.Sprintf("%v", LO[0]))
			if nil != err {
				return []interface{}{1}, systemStack, err
			}
		}

		return []interface{}{0}, systemStack, nil
	} else if "REROUTE" == OP {
		return []interface{}{"REROUTE", 0}, systemStack, nil
	} else if "UNSET_SOURCE" == OP {
		return []interface{}{"UNSET_SOURCE", 0}, systemStack, nil
	} else if "UNSET_DEST" == OP {
		return []interface{}{"UNSET_DEST", 0}, systemStack, nil
	} else if "RESET_SOURCE" == OP {
		return []interface{}{"RESET_SOURCE", 0}, systemStack, nil
	} else if "next_command" == OP {
		return []interface{}{"next_command", LO}, systemStack, nil
	} else if "get_root_source" == OP {
		return []interface{}{"get_root_source", LO}, systemStack, nil
	} else if "get_root_dest" == OP {
		return []interface{}{"get_root_dest", LO}, systemStack, nil
	} else if "send_command" == OP {
		return []interface{}{"send_command", LO}, systemStack, nil
	} else if "push" == OP {
		systemStack = append(systemStack, LO)
		return []interface{}{0}, systemStack, nil
	} else if "pop" == OP {
		res := systemStack[len(systemStack)-1]
		if "end" != res {
			systemStack = systemStack[:len(systemStack)-1]
		}
		if "string" == fmt.Sprintf("%T", ValueFoldInterface(res)) {
			return []interface{}{ValueFoldInterface(res)}, systemStack, nil
		} else {
			return ValueFoldInterface(res).([]interface{}), systemStack, nil
		}
	}
	err := errors.New("execute: ERROR: wrong syntax: OP=\"" + OP + "\", " +
		"LO=\"" + fmt.Sprintf("%v", LO) + "\", " + "RO=\"" + fmt.Sprintf("%v", RO) + "\"")
	return []interface{}{1}, systemStack, err

}

func sysExecuteTree(infoList []interface{}, variables [][]interface{}, systemStack []interface{},
	OPPointer int, toTranspile bool, toPrimitive bool, primitiveDest *os.File, transpileDest *os.File) ([]interface{}, [][]interface{}, []interface{}, int) {
	// заканчивает свою работу, когда выполнен первый оператор
	OP := fmt.Sprintf("%v", infoList[OPPointer])

	var LO []interface{}
	var RO []interface{}

	if toPrimitive && 0 == OPPointer {
		var err error
		for i := 0; i < len(infoList); i++ {
			el := infoList[i]
			if "\"" == string(fmt.Sprintf("%v", el)[0]) && "\"" == string(fmt.Sprintf("%v", el)[len(fmt.Sprintf("%v", el))-1]) {
				el = strings.Replace(fmt.Sprintf("%v", el)[1:len(fmt.Sprintf("%v", el))-1], "\"", "\\\"", -1)
				el = "\"" + fmt.Sprintf("%v", el) + "\""
			}
			if i < len(infoList)-1 {
				_, err = primitiveDest.WriteString(fmt.Sprintf("%v", el) + options.BendSep)
			} else {
				_, err = primitiveDest.WriteString(fmt.Sprintf("%v", el))
			}
			if nil != err {
				panic(err)
			}
		}

		_, err = primitiveDest.WriteString(";\n")
		if nil != err {
			panic(err)
		}

	}
	if 1 == len(infoList) && ("UNSET_SOURCE" == infoList[0] || "RESET_SOURCE" == infoList[0] ||
		"UNSET_DEST" == infoList[0] || "REROUTE" == infoList[0]) {
		infoList = append(infoList, "null")
		infoList = append(infoList, "null")
	}

	if "True" == infoList[OPPointer+1] {
		LO = []interface{}{"True"}
		OPPointer += 1
	} else if "False" == infoList[OPPointer+1] {
		LO = []interface{}{"False"}
		OPPointer += 1
	} else if !IsOp(fmt.Sprintf("%v", infoList[OPPointer+1])) {
		LO = []interface{}{infoList[OPPointer+1]}
		OPPointer += 1
	} else {
		// операция
		OPPointer += 1
		LO, variables, _, OPPointer = sysExecuteTree(infoList, variables, systemStack, OPPointer, toTranspile, toPrimitive, primitiveDest, transpileDest)
	}

	if "True" == infoList[OPPointer+1] {
		RO = []interface{}{"True"}
		OPPointer += 1
	} else if "False" == infoList[OPPointer+1] {
		RO = []interface{}{"False"}
		OPPointer += 1
	} else if !IsOp(fmt.Sprintf("%v", infoList[OPPointer+1])) {
		RO = []interface{}{infoList[OPPointer+1]}
		OPPointer += 1
	} else {
		// операция
		OPPointer += 1
		RO, variables, _, OPPointer = sysExecuteTree(infoList, variables, systemStack, OPPointer, toTranspile, toPrimitive, primitiveDest, transpileDest)
	}

	if "UNDEFINE" == OP {
		if !toTranspile {
			i := len(variables) - 1
			newVariable := EachVariable(variables)
			for v := newVariable(); "end" != v[0]; v = newVariable() {
				if fmt.Sprintf("%v", v[1]) == fmt.Sprintf("%v", LO[0]) {
					copy(variables[i:], variables[i+1:])
					variables = variables[:len(variables)-1]
					break
				}
				i -= 1
			}
		}
		if toTranspile {
			_, err := transpileDest.WriteString("undefineVar(\"" + fmt.Sprintf("%v", LO[0]) + "\")\n")
			if nil != err {
				panic(err)
			}
		}

	}

	if "input" == OP && !toPrimitive {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", v[1]) == fmt.Sprintf("%v", LO[0]) {
				var err error
				if "string" != v[0] {
					err = errors.New("executor: input: ERROR: data type mismatch")
					panic(err)
				}
				v[2], systemStack, err = execute(systemStack, OP, LO, RO)
				if nil != err {
					panic(err)
				}
				break
			}
		}
	}

	if "next_command" == OP || "get_root_source" == OP || "get_root_dest" == OP {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", v[1]) == fmt.Sprintf("%v", LO[0]) {

				v[2] = v[1]
				break
			}
		}
	}

	if "pop" == OP && !toPrimitive {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", v[1]) == fmt.Sprintf("%v", LO[0]) {

				var err error
				if !toTranspile {
					v[2], systemStack, err = execute(systemStack, OP, LO, RO)

					if nil != err {
						panic(err)
					}

					tempV2 := fmt.Sprintf("%v", ValueFoldInterface(v[2]))
					if "\"" == string(tempV2[0]) && "\"" == string(tempV2[len(tempV2)-1]) {
						tempV2 = tempV2[1 : len(tempV2)-1]
					}

					if fmt.Sprintf("%v", v[0]) != WhatsType(tempV2) && "string" != fmt.Sprintf("%v", v[0]) &&
						(fmt.Sprintf("%v", v[0]) != "float" && WhatsType(tempV2) != "int") &&
						!("stack" == fmt.Sprintf("%v", v[0]) && "[]interface {}" == fmt.Sprintf("%T", v[2]) &&
							"end" == ValueFoldInterface(v[2].([]interface{})[0])) {
						panic("pop: data type mismatch: " + fmt.Sprintf("%v", v[0]) + " and " +
							WhatsType(tempV2))
					}
				} else {
					_, err = transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", LO[0]) + "\", " + "systemStack[len(systemStack)-1])\n" +
						"if \"end\" != systemStack[len(systemStack)-1] {\n" +
						"systemStack = systemStack[:len(systemStack)-1]\n}\n")
					if nil != err {
						panic(err)
					}
				}
				break
			}
		}
	}
	var transpileVar interface{}
	var transpileVarT interface{}
	var stranspileVar string

	if "=" == OP && !toPrimitive {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", LO[0]) == fmt.Sprintf("%v", v[1]) {
				typeLO := fmt.Sprintf("%v", v[0])
				var typeRO string
				newRightVar := EachVariable(variables)
				for rightVar := newRightVar(); "end" != rightVar[0]; rightVar = newRightVar() {
					if fmt.Sprintf("%v", RO[0]) == fmt.Sprintf("%v", rightVar[1]) {
						if "[]interface {}" == fmt.Sprintf("%T", rightVar[2]) &&
							"string" == fmt.Sprintf("%T", rightVar[2].([]interface{})[0]) &&
							"end" == fmt.Sprintf("%v", rightVar[2].([]interface{})[0]) {
							rightVar[2].([]interface{})[0] = []interface{}{"end"}
						}
						RO[0] = ValueFoldInterface(rightVar[2])
						transpileVar = rightVar[1]
						typeRO = fmt.Sprintf("%v", rightVar[0])
						break
					}
				}
				if "" == typeRO {
					typeRO = WhatsType(fmt.Sprintf("%v", RO[0]))
				}
				if typeLO == typeRO || ("float" == typeLO && "int" == typeRO) || toTranspile {
					v[2] = RO
					if toTranspile {
						if nil == transpileVar {
							transpileVar = RO[0]
							stranspileVar = fmt.Sprintf("%v", transpileVar)

							if len(stranspileVar) > 7 && "\"getVar" == stranspileVar[0:7] {
								stranspileVar = stranspileVar[1 : len(stranspileVar)-1]
							} else if "\"" != string(stranspileVar[0]) &&
								!(len(stranspileVar) > 6 && "string" == stranspileVar[0:6]) &&
								!(len(stranspileVar) > 6 && "getVar" == stranspileVar[0:6]) &&
								!(len(stranspileVar) > 8 && "unicode." == stranspileVar[0:8]) &&
								!(len(stranspileVar) > 6 && "Exists" == stranspileVar[0:6]) {
								stranspileVar = "\"" + stranspileVar + "\""
							} else if "\"" == string(stranspileVar[0]) {
								stranspileVar = "\"" + stranspileVar[1:len(stranspileVar)-1] + "\""
							}

							if len(stranspileVar) > 4 && ("\"sum" == stranspileVar[0:4] || "\"len" == stranspileVar[0:4]) {
								stranspileVar = stranspileVar[1 : len(stranspileVar)-1]
							}

							if len(stranspileVar) > 7 && ("\"Exists" == stranspileVar[0:7]) {
								stranspileVar = stranspileVar[1 : len(stranspileVar)-1]
							}

							if len(stranspileVar) > 8 && ("\"toInt" == stranspileVar[0:6] || "\"toFloat" == stranspileVar[0:8]) {
								stranspileVar = stranspileVar[1 : len(stranspileVar)-1]
							}
							if len(stranspileVar) > 8 && ("\"toBool" == stranspileVar[0:7] || "\"!toBool" == stranspileVar[0:8]) {
								stranspileVar = stranspileVar[1 : len(stranspileVar)-1]
							}

							if -1 == strings.Index(stranspileVar, "getVar") && -1 == strings.Index(stranspileVar, "Exists") {
								stranspileVar = string(stranspileVar[0]) +
									strings.Replace(stranspileVar[1:len(stranspileVar)-1], "\"", `\"`, -1) +
									string(stranspileVar[len(stranspileVar)-1])
							}

							/*if "\"" == string(stranspileVar[0]) {
								stranspileVar = "`" + stranspileVar[1:len(stranspileVar)-1] + "`"
							}*/
							//re, err := regexp.Compile(`(?:\".*?\\[^\"n].*?\")`)
							re, err := regexp.Compile(`\"(\\.|[^\"])*\"`)

							if nil != err {
								panic(err)
							}

							locs := re.FindAllIndex([]byte(stranspileVar), -1)

							for _, loc := range locs {
								re, err := regexp.Compile(`\\[^\"nt]`)
								if nil != err {
									panic(err)
								}
								if nil != re.FindIndex([]byte(stranspileVar[loc[0]:loc[1]])) {
									stranspileVar = stranspileVar[:loc[0]] + "`" +
										stranspileVar[loc[0]+1:loc[1]-1] + "`" + stranspileVar[loc[1]:]
								}

							}

							_, err = transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", LO[0]) +
								"\"," + stranspileVar + ")\n")
							if nil != err {
								panic(err)
							}
							transpileVar = nil
						} else {
							stranspileVar = fmt.Sprintf("%v", transpileVar)
							if len(stranspileVar) > 7 && "\"getVar" == string(stranspileVar[0:7]) {
								stranspileVar = stranspileVar[1 : len(stranspileVar)-1]
							}

							_, err := transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", LO[0]) + "\", getVar(\"" +
								stranspileVar + "\"))\n")
							if nil != err {
								panic(err)
							}
						}
					}
					break
				} else {
					err := errors.New("executor: ERROR: data type mismatch in assignment")
					panic(err)
				}
			}
		}

	}

	if "." == OP {
		newVariable := EachVariable(variables)
		breakFlag := false
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if breakFlag {
				break
			}
			if fmt.Sprintf("%v", LO[0]) == fmt.Sprintf("%v", v[1]) {
				typeLO := fmt.Sprintf("%v", v[0])

				if "stack" != typeLO {
					err := errors.New("executor: ERROR: the left operand of the operation \".\" must be of type stack")
					panic(err)
				}

				if "pop" != fmt.Sprintf("%v", RO[0].([]string)[0]) {
					newRightVar := EachVariable(variables)
					for rightVar := newRightVar(); "end" != fmt.Sprintf("%v", rightVar[0]); rightVar = newRightVar() {
						if fmt.Sprintf("%v", RO[0].([]string)[1]) == fmt.Sprintf("%v", rightVar[1]) {

							if "string" == fmt.Sprintf("%T", ValueFoldInterface(rightVar[2])) &&
								"stack" != fmt.Sprintf("%v", rightVar[0]) {
								if !toTranspile {
									RO[0].([]string)[1] = fmt.Sprintf("%v", rightVar[2].([]interface{})[0])
								} else {
									RO[0].([]string)[1] = fmt.Sprintf("%v", rightVar[2])
								}
								transpileVar = rightVar[1]
								transpileVarT = rightVar[0]
							} else {
								RO = nil
								RO = append(RO, []string{"push"})
								RO = append(RO, rightVar)
								transpileVar = rightVar[1]
								transpileVarT = rightVar[0]
							}
							break

						}
					}
				}

				if "push" == fmt.Sprintf("%v", RO[0].([]string)[0]) {

					if 1 == len(RO) {
						v[2] = append(v[2].([]interface{}), []interface{}{RO[0].([]string)[1]})

						if toTranspile && nil != transpileVar {
							_, err := transpileDest.WriteString("setVar(\"" +
								fmt.Sprintf("%v", v[1]) + "\", append(getVar(\"" +
								fmt.Sprintf("%v", v[1]) + "\").([]interface{}), " + "getVar(\"" + fmt.Sprintf("%v", transpileVar) + "\")))\n")
							if nil != err {
								panic(err)
							}
						} else if toTranspile && nil == transpileVar {
							if `"` != string(RO[0].([]string)[1][0]) {
								RO[0].([]string)[1] = `"` + RO[0].([]string)[1] + `"`
							}
							_, err := transpileDest.WriteString("setVar(\"" +
								fmt.Sprintf("%v", v[1]) + "\", append(getVar(\"" +
								fmt.Sprintf("%v", v[1]) + "\").([]interface{}), " + RO[0].([]string)[1] + "))\n")
							if nil != err {
								panic(err)
							}
						}
					} else {
						v[2] = append(v[2].([]interface{}), RO[1].([]interface{})[2])

						if toTranspile {
							if "stack" == transpileVarT {

								goCommand := "if (\"end\" == fmt.Sprintf(\"%v\", getVar(\"" + fmt.Sprintf("%v", v[1]) + "\").([]interface{})[0])){\n"
								_, err := transpileDest.WriteString(goCommand)
								if nil != err {
									panic(err)
								}

								goCommand = "setVar(\"" + fmt.Sprintf("%v", v[1]) +
									"\", append([]interface{}{[]interface{}{\"end\"}}, getVar(\"" + fmt.Sprintf("%v", v[1]) + "\").([]interface{})[1:]...))" + "\n}\n"
								_, err = transpileDest.WriteString(goCommand)
								if nil != err {
									panic(err)
								}
							}
							goCommand := "setVar(\"" +
								fmt.Sprintf("%v", v[1]) + "\", append(getVar(\"" +
								fmt.Sprintf("%v", v[1]) + "\").([]interface{}), " + "getVar(\"" + fmt.Sprintf("%v", transpileVar) + "\")))\n"
							_, err := transpileDest.WriteString(goCommand)
							if nil != err {
								panic(err)
							}
						}
					}
					break
				} else if "pop" == fmt.Sprintf("%v", RO[0].([]string)[0]) {
					newPopVariable := EachVariable(variables)
					for popVar := newPopVariable(); "end" != fmt.Sprintf("%v", popVar[0]); popVar = newPopVariable() {
						if popVar[1] == RO[0].([]string)[1] {

							if "end" == fmt.Sprintf("%v", v[2].([]interface{})[0]) {
								v[2].([]interface{})[0] = []interface{}{"end"}
							}
							VFI := ValueFoldInterface(v[2])
							if "string" == fmt.Sprintf("%T", VFI) {
								VFI = []interface{}{VFI}
							}
							popVar[2] = VFI.([]interface{})[len(VFI.([]interface{}))-1]

							tempPopVar2 := fmt.Sprintf("%v", ValueFoldInterface(popVar[2]))
							if "\"" == string(tempPopVar2[0]) && "\"" == string(tempPopVar2[len(tempPopVar2)-1]) {
								tempPopVar2 = tempPopVar2[1 : len(tempPopVar2)-1]
							}

							if fmt.Sprintf("%v", popVar[0]) != WhatsType(tempPopVar2) && "string" != fmt.Sprintf("%v", popVar[0]) &&
								(fmt.Sprintf("%v", popVar[0]) != "float" && WhatsType(tempPopVar2) != "int") && !("stack" == fmt.Sprintf("%v", popVar[0]) &&
								"[]interface {}" == fmt.Sprintf("%T", popVar[2]) &&
								"end" == ValueFoldInterface(popVar[2].([]interface{})[0])) {
								if !("stack" == fmt.Sprintf("%v", popVar[0]) && "string" == WhatsType(fmt.Sprintf("%v", tempPopVar2)) &&
									"end" == tempPopVar2) && !toTranspile {
									panic("pop: data type mismatch: " + fmt.Sprintf("%v", popVar[0]) + " and " +
										WhatsType(fmt.Sprintf("%v", tempPopVar2)))
								} else {
									VFI = []interface{}{[]interface{}{"end"}}
									popVar[2] = []interface{}{[]interface{}{"end"}}

								}
							}

							v[2] = VFI

							if len(VFI.([]interface{})) > 1 {
								v[2] = v[2].([]interface{})[:len(v[2].([]interface{}))-1]
							}
							breakFlag = true
							if toTranspile {
								_, err := transpileDest.WriteString("if \"[]interface {}\" == fmt.Sprintf(\"%T\", getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\")) && len(getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\").([]interface{})) > 1{\n")
								if nil != err {
									panic(err)
								}
								_, err = transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", RO[0].([]string)[1]) + "\", getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\").([]interface{})[len(getVar(\"" + fmt.Sprintf("%v", v[1]) +
									"\").([]interface{})) - 1])\n}")
								if nil != err {
									panic(err)
								}

								_, err = transpileDest.WriteString("else if fmt.Sprintf(\"%T\", getVar(\"" +
									fmt.Sprintf("%v", RO[0].([]string)[1]) + "\")) == \"[]interface {}\"{\nsetVar(\"" +
									fmt.Sprintf("%v", RO[0].([]string)[1]) + "\", []interface{}{[]interface{}{\"end\"}})\n}else{\n")
								if nil != err {
									panic(err)
								}

								_, err = transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", RO[0].([]string)[1]) + "\", \"end\")\n}\n")

								if nil != err {
									panic(err)
								}

								_, err = transpileDest.WriteString("if \"[]interface {}\" == fmt.Sprintf(\"%T\", getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\")) && !isEqual(\"end\", " + "getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\").([]interface{})[len(getVar(\"" + fmt.Sprintf("%v", v[1]) +
									"\").([]interface{})) - 1]) && " + "!isEqual(\"[end]\", " + "getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\").([]interface{})[len(getVar(\"" + fmt.Sprintf("%v", v[1]) +
									"\").([]interface{})) - 1]){\n")
								if nil != err {
									panic(err)
								}
								_, err = transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", v[1]) + "\", getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\").([]interface{})[:len(getVar(\"" + fmt.Sprintf("%v", v[1]) +
									"\").([]interface{})) - 1])\n}\n")
								if nil != err {
									panic(err)
								}
							}

							break

						}
					}
				} else {
					err := errors.New("executor: ERROR: after the operation \".\" must follow \"push\" or \"pop\"")
					panic(err)
				}
			}
		}
	}

	if "goto" == OP {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != fmt.Sprintf("%v", v[0]); v = newVariable() {
			if fmt.Sprintf("%v", LO[0]) == fmt.Sprintf("%v", v[1]) {
				if !toTranspile {
					LO = v[2].([]interface{})
				} else {
					LO[0] = "#getVar(\"" + fmt.Sprintf("%v", v[1]) + "\")"
				}
			}

		}
	} else {
		newVariable := EachVariable(variables)
		wasLO := false
		wasRO := false

		for v := newVariable(); "end" != fmt.Sprintf("%v", v[0]); v = newVariable() {
			if fmt.Sprintf("%v", LO[0]) == fmt.Sprintf("%v", v[1]) {
				if "print" == OP && "string" != v[0] {
					panic(errors.New("sysExecuteTree: ERROR: print: dataTypeMismatch"))
				}
				if "len" == OP && "string" != v[0] {
					panic(errors.New("sysExecuteTree: ERROR: len: dataTypeMismatch"))
				}

				if "string" == fmt.Sprintf("%T", v[2]) && !wasLO {
					if !toTranspile {
						LO = []interface{}{v[2]}
					} else {
						LO = []interface{}{"getVar(\"" + fmt.Sprintf("%v", v[1]) + "\")"}
					}
					wasLO = true
				} else if !wasLO {
					if !toTranspile {
						LO = v[2].([]interface{})
					} else {
						LO = []interface{}{"getVar(\"" + fmt.Sprintf("%v", v[1]) + "\")"}
					}
					wasLO = true
				}
			}
			if fmt.Sprintf("%v", RO[0]) == fmt.Sprintf("%v", v[1]) {
				if "string" == fmt.Sprintf("%T", v[2]) && !wasRO {
					if !toTranspile {
						RO = []interface{}{v[2]}
					} else {
						RO = []interface{}{"getVar(\"" + fmt.Sprintf("%v", v[1]) + "\")"}
					}
					wasRO = true

				} else if !wasRO {
					if !toTranspile {
						RO = v[2].([]interface{})
					} else {
						RO = []interface{}{"getVar(\"" + fmt.Sprintf("%v", v[1]) + "\")"}
					}
					wasRO = true
				}
			}
		}
	}

	var res []interface{}
	if "input" != OP && "pop" != OP {
		var err error
		var passLO []interface{}
		var passRO []interface{}

		for _, el := range LO {
			passLO = append(passLO, el)
		}
		for _, el := range RO {
			passRO = append(passRO, el)
		}
		if !toTranspile && !toPrimitive {
			res, systemStack, err = execute(systemStack, OP, passLO, passRO)
			if nil != err {
				panic(err)
			}
		} else if toTranspile {
			res, systemStack, err = Transpile(systemStack, OP, passLO, passRO, transpileDest)
			if nil != err {
				panic(err)
			}
		} else if toPrimitive {
			res = []interface{}{0}
		}
	} else {
		res = []interface{}{0}
	}

	return res, variables, systemStack, OPPointer
}

func ExecuteTree(infoList []interface{}, variables [][]interface{},
	systemStack []interface{}, toTranspile bool, toPrimitive bool, primitiveDest *os.File, transpileDest *os.File) ([]interface{}, [][]interface{}, []interface{}) {
	res, variables, systemStack, _ := sysExecuteTree(infoList, variables, systemStack, 0, toTranspile, toPrimitive, primitiveDest, transpileDest)
	return res, variables, systemStack
}
