package decryptor

import (
	"bint.com/internal/options"
	"strconv"
	"strings"
)

func Decrypt(inputedCode string, key string) string {
	var output string
	var lexemeList []string
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
		lexemeList = append(lexemeList, codeList[intKeyVal])
	}

	for i := 0; i < len(lexemeList); i++ {
		if i < len(lexemeList)-1 {
			output += lexemeList[i] + options.BendSep
		} else {
			output += lexemeList[i]
		}
	}

	return output
}
