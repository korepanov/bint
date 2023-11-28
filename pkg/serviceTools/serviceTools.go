package serviceTools

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"
)

var LineCounter int
var CommandToExecute string
var commentBegin bool

func IsValidString(s string) bool {
	if "$" == string(s[0]) {
		s = s[1:]
	}
	for _, symbol := range s {
		if !IsValidSymbol(symbol) {
			return false
		}
	}
	return true
}

func IsValidSymbol(symbol rune) bool {
	if unicode.IsLetter(symbol) {
		return true
	}
	if unicode.IsDigit(symbol) {
		return true
	}
	if "_" == string(symbol) {
		return true
	}
	return false
}

func LookBehind(reg string, s string) ([]string, error) {
	re, err := regexp.Compile(reg)
	if nil != err {
		return nil, err
	}
	locArr := re.FindAllIndex([]byte(s), -1)

	var res []string

	for _, loc := range locArr {
		if loc[0] > 0 {
			res = append(res, string(s[loc[0]-1]))
		} else {
			res = append(res, "%")
		}
	}

	return res, nil
}

func LookAhead(reg string, s string) ([]string, error) {
	re, err := regexp.Compile(reg)
	if nil != err {
		return nil, err
	}
	locArr := re.FindAllIndex([]byte(s), -1)

	var res []string

	for _, loc := range locArr {
		if loc[1] < len(s)-1 {
			res = append(res, string(s[loc[1]]))
		} else {
			res = append(res, "$")
		}
	}

	return res, nil
}

func CheckEntry(reg string, command string) (bool, error) {
	behindSymbols, err := LookBehind(reg, command)
	if nil != err {
		return false, err
	}

	re, err := regexp.Compile(reg)

	if nil != err {
		return false, err
	}

	locArr := re.FindAllIndex([]byte(command), -1)

	var funcLocArr [][]int

	for i := 0; i < len(behindSymbols); i++ {
		if !("_" == behindSymbols[i] || unicode.IsLetter(rune(behindSymbols[i][0])) || unicode.IsDigit(rune(behindSymbols[i][0]))) {
			funcLocArr = append(funcLocArr, locArr[i])
		}
	}

	if nil != funcLocArr {
		return true, nil
	}

	return false, nil
}

func GetFuncNameEntry(funcName string, command string) [][]int {
	reg := funcName

	behindSymbols, err := LookBehind(reg, command)
	if nil != err {
		panic(err)
	}

	aheadSymbols, err := LookAhead(reg, command)
	if nil != err {
		panic(err)
	}

	re, err := regexp.Compile(reg)

	if nil != err {
		panic(err)
	}

	locArr := re.FindAllIndex([]byte(command), -1)

	var funcLocArr [][]int

	for i := 0; i < len(behindSymbols); i++ {
		if !("_" == behindSymbols[i] || unicode.IsLetter(rune(behindSymbols[i][0])) || unicode.IsDigit(rune(behindSymbols[i][0])) ||
			"$" == behindSymbols[i]) {
			if !("_" == aheadSymbols[i] || unicode.IsLetter(rune(aheadSymbols[i][0])) || unicode.IsDigit(rune(aheadSymbols[i][0]))) {
				funcLocArr = append(funcLocArr, locArr[i])
			}
		}
	}

	return funcLocArr
}

func ReplaceFunc(reg string, command string) (string, error) {
	behindSymbols, err := LookBehind(reg, command)
	if nil != err {
		return ``, err
	}

	re, err := regexp.Compile(reg)

	if nil != err {
		return ``, err
	}

	locArr := re.FindAllIndex([]byte(command), -1)

	var funcLocArr [][]int

	for i := 0; i < len(behindSymbols); i++ {
		if !("_" == behindSymbols[i] || unicode.IsLetter(rune(behindSymbols[i][0]))) {
			funcLocArr = append(funcLocArr, locArr[i])
		}
	}

	var replacerArgs []string

	for _, loc := range funcLocArr {
		replacerArgs = append(replacerArgs, command[loc[0]:loc[1]])
		replacerArgs = append(replacerArgs, "val")
	}

	newCommand := command

	if nil != replacerArgs {
		r := strings.NewReplacer(replacerArgs...)
		newCommand = r.Replace(command)
	}

	return newCommand, nil
}

func FindExprListEnd(exprList [][]interface{}, exprBegin int) int {
	openedBraces := 1
	closedBraces := 0
	//exprBegin - первый токен после "("
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

func FindExprEnd(expr string, exprBegin int) int {
	openedBraces := 1
	closedBraces := 0
	//exprBegin - первый символ после "("
	i := exprBegin

	for openedBraces != closedBraces {
		if "(" == string(expr[i]) {
			openedBraces += 1
		}
		if ")" == string(expr[i]) {
			closedBraces += 1
		}
		i += 1
	}
	// возвращает символ, следующий за финалным
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
		"get_root_dest", "is_letter", "is_digit", "exists"}
	return StringInSlice(OP, operations)
}

func CanBePartOfBoolExpr(subExpr string) bool {
	partsOfBoolExpr := []string{"(", ")", "AND", "OR", "XOR", "NOT", "<", "<=", "==", ">", ">=", "True", "False",
		"+", "-", "*", "/", "^", "@", "[", "]", "exists"}
	return StringInSlice(subExpr, partsOfBoolExpr) || IsNumber(subExpr)
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

func EachChunk(file *os.File) func() (string, error) {
	const chunkSize = 200
	chunk := make([]byte, chunkSize)
	var buffer string
	var resList [2]string
	var part string

	//pattern, err := regexp.Compile("((?:[^;\"']|\"[^\"]*\"|'[^']*'|\".*)+)")
	pattern, err := regexp.Compile(`((?:[^;"']|"(?:\\"|[^"])*"|'[^']*'|".*)+)`)

	if nil != err {
		panic(err)
	}

	_, err = file.Read(chunk)

	if nil != err && io.EOF != err {
		panic(err)
	}
	buffer = string(chunk)

	return func() (string, error) {
		var wasSemicolon bool

		trimBuffer := strings.TrimSpace(buffer)

		if len(trimBuffer) > 0 && ";" == string(trimBuffer[len(trimBuffer)-1]) {
			wasSemicolon = true
		}

		if "" == strings.TrimSpace(buffer) {
			buffer = part + buffer
			chunk = make([]byte, chunkSize)
			_, err := file.Read(chunk)

			if io.EOF == err {
				return "end", nil
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
				return "end", nil
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
			if "" == resList[1] {
				return "", errors.New("too long command")
			}
		}
		part = resList[0]
		buffer = resList[1]

		if wasSemicolon {
			if "" != strings.TrimSpace(buffer) {
				buffer += ";"
			}

			part += ";"

			return part[:len(part)-1], nil
		}

		return part, nil
	}
}

func SetCommandCounter(file *os.File, COMMAND_COUNTER int) (func() (string, error), error) {
	_, err := file.Seek(0, 0)

	if nil != err {
		return nil, err
	}

	newChunk := EachChunk(file)

	if 0 == COMMAND_COUNTER {
		return newChunk, nil
	}

	i := 1

	for _, err = newChunk(); i < COMMAND_COUNTER-1; _, err = newChunk() {
		if nil != err {
			return nil, err
		}
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

	for chunk, err := newChunk(); "end" != chunk; chunk, err = newChunk() {
		if nil != err {
			return i, f, err
		}
		expr := CodeInput(chunk, true)
		if len(expr) > len(mark) && expr[0:len(mark)] == mark && ":" == string(expr[len(mark)]) {
			return i, f, nil
		}
		i++
	}
	LineCounter = shadowLineCounter
	err = errors.New("serviceTools: getCommandCounterByMark: ERROR: no such mark: " + mark)
	return i, f, err
}

func GetMark(expr string) string {
	var mark string

	if "#" == string(expr[0]) {
		mark = "#"
		i := 1
		for ":" != string(expr[i]) {
			mark += string(expr[i])
			i++
		}
	}
	return mark
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

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func IsKeyWordWithAssignment(word string) bool {
	words := []string{"len", "exists", "int", "bool", "float", "index", "is_letter", "is_digit", "reg_find"}
	if StringInSlice(word, words) {
		return true
	}
	return false
}

func IsKeyWord(word string) bool {
	words := []string{"AND", "OR", "NOT", "XOR", "print", "len", "exists", "index", "is_letter", "is_digit", "reg_find",
		"pop", "push", "input", "next_command", "get_root_source", "get_root_dest", "send_command", "goto", "exit",
		"SET_SOURCE", "SET_DEST", "SEND_DEST", "DEL_DEST", "UNDEFINE", "UNSET_SOURCE", "REROUTE", "UNSET_DEST",
		"RESET_SOURCE", "str", "int", "float", "bool", "string", "stack", "do", "while", "for", "void",
		"if", "index", "try", "error"}
	if StringInSlice(word, words) {
		return true
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
		"<=", "==", ">", ">=", "=", "+", "-", "*", "/", "^", "@", "int", "bool", "float",
		"goto", "SET_SOURCE", "UNSET_SOURCE", "RESET_SOURCE",
		"SET_DEST", "UNSET_DEST", "next_command", "send_command", "UNDEFINE", "pop", "push",
		"DEL_DEST", "SEND_DEST", "REROUTE", ".", "len", "index", "get_root_source", "get_root_dest",
		"is_letter", "is_digit", "reg_find", "exists"}

	if StringInSlice(s, ops) {
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
	if "True" == s || "true" == s {
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

	if len(val) > 0 && `[` == string(val[0]) && `]` == string(val[len(val)-1]) {
		return "stack"
	}

	return "string"
}

func replace(re *regexp.Regexp, src string, repl string) string {
	loc := re.FindIndex([]byte(src))

	res := src

	if nil != loc {
		res = src[:loc[0]] + repl + src[loc[1]:]
	}

	return res
}

func cutComment(expr string) string {
	re, err := regexp.Compile(`"(\\.|[^"])*"`)
	if nil != err {
		panic(err)
	}

	var strs []string

	loc := re.FindIndex([]byte(expr))

	var i int
	res := expr

	for nil != loc {
		strs = append(strs, res[loc[0]:loc[1]])
		res = replace(re, res, "$"+fmt.Sprintf("%v", i))
		loc = re.FindIndex([]byte(res))
		i++
	}

	re, err = regexp.Compile(`(?s).*\*/`)

	if nil != err {
		panic(err)
	}

	if commentBegin {
		if nil == re.FindIndex([]byte(res)) {
			return ""
		} else {
			res = string(re.ReplaceAll([]byte(res), []byte("")))
			commentBegin = false
		}
	}

	re, err = regexp.Compile(`//.*`)
	if nil != err {
		panic(err)
	}

	res = string(re.ReplaceAll([]byte(res), []byte("")))

	re, err = regexp.Compile(`/\*(?s).*\*/`)
	if nil != err {
		panic(err)
	}
	res = string(re.ReplaceAll([]byte(res), []byte("")))

	re, err = regexp.Compile(`/\*(?s).*`)
	if nil != err {
		panic(err)
	}
	if nil != re.FindIndex([]byte(res)) {
		commentBegin = true
	}

	res = string(re.ReplaceAll([]byte(res), []byte("")))

	if len(strings.TrimSpace(res)) > 0 && "{" == string(strings.TrimSpace(res)[len(strings.TrimSpace(res))-1]) {
		res = res + `print("")`
	}

	for i = 0; i < len(strs); i++ {
		res = strings.Replace(res, "$"+fmt.Sprintf("%v", i), strs[i], 1)
	}

	return res
}

func CodeInput(expr string, lineIncrement bool) string {
	//fmt.Println(expr)
	var stringsInside []string
	var poses []int
	var pos int
	var startFlag bool
	var stringInside string

	var i int

	if lineIncrement {
		LineCounter += strings.Count(expr, "\n")
	}

	expr = cutComment(expr)

	if "" == strings.TrimSpace(expr) {
		return `print("")`
	}

	//запоминаем строки, чтобы оставить в них пробелы
	for i = 0; i < len([]rune(expr)); i++ {
		if startFlag {
			if `"` != string([]rune(expr)[i]) || (i > 0 && string([]rune(expr)[i-1]) == `\` && `"` == string([]rune(expr)[i])) {
				stringInside += string([]rune(expr)[i])

			} else {
				startFlag = false
				stringsInside = append(stringsInside, stringInside)
				stringInside = ""
				continue
			}
		}
		if `"` == string([]rune(expr)[i]) {
			startFlag = true
		}
	}

	expr = strings.Replace(expr, " ", "", -1)
	expr = strings.Replace(expr, "\t", "", -1)
	expr = strings.Replace(expr, "\n", "", -1)

	// запоминаем местоположение строк
	for i = 0; i < len([]rune(expr)); i++ {
		pos += 1
		if startFlag {
			if i > 0 && `\` != string([]rune(expr)[i-1]) && `"` == string([]rune(expr)[i]) {
				startFlag = false
				continue
			}
		}

		if i > 0 && `\` != string([]rune(expr)[i-1]) && `"` == string([]rune(expr)[i]) {
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

		lenStringInside = len([]rune(stringInside))

		for j := poses[i]; j < poses[i]+lenStringInside; j++ {
			expr = string([]rune(expr)[:poses[i]]) + string([]rune(expr)[poses[i]+1:])
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
		leftExpr = string([]rune(expr)[:poses[i]])
		rightExpr = string([]rune(expr)[poses[i]:])
		expr = leftExpr + str + rightExpr
		// пересчитываем позиции из-за изменившегося выражения
		for k := i + 1; k < len(poses); k++ {
			poses[k] += len([]rune(str))
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

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyDirectory(scrDir, dest string) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if _, err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		if "bint" != name && "benv" != name {
			err = os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
