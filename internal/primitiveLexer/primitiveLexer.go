package primitiveLexer

import (
	"bint.com/internal/const/options"
	"errors"
	"strings"
)

func PrimitiveLexicalAnalyze(expr string, variables [][]interface{}) ([][]interface{}, [][]interface{}, error) {

	var res [][]interface{}

	strRes := strings.Split(expr, options.BendSep)
	if "string" == strRes[0] || "float" == strRes[0] || "int" == strRes[0] || "bool" == strRes[0] {
		variables = append(variables, []interface{}{strRes[0], strRes[1], []interface{}{"var_val"}})
	} else if "stack" == strRes[0] {
		variables = append(variables, []interface{}{strRes[0], strRes[1], []interface{}{[]interface{}{"end"}}})
	} else {
		var temp []interface{}

		if "#" == string(strRes[0][0]) {
			pos := strings.Index(strRes[0], ":")
			if -1 == pos {
				panic(errors.New("ERROR: mark must end with \":\""))
			}
			pos++
			strRes[0] = strRes[0][pos:]
		}
		for _, el := range strRes {
			if "\"" == string(el[0]) && "\"" == string(el[len(el)-1]) {
				el = strings.Replace(el[1:len(el)-1], "\\\"", "\"", -1)
				el = "\"" + el + "\""
			}
			temp = append(temp, el)
		}
		res = append(res, temp)
	}

	if 0 == len(res) {
		res = append(res, []interface{}{"res", 0})
	}

	return res, variables, nil

}
