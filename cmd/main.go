package main

import (
	"bint.com/internal/const/options"
	. "bint.com/internal/internalTools"
	. "bint.com/internal/validator"
	"bint.com/pkg/serviceTools"
	"fmt"
)

func main() {
	// эти опции можно менять для системной отладки
	sysMode := options.ExecBasm
	benvMode := options.ExecBenv

	var filesListToExecute []string

	toTranslate, rootSource, rootDest, keyDest, err := ParseArgs()

	if nil != err {
		return
	}

	rootSource, rootDest, keyDest, filesListToExecute = SetConf(toTranslate, rootSource, rootDest, keyDest, sysMode, benvMode)

	if options.UserTranslate == toTranslate ||
		(options.Internal == toTranslate && (options.Internal == sysMode || options.UserTranslate == sysMode)) {
		err = StaticValidate(rootSource)
		if nil != err {
			fmt.Println("ERROR in " + rootSource + " at near line " +
				fmt.Sprintf("%v", serviceTools.LineCounter))
			fmt.Println(serviceTools.CommandToExecute)
			fmt.Println(err)
			return
		}
		//Start(options.Validate, filesListToExecute, rootSource, rootDest, keyDest, sysMode, benvMode)
	}

	Start(toTranslate, filesListToExecute, rootSource, rootDest, keyDest, sysMode, benvMode)
}
