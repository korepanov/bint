package serviceTools

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var LineCounter = 0
var CommandToExecute string

func FindExprListEnd(exprList [][]interface{}, exprBegin int) int {
	openedBraces := 1
	closedBraces := 0
	//exprBegin - первы токен после "("
	i := exprBegin

	for openedBraces != closedBraces {
		if "(" == exprList[i][1].(string) {
			openedBraces += 1
		}
		if ")" == exprList[i][1].(string) {
			closedBraces += 1
		}
		i += 1
	}
	// возвращает токен, следующий за финалным
	return i
}

func Pop(list [][]interface{}, i int) [][]interface{} {
	copy(list[i:], list[i+1:])
	list = list[:len(list)-1]
	return list
}

func Insert(a [][]interface{}, index int, value []interface{}) [][]interface{} {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func IsUnaryOperation(OP string) bool {
	operations := []string{"print", "str", "input", "int", "float", "bool", "goto", "SET_SOURCE", "SET_DEST",
		"next_command", "send_command", "UNDEFINE", "pop", "push", "DEL_DEST", "SEND_DEST", "len", "get_root_source",
		"get_root_dest", "is_letter", "is_digit"}
	return stringInSlice(OP, operations)
}

func CanBePartOfBoolExpr(subExpr string) bool {
	partsOfBoolExpr := []string{"(", ")", "AND", "OR", "XOR", "NOT", "<", "<=", "==", ">", ">=", "True", "False",
		"+", "-", "*", "/", "^", "[", "]"}
	return stringInSlice(subExpr, partsOfBoolExpr) || IsNumber(subExpr)
}

func mySplit(buffer string, pattern *regexp.Regexp) [2]string {
	findList := pattern.FindAllString(buffer, -1)
	var resList [2]string
	resList[0] = findList[0]
	for i := 1; i < len(findList); i++ {
		resList[1] += findList[i]
		if i != len(findList)-1 {
			resList[1] += ";"
		}
	}

	trimBuffer := strings.TrimSpace(buffer)
	trimRes := strings.TrimSpace(resList[1])

	if len(trimBuffer) > 0 && ";" == string(trimBuffer[len(trimBuffer)-1]) &&
		len(trimRes) > 0 && ";" != string(trimRes[len(trimRes)-1]) {
		resList[1] += ";"
	}

	return resList
}

func EachChunk(file *os.File) func() string {
	const chunkSize = 100
	chunk := make([]byte, chunkSize)
	var buffer string
	var resList [2]string
	var part string

	pattern, err := regexp.Compile("((?:[^;\"']|\"[^\"]*\"|'[^']*'|\".*)+)")

	if nil != err {
		panic(err)
	}

	_, err = file.Read(chunk)

	if nil != err && io.EOF != err {
		panic(err)
	}
	buffer = string(chunk)

	return func() string {
		var wasSemicolon bool

		trimBuffer := strings.TrimSpace(buffer)

		if len(trimBuffer) > 0 && ";" == string(trimBuffer[len(trimBuffer)-1]) {
			wasSemicolon = true
		}

		if "" == strings.TrimSpace(buffer) {
			buffer = part
			chunk = make([]byte, chunkSize)
			_, err := file.Read(chunk)

			if io.EOF == err {
				return "end"
			}
			if nil != err && io.EOF != err {
				panic(err)
			}

			if !wasSemicolon {
				buffer += string(chunk)
			}
			resList = mySplit(buffer, pattern)
			part = resList[0]
			buffer = resList[1]
		}

		if -1 != strings.Index(buffer, ";") {
			resList = mySplit(buffer, pattern)
		} else {
			resList[0] = buffer
			resList[1] = ""
			part = resList[0]
			buffer = resList[1]

			buffer = part
			chunk = make([]byte, chunkSize)
			_, err := file.Read(chunk)

			if io.EOF == err {
				return "end"
			}
			if nil != err && io.EOF != err {
				panic(err)
			}

			if !wasSemicolon {
				buffer += string(chunk)
			}
			if len(buffer) > 0 && ";" == string(buffer[len(buffer)-1]) {
				wasSemicolon = true
			}
			resList = mySplit(buffer, pattern)

		}
		part = resList[0]
		buffer = resList[1]

		if wasSemicolon {
			if "" != strings.TrimSpace(buffer) {
				buffer += ";"
			}

			part += ";"

			return part[:len(part)-1]
		}

		return part
	}
}

func SetCommandCounter(file *os.File, COMMAND_COUNTER int) (func() string, error) {
	_, err := file.Seek(0, 0)
	newChunk := EachChunk(file)

	if nil != err {
		return newChunk, err
	}

	i := 1

	for _ = newChunk(); i < COMMAND_COUNTER-1; _ = newChunk() {
		i++
	}

	return newChunk, nil
}

func GetCommandCounterByMark(f *os.File, mark string) (int, *os.File, error) {
	i := 1
	shadowLineCounter := LineCounter
	LineCounter = 0

	_, err := f.Seek(0, 0)
	if nil != err {
		return i, f, err
	}
	newChunk := EachChunk(f)
	for chunk := newChunk(); "end" != chunk; chunk = newChunk() {
		expr := CodeInput(chunk, true)
		if len(expr) > len(mark) && expr[0:len(mark)] == mark && ":" == string(expr[len(mark)]) {
			return i, f, nil
		}
		i++
	}
	LineCounter = shadowLineCounter
	err = errors.New("serviceTools: get_command_counter_by_mark: ERROR: no such mark: " + mark)
	return i, f, err
}

// Exists проверка существования файла
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if nil == err {
		return true
	}
	return false
}

func IsOp(s string) bool {
	ops := []string{"AND", "OR", "XOR", "NOT", "print", "input", "L: True", "L: False", "str", "<",
		"<=", "==", ">", ">=", "=", "+", "-", "*", "/", "^", "int", "bool", "float",
		"goto", "SET_SOURCE", "UNSET_SOURCE", "RESET_SOURCE",
		"SET_DEST", "UNSET_DEST", "next_command", "send_command", "UNDEFINE", "pop", "push",
		"DEL_DEST", "SEND_DEST", "REROUTE", ".", "len", "index", "get_root_source", "get_root_dest",
		"is_letter", "is_digit"}

	if stringInSlice(s, ops) {
		return true
	}
	return false
}

func EachVariable(variables [][]interface{}) func() []interface{} {
	i := len(variables)
	endInterf := []interface{}{"end"}
	return func() []interface{} {
		i--
		if i >= 0 {
			return variables[i]
		}
		return endInterf
	}
}

func StrToBool(s string) bool {
	if "True" == s {
		return true
	}
	return false
}

func BoolToStr(val bool) string {
	if val {
		return "True"
	}
	return "False"
}

func WhatsType(val string) string {
	if len(val) > 0 && `"` == string(val[0]) && `"` == string(val[len(val)-1]) {
		return "string"
	}

	if "True" == val || "False" == val || "true" == val || "false" == val {
		return "bool"
	}

	_, err := strconv.Atoi(val)
	if nil == err {
		return "int"
	}

	_, err = strconv.ParseFloat(val, 64)
	if nil == err {
		return "float"
	}

	return "string"
}

func CodeInput(expr string, lineIncrement bool) string {
	var stringsInside []string
	var poses []int
	var pos int
	var startFlag bool
	var stringInside string

	// базовые комментарии в одну строку
	lineComPos := strings.Index(expr, "//")

	var i int

	if -1 != lineComPos {
		i = lineComPos
		for i < len(expr) && "\n" != string(expr[i]) {
			expr = expr[:i] + expr[i+1:]
		}
	}

	//запоминаем стоки, чтобы оставить в них пробелы
	for _, ch := range expr {
		if startFlag {
			if `"` != string(ch) {
				stringInside += string(ch)
			} else {
				startFlag = false
				stringsInside = append(stringsInside, stringInside)
				stringInside = ""
				continue
			}
		}
		if `"` == string(ch) {
			startFlag = true
		}
	}
	if lineIncrement {
		if 0 == strings.Count(expr, "\n") {
			LineCounter++
		} else {
			LineCounter += strings.Count(expr, "\n")
		}
	}
	expr = strings.Replace(expr, " ", "", -1)
	expr = strings.Replace(expr, "\t", "", -1)
	expr = strings.Replace(expr, "\n", "", -1)

	// запоминаем местоположение строк
	for _, ch := range expr {
		pos += 1
		if startFlag {
			if `"` == string(ch) {
				startFlag = false
				continue
			}
		}

		if `"` == string(ch) {
			poses = append(poses, pos)
			startFlag = true
		}
	}

	// вырезаем строки из выражений и перерассчитываем их местоположение
	i = 0
	var lenStringInside int

	for _, str := range stringsInside {
		stringInside = strings.Replace(str, " ", "", -1)
		stringInside = strings.Replace(stringInside, "\t", "", -1)
		stringInside = strings.Replace(stringInside, "\n", "", -1)

		lenStringInside = len(stringInside)

		for j := poses[i]; j < poses[i]+lenStringInside; j++ {
			expr = expr[:poses[i]] + expr[poses[i]+1:]
			// переситываем позиции из-за изменившегося выражения
			for k := i + 1; k < len(poses); k++ {
				poses[k] -= 1
			}
		}
		i++
	}

	i = 0
	var leftExpr string
	var rightExpr string
	for _, str := range stringsInside {
		leftExpr = expr[:poses[i]]
		rightExpr = expr[poses[i]:]
		expr = leftExpr + str + rightExpr
		// пересчитываем позиции из-за изменившегося выражения
		for k := i + 1; k < len(poses); k++ {
			poses[k] += len(str)
		}
		i += 1
	}

	return expr
}

func UnfoldInterfaceSlice(exprList []interface{}) [][]interface{} {
	var res [][]interface{}

	for _, el := range exprList {
		res = append(res, el.([]interface{}))
	}

	return res
}

func ValueFoldInterface(exprList interface{}) interface{} {
	if "string" == fmt.Sprintf("%T", exprList) {
		return exprList
	}
	for "[]interface {}" == fmt.Sprintf("%T", exprList.([]interface{})[0]) {
		if 1 == len(exprList.([]interface{})) {
			exprList = exprList.([]interface{})[0]
		} else {
			return exprList
		}
	}

	return fmt.Sprintf("%v", exprList.([]interface{})[0])
}
