package main

import (
	"bint.com/internal/const/options"
	. "bint.com/internal/internalTools"
	"bint.com/internal/validator"
)

func main() {
	// эти опции можно менять для системной отладки
	sysMode := options.UserTranslate
	benvMode := options.ExecBenv

	var filesListToExecute []string
	toTranslate, rootSource, rootDest, keyDest, err := ParseArgs()

	if nil != err {
		return
	}

	rootSource, rootDest, keyDest, filesListToExecute = SetConf(toTranslate, rootSource, rootDest, keyDest, sysMode, benvMode)
	validator.Validate(toTranslate, filesListToExecute, rootSource, rootDest, keyDest, sysMode, benvMode)
	Start(toTranslate, filesListToExecute, rootSource, rootDest, keyDest, sysMode, benvMode)
}
