package internalTools

import (
	"bint.com/internal/executor"
	. "bint.com/internal/lexer"
	"bint.com/internal/options"
	"bint.com/internal/parser"
	. "bint.com/pkg/serviceTools"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var help = flag.Bool("help", false, "Show help")
var iFlag = "-i"
var oFlag = "-o"
var sFlag = flag.Bool("s", false, "Use system options (for system debug only)")
var eFlag = "-e"
var fileToExecute string

func ParseArgs() (int, string, string, error) {
	var toTranslate int
	var rootSource string
	var rootDest string

	flag.StringVar(&iFlag, "i", "", "-i input.b")
	flag.StringVar(&oFlag, "o", "", "-o output.basm")
	flag.StringVar(&eFlag, "e", "", "-e program.basm")
	flag.Parse()

	if *help {
		flag.Usage()
		err := errors.New("help")
		return toTranslate, rootSource, rootDest, err
	} else if "" == iFlag && "" == oFlag && "" != eFlag {
		toTranslate = options.No
		rootDest = eFlag
	} else if "" != iFlag && "" != oFlag && "" == eFlag {
		toTranslate = options.User
		rootSource = iFlag
		rootDest = oFlag
	} else if "" == iFlag && "" == oFlag && "" == eFlag && *sFlag {
		toTranslate = options.Internal
	} else {
		flag.Usage()
		err := errors.New("invalid arguments")
		return toTranslate, rootSource, rootDest, err
	}

	return toTranslate, rootSource, rootDest, nil
}

func SetConf(toTranslate int, rootSource string, rootDest string, toTranslateInternal int, execBenv bool) (string, string, []string) {
	var filesListToExecute []string

	if options.Internal == toTranslate {
		toTranslate = toTranslateInternal

		if options.Internal == toTranslate {
			rootSource = "benv/import.b"
			rootDest = "benv/import.basm"
			if execBenv {
				filesListToExecute = []string{"benv/internal/build/prep_func",
					"benv/internal/build/long_function", "benv/internal/build/func"}

			} else {
				filesListToExecute = []string{"benv/internal/prep_func.basm", "benv/internal/long_function.basm",
					"benv/internal/func.basm"}
			}
		} else if options.User == toTranslate {
			rootSource = "program.b"
			rootDest = "program.basm"
			//filesListToExecute = []string{"benv/import.basm"}
			if execBenv {
				filesListToExecute = []string{"benv/build/import", "benv/build/prep_func", "benv/build/long_function",
					"benv/build/func"}
			} else {
				filesListToExecute = []string{"benv/import.basm", "benv/prep_func.basm", "benv/long_function.basm",
					"benv/func.basm"}
			}
		} else if options.Transpile == toTranslate {
			rootSource = "benv/internal/func.basm"
			rootDest = "benv/internal/build/main.go"

			source, err := os.Open("benv/build/pattern.p")
			if nil != err {
				panic(err)
			}
			dest, err := os.Create(rootDest)

			if nil != err {
				panic(err)
			}
			_, err = io.Copy(dest, source)
			if nil != err {
				panic(err)
			}
			err = source.Close()
			if nil != err {
				panic(err)
			}

			err = dest.Close()
			if nil != err {
				panic(err)
			}
			filesListToExecute = []string{rootSource}
		} else if options.No == toTranslate {
			rootSource = "program.b"
			rootDest = "benv/import.basm"
			filesListToExecute = []string{rootDest}
		} else {
			panic(errors.New("set option to translate"))
		}

	} else if options.User == toTranslate {
		filesListToExecute = []string{"benv/build/import", "benv/build/prep_func",
			"benv/build/long_function", "benv/build/func"}
	} else if options.No == toTranslate {
		filesListToExecute = []string{rootDest}
	} else {
		panic(errors.New("set option to translate"))
	}

	return rootSource, rootDest, filesListToExecute
}
func ExecBenv(filesListToExecute []string, rootSource string, rootDest string) {
	fileToExecute = filesListToExecute[0]

	cmd := exec.Command(filesListToExecute[0], "-i", rootSource)

	err := cmd.Start()
	if nil != err {
		panic(err)
	}
	err = cmd.Wait()
	if nil != err {
		panic(err)
	}

	for i := 1; i < len(filesListToExecute)-1; i++ {
		fileToExecute = filesListToExecute[i]
		cmd := exec.Command(filesListToExecute[i])
		err := cmd.Start()
		if nil != err {
			panic(err)
		}
		err = cmd.Wait()
		if nil != err {
			panic(err)
		}
	}

	fileToExecute = filesListToExecute[len(filesListToExecute)-1]

	cmd = exec.Command(filesListToExecute[len(filesListToExecute)-1], "-o", rootDest)
	err = cmd.Start()
	if nil != err {
		panic(err)
	}
	err = cmd.Wait()
	if nil != err {
		panic(err)
	}
}

func Start(toTranslate int, filesListToExecute []string, rootSource string, rootDest string, sysMod int, execBenv bool) {
	COMMAND_COUNTER := 1
	var variables [][]interface{}
	var exprList [][]interface{}
	var infoListList [][]interface{}
	var source *os.File
	var dest *os.File
	var transpileDest *os.File
	var sourceNewChunk func() string
	var wasGetCommandCounterByMark bool

	systemStack := []interface{}{"end"}
	sourceCommandCounter := 1

	defer func() {
		if r := recover(); nil != r {
			fmt.Println("ERROR in " + fileToExecute + " at near line " + fmt.Sprintf("%v", LineCounter))
			fmt.Println(CommandToExecute)
			fmt.Println(r)
		}
	}()

	var err error
	if (options.Internal == toTranslate && (options.User == sysMod || options.Internal == sysMod) && execBenv) ||
		options.User == toTranslate {
		ExecBenv(filesListToExecute, rootSource, rootDest)
		return
	}

	if options.Transpile == sysMod {
		transpileDest, err = os.OpenFile(rootDest, os.O_APPEND|os.O_WRONLY, 0644)
		if nil != err {
			panic(err)
		}
	}

	for _, fileToExecute = range filesListToExecute {
		COMMAND_COUNTER = 1
		sourceCommandCounter = 1
		LineCounter = 0
		variables = nil
		systemStack = []interface{}{"end"}
		var marks []string
		var inputedCode string

		f, err := os.Open(fileToExecute)

		if nil != err {
			panic(err)
		}
		newChunk := EachChunk(f)
		// собираем все метки файла в marks
		if options.Transpile == sysMod {
			for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
				CommandToExecute = strings.TrimSpace(chunk)
				inputedCode = CodeInput(chunk, false)

				mark := GetMark(inputedCode)
				if "" != mark {
					marks = append(marks, mark)
				}

			}

			_, err := f.Seek(0, 0)

			if nil != err {
				panic(err)
			}

			newChunk = EachChunk(f)
		}
		for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
			CommandToExecute = strings.TrimSpace(chunk)
			inputedCode = CodeInput(chunk, !wasGetCommandCounterByMark)

			if options.Transpile == sysMod {
				mark := GetMark(inputedCode)
				if "" != mark {
					_, err = transpileDest.WriteString("goto " + mark[1:] + "\n" + mark[1:] + ":\n")
					if nil != err {
						panic(err)
					}
				}
			}

			exprList, variables, err = LexicalAnalyze(inputedCode,
				variables, options.Transpile == sysMod, transpileDest)
			wasGetCommandCounterByMark = false

			if nil != err {
				panic(err)
			}

			if 0 != exprList[0][1] { // выражение содержит команды
				_, infoListList, systemStack = parser.Parse(exprList, variables, systemStack, options.HideTree,
					options.Transpile == sysMod, transpileDest)
			} else {
				continue
			}

			for _, infoList := range infoListList {
				var res []interface{}
				res, variables, systemStack = executor.ExecuteTree(infoList, variables, systemStack,
					options.Transpile == sysMod, transpileDest)
				if "print" != res[0] {
					if "goto" == fmt.Sprintf("%v", res[0]) {
						if options.Transpile != sysMod {
							wasGetCommandCounterByMark = true
							COMMAND_COUNTER, f, err = GetCommandCounterByMark(f, fmt.Sprintf("%v", res[1]))

							if nil != err {
								panic(err)
							}

							newChunk, err = SetCommandCounter(f, COMMAND_COUNTER)
							if nil != err {
								panic(err)
							}
						} else {
							if (len(fmt.Sprintf("%v", res[1])) > 7 && "#getVar" != fmt.Sprintf("%v", res[1])[0:7]) ||
								len(fmt.Sprintf("%v", res[1])) <= 7 {
								_, err = transpileDest.WriteString(fmt.Sprintf("%v", res[0]) + " " + fmt.Sprintf("%v", res[1])[1:] + "\n")
								if nil != err {
									panic(err)
								}
							} else {
								for _, mark := range marks {
									_, err = transpileDest.WriteString("if \"" + mark + "\" == fmt.Sprintf(\"%v\", " +
										fmt.Sprintf("%v", res[1])[1:] + "){\ngoto " + mark[1:] + "\n}\n")
									if nil != err {
										panic(err)
									}
								}
							}
						}
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
						_, err := source.Seek(0, 0)
						if nil != err {
							panic(err)
						}
						sourceNewChunk = EachChunk(source)
					}
					if "SET_DEST" == fmt.Sprintf("%v", res[0]) {
						dest = res[1].(*os.File)
						continue
					}
					if "REROUTE" == fmt.Sprintf("%v", res[0]) {
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
								v[2] = CodeInput(code, false)
								break
							}
						}

						sourceCommandCounter++
						continue

					}
					if "get_root_source" == fmt.Sprintf("%v", res[0]) {
						varName := res[1].([]interface{})[0]
						newVariable := EachVariable(variables)
						for v := newVariable(); "end" != v[0]; v = newVariable() {
							if fmt.Sprintf("%v", varName) == fmt.Sprintf("%v", v[1]) {
								v[2] = rootSource
								if options.Transpile == sysMod {
									_, err = transpileDest.WriteString("setVar(" +
										fmt.Sprintf("%v", varName) + ", " + v[2].(string))
									if nil != err {
										panic(err)
									}
								}
								break
							}
						}
					}
					if "get_root_dest" == fmt.Sprintf("%v", res[0]) {
						varName := res[1].([]interface{})[0]
						newVariable := EachVariable(variables)
						for v := newVariable(); "end" != v[0]; v = newVariable() {
							if fmt.Sprintf("%v", varName) == fmt.Sprintf("%v", v[1]) {
								v[2] = rootDest
								if options.Transpile == sysMod {
									_, err = transpileDest.WriteString("setVar(" +
										fmt.Sprintf("%v", varName) + ", " + v[2].(string))
									if nil != err {
										panic(err)
									}
								}
								break
							}
						}
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
					strRes = fmt.Sprintf("%v", ValueFoldInterface(res[1]))
					if options.Transpile == sysMod {
						goCommand := "fmt.Print(" + strRes + ")\n"
						_, err = transpileDest.WriteString(goCommand)
						if nil != err {
							panic(err)
						}
					}

					strRes = strings.Replace(strRes, "\\n", "\n", -1)

					if "" != strRes && "\"" == string(strRes[0]) && "\"" == string(strRes[len(strRes)-1]) {
						strRes = strRes[1 : len(strRes)-1]
					}

					if options.Transpile != sysMod {
						fmt.Print(strRes)
					}
				}
			}

		}
		err = f.Close()

		if nil != err {
			panic(err)
		}

		if options.Transpile == sysMod {
			_, err = transpileDest.WriteString("}\n")
			if nil != err {
				panic(err)
			}
			err = transpileDest.Close()

			if nil != err {
				panic(err)
			}
		}
	}
}
