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
if "#is_while" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto finish_res0
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
setVar("s_len",len(fmt.Sprintf("%v",getVar("s"))))
setVar("sub_len",len(fmt.Sprintf("%v",getVar("sub_s"))))
setVar("s_len_old", getVar("s_len"))
setVar("i",strings.Index(fmt.Sprintf("%v",getVar("s")), fmt.Sprintf("%v",getVar("sub_s"))))
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
setVar("s",getVar("s").(string)[toInt(getVar("pointer")):toInt(getVar("s_len"))])
setVar("s_len",len(fmt.Sprintf("%v",getVar("s"))))
setVar("i",strings.Index(fmt.Sprintf("%v",getVar("s")), fmt.Sprintf("%v",getVar("sub_s"))))
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
systemStack = append(systemStack, getVar("res"))
undefineVar("$reverse_res0")
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
if "#is_while" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto finish_res0
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
if "#is_while" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$SET_COMMAND_COUNTER_return_var")){
goto finish_res0
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
if "#is_while" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$get_command_return_var")){
goto finish_res0
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
if "#is_while" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$switch_command_return_var")){
goto finish_res0
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
if "#is_while" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$stack_len_return_var")){
goto finish_res0
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
setVar("op_num",toInt(getVar("buf")))
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
setVar("quotes",getVar("$regRes"))
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
setVar("num1",toInt(getVar("buf")))
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
setVar("num2",toInt(getVar("buf")))
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
systemStack = append(systemStack, getVar("res"))
undefineVar("$reverse_res1")
undefineVar("$indexes_res0")
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
if "#is_while" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$ops_return_var")){
goto finish_res0
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
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command",getVar("command").(string)[1:toInt(getVar("command_len"))])
setVar("code", getVar("command"))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("counter",sum(getVar("counter"), 1))
setVar("code",sum(getVar("code"), string(getVar("command").(string)[0])))
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
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_buf",getVar("command").(string)[1:toInt(getVar("command_len"))])
setVar("code",sum(getVar("code"), getVar("command_buf")))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("counter",sum(getVar("counter"), 1))
setVar("command_buf",string(getVar("command").(string)[0]))
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
undefineVar("$stack_len_res1")
undefineVar("$stack_len_res0")
undefineVar("$ops_res1")
undefineVar("$ops_res0")
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
if "#is_while" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$block_end_return_var")){
goto finish_res0
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
if "#is_while" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$println_return_var")){
goto finish_res0
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
if "#is_while" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$in_stack_return_var")){
goto finish_res0
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
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
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
goto _braces_s
_braces_s:
fmt.Print("")
if toFloat(getVar("pos"))<toFloat(getVar("command_len")){
print("")
}else{
goto _braces_e
}
setVar("spos",getVar("pos"))
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
undefineVar("$in_stack_res1")
undefineVar("$in_stack_res0")
undefineVar("$ops_res3")
undefineVar("$ops_res2")
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
if "#is_while" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$func_end_return_var")){
goto finish_res0
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
defineVar("$init_return_var")
setVar("$init_return_var", "")
goto init_end
goto init
init:
fmt.Print("")
fmt.Print("")
setVar("root_source","benv/prep_try_program.b")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", getVar("root_source")))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/while_program.b"))))
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
if "#is_while" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto finish_res0
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
os.Remove(getRootSource(fmt.Sprintf("%v", getVar("root_source"))))
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
if "#is_while" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto finish_res0
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
defineVar("$is_while_return_var")
setVar("$is_while_return_var", "")
defineVar("$is_while_res")
setVar("$is_while_res", false)
goto is_while_end
goto is_while
is_while:
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
setVar("buf","while(")
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
undefineVar("$ops_res4")
undefineVar("buf")
undefineVar("s")
undefineVar("command")
if "#reverse" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto reverse
}
if "#_reverse_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _reverse_s
}
if "#_reverse_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto indexes
}
if "#_indexes_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _indexes_s
}
if "#_indexes_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto indexes_end
}
if "#SET_COMMAND_COUNTER" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto SET_COMMAND_COUNTER
}
if "#_set_start" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _set_start
}
if "#_set_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _set_end
}
if "#SET_COMMAND_COUNTER_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto SET_COMMAND_COUNTER_end
}
if "#get_command" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto get_command
}
if "#_get_command_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _get_command_s
}
if "#_get_command_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _get_command_e
}
if "#SET_COMMAND_COUNTER_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto SET_COMMAND_COUNTER_res0
}
if "#get_command_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto get_command_end
}
if "#switch_command" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto switch_command
}
if "#switch_command_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto switch_command_end
}
if "#stack_len" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto stack_len
}
if "#_stack_len_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _stack_len_s
}
if "#_stack_len_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _stack_len_e
}
if "#stack_len_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto stack_len_end
}
if "#ops" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto indexes_res0
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _op_nums_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto reverse_res1
}
if "#ops_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops_end
}
if "#block_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto block_end
}
if "#_block_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _block_s
}
if "#ops_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops_res0
}
if "#ops_res1" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops_res1
}
if "#stack_len_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto stack_len_res0
}
if "#stack_len_res1" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto stack_len_res1
}
if "#_block_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _block_e
}
if "#SET_COMMAND_COUNTER_res1" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto SET_COMMAND_COUNTER_res1
}
if "#block_end_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto block_end_end
}
if "#println" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto println
}
if "#println_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto println_end
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _in_stack_e
}
if "#in_stack_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto in_stack_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto func_end
}
if "#ops_res2" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops_res2
}
if "#ops_res3" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops_res3
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _braces_s
}
if "#in_stack_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto in_stack_res0
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _o_plus_end
}
if "#in_stack_res1" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto in_stack_res1
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _braces_e
}
if "#func_end_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto func_end_end
}
if "#init" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto finish_end
}
if "#is_while" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$is_while_return_var")){
goto main_res0
}
goto is_while_end
is_while_end:
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
setVar("op","while(")
defineVar("$ops_res5")
setVar("$ops_res5", []interface{}{"end"})
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("op"))
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
setVar("start_pos",toInt(getVar("buf")))
setVar("start_pos",sum(getVar("start_pos"), 5))
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
setVar("start_pos",sum(getVar("start_pos"), 1))
setVar("buf",getVar("command").(string)[toInt(getVar("start_pos")):toInt(getVar("end_pos"))])
systemStack = append(systemStack, getVar("buf"))
undefineVar("$func_end_res0")
undefineVar("$ops_res5")
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
if "#is_while" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$get_cond_return_var")){
goto finish_res0
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
defineVar("$main_return_var")
setVar("$main_return_var", "")
goto main_end
goto main
main:
fmt.Print("")
fmt.Print("")
defineVar("buf")
setVar("buf", "")
setVar("$init_return_var","#init_res0")
goto init
goto init_res0
init_res0:
fmt.Print("")
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res0")
goto switch_command
goto switch_command_res0
switch_command_res0:
fmt.Print("")
fmt.Print("")
if true{
print("")
}else{
goto _cond0_end
}
fmt.Print("")
fmt.Print("")
goto for0
for0:
fmt.Print("")
if !toBool(isEqual(ValueFoldInterface(getVar("command")), ValueFoldInterface("end"))){
print("")
}else{
goto _cond1_end
}
defineVar("$is_while_res0")
setVar("$is_while_res0", false)
systemStack = append(systemStack, getVar("command"))
setVar("$is_while_return_var","#is_while_res0")
goto is_while
goto is_while_res0
is_while_res0:
setVar("$is_while_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$is_while_res0", getVar("$is_while_res"))
if toBool(getVar("$is_while_res0")){
print("")
}else{
goto _cond2_end
}
setVar("buf","for (print(\"\")")
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("buf","print(\"\")")
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
defineVar("$get_cond_res0")
setVar("$get_cond_res0", "")
systemStack = append(systemStack, getVar("command"))
setVar("$get_cond_return_var","#get_cond_res0")
goto get_cond
goto get_cond_res0
get_cond_res0:
setVar("$get_cond_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$get_cond_res0", getVar("$get_cond_res"))
getVar("$DEST").(*os.File).WriteString(getVar("$get_cond_res0").(string) + ";\n")
setVar("buf","print(\"\")){print(\"\")")
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
goto _cond_exit0
goto _cond2_end
_cond2_end:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
fmt.Print("")
goto _cond_exit0
_cond_exit0:
fmt.Print("")
fmt.Print("")
setVar("$switch_command_return_var","#switch_command_res1")
goto switch_command
goto switch_command_res1
switch_command_res1:
fmt.Print("")
fmt.Print("")
goto for0
goto _cond1_end
_cond1_end:
fmt.Print("")
goto for0_end
for0_end:
fmt.Print("")
goto _cond0_end
_cond0_end:
fmt.Print("")
fmt.Print("")
setVar("$finish_return_var","#finish_res0")
goto finish
goto finish_res0
finish_res0:
fmt.Print("")
fmt.Print("")
undefineVar("$get_cond_res0")
undefineVar("$is_while_res0")
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
if "#is_while" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_while
}
if "#ops_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res4
}
if "#is_while_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_while_end
}
if "#get_cond" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_cond
}
if "#ops_res5" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ops_res5
}
if "#func_end_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_end_res0
}
if "#get_cond_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_cond_end
}
if "#main" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto init_res0
}
if "#switch_command_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res0
}
if "#for0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto for0
}
if "#is_while_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_while_res0
}
if "#get_cond_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_cond_res0
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond2_end
}
if "#_cond_exit0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond_exit0
}
if "#switch_command_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto switch_command_res1
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond1_end
}
if "#for0_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto for0_end
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond0_end
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto finish_res0
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
