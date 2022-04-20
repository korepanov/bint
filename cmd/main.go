package main

import (
	. "bint.com/internal/internalTools"
	"bint.com/internal/options"
)

func main() {
	// эти опции можно менять для системной отладки
	sysMode := options.No
	benvMode := options.ExecBenv

	var filesListToExecute []string

	toTranslate, rootSource, rootDest, err := ParseArgs()

	if nil != err {
		return
	}

	rootSource, rootDest, filesListToExecute = SetConf(toTranslate, rootSource, rootDest, sysMode, benvMode)

	Start(toTranslate, filesListToExecute, rootSource, rootDest, sysMode, benvMode)
}
