package main 

import(
"strings"
"fmt"
"io"
"os"
"strconv"
"unicode"
"math"
"errors"
"path/filepath"
"regexp"
"flag"
)
var iFlag = "-i"
var oFlag = "-o"
var systemStack = []interface{}{"end"}
var vars = make(map[string][]interface{})
var commentBegin bool

func ValueFoldInterface(exprList interface{}) interface{} {
	if "[]interface {}" != fmt.Sprintf("%T", exprList) {
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

func Exists(name string) bool {
	if _, err := os.Stat(getRootSource(name)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func defineVar(varName string){
vars[varName] = append(vars[varName], interface{}(nil))
}

func undefineVar(varName string){
if len(vars[varName]) > 0{
vars[varName] = vars[varName][:len(vars[varName]) - 1]
}
}

func setVar(varName string, varVal interface{}){
vars[varName][len(vars[varName]) - 1] = varVal
}

func getVar(varName string) interface{}{
return vars[varName][len(vars[varName]) - 1]
}

func toBool(s interface{}) bool{
if "bool" == fmt.Sprintf("%T", s){
	return s.(bool)
}
if "string" == fmt.Sprintf("%T", s){
if "True" == s.(string) || "true" == s.(string){
return true
}else{
return false
}
}

panic(errors.New("toBool: failed to parse bool, arg: " + fmt.Sprintf("%v", s)))
} 

func toFloat(n interface{}) float64{
res, err := strconv.ParseFloat(fmt.Sprintf("%v", n), 64)
if nil != err{
panic(errors.New("toFloat: failed to parse float, arg: " + fmt.Sprintf("%v", n)))
}
return res 
}

func toInt(n interface{}) int{
res, err := strconv.Atoi(fmt.Sprintf("%v", n))
if nil != err{
panic(errors.New("toInt: failed to parse int, arg: " + fmt.Sprintf("%v", n)))
}
return res 
}

func isEqual(a interface{}, b interface{}) bool{
if len(fmt.Sprintf("%v", a)) > 0 && len(fmt.Sprintf("%v", b)) > 0 && "\"" == string(fmt.Sprintf("%v", a)[0]) && "\"" == string(fmt.Sprintf("%v", b)[0]) {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

resA, errA := strconv.Atoi(fmt.Sprintf("%v", a))
resB, errB := strconv.Atoi(fmt.Sprintf("%v", b))

if errA != nil || errB != nil{

resA, errA := strconv.ParseFloat(fmt.Sprintf("%v", a), 64)
resB, errB := strconv.ParseFloat(fmt.Sprintf("%v", b), 64)

if nil != errA || nil != errB{
return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

return resA == resB
}

return resA == resB 
}

func sum(a interface{}, b interface{}) interface{}{
if len(fmt.Sprintf("%v", a)) > 0 && len(fmt.Sprintf("%v", b)) > 0 && "\"" == string(fmt.Sprintf("%v", a)[0]) && "\"" == string(fmt.Sprintf("%v", b)[0]) {
	return fmt.Sprintf("%v", a) + fmt.Sprintf("%v", b)
}

resA, errA := strconv.Atoi(fmt.Sprintf("%v", a))
resB, errB := strconv.Atoi(fmt.Sprintf("%v", b))

if errA != nil || errB != nil{

resA, errA := strconv.ParseFloat(fmt.Sprintf("%v", a), 64)
resB, errB := strconv.ParseFloat(fmt.Sprintf("%v", b), 64)

if nil != errA || nil != errB{
return fmt.Sprintf("%v", a)+ fmt.Sprintf("%v", b)
}

return resA + resB
}

return resA + resB 
}

func openFile(fileName interface{}) *os.File{
f, err := os.Open(fmt.Sprintf("%v", fileName))
if nil != err{
panic(errors.New("could not open file " + fmt.Sprintf("%v", fileName) + ": " + err.Error()))
}
return f 
}

func openFile666(fileName interface{}) *os.File{
f, err := os.OpenFile(fmt.Sprintf("%v", fileName), os.O_WRONLY, 0666)
if nil != err{
panic(errors.New("could not open file " + fmt.Sprintf("%v", fileName) + ": " + err.Error()))
}
return f 
}

func createFile(fileName interface{}) *os.File{
f, err := os.Create(fmt.Sprintf("%v", fileName))
if nil != err{
panic(errors.New("could not create file " + fmt.Sprintf("%v", fileName) + ": " + err.Error()))
}
return f 
}

func getRootSource(sourceI interface{}) string{
source := fmt.Sprintf("%v", sourceI)
ex, err := os.Executable()
if err != nil {
panic(errors.New("getRootSource: " + err.Error()))
}
exPath := filepath.Dir(ex)
i := strings.Index(exPath, "benv")
return exPath[:i] + source 
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
	const chunkSize = 256
	chunk := make([]byte, chunkSize)
	var buffer string
	var resList [2]string
	var part string

	pattern, err := regexp.Compile(`((?:[^;"']|"(?:\\"|[^"])*"|'[^']*'|".*)+)`)

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
	
	if len(strings.TrimSpace(res)) > 0 && "{" == string(strings.TrimSpace(res)[len(strings.TrimSpace(res)) - 1]){
		res = res + `print("")`
	}

	for i = 0; i < len(strs); i++ {
		res = strings.Replace(res, "$"+fmt.Sprintf("%v", i), strs[i], 1)
	}

	return res
}

func CodeInput(expr string, lineIncrement bool) string {
	var stringsInside []string
	var poses []int
	var pos int
	var startFlag bool
	var stringInside string

	var i int

	expr = cutComment(expr)
	
	if "" == strings.TrimSpace(expr) {
		return `print("")`
	}
	//запоминаем стоки, чтобы оставить в них пробелы
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

func compileRegexp(r string) *regexp.Regexp{
	pattern, err := regexp.Compile(r)
	if nil != err{
		panic(errors.New("could compile regexp " + r + ": " + err.Error()))
	}
	
	return pattern 
}

func main(){
flag.StringVar(&iFlag, "i", "", "-i input.b")
flag.StringVar(&oFlag, "o", "", "-o output.b")
flag.Parse()

_ = strings.Index("", "")
fmt.Printf("")
f, _ := os.Create("temp.b")
f2, _ := os.Create("temp2.b")
io.Copy(f, f2)
os.Remove(f.Name())
os.Remove(f2.Name());
_, _ = strconv.ParseFloat("64", 64)
_ = unicode.IsLetter('a')
_ = math.Pow(2, 2);
	
	
defineVar("$ret")
setVar("$ret", "")
defineVar("COMMAND_COUNTER")
setVar("COMMAND_COUNTER", 0)
defineVar("command")
setVar("command", "")
setVar("COMMAND_COUNTER","0")
defineVar("$reverse_return_var")
setVar("$reverse_return_var", "")
defineVar("$reverse_res")
setVar("$reverse_res", []interface{}{"end"})
goto reverse_end
goto reverse
reverse:
fmt.Print("")
defineVar("s")
setVar("s", []interface{}{"end"})
fmt.Print("")
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("buf")
setVar("buf", "")
defineVar("res")
setVar("res", []interface{}{"end"})
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
goto _reverse_s
_reverse_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto _reverse_e
}else{
print("")
}
setVar("res", append(getVar("res").([]interface{}), getVar("buf")))
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
goto _reverse_s
goto _reverse_e
_reverse_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
undefineVar("res")
undefineVar("buf")
undefineVar("s")
if "#reverse" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto main_res0
}
goto reverse_end
reverse_end:
fmt.Print("")
defineVar("$indexes_return_var")
setVar("$indexes_return_var", "")
defineVar("$indexes_res")
setVar("$indexes_res", []interface{}{"end"})
goto indexes_end
goto indexes
indexes:
fmt.Print("")
defineVar("sub_s")
setVar("sub_s", "")
defineVar("s")
setVar("s", "")
fmt.Print("")
setVar("sub_s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("res")
setVar("res", []interface{}{"end"})
defineVar("i")
setVar("i", 0)
defineVar("pointer")
setVar("pointer", 0)
defineVar("s_len")
setVar("s_len", 0)
defineVar("s_len_old")
setVar("s_len_old", 0)
defineVar("sub_len")
setVar("sub_len", 0)
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("s"))))
setVar("s_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("sub_s"))))
setVar("sub_len", getVar("$l0"))
undefineVar("$l0")
setVar("s_len_old", getVar("s_len"))
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("s")), fmt.Sprintf("%v",getVar("sub_s"))))
setVar("i", getVar("$ind0"))
undefineVar("$ind0")
setVar("pointer", getVar("i"))
goto _indexes_s
_indexes_s:
fmt.Print("")
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("i"))){
goto _indexes_e
}else{
print("")
}
setVar("i",sum(getVar("i"), toFloat(getVar("s_len_old"))-toFloat(getVar("s_len"))))
setVar("res", append(getVar("res").([]interface{}), getVar("i")))
setVar("pointer",sum(getVar("pointer"), getVar("sub_len")))
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("pointer"))
setVar("$sl_right0", getVar("s_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("s").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("s", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("s"))))
setVar("s_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("s")), fmt.Sprintf("%v",getVar("sub_s"))))
setVar("i", getVar("$ind0"))
undefineVar("$ind0")
setVar("pointer", getVar("i"))
goto _indexes_s
goto _indexes_e
_indexes_e:
fmt.Print("")
defineVar("$reverse_res0")
setVar("$reverse_res0", []interface{}{"end"})
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_return_var","#reverse_res0")
goto reverse
goto reverse_res0
reverse_res0:
setVar("$reverse_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$reverse_res0", getVar("$reverse_res"))
setVar("res", getVar("$reverse_res0"))
undefineVar("$reverse_res0")
systemStack = append(systemStack, getVar("res"))
undefineVar("sub_len")
undefineVar("s_len_old")
undefineVar("s_len")
undefineVar("pointer")
undefineVar("i")
undefineVar("res")
undefineVar("sub_s")
undefineVar("s")
if "#reverse" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto main_res0
}
goto indexes_end
indexes_end:
fmt.Print("")
defineVar("$SET_COMMAND_COUNTER_return_var")
setVar("$SET_COMMAND_COUNTER_return_var", "")
goto SET_COMMAND_COUNTER_end
goto SET_COMMAND_COUNTER
SET_COMMAND_COUNTER:
fmt.Print("")
defineVar("counter")
setVar("counter", 0)
fmt.Print("")
setVar("counter", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("i")
setVar("i", 0)
setVar("i","0")
defineVar("command")
setVar("command", "")
getVar("$SOURCE").(*os.File).Seek(0, 0)
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
setVar("COMMAND_COUNTER", getVar("counter"))
goto _set_start
_set_start:
fmt.Print("")
if toFloat(getVar("i"))<toFloat(getVar("counter")){
print("")
}else{
goto _set_end
}
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("i",sum(getVar("i"), 1))
goto _set_start
goto _set_end
_set_end:
fmt.Print("")
fmt.Print("")
undefineVar("command")
undefineVar("i")
undefineVar("counter")
if "#reverse" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto main_res0
}
goto SET_COMMAND_COUNTER_end
SET_COMMAND_COUNTER_end:
fmt.Print("")
defineVar("$get_command_return_var")
setVar("$get_command_return_var", "")
defineVar("$get_command_res")
setVar("$get_command_res", "")
goto get_command_end
goto get_command
get_command:
fmt.Print("")
defineVar("counter")
setVar("counter", 0)
fmt.Print("")
setVar("counter", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("i")
setVar("i", 0)
defineVar("buf")
setVar("buf", "")
defineVar("command")
setVar("command", "")
setVar("i","0")
getVar("$SOURCE").(*os.File).Seek(0, 0)
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
goto _get_command_s
_get_command_s:
fmt.Print("")
if toFloat(getVar("i"))<toFloat(getVar("counter")){
print("")
}else{
goto _get_command_e
}
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("i",sum(getVar("i"), 1))
goto _get_command_s
goto _get_command_e
_get_command_e:
fmt.Print("")
systemStack = append(systemStack, getVar("COMMAND_COUNTER"))
setVar("$SET_COMMAND_COUNTER_return_var","#SET_COMMAND_COUNTER_res0")
goto SET_COMMAND_COUNTER
goto SET_COMMAND_COUNTER_res0
SET_COMMAND_COUNTER_res0:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
undefineVar("command")
undefineVar("buf")
undefineVar("i")
undefineVar("counter")
if "#reverse" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto main_res0
}
goto get_command_end
get_command_end:
fmt.Print("")
defineVar("$switch_command_return_var")
setVar("$switch_command_return_var", "")
goto switch_command_end
goto switch_command
switch_command:
fmt.Print("")
fmt.Print("")
setVar("COMMAND_COUNTER",sum(getVar("COMMAND_COUNTER"), 1))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
if "#reverse" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto main_res0
}
goto switch_command_end
switch_command_end:
fmt.Print("")
defineVar("$stack_len_return_var")
setVar("$stack_len_return_var", "")
defineVar("$stack_len_res")
setVar("$stack_len_res", 0)
goto stack_len_end
goto stack_len
stack_len:
fmt.Print("")
defineVar("s")
setVar("s", []interface{}{"end"})
fmt.Print("")
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("res")
setVar("res", 0)
defineVar("buf")
setVar("buf", "")
setVar("res","0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
goto _stack_len_s
_stack_len_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto _stack_len_e
}else{
print("")
}
setVar("res",sum(getVar("res"), 1))
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
goto _stack_len_s
goto _stack_len_e
_stack_len_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
undefineVar("buf")
undefineVar("res")
undefineVar("s")
if "#reverse" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto main_res0
}
goto stack_len_end
stack_len_end:
fmt.Print("")
defineVar("$ops_return_var")
setVar("$ops_return_var", "")
defineVar("$ops_res")
setVar("$ops_res", []interface{}{"end"})
goto ops_end
goto ops
ops:
fmt.Print("")
defineVar("op")
setVar("op", "")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("op", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("quotes")
setVar("quotes", []interface{}{"end"})
defineVar("these_quotes")
setVar("these_quotes", []interface{}{"end"})
defineVar("op_nums")
setVar("op_nums", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("res")
setVar("res", []interface{}{"end"})
defineVar("num1")
setVar("num1", 0)
defineVar("num2")
setVar("num2", 0)
defineVar("op_num")
setVar("op_num", 0)
defineVar("was_quote")
setVar("was_quote", false)
defineVar("to_add")
setVar("to_add", false)
setVar("was_quote","False")
setVar("to_add","True")
defineVar("$indexes_res0")
setVar("$indexes_res0", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("op"))
setVar("$indexes_return_var","#indexes_res0")
goto indexes
goto indexes_res0
indexes_res0:
setVar("$indexes_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$indexes_res0", getVar("$indexes_res"))
setVar("op_nums", getVar("$indexes_res0"))
undefineVar("$indexes_res0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("op_nums")) && len(getVar("op_nums").([]interface{})) > 1{
setVar("buf", getVar("op_nums").([]interface{})[len(getVar("op_nums").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("op_nums")) && !isEqual("end", getVar("op_nums").([]interface{})[len(getVar("op_nums").([]interface{})) - 1]) && !isEqual("[end]", getVar("op_nums").([]interface{})[len(getVar("op_nums").([]interface{})) - 1]){
setVar("op_nums", getVar("op_nums").([]interface{})[:len(getVar("op_nums").([]interface{})) - 1])
}
goto _op_nums_s
_op_nums_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto _op_nums_e
}else{
print("")
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("op_num", getVar("$I0"))
undefineVar("$I0")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`"(\\.|[^"])*"`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("quotes", getVar("$stack_var0"))
undefineVar("$stack_var0")
goto _quotes_s
_quotes_s:
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("quotes")) && len(getVar("quotes").([]interface{})) > 1{
setVar("these_quotes", getVar("quotes").([]interface{})[len(getVar("quotes").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("these_quotes")) == "[]interface {}"{
setVar("these_quotes", []interface{}{[]interface{}{"end"}})
}else{
setVar("these_quotes", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("quotes")) && !isEqual("end", getVar("quotes").([]interface{})[len(getVar("quotes").([]interface{})) - 1]) && !isEqual("[end]", getVar("quotes").([]interface{})[len(getVar("quotes").([]interface{})) - 1]){
setVar("quotes", getVar("quotes").([]interface{})[:len(getVar("quotes").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("these_quotes")) && len(getVar("these_quotes").([]interface{})) > 1{
setVar("buf", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("these_quotes")) && !isEqual("end", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1]) && !isEqual("[end]", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1]){
setVar("these_quotes", getVar("these_quotes").([]interface{})[:len(getVar("these_quotes").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto _quotes_e
}else{
print("")
}
goto _these_quotes_s
_these_quotes_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto _these_quotes_e
}else{
print("")
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("num1", getVar("$I0"))
undefineVar("$I0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("these_quotes")) && len(getVar("these_quotes").([]interface{})) > 1{
setVar("buf", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("these_quotes")) && !isEqual("end", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1]) && !isEqual("[end]", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1]){
setVar("these_quotes", getVar("these_quotes").([]interface{})[:len(getVar("these_quotes").([]interface{})) - 1])
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("num2", getVar("$I0"))
undefineVar("$I0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("these_quotes")) && len(getVar("these_quotes").([]interface{})) > 1{
setVar("buf", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("these_quotes")) && !isEqual("end", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1]) && !isEqual("[end]", getVar("these_quotes").([]interface{})[len(getVar("these_quotes").([]interface{})) - 1]){
setVar("these_quotes", getVar("these_quotes").([]interface{})[:len(getVar("these_quotes").([]interface{})) - 1])
}
setVar("was_quote","True")
goto _these_quotes_s
goto _these_quotes_e
_these_quotes_e:
fmt.Print("")
if toBool(toFloat(getVar("op_num"))>toFloat(getVar("num1")))&&toBool(toFloat(getVar("op_num"))<toFloat(getVar("num2"))){
print("")
}else{
goto _is_op_end
}
setVar("to_add","False")
goto _push_op_end
goto _is_op_end
_is_op_end:
fmt.Print("")
goto _quotes_s
goto _quotes_e
_quotes_e:
fmt.Print("")
if toBool(getVar("was_quote"))&&toBool(!toBool(getVar("to_add"))){
goto _push_op_end
}else{
print("")
}
setVar("res", append(getVar("res").([]interface{}), getVar("op_num")))
setVar("was_quote","False")
setVar("to_add","True")
goto _push_op_end
_push_op_end:
fmt.Print("")
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("op_nums")) && len(getVar("op_nums").([]interface{})) > 1{
setVar("buf", getVar("op_nums").([]interface{})[len(getVar("op_nums").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("op_nums")) && !isEqual("end", getVar("op_nums").([]interface{})[len(getVar("op_nums").([]interface{})) - 1]) && !isEqual("[end]", getVar("op_nums").([]interface{})[len(getVar("op_nums").([]interface{})) - 1]){
setVar("op_nums", getVar("op_nums").([]interface{})[:len(getVar("op_nums").([]interface{})) - 1])
}
setVar("was_quote","False")
setVar("to_add","True")
goto _op_nums_s
goto _op_nums_e
_op_nums_e:
fmt.Print("")
defineVar("$reverse_res1")
setVar("$reverse_res1", []interface{}{"end"})
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_return_var","#reverse_res1")
goto reverse
goto reverse_res1
reverse_res1:
setVar("$reverse_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$reverse_res1", getVar("$reverse_res"))
setVar("res", getVar("$reverse_res1"))
undefineVar("$reverse_res1")
systemStack = append(systemStack, getVar("res"))
undefineVar("to_add")
undefineVar("was_quote")
undefineVar("op_num")
undefineVar("num2")
undefineVar("num1")
undefineVar("res")
undefineVar("buf")
undefineVar("op_nums")
undefineVar("these_quotes")
undefineVar("quotes")
undefineVar("op")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto main_res0
}
goto ops_end
ops_end:
fmt.Print("")
defineVar("$block_end_return_var")
setVar("$block_end_return_var", "")
defineVar("$block_end_res")
setVar("$block_end_res", 0)
goto block_end_end
goto block_end
block_end:
fmt.Print("")
fmt.Print("")
defineVar("op1")
setVar("op1", "")
defineVar("op2")
setVar("op2", "")
defineVar("code")
setVar("code", "")
defineVar("command_buf")
setVar("command_buf", "")
defineVar("o_sum")
setVar("o_sum", 0)
defineVar("c_sum")
setVar("c_sum", 0)
defineVar("command_len")
setVar("command_len", 0)
defineVar("obraces")
setVar("obraces", []interface{}{"end"})
defineVar("cbraces")
setVar("cbraces", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("spos")
setVar("spos", "")
defineVar("counter")
setVar("counter", 0)
defineVar("buf_counter")
setVar("buf_counter", 0)
defineVar("pos")
setVar("pos", 0)
setVar("counter", getVar("COMMAND_COUNTER"))
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0","1")
setVar("$sl_right0", getVar("command_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("command", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("code", getVar("command"))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("counter",sum(getVar("counter"), 1))
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_internal0","0")
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",string(getVar("command").(string)[toInt(getVar("$sl_internal0"))]))
setVar("code",sum(getVar("code"), getVar("$sl0")))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("o_sum","1")
setVar("c_sum","0")
goto _block_s
_block_s:
fmt.Print("")
if isEqual(ValueFoldInterface(getVar("o_sum")), ValueFoldInterface(getVar("c_sum"))){
goto _block_e
}else{
print("")
}
defineVar("$ops_res0")
setVar("$ops_res0", []interface{}{"end"})
systemStack = append(systemStack, getVar("code"))
systemStack = append(systemStack, "{")
setVar("$ops_return_var","#ops_res0")
goto ops
goto ops_res0
ops_res0:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res0", getVar("$ops_res"))
setVar("obraces", getVar("$ops_res0"))
undefineVar("$ops_res0")
defineVar("$ops_res1")
setVar("$ops_res1", []interface{}{"end"})
systemStack = append(systemStack, getVar("code"))
systemStack = append(systemStack, "}")
setVar("$ops_return_var","#ops_res1")
goto ops
goto ops_res1
ops_res1:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res1", getVar("$ops_res"))
setVar("cbraces", getVar("$ops_res1"))
undefineVar("$ops_res1")
defineVar("$stack_len_res0")
setVar("$stack_len_res0", 0)
systemStack = append(systemStack, getVar("obraces"))
setVar("$stack_len_return_var","#stack_len_res0")
goto stack_len
goto stack_len_res0
stack_len_res0:
setVar("$stack_len_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$stack_len_res0", getVar("$stack_len_res"))
setVar("o_sum", getVar("$stack_len_res0"))
undefineVar("$stack_len_res0")
defineVar("$stack_len_res1")
setVar("$stack_len_res1", 0)
systemStack = append(systemStack, getVar("cbraces"))
setVar("$stack_len_return_var","#stack_len_res1")
goto stack_len
goto stack_len_res1
stack_len_res1:
setVar("$stack_len_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$stack_len_res1", getVar("$stack_len_res"))
setVar("c_sum", getVar("$stack_len_res1"))
undefineVar("$stack_len_res1")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0","1")
setVar("$sl_right0", getVar("command_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("command_buf", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("code",sum(getVar("code"), getVar("command_buf")))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("counter",sum(getVar("counter"), 1))
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_internal0","0")
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",string(getVar("command").(string)[toInt(getVar("$sl_internal0"))]))
setVar("command_buf", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("code",sum(getVar("code"), getVar("command_buf")))
goto _block_s
goto _block_e
_block_e:
fmt.Print("")
setVar("buf_counter",toFloat(getVar("COMMAND_COUNTER"))-1)
systemStack = append(systemStack, getVar("buf_counter"))
setVar("$SET_COMMAND_COUNTER_return_var","#SET_COMMAND_COUNTER_res1")
goto SET_COMMAND_COUNTER
goto SET_COMMAND_COUNTER_res1
SET_COMMAND_COUNTER_res1:
fmt.Print("")
fmt.Print("")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("COMMAND_COUNTER",sum(getVar("COMMAND_COUNTER"), 1))
systemStack = append(systemStack, toFloat(getVar("counter"))-1)
undefineVar("pos")
undefineVar("buf_counter")
undefineVar("counter")
undefineVar("spos")
undefineVar("buf")
undefineVar("cbraces")
undefineVar("obraces")
undefineVar("command_len")
undefineVar("c_sum")
undefineVar("o_sum")
undefineVar("command_buf")
undefineVar("code")
undefineVar("op2")
undefineVar("op1")
if "#reverse" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto main_res0
}
goto block_end_end
block_end_end:
fmt.Print("")
defineVar("$println_return_var")
setVar("$println_return_var", "")
goto println_end
goto println
println:
fmt.Print("")
defineVar("s")
setVar("s", "")
fmt.Print("")
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("$print_arg0")
setVar("$print_arg0", "")
setVar("$print_arg0", getVar("s"))
fmt.Print(getVar("$print_arg0"))
undefineVar("$print_arg0")
defineVar("$print_arg0")
setVar("$print_arg0", "")
setVar("$print_arg0","\n")
fmt.Print(getVar("$print_arg0"))
undefineVar("$print_arg0")
undefineVar("s")
if "#reverse" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto main_res0
}
goto println_end
println_end:
fmt.Print("")
defineVar("$in_stack_return_var")
setVar("$in_stack_return_var", "")
defineVar("$in_stack_res")
setVar("$in_stack_res", false)
goto in_stack_end
goto in_stack
in_stack:
fmt.Print("")
defineVar("el")
setVar("el", "")
defineVar("s")
setVar("s", []interface{}{"end"})
fmt.Print("")
setVar("el", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("buf")
setVar("buf", "")
defineVar("res")
setVar("res", false)
setVar("res","False")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
goto _in_stack_s
_in_stack_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto _in_stack_e
}else{
print("")
}
if isEqual(ValueFoldInterface(getVar("el")), ValueFoldInterface(getVar("buf"))){
print("")
}else{
goto _no
}
setVar("res","True")
goto _in_stack_e
goto _no
_no:
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
goto _in_stack_s
goto _in_stack_e
_in_stack_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
undefineVar("res")
undefineVar("buf")
undefineVar("el")
undefineVar("s")
if "#reverse" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto main_res0
}
goto in_stack_end
in_stack_end:
fmt.Print("")
defineVar("$func_end_return_var")
setVar("$func_end_return_var", "")
defineVar("$func_end_res")
setVar("$func_end_res", 0)
goto func_end_end
goto func_end
func_end:
fmt.Print("")
defineVar("func_begin")
setVar("func_begin", 0)
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("func_begin", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("obraces")
setVar("obraces", []interface{}{"end"})
defineVar("cbraces")
setVar("cbraces", []interface{}{"end"})
defineVar("obrace")
setVar("obrace", "")
defineVar("cbrace")
setVar("cbrace", "")
defineVar("symbol")
setVar("symbol", "")
defineVar("o_sum")
setVar("o_sum", 0)
defineVar("c_sum")
setVar("c_sum", 0)
defineVar("pos")
setVar("pos", 0)
defineVar("spos")
setVar("spos", "")
defineVar("command_len")
setVar("command_len", 0)
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
setVar("obrace","(")
setVar("cbrace",")")
setVar("o_sum","1")
setVar("c_sum","0")
setVar("pos",sum(getVar("func_begin"), 1))
defineVar("$ops_res2")
setVar("$ops_res2", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("obrace"))
setVar("$ops_return_var","#ops_res2")
goto ops
goto ops_res2
ops_res2:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res2", getVar("$ops_res"))
setVar("obraces", getVar("$ops_res2"))
undefineVar("$ops_res2")
defineVar("$ops_res3")
setVar("$ops_res3", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("cbrace"))
setVar("$ops_return_var","#ops_res3")
goto ops
goto ops_res3
ops_res3:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res3", getVar("$ops_res"))
setVar("cbraces", getVar("$ops_res3"))
undefineVar("$ops_res3")
goto _braces_s
_braces_s:
fmt.Print("")
if toFloat(getVar("pos"))<toFloat(getVar("command_len")){
print("")
}else{
goto _braces_e
}
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("pos"))
setVar("spos", getVar("$s0"))
undefineVar("$s0")
defineVar("$in_stack_res0")
setVar("$in_stack_res0", false)
systemStack = append(systemStack, getVar("obraces"))
systemStack = append(systemStack, getVar("spos"))
setVar("$in_stack_return_var","#in_stack_res0")
goto in_stack
goto in_stack_res0
in_stack_res0:
setVar("$in_stack_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$in_stack_res0", getVar("$in_stack_res"))
if toBool(getVar("$in_stack_res0")){
print("")
}else{
goto _o_plus_end
}
undefineVar("$in_stack_res0")
setVar("o_sum",sum(getVar("o_sum"), 1))
goto _o_plus_end
_o_plus_end:
fmt.Print("")
defineVar("$in_stack_res1")
setVar("$in_stack_res1", false)
systemStack = append(systemStack, getVar("cbraces"))
systemStack = append(systemStack, getVar("spos"))
setVar("$in_stack_return_var","#in_stack_res1")
goto in_stack
goto in_stack_res1
in_stack_res1:
setVar("$in_stack_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$in_stack_res1", getVar("$in_stack_res"))
if toBool(getVar("$in_stack_res1")){
print("")
}else{
goto _c_plus_end
}
undefineVar("$in_stack_res1")
setVar("c_sum",sum(getVar("c_sum"), 1))
goto _c_plus_end
_c_plus_end:
fmt.Print("")
if isEqual(ValueFoldInterface(getVar("o_sum")), ValueFoldInterface(getVar("c_sum"))){
goto _braces_e
}else{
print("")
}
setVar("pos",sum(getVar("pos"), 1))
goto _braces_s
goto _braces_e
_braces_e:
fmt.Print("")
systemStack = append(systemStack, getVar("pos"))
undefineVar("command_len")
undefineVar("spos")
undefineVar("pos")
undefineVar("c_sum")
undefineVar("o_sum")
undefineVar("symbol")
undefineVar("cbrace")
undefineVar("obrace")
undefineVar("cbraces")
undefineVar("obraces")
undefineVar("func_begin")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto main_res0
}
goto func_end_end
func_end_end:
fmt.Print("")
defineVar("root_source")
setVar("root_source", "")
defineVar("command")
setVar("command", "")
defineVar("first_file")
setVar("first_file", false)
defineVar("num")
setVar("num", 0)
defineVar("exit_num")
setVar("exit_num", 0)
defineVar("br_closed")
setVar("br_closed", 0)
defineVar("br_opened")
setVar("br_opened", 0)
defineVar("$init_return_var")
setVar("$init_return_var", "")
goto init_end
goto init
init:
fmt.Print("")
fmt.Print("")
setVar("br_closed","0")
setVar("br_opened","0")
setVar("root_source","benv/for_program.b")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", getVar("root_source")))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/if_program.b"))))
if "#reverse" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto main_res0
}
goto init_end
init_end:
fmt.Print("")
defineVar("$finish_return_var")
setVar("$finish_return_var", "")
goto finish_end
goto finish
finish:
fmt.Print("")
fmt.Print("")
getVar("$SOURCE").(*os.File).Close()
getVar("$DEST").(*os.File).Close()
if "#reverse" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto main_res0
}
goto finish_end
finish_end:
fmt.Print("")
defineVar("$is_var_def_return_var")
setVar("$is_var_def_return_var", "")
defineVar("$is_var_def_res")
setVar("$is_var_def_res", false)
goto is_var_def_end
goto is_var_def
is_var_def:
fmt.Print("")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("s")
setVar("s", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`^(?:(int|float|bool|stack|string)[^\(]*)`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("s", getVar("$stack_var0"))
undefineVar("$stack_var0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
systemStack = append(systemStack, !toBool(isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf")))))
undefineVar("buf")
undefineVar("s")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto main_res0
}
goto is_var_def_end
is_var_def_end:
fmt.Print("")
defineVar("$Type_return_var")
setVar("$Type_return_var", "")
defineVar("$Type_res")
setVar("$Type_res", "")
goto Type_end
goto Type
Type:
fmt.Print("")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("s")
setVar("s", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`(?:(^int))`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("s", getVar("$stack_var0"))
undefineVar("$stack_var0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto int_end
}else{
print("")
}
systemStack = append(systemStack, "int")
undefineVar("buf")
undefineVar("s")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_res0
}
goto int_end
int_end:
fmt.Print("")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`(?:(^float))`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("s", getVar("$stack_var0"))
undefineVar("$stack_var0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto float_end
}else{
print("")
}
systemStack = append(systemStack, "float")
undefineVar("s")
undefineVar("buf")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_res0
}
goto float_end
float_end:
fmt.Print("")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`(?:(^bool))`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("s", getVar("$stack_var0"))
undefineVar("$stack_var0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto bool_end
}else{
print("")
}
systemStack = append(systemStack, "bool")
undefineVar("buf")
undefineVar("s")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_res0
}
goto bool_end
bool_end:
fmt.Print("")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`(?:(^stack))`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("s", getVar("$stack_var0"))
undefineVar("$stack_var0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto stack_end
}else{
print("")
}
systemStack = append(systemStack, "stack")
undefineVar("s")
undefineVar("buf")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_res0
}
goto stack_end
stack_end:
fmt.Print("")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`(?:(^string))`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("s", getVar("$stack_var0"))
undefineVar("$stack_var0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto string_end
}else{
print("")
}
systemStack = append(systemStack, "string")
undefineVar("buf")
undefineVar("s")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_res0
}
goto string_end
string_end:
fmt.Print("")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(`(?:(^void))`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("command"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("$stack_var0",getVar("$regRes"))
setVar("s", getVar("$stack_var0"))
undefineVar("$stack_var0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto void_end
}else{
print("")
}
systemStack = append(systemStack, "void")
undefineVar("s")
undefineVar("buf")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto main_res0
}
goto void_end
void_end:
fmt.Print("")
defineVar("$print_arg0")
setVar("$print_arg0", "")
setVar("$print_arg0", getVar("command"))
fmt.Print(getVar("$print_arg0"))
undefineVar("$print_arg0")
defineVar("$print_arg0")
setVar("$print_arg0", "")
setVar("$print_arg0","\n")
fmt.Print(getVar("$print_arg0"))
undefineVar("$print_arg0")
defineVar("$print_arg0")
setVar("$print_arg0", "")
setVar("$print_arg0","Type: ERROR\n")
fmt.Print(getVar("$print_arg0"))
undefineVar("$print_arg0")
os.Exit(1)
goto Type_end
Type_end:
fmt.Print("")
defineVar("$check_br_return_var")
setVar("$check_br_return_var", "")
goto check_br_end
goto check_br
check_br:
fmt.Print("")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("symbol")
setVar("symbol", "")
defineVar("s")
setVar("s", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
setVar("symbol","{")
defineVar("$ops_res4")
setVar("$ops_res4", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("symbol"))
setVar("$ops_return_var","#ops_res4")
goto ops
goto ops_res4
ops_res4:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res4", getVar("$ops_res"))
setVar("s", getVar("$ops_res4"))
undefineVar("$ops_res4")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto br_opened_end
}else{
print("")
}
setVar("br_opened",sum(getVar("br_opened"), 1))
goto br_opened_end
br_opened_end:
fmt.Print("")
setVar("symbol","}")
defineVar("$ops_res5")
setVar("$ops_res5", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("symbol"))
setVar("$ops_return_var","#ops_res5")
goto ops
goto ops_res5
ops_res5:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res5", getVar("$ops_res"))
setVar("s", getVar("$ops_res5"))
undefineVar("$ops_res5")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto br_closed_end
}else{
print("")
}
setVar("br_closed",sum(getVar("br_closed"), 1))
goto br_closed_end
br_closed_end:
fmt.Print("")
fmt.Print("")
undefineVar("buf")
undefineVar("s")
undefineVar("symbol")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$check_br_return_var")){
goto main_res0
}
goto check_br_end
check_br_end:
fmt.Print("")
defineVar("$reset_br_return_var")
setVar("$reset_br_return_var", "")
goto reset_br_end
goto reset_br
reset_br:
fmt.Print("")
fmt.Print("")
setVar("br_opened","0")
setVar("br_closed","0")
if "#reverse" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$reset_br_return_var")){
goto main_res0
}
goto reset_br_end
reset_br_end:
fmt.Print("")
defineVar("$is_if_return_var")
setVar("$is_if_return_var", "")
defineVar("$is_if_res")
setVar("$is_if_res", false)
goto is_if_end
goto is_if
is_if:
fmt.Print("")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("op")
setVar("op", "")
defineVar("s")
setVar("s", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("ibuf")
setVar("ibuf", 0)
defineVar("start_pos")
setVar("start_pos", 0)
setVar("op","if(")
defineVar("$ops_res6")
setVar("$ops_res6", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("op"))
setVar("$ops_return_var","#ops_res6")
goto ops
goto ops_res6
ops_res6:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res6", getVar("$ops_res"))
setVar("s", getVar("$ops_res6"))
undefineVar("$ops_res6")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto end_true
}else{
print("")
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("ibuf", getVar("$I0"))
undefineVar("$I0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("ibuf"))){
goto start_true
}else{
print("")
}
setVar("start_pos",toFloat(getVar("ibuf"))-4)
if toFloat(getVar("start_pos"))<0{
goto end_true
}else{
print("")
}
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("start_pos"))
setVar("$sl_right0", getVar("ibuf"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("buf", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
if isEqual(ValueFoldInterface("else"), ValueFoldInterface(getVar("buf"))){
print("")
}else{
goto end_true
}
goto start_true
start_true:
fmt.Print("")
systemStack = append(systemStack, "True")
undefineVar("start_pos")
undefineVar("ibuf")
undefineVar("buf")
undefineVar("s")
undefineVar("op")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_res0
}
goto end_true
end_true:
fmt.Print("")
systemStack = append(systemStack, "False")
undefineVar("op")
undefineVar("s")
undefineVar("buf")
undefineVar("ibuf")
undefineVar("start_pos")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$is_if_return_var")){
goto main_res0
}
fmt.Print("")
goto is_if_end
is_if_end:
fmt.Print("")
defineVar("$get_cond_return_var")
setVar("$get_cond_return_var", "")
defineVar("$get_cond_res")
setVar("$get_cond_res", "")
goto get_cond_end
goto get_cond
get_cond:
fmt.Print("")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("op")
setVar("op", "")
defineVar("s")
setVar("s", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("start_pos")
setVar("start_pos", 0)
defineVar("end_pos")
setVar("end_pos", 0)
setVar("op","if(")
defineVar("$ops_res7")
setVar("$ops_res7", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("op"))
setVar("$ops_return_var","#ops_res7")
goto ops
goto ops_res7
ops_res7:
setVar("$ops_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_res7", getVar("$ops_res"))
setVar("s", getVar("$ops_res7"))
undefineVar("$ops_res7")
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && len(getVar("s").([]interface{})) > 1{
setVar("buf", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("s")) && !isEqual("end", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]) && !isEqual("[end]", getVar("s").([]interface{})[len(getVar("s").([]interface{})) - 1]){
setVar("s", getVar("s").([]interface{})[:len(getVar("s").([]interface{})) - 1])
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("start_pos", getVar("$I0"))
undefineVar("$I0")
setVar("start_pos",sum(getVar("start_pos"), 2))
defineVar("$func_end_res0")
setVar("$func_end_res0", 0)
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("start_pos"))
setVar("$func_end_return_var","#func_end_res0")
goto func_end
goto func_end_res0
func_end_res0:
setVar("$func_end_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$func_end_res0", getVar("$func_end_res"))
setVar("end_pos", getVar("$func_end_res0"))
undefineVar("$func_end_res0")
setVar("end_pos",sum(getVar("end_pos"), 1))
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("start_pos"))
setVar("$sl_right0", getVar("end_pos"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("buf", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
systemStack = append(systemStack, getVar("buf"))
undefineVar("end_pos")
undefineVar("start_pos")
undefineVar("buf")
undefineVar("s")
undefineVar("op")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto main_res0
}
goto get_cond_end
get_cond_end:
fmt.Print("")
defineVar("$if_type_return_var")
setVar("$if_type_return_var", "")
defineVar("$if_type_res")
setVar("$if_type_res", "")
goto if_type_end
goto if_type
if_type:
fmt.Print("")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("command_len")
setVar("command_len", 0)
defineVar("prefix")
setVar("prefix", "")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
if isEqual(ValueFoldInterface(1), ValueFoldInterface(getVar("command_len"))){
print("")
}else{
goto not_clear
}
systemStack = append(systemStack, "clear")
undefineVar("prefix")
undefineVar("command_len")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_res0
}
goto not_clear
not_clear:
fmt.Print("")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0","1")
setVar("$sl_right0","7")
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("prefix", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
if isEqual(ValueFoldInterface("elseif"), ValueFoldInterface(getVar("prefix"))){
print("")
}else{
goto not_elseif
}
systemStack = append(systemStack, "elseif")
undefineVar("command_len")
undefineVar("prefix")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_res0
}
goto not_elseif
not_elseif:
fmt.Print("")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0","1")
setVar("$sl_right0","5")
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("prefix", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
if isEqual(ValueFoldInterface("else"), ValueFoldInterface(getVar("prefix"))){
print("")
}else{
goto if_type_error
}
systemStack = append(systemStack, "else")
undefineVar("prefix")
undefineVar("command_len")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_res0
}
goto if_type_error
if_type_error:
fmt.Print("")
systemStack = append(systemStack, "error")
undefineVar("command_len")
undefineVar("prefix")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$if_type_return_var")){
goto main_res0
}
goto if_type_end
if_type_end:
fmt.Print("")
defineVar("$switch_files_return_var")
setVar("$switch_files_return_var", "")
goto switch_files_end
goto switch_files
switch_files:
fmt.Print("")
fmt.Print("")
setVar("$finish_return_var","#finish_res0")
goto finish
goto finish_res0
finish_res0:
fmt.Print("")
fmt.Print("")
if toBool(getVar("first_file")){
print("")
}else{
goto first_end
}
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/if_program.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/if_program2.b"))))
setVar("first_file","False")
goto switch_files_e
goto first_end
first_end:
fmt.Print("")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/if_program2.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/if_program.b"))))
setVar("first_file","True")
goto switch_files_e
switch_files_e:
fmt.Print("")
fmt.Print("")
if "#reverse" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto main_res0
}
goto switch_files_end
switch_files_end:
fmt.Print("")
defineVar("$replace_if_return_var")
setVar("$replace_if_return_var", "")
goto replace_if_end
goto replace_if
replace_if:
fmt.Print("")
defineVar("stop_pos")
setVar("stop_pos", 0)
defineVar("cond")
setVar("cond", "")
fmt.Print("")
setVar("stop_pos", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("cond", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("buf")
setVar("buf", "")
defineVar("snum")
setVar("snum", "")
defineVar("args_to_undefine")
setVar("args_to_undefine", []interface{}{"end"})
defineVar("arg_type")
setVar("arg_type", "")
defineVar("command_len")
setVar("command_len", 0)
defineVar("arg_name")
setVar("arg_name", "")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("num"))
setVar("snum", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum(sum(sum(sum("[print(\"\"), ", getVar("cond")), ", goto(#_cond"), getVar("snum")), "_end)]"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res0")
goto switch_command
goto switch_command_res0
switch_command_res0:
fmt.Print("")
fmt.Print("")
goto replace_clear_if_s
replace_clear_if_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto replace_clear_if_e
}else{
print("")
}
if isEqual(ValueFoldInterface(getVar("stop_pos")), ValueFoldInterface(getVar("COMMAND_COUNTER"))){
print("")
}else{
goto add_replace_clear_if_mark
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un
un:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto un_end
}else{
print("")
}
setVar("buf",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un
goto un_end
un_end:
fmt.Print("")
setVar("buf",sum(sum("#_cond", getVar("snum")), "_end:print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res1")
goto switch_command
goto switch_command_res1
switch_command_res1:
fmt.Print("")
fmt.Print("")
goto add_replace_clear_if_mark
add_replace_clear_if_mark:
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res0")
goto check_br
goto check_br_res0
check_br_res0:
fmt.Print("")
fmt.Print("")
defineVar("$is_var_def_res0")
setVar("$is_var_def_res0", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_var_def_return_var","#is_var_def_res0")
goto is_var_def
goto is_var_def_res0
is_var_def_res0:
setVar("$is_var_def_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_var_def_res0", getVar("$is_var_def_res"))
if toBool(getVar("$is_var_def_res0"))&&toBool(isEqual(ValueFoldInterface(getVar("br_closed")), ValueFoldInterface(getVar("br_opened")))){
print("")
}else{
goto pop_e
}
undefineVar("$is_var_def_res0")
defineVar("$Type_res0")
setVar("$Type_res0", "")
systemStack = append(systemStack, getVar("command"))
setVar("$Type_return_var","#Type_res0")
goto Type
goto Type_res0
Type_res0:
setVar("$Type_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$Type_res0", getVar("$Type_res"))
setVar("arg_type", getVar("$Type_res0"))
undefineVar("$Type_res0")
defineVar("type_len")
setVar("type_len", 0)
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("arg_type"))))
setVar("type_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("type_len"))
setVar("$sl_right0", getVar("command_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("arg_name", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("args_to_undefine", append(getVar("args_to_undefine").([]interface{}), getVar("arg_name")))
goto pop_e
pop_e:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res2")
goto switch_command
goto switch_command_res2
switch_command_res2:
fmt.Print("")
fmt.Print("")
goto replace_clear_if_s
goto replace_clear_if_e
replace_clear_if_e:
fmt.Print("")
setVar("$reset_br_return_var","#reset_br_res0")
goto reset_br
goto reset_br_res0
reset_br_res0:
fmt.Print("")
fmt.Print("")
setVar("COMMAND_COUNTER","0")
setVar("$switch_files_return_var","#switch_files_res0")
goto switch_files
goto switch_files_res0
switch_files_res0:
fmt.Print("")
fmt.Print("")
undefineVar("type_len")
undefineVar("arg_name")
undefineVar("command_len")
undefineVar("arg_type")
undefineVar("args_to_undefine")
undefineVar("snum")
undefineVar("buf")
undefineVar("stop_pos")
undefineVar("cond")
if "#reverse" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$replace_if_return_var")){
goto main_res0
}
goto replace_if_end
replace_if_end:
fmt.Print("")
defineVar("$replace_elseif_return_var")
setVar("$replace_elseif_return_var", "")
goto replace_elseif_end
goto replace_elseif
replace_elseif:
fmt.Print("")
defineVar("stop_pos")
setVar("stop_pos", 0)
defineVar("cond")
setVar("cond", "")
fmt.Print("")
setVar("stop_pos", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("cond", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
fmt.Print("")
defineVar("buf")
setVar("buf", "")
defineVar("ibuf")
setVar("ibuf", 0)
defineVar("snum")
setVar("snum", "")
defineVar("t")
setVar("t", "")
defineVar("bcommand")
setVar("bcommand", "")
defineVar("this_stop_pos")
setVar("this_stop_pos", 0)
defineVar("sexit_num")
setVar("sexit_num", "")
defineVar("command_len")
setVar("command_len", 0)
defineVar("arg_name")
setVar("arg_name", "")
defineVar("arg_type")
setVar("arg_type", "")
defineVar("args_to_undefine")
setVar("args_to_undefine", []interface{}{"end"})
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("exit_num"))
setVar("sexit_num", getVar("$s0"))
undefineVar("$s0")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("num"))
setVar("snum", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum(sum(sum(sum("[print(\"\"), ", getVar("cond")), ", goto(#_cond"), getVar("snum")), "_end)]"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res3")
goto switch_command
goto switch_command_res3
switch_command_res3:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res1")
goto check_br
goto check_br_res1
check_br_res1:
fmt.Print("")
fmt.Print("")
goto replace_elseif_s
replace_elseif_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto replace_elseif_e
}else{
print("")
}
if isEqual(ValueFoldInterface(getVar("stop_pos")), ValueFoldInterface(getVar("COMMAND_COUNTER"))){
print("")
}else{
goto add_replace_elseif_mark
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un2
un2:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto un_end2
}else{
print("")
}
setVar("buf",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un2
goto un_end2
un_end2:
fmt.Print("")
setVar("$reset_br_return_var","#reset_br_res1")
goto reset_br
goto reset_br_res1
reset_br_res1:
fmt.Print("")
fmt.Print("")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("exit_num"))
setVar("sexit_num", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum(sum("goto(#_cond_exit", getVar("sexit_num")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf",sum(sum("#_cond", getVar("snum")), "_end:print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("num",sum(getVar("num"), 1))
defineVar("$if_type_res0")
setVar("$if_type_res0", "")
systemStack = append(systemStack, getVar("command"))
setVar("$if_type_return_var","#if_type_res0")
goto if_type
goto if_type_res0
if_type_res0:
setVar("$if_type_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$if_type_res0", getVar("$if_type_res"))
setVar("t", getVar("$if_type_res0"))
undefineVar("$if_type_res0")
if isEqual(ValueFoldInterface("elseif"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto find_block_e
}
defineVar("$block_end_res0")
setVar("$block_end_res0", 0)
setVar("$block_end_return_var","#block_end_res0")
goto block_end
goto block_end_res0
block_end_res0:
setVar("$block_end_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$block_end_res0", getVar("$block_end_res"))
setVar("this_stop_pos", getVar("$block_end_res0"))
undefineVar("$block_end_res0")
goto find_block_e
find_block_e:
fmt.Print("")
setVar("bcommand", getVar("command"))
setVar("$switch_command_return_var","#switch_command_res4")
goto switch_command
goto switch_command_res4
switch_command_res4:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res2")
goto check_br
goto check_br_res2
check_br_res2:
fmt.Print("")
fmt.Print("")
if !toBool(isEqual(ValueFoldInterface("elseif"), ValueFoldInterface(getVar("t")))){
goto replace_elseif_e
}else{
print("")
}
if isEqual(ValueFoldInterface("elseif"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto elif_end
}
defineVar("$get_cond_res0")
setVar("$get_cond_res0", "")
systemStack = append(systemStack, getVar("bcommand"))
setVar("$get_cond_return_var","#get_cond_res0")
goto get_cond
goto get_cond_res0
get_cond_res0:
setVar("$get_cond_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$get_cond_res0", getVar("$get_cond_res"))
setVar("cond", getVar("$get_cond_res0"))
undefineVar("$get_cond_res0")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("num"))
setVar("snum", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum(sum(sum(sum("[print(\"\"), ", getVar("cond")), ", goto(#_cond"), getVar("snum")), "_end)]"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("stop_pos", getVar("this_stop_pos"))
goto replace_elseif_s
goto elif_end
elif_end:
fmt.Print("")
fmt.Print("")
goto add_replace_elseif_mark
add_replace_elseif_mark:
fmt.Print("")
fmt.Print("")
defineVar("$is_var_def_res1")
setVar("$is_var_def_res1", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_var_def_return_var","#is_var_def_res1")
goto is_var_def
goto is_var_def_res1
is_var_def_res1:
setVar("$is_var_def_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_var_def_res1", getVar("$is_var_def_res"))
if toBool(getVar("$is_var_def_res1"))&&toBool(isEqual(ValueFoldInterface(getVar("br_closed")), ValueFoldInterface(getVar("br_opened")))){
print("")
}else{
goto pop_e2
}
undefineVar("$is_var_def_res1")
defineVar("$Type_res1")
setVar("$Type_res1", "")
systemStack = append(systemStack, getVar("command"))
setVar("$Type_return_var","#Type_res1")
goto Type
goto Type_res1
Type_res1:
setVar("$Type_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$Type_res1", getVar("$Type_res"))
setVar("arg_type", getVar("$Type_res1"))
undefineVar("$Type_res1")
defineVar("type_len")
setVar("type_len", 0)
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("arg_type"))))
setVar("type_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("type_len"))
setVar("$sl_right0", getVar("command_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("arg_name", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("args_to_undefine", append(getVar("args_to_undefine").([]interface{}), getVar("arg_name")))
goto pop_e2
pop_e2:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res5")
goto switch_command
goto switch_command_res5
switch_command_res5:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res3")
goto check_br
goto check_br_res3
check_br_res3:
fmt.Print("")
fmt.Print("")
goto replace_elseif_s
goto replace_elseif_e
replace_elseif_e:
fmt.Print("")
defineVar("$if_type_res1")
setVar("$if_type_res1", "")
systemStack = append(systemStack, getVar("bcommand"))
setVar("$if_type_return_var","#if_type_res1")
goto if_type
goto if_type_res1
if_type_res1:
setVar("$if_type_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$if_type_res1", getVar("$if_type_res"))
setVar("t", getVar("$if_type_res1"))
undefineVar("$if_type_res1")
if isEqual(ValueFoldInterface("else"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto restore_end
}
setVar("ibuf",toFloat(getVar("COMMAND_COUNTER"))-2)
systemStack = append(systemStack, getVar("ibuf"))
setVar("$SET_COMMAND_COUNTER_return_var","#SET_COMMAND_COUNTER_res2")
goto SET_COMMAND_COUNTER
goto SET_COMMAND_COUNTER_res2
SET_COMMAND_COUNTER_res2:
fmt.Print("")
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res6")
goto switch_command
goto switch_command_res6
switch_command_res6:
fmt.Print("")
fmt.Print("")
defineVar("$block_end_res1")
setVar("$block_end_res1", 0)
setVar("$block_end_return_var","#block_end_res1")
goto block_end
goto block_end_res1
block_end_res1:
setVar("$block_end_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$block_end_res1", getVar("$block_end_res"))
setVar("stop_pos", getVar("$block_end_res1"))
undefineVar("$block_end_res1")
setVar("$switch_command_return_var","#switch_command_res7")
goto switch_command
goto switch_command_res7
switch_command_res7:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res4")
goto check_br
goto check_br_res4
check_br_res4:
fmt.Print("")
fmt.Print("")
goto restore_end
restore_end:
fmt.Print("")
if isEqual(ValueFoldInterface("else"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto final_cond_end
}
goto final_cond_s
final_cond_s:
fmt.Print("")
if isEqual(ValueFoldInterface(getVar("stop_pos")), ValueFoldInterface(getVar("COMMAND_COUNTER"))){
goto final_cond_end
}else{
print("")
}
fmt.Print("")
defineVar("$is_var_def_res2")
setVar("$is_var_def_res2", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_var_def_return_var","#is_var_def_res2")
goto is_var_def
goto is_var_def_res2
is_var_def_res2:
setVar("$is_var_def_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_var_def_res2", getVar("$is_var_def_res"))
if toBool(getVar("$is_var_def_res2"))&&toBool(isEqual(ValueFoldInterface(getVar("br_closed")), ValueFoldInterface(getVar("br_opened")))){
print("")
}else{
goto pop_e3
}
undefineVar("$is_var_def_res2")
defineVar("$Type_res2")
setVar("$Type_res2", "")
systemStack = append(systemStack, getVar("command"))
setVar("$Type_return_var","#Type_res2")
goto Type
goto Type_res2
Type_res2:
setVar("$Type_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$Type_res2", getVar("$Type_res"))
setVar("arg_type", getVar("$Type_res2"))
undefineVar("$Type_res2")
defineVar("type_len")
setVar("type_len", 0)
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("arg_type"))))
setVar("type_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("type_len"))
setVar("$sl_right0", getVar("command_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("arg_name", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("args_to_undefine", append(getVar("args_to_undefine").([]interface{}), getVar("arg_name")))
goto pop_e3
pop_e3:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res8")
goto switch_command
goto switch_command_res8
switch_command_res8:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res5")
goto check_br
goto check_br_res5
check_br_res5:
fmt.Print("")
fmt.Print("")
goto final_cond_s
goto final_cond_end
final_cond_end:
fmt.Print("")
if isEqual(ValueFoldInterface("else"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto else_end
}
setVar("$switch_command_return_var","#switch_command_res9")
goto switch_command
goto switch_command_res9
switch_command_res9:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res6")
goto check_br
goto check_br_res6
check_br_res6:
fmt.Print("")
fmt.Print("")
goto else_end
else_end:
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un3
un3:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto un_end3
}else{
print("")
}
setVar("buf",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un3
goto un_end3
un_end3:
fmt.Print("")
setVar("$reset_br_return_var","#reset_br_res2")
goto reset_br
goto reset_br_res2
reset_br_res2:
fmt.Print("")
fmt.Print("")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("exit_num"))
setVar("sexit_num", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum(sum("#_cond_exit", getVar("sexit_num")), ":print(\"\")"))
setVar("exit_num",sum(getVar("exit_num"), 1))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
goto ts
ts:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto te
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res10")
goto switch_command
goto switch_command_res10
switch_command_res10:
fmt.Print("")
fmt.Print("")
goto ts
goto te
te:
fmt.Print("")
setVar("COMMAND_COUNTER","0")
setVar("$reset_br_return_var","#reset_br_res3")
goto reset_br
goto reset_br_res3
reset_br_res3:
fmt.Print("")
fmt.Print("")
setVar("$switch_files_return_var","#switch_files_res1")
goto switch_files
goto switch_files_res1
switch_files_res1:
fmt.Print("")
fmt.Print("")
undefineVar("type_len")
undefineVar("type_len")
undefineVar("args_to_undefine")
undefineVar("arg_type")
undefineVar("arg_name")
undefineVar("command_len")
undefineVar("sexit_num")
undefineVar("this_stop_pos")
undefineVar("bcommand")
undefineVar("t")
undefineVar("snum")
undefineVar("ibuf")
undefineVar("buf")
undefineVar("stop_pos")
undefineVar("cond")
if "#reverse" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$replace_elseif_return_var")){
goto main_res0
}
goto replace_elseif_end
replace_elseif_end:
fmt.Print("")
defineVar("$replace_else_return_var")
setVar("$replace_else_return_var", "")
goto replace_else_end
goto replace_else
replace_else:
fmt.Print("")
defineVar("stop_pos")
setVar("stop_pos", 0)
defineVar("cond")
setVar("cond", "")
fmt.Print("")
setVar("stop_pos", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("cond", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("buf")
setVar("buf", "")
defineVar("snum")
setVar("snum", "")
defineVar("sexit_num")
setVar("sexit_num", "")
defineVar("ibuf")
setVar("ibuf", 0)
defineVar("pos")
setVar("pos", 0)
defineVar("command_len")
setVar("command_len", 0)
defineVar("arg_name")
setVar("arg_name", "")
defineVar("arg_type")
setVar("arg_type", "")
defineVar("args_to_undefine")
setVar("args_to_undefine", []interface{}{"end"})
setVar("pos","-1")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("num"))
setVar("snum", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum(sum(sum(sum("[print(\"\"), ", getVar("cond")), ", goto(#_cond"), getVar("snum")), "_end)]"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res11")
goto switch_command
goto switch_command_res11
switch_command_res11:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res7")
goto check_br
goto check_br_res7
check_br_res7:
fmt.Print("")
fmt.Print("")
goto replace_else_s
replace_else_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto replace_else_e
}else{
print("")
}
if isEqual(ValueFoldInterface(getVar("stop_pos")), ValueFoldInterface(getVar("COMMAND_COUNTER"))){
print("")
}else{
goto add_replace_else_mark
}
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("exit_num"))
setVar("sexit_num", getVar("$s0"))
undefineVar("$s0")
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un4
un4:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto un_end4
}else{
print("")
}
setVar("buf",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un4
goto un_end4
un_end4:
fmt.Print("")
setVar("$reset_br_return_var","#reset_br_res4")
goto reset_br
goto reset_br_res4
reset_br_res4:
fmt.Print("")
fmt.Print("")
setVar("buf",sum(sum("goto(#_cond_exit", getVar("sexit_num")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf",sum(sum("#_cond", getVar("snum")), "_end:print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res12")
goto switch_command
goto switch_command_res12
switch_command_res12:
fmt.Print("")
fmt.Print("")
setVar("ibuf",toFloat(getVar("COMMAND_COUNTER"))-2)
systemStack = append(systemStack, getVar("ibuf"))
setVar("$SET_COMMAND_COUNTER_return_var","#SET_COMMAND_COUNTER_res3")
goto SET_COMMAND_COUNTER
goto SET_COMMAND_COUNTER_res3
SET_COMMAND_COUNTER_res3:
fmt.Print("")
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res13")
goto switch_command
goto switch_command_res13
switch_command_res13:
fmt.Print("")
fmt.Print("")
defineVar("$block_end_res2")
setVar("$block_end_res2", 0)
setVar("$block_end_return_var","#block_end_res2")
goto block_end
goto block_end_res2
block_end_res2:
setVar("$block_end_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$block_end_res2", getVar("$block_end_res"))
setVar("pos", getVar("$block_end_res2"))
undefineVar("$block_end_res2")
setVar("$switch_command_return_var","#switch_command_res14")
goto switch_command
goto switch_command_res14
switch_command_res14:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res8")
goto check_br
goto check_br_res8
check_br_res8:
fmt.Print("")
fmt.Print("")
goto add_replace_else_mark
add_replace_else_mark:
fmt.Print("")
if isEqual(ValueFoldInterface(getVar("pos")), ValueFoldInterface(getVar("COMMAND_COUNTER"))){
print("")
}else{
goto figure_brace_end
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un5
un5:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto un_end5
}else{
print("")
}
setVar("buf",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && len(getVar("args_to_undefine").([]interface{})) > 1{
setVar("arg_name", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("args_to_undefine")) && !isEqual("end", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]) && !isEqual("[end]", getVar("args_to_undefine").([]interface{})[len(getVar("args_to_undefine").([]interface{})) - 1]){
setVar("args_to_undefine", getVar("args_to_undefine").([]interface{})[:len(getVar("args_to_undefine").([]interface{})) - 1])
}
goto un5
goto un_end5
un_end5:
fmt.Print("")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("exit_num"))
setVar("sexit_num", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum(sum("#_cond_exit", getVar("sexit_num")), ":print(\"\")"))
setVar("exit_num",sum(getVar("exit_num"), 1))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res15")
goto switch_command
goto switch_command_res15
switch_command_res15:
fmt.Print("")
fmt.Print("")
goto ets
ets:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto ete
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res16")
goto switch_command
goto switch_command_res16
switch_command_res16:
fmt.Print("")
fmt.Print("")
goto ets
goto ete
ete:
fmt.Print("")
goto replace_else_e
goto figure_brace_end
figure_brace_end:
fmt.Print("")
defineVar("$is_var_def_res3")
setVar("$is_var_def_res3", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_var_def_return_var","#is_var_def_res3")
goto is_var_def
goto is_var_def_res3
is_var_def_res3:
setVar("$is_var_def_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_var_def_res3", getVar("$is_var_def_res"))
if toBool(getVar("$is_var_def_res3"))&&toBool(isEqual(ValueFoldInterface(getVar("br_closed")), ValueFoldInterface(getVar("br_opened")))){
print("")
}else{
goto pop_e4
}
undefineVar("$is_var_def_res3")
defineVar("$Type_res3")
setVar("$Type_res3", "")
systemStack = append(systemStack, getVar("command"))
setVar("$Type_return_var","#Type_res3")
goto Type
goto Type_res3
Type_res3:
setVar("$Type_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$Type_res3", getVar("$Type_res"))
setVar("arg_type", getVar("$Type_res3"))
undefineVar("$Type_res3")
defineVar("type_len")
setVar("type_len", 0)
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("arg_type"))))
setVar("type_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("type_len"))
setVar("$sl_right0", getVar("command_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("arg_name", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("args_to_undefine", append(getVar("args_to_undefine").([]interface{}), getVar("arg_name")))
goto pop_e4
pop_e4:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res17")
goto switch_command
goto switch_command_res17
switch_command_res17:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
setVar("$check_br_return_var","#check_br_res9")
goto check_br
goto check_br_res9
check_br_res9:
fmt.Print("")
fmt.Print("")
goto replace_else_s
goto replace_else_e
replace_else_e:
fmt.Print("")
setVar("COMMAND_COUNTER","0")
setVar("$reset_br_return_var","#reset_br_res5")
goto reset_br
goto reset_br_res5
reset_br_res5:
fmt.Print("")
fmt.Print("")
setVar("$switch_files_return_var","#switch_files_res2")
goto switch_files
goto switch_files_res2
switch_files_res2:
fmt.Print("")
fmt.Print("")
undefineVar("type_len")
undefineVar("args_to_undefine")
undefineVar("arg_type")
undefineVar("arg_name")
undefineVar("command_len")
undefineVar("pos")
undefineVar("ibuf")
undefineVar("sexit_num")
undefineVar("snum")
undefineVar("buf")
undefineVar("stop_pos")
undefineVar("cond")
if "#reverse" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$replace_else_return_var")){
goto main_res0
}
goto replace_else_end
replace_else_end:
fmt.Print("")
defineVar("$clear_files_return_var")
setVar("$clear_files_return_var", "")
goto clear_files_end
goto clear_files
clear_files:
fmt.Print("")
fmt.Print("")
if toBool(getVar("first_file")){
goto first_file
}else{
print("")
}
setVar("$switch_files_return_var","#switch_files_res3")
goto switch_files
goto switch_files_res3
switch_files_res3:
fmt.Print("")
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res18")
goto switch_command
goto switch_command_res18
switch_command_res18:
fmt.Print("")
fmt.Print("")
goto clear_files_s
clear_files_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto clear_files_e
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res19")
goto switch_command
goto switch_command_res19
switch_command_res19:
fmt.Print("")
fmt.Print("")
goto clear_files_s
goto first_file
first_file:
fmt.Print("")
fmt.Print("")
goto clear_files_e
clear_files_e:
fmt.Print("")
setVar("$finish_return_var","#finish_res1")
goto finish
goto finish_res1
finish_res1:
fmt.Print("")
fmt.Print("")
os.Remove(getRootSource("benv/if_program2.b"))
os.Remove(getRootSource(fmt.Sprintf("%v", getVar("root_source"))))
if "#reverse" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto main_res0
}
goto clear_files_end
clear_files_end:
fmt.Print("")
defineVar("$main_return_var")
setVar("$main_return_var", "")
goto main_end
goto main
main:
fmt.Print("")
fmt.Print("")
defineVar("buf")
setVar("buf", "")
defineVar("cond")
setVar("cond", "")
defineVar("counter")
setVar("counter", 0)
defineVar("buf_counter")
setVar("buf_counter", 0)
defineVar("snum")
setVar("snum", "")
defineVar("t")
setVar("t", "")
defineVar("is_stop")
setVar("is_stop", false)
setVar("first_file","True")
setVar("$init_return_var","#init_res0")
goto init
goto init_res0
init_res0:
fmt.Print("")
fmt.Print("")
setVar("num","0")
setVar("exit_num","0")
setVar("COMMAND_COUNTER","0")
goto again_s
again_s:
fmt.Print("")
setVar("is_stop","True")
systemStack = append(systemStack, getVar("command"))
setVar("$switch_command_return_var","#switch_command_res20")
goto switch_command
goto switch_command_res20
switch_command_res20:
fmt.Print("")
fmt.Print("")
goto main_s
main_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto main_e
}else{
print("")
}
defineVar("$is_if_res0")
setVar("$is_if_res0", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_if_return_var","#is_if_res0")
goto is_if
goto is_if_res0
is_if_res0:
setVar("$is_if_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_if_res0", getVar("$is_if_res"))
if toBool(getVar("$is_if_res0")){
print("")
}else{
goto next
}
undefineVar("$is_if_res0")
setVar("is_stop","False")
defineVar("$get_cond_res1")
setVar("$get_cond_res1", "")
systemStack = append(systemStack, getVar("command"))
setVar("$get_cond_return_var","#get_cond_res1")
goto get_cond
goto get_cond_res1
get_cond_res1:
setVar("$get_cond_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$get_cond_res1", getVar("$get_cond_res"))
setVar("cond", getVar("$get_cond_res1"))
undefineVar("$get_cond_res1")
defineVar("$block_end_res3")
setVar("$block_end_res3", 0)
setVar("$block_end_return_var","#block_end_res3")
goto block_end
goto block_end_res3
block_end_res3:
setVar("$block_end_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$block_end_res3", getVar("$block_end_res"))
setVar("counter", getVar("$block_end_res3"))
undefineVar("$block_end_res3")
defineVar("$get_command_res0")
setVar("$get_command_res0", "")
systemStack = append(systemStack, getVar("counter"))
setVar("$get_command_return_var","#get_command_res0")
goto get_command
goto get_command_res0
get_command_res0:
setVar("$get_command_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$get_command_res0", getVar("$get_command_res"))
setVar("buf", getVar("$get_command_res0"))
undefineVar("$get_command_res0")
defineVar("$if_type_res2")
setVar("$if_type_res2", "")
systemStack = append(systemStack, getVar("buf"))
setVar("$if_type_return_var","#if_type_res2")
goto if_type
goto if_type_res2
if_type_res2:
setVar("$if_type_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$if_type_res2", getVar("$if_type_res"))
setVar("t", getVar("$if_type_res2"))
undefineVar("$if_type_res2")
if isEqual(ValueFoldInterface("error"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto error_end
}
systemStack = append(systemStack, "ERROR in the if type\n")
setVar("$println_return_var","#println_res0")
goto println
goto println_res0
println_res0:
fmt.Print("")
fmt.Print("")
goto total_e
goto error_end
error_end:
fmt.Print("")
if isEqual(ValueFoldInterface("clear"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto if_end
}
systemStack = append(systemStack, getVar("cond"))
systemStack = append(systemStack, getVar("counter"))
setVar("$replace_if_return_var","#replace_if_res0")
goto replace_if
goto replace_if_res0
replace_if_res0:
fmt.Print("")
fmt.Print("")
setVar("num",sum(getVar("num"), 1))
goto main_e
goto if_end
if_end:
fmt.Print("")
if isEqual(ValueFoldInterface("elseif"), ValueFoldInterface(getVar("t"))){
print("")
}else{
goto elseif_end
}
systemStack = append(systemStack, getVar("cond"))
systemStack = append(systemStack, getVar("counter"))
setVar("$replace_elseif_return_var","#replace_elseif_res0")
goto replace_elseif
goto replace_elseif_res0
replace_elseif_res0:
fmt.Print("")
fmt.Print("")
goto main_e
goto elseif_end
elseif_end:
fmt.Print("")
systemStack = append(systemStack, getVar("cond"))
systemStack = append(systemStack, getVar("counter"))
setVar("$replace_else_return_var","#replace_else_res0")
goto replace_else
goto replace_else_res0
replace_else_res0:
fmt.Print("")
fmt.Print("")
setVar("num",sum(getVar("num"), 1))
goto main_e
goto next
next:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res21")
goto switch_command
goto switch_command_res21
switch_command_res21:
fmt.Print("")
fmt.Print("")
goto main_s
goto main_e
main_e:
fmt.Print("")
if toBool(getVar("is_stop")){
print("")
}else{
goto again_s
}
goto total_e
total_e:
fmt.Print("")
setVar("$clear_files_return_var","#clear_files_res0")
goto clear_files
goto clear_files_res0
clear_files_res0:
fmt.Print("")
fmt.Print("")
undefineVar("is_stop")
undefineVar("t")
undefineVar("snum")
undefineVar("buf_counter")
undefineVar("counter")
undefineVar("cond")
undefineVar("buf")
if "#reverse" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto finish_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto Type_end
}
if "#check_br" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res4
}
if "#br_opened_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto br_opened_end
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res5
}
if "#br_closed_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto br_closed_end
}
if "#check_br_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_end
}
if "#reset_br" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br
}
if "#reset_br_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br_end
}
if "#is_if" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_if
}
if "#ops_res6" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res6
}
if "#start_true" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto start_true
}
if "#end_true" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto end_true
}
if "#is_if_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_if_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_cond
}
if "#ops_res7" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res7
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_cond_end
}
if "#if_type" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto if_type
}
if "#not_clear" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto not_clear
}
if "#not_elseif" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto not_elseif
}
if "#if_type_error" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto if_type_error
}
if "#if_type_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto if_type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_end
}
if "#replace_if" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_if
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res0
}
if "#replace_clear_if_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_clear_if_s
}
if "#un" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un
}
if "#un_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un_end
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res1
}
if "#add_replace_clear_if_mark" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto add_replace_clear_if_mark
}
if "#check_br_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res0
}
if "#is_var_def_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_var_def_res0
}
if "#Type_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto Type_res0
}
if "#pop_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto pop_e
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res2
}
if "#replace_clear_if_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_clear_if_e
}
if "#reset_br_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br_res0
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_res0
}
if "#replace_if_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_if_end
}
if "#replace_elseif" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_elseif
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res3
}
if "#check_br_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res1
}
if "#replace_elseif_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_elseif_s
}
if "#un2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un2
}
if "#un_end2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un_end2
}
if "#reset_br_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br_res1
}
if "#if_type_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto if_type_res0
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_res0
}
if "#find_block_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto find_block_e
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res4
}
if "#check_br_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res2
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_cond_res0
}
if "#elif_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto elif_end
}
if "#add_replace_elseif_mark" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto add_replace_elseif_mark
}
if "#is_var_def_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_var_def_res1
}
if "#Type_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto Type_res1
}
if "#pop_e2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto pop_e2
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res5
}
if "#check_br_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res3
}
if "#replace_elseif_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_elseif_e
}
if "#if_type_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto if_type_res1
}
if "#SET_COMMAND_COUNTER_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto SET_COMMAND_COUNTER_res2
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res6
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res7
}
if "#check_br_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res4
}
if "#restore_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto restore_end
}
if "#final_cond_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto final_cond_s
}
if "#is_var_def_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_var_def_res2
}
if "#Type_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto Type_res2
}
if "#pop_e3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto pop_e3
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res8
}
if "#check_br_res5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res5
}
if "#final_cond_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto final_cond_end
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res9
}
if "#check_br_res6" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res6
}
if "#else_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto else_end
}
if "#un3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un3
}
if "#un_end3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un_end3
}
if "#reset_br_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br_res2
}
if "#ts" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ts
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res10
}
if "#te" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto te
}
if "#reset_br_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br_res3
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_res1
}
if "#replace_elseif_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_elseif_end
}
if "#replace_else" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_else
}
if "#switch_command_res11" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res11
}
if "#check_br_res7" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res7
}
if "#replace_else_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_else_s
}
if "#un4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un4
}
if "#un_end4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un_end4
}
if "#reset_br_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br_res4
}
if "#switch_command_res12" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res12
}
if "#SET_COMMAND_COUNTER_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto SET_COMMAND_COUNTER_res3
}
if "#switch_command_res13" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res13
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_res2
}
if "#switch_command_res14" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res14
}
if "#check_br_res8" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res8
}
if "#add_replace_else_mark" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto add_replace_else_mark
}
if "#un5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un5
}
if "#un_end5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto un_end5
}
if "#switch_command_res15" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res15
}
if "#ets" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ets
}
if "#switch_command_res16" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res16
}
if "#ete" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ete
}
if "#figure_brace_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto figure_brace_end
}
if "#is_var_def_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_var_def_res3
}
if "#Type_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto Type_res3
}
if "#pop_e4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto pop_e4
}
if "#switch_command_res17" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res17
}
if "#check_br_res9" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto check_br_res9
}
if "#replace_else_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_else_e
}
if "#reset_br_res5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reset_br_res5
}
if "#switch_files_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_res2
}
if "#replace_else_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_else_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto clear_files
}
if "#switch_files_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_res3
}
if "#switch_command_res18" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res18
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto clear_files_s
}
if "#switch_command_res19" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res19
}
if "#first_file" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto init_res0
}
if "#again_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto again_s
}
if "#switch_command_res20" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res20
}
if "#main_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto main_s
}
if "#is_if_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_if_res0
}
if "#get_cond_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_cond_res1
}
if "#block_end_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_res3
}
if "#get_command_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_command_res0
}
if "#if_type_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto if_type_res2
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto println_res0
}
if "#error_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto error_end
}
if "#replace_if_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_if_res0
}
if "#if_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto if_end
}
if "#replace_elseif_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_elseif_res0
}
if "#elseif_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto elseif_end
}
if "#replace_else_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_else_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next
}
if "#switch_command_res21" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res21
}
if "#main_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto main_e
}
if "#total_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto total_e
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto main_res0
}
goto main_end
main_end:
fmt.Print("")
setVar("$main_return_var","#main_res0")
goto main
goto main_res0
main_res0:
fmt.Print("")
fmt.Print("")
}
