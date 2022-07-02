package validator

import (
	"bint.com/internal/const/status"
	. "bint.com/pkg/serviceTools"
	"errors"
	"os"
	"regexp"
	"strings"
)

func ValidateFuncDefinition(command string) (tail string, stat int, err error) {
	re, err := regexp.Compile(`(?m)(?:(int|float|bool|string|stack|void)[[:alnum:]]*?\` +
		`((?:((int|float|bool|string|stack))[[:alnum:]]*?\,)+)(int|float|bool|string|stack)[[:alnum:]]*?\){`)
	if nil != err {
		panic(err)
	}
	loc := re.FindIndex([]byte(command))
	if nil != loc {
		tail = command[loc[1]:]
		return tail, status.Yes, nil
	}
	return tail, status.No, nil
}

func ValidateCommand(command string) error {
	tail, stat, err := ValidateFuncDefinition(command)
	if nil != err {
		return err
	}
	if status.Yes == stat {
		return ValidateCommand(tail)
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
		err = ValidateCommand(inputedCode)

		if nil != err {
			return err
		}
	}

	LineCounter = 0
	return nil
}
