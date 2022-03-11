package main

import (
	"bint.com/internal/executor"
	. "bint.com/internal/lexer"
	"bint.com/internal/parser"
	. "bint.com/pkg/serviceTools"
	"fmt"
	"os"
	"strings"
)

func main() {
	COMMAND_COUNTER := 1
	var variables [][]interface{}
	var exprList [][]interface{}
	var infoListList [][]interface{}
	var source *os.File
	var dest *os.File
	var filesListToExecute []string
	var sourceNewChunk func() string

	toTranslate := false

	if toTranslate {
		fmt.Println("Translating...")
		filesListToExecute = []string{"benv/func.basm"}
	} else {
		filesListToExecute = []string{"program.basm"}
	}

	systemStack := []interface{}{"end"}
	sourceCommandCounter := 1

	for _, fileToExecute := range filesListToExecute {
		COMMAND_COUNTER = 1
		sourceCommandCounter = 1
		variables = nil
		systemStack = []interface{}{"end"}

		f, err := os.Open(fileToExecute)

		if nil != err {
			panic(err)
		}
		newChunk := EachChunk(f)
		for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
			exprList, variables, err = LexicalAnalyze(CodeInput(chunk), variables)
			if nil != err {
				panic(err)
			}

			if 0 != exprList[0][1] { // выражение содержит команды
				_, infoListList, systemStack = parser.Parse(exprList, variables, systemStack, false)
			} else {
				continue
			}

			for _, infoList := range infoListList {
				var res []interface{}
				res, variables, systemStack = executor.ExecuteTree(infoList, variables, systemStack)
				if "print" != res[0] {
					if "goto" == fmt.Sprintf("%v", res[0]) {
						//_, err = f.Seek(0, 0)

						//if nil != err {
						//	panic(err)
						//}

						COMMAND_COUNTER, f, err = GetCommandCounterByMark(f, fmt.Sprintf("%v", res[1]))

						if nil != err {
							panic(err)
						}

						newChunk, err = SetCommandCounter(f, COMMAND_COUNTER)
						if nil != err {
							panic(err)
						}

						//COMMAND_COUNTER -= 1 // в конце цикла счетчик команд увеличивается на 1
						continue
					}
					if "SET_SOURCE" == fmt.Sprintf("%v", res[0]) {
						source = res[1].(*os.File)
						sourceNewChunk = EachChunk(source)
						continue
					}
					if "UNSET_SOURCE" == fmt.Sprintf("%v", res[0]) {
						sourceCommandCounter = 1
						err := source.Close()
						if nil != err {
							panic(err)
						}

					}
					if "UNSET_DEST" == fmt.Sprintf("%v", res[0]) {
						err := dest.Close()
						if nil != err {
							panic(err)
						}
					}
					if "RESET_SOURCE" == fmt.Sprintf("%v", res[0]) {
						sourceCommandCounter = 1
					}
					if "SET_DEST" == fmt.Sprintf("%v", res[0]) {
						dest = res[1].(*os.File)
						continue
					}
					if "REROUTE" == fmt.Sprintf("%v", res[0]) {
						//_, err = dest.Seek(0, 0)
						//if nil != err{
						//	panic(err)
						//}
						//_, err = source.Seek(0, 0)
						//if nil != err{
						//	panic(err)
						//}
						err = dest.Close()
						if nil != err {
							panic(err)
						}
						err = source.Close()
						if nil != err {
							panic(err)
						}
						temp := dest
						dest = source
						source = temp

						source, err = os.Open(source.Name())
						if nil != err {
							panic(err)
						}

						dest, err = os.OpenFile(dest.Name(), os.O_WRONLY, 0666)
						if nil != err {
							panic(err)
						}
						sourceNewChunk = EachChunk(source)
						continue
					}
					if "next_command" == fmt.Sprintf("%v", res[0]) {
						varName := res[1].([]interface{})[0]
						newVariable := EachVariable(variables)
						for v := newVariable(); "end" != v[0]; v = newVariable() {
							if fmt.Sprintf("%v", varName) == fmt.Sprintf("%v", v[1]) {
								var code string
								code = sourceNewChunk()
								v[2] = CodeInput(code)
								break
							}
						}

						sourceCommandCounter++
						continue

					}
					if "send_command" == fmt.Sprintf("%v", res[0]) {
						result := ValueFoldInterface(res[1]).(string)
						if "" != result && ("\n" != result) {
							if "\"" == string(result[0]) && "\"" == string(result[len(result)-1]) {
								result = result[1 : len(result)-1]
							}
							result = strings.Replace(result, "\\n", "\n", -1)
							_, err := dest.WriteString(result + ";\n")
							if nil != err {
								panic(err)
							}
						}
					}
				}
				var strRes string

				if "print" == res[0] {
					strRes = strings.Replace(fmt.Sprintf("%v", ValueFoldInterface(res[1])), "\\n", "\n", -1)

					if "" != strRes && "\"" == string(strRes[0]) && "\"" == string(strRes[len(strRes)-1]) {
						strRes = strRes[1 : len(strRes)-1]
					}

					fmt.Print(strRes)
				}
			}

		}
		err = f.Close()

		if nil != err {
			panic(err)
		}
	}
}
