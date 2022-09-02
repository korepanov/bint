package internalTools

import (
	"bint.com/internal/const/options"
	"bint.com/internal/decryptor"
	"bint.com/internal/encrypter"
	"bint.com/internal/executor"
	. "bint.com/internal/lexer"
	"bint.com/internal/parser"
	. "bint.com/internal/primitiveLexer"
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
var piFlag = "-pi"
var poFlag = "-po"
var peFlag = "-pe"
var ciFlag = "-ci"
var coFlag = "-co"
var ceFlag = "-ce"
var kFlag = "-k"
var rFlag = flag.Bool("r", false, "Translation from B to Basm without debug information")

var FileToExecute string

func ParseArgs() (int, string, string, string, error) {
	var toTranslate int
	var rootSource string
	var rootDest string
	var keyDest string

	flag.StringVar(&iFlag, "i", "", "-i input.b")
	flag.StringVar(&oFlag, "o", "", "-o output.basm")
	flag.StringVar(&eFlag, "e", "", "-e program.basm")
	flag.StringVar(&piFlag, "pi", "", "-pi input.basm")
	flag.StringVar(&poFlag, "po", "", "-po output.bend")
	flag.StringVar(&peFlag, "pe", "", "-pe program.bend")
	flag.StringVar(&ciFlag, "ci", "", "-ci input.bend")
	flag.StringVar(&coFlag, "co", "", "-co output.benc")
	flag.StringVar(&ceFlag, "ce", "", "-ce program.benc")
	flag.StringVar(&kFlag, "k", "", "-k key.k")
	flag.Parse()

	if *help {
		flag.Usage()
		err := errors.New("help")
		return toTranslate, rootSource, rootDest, keyDest, err
	} else if "" == iFlag && "" == oFlag && "" != eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag {
		toTranslate = options.ExecBasm
		rootDest = eFlag
	} else if "" != iFlag && "" != oFlag && "" == eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag {
		toTranslate = options.UserTranslate
		rootSource = iFlag
		rootDest = oFlag
	} else if "" != piFlag && "" != poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag {
		toTranslate = options.Primitive
		rootSource = piFlag
		rootDest = poFlag
	} else if "" == piFlag && "" == poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" != peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag {
		toTranslate = options.InterpPrimitive
		rootDest = peFlag
	} else if "" == piFlag && "" == poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" == peFlag &&
		"" != ciFlag && "" != coFlag && "" == ceFlag && "" != kFlag {
		toTranslate = options.Encrypt
		rootSource = ciFlag
		rootDest = coFlag
		keyDest = kFlag
	} else if "" == piFlag && "" == poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" != ceFlag && "" != kFlag {
		toTranslate = options.ExecEncrypt
		rootDest = ceFlag
		keyDest = kFlag
	} else if "" == iFlag && "" == oFlag && "" == eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && *sFlag {
		toTranslate = options.Internal
	} else {
		flag.Usage()
		err := errors.New("invalid arguments")
		return toTranslate, rootSource, rootDest, keyDest, err
	}

	return toTranslate, rootSource, rootDest, keyDest, nil
}

func SetConf(toTranslate int, rootSource string, rootDest string, keyDest string, toTranslateInternal int, execBenv bool) (string, string, string, []string) {
	var filesListToExecute []string

	if options.Internal == toTranslate {
		toTranslate = toTranslateInternal

		if options.Internal == toTranslate {
			rootSource = "benv/for.b"
			rootDest = "benv/for.basm"
			if execBenv {
				filesListToExecute = []string{"benv/internal/build/import",
					"benv/internal/build/prep_func",
					"benv/internal/build/prep_if",
					"benv/internal/build/if",
					"benv/internal/build/long_function",
					"benv/internal/build/recurs",
					"benv/internal/build/func",
					"benv/internal/build/print_format"}

			} else {
				filesListToExecute = []string{"benv/internal/import.basm",
					"benv/internal/prep_func.basm",
					"benv/internal/prep_if.basm",
					"benv/internal/if.basm",
					"benv/internal/long_function.basm",
					"benv/internal/func.basm",
					"benv/internal/print_format.basm"}
			}
		} else if options.UserTranslate == toTranslate {
			rootSource = "program.b"
			rootDest = "program.basm"
			//filesListToExecute = []string{"benv/import.basm"}
			if execBenv {
				filesListToExecute = []string{"benv/build/import",
					"benv/build/prep_func",
					"benv/build/prep_if",
					"benv/build/if",
					"benv/build/long_function",
					"benv/build/recurs",
					"benv/build/func",
					"benv/build/print_format"}
			} else {
				filesListToExecute = []string{"benv/import.basm",
					//	"benv/trace.basm"}
					"benv/prep_func.basm",
					"benv/prep_if.basm",
					"benv/if.basm",
					"benv/long_function.basm",
					"benv/func.basm",
					"benv/print_format.basm"}
			}
		} else if options.Transpile == toTranslate {
			rootSource = "benv/internal/build/func.basm"
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
		} else if options.ExecBasm == toTranslate {
			//rootSource = "benv/internal/recurs_program.b"
			//rootDest = "benv/func.basm"
			rootDest = "benv/for.basm"
			filesListToExecute = []string{rootDest}
		} else if options.Primitive == toTranslate {
			rootSource = "bendBenv/func.basm"
			rootDest = "bendBenv/func.bend"
			//rootSource = "program.basm"
			//rootDest = "program.bend"
			filesListToExecute = []string{rootSource}
		} else if options.InterpPrimitive == toTranslate {
			rootSource = "program.b"
			rootDest = "program.basm"
			filesListToExecute = []string{"bendBenv/import.bend",
				"bendBenv/if.bend",
				"bendBenv/prep_func.bend",
				"bendBenv/long_function.bend",
				"bendBenv/func.bend",
				"bendBenv/print_format.bend"}
			//rootDest = "program.bend"
			//filesListToExecute = []string{rootDest}
		} else if options.Encrypt == toTranslate {
			rootSource = "bencBenv/print_format.bend"
			rootDest = "bencBenv/print_format.benc"
			keyDest = "bencBenv/print_format.k"
			//rootSource = "program.bend"
			//rootDest = "program.benc"
			//keyDest = "key.k"
		} else if options.ExecEncrypt == toTranslate {
			rootSource = "program.b"

			rootDest = "bencBenv/print_format.benc"
			keyDest = "bencBenv/print_format.k"
			filesListToExecute = []string{rootDest}

			rootDest = "program.basm"
		} else {
			panic(errors.New("set option to translate"))
		}

	} else if options.UserTranslate == toTranslate {
		filesListToExecute = []string{"benv/build/import", "benv/build/prep_func",
			"benv/build/long_function", "benv/build/func"}
	} else if options.ExecBasm == toTranslate {
		filesListToExecute = []string{rootDest}
	} else if options.Primitive == toTranslate {
		rootSource = piFlag
		rootDest = poFlag
		filesListToExecute = []string{rootSource}
	} else if options.InterpPrimitive == toTranslate {
		rootDest = peFlag
		filesListToExecute = []string{rootDest}
	} else if options.Encrypt == toTranslate {
		rootSource = ciFlag
		rootDest = coFlag
		keyDest = kFlag
	} else if options.ExecEncrypt == toTranslate {
		rootDest = ceFlag
		keyDest = kFlag
		filesListToExecute = []string{rootDest}
	} else {
		panic(errors.New("set option to translate"))
	}

	return rootSource, rootDest, keyDest, filesListToExecute
}
func ExecBenv(filesListToExecute []string, rootSource string, rootDest string) {
	FileToExecute = filesListToExecute[0]

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
		FileToExecute = filesListToExecute[i]
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

	FileToExecute = filesListToExecute[len(filesListToExecute)-1]

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

func Start(toTranslate int, filesListToExecute []string, rootSource string, rootDest string, keyDest string, sysMod int, execBenv bool) {
	COMMAND_COUNTER := 1
	var variables [][]interface{}
	var exprList [][]interface{}
	var infoListList [][]interface{}
	var source *os.File
	var dest *os.File
	var key *os.File
	var transpileDest *os.File
	var primitiveDest *os.File
	var sourceNewChunk func() string
	var wasGetCommandCounterByMark bool
	var newKey func() string

	systemStack := []interface{}{"end"}
	sourceCommandCounter := 1

	defer func() {
		if r := recover(); nil != r {
			fmt.Println("ERROR in " + FileToExecute + " at near line " + fmt.Sprintf("%v", LineCounter))
			fmt.Println(CommandToExecute)
			fmt.Println(r)
			os.Exit(1)
		}
	}()

	var err error

	if options.Validate == toTranslate {
		rootDest = ".temp.b"
		if execBenv {
			filesListToExecute = []string{"benv/build/import",
				"benv/build/prep_func",
				"benv/build/prep_if",
				"benv/build/if",
				"benv/build/long_function",
				"benv/build/func",
				"benv/build/print_format"}
		} else {
			filesListToExecute = []string{"benv/import.basm",
				"benv/trace.basm"}
			//	"benv/prep_func.basm",
			//	"benv/prep_if.basm",
			//	"benv/if.basm"}
			//"benv/long_function.basm",
			//"benv/func.basm",
			//"benv/print_format.basm"}
		}
	}

	if options.Internal != toTranslate {
		sysMod = toTranslate
	}

	if (options.Internal == toTranslate && (options.UserTranslate == sysMod || options.Internal == sysMod) && execBenv) ||
		options.UserTranslate == toTranslate {
		ExecBenv(filesListToExecute, rootSource, rootDest)
		return
	}

	if options.Transpile == sysMod {
		transpileDest, err = os.OpenFile(rootDest, os.O_APPEND|os.O_WRONLY, 0644)
		if nil != err {
			panic(err)
		}
	}

	if options.Primitive == sysMod {
		primitiveDest, err = os.Create(rootDest)
		if nil != err {
			panic(err)
		}
	}

	if options.Encrypt == sysMod {
		encrypter.Encrypt(rootSource, rootDest, keyDest)
		return
	}
	if options.ExecEncrypt == sysMod {
		key, err = os.Open(keyDest)
		if nil != err {
			panic(err)
		}
		newKey = EachChunk(key)
	}

	for _, FileToExecute = range filesListToExecute {
		COMMAND_COUNTER = 1
		sourceCommandCounter = 1
		LineCounter = 0
		variables = nil
		systemStack = []interface{}{"end"}
		var marks []string
		var inputedCode string
		var isBasmStyle bool

		f, err := os.Open(FileToExecute)

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
			//ввод команд
			CommandToExecute = strings.TrimSpace(chunk)
			inputedCode = CodeInput(chunk, !wasGetCommandCounterByMark)

			if options.ExecEncrypt == sysMod {
				inputedCode = decryptor.Decrypt(inputedCode, newKey())
			}

			if options.Transpile == sysMod {
				mark := GetMark(inputedCode)
				if "" != mark {
					_, err = transpileDest.WriteString("goto " + mark[1:] + "\n" + mark[1:] + ":\n")
					if nil != err {
						panic(err)
					}
				}
			}
			if options.Primitive == sysMod {
				mark := GetMark(inputedCode)
				if "" != mark {
					_, err = primitiveDest.WriteString(mark + ":\n")
					if nil != err {
						panic(err)
					}
				}
			}

			if options.InterpPrimitive != sysMod && options.ExecEncrypt != sysMod {
				var temp string
				if options.Primitive == sysMod {
					if "#" == string(inputedCode[0]) {
						pos := strings.Index(inputedCode, ":")
						if -1 == pos {
							panic(errors.New("ERROR: mark must end with \":\""))
						}
						pos++
						temp = inputedCode[pos:]
					} else {
						temp = inputedCode
					}
				}
				if len(temp) > 0 && "[" == string(temp[0]) {
					_, err = primitiveDest.WriteString(temp + ";\n")
					if nil != err {
						panic(err)
					}
					continue
				} else if -1 != strings.Index(temp, "[") &&
					(-1 == strings.Index(temp, "\"") || -1 != strings.Index(temp[:strings.Index(temp, "\"")], "[")) {
					_, err = primitiveDest.WriteString(temp + ";\n")
					if nil != err {
						panic(err)
					}
					continue
				} else if -1 != strings.Index(temp, ".") && -1 == strings.Index(temp[:strings.Index(temp, ".")], "\"") {
					_, err = primitiveDest.WriteString(temp + ";\n")
					if nil != err {
						panic(err)
					}
					continue

				} else {
					exprList, variables, err = LexicalAnalyze(inputedCode,
						variables, options.Transpile == sysMod, transpileDest, options.Primitive == sysMod, primitiveDest)
				}
			} else {
				var temp string

				if "#" == string(inputedCode[0]) {
					pos := strings.Index(inputedCode, ":")
					if -1 == pos {
						panic(errors.New("ERROR: mark must end with \":\""))
					}
					pos++
					temp = inputedCode[pos:]
				} else {
					temp = inputedCode
				}

				if len(temp) > 0 && "[" == string(temp[0]) {
					isBasmStyle = true
				} else if -1 != strings.Index(temp, "[") &&
					(-1 == strings.Index(temp, "\"") || -1 != strings.Index(temp[:strings.Index(temp, "\"")], "[")) {
					isBasmStyle = true
				} else if -1 != strings.Index(temp, ".") && -1 == strings.Index(temp[:strings.Index(temp, ".")], "\"") {
					isBasmStyle = true
				} else {
					isBasmStyle = false
					exprList, variables, err = PrimitiveLexicalAnalyze(inputedCode, variables)
					if nil != err {
						panic(err)
					}
				}

			}
			wasGetCommandCounterByMark = false

			if options.InterpPrimitive != sysMod && options.ExecEncrypt != sysMod && 0 != exprList[0][1] { // выражение содержит команды
				_, infoListList, systemStack = parser.Parse(exprList, variables, systemStack, options.HideTree,
					options.Transpile == sysMod, options.Primitive == sysMod, primitiveDest, transpileDest)
			} else if options.InterpPrimitive != sysMod && options.ExecEncrypt != sysMod && 0 == exprList[0][1] {
				continue
			} else if options.InterpPrimitive == sysMod || options.ExecEncrypt == sysMod {
				if isBasmStyle {
					exprList, variables, err = LexicalAnalyze(inputedCode,
						variables, options.Transpile == sysMod, transpileDest, options.Primitive == sysMod, primitiveDest)
					_, infoListList, systemStack = parser.Parse(exprList, variables, systemStack, options.HideTree,
						options.Transpile == sysMod, options.Primitive == sysMod, primitiveDest, transpileDest)
				} else {
					infoListList = exprList
				}
			} else if (options.InterpPrimitive == sysMod || options.ExecEncrypt == sysMod) && !isBasmStyle {
				continue
			}
			if 2 == len(infoListList[0]) && "res" == fmt.Sprintf("%v", infoListList[0][0]) &&
				0 == infoListList[0][1] {
				continue // выражение не содержит команд
			}
			for _, infoList := range infoListList {
				var res []interface{}
				res, variables, systemStack = executor.ExecuteTree(infoList, variables, systemStack,
					options.Transpile == sysMod, options.Primitive == sysMod, primitiveDest, transpileDest)
				if "print" != res[0] {
					if "goto" == fmt.Sprintf("%v", res[0]) {
						if options.Transpile != sysMod && options.Primitive != sysMod {
							wasGetCommandCounterByMark = true
							COMMAND_COUNTER, f, err = GetCommandCounterByMark(f, fmt.Sprintf("%v", res[1]))

							if nil != err {
								panic(err)
							}

							newChunk, err = SetCommandCounter(f, COMMAND_COUNTER)
							if nil != err {
								panic(err)
							}
							if options.ExecEncrypt == sysMod {
								newKey, err = SetCommandCounter(key, COMMAND_COUNTER)
								if nil != err {
									panic(err)
								}
							}
						} else if options.Transpile == sysMod {
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

								if options.Validate == toTranslate {
									v[2] = rootSource + "v"
								} else if *rFlag {
									v[2] = rootSource
								} else {
									v[2] = rootSource + "d"
								}

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
						if len(strRes) > 9 && `getVar(\"` == strRes[0:9] {
							strRes = `getVar("` + strRes[9:len(strRes)-3] + `")`
						}
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
		if options.ExecEncrypt == sysMod {
			err = key.Close()
			if nil != err {
				panic(err)
			}
		}
	}
}
