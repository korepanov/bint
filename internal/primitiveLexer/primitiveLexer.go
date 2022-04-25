package primitiveLexer

import "strings"

func PrimitiveLexicalAnalyze(expr string, variables [][]interface{}) ([][]interface{}, [][]interface{}, error) {

	var res [][]interface{}

	strRes := strings.Split(expr, ",")
	if "string" == strRes[0] || "float" == strRes[0] || "int" == strRes[0] || "bool" == strRes[0] || "stack" == strRes[0] {
		variables = append(variables, []interface{}{strRes[0], strRes[1], []interface{}{"var_val"}})
	} else {
		var temp []interface{}

		for _, el := range strRes {
			temp = append(temp, el)
		}
		res = append(res, temp)
	}

	if 0 == len(res) {
		res = append(res, []interface{}{"res", 0})
	}

	return res, variables, nil

}
