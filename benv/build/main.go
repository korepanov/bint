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
defineVar("command")
setVar("command", "")
defineVar("return_type")
setVar("return_type", "")
defineVar("func_name")
setVar("func_name", "")
defineVar("func_stack")
setVar("func_stack", []interface{}{"end"})
defineVar("bool_res")
setVar("bool_res", false)
defineVar("root_source")
setVar("root_source", "")
defineVar("$init_return_var")
setVar("$init_return_var", "")
defineVar("$init_res")
setVar("$init_res", 0)
goto init_end
goto init
init:
fmt.Print("")
fmt.Print("")
setVar("root_source","benv/nested_call_program.b")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", getVar("root_source")))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/long_function_program.b"))))
systemStack = append(systemStack, 0)
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
if "#reverse" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$init_return_var")){
goto replace_res0
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
defineVar("$finish_res")
setVar("$finish_res", 0)
goto finish_end
goto finish
finish:
fmt.Print("")
fmt.Print("")
os.Remove(getRootSource(fmt.Sprintf("%v", getVar("root_source"))))
systemStack = append(systemStack, 0)
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
if "#reverse" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$finish_return_var")){
goto replace_res0
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
goto reverse_s
reverse_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto reverse_e
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
goto reverse_s
goto reverse_e
reverse_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
undefineVar("res")
undefineVar("buf")
undefineVar("s")
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
if "#reverse" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$reverse_return_var")){
goto replace_res0
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
goto indexes_s
indexes_s:
fmt.Print("")
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("i"))){
goto indexes_e
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
goto indexes_s
goto indexes_e
indexes_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_return_var","#reverse_res0")
goto reverse
goto reverse_res0
reverse_res0:
setVar("$reverse_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("res", getVar("$reverse_res"))
systemStack = append(systemStack, getVar("res"))
undefineVar("sub_len")
undefineVar("s_len_old")
undefineVar("s_len")
undefineVar("pointer")
undefineVar("i")
undefineVar("res")
undefineVar("sub_s")
undefineVar("s")
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
if "#reverse" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$indexes_return_var")){
goto replace_res0
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
defineVar("$look_behind_return_var")
setVar("$look_behind_return_var", "")
defineVar("$look_behind_res")
setVar("$look_behind_res", []interface{}{"end"})
goto look_behind_end
goto look_behind
look_behind:
fmt.Print("")
defineVar("s")
setVar("s", "")
defineVar("reg")
setVar("reg", "")
fmt.Print("")
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("reg", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("st")
setVar("st", []interface{}{"end"})
defineVar("this")
setVar("this", []interface{}{"end"})
defineVar("res")
setVar("res", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("pos")
setVar("pos", 0)
defineVar("symbol")
setVar("symbol", "")
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(fmt.Sprintf("%v", getVar("reg"))).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("s"))), -1)
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
setVar("st", getVar("$stack_var0"))
undefineVar("$stack_var0")
goto look_behind_s
look_behind_s:
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("st")) && len(getVar("st").([]interface{})) > 1{
setVar("this", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("this")) == "[]interface {}"{
setVar("this", []interface{}{[]interface{}{"end"}})
}else{
setVar("this", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("st")) && !isEqual("end", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1]) && !isEqual("[end]", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1]){
setVar("st", getVar("st").([]interface{})[:len(getVar("st").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("this")) && len(getVar("this").([]interface{})) > 1{
setVar("buf", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("this")) && !isEqual("end", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1]) && !isEqual("[end]", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1]){
setVar("this", getVar("this").([]interface{})[:len(getVar("this").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto look_behind_e
}else{
print("")
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("pos", getVar("$I0"))
undefineVar("$I0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("pos"))){
print("")
}else{
goto is_not_zero
}
setVar("res", append(getVar("res").([]interface{}), "$"))
goto look_behind_s
goto is_not_zero
is_not_zero:
fmt.Print("")
setVar("pos",toFloat(getVar("pos"))-1)
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_internal0", getVar("pos"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",string(getVar("s").(string)[toInt(getVar("$sl_internal0"))]))
setVar("symbol", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("res", append(getVar("res").([]interface{}), getVar("symbol")))
goto look_behind_s
goto look_behind_e
look_behind_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_return_var","#reverse_res1")
goto reverse
goto reverse_res1
reverse_res1:
setVar("$reverse_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("res", getVar("$reverse_res"))
systemStack = append(systemStack, getVar("res"))
undefineVar("symbol")
undefineVar("pos")
undefineVar("buf")
undefineVar("res")
undefineVar("this")
undefineVar("st")
undefineVar("s")
undefineVar("reg")
if "#init" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$look_behind_return_var")){
goto main_res0
}
goto look_behind_end
look_behind_end:
fmt.Print("")
defineVar("$look_ahead_return_var")
setVar("$look_ahead_return_var", "")
defineVar("$look_ahead_res")
setVar("$look_ahead_res", []interface{}{"end"})
goto look_ahead_end
goto look_ahead
look_ahead:
fmt.Print("")
defineVar("s")
setVar("s", "")
defineVar("reg")
setVar("reg", "")
fmt.Print("")
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("reg", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("st")
setVar("st", []interface{}{"end"})
defineVar("this")
setVar("this", []interface{}{"end"})
defineVar("res")
setVar("res", []interface{}{"end"})
defineVar("buf")
setVar("buf", "")
defineVar("pos")
setVar("pos", 0)
defineVar("symbol")
setVar("symbol", "")
defineVar("last_pos")
setVar("last_pos", 0)
defineVar("$stack_var0")
setVar("$stack_var0", []interface{}{"end"})
defineVar("$regRes")
{intListList := compileRegexp(fmt.Sprintf("%v", getVar("reg"))).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("s"))), -1)
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
setVar("st", getVar("$stack_var0"))
undefineVar("$stack_var0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("s"))))
setVar("last_pos", getVar("$l0"))
undefineVar("$l0")
setVar("last_pos",toFloat(getVar("last_pos"))-1)
goto look_ahead_s
look_ahead_s:
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("st")) && len(getVar("st").([]interface{})) > 1{
setVar("this", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("this")) == "[]interface {}"{
setVar("this", []interface{}{[]interface{}{"end"}})
}else{
setVar("this", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("st")) && !isEqual("end", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1]) && !isEqual("[end]", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1]){
setVar("st", getVar("st").([]interface{})[:len(getVar("st").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("this")) && len(getVar("this").([]interface{})) > 1{
setVar("buf", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("this")) && !isEqual("end", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1]) && !isEqual("[end]", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1]){
setVar("this", getVar("this").([]interface{})[:len(getVar("this").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("this")) && len(getVar("this").([]interface{})) > 1{
setVar("buf", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("this")) && !isEqual("end", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1]) && !isEqual("[end]", getVar("this").([]interface{})[len(getVar("this").([]interface{})) - 1]){
setVar("this", getVar("this").([]interface{})[:len(getVar("this").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto look_ahead_e
}else{
print("")
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("pos", getVar("$I0"))
undefineVar("$I0")
if isEqual(ValueFoldInterface(getVar("last_pos")), ValueFoldInterface(getVar("pos"))){
print("")
}else{
goto ais_not_zero
}
setVar("res", append(getVar("res").([]interface{}), "$"))
goto look_ahead_s
goto ais_not_zero
ais_not_zero:
fmt.Print("")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_internal0", getVar("pos"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",string(getVar("s").(string)[toInt(getVar("$sl_internal0"))]))
setVar("symbol", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("res", append(getVar("res").([]interface{}), getVar("symbol")))
goto look_ahead_s
goto look_ahead_e
look_ahead_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_return_var","#reverse_res2")
goto reverse
goto reverse_res2
reverse_res2:
setVar("$reverse_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("res", getVar("$reverse_res"))
systemStack = append(systemStack, getVar("res"))
undefineVar("last_pos")
undefineVar("symbol")
undefineVar("pos")
undefineVar("buf")
undefineVar("res")
undefineVar("this")
undefineVar("st")
undefineVar("s")
undefineVar("reg")
if "#init" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$look_ahead_return_var")){
goto main_res0
}
goto look_ahead_end
look_ahead_end:
fmt.Print("")
defineVar("$func_call_index_return_var")
setVar("$func_call_index_return_var", "")
defineVar("$func_call_index_res")
setVar("$func_call_index_res", 0)
goto func_call_index_end
goto func_call_index
func_call_index:
fmt.Print("")
defineVar("func_name")
setVar("func_name", "")
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("func_name", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("reg")
setVar("reg", "")
defineVar("st")
setVar("st", []interface{}{"end"})
defineVar("st2")
setVar("st2", []interface{}{"end"})
defineVar("ist")
setVar("ist", []interface{}{"end"})
defineVar("res")
setVar("res", 0)
defineVar("buf")
setVar("buf", "")
defineVar("symbol")
setVar("symbol", "")
defineVar("symbol2")
setVar("symbol2", "")
defineVar("letter")
setVar("letter", false)
defineVar("digit")
setVar("digit", false)
setVar("reg",sum(sum("(?:", getVar("func_name")), ")"))
systemStack = append(systemStack, getVar("reg"))
systemStack = append(systemStack, getVar("command"))
setVar("$look_behind_return_var","#look_behind_res0")
goto look_behind
goto look_behind_res0
look_behind_res0:
setVar("$look_behind_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("st", getVar("$look_behind_res"))
systemStack = append(systemStack, getVar("reg"))
systemStack = append(systemStack, getVar("command"))
setVar("$look_ahead_return_var","#look_ahead_res0")
goto look_ahead
goto look_ahead_res0
look_ahead_res0:
setVar("$look_ahead_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("st2", getVar("$look_ahead_res"))
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("func_name"))
setVar("$indexes_return_var","#indexes_res0")
goto indexes
goto indexes_res0
indexes_res0:
setVar("$indexes_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("ist", getVar("$indexes_res"))
goto func_call_index_s
func_call_index_s:
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("st")) && len(getVar("st").([]interface{})) > 1{
setVar("symbol", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("symbol")) == "[]interface {}"{
setVar("symbol", []interface{}{[]interface{}{"end"}})
}else{
setVar("symbol", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("st")) && !isEqual("end", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1]) && !isEqual("[end]", getVar("st").([]interface{})[len(getVar("st").([]interface{})) - 1]){
setVar("st", getVar("st").([]interface{})[:len(getVar("st").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("symbol"))){
print("")
}else{
goto _cond0_end
}
systemStack = append(systemStack, -1)
undefineVar("digit")
undefineVar("letter")
undefineVar("symbol2")
undefineVar("symbol")
undefineVar("buf")
undefineVar("res")
undefineVar("ist")
undefineVar("st2")
undefineVar("st")
undefineVar("reg")
undefineVar("func_name")
undefineVar("command")
if "#init" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto main_res0
}
goto _cond0_end
_cond0_end:
fmt.Print("")
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("st2")) && len(getVar("st2").([]interface{})) > 1{
setVar("symbol2", getVar("st2").([]interface{})[len(getVar("st2").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("symbol2")) == "[]interface {}"{
setVar("symbol2", []interface{}{[]interface{}{"end"}})
}else{
setVar("symbol2", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("st2")) && !isEqual("end", getVar("st2").([]interface{})[len(getVar("st2").([]interface{})) - 1]) && !isEqual("[end]", getVar("st2").([]interface{})[len(getVar("st2").([]interface{})) - 1]){
setVar("st2", getVar("st2").([]interface{})[:len(getVar("st2").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("symbol2"))){
print("")
}else{
goto _cond1_end
}
defineVar("$print_arg0")
setVar("$print_arg0", "")
setVar("$print_arg0","long_function: func_call_index: ERROR\n")
fmt.Print(getVar("$print_arg0"))
undefineVar("$print_arg0")
os.Exit(1)
goto _cond1_end
_cond1_end:
fmt.Print("")
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("ist")) && len(getVar("ist").([]interface{})) > 1{
setVar("buf", getVar("ist").([]interface{})[len(getVar("ist").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("ist")) && !isEqual("end", getVar("ist").([]interface{})[len(getVar("ist").([]interface{})) - 1]) && !isEqual("[end]", getVar("ist").([]interface{})[len(getVar("ist").([]interface{})) - 1]){
setVar("ist", getVar("ist").([]interface{})[:len(getVar("ist").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
print("")
}else{
goto _cond2_end
}
defineVar("$print_arg0")
setVar("$print_arg0", "")
setVar("$print_arg0","func_call_index ERROR\n")
fmt.Print(getVar("$print_arg0"))
undefineVar("$print_arg0")
os.Exit(1)
goto _cond2_end
_cond2_end:
fmt.Print("")
fmt.Print("")
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("res", getVar("$I0"))
undefineVar("$I0")
defineVar("$let0")
setVar("$let0", false)
setVar("$let0",unicode.IsLetter([]rune(fmt.Sprintf("%v",getVar("symbol")))[0]))
setVar("letter", getVar("$let0"))
undefineVar("$let0")
defineVar("$dig0")
setVar("$dig0", false)
setVar("$dig0",unicode.IsDigit([]rune(fmt.Sprintf("%v",getVar("symbol")))[0]))
setVar("digit", getVar("$dig0"))
undefineVar("$dig0")
if !toBool(toBool(toBool(getVar("letter"))||toBool(getVar("digit")))||toBool(isEqual(ValueFoldInterface("_"), ValueFoldInterface(getVar("symbol"))))){
print("")
}else{
goto _cond3_end
}
defineVar("$let0")
setVar("$let0", false)
setVar("$let0",unicode.IsLetter([]rune(fmt.Sprintf("%v",getVar("symbol2")))[0]))
setVar("letter", getVar("$let0"))
undefineVar("$let0")
defineVar("$dig0")
setVar("$dig0", false)
setVar("$dig0",unicode.IsDigit([]rune(fmt.Sprintf("%v",getVar("symbol2")))[0]))
setVar("digit", getVar("$dig0"))
undefineVar("$dig0")
if !toBool(toBool(toBool(getVar("letter"))||toBool(getVar("digit")))||toBool(isEqual(ValueFoldInterface("_"), ValueFoldInterface(getVar("symbol2"))))){
print("")
}else{
goto _cond4_end
}
systemStack = append(systemStack, getVar("res"))
undefineVar("reg")
undefineVar("st")
undefineVar("st2")
undefineVar("ist")
undefineVar("res")
undefineVar("buf")
undefineVar("symbol")
undefineVar("symbol2")
undefineVar("letter")
undefineVar("digit")
undefineVar("command")
undefineVar("func_name")
if "#init" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$func_call_index_return_var")){
goto main_res0
}
goto _cond4_end
_cond4_end:
fmt.Print("")
fmt.Print("")
goto _cond3_end
_cond3_end:
fmt.Print("")
fmt.Print("")
goto func_call_index_s
goto func_call_index_end
func_call_index_end:
fmt.Print("")
defineVar("$func_calls_return_var")
setVar("$func_calls_return_var", "")
defineVar("$func_calls_res")
setVar("$func_calls_res", []interface{}{"end"})
goto func_calls_end
goto func_calls
func_calls:
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
systemStack = append(systemStack, getVar("s"))
systemStack = append(systemStack, getVar("sub_s"))
setVar("$func_call_index_return_var","#func_call_index_res0")
goto func_call_index
goto func_call_index_res0
func_call_index_res0:
setVar("$func_call_index_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("i", getVar("$func_call_index_res"))
setVar("pointer", getVar("i"))
goto findexes_s
findexes_s:
fmt.Print("")
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("i"))){
goto findexes_e
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
systemStack = append(systemStack, getVar("s"))
systemStack = append(systemStack, getVar("sub_s"))
setVar("$func_call_index_return_var","#func_call_index_res1")
goto func_call_index
goto func_call_index_res1
func_call_index_res1:
setVar("$func_call_index_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("i", getVar("$func_call_index_res"))
setVar("pointer", getVar("i"))
goto findexes_s
goto findexes_e
findexes_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_return_var","#reverse_res3")
goto reverse
goto reverse_res3
reverse_res3:
setVar("$reverse_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("res", getVar("$reverse_res"))
systemStack = append(systemStack, getVar("res"))
undefineVar("sub_len")
undefineVar("s_len_old")
undefineVar("s_len")
undefineVar("pointer")
undefineVar("i")
undefineVar("res")
undefineVar("sub_s")
undefineVar("s")
if "#init" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$func_calls_return_var")){
goto main_res0
}
goto func_calls_end
func_calls_end:
fmt.Print("")
defineVar("$next_func_return_var")
setVar("$next_func_return_var", "")
defineVar("$next_func_res")
setVar("$next_func_res", []interface{}{"end"})
goto next_func_end
goto next_func
next_func:
fmt.Print("")
fmt.Print("")
defineVar("number")
setVar("number", 0)
defineVar("left_border")
setVar("left_border", 0)
defineVar("right_border")
setVar("right_border", 0)
defineVar("func_stack")
setVar("func_stack", []interface{}{"end"})
defineVar("func_name")
setVar("func_name", "")
defineVar("type_len")
setVar("type_len", 0)
defineVar("command")
setVar("command", "")
defineVar("arg_type")
setVar("arg_type", "")
goto next_func_st
next_func_st:
fmt.Print("")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto end_file
}else{
print("")
}
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","{")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
goto end_clause
}else{
print("")
}
setVar("arg_type","int")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","int")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto next_func_e
}else{
print("")
}
setVar("arg_type","bool")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","bool")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto next_func_e
}else{
print("")
}
setVar("arg_type","float")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","float")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto next_func_e
}else{
print("")
}
setVar("arg_type","stack")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","stack")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto next_func_e
}else{
print("")
}
setVar("arg_type","string")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","string")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto next_func_e
}else{
print("")
}
goto end_clause
end_clause:
fmt.Print("")
goto next_func_st
goto next_func_e
next_func_e:
fmt.Print("")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("arg_type"))))
setVar("type_len", getVar("$l0"))
undefineVar("$l0")
setVar("left_border", getVar("type_len"))
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","(")))
setVar("right_border", getVar("$ind0"))
undefineVar("$ind0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("left_border"))
setVar("$sl_right0", getVar("right_border"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("func_name", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("func_stack", append(getVar("func_stack").([]interface{}), getVar("func_name")))
setVar("func_stack", append(getVar("func_stack").([]interface{}), getVar("arg_type")))
goto end_file
end_file:
fmt.Print("")
systemStack = append(systemStack, getVar("func_stack"))
undefineVar("arg_type")
undefineVar("command")
undefineVar("type_len")
undefineVar("func_name")
undefineVar("func_stack")
undefineVar("right_border")
undefineVar("left_border")
undefineVar("number")
if "#init" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$next_func_return_var")){
goto main_res0
}
goto next_func_end
next_func_end:
fmt.Print("")
defineVar("$get_funcs_return_var")
setVar("$get_funcs_return_var", "")
defineVar("$get_funcs_res")
setVar("$get_funcs_res", []interface{}{"end"})
goto get_funcs_end
goto get_funcs
get_funcs:
fmt.Print("")
fmt.Print("")
defineVar("res_stack")
setVar("res_stack", []interface{}{"end"})
defineVar("return_type")
setVar("return_type", "")
defineVar("func_name")
setVar("func_name", "")
goto get_funcs_s
get_funcs_s:
fmt.Print("")
setVar("$next_func_return_var","#next_func_res0")
goto next_func
goto next_func_res0
next_func_res0:
setVar("$next_func_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_stack", getVar("$next_func_res"))
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && len(getVar("func_stack").([]interface{})) > 1{
setVar("return_type", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("return_type")) == "[]interface {}"{
setVar("return_type", []interface{}{[]interface{}{"end"}})
}else{
setVar("return_type", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && !isEqual("end", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]){
setVar("func_stack", getVar("func_stack").([]interface{})[:len(getVar("func_stack").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && len(getVar("func_stack").([]interface{})) > 1{
setVar("func_name", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("func_name")) == "[]interface {}"{
setVar("func_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("func_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && !isEqual("end", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]){
setVar("func_stack", getVar("func_stack").([]interface{})[:len(getVar("func_stack").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("func_name"))){
goto get_funcs_e
}else{
print("")
}
setVar("res_stack", append(getVar("res_stack").([]interface{}), getVar("func_name")))
setVar("res_stack", append(getVar("res_stack").([]interface{}), getVar("return_type")))
goto get_funcs_s
goto get_funcs_e
get_funcs_e:
fmt.Print("")
getVar("$SOURCE").(*os.File).Seek(0, 0)
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
systemStack = append(systemStack, getVar("res_stack"))
undefineVar("func_name")
undefineVar("return_type")
undefineVar("res_stack")
if "#init" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$get_funcs_return_var")){
goto main_res0
}
goto get_funcs_end
get_funcs_end:
fmt.Print("")
defineVar("$func_ends_return_var")
setVar("$func_ends_return_var", "")
defineVar("$func_ends_res")
setVar("$func_ends_res", []interface{}{"end"})
goto func_ends_end
goto func_ends
func_ends:
fmt.Print("")
defineVar("func_len")
setVar("func_len", 0)
defineVar("func_begins")
setVar("func_begins", []interface{}{"end"})
defineVar("command")
setVar("command", "")
fmt.Print("")
setVar("func_len", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_begins", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("buf")
setVar("buf", "")
defineVar("symbol")
setVar("symbol", "")
defineVar("i")
setVar("i", 0)
defineVar("br_begin")
setVar("br_begin", 0)
defineVar("br_end")
setVar("br_end", 0)
defineVar("command_len")
setVar("command_len", 0)
defineVar("opened_braces")
setVar("opened_braces", 0)
defineVar("closed_braces")
setVar("closed_braces", 0)
defineVar("res")
setVar("res", []interface{}{"end"})
defineVar("temp")
setVar("temp", "")
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_begins")) && len(getVar("func_begins").([]interface{})) > 1{
setVar("buf", getVar("func_begins").([]interface{})[len(getVar("func_begins").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_begins")) && !isEqual("end", getVar("func_begins").([]interface{})[len(getVar("func_begins").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_begins").([]interface{})[len(getVar("func_begins").([]interface{})) - 1]){
setVar("func_begins", getVar("func_begins").([]interface{})[:len(getVar("func_begins").([]interface{})) - 1])
}
goto func_ends_s
func_ends_s:
fmt.Print("")
setVar("closed_braces","0")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("buf"))){
goto func_ends_e
}else{
print("")
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("buf")))
setVar("i", getVar("$I0"))
undefineVar("$I0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command_len", getVar("$l0"))
undefineVar("$l0")
setVar("br_begin",sum(getVar("i"), getVar("func_len")))
setVar("opened_braces","1")
setVar("br_end",sum(getVar("br_begin"), 1))
goto counter_s
counter_s:
fmt.Print("")
if isEqual(ValueFoldInterface(getVar("opened_braces")), ValueFoldInterface(getVar("closed_braces"))){
goto counter_e
}else{
print("")
}
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_internal0", getVar("br_end"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",string(getVar("command").(string)[toInt(getVar("$sl_internal0"))]))
setVar("symbol", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
if isEqual(ValueFoldInterface("("), ValueFoldInterface(getVar("symbol"))){
print("")
}else{
goto inc_o_end
}
setVar("opened_braces",sum(getVar("opened_braces"), 1))
goto inc_o_end
inc_o_end:
fmt.Print("")
if isEqual(ValueFoldInterface(")"), ValueFoldInterface(getVar("symbol"))){
print("")
}else{
goto inc_c_end
}
setVar("closed_braces",sum(getVar("closed_braces"), 1))
goto inc_c_end
inc_c_end:
fmt.Print("")
setVar("br_end",sum(getVar("br_end"), 1))
goto counter_s
goto counter_e
counter_e:
fmt.Print("")
setVar("res", append(getVar("res").([]interface{}), getVar("br_end")))
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_begins")) && len(getVar("func_begins").([]interface{})) > 1{
setVar("buf", getVar("func_begins").([]interface{})[len(getVar("func_begins").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("buf")) == "[]interface {}"{
setVar("buf", []interface{}{[]interface{}{"end"}})
}else{
setVar("buf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_begins")) && !isEqual("end", getVar("func_begins").([]interface{})[len(getVar("func_begins").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_begins").([]interface{})[len(getVar("func_begins").([]interface{})) - 1]){
setVar("func_begins", getVar("func_begins").([]interface{})[:len(getVar("func_begins").([]interface{})) - 1])
}
goto func_ends_s
goto func_ends_e
func_ends_e:
fmt.Print("")
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_return_var","#reverse_res4")
goto reverse
goto reverse_res4
reverse_res4:
setVar("$reverse_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("res", getVar("$reverse_res"))
systemStack = append(systemStack, getVar("res"))
undefineVar("temp")
undefineVar("res")
undefineVar("closed_braces")
undefineVar("opened_braces")
undefineVar("command_len")
undefineVar("br_end")
undefineVar("br_begin")
undefineVar("i")
undefineVar("symbol")
undefineVar("buf")
undefineVar("func_len")
undefineVar("func_begins")
undefineVar("command")
if "#init" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$func_ends_return_var")){
goto main_res0
}
goto func_ends_end
func_ends_end:
fmt.Print("")
defineVar("$del_file_return_var")
setVar("$del_file_return_var", "")
goto del_file_end
goto del_file
del_file:
fmt.Print("")
defineVar("change_flag")
setVar("change_flag", false)
fmt.Print("")
setVar("change_flag", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
fmt.Print("")
defineVar("command")
setVar("command", "")
if toBool(getVar("change_flag")){
print("")
}else{
goto copy_e
}
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/long_function_program2.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/long_function_program.b"))))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto copy_s
copy_s:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto copy_e
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto copy_s
goto copy_e
copy_e:
fmt.Print("")
os.Remove(getRootSource("benv/long_function_program2.b"))
undefineVar("command")
undefineVar("change_flag")
if "#init" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$del_file_return_var")){
goto main_res0
}
goto del_file_end
del_file_end:
fmt.Print("")
defineVar("$replace_return_var")
setVar("$replace_return_var", "")
goto replace_end
goto replace
replace:
fmt.Print("")
fmt.Print("")
defineVar("command")
setVar("command", "")
defineVar("command_to_send")
setVar("command_to_send", "")
defineVar("replaced_command")
setVar("replaced_command", "")
defineVar("left_part")
setVar("left_part", "")
defineVar("right_part")
setVar("right_part", "")
defineVar("command_len")
setVar("command_len", 0)
defineVar("number")
setVar("number", 0)
defineVar("func_pos")
setVar("func_pos", 0)
defineVar("itemp")
setVar("itemp", 0)
defineVar("offset")
setVar("offset", 0)
defineVar("func_pos_stack")
setVar("func_pos_stack", []interface{}{"end"})
defineVar("func_ends_stack")
setVar("func_ends_stack", []interface{}{"end"})
defineVar("change_flag")
setVar("change_flag", false)
defineVar("func_len")
setVar("func_len", 0)
defineVar("symbol")
setVar("symbol", "")
defineVar("return_type")
setVar("return_type", "")
defineVar("func_entry")
setVar("func_entry", 0)
defineVar("str_func_entry")
setVar("str_func_entry", "")
defineVar("sleft_border")
setVar("sleft_border", "")
defineVar("sright_border")
setVar("sright_border", "")
defineVar("left_border")
setVar("left_border", 0)
defineVar("right_border")
setVar("right_border", 0)
defineVar("left_border_reserv")
setVar("left_border_reserv", 0)
defineVar("right_border_reserv")
setVar("right_border_reserv", 0)
defineVar("func_call")
setVar("func_call", "")
defineVar("stemp")
setVar("stemp", "")
defineVar("itemp")
setVar("itemp", 0)
defineVar("temp")
setVar("temp", "")
defineVar("stemp_len")
setVar("stemp_len", 0)
defineVar("arg_type")
setVar("arg_type", "")
setVar("func_entry","0")
setVar("offset","0")
setVar("change_flag","False")
setVar("$get_funcs_return_var","#get_funcs_res0")
goto get_funcs
goto get_funcs_res0
get_funcs_res0:
setVar("$get_funcs_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_stack", getVar("$get_funcs_res"))
setVar("func_stack", append(getVar("func_stack").([]interface{}), "$temp"))
setVar("func_stack", append(getVar("func_stack").([]interface{}), "$temp"))
goto replace_s
replace_s:
fmt.Print("")
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && len(getVar("func_stack").([]interface{})) > 1{
setVar("return_type", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("return_type")) == "[]interface {}"{
setVar("return_type", []interface{}{[]interface{}{"end"}})
}else{
setVar("return_type", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && !isEqual("end", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]){
setVar("func_stack", getVar("func_stack").([]interface{})[:len(getVar("func_stack").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && len(getVar("func_stack").([]interface{})) > 1{
setVar("func_name", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("func_name")) == "[]interface {}"{
setVar("func_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("func_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && !isEqual("end", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]){
setVar("func_stack", getVar("func_stack").([]interface{})[:len(getVar("func_stack").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("func_name"))){
goto replace_e
}else{
print("")
}
goto next
next:
fmt.Print("")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto next_end
}else{
print("")
}
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("func_name"))
setVar("$func_call_index_return_var","#func_call_index_res2")
goto func_call_index
goto func_call_index_res2
func_call_index_res2:
setVar("$func_call_index_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("number", getVar("$func_call_index_res"))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto not_send
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto next
goto not_send
not_send:
fmt.Print("")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("func_name"))))
setVar("func_len", getVar("$l0"))
undefineVar("$l0")
setVar("number",sum(getVar("number"), getVar("func_len")))
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_internal0", getVar("number"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",string(getVar("command").(string)[toInt(getVar("$sl_internal0"))]))
setVar("symbol", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
if isEqual(ValueFoldInterface("("), ValueFoldInterface(getVar("symbol"))){
print("")
}else{
goto to_next_start
}
setVar("arg_type","int")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","int")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto to_next_start
}else{
print("")
}
setVar("arg_type","bool")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","bool")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto to_next_start
}else{
print("")
}
setVar("arg_type","float")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","float")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto to_next_start
}else{
print("")
}
setVar("arg_type","stack")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","stack")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto to_next_start
}else{
print("")
}
setVar("arg_type","string")
defineVar("$ind0")
setVar("$ind0", 0)
setVar("$ind0",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","string")))
setVar("number", getVar("$ind0"))
undefineVar("$ind0")
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto to_next_start
}else{
goto to_next_end
}
goto to_next_start
to_next_start:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto next
goto to_next_end
to_next_end:
fmt.Print("")
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("func_name"))
setVar("$func_calls_return_var","#func_calls_res0")
goto func_calls
goto func_calls_res0
func_calls_res0:
setVar("$func_calls_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_pos_stack", getVar("$func_calls_res"))
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("func_pos_stack"))
systemStack = append(systemStack, getVar("func_len"))
setVar("$func_ends_return_var","#func_ends_res0")
goto func_ends
goto func_ends_res0
func_ends_res0:
setVar("$func_ends_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_ends_stack", getVar("$func_ends_res"))
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_pos_stack")) && len(getVar("func_pos_stack").([]interface{})) > 1{
setVar("sleft_border", getVar("func_pos_stack").([]interface{})[len(getVar("func_pos_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("sleft_border")) == "[]interface {}"{
setVar("sleft_border", []interface{}{[]interface{}{"end"}})
}else{
setVar("sleft_border", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_pos_stack")) && !isEqual("end", getVar("func_pos_stack").([]interface{})[len(getVar("func_pos_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_pos_stack").([]interface{})[len(getVar("func_pos_stack").([]interface{})) - 1]){
setVar("func_pos_stack", getVar("func_pos_stack").([]interface{})[:len(getVar("func_pos_stack").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_ends_stack")) && len(getVar("func_ends_stack").([]interface{})) > 1{
setVar("sright_border", getVar("func_ends_stack").([]interface{})[len(getVar("func_ends_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("sright_border")) == "[]interface {}"{
setVar("sright_border", []interface{}{[]interface{}{"end"}})
}else{
setVar("sright_border", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_ends_stack")) && !isEqual("end", getVar("func_ends_stack").([]interface{})[len(getVar("func_ends_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_ends_stack").([]interface{})[len(getVar("func_ends_stack").([]interface{})) - 1]){
setVar("func_ends_stack", getVar("func_ends_stack").([]interface{})[:len(getVar("func_ends_stack").([]interface{})) - 1])
}
setVar("replaced_command", getVar("command"))
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("command"))))
setVar("itemp", getVar("$l0"))
undefineVar("$l0")
setVar("stemp_len","0")
goto pop_func_pos_start
pop_func_pos_start:
fmt.Print("")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("sleft_border"))){
goto pop_func_pos_end
}else{
print("")
}
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("sleft_border")))
setVar("left_border", getVar("$I0"))
undefineVar("$I0")
defineVar("$I0")
setVar("$I0", 0)
setVar("$I0",toInt(getVar("sright_border")))
setVar("right_border", getVar("$I0"))
undefineVar("$I0")
setVar("left_border_reserv", getVar("left_border"))
setVar("right_border_reserv", getVar("right_border"))
setVar("left_border",sum(getVar("left_border"), getVar("offset")))
setVar("right_border",sum(getVar("right_border"), getVar("offset")))
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("func_entry"))
setVar("str_func_entry", getVar("$s0"))
undefineVar("$s0")
setVar("command_to_send",sum(sum(sum(sum(getVar("return_type"), "$"), getVar("func_name")), "_res"), getVar("str_func_entry")))
getVar("$DEST").(*os.File).WriteString(getVar("command_to_send").(string) + ";\n")
setVar("command_to_send",sum(sum(sum(sum("$", getVar("func_name")), "_res"), getVar("str_func_entry")), "="))
setVar("func_entry",sum(getVar("func_entry"), 1))
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("left_border_reserv"))
setVar("$sl_right0", getVar("right_border_reserv"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("func_call", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0","0")
setVar("$sl_right0", getVar("left_border"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("replaced_command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("left_part", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
defineVar("$sl_internal0")
setVar("$sl_internal0", 0)
defineVar("$sl_left0")
setVar("$sl_left0", 0)
defineVar("$sl_right0")
setVar("$sl_right0", 0)
setVar("$sl_left0", getVar("right_border"))
setVar("$sl_right0", getVar("itemp"))
defineVar("$sl0")
setVar("$sl0", "")
setVar("$sl0",getVar("replaced_command").(string)[toInt(getVar("$sl_left0")):toInt(getVar("$sl_right0"))])
setVar("right_part", getVar("$sl0"))
undefineVar("$sl_internal0")
undefineVar("$sl0")
undefineVar("$sl_left0")
undefineVar("$sl_right0")
setVar("replaced_command",sum(sum(sum(sum(sum(getVar("left_part"), "$"), getVar("func_name")), "_res"), getVar("str_func_entry")), getVar("right_part")))
setVar("stemp",sum(sum(sum("$", getVar("func_name")), "_res"), getVar("str_func_entry")))
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("stemp"))))
setVar("stemp_len", getVar("$l0"))
undefineVar("$l0")
setVar("offset",sum(getVar("offset"), toFloat(getVar("stemp_len"))-toFloat(toFloat(getVar("right_border"))-toFloat(getVar("left_border")))))
defineVar("$s0")
setVar("$s0", "")
setVar("$s0",getVar("offset"))
setVar("temp", getVar("$s0"))
undefineVar("$s0")
defineVar("$l0")
setVar("$l0", 0)
setVar("$l0",len(fmt.Sprintf("%v",getVar("replaced_command"))))
setVar("itemp", getVar("$l0"))
undefineVar("$l0")
setVar("command_to_send",sum(getVar("command_to_send"), getVar("func_call")))
getVar("$DEST").(*os.File).WriteString(getVar("command_to_send").(string) + ";\n")
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_pos_stack")) && len(getVar("func_pos_stack").([]interface{})) > 1{
setVar("sleft_border", getVar("func_pos_stack").([]interface{})[len(getVar("func_pos_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("sleft_border")) == "[]interface {}"{
setVar("sleft_border", []interface{}{[]interface{}{"end"}})
}else{
setVar("sleft_border", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_pos_stack")) && !isEqual("end", getVar("func_pos_stack").([]interface{})[len(getVar("func_pos_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_pos_stack").([]interface{})[len(getVar("func_pos_stack").([]interface{})) - 1]){
setVar("func_pos_stack", getVar("func_pos_stack").([]interface{})[:len(getVar("func_pos_stack").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_ends_stack")) && len(getVar("func_ends_stack").([]interface{})) > 1{
setVar("sright_border", getVar("func_ends_stack").([]interface{})[len(getVar("func_ends_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("sright_border")) == "[]interface {}"{
setVar("sright_border", []interface{}{[]interface{}{"end"}})
}else{
setVar("sright_border", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_ends_stack")) && !isEqual("end", getVar("func_ends_stack").([]interface{})[len(getVar("func_ends_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_ends_stack").([]interface{})[len(getVar("func_ends_stack").([]interface{})) - 1]){
setVar("func_ends_stack", getVar("func_ends_stack").([]interface{})[:len(getVar("func_ends_stack").([]interface{})) - 1])
}
goto pop_func_pos_start
goto pop_func_pos_end
pop_func_pos_end:
fmt.Print("")
getVar("$DEST").(*os.File).WriteString(getVar("replaced_command").(string) + ";\n")
setVar("command",sum(sum("UNDEFINE(", sum(sum(sum("$", getVar("func_name")), "_res"), getVar("str_func_entry"))), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("offset","0")
goto next
goto next_end
next_end:
fmt.Print("")
getVar("$SOURCE").(*os.File).Close()
getVar("$DEST").(*os.File).Close()
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && len(getVar("func_stack").([]interface{})) > 1{
setVar("func_name", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("func_name")) == "[]interface {}"{
setVar("func_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("func_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("func_stack")) && !isEqual("end", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("func_stack").([]interface{})[len(getVar("func_stack").([]interface{})) - 1]){
setVar("func_stack", getVar("func_stack").([]interface{})[:len(getVar("func_stack").([]interface{})) - 1])
}
setVar("func_entry","0")
setVar("offset","0")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("func_name"))){
goto replace_e
}else{
print("")
}
setVar("func_stack", append(getVar("func_stack").([]interface{}), getVar("func_name")))
if toBool(getVar("change_flag")){
goto change
}else{
print("")
}
setVar("change_flag","True")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/long_function_program.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/long_function_program2.b"))))
goto replace_s
goto change
change:
fmt.Print("")
setVar("change_flag","False")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/long_function_program2.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/long_function_program.b"))))
goto replace_s
goto replace_e
replace_e:
fmt.Print("")
systemStack = append(systemStack, getVar("change_flag"))
setVar("$del_file_return_var","#del_file_res0")
goto del_file
goto del_file_res0
del_file_res0:
fmt.Print("")
fmt.Print("")
undefineVar("arg_type")
undefineVar("stemp_len")
undefineVar("temp")
undefineVar("itemp")
undefineVar("stemp")
undefineVar("func_call")
undefineVar("right_border_reserv")
undefineVar("left_border_reserv")
undefineVar("right_border")
undefineVar("left_border")
undefineVar("sright_border")
undefineVar("sleft_border")
undefineVar("str_func_entry")
undefineVar("func_entry")
undefineVar("return_type")
undefineVar("symbol")
undefineVar("func_len")
undefineVar("change_flag")
undefineVar("func_ends_stack")
undefineVar("func_pos_stack")
undefineVar("offset")
undefineVar("itemp")
undefineVar("func_pos")
undefineVar("number")
undefineVar("command_len")
undefineVar("right_part")
undefineVar("left_part")
undefineVar("replaced_command")
undefineVar("command_to_send")
undefineVar("command")
if "#init" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto init
}
if "#init_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto init_end
}
if "#finish" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto finish
}
if "#finish_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto finish_end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto replace_res0
}
if "#finish_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto finish_res0
}
if "#main_end" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto main_end
}
if "#main_res0" == fmt.Sprintf("%v", getVar("$replace_return_var")){
goto main_res0
}
goto replace_end
replace_end:
fmt.Print("")
defineVar("$main_return_var")
setVar("$main_return_var", "")
goto main_end
goto main
main:
fmt.Print("")
fmt.Print("")
defineVar("res")
setVar("res", 0)
setVar("$init_return_var","#init_res0")
goto init
goto init_res0
init_res0:
setVar("$init_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("res", getVar("$init_res"))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("res"))){
print("")
}else{
print("INIT ERROR\n")
}
setVar("$replace_return_var","#replace_res0")
goto replace
goto replace_res0
replace_res0:
fmt.Print("")
fmt.Print("")
setVar("$finish_return_var","#finish_res0")
goto finish
goto finish_res0
finish_res0:
setVar("$finish_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("res", getVar("$finish_res"))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("res"))){
print("")
}else{
print("FINISH ERROR\n")
}
undefineVar("res")
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
if "#reverse" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_e
}
if "#reverse_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_end
}
if "#indexes" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes_s
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes_e
}
if "#reverse_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_res0
}
if "#indexes_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes_end
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_behind_e
}
if "#reverse_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_res1
}
if "#look_behind_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_behind_end
}
if "#look_ahead" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_ahead
}
if "#look_ahead_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_ahead_s
}
if "#ais_not_zero" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto ais_not_zero
}
if "#look_ahead_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_ahead_e
}
if "#reverse_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_res2
}
if "#look_ahead_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_ahead_end
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_call_index
}
if "#look_behind_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_behind_res0
}
if "#look_ahead_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto look_ahead_res0
}
if "#indexes_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto indexes_res0
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_call_index_s
}
if "#_cond0_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond0_end
}
if "#_cond1_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond1_end
}
if "#_cond2_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond2_end
}
if "#_cond4_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond4_end
}
if "#_cond3_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto _cond3_end
}
if "#func_call_index_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_call_index_end
}
if "#func_calls" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_calls
}
if "#func_call_index_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_call_index_res0
}
if "#findexes_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto findexes_s
}
if "#func_call_index_res1" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_call_index_res1
}
if "#findexes_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto findexes_e
}
if "#reverse_res3" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_res3
}
if "#func_calls_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_calls_end
}
if "#next_func" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next_func
}
if "#next_func_st" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next_func_st
}
if "#end_clause" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto end_clause
}
if "#next_func_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next_func_e
}
if "#end_file" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto end_file
}
if "#next_func_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next_func_end
}
if "#get_funcs" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_funcs
}
if "#get_funcs_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_funcs_s
}
if "#next_func_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next_func_res0
}
if "#get_funcs_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_funcs_e
}
if "#get_funcs_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_funcs_end
}
if "#func_ends" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_ends
}
if "#func_ends_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_ends_s
}
if "#counter_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto counter_s
}
if "#inc_o_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto inc_o_end
}
if "#inc_c_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto inc_c_end
}
if "#counter_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto counter_e
}
if "#func_ends_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_ends_e
}
if "#reverse_res4" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto reverse_res4
}
if "#func_ends_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_ends_end
}
if "#del_file" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto del_file
}
if "#copy_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto copy_s
}
if "#copy_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto copy_e
}
if "#del_file_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto del_file_end
}
if "#replace" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace
}
if "#get_funcs_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto get_funcs_res0
}
if "#replace_s" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_s
}
if "#next" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next
}
if "#func_call_index_res2" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_call_index_res2
}
if "#not_send" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto not_send
}
if "#to_next_start" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto to_next_start
}
if "#to_next_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto to_next_end
}
if "#func_calls_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_calls_res0
}
if "#func_ends_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto func_ends_res0
}
if "#pop_func_pos_start" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto pop_func_pos_start
}
if "#pop_func_pos_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto pop_func_pos_end
}
if "#next_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto next_end
}
if "#change" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto change
}
if "#replace_e" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_e
}
if "#del_file_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto del_file_res0
}
if "#replace_end" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_end
}
if "#main" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto main
}
if "#init_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto init_res0
}
if "#replace_res0" == fmt.Sprintf("%v", getVar("$main_return_var")){
goto replace_res0
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
