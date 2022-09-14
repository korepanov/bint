package validator

import (
	"bint.com/internal/const/status"
	. "bint.com/pkg/serviceTools"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var COMMAND_COUNTER int
var fileName string

func filter(command string) bool {
	if "$" == string(command[0]) && "$" == string(command[len(command)-1]) {
		return false
	}
	return true
}

func handleError(errMessage string) {
	var inputedCode string

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

	fmt.Println(inputedCode)
	err = f.Close()
	if nil != err {
		panic(err)
	}
}
func dValidateFunc(command string) (tail string, stat int) {
	var retVal string

	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){`, command)
	if status.Yes == stat {
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

		for _, loc := range locs {
			fmt.Println(tail[loc[0]:loc[1]])
		}

		tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
			`((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){`, command)

		return tail, stat
	}
	tail, stat = check(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\`+
		`(\){)`, command)
	if status.Yes == stat {
		reg, err := regexp.Compile(`^(?:(int|float|bool|string|stack|void))`)
		if nil != err {
			panic(err)
		}
		loc := reg.FindIndex([]byte(command))
		retVal = command[loc[0]:loc[1]]
		fmt.Println(retVal)

		return tail, stat
	}

	return tail, stat
}

func dynamicValidateCommand(command string) error {
	dValidateFunc(command)
	return nil
}

func DynamicValidate(validatingFile string, rootSource string) (err error) {
	COMMAND_COUNTER = 1
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
			err = dynamicValidateCommand(inputedCode)
			if nil != err {
				return err
			}
		}
	}

	err = f.Close()
	if nil != err {
		panic(err)
	}

	return nil
}
