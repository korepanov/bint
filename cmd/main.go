package main

import (
	. "bint.com/internal/internalTools"
	"bint.com/internal/options"
)

func main() {
	// эту опцию можно менять для системной отладки
	sysMode := options.Transpile

	var filesListToExecute []string

	toTranslate, rootSource, rootDest, err := ParseArgs()

	if nil != err {
		return
	}

	// последнюю опцию можно менять для системной отладки
	rootSource, rootDest, filesListToExecute = SetConf(toTranslate, rootSource, rootDest, sysMode)

	Start(filesListToExecute, rootSource, rootDest, sysMode)
}
