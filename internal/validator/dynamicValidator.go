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
var forCounter int

func checkVars(exprList interface{}, allVariables [][]interface{}) {
	var isVar bool

	for _, expr := range exprList.([][]interface{}) {
		if "SUBEXPR" == fmt.Sprintf("%v", expr[0]) {
			checkVars(expr, allVariables)
		}
		if "VAR" == fmt.Sprintf("%v", expr[0]) {
			newVariable := EachVariable(allVariables)
			for v := newVariable(); "end" != v[0]; v = newVariable() {
				if fmt.Sprintf("%v", expr[1]) == fmt.Sprintf("%v", v[1]) {
					isVar = true
					break
				}
			}
			if !isVar {
				handleError("unresolved reference: " + fmt.Sprintf("%v", expr[1]))
			} else {
				isVar = false
			}

		}

	}
}

func getExprEnd(command string, startPos int) int {
	brOpened := 1
	brClosed := 0
	pos := startPos + 1

	for i := pos; i < len(command); i++ {
		if "(" == string(command[i]) {
			brOpened++
		}
		if ")" == string(command[i]) {
			brClosed++
		}

		if brOpened == brClosed {
			return i + 1
		}
	}

	panic(errors.New("invalid brace number"))
}

func getExprType(command string, variables [][][]interface{}) string {
	var exprList [][]interface{}
	var err error

	var allVariables [][]interface{}

	for _, v := range variables {
		allVariables = append(allVariables, v...)
	}

	exprList, allVariables, err = lexer.LexicalAnalyze(command,
		allVariables, false, nil, false, nil)
	if nil != err {
		handleError(err.Error())
	}

	checkVars(exprList, allVariables)

	_, infoListList, _ := parser.Parse(exprList, allVariables, nil, false, false, false, nil, nil)

	var res []interface{}

	if 1 == len(infoListList[0]) {
		res = infoListList[0]

		newVariable := EachVariable(allVariables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if fmt.Sprintf("%v", res[0]) == fmt.Sprintf("%v", v[1]) {
				return fmt.Sprintf("%v", v[0])
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

func dValidateString(command string) (string, int) {
	tail := command
	re, err := regexp.Compile(`"(\\.|[^"])*"`)
	if nil != err {
		panic(err)
	}
	if nil != re.FindIndex([]byte(tail)) {
		tail = string(re.ReplaceAll([]byte(command), []byte(`$$val`)))
		return tail, status.Yes
	}

	return tail, status.No
}

func dValidateStr(command string, variables [][][]interface{}) (string, [][][]interface{}) {
	tail := command
	re, err := regexp.Compile(`str\(`)
	if nil != err {
		panic(err)
	}
	if nil != re.FindIndex([]byte(tail)) {
		var poses [][]int
		poses = re.FindAllIndex([]byte(tail), -1)
		var replacerArgs []string
		for _, pos := range poses {
			exprEnd := getExprEnd(tail, pos[1]-1)
			t := getExprType(tail[pos[0]+4:exprEnd-1], variables)
			if "stack" == t {
				handleError("data type mismatch in str: stack")
			}
			replacerArgs = append(replacerArgs, tail[pos[0]:exprEnd])
			replacerArgs = append(replacerArgs, `$val`)
		}

		r := strings.NewReplacer(replacerArgs...)
		tail = r.Replace(tail)
	}
	return tail, variables
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
	tail2, stat2 := check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`(\){)`, command)
	if "" != tail2 {
		tail = tail2
	}
	if status.Yes == stat2 {
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
			var err error
			if "void" != retVal {
				_, variables[0], err = lexer.LexicalAnalyze(funcTable[funcName]+"$"+funcName,
					variables[0], false, nil, false, nil)
			} else {
				variables[0] = append(variables[0], []interface{}{"void", funcName, []interface{}{"func"}})
			}

			if nil != err {
				handleError(err.Error())
			}

			if "float" == funcTable[funcName] {
				variables[0][len(variables[0])-1][2] = "0.1"
			}
			if "int" == funcTable[funcName] {
				variables[0][len(variables[0])-1][2] = "1"
			}
		}
	}

	if status.Yes == stat || status.Yes == stat2 {
		stat = status.Yes
	} else {
		stat = status.No
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
		variables = append(variables, [][]interface{}{})
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

func dValidateElseIf(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`(?:^}elseif\([^{]+\){)`, command)

	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{elseIfCond, LineCounter, CommandToExecute})
		variables = variables[:len(variables)-1]
		variables = append(variables, [][]interface{}{})
		re, err := regexp.Compile(`(?:^}elseif\([^{]+\){)`)
		if nil != err {
			panic(err)
		}
		loc := re.FindIndex([]byte(command))
		elseIfStruct := command[:loc[1]]

		if "bool" != getExprType(elseIfStruct[7:len(elseIfStruct)-1], variables) {
			handleError("the expression inside if is not a boolean expression")
		}

	}

	return tail, stat, variables
}

func dValidateElse(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`(?:^}else{)`, command)

	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{elseCond, LineCounter, CommandToExecute})
		variables = variables[:len(variables)-1]
		variables = append(variables, [][]interface{}{})
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
		tail, _, variables = dValidateFuncCall(tail, variables)

		exprType := getExprType(tail, variables)

		if retVal != exprType {
			handleError("data type mismatch in func definition and return statement: " + retVal + " and " + exprType)
		}
	}
	return "", stat, variables
}

func dValidateFuncCall(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	var replacerArgs []string
	var thisFuncName string

	tail := command

	for funcName := range funcTable {
		locArr := GetFuncNameEntry(funcName, tail)

		for _, loc := range locArr {
			thisFuncName = funcName
			loc[1] = getExprEnd(tail, loc[1])

			replacerArgs = append(replacerArgs, tail[loc[0]:loc[1]])
			replacerArgs = append(replacerArgs, "$"+funcName)
		}

	}

	if nil != replacerArgs {
		r := strings.NewReplacer(replacerArgs...)
		tail = r.Replace(tail)
		if tail[1:] == thisFuncName && "void" == funcTable[thisFuncName] {
			return ``, status.Yes, variables
		}
		return tail, status.Yes, variables
	}

	return tail, status.No, variables
}

func dValidatePrint(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`(?:print\()`, command)
	if status.Yes == stat {
		T := getExprType(tail[:len(tail)-1], variables)
		if "string" != T {
			handleError("print: data type mismatch: string and " + T)
		}
		return "", stat, variables
	}
	return command, status.No, variables
}

func dValidateAssignment(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	_, stat := check(`(?:[[:alpha:]][[:alnum:]|_]*={1}[^=]+)`, command)
	if status.Yes == stat {

		thisVar := command[0:strings.Index(command, "=")]
		expr := command[strings.Index(command, "=")+1:]

		var allVariables [][]interface{}

		for _, v := range variables {
			allVariables = append(allVariables, v...)
		}

		newVariable := EachVariable(allVariables)
		for v := newVariable(); "end" != v[0]; v = newVariable() {
			if thisVar == fmt.Sprintf("%v", v[1]) {
				T := getExprType(expr, variables)
				if v[0] != T {
					handleError("data type mismatch: " + fmt.Sprintf("%v", v[0]) + " and " + T)
				} else {
					return ``, status.Yes, variables
				}
			}
		}
		handleError("unresolved reference: " + thisVar)
	}

	return ``, status.No, variables
}

func dValidateWhile(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`(?:^while\([^{]+\){)`, command)

	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{ifCond, LineCounter, CommandToExecute})
		variables = append(variables, [][]interface{}{})
		re, err := regexp.Compile(`(?:^while\([^{]+\){)`)
		if nil != err {
			panic(err)
		}
		loc := re.FindIndex([]byte(command))
		whileStruct := command[:loc[1]]

		if "bool" != getExprType(whileStruct[5:len(whileStruct)-1], variables) {
			handleError("the expression inside while is not a boolean expression")
		}
	}
	return tail, stat, variables
}
func dValidateBreak(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	var isLoop bool
	tail, stat := check(`^break`, command)
	if "" == tail && status.Yes == stat {
		for _, brace := range closureHistory {
			if loop == brace.T {
				isLoop = true
				break
			}
		}
		if !isLoop {
			handleError("break is outside of loop")
		}
	}
	return tail, stat, variables
}

func dValidateDoWhile(command string, variables [][][]interface{}) (string, int, [][][]interface{}) {
	tail, stat := check(`^do{`, command)
	if stat == status.Yes {
		closureHistory = append(closureHistory, brace{loop, LineCounter, CommandToExecute})
		variables = append(variables, [][]interface{}{})
		return tail, stat, variables
	}
	tail, stat = check(`(?:^}while\([^{]+\))`, command)

	if status.Yes == stat {
		variables = variables[:len(variables)-1]
		re, err := regexp.Compile(`(?:^}while\([^{]+\))`)
		if nil != err {
			panic(err)
		}
		loc := re.FindIndex([]byte(command))
		doWhileStruct := command[:loc[1]]

		if "bool" != getExprType(doWhileStruct[6:], variables) {
			handleError("the expression inside while is not a boolean expression")
		}

	}

	return tail, stat, variables
}

func dValidateFor(command string, variables [][][]interface{}) (string, [][][]interface{}) {
	tail, stat := check(`^for\(`, command)
	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{loop, LineCounter, CommandToExecute})
		variables = append(variables, [][]interface{}{})
		forCounter = 1
		return tail, variables
	}
	if 1 == forCounter {
		forCounter++
		return command, variables
	}
	if 2 == forCounter {
		forCounter++
		t := getExprType("("+command+")", variables)
		if "bool" != t {
			handleError("data type mismatch in for: bool and " + t)
		}
		return ``, variables
	}
	if 3 == forCounter {
		forCounter = 0

		re, err := regexp.Compile(`\){`)

		loc := re.FindIndex([]byte(command))

		if nil == loc {
			handleError("invalid for syntax")
		}

		tempHistory := closureHistory
		variables, err = dynamicValidateCommand(command[0:loc[0]], variables)
		closureHistory = tempHistory
		if nil != err {
			return command, variables
		}
		tail = command[loc[1]:]

		return tail, variables
	}

	return command, variables
}

func dynamicValidateCommand(command string, variables [][][]interface{}) ([][][]interface{}, error) {
	if isFunc && "void" != retVal && !wasRet && len(closureHistory) < 1 {
		COMMAND_COUNTER = funcCommandCounter
		handleError("missing return statement in function")
	}

	if isFunc && len(closureHistory) < 1 {
		retVal = ""
		isFunc = false
		wasRet = false
	}

	command, _ = dValidateString(command)
	command, variables = dValidateStr(command, variables)

	var tail string
	var stat int

	command, variables = dValidateFor(command, variables)

	if "" == command {
		return variables, nil
	}

	tail, stat, variables = dValidateFuncDefinition(command, variables)
	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	tail, stat, variables = dValidateBreak(command, variables)
	if status.Yes == stat {
		if `` == tail {
			return variables, nil
		}
	}

	tail, stat, variables = dValidateVarDef(command, variables)

	if status.Yes == stat {
		if `` == tail {
			return variables, nil
		}
	}

	tail, stat, variables = dValidateWhile(command, variables)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	tail, stat, variables = dValidateDoWhile(command, variables)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	tail, stat, variables = dValidateIf(command, variables)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	tail, stat, variables = dValidateElseIf(command, variables)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	tail, stat, variables = dValidateElse(command, variables)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	_, stat, variables = dValidateReturn(command, variables)

	if status.Yes == stat {
		return variables, nil
	}
	_, stat, variables = dValidatePrint(command, variables)

	if status.Yes == stat {
		return variables, nil
	}

	tail, stat, variables = dValidateFuncCall(command, variables)

	if status.Yes == stat {
		return dynamicValidateCommand(tail, variables)
	}

	_, stat, variables = dValidateAssignment(command, variables)
	if status.Yes == stat {
		return variables, nil
	}

	tail, stat, err := validateFigureBrace(command)
	if nil != err {
		panic(err)
	}
	if status.Yes == stat {

		variables = variables[:len(variables)-1]

		if tail == `` {
			return variables, nil
		}
	}

	tail, variables = dValidateFor(command, variables)

	return variables, errors.New("unresolved command")
}

func DynamicValidate(validatingFile string, rootSource string) {
	fileName = validatingFile
	funcTable = make(map[string]string)

	defer func() {
		if r := recover(); nil != r {
			handleError(fmt.Sprintf("%v", r))
			os.Exit(1)
		}
	}()

	var variables [][][]interface{}
	variables = append(variables, [][]interface{}{})
	var err error

	_, variables[len(variables)-1], err = lexer.LexicalAnalyze("string$val",
		variables[len(variables)-1], false, nil, false, nil)
	if nil != err {
		panic(err)
	}
	variables[0][0][2] = "val"

	sourceFile = rootSource
	COMMAND_COUNTER = 1

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
			variables, err = dynamicValidateCommand(inputedCode, variables)
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
