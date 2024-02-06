package internalTools

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bint.com/internal/compiler"
	"bint.com/internal/compilerVars"
	"bint.com/internal/const/options"
	"bint.com/internal/decryptor"
	"bint.com/internal/encrypter"
	"bint.com/internal/executor"
	. "bint.com/internal/lexer"
	"bint.com/internal/parser"
	. "bint.com/internal/primitiveLexer"
	. "bint.com/pkg/serviceTools"
)

var help = flag.Bool("help", false, "show help")
var iFlag = "-i"
var oFlag = "-o"
var sFlag = flag.Bool("s", false, "use system options (for system debug only)")
var eFlag = "-e"
var piFlag = "-pi"
var poFlag = "-po"
var peFlag = "-pe"
var ciFlag = "-ci"
var coFlag = "-co"
var ceFlag = "-ce"
var biFlag = "-bi"
var boFlag = "bo"
var kFlag = "-k"
var vFlag = flag.Bool("version", false, "print version")

var rFlag = flag.Bool("r", false, "translation from b to basm without debug information")

var FileToExecute string

func ParseArgs() (int, string, string, string, error) {
	var toTranslate int
	var rootSource string
	var rootDest string
	var keyDest string

	flag.StringVar(&iFlag, "i", "", "file to translate from b to basm")
	flag.StringVar(&oFlag, "o", "", "output file translating b to basm")
	flag.StringVar(&eFlag, "e", "", "execute basm file")
	flag.StringVar(&piFlag, "pi", "", "file to translate from basm to bend")
	flag.StringVar(&poFlag, "po", "", "output file translating basm to bend")
	flag.StringVar(&peFlag, "pe", "", "execute bend file")
	flag.StringVar(&ciFlag, "ci", "", "file to translate from bend to benc (need to specify key file!)")
	flag.StringVar(&coFlag, "co", "", "output file translating bend to benc (need to specify key file!)")
	flag.StringVar(&ceFlag, "ce", "", "execute benc file (need to specify key file!)")
	flag.StringVar(&biFlag, "bi", "", "basm file to compile")
	flag.StringVar(&boFlag, "bo", "", "output file compiling basm")
	flag.StringVar(&kFlag, "k", "", "specify key file")
	flag.Parse()

	mydir, err := os.Getwd()
	if nil != err {
		panic(err)
	}

	if *help {
		flag.Usage()
		err := errors.New("help")
		return toTranslate, rootSource, rootDest, keyDest, err
	} else if "" == iFlag && "" == oFlag && "" != eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && !*sFlag && !*vFlag && "" == biFlag && "" == boFlag {
		toTranslate = options.ExecBasm
		rootDest = mydir + "/" + eFlag
	} else if "" != iFlag && "" != oFlag && "" == eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && !*sFlag && !*vFlag && "" == biFlag && "" == boFlag {
		toTranslate = options.UserTranslate
		rootSource = mydir + "/" + iFlag
		rootDest = mydir + "/" + oFlag
	} else if "" != piFlag && "" != poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && !*sFlag && !*vFlag && "" == biFlag && "" == boFlag {
		toTranslate = options.Primitive
		rootSource = mydir + "/" + piFlag
		rootDest = mydir + "/" + poFlag
	} else if "" == piFlag && "" == poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" != peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && !*sFlag && !*vFlag && "" == biFlag && "" == boFlag {
		toTranslate = options.InterpPrimitive
		rootDest = mydir + "/" + peFlag
	} else if "" == piFlag && "" == poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" == peFlag &&
		"" != ciFlag && "" != coFlag && "" == ceFlag && "" != kFlag && !*sFlag && !*vFlag && "" == biFlag && "" == boFlag {
		toTranslate = options.Encrypt
		rootSource = mydir + "/" + ciFlag
		rootDest = mydir + "/" + coFlag
		keyDest = mydir + "/" + kFlag
	} else if "" == piFlag && "" == poFlag && "" == iFlag && "" == oFlag && "" == eFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" != ceFlag && "" != kFlag && !*sFlag && !*vFlag && "" == biFlag && "" == boFlag {
		toTranslate = options.ExecEncrypt
		rootDest = mydir + "/" + ceFlag
		keyDest = mydir + "/" + kFlag
	} else if "" == iFlag && "" == oFlag && "" == eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && *sFlag && !*vFlag && "" == biFlag && "" == boFlag {
		toTranslate = options.Internal
	} else if "" == iFlag && "" == oFlag && "" == eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && !*sFlag && *vFlag && "" == biFlag && "" == boFlag {
		fmt.Println("bint version 2.2")
		os.Exit(0)
	} else if "" == iFlag && "" == oFlag && "" == eFlag && "" == piFlag && "" == poFlag && "" == peFlag &&
		"" == ciFlag && "" == coFlag && "" == ceFlag && "" == kFlag && !*sFlag && !*vFlag && "" != biFlag && "" != boFlag {
		toTranslate = options.Compile
		rootSource = mydir + "/" + biFlag
		rootDest = mydir + "/" + boFlag
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
			rootSource = "benv/if.b"
			rootDest = "benv/if.basm"
			if execBenv {
				filesListToExecute = []string{"benv/internal/build/import",
					"benv/internal/build/prep_func",
					"benv/internal/build/prep_for",
					"benv/internal/build/prep_dowhile",
					"benv/internal/build/dowhile",
					"benv/internal/build/prep_while",
					"benv/internal/build/prep_if",
					"benv/internal/build/prep_try",
					"benv/internal/build/try",
					"benv/internal/build/while",
					"benv/internal/build/for",
					"benv/internal/build/if_plug",
					"benv/internal/build/nested_call",
					"benv/internal/build/long_function",
					"benv/internal/build/recurs",
					"benv/internal/build/func",
					"benv/internal/build/slice",
					"benv/internal/build/int",
					"benv/internal/build/float",
					"benv/internal/build/bool",
					"benv/internal/build/len",
					"benv/internal/build/str",
					"benv/internal/build/index",
					"benv/internal/build/is_letter",
					"benv/internal/build/is_digit",
					"benv/internal/build/reg_find",
					"benv/internal/build/exists",
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
					"benv/build/prep_for",
					"benv/build/prep_dowhile",
					"benv/build/dowhile",
					"benv/build/prep_while",
					"benv/build/prep_if",
					"benv/build/prep_try",
					"benv/build/try",
					"benv/build/while",
					"benv/build/for",
					"benv/build/if",
					"benv/build/nested_call",
					"benv/build/long_function",
					"benv/build/recurs",
					"benv/build/func",
					"benv/build/slice",
					"benv/build/int",
					"benv/build/float",
					"benv/build/bool",
					"benv/build/len",
					"benv/build/str",
					"benv/build/index",
					"benv/build/is_letter",
					"benv/build/is_digit",
					"benv/build/reg_find",
					"benv/build/exists",
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
			rootSource = "benv/build/if.basm"
			rootDest = "benv/build/main.go"

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
			//rootSource = "long_function_program.b"
			//rootDest = "benv/dowhile.basm"
			rootDest = "benv/try.basm"
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
		} else if options.Compile == toTranslate {
			rootSource = "program.basm"
			rootDest = "program"
		} else {
			panic(errors.New("set option to translate"))
		}

	} else if options.UserTranslate == toTranslate {
		ex, err := os.Executable()
		if nil != err {
			panic(err)
		}
		exPath := filepath.Dir(ex) + "/"
		filesListToExecute = []string{exPath + "benv/build/import",
			exPath + "benv/build/prep_func",
			exPath + "benv/build/prep_for",
			exPath + "benv/build/prep_dowhile",
			exPath + "benv/build/dowhile",
			exPath + "benv/build/prep_while",
			exPath + "benv/build/prep_if",
			exPath + "benv/build/while",
			exPath + "benv/build/for",
			exPath + "benv/build/if",
			exPath + "benv/build/nested_call",
			exPath + "benv/build/long_function",
			exPath + "benv/build/recurs",
			exPath + "benv/build/func",
			exPath + "benv/build/slice",
			exPath + "benv/build/int",
			exPath + "benv/build/float",
			exPath + "benv/build/bool",
			exPath + "benv/build/len",
			exPath + "benv/build/str",
			exPath + "benv/build/index",
			exPath + "benv/build/is_letter",
			exPath + "benv/build/is_digit",
			exPath + "benv/build/reg_find",
			exPath + "benv/build/exists",
			exPath + "benv/build/print_format"}
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
func ExecBenv(filesListToExecute []string, rootSource string, rootDest string, toTranslate int) {
	FileToExecute = filesListToExecute[0]
	bufRootSource := rootSource
	userDest := rootDest

	if options.UserValidate == toTranslate || options.InternalValidate == toTranslate {
		bufRootSource += "v"
	}
	ex, err := os.Executable()
	if nil != err {
		panic(err)
	}
	exPath := filepath.Dir(ex) + "/"

	if !Exists("benv/build/import") {
		err = RemoveContents(exPath)
		if nil != err {
			fmt.Println(err)
		}
		err = CopyDirectory(rootSource[:strings.LastIndex(rootSource, "/")], exPath)
		//_, err = Copy(rootSource, exPath+rootSource[strings.LastIndex(rootSource, "/")+1:])
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		rootSource = rootSource[strings.LastIndex(rootSource, "/")+1:]
		rootDest = rootDest[strings.LastIndex(rootDest, "/")+1:]

		bufRootSource = rootSource

		if options.UserValidate == toTranslate || options.InternalValidate == toTranslate {
			bufRootSource += "v"
		}
	}

	cmd := exec.Command(filesListToExecute[0], "-i", bufRootSource)

	err = cmd.Start()
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

	if len(filesListToExecute) > 2 && !Exists("benv/build/import") {
		_, err = Copy(exPath+rootDest, userDest)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
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
	var sourceNewChunk func() (string, error)
	var wasGetCommandCounterByMark bool
	var newKey func() (string, error)

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

	if options.InternalValidate == toTranslate {
		filesListToExecute = []string{"benv/internal/build/import",
			"benv/internal/build/trace",
		}
	}
	var err error

	if options.UserValidate == toTranslate {
		if Exists("benv/build/import") {
			filesListToExecute = []string{"benv/build/import",
				"benv/build/trace",
			}
		} else {
			ex, err := os.Executable()
			if nil != err {
				panic(err)
			}
			exPath := filepath.Dir(ex) + "/"

			filesListToExecute = []string{exPath + "benv/build/import",
				exPath + "benv/build/trace",
			}
		}
	}

	if options.Internal != toTranslate {
		sysMod = toTranslate
	}

	if (options.Internal == toTranslate && (options.UserTranslate == sysMod || options.Internal == sysMod) && execBenv) ||
		options.UserTranslate == toTranslate || options.UserValidate == toTranslate || options.InternalValidate == toTranslate {
		ExecBenv(filesListToExecute, rootSource, rootDest, toTranslate)
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
			for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
				if nil != err {
					panic(err)
				}
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
		for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
			if nil != err {
				panic(err)
			}
			//ввод команд
			CommandToExecute = strings.TrimSpace(chunk)
			inputedCode = CodeInput(chunk, !wasGetCommandCounterByMark)

			if options.ExecEncrypt == sysMod {
				k, err := newKey()
				if nil != err {
					panic(err)
				}
				inputedCode = decryptor.Decrypt(inputedCode, k)
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
						variables, options.Transpile == sysMod, false, transpileDest, options.Primitive == sysMod, primitiveDest, nil, nil, nil)
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
				_, infoListList, systemStack, err = parser.Parse(exprList, variables, systemStack, options.HideTree,
					options.Transpile == sysMod, options.Primitive == sysMod, primitiveDest, transpileDest, nil, nil)
				if nil != err {
					panic(err)
				}
			} else if options.InterpPrimitive != sysMod && options.ExecEncrypt != sysMod && 0 == exprList[0][1] {
				continue
			} else if options.InterpPrimitive == sysMod || options.ExecEncrypt == sysMod {
				if isBasmStyle {
					exprList, variables, err = LexicalAnalyze(inputedCode,
						variables, options.Transpile == sysMod, false, transpileDest, options.Primitive == sysMod, primitiveDest, nil, nil, nil)
					_, infoListList, systemStack, err = parser.Parse(exprList, variables, systemStack, options.HideTree,
						options.Transpile == sysMod, options.Primitive == sysMod, primitiveDest, transpileDest, nil, nil)
					if nil != err {
						panic(err)
					}
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
								code, err = sourceNewChunk()
								if nil != err {
									panic(err)
								}
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

								if options.UserValidate == toTranslate || options.InternalValidate == toTranslate {
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

	if options.Compile == sysMod {
		compilerVars.CompilerVars = map[string]string{}
		compilerVars.CompilerVars["error"] = "Error"
		compilerVars.CompilerVars["$toPanic"] = "Panic"

		_, variables, err = LexicalAnalyze("stringerror", variables, false, false, nil, false, nil, nil, nil, nil)

		if nil != err {
			panic(err)
		}

		_, variables, err = LexicalAnalyze("bool$toPanic", variables, false, false, nil, false, nil, nil, nil, nil)

		if nil != err {
			panic(err)
		}

		dataFile, err := compiler.InitData()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

		progFile, err := compiler.InitProg()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

		labelsFile, err := compiler.InitLabels()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

		f, err := os.Open(rootSource)
		if nil != err {
			fmt.Println("could not open file " + rootSource)
			os.Exit(1)
		}
		newChunk := EachChunk(f)

		for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
			if nil != err {
				panic(err)
			}
			inputedCode := CodeInput(chunk, false)
			exprList, variables, err = LexicalAnalyze(inputedCode, variables, false, true, nil, false, nil, dataFile, labelsFile, progFile)
			if nil != err {
				panic(err)
			}
			_, infoListList, systemStack, err =
				parser.Parse(exprList, variables, systemStack, options.HideTree, false, false, nil, nil, dataFile, progFile)
			if nil != err {
				panic(err)
			}
			for _, infoList := range infoListList {
				if "0" != infoList[0] {
					compiler.CompileTree(infoList, variables, systemStack, dataFile, progFile)
				}
			}
		}
		err = compiler.FinishData(dataFile)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		err = compiler.FinishProg(progFile)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

		err = compiler.FinishLabels(labelsFile)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		err = dataFile.Close()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		err = progFile.Close()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		err = labelsFile.Close()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

		err = compiler.CompileAsm(rootDest)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

	}
}
