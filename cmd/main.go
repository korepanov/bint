package main

import (
	. "bint.com/internal/internalTools"
	"bint.com/internal/options"
)

func main() {
	// эти опции можно менять для системной отладки
	sysMode := options.Encrypt
	benvMode := options.InterpBenv

	var filesListToExecute []string

	toTranslate, rootSource, rootDest, keyDest, err := ParseArgs()

	if nil != err {
		return
	}

	rootSource, rootDest, keyDest, filesListToExecute = SetConf(toTranslate, rootSource, rootDest, keyDest, sysMode, benvMode)

	Start(toTranslate, filesListToExecute, rootSource, rootDest, keyDest, sysMode, benvMode)
}
