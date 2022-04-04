package main

import (
	. "bint.com/internal/internalTools"
	"bint.com/internal/options"
)

func main() {
	var filesListToExecute []string
	toTranslate, rootSource, rootDest, err := ParseArgs()

	if nil != err {
		return
	}

	// последнюю опцию можно менять для системной отладки
	rootSource, rootDest, filesListToExecute = SetConf(toTranslate, rootSource, rootDest, options.User)

	Start(filesListToExecute, rootSource, rootDest)
}
