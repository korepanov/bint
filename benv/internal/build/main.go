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
if "#is_for" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond11_end
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
defineVar("num")
setVar("num", 0)
defineVar("first_file")
setVar("first_file", false)
defineVar("condition")
setVar("condition", "")
defineVar("$init_return_var")
setVar("$init_return_var", "")
goto init_end
goto init
init:
fmt.Print("")
fmt.Print("")
setVar("num","0")
setVar("first_file","True")
setVar("root_source","benv/internal/while_program.b")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", getVar("root_source")))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/internal/for_program.b"))))
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
if "#is_for" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond11_end
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
defineVar("$is_for_return_var")
setVar("$is_for_return_var", "")
defineVar("$is_for_res")
setVar("$is_for_res", false)
goto is_for_end
goto is_for
is_for:
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
defineVar("pos")
setVar("pos", 0)
setVar("buf","for(")
defineVar("$ops_res4")
setVar("$ops_res4", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("buf"))
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
if !toBool(isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf")))){
print("")
}else{
goto _cond0_end
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("pos", getVar("$I0"))
undefineVar("$I0")
if !toBool(isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("pos")))){
print("")
}else{
goto _cond1_end
}
systemStack = append(systemStack, "for: ERROR")
setVar("$println_return_var","#println_res0")
goto println
goto println_res0
println_res0:
fmt.Print("")
fmt.Print("")
os.Exit(1)
goto _cond1_end
_cond1_end:
fmt.Print("")
fmt.Print("")
goto _cond0_end
_cond0_end:
fmt.Print("")
fmt.Print("")
systemStack = append(systemStack, !toBool(isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf")))))
undefineVar("pos")
undefineVar("buf")
undefineVar("s")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto finish_end
}
if "#is_for" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto is_for_end
}
if "#is_var_def" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto is_var_def
}
if "#is_var_def_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto is_var_def_end
}
if "#Type" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto Type
}
if "#int_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto int_end
}
if "#float_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto float_end
}
if "#bool_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto bool_end
}
if "#stack_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto stack_end
}
if "#string_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto string_end
}
if "#void_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto void_end
}
if "#Type_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto Type_end
}
if "#switch_files" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_files
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto finish_res0
}
if "#first_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto first_end
}
if "#switch_files_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_files_e
}
if "#switch_files_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_files_end
}
if "#clear_files" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res1
}
if "#first_file" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto first_file
}
if "#clear_files_e" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto clear_files_e
}
if "#finish_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto finish_res1
}
if "#clear_files_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto clear_files_end
}
if "#main" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto init_res0
}
if "#next" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto _cond11_end
}
if "#clear_files_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto clear_files_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$is_for_return_var")){
goto main_res0
}
goto is_for_end
is_for_end:
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
if "#is_for" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$is_var_def_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond11_end
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
if "#is_for" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$Type_return_var")){
goto _cond11_end
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
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/internal/for_program.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/internal/for_program2.b"))))
setVar("first_file","False")
goto switch_files_e
goto first_end
first_end:
fmt.Print("")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/internal/for_program2.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/internal/for_program.b"))))
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
if "#is_for" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$switch_files_return_var")){
goto _cond11_end
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
setVar("$switch_files_return_var","#switch_files_res0")
goto switch_files
goto switch_files_res0
switch_files_res0:
fmt.Print("")
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res0")
goto switch_command
goto switch_command_res0
switch_command_res0:
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
setVar("$switch_command_return_var","#switch_command_res1")
goto switch_command
goto switch_command_res1
switch_command_res1:
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
os.Remove(getRootSource("benv/internal/for_program2.b"))
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
if "#is_for" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$clear_files_return_var")){
goto _cond11_end
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
setVar("$init_return_var","#init_res0")
goto init
goto init_res0
init_res0:
fmt.Print("")
fmt.Print("")
defineVar("counter")
setVar("counter", 0)
defineVar("internal_counter")
setVar("internal_counter", 0)
defineVar("command_len")
setVar("command_len", 0)
defineVar("command_len2")
setVar("command_len2", 0)
defineVar("buf_len")
setVar("buf_len", 0)
defineVar("old_num")
setVar("old_num", 0)
defineVar("snum")
setVar("snum", "")
defineVar("buf")
setVar("buf", "")
defineVar("inc")
setVar("inc", "")
defineVar("old_command")
setVar("old_command", "")
defineVar("internal_old_command")
setVar("internal_old_command", "")
defineVar("pos")
setVar("pos", 0)
defineVar("was_for")
setVar("was_for", false)
defineVar("was_internal_for")
setVar("was_internal_for", false)
defineVar("arg_type")
setVar("arg_type", "")
defineVar("arg_name")
setVar("arg_name", "")
setVar("was_for","False")
setVar("was_internal_for","False")
goto next
next:
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res2")
goto switch_command
goto switch_command_res2
switch_command_res2:
fmt.Print("")
fmt.Print("")
if !toBool(isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command")))){
print("")
}else{
goto _cond2_end
}
defineVar("$is_for_res0")
setVar("$is_for_res0", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_for_return_var","#is_for_res0")
goto is_for
goto is_for_res0
is_for_res0:
setVar("$is_for_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_for_res0", getVar("$is_for_res"))
if toBool(getVar("$is_for_res0")){
print("")
}else{
goto _cond3_end
}
undefineVar("$is_for_res0")
setVar("was_for","True")
setVar("was_internal_for","False")
setVar("buf","if(True){print(\"\")")
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
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
setVar("$sl_left0","4")
setVar("$sl_right0", getVar("command_len"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("command", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("old_command", getVar("command"))
setVar("$switch_command_return_var","#switch_command_res3")
goto switch_command
goto switch_command_res3
switch_command_res3:
fmt.Print("")
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res4")
goto switch_command
goto switch_command_res4
switch_command_res4:
fmt.Print("")
fmt.Print("")
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("num"))
setVar("snum", getVar("$s0"))
undefineVar("$s0")
setVar("buf",sum("bool $for", getVar("snum")))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf",sum(sum("#_for", getVar("snum")), ":print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf",sum(sum("if(", getVar("command")), "){print(\"\")"))
setVar("condition", getVar("buf"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res5")
goto switch_command
goto switch_command_res5
switch_command_res5:
fmt.Print("")
fmt.Print("")
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
setVar("counter", getVar("$block_end_res0"))
undefineVar("$block_end_res0")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","{")))
setVar("pos", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("pos"))){
print("")
}else{
goto _cond4_end
}
systemStack = append(systemStack, "for: ERROR")
setVar("$println_return_var","#println_res1")
goto println
goto println_res1
println_res1:
fmt.Print("")
fmt.Print("")
os.Exit(1)
goto _cond4_end
_cond4_end:
fmt.Print("")
fmt.Print("")
setVar("pos",toFloat(getVar("pos"))-1)
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0","0")
setVar("$sl_right0", getVar("pos"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("inc", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
goto next_internal
next_internal:
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res6")
goto switch_command
goto switch_command_res6
switch_command_res6:
fmt.Print("")
fmt.Print("")
if toFloat(getVar("COMMAND_COUNTER"))<toFloat(getVar("counter")){
print("")
}else{
goto _cond5_end
}
defineVar("$is_for_res1")
setVar("$is_for_res1", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_for_return_var","#is_for_res1")
goto is_for
goto is_for_res1
is_for_res1:
setVar("$is_for_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_for_res1", getVar("$is_for_res"))
if toBool(getVar("$is_for_res1")){
print("")
}else{
goto _cond6_end
}
undefineVar("$is_for_res1")
setVar("was_internal_for","True")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res7")
goto switch_command
goto switch_command_res7
switch_command_res7:
fmt.Print("")
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res8")
goto switch_command
goto switch_command_res8
switch_command_res8:
fmt.Print("")
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$switch_command_return_var","#switch_command_res9")
goto switch_command
goto switch_command_res9
switch_command_res9:
fmt.Print("")
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("$block_end_return_var","#block_end_res1")
goto block_end
goto block_end_res1
block_end_res1:
setVar("$block_end_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("internal_counter", getVar("$block_end_res"))
setVar("$switch_command_return_var","#switch_command_res10")
goto switch_command
goto switch_command_res10
switch_command_res10:
fmt.Print("")
fmt.Print("")
defineVar("$block_end_res1")
setVar("$block_end_res1", 0)
setVar("$block_end_return_var","#block_end_res2")
goto block_end
goto block_end_res2
block_end_res2:
setVar("$block_end_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$block_end_res1", getVar("$block_end_res"))
undefineVar("0")
undefineVar("$block_end_res1")
goto _cond6_end
_cond6_end:
fmt.Print("")
fmt.Print("")
if toBool(getVar("was_internal_for")){
print("")
}else{
goto _cond7_end
}
if toFloat(getVar("COMMAND_COUNTER"))>toFloat(getVar("internal_counter")){
print("")
}else{
goto _cond8_end
}
setVar("was_internal_for","False")
goto _cond8_end
_cond8_end:
fmt.Print("")
fmt.Print("")
goto _cond7_end
_cond7_end:
fmt.Print("")
fmt.Print("")
if toBool(isEqual(ValueFoldInterface("break"), ValueFoldInterface(getVar("command"))))&&toBool(!toBool(getVar("was_internal_for"))){
print("")
}else{
goto _cond9_end
}
setVar("command",sum(sum("$for", getVar("snum")), "=True"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("command",sum(sum("goto(#_undef_for", getVar("snum")), ")"))
goto _cond9_end
_cond9_end:
fmt.Print("")
fmt.Print("")
if toBool(isEqual(ValueFoldInterface("continue"), ValueFoldInterface(getVar("command"))))&&toBool(!toBool(getVar("was_internal_for"))){
print("")
}else{
goto _cond10_end
}
getVar("$DEST").(*os.File).WriteString(getVar("inc").(string) + ";\n")
setVar("command",sum(sum("goto(#_undef_for", getVar("snum")), ")"))
goto _cond10_end
_cond10_end:
fmt.Print("")
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto next_internal
goto _cond5_end
_cond5_end:
fmt.Print("")
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("inc").(string) + ";\n")
setVar("buf",sum(sum("#_undef_for", getVar("snum")), ":print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("buf",sum(sum("if($for", getVar("snum")), "){print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf",sum(sum("goto(#_for", getVar("snum")), "_end)"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf","}")
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
getVar("$DEST").(*os.File).WriteString(getVar("condition").(string) + ";\n")
setVar("buf","CLEAR()")
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf",sum(sum("goto(#_for", getVar("snum")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf","}")
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("command",sum(sum("#_for", getVar("snum")), "_end:print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
fmt.Print("")
fmt.Print("")
setVar("command","}")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
fmt.Print("")
setVar("old_num", getVar("num"))
setVar("num",sum(getVar("num"), 1))
goto _cond_exit0
goto _cond3_end
_cond3_end:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
fmt.Print("")
goto _cond_exit0
_cond_exit0:
fmt.Print("")
fmt.Print("")
goto next
goto _cond2_end
_cond2_end:
fmt.Print("")
fmt.Print("")
if toBool(getVar("was_for")){
print("")
}else{
goto _cond11_end
}
setVar("was_for","False")
setVar("$switch_files_return_var","#switch_files_res1")
goto switch_files
goto switch_files_res1
switch_files_res1:
fmt.Print("")
fmt.Print("")
setVar("COMMAND_COUNTER","0")
goto next
goto _cond11_end
_cond11_end:
fmt.Print("")
fmt.Print("")
setVar("$clear_files_return_var","#clear_files_res0")
goto clear_files
goto clear_files_res0
clear_files_res0:
fmt.Print("")
fmt.Print("")
undefineVar("0")
undefineVar("arg_name")
undefineVar("arg_type")
undefineVar("was_internal_for")
undefineVar("was_for")
undefineVar("pos")
undefineVar("internal_old_command")
undefineVar("old_command")
undefineVar("inc")
undefineVar("buf")
undefineVar("snum")
undefineVar("old_num")
undefineVar("buf_len")
undefineVar("command_len2")
undefineVar("command_len")
undefineVar("internal_counter")
undefineVar("counter")
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
if "#is_for" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_for
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res4
}
if "#println_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto println_res0
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond1_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond0_end
}
if "#is_for_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_for_end
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
if "#clear_files" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto clear_files
}
if "#switch_files_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res0
}
if "#clear_files_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto clear_files_s
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res1
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
if "#next" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next
}
if "#switch_command_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res2
}
if "#is_for_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_for_res0
}
if "#switch_command_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res3
}
if "#switch_command_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res4
}
if "#switch_command_res5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res5
}
if "#block_end_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_res0
}
if "#println_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto println_res1
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond4_end
}
if "#next_internal" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next_internal
}
if "#switch_command_res6" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res6
}
if "#is_for_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_for_res1
}
if "#switch_command_res7" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res7
}
if "#switch_command_res8" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res8
}
if "#switch_command_res9" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res9
}
if "#block_end_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_res1
}
if "#switch_command_res10" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res10
}
if "#block_end_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto block_end_res2
}
if "#_cond6_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond6_end
}
if "#_cond8_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond8_end
}
if "#_cond7_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond7_end
}
if "#_cond9_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond9_end
}
if "#_cond10_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond10_end
}
if "#_cond5_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond5_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond3_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond_exit0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond2_end
}
if "#switch_files_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_files_res1
}
if "#_cond11_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond11_end
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
