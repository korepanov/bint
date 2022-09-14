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

	if 1 == len(infoListList) {
		res = infoListList[0]
	} else {
		res, _, _ = executor.ExecuteTree(infoListList[0], allVariables, nil, false, false, nil, nil)
	}

	fmt.Println(res[0])
	return "bool"
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

	/*for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
		CommandToExecute = strings.TrimSpace(chunk)
		errorFile = CodeInput(chunk, false)

		if filter(errorFile) {
			COMMAND_COUNTER--
			newChunk, err = SetCommandCounter(f, COMMAND_COUNTER)
			if nil != err {
				panic(err)
			}
		} else {
			if len(errorFile) > 6 && "$file$" == errorFile[0:6]{
				errorFile = errorFile[6:]
				break
			}
		}
	}*/
	fmt.Println(errorFile)
	if "$trace" != inputedCode[0:6] {
		sourceCommandCounter = 1
	} else {
		sourceCommandCounter, err = strconv.Atoi(inputedCode[6 : len(inputedCode)-1])
		if nil != err {
			sourceCommandCounter = 1
		}
	}

	fmt.Println(sourceCommandCounter)
	err = f.Close()
	if nil != err {
		panic(err)
	}

	/*f, err = os.Open(sourceFile)
	if nil != err {
		panic(err)
	}
	newChunk, err = SetCommandCounter(f, sourceCommandCounter)
	if nil != err{
		panic(err)
	}
	chunk := newChunk()

	fmt.Println("ERROR in " + rootSource + " at near line " +
		fmt.Sprintf("%v", serviceTools.LineCounter))
	fmt.Println(serviceTools.CommandToExecute)
	fmt.Println(err)*/

	err = f.Close()
	if nil != err {
		panic(err)
	}
}
func dValidateFuncDefinition(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	var retVal string

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
		fmt.Println(retVal)

		tail, stat = check(`^(?:(int|float|bool|string|stack|void)[[:alnum:]|_]+)`, command)
		if status.Yes != stat {
			handleError("err")
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
		fmt.Println(retVal)

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
