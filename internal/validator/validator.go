package validator

import (
	. "bint.com/pkg/serviceTools"
	"os"
	"strings"
)

func ValidateCommand(command string) error {

	return nil
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
