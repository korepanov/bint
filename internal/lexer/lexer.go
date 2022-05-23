package lexer

import (
	"bint.com/internal/options"
	. "bint.com/pkg/serviceTools"
	"errors"
	"fmt"
	"os"
	"unicode"
)

func LexicalAnalyze(expr string, variables [][]interface{}, toTranspile bool,
	transpileDest *os.File, toPrimitive bool, primitiveDest *os.File) ([][]interface{}, [][]interface{}, error) {

	var res [][]interface{}

	i := 0
	isType := false

	for i < len(expr) {

		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if len(expr) > i+len(v[1].(string)) {
				if v[1].(string) == expr[i:i+len(v[1].(string))] {
					// проверяем, что имя переменной не является частью имени другой переменной
					if i+len(v[1].(string)) < len(expr) && !unicode.IsDigit(rune(expr[i+len(v[1].(string))])) &&
						!unicode.IsLetter(rune(expr[i+len(v[1].(string))])) && "_" != string(expr[i+len(v[1].(string))]) {
						// проверяем, что это не присвоение
						if i+1 < len(expr) && "=" != string(expr[i+len(v[1].(string))]) {
							// проверяем, что найденная переменная не является частью какого-либо другого слова
							if i-1 >= 0 && !unicode.IsDigit(rune(expr[i-1])) &&
								!unicode.IsLetter(rune(expr[i-1])) &&
								!unicode.IsDigit(rune(expr[i+len(v[1].(string))])) &&
								!unicode.IsLetter(rune(expr[i+len(v[1].(string))])) {
								res = append(res, []interface{}{"VAR", v[1]})
								i += len(v[1].(string))
							}
						}
					}
				}
			}
		}

		if len(expr) > i+3 && "AND" == expr[i:i+3] {
			res = append(res, []interface{}{"OP", "AND"})
			i += 2
		} else if len(expr) > i+2 && "OR" == expr[i:i+2] {
			res = append(res, []interface{}{"OP", "OR"})
			i += 1
		} else if len(expr) > i+3 && "NOT" == expr[i:i+3] {
			res = append(res, []interface{}{"OP", "NOT"})
			i += 2
		} else if len(expr) > i+3 && "XOR" == expr[i:i+3] {
			res = append(res, []interface{}{"OP", "XOR"})
			i += 2
		} else if len(expr) > i && "." == string(expr[i]) {
			res = append(res, []interface{}{"OP", "."})
		} else if len(expr) > i && "+" == string(expr[i]) {
			res = append(res, []interface{}{"OP", "+"})
		} else if len(expr) > i+1 && "-" == string(expr[i]) && (unicode.IsDigit(rune(expr[i-1])) ||
			unicode.IsLetter(rune(expr[i-1])) ||
			")" == string(expr[i-1])) && (unicode.IsDigit(rune(expr[i+1])) || unicode.IsLetter(rune(expr[i+1])) ||
			"(" == string(expr[i+1])) {
			res = append(res, []interface{}{"OP", "-"})
		} else if len(expr) > i && "*" == string(expr[i]) {
			res = append(res, []interface{}{"OP", "*"})
		} else if len(expr) > i && "/" == string(expr[i]) {
			res = append(res, []interface{}{"OP", "/"})
		} else if len(expr) > i && "^" == string(expr[i]) {
			res = append(res, []interface{}{"OP", "^"})
		} else if len(expr) > i+6 && "print(" == expr[i:i+6] {
			res = append(res, []interface{}{"OP", "print"})
			i += 4
		} else if len(expr) > i+3 && "len" == expr[i:i+3] {
			res = append(res, []interface{}{"OP", "len"})
			i += 2
		} else if len(expr) > i+6 && "index(" == expr[i:i+6] {
			res = append(res, []interface{}{"OP", "index"})
			i += 4
		} else if len(expr) > i+9 && "is_letter" == expr[i:i+9] {
			res = append(res, []interface{}{"OP", "is_letter"})
			i += 8
		} else if len(expr) > i+8 && "is_digit" == expr[i:i+8] {
			res = append(res, []interface{}{"OP", "is_digit"})
			i += 7
		} else if len(expr) > i+8 && "reg_find" == expr[i:i+8] {
			res = append(res, []interface{}{"OP", "reg_find"})
			i += 7
		} else if len(expr) > i+3 && "pop" == expr[i:i+3] {
			res = append(res, []interface{}{"OP", "pop"})
			i += 2
		} else if len(expr) > i+4 && "push" == expr[i:i+4] {
			res = append(res, []interface{}{"OP", "push"})
			i += 3
		} else if len(expr) > i+5 && "input" == expr[i:i+5] {
			res = append(res, []interface{}{"OP", "input"})
			i += 4
		} else if len(expr) > i+12 && "next_command" == expr[i:i+12] {
			res = append(res, []interface{}{"OP", "next_command"})
			i += 11
		} else if len(expr) > i+15 && "get_root_source" == expr[i:i+15] {
			res = append(res, []interface{}{"OP", "get_root_source"})
			i += 14
		} else if len(expr) > i+13 && "get_root_dest" == expr[i:i+13] {
			res = append(res, []interface{}{"OP", "get_root_dest"})
			i += 12
		} else if len(expr) > i+12 && "send_command" == expr[i:i+12] {
			res = append(res, []interface{}{"OP", "send_command"})
			i += 11
		} else if len(expr) > i+4 && "goto" == expr[i:i+4] {
			res = append(res, []interface{}{"OP", "goto"})
			i += 3
		} else if len(expr) > i && "#" == string(expr[i]) {
			mark := string(expr[i])
			i += 1

			for ":" != string(expr[i]) && ")" != string(expr[i]) {
				mark += string(expr[i])
				i += 1
			}

			if ")" == string(expr[i]) {
				res = append(res, []interface{}{"VAR", mark})
				res = append(res, []interface{}{"BR", ")"})
			}
		} else if len(expr) > i+10 && "SET_SOURCE" == expr[i:i+10] {
			res = append(res, []interface{}{"OP", "SET_SOURCE"})
			i += 9
		} else if len(expr) > i+8 && "SET_DEST" == expr[i:i+8] {
			res = append(res, []interface{}{"OP", "SET_DEST"})
			i += 7
		} else if len(expr) > i+9 && "SEND_DEST" == expr[i:i+9] {
			res = append(res, []interface{}{"OP", "SEND_DEST"})
			i += 8
		} else if len(expr) > i+8 && "DEL_DEST" == expr[i:i+8] {
			res = append(res, []interface{}{"OP", "DEL_DEST"})
			i += 7
		} else if len(expr) > i+8 && "UNDEFINE" == expr[i:i+8] {
			res = append(res, []interface{}{"OP", "UNDEFINE"})
			i += 7
		} else if len(expr) > i+12 && "UNSET_SOURCE" == expr[i:i+12] {
			res = append(res, []interface{}{"OP", "UNSET_SOURCE"})
			i += 13 // плюс скобки, которые мы не разбираем
		} else if len(expr) > i+7 && "REROUTE" == expr[i:i+7] {
			res = append(res, []interface{}{"OP", "REROUTE"})
			i += 8
		} else if len(expr) > i+10 && "UNSET_DEST" == expr[i:i+10] {
			res = append(res, []interface{}{"OP", "UNSET_DEST"})
			i += 11
		} else if len(expr) > i+12 && "RESET_SOURCE" == expr[i:i+12] {
			res = append(res, []interface{}{"OP", "RESET_SOURCE"})
			i += 13
		} else if len(expr) > i+3 && "str" == expr[i:i+3] && "(" == string(expr[i+3]) {
			res = append(res, []interface{}{"OP", "str"})
			i += 2
		} else if len(expr) > i+3 && "int" == expr[i:i+3] && "(" == string(expr[i+3]) {
			res = append(res, []interface{}{"OP", "int"})
			i += 2
		} else if len(expr) > i+5 && "float" == expr[i:i+5] && "(" == string(expr[i+5]) {
			res = append(res, []interface{}{"OP", "float"})
			i += 4
		} else if len(expr) > i+4 && "bool" == expr[i:i+4] && "(" == string(expr[i+4]) {
			res = append(res, []interface{}{"OP", "bool"})
			i += 3
		} else if len(expr) > i && "(" == string(expr[i]) {
			res = append(res, []interface{}{"BR", "("})
		} else if len(expr) > i && ")" == string(expr[i]) {
			res = append(res, []interface{}{"BR", ")"})
		} else if len(expr) > i && ":" == string(expr[i]) {
			res = append(res, []interface{}{"OP", ":"})
		} else if len(expr) > i+4 && "True" == expr[i:i+4] {
			// крайнее правое значение - операция, которая была выполнена
			// значение посередине - результат этой операции
			res = append(res, []interface{}{"VAL", "True", "True"})
			i += 3
		} else if len(expr) >= i+5 && "False" == expr[i:i+5] {
			res = append(res, []interface{}{"VAL", "False", "False"})
			i += 4
		} else if len(expr) > i && "=" == string(expr[i]) && "=" != string(expr[i+1]) {
			res = append(res, []interface{}{"OP", "="})
		} else if len(expr) > i && "[" == string(expr[i]) {
			res = append(res, []interface{}{"CD_BR", "["})
		} else if len(expr) > i && "]" == string(expr[i]) {
			res = append(res, []interface{}{"CD_BR", "]"})
		} else if len(expr) > i+2 && "<=" == expr[i:i+2] {
			res = append(res, []interface{}{"OP", "<="})
			i += 1
		} else if len(expr) > i+2 && ">=" == expr[i:i+2] {
			res = append(res, []interface{}{"OP", ">="})
			i += 1
		} else if len(expr) > i+2 && "==" == expr[i:i+2] {
			res = append(res, []interface{}{"OP", "=="})
			i += 1
		} else if len(expr) > i && "<" == string(expr[i]) {
			res = append(res, []interface{}{"OP", "<"})
		} else if len(expr) > i && ">" == string(expr[i]) {
			res = append(res, []interface{}{"OP", ">"})
		} else if len(expr) > i && "," == string(expr[i]) {
			res = append(res, []interface{}{"SEP", ","})
		} else if len(expr) > i+3 && "int" == expr[i:i+3] && "_" != string(expr[i+3]) {
			isType = true
			variables = append(variables, []interface{}{"int", "var_name", []interface{}{"var_val"}})
			i += 2
		} else if len(expr) > i+5 && "float" == expr[i:i+5] && "_" != string(expr[i+5]) {
			isType = true
			variables = append(variables, []interface{}{"float", "var_name", []interface{}{"var_val"}})
			i += 4
		} else if len(expr) > i+4 && "bool" == expr[i:i+4] && "_" != string(expr[i+4]) {
			isType = true
			variables = append(variables, []interface{}{"bool", "var_name", []interface{}{"var_val"}})
			i += 3
		} else if len(expr) > i+6 && "string" == expr[i:i+6] && "_" != string(expr[i+6]) {
			isType = true
			variables = append(variables, []interface{}{"string", "var_name", []interface{}{"var_val"}})
			i += 5
		} else if len(expr) > i+5 && "stack" == expr[i:i+5] && "_" != string(expr[i+5]) {
			isType = true
			variables = append(variables, []interface{}{"stack", "var_name", []interface{}{[]interface{}{"end"}}})
			i += 4
		} else {
			// число, либо переменная, либо строка
			// число
			// выражение "-(переменная)" недопустимо
			if unicode.IsDigit(rune(expr[i])) || "-" == string(expr[i]) {
				number := string(expr[i])
				exprInside := expr[i+1:]
				for _, ch := range exprInside {
					if unicode.IsDigit(ch) || "." == string(ch) {
						number += string(ch)
					} else {
						break
					}
				}

				res = append(res, []interface{}{"VAL", number})
				i += len(number) - 1
			} else if unicode.IsLetter(rune(expr[i])) || "$" == string(expr[i]) {
				// переменная
				// переменная может состоять только из:
				// а) латинских букв;
				// б) цифр;
				// в) символа нижнего подчеркивания.
				// переменная всегда начинается с буквы
				// системная переменная начинается с символа "$"

				varName := string(expr[i])
				exprInside := expr[i+1:]

				for _, ch := range exprInside {
					if unicode.IsLetter(ch) || unicode.IsDigit(ch) || "_" == string(ch) {
						varName += string(ch)
					} else {
						break
					}
				}

				if !isType {
					res = append(res, []interface{}{"VAR", varName})
				}

				i += len(varName) - 1

				if isType {
					isType = false
					variables[len(variables)-1][1] = varName

					if toTranspile {
						var err error

						goCommand := "defineVar(\"" + varName + "\")\n"

						_, err = transpileDest.WriteString(goCommand)
						if nil != err {
							return res, variables, err
						}

						if "stack" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							goCommand = "setVar(\"" + varName + "\", []interface{}{\"end\"})\n"

							_, err = transpileDest.WriteString(goCommand)
							if nil != err {
								return res, variables, err
							}
						}
					}
					if toPrimitive {
						var err error

						varType := fmt.Sprintf("%v", variables[len(variables)-1][0])
						_, err = primitiveDest.WriteString(varType + options.BendSep + varName + ";\n")
						if nil != err {
							panic(err)
						}

					}
				}
			} else if len([]rune(expr)) > i && "\"" == string([]rune(expr)[i]) {
				// строка
				var stringInside string

				j := i + 1
				offset := 1

				for !("\"" == string([]rune(expr)[j]) && "\\" != string([]rune(expr)[j-1])) {
					if !("\\" == string([]rune(expr)[j]) && "\"" == string([]rune(expr)[j+1])) {
						stringInside += string([]rune(expr)[j])
					} else {
						offset += 1
					}

					j += 1
				}

				res = append(res, []interface{}{"VAL", "\"" + stringInside + "\""})
				i += len(stringInside) + offset
			} else {
				err := errors.New("lexer: ERROR: can not recognize symbol " + "\"" + string(expr[i]) + "\"")
				return res, variables, err
			}

		}
		i += 1
	}

	if 0 == len(res) {
		res = append(res, []interface{}{"res", 0})
	}

	return res, variables, nil

}
