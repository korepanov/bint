package validator

import (
	. "bint.com/pkg/serviceTools"
	"errors"
	"os"
	"strings"
)

func ValidateCommand(command string) error {
	return errors.New("test")
}

func Validate(rootSource string) error {
	f, err := os.Open(rootSource)

	if nil != err {
		panic(err)
	}
	newChunk := EachChunk(f)

	for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
		CommandToExecute = strings.TrimSpace(chunk)
		inputedCode := CodeInput(chunk, false)
		err = ValidateCommand(inputedCode)

		if nil != err {
			return err
		}
	}

	return nil
}
