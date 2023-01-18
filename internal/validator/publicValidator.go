package validator

import (
	"bint.com/internal/const/options"
	. "bint.com/internal/internalTools"
	. "bint.com/pkg/serviceTools"
	"fmt"
	"os"
)

func Validate(toTranslate int, filesListToExecute []string, rootSource string,
	rootDest string, keyDest string, sysMode int, benvMode bool) {

	defer func() {
		if r := recover(); nil != r {
			fmt.Println("ERROR in " + FileToExecute + " at near line " + fmt.Sprintf("%v", LineCounter))
			fmt.Println(CommandToExecute)
			fmt.Println(r)
			os.Exit(1)
		}
	}()

	var err error
	var name string

	if options.UserTranslate == toTranslate ||
		(options.Internal == toTranslate && (options.Internal == sysMode || options.UserTranslate == sysMode)) {
		// статическая валидация
		name, err = StaticValidate(rootSource)

		if nil != err {
			FileToExecute = name
			panic(err)
		}
		// динамическая валидация

		if options.UserTranslate == toTranslate || (options.Internal == toTranslate && options.UserTranslate == sysMode) {
			Start(options.UserValidate, filesListToExecute, rootSource, rootDest, keyDest, sysMode, benvMode)
			var number int
			file := "benv/trace/trace_program" + fmt.Sprintf("%v", number) + ".b"
			for Exists(file) {
				DynamicValidate(file, rootSource)
				number++
				file = "benv/trace/trace_program" + fmt.Sprintf("%v", number) + ".b"
			}
			validatingFile := "benv/trace_program.b"
			DynamicValidate(validatingFile, rootSource)
		} /*else {
			Start(options.InternalValidate, filesListToExecute, rootSource, rootDest, keyDest, sysMode, benvMode)
			var number int
			file := "benv/internal/trace/trace_program" + fmt.Sprintf("%v", number) + ".b"
			for Exists(file) {
				DynamicValidate(file, rootSource)
				number++
				file = "benv/internal/trace/trace_program" + fmt.Sprintf("%v", number) + ".b"
			}
			validatingFile := "benv/internal/trace_program.b"
			DynamicValidate(validatingFile, rootSource)
		}*/

	}
}
