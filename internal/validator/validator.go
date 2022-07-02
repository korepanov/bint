package validator

import (
	"bint.com/internal/const/status"
	. "bint.com/pkg/serviceTools"
	"errors"
	"os"
	"regexp"
	"strings"
)

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

func validateFuncDefinition(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]]*?\,)+)(int|float|bool|string|stack)[[:alnum:]]*?\){`, command)
	if status.Yes == stat {
		return tail, stat, nil
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]]*?\`+
		`(\){)`, command)
	if status.Yes == stat {
		return tail, stat, nil
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void).*?\`+
		`((?:((int|float|bool|string|stack)).*?\,)+)(int|float|bool|string|stack).*?\){`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`invalid characters in entity names`)
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]]*?\`+
		`((?:((int|float|bool|string|stack|void))[[:alnum:]]*?\,)+)(int|float|bool|string|stack|void)[[:alnum:]]*?\){`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`function parameter cannot be of type void`)
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void).*?`+
		`(?:((int|float|bool|string|stack))[[:alnum:]]*?\,)+)(int|float|bool|string|stack).*?{`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`missing '(' or ')' in function definition`)
	}

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]]*?\,)+)(int|float|bool|string|stack)[[:alnum:]]*?\)`, command)
	if status.Yes == stat {
		return ``, status.Err, errors.New(`missing '{' in function definition`)
	}
	return ``, status.No, nil
}

func validateReturn(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?m)(?:return)`, command)
	return tail, stat, nil
}

func validateFigureBrace(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?m)(?:})`, command)
	return tail, stat, nil
}

func validateVarDef(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack)[[:alnum:]]*)`, command)
	if status.Yes == stat && `` == tail {
		return tail, stat, nil
	}
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack).*)`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`invalid variable name`)
	}
	return ``, status.No, nil
}

func validateAssignment(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:[[:alnum:]]*=.*)`, command)
	if status.Yes == stat && `` == tail {
		return tail, stat, nil
	}
	tail, stat = check(`(?:.*=.*)`, command)
	if status.Yes == stat && `` == tail {
		return ``, status.Err, errors.New(`unresolved reference`)
	}

	return ``, status.No, nil
}

//func validateStr(command string) (tail string, stat int, err error) {

//}

func validateVar(command string) (tail string, stat int, err error) {
	tail, stat = check(`(?:[[:alnum:]]+)`, command)
	return tail, stat, nil
}

func validateArithmetic(command string, isOp bool) (tail string, stat int, err error) {
	re, err := regexp.Compile(`(?m)(?:\(([[[:alnum:]]|\[|])*[+|\-|*|/|^]([[[:alnum:]]|\[|])*\))`)
	if nil != err {
		panic(err)
	}
	loc := re.FindIndex([]byte(command))
	if nil == loc {
		if isOp {
			return command, status.Yes, nil
		}
		return ``, status.No, nil
	}
	command = string(re.ReplaceAll([]byte(command), []byte("val")))
	return validateArithmetic(command, true)
}

func validateCommand(command string) error {
	if !isValidBracesNum(command) {
		return errors.New("number of '(' doest not equal number of ')'")
	}
	tail, stat, err := validateFuncDefinition(command)
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

	tail, stat, err = validateArithmetic(command, false)

	if nil != err {
		return err
	}

	if status.Yes == stat && "val" == tail {
		return nil
	}
	return errors.New("unresolved command")
}

func Validate(rootSource string) error {
	f, err := os.Open(rootSource)

	if nil != err {
		panic(err)
	}
	newChunk := EachChunk(f)

	LineCounter = 1

	for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
		CommandToExecute = strings.TrimSpace(chunk)
		inputedCode := CodeInput(chunk, true)
		err = validateCommand(inputedCode)

		if nil != err {
			return err
		}
	}

	LineCounter = 0
	return nil
}
