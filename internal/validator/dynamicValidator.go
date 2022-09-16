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
var fileName string
var sourceFile string
var retVal string

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
	tail, stat := check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){`, command)
	if status.Yes == stat {
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

		return tail, stat, variables
	}
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`(\){)`, command)
	if status.Yes == stat {
		variables = append(variables, [][]interface{}{})
		closureHistory = append(closureHistory, brace{funcDefinition, LineCounter, CommandToExecute})
		reg, err := regexp.Compile(`^(?:(int|float|bool|string|stack|void))`)
		if nil != err {
			panic(err)
		}
		loc := reg.FindIndex([]byte(command))
		retVal = command[loc[0]:loc[1]]

		return tail, stat, variables
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

func dynamicValidateCommand(command string, variables [][][]interface{}) error {
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

	tail, stat, variables = dValidateReturn(command, variables)

	if status.Yes == stat {
		return nil
	}

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
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
