package lexer

import (
	"errors"
	"fmt"
	"os"
	"unicode"

	. "bint.com/internal/compilerVars"
	"bint.com/internal/const/options"
	. "bint.com/pkg/serviceTools"
)

func LexicalAnalyze(expr string, variables [][]interface{}, toTranspile bool, toCompile bool,
	transpileDest *os.File, toPrimitive bool, primitiveDest *os.File, dataFile *os.File, labelsFile *os.File,
	progFile *os.File) ([][]interface{}, [][]interface{}, error) {

	var res [][]interface{}

	i := 0
	isType := false

	for i < len([]rune(expr)) {
		newVariable := EachVariable(variables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if len([]rune(expr)) > i+len(v[1].(string)) {
				if v[1].(string) == string([]rune(expr)[i:i+len(v[1].(string))]) {
					// проверяем, что имя переменной не является частью имени другой переменной
					if i+len(v[1].(string)) < len([]rune(expr)) && !unicode.IsDigit(rune([]rune(expr)[i+len(v[1].(string))])) &&
						!unicode.IsLetter(rune([]rune(expr)[i+len(v[1].(string))])) && "_" != string([]rune(expr)[i+len(v[1].(string))]) {
						// проверяем, что это не присвоение
						if i+1 < len([]rune(expr)) && "=" != string([]rune(expr)[i+len(v[1].(string))]) {
							// проверяем, что найденная переменная не является частью какого-либо другого слова
							if i-1 >= 0 && !unicode.IsDigit(rune(expr[i-1])) &&
								!unicode.IsLetter(rune([]rune(expr)[i-1])) &&
								!unicode.IsDigit(rune([]rune(expr)[i+len(v[1].(string))])) &&
								!unicode.IsLetter(rune([]rune(expr)[i+len(v[1].(string))])) {
								res = append(res, []interface{}{"VAR", v[1]})
								i += len(v[1].(string))
							}
						}
					}
				}
			}
		}

		if len([]rune(expr)) > i+3 && "AND" == string([]rune(expr)[i:i+3]) {
			res = append(res, []interface{}{"OP", "AND"})
			i += 2
		} else if len([]rune(expr)) > i+2 && "OR" == string([]rune(expr)[i:i+2]) {
			res = append(res, []interface{}{"OP", "OR"})
			i += 1
		} else if len([]rune(expr)) > i+3 && "NOT" == string([]rune(expr)[i:i+3]) {
			res = append(res, []interface{}{"OP", "NOT"})
			i += 2
		} else if len([]rune(expr)) > i+3 && "XOR" == string([]rune(expr)[i:i+3]) {
			res = append(res, []interface{}{"OP", "XOR"})
			i += 2
		} else if len([]rune(expr)) > i && "." == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", "."})
		} else if len([]rune(expr)) > i && "+" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", "+"})
		} else if i > 0 && len([]rune(expr)) > i+1 && "-" == string([]rune(expr)[i]) && (unicode.IsDigit(rune([]rune(expr)[i-1])) ||
			unicode.IsLetter(rune([]rune(expr)[i-1])) ||
			")" == string([]rune(expr)[i-1])) && (unicode.IsDigit(rune([]rune(expr)[i+1])) || unicode.IsLetter(rune([]rune(expr)[i+1])) ||
			"(" == string([]rune(expr)[i+1]) || "$" == string([]rune(expr)[i+1])) {
			res = append(res, []interface{}{"OP", "-"})
		} else if len([]rune(expr)) > i && "*" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", "*"})
		} else if len([]rune(expr)) > i && "/" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", "/"})
		} else if len([]rune(expr)) > i && "@" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", "@"})
		} else if len([]rune(expr)) > i && "^" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", "^"})
		} else if len([]rune(expr)) > i+6 && "print(" == string([]rune(expr)[i:i+6]) {
			res = append(res, []interface{}{"OP", "print"})
			i += 4
		} else if len([]rune(expr)) > i+4 && "len(" == string([]rune(expr)[i:i+4]) {
			res = append(res, []interface{}{"OP", "len"})
			i += 2
		} else if len([]rune(expr)) > i+6 && "exists" == string([]rune(expr)[i:i+6]) {
			res = append(res, []interface{}{"OP", "exists"})
			i += 5
		} else if len([]rune(expr)) > i+6 && "index(" == string([]rune(expr)[i:i+6]) {
			res = append(res, []interface{}{"OP", "index"})
			i += 4
		} else if len([]rune(expr)) > i+9 && "is_letter" == string([]rune(expr)[i:i+9]) {
			res = append(res, []interface{}{"OP", "is_letter"})
			i += 8
		} else if len([]rune(expr)) > i+8 && "is_digit" == string([]rune(expr)[i:i+8]) {
			res = append(res, []interface{}{"OP", "is_digit"})
			i += 7
		} else if len([]rune(expr)) > i+8 && "reg_find" == string([]rune(expr)[i:i+8]) {
			res = append(res, []interface{}{"OP", "reg_find"})
			i += 7
		} else if len([]rune(expr)) > i+3 && "pop" == string([]rune(expr)[i:i+3]) {
			res = append(res, []interface{}{"OP", "pop"})
			i += 2
		} else if len([]rune(expr)) > i+4 && "push" == string([]rune(expr)[i:i+4]) {
			res = append(res, []interface{}{"OP", "push"})
			i += 3
		} else if len([]rune(expr)) > i+5 && "input" == string([]rune(expr)[i:i+5]) {
			res = append(res, []interface{}{"OP", "input"})
			i += 4
		} else if len([]rune(expr)) > i+8 && "$open_f(" == string([]rune(expr)[i:i+8]) {
			res = append(res, []interface{}{"OP", "$open_f"})
			i += 6
		} else if len([]rune(expr)) > i+8 && "$read_f(" == string([]rune(expr)[i:i+8]) {
			res = append(res, []interface{}{"OP", "$read_f"})
			i += 6
		} else if len([]rune(expr)) > i+9 && "$close_f(" == string(expr[i:i+9]) {
			res = append(res, []interface{}{"OP", "$close_f"})
			i += 7
		} else if len([]rune(expr)) > i+12 && "next_command" == string(expr[i:i+12]) {
			res = append(res, []interface{}{"OP", "next_command"})
			i += 11
		} else if len([]rune(expr)) > i+15 && "get_root_source" == string([]rune(expr)[i:i+15]) {
			res = append(res, []interface{}{"OP", "get_root_source"})
			i += 14
		} else if len([]rune(expr)) > i+13 && "get_root_dest" == string([]rune(expr)[i:i+13]) {
			res = append(res, []interface{}{"OP", "get_root_dest"})
			i += 12
		} else if len([]rune(expr)) > i+12 && "send_command" == string([]rune(expr)[i:i+12]) {
			res = append(res, []interface{}{"OP", "send_command"})
			i += 11
		} else if len([]rune(expr)) > i+4 && "goto" == string([]rune(expr)[i:i+4]) {
			res = append(res, []interface{}{"OP", "goto"})
			i += 3
		} else if len([]rune(expr)) > i+5 && "exit(" == string([]rune(expr)[i:i+5]) {
			res = append(res, []interface{}{"OP", "exit"})
			i += 3
		} else if len([]rune(expr)) > i && "#" == string([]rune(expr)[i]) {
			mark := string([]rune(expr)[i])
			i += 1

			for ":" != string([]rune(expr)[i]) && ")" != string([]rune(expr)[i]) {
				mark += string([]rune(expr)[i])
				i += 1
			}
			if ":" == string([]rune(expr)[i]) && toCompile {
				_, err := progFile.Write([]byte("\n." + mark[1:] + ":\n"))
				if nil != err {
					fmt.Println(err)
					os.Exit(1)
				}
				_, err = dataFile.Write([]byte("\nlabel" + fmt.Sprintf("%v", LabelCounter) + ":\n .quad ." + mark[1:] +
					"\nlabelName" + fmt.Sprintf("%v", LabelCounter) + ":" + "\n.ascii \"." + mark[1:] + "\"" + "\n.space 1,0"))
				if nil != err {
					fmt.Println(err)
					os.Exit(1)
				}

				_, err = labelsFile.Write([]byte("\nmov $labelName" + fmt.Sprintf("%v", LabelCounter) + ", %rbx" +
					"\n __initLabelsName" + fmt.Sprintf("%v", LabelCounter) + ": \n mov (%rbx), %dl \n cmp $0, %dl \n" +
					"jz __initLabelsNameEx" + fmt.Sprintf("%v", LabelCounter) + "\n mov %dl, (%rdi) \n inc %rbx \n inc %rdi" +
					"\n jmp __initLabelsName" + fmt.Sprintf("%v", LabelCounter) +
					"\n __initLabelsNameEx" + fmt.Sprintf("%v", LabelCounter) + ":\n movb $0, (%rdi)\n\n mov (label" +
					fmt.Sprintf("%v", LabelCounter) + "), %rax \n call __toStr\n add (valSize), %r9\n mov %r9, %rdi" +
					"\n mov $buf2, %rbx \n__initLabelsAddr" + fmt.Sprintf("%v", LabelCounter) + ":\n mov (%rbx), %dl " +
					"\n cmp $0, %dl \n jz __initLabelsAddrEx" + fmt.Sprintf("%v", LabelCounter) + "\n mov %dl, (%rdi)\n inc %rbx" +
					"\n inc %rdi \n jmp __initLabelsAddr" + fmt.Sprintf("%v", LabelCounter) +
					"\n __initLabelsAddrEx" + fmt.Sprintf("%v", LabelCounter) + ":\n movb $0, (%rdi)" +
					"\n add (valSize), %r9 \n mov %r9, %rdi \n mov %rdi, %r10 \n mov %r9, %rsi \n mov %r12, %rax\n" +
					"call __newLabelMem\n add (labelSize), %r12 \n\n mov %r10, %rdi \n mov %rsi, %r9"))
				LabelCounter++
			}
			if ")" == string([]rune(expr)[i]) {
				res = append(res, []interface{}{"VAR", mark})
				res = append(res, []interface{}{"BR", ")"})
			}
		} else if len([]rune(expr)) > i+10 && "SET_SOURCE" == string([]rune(expr)[i:i+10]) {
			res = append(res, []interface{}{"OP", "SET_SOURCE"})
			i += 9
		} else if len([]rune(expr)) > i+8 && "SET_DEST" == string([]rune(expr)[i:i+8]) {
			res = append(res, []interface{}{"OP", "SET_DEST"})
			i += 7
		} else if len([]rune(expr)) > i+9 && "SEND_DEST" == string([]rune(expr)[i:i+9]) {
			res = append(res, []interface{}{"OP", "SEND_DEST"})
			i += 8
		} else if len([]rune(expr)) > i+8 && "DEL_DEST" == string([]rune(expr)[i:i+8]) {
			res = append(res, []interface{}{"OP", "DEL_DEST"})
			i += 7
		} else if len([]rune(expr)) > i+8 && "UNDEFINE" == string([]rune(expr)[i:i+8]) {
			res = append(res, []interface{}{"OP", "UNDEFINE"})
			i += 7
		} else if len([]rune(expr)) > i+5 && "CLEAR" == string([]rune(expr)[i:i+5]) {
			res = append(res, []interface{}{"OP", "CLEAR"})
			i += 4
		} else if len([]rune(expr)) > i+12 && "UNSET_SOURCE" == string([]rune(expr)[i:i+12]) {
			res = append(res, []interface{}{"OP", "UNSET_SOURCE"})
			i += 13 // плюс скобки, которые мы не разбираем
		} else if len([]rune(expr)) > i+7 && "REROUTE" == string([]rune(expr)[i:i+7]) {
			res = append(res, []interface{}{"OP", "REROUTE"})
			i += 8
		} else if len([]rune(expr)) > i+10 && "UNSET_DEST" == string([]rune(expr)[i:i+10]) {
			res = append(res, []interface{}{"OP", "UNSET_DEST"})
			i += 11
		} else if len([]rune(expr)) > i+12 && "RESET_SOURCE" == string([]rune(expr)[i:i+12]) {
			res = append(res, []interface{}{"OP", "RESET_SOURCE"})
			i += 13
		} else if len([]rune(expr)) > i+3 && "str" == string([]rune(expr)[i:i+3]) && "(" == string([]rune(expr)[i+3]) {
			res = append(res, []interface{}{"OP", "str"})
			i += 2
		} else if len([]rune(expr)) > i+3 && "int" == string([]rune(expr)[i:i+3]) && "(" == string([]rune(expr)[i+3]) {
			res = append(res, []interface{}{"OP", "int"})
			i += 2
		} else if len([]rune(expr)) > i+5 && "float" == string([]rune(expr)[i:i+5]) && "(" == string([]rune(expr)[i+5]) {
			res = append(res, []interface{}{"OP", "float"})
			i += 4
		} else if len([]rune(expr)) > i+4 && "bool" == string([]rune(expr)[i:i+4]) && "(" == string([]rune(expr)[i+4]) {
			res = append(res, []interface{}{"OP", "bool"})
			i += 3
		} else if len([]rune(expr)) > i && "(" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"BR", "("})
		} else if len([]rune(expr)) > i && ")" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"BR", ")"})
		} else if len([]rune(expr)) > i && ":" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", ":"})
		} else if len([]rune(expr)) >= i+4 && "True" == string([]rune(expr)[i:i+4]) {
			// крайнее правое значение - операция, которая была выполнена
			// значение посередине - результат этой операции
			res = append(res, []interface{}{"VAL", "True", "True"})
			i += 3
		} else if len([]rune(expr)) >= i+5 && "False" == string([]rune(expr)[i:i+5]) {
			res = append(res, []interface{}{"VAL", "False", "False"})
			i += 4
		} else if len([]rune(expr)) > i && "=" == string([]rune(expr)[i]) && "=" != string([]rune(expr)[i+1]) {
			res = append(res, []interface{}{"OP", "="})
		} else if len([]rune(expr)) > i && "[" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"CD_BR", "["})
		} else if len([]rune(expr)) > i && "]" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"CD_BR", "]"})
		} else if len([]rune(expr)) > i+2 && "<=" == string([]rune(expr)[i:i+2]) {
			res = append(res, []interface{}{"OP", "<="})
			i += 1
		} else if len([]rune(expr)) > i+2 && ">=" == string([]rune(expr)[i:i+2]) {
			res = append(res, []interface{}{"OP", ">="})
			i += 1
		} else if len([]rune(expr)) > i+2 && "==" == string([]rune(expr)[i:i+2]) {
			res = append(res, []interface{}{"OP", "=="})
			i += 1
		} else if len([]rune(expr)) > i && "<" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", "<"})
		} else if len([]rune(expr)) > i && ">" == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"OP", ">"})
		} else if len([]rune(expr)) > i && "," == string([]rune(expr)[i]) {
			res = append(res, []interface{}{"SEP", ","})
		} else if len([]rune(expr)) > i+3 && "int" == string([]rune(expr)[i:i+3]) && "_" != string([]rune(expr)[i+3]) && IsValidString(string([]rune(expr)[i+3:])) && !isType {
			isType = true
			variables = append(variables, []interface{}{"int", "var_name", []interface{}{"0"}})
			i += 2
		} else if len([]rune(expr)) > i+5 && "float" == string([]rune(expr)[i:i+5]) && "_" != string([]rune(expr)[i+5]) && IsValidString(string([]rune(expr)[i+5:])) && !isType {
			isType = true
			variables = append(variables, []interface{}{"float", "var_name", []interface{}{"0"}})
			i += 4
		} else if len([]rune(expr)) > i+4 && "bool" == string([]rune(expr)[i:i+4]) && "_" != string([]rune(expr)[i+4]) && IsValidString(string([]rune(expr)[i+4:])) && !isType {
			isType = true
			variables = append(variables, []interface{}{"bool", "var_name", []interface{}{"False"}})
			i += 3
		} else if len([]rune(expr)) > i+6 && "string" == string([]rune(expr)[i:i+6]) && "_" != string([]rune(expr)[i+6]) && IsValidString(string([]rune(expr)[i+6:])) && !isType {
			isType = true
			variables = append(variables, []interface{}{"string", "var_name", []interface{}{""}})
			i += 5
		} else if len([]rune(expr)) > i+5 && "stack" == string([]rune(expr)[i:i+5]) && "_" != string([]rune(expr)[i+5]) && IsValidString(string([]rune(expr)[i+5:])) && !isType {
			isType = true
			variables = append(variables, []interface{}{"stack", "var_name", []interface{}{[]interface{}{"end"}}})
			i += 4
		} else {

			// число, либо переменная, либо строка
			// число
			// выражение "-(переменная)" недопустимо
			if unicode.IsDigit(rune([]rune(expr)[i])) || "-" == string([]rune(expr)[i]) {
				number := string([]rune(expr)[i])
				exprInside := string([]rune(expr)[i+1:])
				for _, ch := range exprInside {
					if unicode.IsDigit(ch) || "." == string(ch) {
						number += string(ch)
					} else {
						break
					}
				}

				res = append(res, []interface{}{"VAL", number})
				i += len(number) - 1
			} else if unicode.IsLetter(rune([]rune(expr)[i])) || "$" == string([]rune(expr)[i]) {
				// переменная
				// переменная может состоять только из:
				// а) латинских букв;
				// б) цифр;
				// в) символа нижнего подчеркивания.
				// переменная всегда начинается с буквы
				// системная переменная начинается с символа "$"

				varName := string([]rune(expr)[i])
				exprInside := string([]rune(expr)[i+1:])

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
						} else if "string" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							goCommand = "setVar(\"" + varName + "\", \"\")\n"

							_, err = transpileDest.WriteString(goCommand)
							if nil != err {
								return res, variables, err
							}
						} else if "int" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							goCommand = "setVar(\"" + varName + "\", 0)\n"

							_, err = transpileDest.WriteString(goCommand)
							if nil != err {
								return res, variables, err
							}
						} else if "float" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							goCommand = "setVar(\"" + varName + "\", 0)\n"

							_, err = transpileDest.WriteString(goCommand)
							if nil != err {
								return res, variables, err
							}
						} else if "bool" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							goCommand = "setVar(\"" + varName + "\", false)\n"

							_, err = transpileDest.WriteString(goCommand)
							if nil != err {
								return res, variables, err
							}
						}
					}
					if toCompile {
						if "stack" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							print("")
						} else if "string" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							_, err := dataFile.Write([]byte("\nvarName" + fmt.Sprintf("%v", VarsCounter) + ":" +
								"\n.ascii \"" + fmt.Sprintf("%v", variables[len(variables)-1][1]) +
								"\"\n.space 1, 0" +
								"\nlenVarName" + fmt.Sprintf("%v", VarsCounter) + " = . - varName" + fmt.Sprintf("%v", VarsCounter)))
							if nil != err {
								fmt.Println(err)
								os.Exit(1)
							}

							_, err = progFile.Write([]byte("\nmov $lenVarName, %rsi \n mov $varName, %rdx \n mov $lenVarName" +
								fmt.Sprintf("%v", VarsCounter) + ", %rax \n mov $varName" + fmt.Sprintf("%v", VarsCounter) +
								", %rdi\n call __set \n mov $lenVarType, %rsi \n mov $varType, %rdx \n mov $lenStringType, %rax" +
								"\n mov $stringType, %rdi\n call __set \n call __defineVar"))

							CompilerVars[fmt.Sprintf("%v", variables[len(variables)-1][1])] = fmt.Sprintf("%v", VarsCounter)
							VarsCounter++
						} else if "int" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							_, err := dataFile.Write([]byte("\nvarName" + fmt.Sprintf("%v", VarsCounter) + ":" +
								"\n.ascii \"" + fmt.Sprintf("%v", variables[len(variables)-1][1]) +
								"\"\n.space 1, 0" +
								"\nlenVarName" + fmt.Sprintf("%v", VarsCounter) + " = . - varName" + fmt.Sprintf("%v", VarsCounter)))
							if nil != err {
								fmt.Println(err)
								os.Exit(1)
							}

							_, err = progFile.Write([]byte("\nmov $lenVarName, %rsi \n mov $varName, %rdx" +
								"\n mov $lenVarName" + fmt.Sprintf("%v", VarsCounter) + ", %rax \n mov $varName" +
								fmt.Sprintf("%v", VarsCounter) + ", %rdi \n call __set \n mov $lenVarType, %rsi \n mov $varType, %rdx " +
								"\n mov $lenIntType, %rax \n mov $intType, %rdi \n call __set \n mov $varName, %rcx \n mov $varType, %rdx  " +
								"\n call __defineVar"))
							if nil != err {
								fmt.Println(err)
								os.Exit(1)
							}
							CompilerVars[fmt.Sprintf("%v", variables[len(variables)-1][1])] = fmt.Sprintf("%v", VarsCounter)
							VarsCounter++
						} else if "float" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							_, err := dataFile.Write([]byte("\nvarName" + fmt.Sprintf("%v", VarsCounter) + ":" +
								"\n.ascii \"" + fmt.Sprintf("%v", variables[len(variables)-1][1]) +
								"\"\n.space 1, 0" +
								"\nlenVarName" + fmt.Sprintf("%v", VarsCounter) + " = . - varName" + fmt.Sprintf("%v", VarsCounter)))
							if nil != err {
								fmt.Println(err)
								os.Exit(1)
							}

							_, err = progFile.Write([]byte("\nmov $lenVarName, %rsi \n mov $varName, %rdx" +
								"\n mov $lenVarName" + fmt.Sprintf("%v", VarsCounter) + ", %rax \n mov $varName" +
								fmt.Sprintf("%v", VarsCounter) + ", %rdi \n call __set \n mov $lenVarType, %rsi \n mov $varType, %rdx " +
								"\n mov $lenFloatType, %rax \n mov $floatType, %rdi \n call __set \n mov $varName, %rcx \n mov $varType, %rdx  " +
								"\n call __defineVar"))
							if nil != err {
								fmt.Println(err)
								os.Exit(1)
							}
							CompilerVars[fmt.Sprintf("%v", variables[len(variables)-1][1])] = fmt.Sprintf("%v", VarsCounter)
							VarsCounter++
						} else if "bool" == fmt.Sprintf("%v", variables[len(variables)-1][0]) {
							_, err := dataFile.Write([]byte("\nvarName" + fmt.Sprintf("%v", VarsCounter) + ":" +
								"\n.ascii \"" + fmt.Sprintf("%v", variables[len(variables)-1][1]) +
								"\"\n.space 1, 0" +
								"\nlenVarName" + fmt.Sprintf("%v", VarsCounter) + " = . - varName" + fmt.Sprintf("%v", VarsCounter)))
							if nil != err {
								fmt.Println(err)
								os.Exit(1)
							}

							_, err = progFile.Write([]byte("\nmov $lenVarName, %rsi \n mov $varName, %rdx" +
								"\n mov $lenVarName" + fmt.Sprintf("%v", VarsCounter) + ", %rax \n mov $varName" +
								fmt.Sprintf("%v", VarsCounter) + ", %rdi \n call __set \n mov $lenVarType, %rsi \n mov $varType, %rdx " +
								"\n mov $lenBoolType, %rax \n mov $boolType, %rdi \n call __set \n mov $varName, %rcx \n mov $varType, %rdx  " +
								"\n call __defineVar"))
							if nil != err {
								fmt.Println(err)
								os.Exit(1)
							}
							CompilerVars[fmt.Sprintf("%v", variables[len(variables)-1][1])] = fmt.Sprintf("%v", VarsCounter)
							VarsCounter++
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
				i += len([]rune(stringInside)) + offset
			} else {
				err := errors.New("lexer: ERROR: can not recognize symbol " + "\"" + string([]rune(expr)[i]) + "\"")
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
