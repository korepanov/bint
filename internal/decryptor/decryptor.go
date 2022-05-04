package decryptor

import (
	"bint.com/internal/options"
	"strconv"
	"strings"
)

func Decrypt(inputedCode string, key string) string {
	var output string
	var lexemList []string
	var keyList []string
	var codeList []string

	keyList = strings.Split(key, ", ")
	codeList = strings.Split(inputedCode, options.BendSep)

	for _, keyVal := range keyList {
		keyVal = strings.TrimSpace(keyVal)
		intKeyVal, err := strconv.Atoi(keyVal)
		if nil != err {
			panic(err)
		}
		lexemList = append(lexemList, codeList[intKeyVal])
	}

	for i := 0; i < len(lexemList); i++ {
		if i < len(lexemList)-1 {
			output += lexemList[i] + options.BendSep
		} else {
			output += lexemList[i]
		}
	}

	return output
}
