package executor

import (
	. "bint.com/internal/transpiler"
	. "bint.com/pkg/serviceTools"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func execute(systemStack []interface{}, OP string, LO []interface{}, RO []interface{}) ([]interface{}, []interface{}, error) {
	if "print" == OP {
		return []interface{}{"print", LO}, systemStack, nil
	} else if "index" == OP {

		if `"` == string(fmt.Sprintf("%v", LO[0])[0]) && `"` == string(fmt.Sprintf("%v", LO[0])[len(fmt.Sprintf("%v", LO[0]))-1]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}
		if `"` == string(fmt.Sprintf("%v", RO[0])[0]) && `"` == string(fmt.Sprintf("%v", RO[0])[len(fmt.Sprintf("%v", RO[0]))-1]) {
			RO[0] = RO[0].(string)[1 : len(RO[0].(string))-1]
		}

		return []interface{}{strings.Index(fmt.Sprintf("%v", LO[0]), fmt.Sprintf("%v", RO[0]))}, systemStack, nil
	} else if "len" == OP {
		if `"` == string(fmt.Sprintf("%v", LO[0])[0]) && `"` == string(fmt.Sprintf("%v", LO[0])[len(fmt.Sprintf("%v", LO[0]))-1]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}

		return []interface{}{len(fmt.Sprintf("%v", LO[0]))}, systemStack, nil
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
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			if !("string" == WhatsType(fmt.Sprintf("%v", LO[0])) && "string" == WhatsType(fmt.Sprintf("%v", RO[0]))) {
				err := errors.New("executor: + : error: data type mismatch")
				return LO, systemStack, err
			}
		}

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

			err = errors.New("executor: + : error: data type mismatch")
			return LO, systemStack, err

		}

		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{floatLO + floatRO}, systemStack, nil
			}

			err = errors.New("executor: + : error: data type mismatch")
			return LO, systemStack, err
		}

		if "string" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			if "string" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				if len(LO) > 0 && "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
					LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
				}
				if len(RO) > 0 && "\"" == string(fmt.Sprintf("%v", RO[0])[0]) {
					RO[0] = RO[0].(string)[1 : len(RO[0].(string))-1]
				}

				return []interface{}{"\"" + fmt.Sprintf("%v", LO[0]) + fmt.Sprintf("%v", RO[0]) + "\""}, systemStack, nil
			}

			err := errors.New("executor: + : error: data type mismatch")
			return LO, systemStack, err
		}

		err := errors.New("executor: + : error: data type mismatch")
		return LO, systemStack, err
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

			err = errors.New("executor: - : error: data type mismatch")
			return LO, systemStack, err

		}

		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{floatLO - floatRO}, systemStack, nil
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

			err = errors.New("executor: * : error: data type mismatch")
			return LO, systemStack, err

		}

		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{floatLO * floatRO}, systemStack, nil
			}

			err = errors.New("executor: * : error: data type mismatch")
			return LO, systemStack, err
		}
	} else if "/" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
			err := errors.New("executor: / : error: data type mismatch")
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

				return []interface{}{intLO / intRO}, systemStack, nil
			}

			err = errors.New("executor: / : error: data type mismatch")
			return LO, systemStack, err

		}

		if "float" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			floatLO, err := strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			if "float" == WhatsType(fmt.Sprintf("%v", RO[0])) {
				floatRO, err := strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)
				if nil != err {
					return LO, systemStack, err
				}

				return []interface{}{floatLO / floatRO}, systemStack, nil
			}

			err = errors.New("executor: / : error: data type mismatch")
			return LO, systemStack, err
		}
	} else if "^" == OP {
		if ("int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0]))) ||
			("int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0]))) {
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

		return []interface{}{math.Pow(floatLO, floatRO)}, systemStack, nil
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
	OPPointer int, toTranspile bool, transpileDest *os.File) ([]interface{}, [][]interface{}, []interface{}, int) {
	// заканчивает свою работу, когда выполнен первый оператор
	OP := fmt.Sprintf("%v", infoList[OPPointer])
	var LO []interface{}
	var RO []interface{}

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
		LO, variables, _, OPPointer = sysExecuteTree(infoList, variables, systemStack, OPPointer, toTranspile, transpileDest)
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
		RO, variables, _, OPPointer = sysExecuteTree(infoList, variables, systemStack, OPPointer, toTranspile, transpileDest)
	}

	if "UNDEFINE" == OP {
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

		if toTranspile {
			_, err := transpileDest.WriteString("undefineVar(\"" + fmt.Sprintf("%v", LO[0]) + "\")\n")
			if nil != err {
				panic(err)
			}
		}
	}

	if "input" == OP {
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

	if "pop" == OP {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", v[1]) == fmt.Sprintf("%v", LO[0]) {
				var err error
				v[2], systemStack, err = execute(systemStack, OP, LO, RO)
				if nil != err {
					panic(err)
				}
				_, err = transpileDest.WriteString("if \"end\" != systemStack[len(systemStack)-1] {\n" +
					"setVar(" + fmt.Sprintf("%v", LO[0]) + ", " + "systemStack[len(systemStack)-1])\n" +
					"systemStack = systemStack[:len(systemStack)-1]\n}\n")
				if nil != err {
					panic(err)
				}
				break
			}
		}
	}
	var transpileVar interface{}

	if "=" == OP {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", LO[0]) == fmt.Sprintf("%v", v[1]) {
				typeLO := fmt.Sprintf("%v", v[0])
				var typeRO string
				newRightVar := EachVariable(variables)
				for rightVar := newRightVar(); "end" != rightVar[0]; rightVar = newRightVar() {
					if fmt.Sprintf("%v", RO[0]) == fmt.Sprintf("%v", rightVar[1]) {
						RO[0] = ValueFoldInterface(rightVar[2])
						transpileVar = rightVar[1]
						typeRO = fmt.Sprintf("%v", rightVar[0])
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
							stranspileVar := fmt.Sprintf("%v", transpileVar)
							if len(stranspileVar) > 7 && "\"getVar" == string(stranspileVar[0:7]) {
								stranspileVar = stranspileVar[1 : len(stranspileVar)-1]
							}
							_, err := transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", LO[0]) +
								"\"," + stranspileVar + ")\n")
							if nil != err {
								panic(err)
							}
						} else {
							stranspileVar := fmt.Sprintf("%v", transpileVar)
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

							if "string" == fmt.Sprintf("%T", ValueFoldInterface(rightVar[2])) {
								RO[0].([]string)[1] = fmt.Sprintf("%v", rightVar[2].([]interface{})[0])
								transpileVar = rightVar[1]
							} else {
								RO = nil
								RO = append(RO, []string{"push"})
								RO = append(RO, rightVar)
								transpileVar = rightVar[1]
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
						} else if nil == transpileVar {
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
							v[2] = VFI

							if "end" != fmt.Sprintf("%v", ValueFoldInterface(VFI.([]interface{})[len(VFI.([]interface{}))-1])) {
								v[2] = v[2].([]interface{})[:len(v[2].([]interface{}))-1]
							}
							breakFlag = true
							if toTranspile {
								_, err := transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", RO[0].([]string)[1]) + "\", getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\").([]interface{})[len(getVar(\"" + fmt.Sprintf("%v", v[1]) +
									"\").([]interface{})) - 1])\n")
								if nil != err {
									panic(err)
								}

								_, err = transpileDest.WriteString("setVar(\"" + fmt.Sprintf("%v", v[1]) + "\", getVar(\"" +
									fmt.Sprintf("%v", v[1]) + "\").([]interface{})[:len(getVar(\"" + fmt.Sprintf("%v", v[1]) +
									"\").([]interface{})) - 1])\n")
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
				LO = v[2].([]interface{})
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
		if !toTranspile {
			res, systemStack, err = execute(systemStack, OP, passLO, passRO)
			if nil != err {
				panic(err)
			}
		} else {
			res, systemStack, err = Transpile(systemStack, OP, passLO, passRO, transpileDest)
			if nil != err {
				panic(err)
			}
		}
	} else {
		res = []interface{}{0}
	}

	return res, variables, systemStack, OPPointer
}

func ExecuteTree(infoList []interface{}, variables [][]interface{},
	systemStack []interface{}, toTranspile bool, transpileDest *os.File) ([]interface{}, [][]interface{}, []interface{}) {
	res, variables, systemStack, _ := sysExecuteTree(infoList, variables, systemStack, 0, toTranspile, transpileDest)
	return res, variables, systemStack
}
