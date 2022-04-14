package transpiler

import (
	. "bint.com/pkg/serviceTools"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func Transpile(systemStack []interface{}, OP string, LO []interface{}, RO []interface{},
	transpileDest *os.File) ([]interface{}, []interface{}, error) {
	if "print" == OP {
		return []interface{}{"print", LO}, systemStack, nil
	} else if "index" == OP {

		return []interface{}{"strings.Index(fmt.Sprintf(\"%v\"," + fmt.Sprintf("%v", LO[0]) + "), fmt.Sprintf(\"%v\"," + fmt.Sprintf("%v", RO[0]) + "))"}, systemStack, nil
	} else if "len" == OP {

		return []interface{}{"len(" + "fmt.Sprintf(\"%v\"," + fmt.Sprintf("%v", LO[0]) + "))"}, systemStack, nil
	} else if "CD" == OP {
		if "goto" == fmt.Sprintf("%v", LO[0]) && "goto" == fmt.Sprintf("%v", RO[0]) {
			_, err := transpileDest.WriteString(fmt.Sprintf("%v", LO[0]) + " " + fmt.Sprintf("%v", LO[1].(string)[1:]) + "\n" + "}else{\n" +
				fmt.Sprintf("%v", RO[0]) + " " + fmt.Sprintf("%v", RO[1].(string)[1:]) + "\n}\n")
			return []interface{}{0}, systemStack, err
		} else if "print" == fmt.Sprintf("%v", LO[0]) && "print" == fmt.Sprintf("%v", RO[0]) {
			_, err := transpileDest.WriteString(fmt.Sprintf("%v", LO[0]) + "(" + fmt.Sprintf("%v", LO[1].([]interface{})[0]) + ")\n" + "}else{\n" +
				fmt.Sprintf("%v", RO[0]) + "(" + fmt.Sprintf("%v", RO[1].([]interface{})[0]) + ")" + "\n}\n")
			return []interface{}{0}, systemStack, err
		} else if "print" == fmt.Sprintf("%v", LO[0]) && "goto" == fmt.Sprintf("%v", RO[0]) {
			_, err := transpileDest.WriteString(fmt.Sprintf("%v", LO[0]) + "(" + fmt.Sprintf("%v", LO[1].([]interface{})[0]) + ")\n" + "}else{\n" +
				fmt.Sprintf("%v", RO[0]) + " " + fmt.Sprintf("%v", RO[1].(string)[1:]) + "\n}\n")
			return []interface{}{0}, systemStack, err
		} else if "goto" == fmt.Sprintf("%v", LO[0]) && "print" == fmt.Sprintf("%v", RO[0]) {
			_, err := transpileDest.WriteString(fmt.Sprintf("%v", LO[0]) + " " + fmt.Sprintf("%v", LO[1].(string)[1:]) + "\n" + "}else{\n" +
				fmt.Sprintf("%v", RO[0]) + "(" + fmt.Sprintf("%v", RO[1].([]interface{})[0]) + ")" + "\n}\n")
			return []interface{}{0}, systemStack, err
		} else {
			return []interface{}{-1}, systemStack,
				errors.New("could not recognize operands of conditional disjunction; left operand: " + fmt.Sprintf("%v", LO[0]) + ", " +
					"right operand:" + fmt.Sprintf("%v", RO[0]))
		}
	} else if "AND" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = StrToBool(fmt.Sprintf("%v", LO[0]))
		} else {
			LO[0] = "toBool(" + fmt.Sprintf("%v", LO[0]) + ")"
		}
		if "bool" == WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = StrToBool(fmt.Sprintf("%v", RO[0]))
		} else {
			RO[0] = "toBool(" + fmt.Sprintf("%v", RO[0]) + ")"
		}

		return []interface{}{fmt.Sprintf("%v", fmt.Sprintf("%v", LO[0])) + "&&" + fmt.Sprintf("%v", fmt.Sprintf("%v", RO[0]))}, systemStack, nil
	} else if "OR" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = StrToBool(fmt.Sprintf("%v", LO[0]))
		} else {
			LO[0] = "toBool(" + fmt.Sprintf("%v", LO[0]) + ")"
		}
		if "bool" == WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = StrToBool(fmt.Sprintf("%v", RO[0]))
		} else {
			RO[0] = "toBool(" + fmt.Sprintf("%v", RO[0]) + ")"
		}

		return []interface{}{fmt.Sprintf("%v", fmt.Sprintf("%v", LO[0])) + "||" + fmt.Sprintf("%v", fmt.Sprintf("%v", RO[0]))}, systemStack, nil

	} else if "XOR" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = StrToBool(fmt.Sprintf("%v", LO[0]))
		} else {
			LO[0] = "toBool(" + fmt.Sprintf("%v", LO[0]) + ")"
		}
		if "bool" == WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = StrToBool(fmt.Sprintf("%v", RO[0]))
		} else {
			RO[0] = "toBool(" + fmt.Sprintf("%v", RO[0]) + ")"
		}

		return []interface{}{fmt.Sprintf("%v", fmt.Sprintf("%v", LO[0])) + "!=" + fmt.Sprintf("%v", fmt.Sprintf("%v", RO[0]))}, systemStack, nil
	} else if "NOT" == OP {
		if "bool" == WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = StrToBool(fmt.Sprintf("%v", LO[0]))
		} else {
			LO[0] = "toBool(" + fmt.Sprintf("%v", LO[0]) + ")"
		}

		return []interface{}{"!" + fmt.Sprintf("%v", fmt.Sprintf("%v", LO[0]))}, systemStack, nil
	} else if "<" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{fmt.Sprintf("%v", LO[0]) + "<" + fmt.Sprintf("%v", RO[0])}, systemStack, nil
	} else if "<=" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{fmt.Sprintf("%v", LO[0]) + "<=" + fmt.Sprintf("%v", RO[0])}, systemStack, nil

	} else if "==" == OP {

		return []interface{}{"isEqual(" + fmt.Sprintf("%v", LO[0]) + ", " + fmt.Sprintf("%v", RO[0]) + ")"}, systemStack, nil
	} else if ">" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{fmt.Sprintf("%v", LO[0]) + ">" + fmt.Sprintf("%v", RO[0])}, systemStack, nil
	} else if ">=" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{fmt.Sprintf("%v", LO[0]) + ">=" + fmt.Sprintf("%v", RO[0])}, systemStack, nil
	} else if "+" == OP {

		return []interface{}{"sum(" + fmt.Sprintf("%v", LO[0]) + ", " + fmt.Sprintf("%v", RO[0]) + ")"}, systemStack, nil
	} else if "-" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{fmt.Sprintf("%v", LO[0]) + "-" + fmt.Sprintf("%v", RO[0])}, systemStack, nil
	} else if "*" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{fmt.Sprintf("%v", LO[0]) + "*" + fmt.Sprintf("%v", RO[0])}, systemStack, nil
	} else if "/" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{fmt.Sprintf("%v", LO[0]) + "/" + fmt.Sprintf("%v", RO[0])}, systemStack, nil
	} else if "^" == OP {
		var floatLO float64
		var floatRO float64
		var err error

		if "int" != WhatsType(fmt.Sprintf("%v", LO[0])) && "float" != WhatsType(fmt.Sprintf("%v", LO[0])) {
			LO[0] = "toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"
		} else {
			floatLO, err = strconv.ParseFloat(fmt.Sprintf("%v", LO[0]), 64)

			if nil != err {
				return LO, systemStack, err
			}

			LO[0] = floatLO
		}
		if "int" != WhatsType(fmt.Sprintf("%v", RO[0])) && "float" != WhatsType(fmt.Sprintf("%v", RO[0])) {
			RO[0] = "toFloat(" + fmt.Sprintf("%v", RO[0]) + ")"
		} else {
			floatRO, err = strconv.ParseFloat(fmt.Sprintf("%v", RO[0]), 64)

			if nil != err {
				return RO, systemStack, err
			}

			RO[0] = floatRO
		}

		return []interface{}{"math.Pow(" + fmt.Sprintf("%v", LO[0]) + ", " + fmt.Sprintf("%v", RO[0]) + ")"}, systemStack, nil
	} else if "str" == OP {
		return []interface{}{fmt.Sprintf("%v", LO[0])}, systemStack, nil
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

		return []interface{}{"toFloat(" + fmt.Sprintf("%v", LO[0]) + ")"}, systemStack, nil
	} else if "bool" == OP {
		if "\"" == string(fmt.Sprintf("%v", LO[0])[0]) {
			LO[0] = LO[0].(string)[1 : len(LO[0].(string))-1]
		}

		return []interface{}{"toBool(" + fmt.Sprintf("%v", LO[0]) + ")"}, systemStack, nil
	} else if "input" == OP {

		return []interface{}{-1}, systemStack, errors.New("can not transpile input")
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

		return []interface{}{"unicode.IsLetter([]rune(fmt.Sprintf(\"%v\"," + fmt.Sprintf("%v", LO[0]) + "))[0])"}, systemStack, nil
	} else if "is_digit" == OP {

		return []interface{}{"unicode.IsDigit([]rune(fmt.Sprintf(\"%v\"," + fmt.Sprintf("%v", LO[0]) + "))[0])"}, systemStack, nil
	} else if "SET_SOURCE" == OP {
		_, err := transpileDest.WriteString("defineVar(\"$SOURCE\")\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}

		_, err = transpileDest.WriteString("setVar(\"$SOURCE\", openFile(getRootSource(" + fmt.Sprintf("%v", LO[0]) + ")))\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}

		_, err = transpileDest.WriteString("defineVar(\"sourceNewChunk\")\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}

		_, err = transpileDest.WriteString("setVar(\"sourceNewChunk\", EachChunk(getVar(\"$SOURCE\").(*os.File)))\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		return []interface{}{0}, systemStack, nil
	} else if "SET_DEST" == OP {
		_, err := transpileDest.WriteString("defineVar(\"$DEST\")\n")
		if nil != err {
			panic(err)
		}

		_, err = transpileDest.WriteString("setVar(\"$DEST\", createFile(getRootSource(" + fmt.Sprintf("%v", LO[0]) + ")))\n")
		if nil != err {
			panic(err)
		}
		return []interface{}{0}, systemStack, nil
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
		_, err := transpileDest.WriteString("getVar(\"$SOURCE\").(*os.File).Close()\n")
		if nil != err {
			return []interface{}{-1}, systemStack, nil
		}

		return []interface{}{0}, systemStack, nil
	} else if "UNSET_DEST" == OP {
		_, err := transpileDest.WriteString("getVar(\"$DEST\").(*os.File).Close()\n")
		if nil != err {
			return []interface{}{-1}, systemStack, nil
		}

		return []interface{}{0}, systemStack, nil
	} else if "RESET_SOURCE" == OP {
		return []interface{}{"RESET_SOURCE", 0}, systemStack, nil
	} else if "next_command" == OP {
		_, err := transpileDest.WriteString("defineVar(\"$CODE\")\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		_, err = transpileDest.WriteString("setVar(\"$CODE\", fmt.Sprintf(\"%v\", getVar(\"sourceNewChunk\").(func () string)()))\n")

		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		_, err = transpileDest.WriteString("setVar(" +
			fmt.Sprintf("%v", LO[0])[7:len(fmt.Sprintf("%v", LO[0]))-1] + ", " + "getVar(\"$CODE\"))\n")
		if nil != err {
			panic(err)
		}

		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		return []interface{}{0}, systemStack, nil
	} else if "get_root_source" == OP {
		_, err := transpileDest.WriteString("setVar(" + fmt.Sprintf("%v", LO[0])[7:len(fmt.Sprintf("%v", LO[0]))-1] + ", iFlag)\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		return []interface{}{0}, systemStack, nil
	} else if "get_root_dest" == OP {
		_, err := transpileDest.WriteString("setVar(" + fmt.Sprintf("%v", LO[0])[7:len(fmt.Sprintf("%v", LO[0]))-1] + ", oFlag)\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		return []interface{}{0}, systemStack, nil
	} else if "send_command" == OP {
		_, err := transpileDest.WriteString("getVar(\"$DEST\").(*os.File).WriteString(" + fmt.Sprintf("%v", LO[0]) + ".(string) + \";\")\n")
		if nil != err {
			return []interface{}{-1}, systemStack, err
		}
		return []interface{}{0}, systemStack, nil
	} else if "push" == OP {
		systemStack = append(systemStack, LO)

		_, err := transpileDest.WriteString("systemStack = append(systemStack, " + fmt.Sprintf("%v", LO[0]) + ")\n")

		return []interface{}{0}, systemStack, err
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
