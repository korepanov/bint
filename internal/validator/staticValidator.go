package validator

import (
	"bint.com/internal/const/status"
	. "bint.com/pkg/serviceTools"
	"errors"
	"os"
	"regexp"
	"strings"
)

const (
	funcDefinition = iota
	ifCond         = iota
	elseIfCond     = iota
	elseCond       = iota
	loop           = iota
)

type brace struct {
	T       int
	line    int
	command string
}

var closureHistory []brace

func check(reg string, command string) (tail string, stat int) {
	re, err := regexp.Compile(reg)
	if nil != err {
		panic(err)
	}
	loc := re.FindIndex([]byte(command))

	if nil != loc && 0 == loc[0] {
		tail = command[loc[1]:]
		return tail, status.Yes
	}
	return ``, status.No
}

func isValidBracesNum(command string) bool {
	if strings.Count(command, "(") != strings.Count(command, ")") {
		return false
	}
	return true
}

func validateStandardFuncCall(command string, funcName string, argNum int,
	randomLocation bool) (tail string, stat int, err error) {

	if !randomLocation {
		if 0 == argNum {
			tail, stat = check(`(?:`+funcName+`\(\))`, command)
			if status.Yes == stat && `` == tail {
				return tail, stat, nil
			}
			tail, stat = check(`(?:`+funcName+`\(.+\))`, command)
			if status.Yes == stat && `` == tail {
				return tail, status.Err, errors.New(funcName + ` must have no arguments`)
			}

		} else if 1 == argNum {
			tail, stat = check(`(?:`+funcName+`\([^\,]*\))`, command)
			if status.Yes == stat && `` == tail {
				return tail, stat, nil
			}
			tail, stat = check(`(?:(`+funcName+`)\(.+\,.*\))`, command)
			if status.Yes == stat && `` == tail {
				return tail, status.Err, errors.New(funcName + ` must have 1 argument`)
			}

			tail, stat = check(`(?:`+funcName+`\(\))`, command)
			if status.Yes == stat && `` == tail {
				return ``, status.Err, errors.New(`input must have 1 argument`)
			}
		} else {
			panic(`validateStandardFuncCall: unsupported case`)
		}
	} else {
		if 0 == argNum {
			panic(`validateStandardFuncCall: unsupported case`)
		} else if 1 == argNum {
			tail = command

			tail, err = ReplaceFunc(`(?:`+funcName+`\([[:alpha:]]+[[:alnum:]|_|\.]*\))`, tail)
			if nil != err {
				panic(err)
			}

			moreArgs, err := CheckEntry(`(?:(`+funcName+`)\(.+\,.*\))`, tail)

			if nil != err {
				panic(err)
			}

			noArgs, err := CheckEntry(`(?:`+funcName+`\(\))`, tail)

			if nil != err {
				panic(err)
			}

			if moreArgs || noArgs {
				return ``, status.Err, errors.New(funcName + ` must have 1 argument`)
			}

		} else if 2 == argNum {
			tail = command
			tail, err = ReplaceFunc(`(?:`+funcName+
				`\([[:alpha:]]+[[:alnum:]|_|\[|\]|\.]*\,[[:alpha:]]+[[:alnum:]|_|\[|\]|\.]*\))`, tail)

			if nil != err {
				panic(err)
			}

			moreArgs, err := CheckEntry(`(?:`+funcName+
				`\(([[:alpha:]]+[[:alnum:]|_|\[|\]|\.]*\,)+[[:alpha:]]+[[:alnum:]|_|\[|\]|\.]*\))`, tail)

			if nil != err {
				panic(err)
			}

			oneArg, err := CheckEntry(`(?:`+funcName+`\([[:alpha:]]+[[:alnum:]|_|\[|\]|\.]*\,?\))`, tail)

			if nil != err {
				panic(err)
			}

			noArgs, err := CheckEntry(`(?:`+funcName+`\(\))`, tail)

			if nil != err {
				panic(err)
			}

			if moreArgs || oneArg || noArgs {
				return ``, status.Err, errors.New(funcName + ` must have 2 arguments`)
			}

		} else {
			panic(`validateStandardFuncCall: unsupported case`)
		}
	}

	return command, status.No, nil
}

func validateFuncDefinition(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){`, command)
	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{funcDefinition, LineCounter, CommandToExecute})
		return tail, stat, nil
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`(\){)`, command)
	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{funcDefinition, LineCounter, CommandToExecute})
		return tail, stat, nil
	}

	tail, stat = check(`(?m)(?:[[:alnum:]|_]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`missing return value in function definition`)
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`((?:((int|float|bool|string|stack|void))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack|void)[[:alnum:]|_]+\){`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`function parameter cannot be of type void`)
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\)`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`missing '{' in function definition`)
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void).*?\`+
		`((?:((int|float|bool|string|stack)).+\,){0,})(int|float|bool|string|stack).+\){`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`invalid characters in entity names`)
	}
	return ``, status.No, nil
}

func validateReturn(command string) (tail string, stat int, err error) {
	re, err := regexp.Compile(`[^=]=[^=]`)
	if nil != err {
		panic(err)
	}
	if nil != re.FindIndex([]byte(command)) {
		return command, status.No, nil
	}
	tail, stat = check(`^return.*`, command)
	if `` == tail {
		return tail, stat, nil
	} else {
		return command, status.No, nil
	}
}

func validateFigureBrace(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?m)(?:})`, command)

	if status.Yes == stat {
		if len(closureHistory) > 0 {
			closureHistory = closureHistory[:len(closureHistory)-1]
		} else {
			return ``, status.Err, errors.New(`'{' is missing`)
		}

	}

	return tail, stat, nil
}

func validateVarDef(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack)[[:alpha:]][[:alnum:]|_]*)`, command)
	if status.Yes == stat && `` == tail {
		return tail, stat, nil
	}
	tail, stat = check(`(?m)(?:(string|stack)[^(]*)`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`invalid variable name`)
	}
	tail, stat = check(`(?m)(?:(int|float|bool)[^(]*)`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`invalid variable name`)
	}
	return ``, status.No, nil
}

func validateAssignment(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:[[:alpha:]][[:alnum:]|_]*={1}[^=]+)`, command)
	if status.Yes == stat {
		return tail, stat, nil
	}

	return ``, status.No, nil
}

func validateNumber(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:\-?[[:digit:]]+\.?[[:digit:]]*)`, command)
	return tail, stat, nil
}

func validateVar(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:\(?[[:alnum:]|_|\[|\]|:]+\)?)`, command)
	return tail, stat, nil
}
func validateString(command string) (tail string, stat int, err error) {
	tail = command
	re, err := regexp.Compile(`"(\\.|[^"])*"`)
	if nil != err {
		panic(err)
	}
	if nil != re.FindIndex([]byte(tail)) {
		tail = string(re.ReplaceAll([]byte(command), []byte(`val`)))
		return tail, status.Yes, nil
	}

	return tail, status.No, nil
}

func checkComparison(command string, reg string) (isComparison bool, newCommand string) {

	re, err := regexp.Compile(reg)
	if nil != err {
		panic(err)
	}
	if nil != re.FindIndex([]byte(command)) {
		isComparison = true
	}
	newCommand = string(re.ReplaceAll([]byte(command), []byte("val")))

	return isComparison, newCommand
}

func validateComparison(command string, isOp bool) (tail string, stat int, err error) {
	var isComparison bool

	if !isOp {
		re, err := regexp.Compile(`(?m)(?:[^[[:alnum:]|_,\)]\-([[:alpha:]]([[:alnum:]|_]+)?))`)
		if nil != err {
			panic(err)
		}
		if nil != re.FindIndex([]byte(command)) {
			return ``, status.Err, errors.New(`unary minus before variable is not allowed, use expression like (-1)*var`)
		}
	}

	temp, command := checkComparison(command, `(?:[^=(]+=={1}[^=)]+)`)

	if temp {
		isComparison = true
	}

	temp, command = checkComparison(command,
		`(?m)(?:([[[:alnum:]|_])*(?:(>=|=<|>|<))([[[:alnum:]|_])*)`)

	if temp {
		isComparison = true
	}

	temp, command = checkComparison(command,
		`(?m)(?:\(([[[:alpha:]|_])+[[:alnum:]]*\)(?:(AND|OR|XOR))\(([[[:alpha:]|_])+[[:alnum:]]*\))`)

	if temp {
		isComparison = true
	}

	temp, command = checkComparison(command, `(?:NOT\([[:alpha:]]+[[:alnum:]]*)\)`)

	if temp {
		isComparison = true
	}

	if `val` == command && isOp {
		return command, status.Yes, nil
	}

	if !isComparison {

		if `(val)` == command {
			command = `val`
		}

		if isOp {
			return command, status.Yes, nil
		}
		return command, status.No, nil
	}

	return validateComparison(command, true)
}

func validateArithmetic(command string, isOp bool) (tail string, stat int, err error) {
	if !isOp {
		re, err := regexp.Compile(`(?m)(?:[^[[:alnum:]|_,\)]\-([[:alpha:]]([[:alnum:]|_]+)?))`)
		if nil != err {
			panic(err)
		}
		if nil != re.FindIndex([]byte(command)) {
			return ``, status.Err, errors.New(`unary minus before variable is not allowed, use expression like (-1)*var`)
		}
	}
	re, err := regexp.Compile(`(?m)(?:\(([[[:alnum:]|_]|\[|])*[+|\-|*|/|^|@]([[[:alnum:]|_]|\[|])*\))`)
	if nil != err {
		panic(err)
	}

	reConc, err := regexp.Compile(`(?m)(?:\(([[[:alnum:]|_]|\[|]|:)*[+]([[[:alnum:]|_]|\[|]|:)*\))`)
	if nil != err {
		panic(err)
	}

	loc := re.FindIndex([]byte(command))
	locConc := reConc.FindIndex([]byte(command))

	if nil == loc && nil == locConc {
		if isOp {
			return command, status.Yes, nil
		}
		return command, status.No, nil
	}

	command = string(re.ReplaceAll([]byte(command), []byte("val")))
	command = string(reConc.ReplaceAll([]byte(command), []byte("val")))

	return validateArithmetic(command, true)
}

func validateFuncCall(command string, isFunc bool) (tail string, stat int, err error) {
	tail = command
	re, err := regexp.Compile(`(?:[)][[:alpha:]|_]+[[:alnum:]|_]*\((([[:alnum:]]*\,){0,})[[:alnum:]]*\))`)
	if nil != err {
		panic(err)
	}
	notFuncLocArr := re.FindAllIndex([]byte(tail), -1)

	for _, loc := range notFuncLocArr {
		loc[0]++
	}
	re, err = regexp.Compile(`(AND|OR|XOR)\(`)

	for _, el := range re.FindAllIndex([]byte(tail), -1) {
		notFuncLocArr = append(notFuncLocArr, []int{el[0], FindExprEnd(tail, el[1]+1)})
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
		replacerArgs = append(replacerArgs, tail[loc[0]:loc[1]])
		replacerArgs = append(replacerArgs, "val")
	}

	if nil != replacerArgs {
		r := strings.NewReplacer(replacerArgs...)
		tail = r.Replace(tail)
		return validateFuncCall(tail, true)
	}

	if isFunc {
		return tail, status.Yes, nil
	}

	return tail, status.No, nil
}

func validateUserStackCall(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:[[:alpha:]]+[[:alnum:]|_]*\.push\(.+\))`, command)
	if status.Yes == stat && `` == tail {
		return tail, stat, err
	}

	tail, stat = check(`(?:[[:alpha:]]+[[:alnum:]|_|\(|\)|,]*\.pop\([[:alpha:]]+[[:alnum:]|_]*\))`, command)
	if status.Yes == stat && `` == tail {
		return tail, stat, err
	}

	tail, stat = check(`(?:[[:alpha:]]+[[:alnum:]|_]*\.pop\(.+\))`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`invalid pop argument`)
	}

	tail, stat = check(`(?:[[:alpha:]]+[[:alnum:]|_]*\.pop\(.*\))`, command)

	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`pop must have an argument`)
	}

	return ``, status.No, nil
}

func validateSystemStackCall(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:push\(.+\))`, command)
	if status.Yes == stat && `` == tail {
		return tail, stat, nil
	}
	tail, stat = check(`(?:pop\([[:alpha:]]+[[:alnum:]|_]*\))`, command)
	if status.Yes == stat && `` == tail {
		return tail, stat, nil
	}

	tail, stat = check(`(?:push\(\))`, command)

	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`push must have an argument`)
	}

	tail, stat = check(`(?:pop\(\))`, command)

	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`pop must have an argument`)
	}

	tail, stat = check(`(?:pop\(.+\))`, command)

	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`invalid pop argument`)
	}

	return command, status.No, nil
}
func validateInput(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:input\([[:alpha:]]+[[:alnum:]|_]*\))`, command)

	if status.Yes == stat && `` == tail {
		return tail, stat, nil
	}

	tail, stat = check(`(?:(input)\(.+\,.*\))`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`input must have 1 argument`)
	}

	tail, stat = check(`(?:input\(\))`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`input must have 1 argument`)
	}

	return ``, status.No, nil
}

func validateGoto(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:goto\([[:alpha:]|_|#]+[[:alnum:]|\]|\[|:|\)|\(|\,|_]*\))`, command)

	if status.Yes == stat && `` == tail {
		return tail, stat, nil
	}

	tail, stat = check(`(?:goto\(.+\))`, command)

	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`unresolved symbols in goto`)
	}

	tail, stat = check(`(?:goto\(\))`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`goto must have 1 argument`)
	}

	return command, status.No, nil
}

func validateMark(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:#[[:alpha:]|_]+[[:alnum:]]*:)`, command)
	if status.Yes == stat {
		return tail, stat, nil
	}
	tail, stat = check(`(?:#.+:)`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`invalid mark name`)
	}

	return command, stat, nil
}

func validateSlice(command string) (tail string, stat int, err error) {
	re, err := regexp.Compile(`]\[`)
	if nil != err {
		panic(err)
	}
	if nil != re.FindIndex([]byte(command)) {
		return ``, status.Err, errors.New("unresolved command")
	}
	return command, status.No, nil
}

func validatePrint(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:print\()`, command)
	if status.Yes == stat {
		return tail[:len(tail)-1], stat, nil
	}
	return command, status.No, nil
}

func validateCD(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:\[(print|goto)\(.+?\)\,.+\,(print|goto)\(.+?\)\])`, command)

	if status.Yes == stat && `` == tail {
		re, err := regexp.Compile(`(?:(print|goto)\(.+?\))`)
		if nil != err {
			panic(err)
		}
		loc := re.FindAllIndex([]byte(command), -1)
		if 2 != len(loc) {
			return ``, status.Err, errors.New(`unresolved command`)
		}
		err = validateCommand(command[loc[0][0]:loc[0][1]])

		if nil != err {
			return ``, status.Err, err
		}

		err = validateCommand(command[loc[1][0]:loc[1][1]])
		if nil != err {
			return ``, status.Err, err
		}

		cond := string(re.ReplaceAll([]byte(command), []byte(``)))
		cond = cond[2 : len(cond)-2]

		err = validateCommand(cond)
		if nil != err {
			// здесь возникает проблема в внешними скобками
			// поэтому делаем повторный запуск для укороченного cond
			if "unresolved command" == err.Error() {
				err2 := validateCommand(cond[1 : len(cond)-1])
				if nil != err2 {
					return ``, status.Err, err
				}
			}
		}

		return ``, status.Yes, nil
	}

	return ``, status.No, nil
}

func validateIf(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:^if\([^{]+\){)`, command)

	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{ifCond, LineCounter, CommandToExecute})
		re, err := regexp.Compile(`(?:^if\([^{]+\){)`)
		if nil != err {
			panic(err)
		}
		loc := re.FindIndex([]byte(command))
		ifStruct := command[:loc[1]]

		err = validateCommand(ifStruct[2 : len(ifStruct)-1])
		if nil != err {
			return ``, status.Err, err
		}
		return tail, status.Yes, nil
	}

	tail, stat = check(`(?:^if\([^{]+\))`, command)

	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`missing '{' in if clause`)
	}

	return command, status.No, nil
}

func validateSimpleClause(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:\((False|True)\))`, command)
	if status.Yes == stat && `` == tail {
		return ``, stat, nil
	}
	return command, status.No, nil
}

func validateElseIf(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:^}elseif\([^{]+\){)`, command)

	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{elseIfCond, LineCounter, CommandToExecute})
		re, err := regexp.Compile(`(?:^}elseif\([^{]+\){)`)
		if nil != err {
			panic(err)
		}
		loc := re.FindIndex([]byte(command))
		elseIfStruct := command[:loc[1]]

		err = validateCommand(elseIfStruct[7 : len(elseIfStruct)-1])
		if nil != err {
			return ``, status.Err, err
		}
		return tail, status.Yes, nil
	}

	tail, stat = check(`(?:^}elseif\([^{]+\))`, command)

	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`missing '{' in else if clause`)
	}

	tail, stat = check(`(?:^elseif\([^{]+\)){`, command)

	if status.Yes == stat {
		return ``, status.Err, errors.New(`missing '}' in else if clause`)
	}

	return command, status.No, nil
}

func validateElse(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:^}else{)`, command)

	if status.Yes == stat {
		closureHistory = append(closureHistory, brace{elseCond, LineCounter, CommandToExecute})
		return tail, stat, nil
	}
	return command, status.No, nil
}

func validateImport(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:#importval)`, command)
	return tail, stat, nil
}

func validateFor(command string) (tail string, stat int, err error) {
	tail, stat = check(`^for\(`, command)
	if status.Yes == stat {
		return tail, stat, nil
	}
	re, err := regexp.Compile(`\){`)

	loc := re.FindIndex([]byte(command))

	if nil != loc {
		tempHistory := closureHistory
		err = validateCommand(command[0:loc[0]])
		closureHistory = tempHistory
		if nil != err {
			return command, status.No, nil
		}
		tail = command[loc[1]:]

		closureHistory = append(closureHistory, brace{loop, LineCounter, CommandToExecute})

		return tail, status.Yes, nil
	}

	return command, status.No, nil

}

func validateWhile(command string) (tail string, stat int, err error) {
	tail, stat = check(`^while\(`, command)
	return tail, stat, nil
}

func validateDoWhile(command string) (tail string, stat int, err error) {
	tail, stat = check(`^do{`, command)
	if stat == status.Yes {
		closureHistory = append(closureHistory, brace{loop, LineCounter, CommandToExecute})
		return tail, stat, nil
	}
	tail, stat = check(`^}while\(`, command)
	if stat == status.Yes {
		return tail[:len(tail)-1], stat, nil
	}
	return command, status.No, nil
}

func replaceVals(command string) (tail string, stat int) {
	var wasReplace bool

	tail = command
	if -1 != strings.Index(tail, "(val)") {
		wasReplace = true
		tail = strings.ReplaceAll(tail, "(val)", "val")
	}
	if wasReplace {
		return tail, status.Yes
	}
	return tail, status.No
}

func ValidateStr(command string, variables [][][]interface{}) (string, [][][]interface{}, error) {
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
			exprEnd, err := getExprEnd(tail, pos[1]-1)
			if nil != err {
				return tail, variables, err
			}

			replacerArgs = append(replacerArgs, tail[pos[0]:exprEnd])
			replacerArgs = append(replacerArgs, `val`)
		}

		r := strings.NewReplacer(replacerArgs...)
		tail = r.Replace(tail)
	}
	return tail, variables, nil
}

func validateCommand(command string) error {
	oldCommand := command

	tail, stat, err := validateFuncDefinition(command)
	if nil != err {
		return err
	}
	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateMark(command)
	if nil != err {
		return err
	}
	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateReturn(command)
	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateFigureBrace(command)
	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateVarDef(command)
	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}
	tail, stat, err = validateVar(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateInput(command)
	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	command, stat, err = validateString(command)

	if nil != err {
		return err
	}

	command, _, err = ValidateStr(command, nil)

	if nil != err {
		return err
	}
	_, _, err = validateSlice(command)
	if nil != err {
		return err
	}

	tail, stat, err = validateDoWhile(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateWhile(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateImport(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateFor(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	if !isValidBracesNum(command) {
		return errors.New("number of '(' does not equal number of ')'")
	}

	tail, stat, err = validateCD(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateIf(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateElseIf(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateSimpleClause(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateElse(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateStandardFuncCall(command, "next_command", 1, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateNumber(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "send_command", 1, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "UNSET_DEST", 0, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "DEL_DEST", 1, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "REROUTE", 0, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "get_root_source", 1, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "get_root_dest", 1, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "UNDEFINE", 1, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "RESET_SOURCE", 0, false)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateStandardFuncCall(command, "index", 2, true)

	if nil != err {
		return err
	}

	tail, stat, err = validateStandardFuncCall(command, "reg_find", 2, true)

	if nil != err {
		return err
	}

	tail, stat, err = validatePrint(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		return validateCommand(tail)
	}

	tail, stat, err = validateGoto(command)

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	if nil != err {
		return err
	}

	tail, stat, err = validateSystemStackCall(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateUserStackCall(command)

	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			return nil
		}
	}

	tail, stat, err = validateAssignment(command)
	if nil != err {
		return err
	}

	if status.Yes == stat {
		if `` == tail {
			if 1 != strings.Count(command, "=") {
				return errors.New("more than one operation '='")
			}

			pos := strings.Index(command, `=`)
			pos++
			command = command[pos:]
			return validateCommand(command)
		}
	}

	command, stat, err = validateFuncCall(command, false)
	if nil != err {
		return err
	}

	tail, stat, err = validateArithmetic(command, false)

	if nil != err {
		return err
	}

	if status.Yes == stat && "val" == tail {
		return nil
	}

	tail, stat, err = validateComparison(tail, false)

	if nil != err {
		return err
	}

	if status.Yes == stat && "val" == tail {
		return nil
	}

	if "val" == command {
		return nil
	}

	command = tail
	if oldCommand != command {
		return validateCommand(command)
	}

	tail, stat = replaceVals(command)
	if status.Yes == stat {
		return validateCommand(tail)
	}

	if "" == command {
		return nil
	}

	return errors.New("unresolved command")
}

func hasCycledImports(imports []string, rootSource string) string {
	var theseImports []string

	f, err := os.Open(rootSource)

	if nil != err {
		panic(err)
	}
	newChunk := EachChunk(f)

	for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
		if nil != err {
			handleError(err)
		}
		inputedCode := CodeInput(chunk, false)
		for len(inputedCode) > 7 && "#import" == inputedCode[0:7] {
			name := inputedCode[8 : strings.Index(inputedCode[8:], "\"")+8]
			if StringInSlice(name, imports) {
				err = f.Close()
				if nil != err {
					panic(err)
				}
				return rootSource
			}
			theseImports = append(theseImports, name)
			inputedCode = inputedCode[strings.Index(inputedCode[8:], "\"")+9:]
		}
	}

	err = f.Close()
	if nil != err {
		panic(err)
	}

	for _, imp := range theseImports {
		if imp != rootSource {
			return hasCycledImports(append(imports, imp), imp)
		}
	}
	return ``
}

func hasRepeatedImports(rootSource string) bool {
	var imports []string

	f, err := os.Open(rootSource)

	if nil != err {
		panic(err)
	}
	newChunk := EachChunk(f)

	for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
		if nil != err {
			handleError(err)
		}
		inputedCode := CodeInput(chunk, false)
		for len(inputedCode) > 7 && "#import" == inputedCode[0:7] {
			name := inputedCode[8 : strings.Index(inputedCode[8:], "\"")+8]
			if StringInSlice(name, imports) {
				err = f.Close()
				if nil != err {
					panic(err)
				}
				return true
			}
			imports = append(imports, name)
			inputedCode = inputedCode[strings.Index(inputedCode[8:], "\"")+9:]
		}
	}

	err = f.Close()
	if nil != err {
		panic(err)
	}

	return false
}

func checkImportPos(rootSource string) error {
	var command_counter int

	LineCounter = 1

	f, err := os.Open(rootSource)

	if nil != err {
		panic(err)
	}
	newChunk := EachChunk(f)

	for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
		command_counter++
		if nil != err {
			handleError(err)
		}
		inputedCode := CodeInput(chunk, true)
		if len(inputedCode) > 7 && "#import" == inputedCode[0:7] && command_counter > 1 {
			err = f.Close()
			if nil != err {
				panic(err)
			}
			CommandToExecute = strings.TrimSpace(chunk)
			return errors.New("imports allowed only in the beginning of the file")
		}
	}

	err = f.Close()
	if nil != err {
		panic(err)
	}
	return nil
}

func StaticValidate(rootSource string) (string, error) {
	var cycledFile string
	var isRepeated bool
	var err error

	err = checkImportPos(rootSource)
	if nil != err {
		return rootSource, err
	}
	isRepeated = hasRepeatedImports(rootSource)
	cycledFile = hasCycledImports([]string{rootSource}, rootSource)
	var f *os.File

	if `` != cycledFile {
		f, err = os.Open(cycledFile)
	} else {
		f, err = os.Open(rootSource)
	}

	if nil != err {
		panic(err)
	}
	newChunk := EachChunk(f)

	LineCounter = 1

	for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
		if nil != err {
			handleError(err)
		}
		CommandToExecute = strings.TrimSpace(chunk)

		if `` != cycledFile {
			return cycledFile, errors.New("import cycle not allowed")
		}
		if isRepeated {
			return rootSource, errors.New("repeated imports not allowed")
		}

		inputedCode := CodeInput(chunk, true)

		err = validateCommand(inputedCode)
		if nil != err {
			return rootSource, err
		}

		for len(inputedCode) > 7 && "#import" == inputedCode[0:7] {
			bufLineCounter := LineCounter
			bufCommandToExecute := CommandToExecute
			var bufClosureHistory []brace
			for _, el := range closureHistory {
				bufClosureHistory = append(bufClosureHistory, el)
			}
			closureHistory = nil
			name := inputedCode[8 : strings.Index(inputedCode[8:], "\"")+8]
			name, err = StaticValidate(name)
			if nil != err {
				return name, err
			}
			LineCounter = bufLineCounter
			CommandToExecute = bufCommandToExecute
			closureHistory = nil
			for _, el := range bufClosureHistory {
				closureHistory = append(closureHistory, el)
			}
			inputedCode = inputedCode[strings.Index(inputedCode[8:], "\"")+9:]
		}
	}

	if len(closureHistory) > 0 {
		LineCounter = closureHistory[len(closureHistory)-1].line
		CommandToExecute = closureHistory[len(closureHistory)-1].command

		return rootSource, errors.New(`'}' is missing`)
	}
	LineCounter = 0

	err = f.Close()
	if nil != err {
		panic(err)
	}

	return rootSource, nil
}
