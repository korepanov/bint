package validator

import (
	"bint.com/internal/const/status"
	"bint.com/internal/executor"
	"bint.com/internal/lexer"
	"bint.com/internal/parser"
	. "bint.com/pkg/serviceTools"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var COMMAND_COUNTER int
var sourceCommandCounter int
var funcCommandCounter int
var fileName string
var sourceFile string
var retVal string
var isFunc bool
var wasRet bool
var funcTable map[string]string

func getExprType(command string, variables [][][]interface{}) string {
	var exprList [][]interface{}
	var err error

	exprList, variables[len(variables)-1], err = lexer.LexicalAnalyze(command,
		variables[len(variables)-1], false, nil, false, nil)
	if nil != err {
		handleError(err.Error())
	}
	var allVariables [][]interface{}

	for _, v := range variables {
		allVariables = append(allVariables, v...)
	}

	_, infoListList, _ := parser.Parse(exprList, allVariables, nil, false, false, false, nil, nil)

	var res []interface{}

	if 1 == len(infoListList[0]) {
		res = infoListList[0]

		newVariable := EachVariable(allVariables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", res[0]) == fmt.Sprintf("%v", v[1]) {
				res[0] = ValueFoldInterface(v[2])
				break
			}
		}
	} else {
		res, _, _ = executor.ExecuteTree(infoListList[0], allVariables, nil, false, false, nil, nil)
	}

	return WhatsType(fmt.Sprintf("%v", res[0]))
}

func filter(command string) bool {
	if "$" == string(command[0]) && "$" == string(command[len(command)-1]) {
		return false
	}
	return true
}

func handleError(errMessage string) {
	var inputedCode string
	var errorFile string

	f, err := os.Open(fileName)
	if nil != err {
		panic(err)
	}
	COMMAND_COUNTER--
	newChunk, err := SetCommandCounter(f, COMMAND_COUNTER)

	if nil != err {
		panic(err)
	}

	for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
		CommandToExecute = strings.TrimSpace(chunk)
		inputedCode = CodeInput(chunk, false)

		if filter(inputedCode) {
			COMMAND_COUNTER--
			newChunk, err = SetCommandCounter(f, COMMAND_COUNTER)
			if nil != err {
				panic(err)
			}
		} else {
			break
		}
	}

	for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
		CommandToExecute = strings.TrimSpace(chunk)
		errorFile = CodeInput(chunk, false)

		if filter(errorFile) {
			COMMAND_COUNTER--
			newChunk, err = SetCommandCounter(f, COMMAND_COUNTER)
			if nil != err {
				panic(err)
			}
		} else {
			if len(errorFile) > 5 && "$file" == errorFile[0:5] {
				errorFile = errorFile[5 : len(errorFile)-1]
				break
			}
		}
	}

	if "$trace" != inputedCode[0:6] {
		sourceCommandCounter = 1
	} else {
		sourceCommandCounter, err = strconv.Atoi(inputedCode[6 : len(inputedCode)-1])
		if nil != err {
			sourceCommandCounter = 1
		}
	}

	f, err = os.Open(sourceFile)
	if nil != err {
		panic(err)
	}

	var chunk string

	newChunk = EachChunk(f)

	if nil != err {
		panic(err)
	}

	for counter := 0; counter < sourceCommandCounter; counter++ {
		chunk = newChunk()
		CodeInput(chunk, true)
	}

	fmt.Println("ERROR in " + errorFile + " at near line " +
		fmt.Sprintf("%v", LineCounter))
	fmt.Println(strings.TrimSpace(chunk))
	fmt.Println(errMessage)

	err = f.Close()
	if nil != err {
		panic(err)
	}

	os.Exit(1)
}
func dValidateFuncDefinition(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	var wasDef bool

	tail, stat := check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){`, command)
	if status.Yes == stat {
		funcCommandCounter = COMMAND_COUNTER
		isFunc = true
		wasDef = true

		closureHistory = append(closureHistory, brace{funcDefinition, LineCounter, CommandToExecute})
		reg, err := regexp.Compile(`^(?:(int|float|bool|string|stack|void))`)
		if nil != err {
			panic(err)
		}
		loc := reg.FindIndex([]byte(command))
		retVal = command[loc[0]:loc[1]]

		tail, stat = check(`^(?:(int|float|bool|string|stack|void)[[:alnum:]|_]+)`, command)
		if status.Yes != stat {
			handleError("is not valid func definition")
		}

		reg, err = regexp.Compile(`(?:(int|float|bool|string|stack|void)[[:alnum:]|_]+)`)
		if nil != err {
			panic(err)
		}

		locs := reg.FindAllIndex([]byte(tail), -1)
		variables = append(variables, [][]interface{}{})

		for _, loc := range locs {
			_, variables[len(variables)-1], err = lexer.LexicalAnalyze(tail[loc[0]:loc[1]],
				variables[len(variables)-1], false, nil, false, nil)
		}

		tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
			`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){`, command)
	}
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`(\){)`, command)
	if status.Yes == stat {
		funcCommandCounter = COMMAND_COUNTER
		isFunc = true
		wasDef = true
		variables = append(variables, [][]interface{}{})
		closureHistory = append(closureHistory, brace{funcDefinition, LineCounter, CommandToExecute})
		reg, err := regexp.Compile(`^(?:(int|float|bool|string|stack|void))`)
		if nil != err {
			panic(err)
		}
		loc := reg.FindIndex([]byte(command))
		retVal = command[loc[0]:loc[1]]
	}

	if wasDef {
		funcName, stat := check(`^(?:(int|float|bool|string|stack|void))`, command)
		if status.Yes != stat {
			handleError("is not valid func definition")
		}

		funcName = funcName[0:strings.Index(funcName, "(")]

		if "" != funcTable[funcName] {
			handleError("function polymorphism is not allowed")
		} else {
			funcTable[funcName] = retVal
		}
	}

	return tail, stat, variables
}
func dValidateVarDef(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`(?m)(?:(int|float|bool|string|stack)[[:alnum:]|_]*)`, command)
	if status.Yes == stat && `` == tail {
		var err error
		_, variables[len(variables)-1], err = lexer.LexicalAnalyze(command,
			variables[len(variables)-1], false, nil, false, nil)
		if nil != err {
			panic(err)
		}
	}
	return tail, stat, variables
}

func dValidateIf(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`(?:^if\([^{]+\){)`, command)

	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{ifCond, LineCounter, CommandToExecute})
		re, err := regexp.Compile(`(?:^if\([^{]+\){)`)
		if nil != err {
			panic(err)
		}
		loc := re.FindIndex([]byte(command))
		ifStruct := command[:loc[1]]

		if "bool" != getExprType(ifStruct[2:len(ifStruct)-1], variables) {
			handleError("the expression inside if is not a boolean expression")
		}
	}
	return tail, stat, variables
}

func dValidateReturn(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`(?m)(?:return)`, command)
	if status.Yes == stat {
		wasRet = true

		if len(closureHistory) < 1 {
			handleError("illegal position of return")
		}

		exprType := getExprType(tail, variables)

		if retVal != exprType {
			handleError("data type mismatch in func definition and return statement: " + retVal + " and " + exprType)
		}
	}
	return "", stat, variables
}

func dValidateFuncCall(command string, variables [][][]interface{}, isFunc bool) (string, int, [][][]interface{}) {
	tail := command
	re, err := regexp.Compile(`(?:[)][[:alpha:]|_]+[[:alnum:]|_]*\((([[:alnum:]]*\,){0,})[[:alnum:]]*\))`)
	if nil != err {
		panic(err)
	}
	notFuncLocArr := re.FindAllIndex([]byte(tail), -1)

	for _, loc := range notFuncLocArr {
		loc[0]++
	}

	re, err = regexp.Compile(`(?:[[:alpha:]|_]+[[:alnum:]|_|\.]*\((([[:alnum:]|_|\[|\]|:|\.]*\,[^,]){0,})[[:alnum:]|_|\[|\]|:|\.]*\))`)
	if nil != err {
		panic(err)
	}

	locArr := re.FindAllIndex([]byte(tail), -1)

	var funcLocArr [][]int
	toAppend := true

	for _, loc := range locArr {
		for _, notFuncLoc := range notFuncLocArr {
			if (loc[0] == notFuncLoc[0]) && (loc[1] == notFuncLoc[1]) {
				toAppend = false
			}
		}
		if toAppend {
			funcLocArr = append(funcLocArr, loc)
		}

		toAppend = true
	}

	var replacerArgs []string

	for _, loc := range funcLocArr {
		funcName := tail[loc[0]:loc[1]]
		funcName = funcName[:strings.Index(funcName, "(")]

		if len(funcName) > 6 && "return" == funcName[0:6] {
			funcName = funcName[6:]
			loc[0] += 6
		}

		replacerArgs = append(replacerArgs, tail[loc[0]:loc[1]])

		_, variables[len(variables)-1], err = lexer.LexicalAnalyze(funcTable[funcName]+"$"+funcName,
			variables[len(variables)-1], false, nil, false, nil)
		if nil != err {
			handleError(err.Error())
		}

		if "float" == funcTable[funcName] {
			variables[len(variables)-1][len(variables[len(variables)-1])-1][2] = "0.1"
		}
		replacerArgs = append(replacerArgs, "$"+funcName)
	}

	if nil != replacerArgs {
		r := strings.NewReplacer(replacerArgs...)
		tail = r.Replace(tail)
		return dValidateFuncCall(tail, variables, true)
	}

	if isFunc {
		return tail, status.Yes, variables
	}

	return tail, status.No, variables
}

func dynamicValidateCommand(command string, variables [][][]interface{}) error {

	if isFunc && "void" != retVal && !wasRet && len(closureHistory) < 1 {
		COMMAND_COUNTER = funcCommandCounter
		handleError("missing return statement in function")
	}

	if isFunc && len(closureHistory) < 1 {
		retVal = ""
		isFunc = false
		wasRet = false
	}

	tail, stat, variables := dValidateFuncDefinition(command, variables)
	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	tail, stat, variables = dValidateVarDef(command, variables)

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, variables = dValidateIf(command, variables)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	tail, stat, variables = dValidateFuncCall(command, variables, false)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	_, stat, variables = dValidateReturn(command, variables)

	if status.Yes == stat {
		return nil
	}

	tail, stat, err := validateFigureBrace(command)
	if nil != err {
		panic(err)
	}
	if status.Yes == stat {
		variables = variables[:len(variables)-1]
		if tail == `` {
			return nil
		}
	}

	return errors.New("unresolved command")
}

func DynamicValidate(validatingFile string, rootSource string) {
	funcTable = make(map[string]string)

	defer func() {
		if r := recover(); nil != r {
			handleError(fmt.Sprintf("%v", r))
			os.Exit(1)
		}
	}()

	sourceFile = rootSource
	COMMAND_COUNTER = 1
	var variables [][][]interface{}

	variables = append(variables, [][]interface{}{})
	fileName = validatingFile

	f, err := os.Open(validatingFile)

	if nil != err {
		panic(err)
	}
	newChunk := EachChunk(f)

	for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
		CommandToExecute = strings.TrimSpace(chunk)
		inputedCode := CodeInput(chunk, false)
		COMMAND_COUNTER++

		if filter(inputedCode) {
			err = dynamicValidateCommand(inputedCode, variables)
			if nil != err {
				handleError(err.Error())
			}
		}
	}

	err = f.Close()
	if nil != err {
		panic(err)
	}
}
