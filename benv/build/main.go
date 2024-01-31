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
	
	
defineVar("command")
setVar("command", "")
defineVar("i")
setVar("i", 0)
defineVar("j")
setVar("j", 0)
defineVar("right_border")
setVar("right_border", 0)
defineVar("command_len")
setVar("command_len", 0)
defineVar("call_number")
setVar("call_number", 0)
defineVar("str_call_number")
setVar("str_call_number", "")
defineVar("func_name_br")
setVar("func_name_br", "")
defineVar("return_type")
setVar("return_type", "")
defineVar("buf")
setVar("buf", "")
defineVar("buf2")
setVar("buf2", "")
defineVar("buf3")
setVar("buf3", "")
defineVar("buf4")
setVar("buf4", "")
defineVar("rbuf")
setVar("rbuf", "")
defineVar("symbol")
setVar("symbol", "")
defineVar("$reverse_ret")
setVar("$reverse_ret", "")
defineVar("indexes_ret")
setVar("indexes_ret", "")
defineVar("dest")
setVar("dest", "")
defineVar("int_buf")
setVar("int_buf", 0)
defineVar("int_buf2")
setVar("int_buf2", 0)
defineVar("int_buf3")
setVar("int_buf3", 0)
defineVar("int_buf4")
setVar("int_buf4", 0)
defineVar("int_buf5")
setVar("int_buf5", 0)
defineVar("type_len")
setVar("type_len", 0)
defineVar("arg_type_len")
setVar("arg_type_len", 0)
defineVar("func_begin")
setVar("func_begin", 0)
defineVar("internal_func_begin")
setVar("internal_func_begin", 0)
defineVar("func_end")
setVar("func_end", 0)
defineVar("func_len")
setVar("func_len", 0)
defineVar("limit_border")
setVar("limit_border", 0)
defineVar("number")
setVar("number", 0)
defineVar("func_len")
setVar("func_len", 0)
defineVar("func_name_len")
setVar("func_name_len", 0)
defineVar("buf_len")
setVar("buf_len", 0)
defineVar("br_number")
setVar("br_number", 0)
defineVar("fig_br_number")
setVar("fig_br_number", 0)
defineVar("func_pos")
setVar("func_pos", 0)
defineVar("func_name")
setVar("func_name", "")
defineVar("arg_type")
setVar("arg_type", "")
defineVar("arg_name")
setVar("arg_name", "")
defineVar("$in_stack_ret")
setVar("$in_stack_ret", "")
defineVar("left_border")
setVar("left_border", 0)
defineVar("comma_pos")
setVar("comma_pos", 0)
defineVar("user_stack")
setVar("user_stack", []interface{}{"end"})
defineVar("buf_stack")
setVar("buf_stack", []interface{}{"end"})
defineVar("buf_stack_reserv")
setVar("buf_stack_reserv", []interface{}{"end"})
defineVar("buf_stack2")
setVar("buf_stack2", []interface{}{"end"})
defineVar("inside_vars_stack")
setVar("inside_vars_stack", []interface{}{"end"})
defineVar("inside_vars_stack_reserv")
setVar("inside_vars_stack_reserv", []interface{}{"end"})
defineVar("null")
setVar("null", []interface{}{"end"})
defineVar("undefined_args")
setVar("undefined_args", []interface{}{"end"})
defineVar("r")
setVar("r", []interface{}{"end"})
defineVar("el")
setVar("el", []interface{}{"end"})
defineVar("is_first")
setVar("is_first", false)
defineVar("was_func")
setVar("was_func", false)
defineVar("was_here")
setVar("was_here", false)
defineVar("void_flag")
setVar("void_flag", false)
defineVar("in_stack_res")
setVar("in_stack_res", false)
defineVar("undefine_checked")
setVar("undefine_checked", false)
defineVar("root_source")
setVar("root_source", "")
defineVar("source")
setVar("source", "")
defineVar("temp_dir")
setVar("temp_dir", "")
defineVar("ops_ret")
setVar("ops_ret", "")
defineVar("temp_buf")
setVar("temp_buf", "")
defineVar("$func_end_res")
setVar("$func_end_res", "")
defineVar("$ops_ret")
setVar("$ops_ret", "")
defineVar("deb")
setVar("deb", "")
setVar("is_first","True")
setVar("was_here","False")
setVar("undefine_checked","False")
setVar("root_source","benv/recurs_program.b")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", getVar("root_source")))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
setVar("source", getVar("root_source"))
setVar("dest","benv/func_program.b")
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", getVar("dest")))))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto again_begin
again_begin:
setVar("i","0")
setVar("j","0")
setVar("right_border","0")
setVar("command_len","0")
setVar("call_number","0")
setVar("str_call_number","")
setVar("return_type","")
setVar("buf","")
setVar("buf2","")
setVar("buf3","")
setVar("int_buf","0")
setVar("int_buf2","0")
setVar("int_buf3","0")
setVar("type_len","0")
setVar("arg_type_len","0")
setVar("func_begin","0")
setVar("func_end","0")
setVar("func_len","0")
setVar("number","0")
setVar("func_len","0")
setVar("br_number","0")
setVar("func_name","")
setVar("arg_type","")
setVar("arg_name","")
setVar("left_border","0")
setVar("comma_pos","0")
setVar("was_func","False")
setVar("buf_stack", getVar("null"))
setVar("buf_stack_reserv", getVar("null"))
setVar("buf_stack2", getVar("null"))
setVar("inside_vars_stack", getVar("null"))
setVar("inside_vars_stack_reserv", getVar("null"))
setVar("call_number","0")
goto begin
begin:
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto again_end
}else{
print("")
}
setVar("void_flag","False")
setVar("br_number",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v","(")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("br_number"))){
goto command_end
}else{
print("")
}
setVar("buf",getVar("command").(string)[0:toInt(getVar("br_number"))])
setVar("return_type","int")
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","int")))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto end_type
}else{
print("")
}
setVar("return_type","bool")
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","bool")))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto end_type
}else{
print("")
}
setVar("return_type","float")
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","float")))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto end_type
}else{
print("")
}
setVar("return_type","stack")
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","stack")))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto end_type
}else{
print("")
}
setVar("return_type","string")
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","string")))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto end_type
}else{
print("")
}
setVar("return_type","void")
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","void")))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("number"))){
goto end_type
}else{
print("")
}
goto command_end
goto end_type
end_type:
if isEqual(ValueFoldInterface("void"), ValueFoldInterface(getVar("return_type"))){
print("")
}else{
goto set_void_flag_end
}
setVar("void_flag","True")
goto set_void_flag_end
set_void_flag_end:
setVar("type_len",len(fmt.Sprintf("%v",getVar("return_type"))))
setVar("func_name",getVar("command").(string)[toInt(getVar("type_len")):toInt(getVar("br_number"))])
setVar("buf",sum(sum("string $", getVar("func_name")), "_return_var"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
setVar("left_border",sum(getVar("br_number"), 1))
setVar("was_func","True")
setVar("buf",sum(sum(sum(getVar("return_type"), " $"), getVar("func_name")), "_res"))
if toBool(getVar("void_flag")){
goto not_define_void_end
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
goto not_define_void_end
not_define_void_end:
fmt.Print("")
goto func_args
func_args:
if isEqual(ValueFoldInterface("{"), ValueFoldInterface(string(getVar("command").(string)[0]))){
goto func_end
}else{
print("")
}
setVar("arg_type","int")
setVar("right_border",sum(getVar("left_border"), 3))
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("command").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","int")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto arg_type_end
}
setVar("arg_type","bool")
setVar("right_border",sum(getVar("left_border"), 4))
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("command").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","bool")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto arg_type_end
}
setVar("arg_type","float")
setVar("right_border",sum(getVar("left_border"), 5))
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("command").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","float")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto arg_type_end
}
setVar("arg_type","stack")
setVar("right_border",sum(getVar("left_border"), 5))
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("command").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","stack")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto arg_type_end
}
setVar("arg_type","string")
setVar("right_border",sum(getVar("left_border"), 6))
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("command").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","string")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto arg_type_end
}
setVar("arg_type","no_args")
goto func_end
goto arg_type_end
arg_type_end:
setVar("arg_type_len",len(fmt.Sprintf("%v",getVar("arg_type"))))
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command",getVar("command").(string)[toInt(getVar("right_border")):toInt(getVar("command_len"))])
setVar("comma_pos",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v",",")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("comma_pos"))){
print("")
}else{
goto br_end
}
setVar("comma_pos",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v",")")))
goto br_end
br_end:
setVar("arg_name",getVar("command").(string)[0:toInt(getVar("comma_pos"))])
setVar("user_stack", append(getVar("user_stack").([]interface{}), getVar("arg_name")))
setVar("buf_stack", append(getVar("buf_stack").([]interface{}), getVar("arg_name")))
setVar("buf_stack2", append(getVar("buf_stack2").([]interface{}), getVar("arg_name")))
setVar("buf_stack2", append(getVar("buf_stack2").([]interface{}), getVar("arg_type")))
setVar("left_border",sum(getVar("comma_pos"), 1))
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
setVar("command",getVar("command").(string)[toInt(getVar("left_border")):toInt(getVar("command_len"))])
setVar("left_border","0")
setVar("right_border",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v",",")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("right_border"))){
print("")
}else{
goto func_args
}
setVar("right_border","3")
goto func_args
goto func_end
func_end:
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
setVar("buf",getVar("command").(string)[1:toInt(getVar("command_len"))])
setVar("command","goto(#")
setVar("command",sum(getVar("command"), getVar("func_name")))
setVar("command",sum(getVar("command"), "_end)"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("command","#")
setVar("command",sum(getVar("command"), getVar("func_name")))
setVar("command",sum(getVar("command"), ":\n"))
setVar("command",sum(getVar("command"), "print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("command","")
goto args_to_pass
args_to_pass:
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack2")) && len(getVar("buf_stack2").([]interface{})) > 1{
setVar("arg_type", getVar("buf_stack2").([]interface{})[len(getVar("buf_stack2").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_type")) == "[]interface {}"{
setVar("arg_type", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_type", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack2")) && !isEqual("end", getVar("buf_stack2").([]interface{})[len(getVar("buf_stack2").([]interface{})) - 1]) && !isEqual("[end]", getVar("buf_stack2").([]interface{})[len(getVar("buf_stack2").([]interface{})) - 1]){
setVar("buf_stack2", getVar("buf_stack2").([]interface{})[:len(getVar("buf_stack2").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack2")) && len(getVar("buf_stack2").([]interface{})) > 1{
setVar("arg_name", getVar("buf_stack2").([]interface{})[len(getVar("buf_stack2").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack2")) && !isEqual("end", getVar("buf_stack2").([]interface{})[len(getVar("buf_stack2").([]interface{})) - 1]) && !isEqual("[end]", getVar("buf_stack2").([]interface{})[len(getVar("buf_stack2").([]interface{})) - 1]){
setVar("buf_stack2", getVar("buf_stack2").([]interface{})[:len(getVar("buf_stack2").([]interface{})) - 1])
}
setVar("command",sum(sum(getVar("command"), getVar("arg_type")), getVar("arg_name")))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_type"))){
goto not_send_pass_args
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto not_send_pass_args
not_send_pass_args:
setVar("command","")
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_type"))){
print("")
}else{
goto args_to_pass
}
setVar("fig_br_number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","{")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("fig_br_number"))){
print("")
}else{
goto not_send_fig_br
}
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
goto not_send_fig_br_end
goto not_send_fig_br
not_send_fig_br:
setVar("buf","print(\"\")")
goto not_send_fig_br_end
not_send_fig_br_end:
setVar("command","")
setVar("arg_name","")
goto pop_start
pop_start:
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto pop_end
}else{
print("")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("user_stack")) && len(getVar("user_stack").([]interface{})) > 1{
setVar("arg_name", getVar("user_stack").([]interface{})[len(getVar("user_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("user_stack")) && !isEqual("end", getVar("user_stack").([]interface{})[len(getVar("user_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("user_stack").([]interface{})[len(getVar("user_stack").([]interface{})) - 1]){
setVar("user_stack", getVar("user_stack").([]interface{})[:len(getVar("user_stack").([]interface{})) - 1])
}
setVar("command",sum(getVar("command"), "pop("))
setVar("command",sum(getVar("command"), getVar("arg_name")))
setVar("command",sum(getVar("command"), ")"))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto not_send_pop
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto not_send_pop
not_send_pop:
setVar("command","")
goto pop_start
goto pop_end
pop_end:
if toBool(getVar("undefine_checked")){
goto check_undefine_end0
}else{
print("")
}
setVar("undefine_checked","True")
goto check_begin
check_begin:
if isEqual(ValueFoldInterface("}"), ValueFoldInterface(getVar("buf"))){
goto check_again
}else{
print("")
}
setVar("command_len",len(fmt.Sprintf("%v",getVar("buf"))))
if toFloat(getVar("command_len"))>8{
print("")
}else{
goto check_next
}
setVar("temp_buf",getVar("buf").(string)[0:8])
if isEqual(ValueFoldInterface("UNDEFINE"), ValueFoldInterface(getVar("temp_buf"))){
print("")
}else{
goto check_next
}
setVar("command_len",toFloat(getVar("command_len"))-1)
setVar("temp_buf",getVar("buf").(string)[9:toInt(getVar("command_len"))])
setVar("undefined_args", append(getVar("undefined_args").([]interface{}), getVar("temp_buf")))
goto check_next
check_next:
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("buf", getVar("$CODE"))
goto check_begin
goto check_again
check_again:
setVar("was_func","False")
getVar("$SOURCE").(*os.File).Seek(0, 0)
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
getVar("$DEST").(*os.File).Close()
os.Remove(getRootSource(fmt.Sprintf("%v", getVar("dest"))))
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", getVar("dest")))))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto again_begin
goto check_undefine_end0
check_undefine_end0:
setVar("undefine_checked","False")
setVar("left_border","0")
goto inside_args
inside_args:
setVar("limit_border",len(fmt.Sprintf("%v",getVar("buf"))))
setVar("limit_border",toFloat(getVar("limit_border"))-1)
if isEqual(ValueFoldInterface(getVar("buf")), ValueFoldInterface("}")){
goto inside_args_end
}else{
print("")
}
setVar("arg_type","int")
setVar("right_border",sum(getVar("left_border"), 3))
if toFloat(getVar("right_border"))>toFloat(getVar("limit_border")){
goto no_args
}else{
print("")
}
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","int")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto inside_arg_type_end
}
setVar("arg_type","bool")
setVar("right_border",sum(getVar("left_border"), 4))
if toFloat(getVar("right_border"))>toFloat(getVar("limit_border")){
goto no_args
}else{
print("")
}
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","bool")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto inside_arg_type_end
}
setVar("arg_type","float")
setVar("right_border",sum(getVar("left_border"), 5))
if toFloat(getVar("right_border"))>toFloat(getVar("limit_border")){
goto no_args
}else{
print("")
}
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","float")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto inside_arg_type_end
}
setVar("arg_type","stack")
setVar("right_border",sum(getVar("left_border"), 5))
if toFloat(getVar("right_border"))>toFloat(getVar("limit_border")){
goto no_args
}else{
print("")
}
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","stack")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto inside_arg_type_end
}
setVar("arg_type","string")
setVar("right_border",sum(getVar("left_border"), 6))
if toFloat(getVar("right_border"))>toFloat(getVar("limit_border")){
goto no_args
}else{
print("")
}
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("buf").(string)[toInt(getVar("left_border")):toInt(getVar("right_border"))]), fmt.Sprintf("%v","string")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
print("")
}else{
goto inside_arg_type_end
}
setVar("arg_type","no_args")
goto inside_arg_type_end
inside_arg_type_end:
if isEqual(ValueFoldInterface("no_args"), ValueFoldInterface(getVar("arg_type"))){
goto no_args
}else{
print("")
}
setVar("type_len",len(fmt.Sprintf("%v",getVar("arg_type"))))
setVar("command_len",len(fmt.Sprintf("%v",getVar("buf"))))
setVar("arg_name",getVar("buf").(string)[toInt(getVar("type_len")):toInt(getVar("command_len"))])
setVar("$in_stack_ret","#undef_ret")
setVar("user_stack", append(getVar("user_stack").([]interface{}), getVar("arg_name")))
systemStack = append(systemStack, getVar("arg_name"))
systemStack = append(systemStack, getVar("undefined_args"))
goto in_stack
goto undef_ret
undef_ret:
setVar("in_stack_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
if toBool(getVar("in_stack_res")){
goto no_args
}else{
print("")
}
setVar("inside_vars_stack", append(getVar("inside_vars_stack").([]interface{}), getVar("arg_name")))
goto no_args
no_args:
setVar("command_len",len(fmt.Sprintf("%v",getVar("buf"))))
if toFloat(getVar("command_len"))>6{
print("")
}else{
goto not_send_return_end
}
defineVar("$regRes")
{intListList := compileRegexp(`[^=]=[^=]`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("buf"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("r",getVar("$regRes"))
if "[]interface {}" == fmt.Sprintf("%T", getVar("r")) && len(getVar("r").([]interface{})) > 1{
setVar("el", getVar("r").([]interface{})[len(getVar("r").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("el")) == "[]interface {}"{
setVar("el", []interface{}{[]interface{}{"end"}})
}else{
setVar("el", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("r")) && !isEqual("end", getVar("r").([]interface{})[len(getVar("r").([]interface{})) - 1]) && !isEqual("[end]", getVar("r").([]interface{})[len(getVar("r").([]interface{})) - 1]){
setVar("r", getVar("r").([]interface{})[:len(getVar("r").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("el")) && len(getVar("el").([]interface{})) > 1{
setVar("rbuf", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("rbuf")) == "[]interface {}"{
setVar("rbuf", []interface{}{[]interface{}{"end"}})
}else{
setVar("rbuf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("el")) && !isEqual("end", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1]) && !isEqual("[end]", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1]){
setVar("el", getVar("el").([]interface{})[:len(getVar("el").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("rbuf"))){
print("")
}else{
goto not_send_return_end
}
defineVar("$regRes")
{intListList := compileRegexp(`^return.*`).FindAllIndex([]byte(fmt.Sprintf("%v", getVar("buf"))), -1)
var res []interface{}
res = append(res, []interface{}{"end"})
for i := len(intListList) - 1; i >= 0; i-- {
res = append(res, []interface{}{[]interface{}{"end"}})
for j := len(intListList[i]) - 1; j >= 0; j-- {
res[len(res)-1] = append(res[len(res)-1].([]interface{}), intListList[i][j])
}
}
setVar("$regRes", res)}
setVar("r",getVar("$regRes"))
if "[]interface {}" == fmt.Sprintf("%T", getVar("r")) && len(getVar("r").([]interface{})) > 1{
setVar("el", getVar("r").([]interface{})[len(getVar("r").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("el")) == "[]interface {}"{
setVar("el", []interface{}{[]interface{}{"end"}})
}else{
setVar("el", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("r")) && !isEqual("end", getVar("r").([]interface{})[len(getVar("r").([]interface{})) - 1]) && !isEqual("[end]", getVar("r").([]interface{})[len(getVar("r").([]interface{})) - 1]){
setVar("r", getVar("r").([]interface{})[:len(getVar("r").([]interface{})) - 1])
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("el")) && len(getVar("el").([]interface{})) > 1{
setVar("rbuf", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("rbuf")) == "[]interface {}"{
setVar("rbuf", []interface{}{[]interface{}{"end"}})
}else{
setVar("rbuf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("el")) && !isEqual("end", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1]) && !isEqual("[end]", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1]){
setVar("el", getVar("el").([]interface{})[:len(getVar("el").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("rbuf"))){
goto not_send_return_end
}else{
print("")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("el")) && len(getVar("el").([]interface{})) > 1{
setVar("rbuf", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("rbuf")) == "[]interface {}"{
setVar("rbuf", []interface{}{[]interface{}{"end"}})
}else{
setVar("rbuf", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("el")) && !isEqual("end", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1]) && !isEqual("[end]", getVar("el").([]interface{})[len(getVar("el").([]interface{})) - 1]){
setVar("el", getVar("el").([]interface{})[:len(getVar("el").([]interface{})) - 1])
}
setVar("int_buf4",toInt(getVar("rbuf")))
setVar("int_buf5",len(fmt.Sprintf("%v",getVar("buf"))))
if isEqual(ValueFoldInterface(getVar("int_buf4")), ValueFoldInterface(getVar("int_buf5"))){
print("")
}else{
goto not_send_return_end
}
if isEqual(ValueFoldInterface("return"), ValueFoldInterface(getVar("buf").(string)[0:6])){
goto not_send_return
}else{
goto not_send_return_end
}
goto not_send_return
not_send_return:
setVar("buf",getVar("buf").(string)[6:toInt(getVar("command_len"))])
setVar("buf",sum(sum("push(", getVar("buf")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("buf", getVar("$CODE"))
goto pop_args_start0
pop_args_start0:
if "[]interface {}" == fmt.Sprintf("%T", getVar("inside_vars_stack")) && len(getVar("inside_vars_stack").([]interface{})) > 1{
setVar("arg_name", getVar("inside_vars_stack").([]interface{})[len(getVar("inside_vars_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("inside_vars_stack")) && !isEqual("end", getVar("inside_vars_stack").([]interface{})[len(getVar("inside_vars_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("inside_vars_stack").([]interface{})[len(getVar("inside_vars_stack").([]interface{})) - 1]){
setVar("inside_vars_stack", getVar("inside_vars_stack").([]interface{})[:len(getVar("inside_vars_stack").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto pop_args_end0
}else{
print("")
}
setVar("command",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto not_send0
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("inside_vars_stack_reserv", append(getVar("inside_vars_stack_reserv").([]interface{}), getVar("arg_name")))
goto not_send0
not_send0:
goto pop_args_start0
goto pop_args_end0
pop_args_end0:
setVar("arg_name","")
goto undefine_start0
undefine_start0:
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto undefine_end0
}else{
print("")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack")) && len(getVar("buf_stack").([]interface{})) > 1{
setVar("arg_name", getVar("buf_stack").([]interface{})[len(getVar("buf_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack")) && !isEqual("end", getVar("buf_stack").([]interface{})[len(getVar("buf_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("buf_stack").([]interface{})[len(getVar("buf_stack").([]interface{})) - 1]){
setVar("buf_stack", getVar("buf_stack").([]interface{})[:len(getVar("buf_stack").([]interface{})) - 1])
}
setVar("command",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto undefine_end0
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("buf_stack_reserv", append(getVar("buf_stack_reserv").([]interface{}), getVar("arg_name")))
goto undefine_start0
goto undefine_end0
undefine_end0:
setVar("inside_vars_stack", getVar("inside_vars_stack_reserv"))
setVar("buf_stack", getVar("buf_stack_reserv"))
setVar("inside_vars_stack_reserv", getVar("null"))
setVar("buf_stack_reserv", getVar("null"))
setVar("arg_name","")
setVar("command",sum(sum("goto($", getVar("func_name")), "_return_var)"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto inside_args
goto not_send_return_end
not_send_return_end:
getVar("$DEST").(*os.File).WriteString(getVar("buf").(string) + ";\n")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("buf", getVar("$CODE"))
goto inside_args
goto inside_args_end
inside_args_end:
if isEqual(ValueFoldInterface("void"), ValueFoldInterface(getVar("return_type"))){
print("")
}else{
goto define_return_end
}
goto pop_args_start
pop_args_start:
if "[]interface {}" == fmt.Sprintf("%T", getVar("inside_vars_stack")) && len(getVar("inside_vars_stack").([]interface{})) > 1{
setVar("arg_name", getVar("inside_vars_stack").([]interface{})[len(getVar("inside_vars_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("inside_vars_stack")) && !isEqual("end", getVar("inside_vars_stack").([]interface{})[len(getVar("inside_vars_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("inside_vars_stack").([]interface{})[len(getVar("inside_vars_stack").([]interface{})) - 1]){
setVar("inside_vars_stack", getVar("inside_vars_stack").([]interface{})[:len(getVar("inside_vars_stack").([]interface{})) - 1])
}
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto pop_args_end
}else{
print("")
}
setVar("command",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto not_send
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("inside_vars_stack_reserv", append(getVar("inside_vars_stack_reserv").([]interface{}), getVar("arg_name")))
goto not_send
not_send:
goto pop_args_start
goto pop_args_end
pop_args_end:
setVar("arg_name","")
goto undefine_start
undefine_start:
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto undefine_end
}else{
print("")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack")) && len(getVar("buf_stack").([]interface{})) > 1{
setVar("arg_name", getVar("buf_stack").([]interface{})[len(getVar("buf_stack").([]interface{})) - 1])
}else if fmt.Sprintf("%T", getVar("arg_name")) == "[]interface {}"{
setVar("arg_name", []interface{}{[]interface{}{"end"}})
}else{
setVar("arg_name", "end")
}
if "[]interface {}" == fmt.Sprintf("%T", getVar("buf_stack")) && !isEqual("end", getVar("buf_stack").([]interface{})[len(getVar("buf_stack").([]interface{})) - 1]) && !isEqual("[end]", getVar("buf_stack").([]interface{})[len(getVar("buf_stack").([]interface{})) - 1]){
setVar("buf_stack", getVar("buf_stack").([]interface{})[:len(getVar("buf_stack").([]interface{})) - 1])
}
setVar("command",sum(sum("UNDEFINE(", getVar("arg_name")), ")"))
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("arg_name"))){
goto undefine_end
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("buf_stack_reserv", append(getVar("buf_stack_reserv").([]interface{}), getVar("arg_name")))
goto undefine_start
goto undefine_end
undefine_end:
setVar("inside_vars_stack", getVar("inside_vars_stack_reserv"))
setVar("buf_stack", getVar("buf_stack_reserv"))
setVar("arg_name","")
setVar("command",sum(sum("goto($", getVar("func_name")), "_return_var)"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto define_return_end
define_return_end:
fmt.Print("")
setVar("inside_vars_stack", getVar("null"))
setVar("buf_stack", getVar("null"))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("command",sum(sum(sum("#", sum(getVar("func_name"), "_end:")), "\n"), getVar("command")))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
setVar("i","0")
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
setVar("func_name_len",len(fmt.Sprintf("%v",getVar("func_name"))))
setVar("func_name_br",sum(getVar("func_name"), "("))
goto start_entry
start_entry:
if isEqual(ValueFoldInterface("end"), ValueFoldInterface(getVar("command"))){
goto again_end
}else{
print("")
}
setVar("number",strings.Index(fmt.Sprintf("%v",getVar("command")), fmt.Sprintf("%v",getVar("func_name_br"))))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("number"))){
goto mod_end
}else{
print("")
}
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("func_name"))
goto func_call_index
goto func_call_index_ret
func_call_index_ret:
setVar("func_pos", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("func_pos"))){
goto mod_end
}else{
print("")
}
setVar("number",sum(getVar("number"), getVar("func_name_len")))
setVar("buf4",string(getVar("command").(string)[toInt(getVar("number"))]))
setVar("number",toFloat(getVar("number"))-toFloat(getVar("func_name_len")))
if isEqual(ValueFoldInterface("("), ValueFoldInterface(getVar("buf4"))){
print("")
}else{
goto mod_end
}
setVar("buf", getVar("command"))
setVar("br_number",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v","(")))
setVar("left_border",sum(getVar("br_number"), 1))
setVar("buf2", getVar("buf"))
goto call_start
call_start:
if isEqual(ValueFoldInterface(""), ValueFoldInterface(getVar("buf"))){
goto call_end
}else{
print("")
}
setVar("comma_pos",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v",",")))
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("comma_pos"))){
print("")
}else{
goto comma_end
}
setVar("buf_len",len(fmt.Sprintf("%v",getVar("buf"))))
setVar("comma_pos",toFloat(getVar("buf_len"))-1)
goto comma_end
comma_end:
if isEqual(ValueFoldInterface(getVar("left_border")), ValueFoldInterface(getVar("comma_pos"))){
print("")
}else{
goto aend
}
setVar("arg_name","")
goto askip
goto aend
aend:
setVar("arg_name",getVar("buf").(string)[toInt(getVar("left_border")):toInt(getVar("comma_pos"))])
goto askip
askip:
setVar("command",sum(sum("push(", getVar("arg_name")), ")"))
if isEqual(ValueFoldInterface(""), ValueFoldInterface(getVar("arg_name"))){
goto not_send_push
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto not_send_push
not_send_push:
setVar("comma_pos",sum(getVar("comma_pos"), 1))
setVar("command_len",len(fmt.Sprintf("%v",getVar("buf"))))
if isEqual(ValueFoldInterface(getVar("comma_pos")), ValueFoldInterface(getVar("command_len"))){
print("")
}else{
goto buf_mod_end
}
setVar("buf","")
setVar("left_border","0")
goto call_start
goto buf_mod_end
buf_mod_end:
setVar("buf",getVar("buf").(string)[toInt(getVar("comma_pos")):toInt(getVar("command_len"))])
setVar("left_border","0")
goto call_start
goto call_end
call_end:
setVar("buf", getVar("buf2"))
setVar("buf3", getVar("buf"))
setVar("func_begin",strings.Index(fmt.Sprintf("%v",getVar("buf")), fmt.Sprintf("%v",getVar("func_name_br"))))
setVar("command_len",len(fmt.Sprintf("%v",getVar("buf"))))
setVar("internal_func_begin",sum(getVar("func_begin"), getVar("func_name_len")))
systemStack = append(systemStack, getVar("internal_func_begin"))
systemStack = append(systemStack, getVar("buf"))
setVar("$func_end_res","#gl_fe")
goto func_end_func
goto gl_fe
gl_fe:
setVar("func_end", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_end",sum(getVar("func_end"), 1))
setVar("func_len", getVar("func_end"))
setVar("buf",getVar("buf").(string)[0:toInt(getVar("func_begin"))])
setVar("buf",sum(sum(sum(getVar("buf"), "$"), getVar("func_name")), "_res"))
if isEqual(ValueFoldInterface(getVar("func_end")), ValueFoldInterface(getVar("command_len"))){
goto nothing
}else{
print("")
}
setVar("buf2",getVar("buf3").(string)[toInt(getVar("func_end")):toInt(getVar("command_len"))])
setVar("buf",sum(getVar("buf"), getVar("buf2")))
goto nothing
nothing:
if toBool(getVar("void_flag")){
print("")
}else{
goto set_buf_nop_end
}
setVar("buf","print(\"\")")
goto set_buf_nop_end
set_buf_nop_end:
setVar("str_call_number",getVar("call_number"))
setVar("command",sum(sum(sum(sum(sum(sum("$", getVar("func_name")), "_return_var=\"#"), getVar("func_name")), "_res"), getVar("str_call_number")), "\""))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("command",sum(sum("goto(#", getVar("func_name")), ")"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
setVar("command",sum(sum(sum(sum("#", getVar("func_name")), "_res"), getVar("str_call_number")), ":\n"))
setVar("command",sum(getVar("command"), sum(sum("pop($", getVar("func_name")), "_res)")))
if toBool(getVar("void_flag")){
goto not_send_void_pop_end
}else{
print("")
}
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto not_send_void_pop_end
not_send_void_pop_end:
if toBool(getVar("void_flag")){
print("")
}else{
goto send_void_mark_end
}
setVar("command",sum(sum(sum(sum("#", getVar("func_name")), "_res"), getVar("str_call_number")), ":\n"))
setVar("command",sum(getVar("command"), "print(\"\")"))
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
goto send_void_mark_end
send_void_mark_end:
setVar("command", getVar("buf"))
setVar("call_number",sum(getVar("call_number"), 1))
goto mod_end
mod_end:
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto start_entry
goto end_entry
end_entry:
fmt.Print("")
goto command_end
command_end:
getVar("$DEST").(*os.File).WriteString(getVar("command").(string) + ";\n")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto begin
goto again_end
again_end:
if toBool(getVar("is_first")){
print("")
}else{
goto not_is_first
}
getVar("$SOURCE").(*os.File).Close()
getVar("$DEST").(*os.File).Close()
setVar("source","benv/func_program.b")
defineVar("$SOURCE")
setVar("$SOURCE", openFile(getRootSource(fmt.Sprintf("%v", "benv/func_program.b"))))
defineVar("$sourceNewChunk")
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
setVar("dest","benv/func_program2.b")
defineVar("$DEST")
setVar("$DEST", createFile(getRootSource(fmt.Sprintf("%v", "benv/func_program2.b"))))
setVar("is_first","False")
setVar("was_func","False")
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto again_begin
goto not_is_first
not_is_first:
if toBool(getVar("was_func")){
print("")
}else{
goto del_file
}
setVar("was_func","False")
getVar("$DEST").(*os.File).Close()
getVar("$SOURCE").(*os.File).Close()
defineVar("$temp")
setVar("$temp", getVar("$DEST"))
setVar("$DEST", getVar("$SOURCE"))
setVar("$SOURCE", getVar("$temp"))
setVar("$SOURCE", openFile(getVar("$SOURCE").(*os.File).Name()))
setVar("$DEST", openFile666(getVar("$DEST").(*os.File).Name()))
setVar("$sourceNewChunk", EachChunk(getVar("$SOURCE").(*os.File)))
setVar("temp_dir", getVar("source"))
setVar("source", getVar("dest"))
setVar("dest", getVar("temp_dir"))
defineVar("$CODE")
setVar("$CODE", CodeInput(fmt.Sprintf("%v", getVar("$sourceNewChunk").(func () string)()), false))
setVar("command", getVar("$CODE"))
goto again_begin
goto del_file
del_file:
os.Remove(getRootSource("benv/func_program2.b"))
os.Remove(getRootSource("benv/recurs_program.b"))
goto end
end:
getVar("$SOURCE").(*os.File).Close()
getVar("$DEST").(*os.File).Close()
os.Exit(0)
goto reverse
reverse:
defineVar("buf")
setVar("buf", "")
defineVar("res")
setVar("res", []interface{}{"end"})
defineVar("s")
setVar("s", []interface{}{"end"})
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
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
systemStack = append(systemStack, getVar("res"))
undefineVar("buf")
undefineVar("res")
undefineVar("s")
if "#again_begin" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto again_begin
}
if "#begin" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto begin
}
if "#end_type" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto end_type
}
if "#set_void_flag_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto set_void_flag_end
}
if "#not_define_void_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_define_void_end
}
if "#func_args" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto func_args
}
if "#arg_type_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto arg_type_end
}
if "#br_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto br_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto func_end
}
if "#args_to_pass" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto args_to_pass
}
if "#not_send_pass_args" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_pass_args
}
if "#not_send_fig_br" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_fig_br
}
if "#not_send_fig_br_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_fig_br_end
}
if "#pop_start" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto pop_start
}
if "#not_send_pop" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_pop
}
if "#pop_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto pop_end
}
if "#check_begin" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto check_begin
}
if "#check_next" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto check_next
}
if "#check_again" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto check_again
}
if "#check_undefine_end0" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto check_undefine_end0
}
if "#inside_args" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto inside_args
}
if "#inside_arg_type_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto inside_arg_type_end
}
if "#undef_ret" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto undef_ret
}
if "#no_args" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto no_args
}
if "#not_send_return" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_return
}
if "#pop_args_start0" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto pop_args_start0
}
if "#not_send0" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send0
}
if "#pop_args_end0" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto pop_args_end0
}
if "#undefine_start0" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto undefine_start0
}
if "#undefine_end0" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto undefine_end0
}
if "#not_send_return_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_return_end
}
if "#inside_args_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto inside_args_end
}
if "#pop_args_start" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto pop_args_start
}
if "#not_send" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send
}
if "#pop_args_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto pop_args_end
}
if "#undefine_start" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto undefine_start
}
if "#undefine_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto undefine_end
}
if "#define_return_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto define_return_end
}
if "#start_entry" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto start_entry
}
if "#func_call_index_ret" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto func_call_index_ret
}
if "#call_start" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto call_start
}
if "#comma_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto comma_end
}
if "#aend" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto aend
}
if "#askip" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto askip
}
if "#not_send_push" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_push
}
if "#buf_mod_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto buf_mod_end
}
if "#call_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto call_end
}
if "#gl_fe" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto gl_fe
}
if "#nothing" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto nothing
}
if "#set_buf_nop_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto set_buf_nop_end
}
if "#not_send_void_pop_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_send_void_pop_end
}
if "#send_void_mark_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto send_void_mark_end
}
if "#mod_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto mod_end
}
if "#end_entry" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto end_entry
}
if "#command_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto command_end
}
if "#again_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto again_end
}
if "#not_is_first" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_is_first
}
if "#del_file" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto del_file
}
if "#end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto reverse_e
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto look_behind_e
}
if "#look_ret" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto look_ret
}
if "#indexes" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto indexes_s
}
if "#iend" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto iend
}
if "#iskip" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto iskip
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto indexes_e
}
if "#iret" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto iret
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto func_call_index
}
if "#look_behind_ret" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto look_behind_ret
}
if "#indexes_ret" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto indexes_ret
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto func_call_index_s
}
if "#is_found" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto is_found
}
if "#not_func" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto not_func
}
if "#end_func" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto end_func
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _in_stack_e
}
if "#ops" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto ops
}
if "#_idx" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _idx
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _op_nums_e
}
if "#_rvs" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _rvs
}
if "#func_end_func" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto func_end_func
}
if "#_fe" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _fe
}
if "#_fe2" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _fe2
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _braces_s
}
if "#_fe3" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _fe3
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _o_plus_end
}
if "#_fe4" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _fe4
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$reverse_ret")){
goto _braces_e
}
goto look_behind
look_behind:
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
defineVar("reg")
setVar("reg", "")
defineVar("s")
setVar("s", "")
setVar("reg", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
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
setVar("st",getVar("$regRes"))
goto look_behind_s
look_behind_s:
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
setVar("pos",toInt(getVar("buf")))
if isEqual(ValueFoldInterface(0), ValueFoldInterface(getVar("pos"))){
print("")
}else{
goto is_not_zero
}
setVar("res", append(getVar("res").([]interface{}), "$"))
goto look_behind_s
goto is_not_zero
is_not_zero:
setVar("pos",toFloat(getVar("pos"))-1)
setVar("symbol",string(getVar("s").(string)[toInt(getVar("pos"))]))
setVar("res", append(getVar("res").([]interface{}), getVar("symbol")))
goto look_behind_s
goto look_behind_e
look_behind_e:
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_ret","#look_ret")
goto look_ret
look_ret:
setVar("res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
systemStack = append(systemStack, getVar("res"))
undefineVar("st")
undefineVar("this")
undefineVar("res")
undefineVar("buf")
undefineVar("pos")
undefineVar("symbol")
undefineVar("reg")
undefineVar("s")
goto look_behind_ret
goto indexes
indexes:
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
defineVar("s")
setVar("s", "")
defineVar("sub_s")
setVar("sub_s", "")
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("sub_s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("s_len",len(fmt.Sprintf("%v",getVar("s"))))
setVar("sub_len",len(fmt.Sprintf("%v",getVar("sub_s"))))
setVar("s_len_old", getVar("s_len"))
setVar("i",strings.Index(fmt.Sprintf("%v",getVar("s")), fmt.Sprintf("%v",getVar("sub_s"))))
setVar("pointer", getVar("i"))
goto indexes_s
indexes_s:
if isEqual(ValueFoldInterface(-1), ValueFoldInterface(getVar("i"))){
goto indexes_e
}else{
print("")
}
setVar("i",sum(getVar("i"), toFloat(getVar("s_len_old"))-toFloat(getVar("s_len"))))
setVar("res", append(getVar("res").([]interface{}), getVar("i")))
setVar("pointer",sum(getVar("pointer"), getVar("sub_len")))
if isEqual(ValueFoldInterface(getVar("pointer")), ValueFoldInterface(getVar("s_len"))){
print("")
}else{
goto iend
}
setVar("s","")
goto iskip
goto iend
iend:
setVar("s",getVar("s").(string)[toInt(getVar("pointer")):toInt(getVar("s_len"))])
goto iskip
iskip:
setVar("s_len",len(fmt.Sprintf("%v",getVar("s"))))
setVar("i",strings.Index(fmt.Sprintf("%v",getVar("s")), fmt.Sprintf("%v",getVar("sub_s"))))
setVar("pointer", getVar("i"))
goto indexes_s
goto indexes_e
indexes_e:
systemStack = append(systemStack, getVar("res"))
setVar("$reverse_ret","#iret")
goto iret
iret:
setVar("res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
systemStack = append(systemStack, getVar("res"))
undefineVar("res")
undefineVar("i")
undefineVar("pointer")
undefineVar("s_len")
undefineVar("s_len_old")
undefineVar("sub_len")
undefineVar("s")
undefineVar("sub_s")
if "#again_begin" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto again_begin
}
if "#begin" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto begin
}
if "#end_type" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto end_type
}
if "#set_void_flag_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto set_void_flag_end
}
if "#not_define_void_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_define_void_end
}
if "#func_args" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto func_args
}
if "#arg_type_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto arg_type_end
}
if "#br_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto br_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto func_end
}
if "#args_to_pass" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto args_to_pass
}
if "#not_send_pass_args" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_pass_args
}
if "#not_send_fig_br" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_fig_br
}
if "#not_send_fig_br_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_fig_br_end
}
if "#pop_start" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto pop_start
}
if "#not_send_pop" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_pop
}
if "#pop_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto pop_end
}
if "#check_begin" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto check_begin
}
if "#check_next" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto check_next
}
if "#check_again" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto check_again
}
if "#check_undefine_end0" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto check_undefine_end0
}
if "#inside_args" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto inside_args
}
if "#inside_arg_type_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto inside_arg_type_end
}
if "#undef_ret" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto undef_ret
}
if "#no_args" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto no_args
}
if "#not_send_return" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_return
}
if "#pop_args_start0" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto pop_args_start0
}
if "#not_send0" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send0
}
if "#pop_args_end0" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto pop_args_end0
}
if "#undefine_start0" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto undefine_start0
}
if "#undefine_end0" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto undefine_end0
}
if "#not_send_return_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_return_end
}
if "#inside_args_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto inside_args_end
}
if "#pop_args_start" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto pop_args_start
}
if "#not_send" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send
}
if "#pop_args_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto pop_args_end
}
if "#undefine_start" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto undefine_start
}
if "#undefine_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto undefine_end
}
if "#define_return_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto define_return_end
}
if "#start_entry" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto start_entry
}
if "#func_call_index_ret" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto func_call_index_ret
}
if "#call_start" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto call_start
}
if "#comma_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto comma_end
}
if "#aend" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto aend
}
if "#askip" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto askip
}
if "#not_send_push" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_push
}
if "#buf_mod_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto buf_mod_end
}
if "#call_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto call_end
}
if "#gl_fe" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto gl_fe
}
if "#nothing" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto nothing
}
if "#set_buf_nop_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto set_buf_nop_end
}
if "#not_send_void_pop_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_send_void_pop_end
}
if "#send_void_mark_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto send_void_mark_end
}
if "#mod_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto mod_end
}
if "#end_entry" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto end_entry
}
if "#command_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto command_end
}
if "#again_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto again_end
}
if "#not_is_first" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_is_first
}
if "#del_file" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto del_file
}
if "#end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto end
}
if "#reverse" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto reverse_e
}
if "#look_behind" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto look_behind_e
}
if "#look_ret" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto look_ret
}
if "#indexes" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto indexes_s
}
if "#iend" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto iend
}
if "#iskip" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto iskip
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto indexes_e
}
if "#iret" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto iret
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto func_call_index
}
if "#look_behind_ret" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto look_behind_ret
}
if "#indexes_ret" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto indexes_ret
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto func_call_index_s
}
if "#is_found" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto is_found
}
if "#not_func" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto not_func
}
if "#end_func" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto end_func
}
if "#in_stack" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _in_stack_e
}
if "#ops" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto ops
}
if "#_idx" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _idx
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _op_nums_e
}
if "#_rvs" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _rvs
}
if "#func_end_func" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto func_end_func
}
if "#_fe" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _fe
}
if "#_fe2" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _fe2
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _braces_s
}
if "#_fe3" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _fe3
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _o_plus_end
}
if "#_fe4" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _fe4
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("indexes_ret")){
goto _braces_e
}
goto func_call_index
func_call_index:
defineVar("reg")
setVar("reg", "")
defineVar("st")
setVar("st", []interface{}{"end"})
defineVar("ist")
setVar("ist", []interface{}{"end"})
defineVar("res")
setVar("res", 0)
defineVar("buf")
setVar("buf", "")
defineVar("symbol")
setVar("symbol", "")
defineVar("letter")
setVar("letter", false)
defineVar("digit")
setVar("digit", false)
defineVar("command")
setVar("command", "")
defineVar("func_name")
setVar("func_name", "")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_name", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("reg",sum(sum("(?:", getVar("func_name")), ")"))
systemStack = append(systemStack, getVar("reg"))
systemStack = append(systemStack, getVar("command"))
goto look_behind
goto look_behind_ret
look_behind_ret:
setVar("st", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
systemStack = append(systemStack, getVar("command"))
systemStack = append(systemStack, getVar("func_name"))
setVar("indexes_ret","#indexes_ret")
goto indexes
goto indexes_ret
indexes_ret:
setVar("ist", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
goto func_call_index_s
func_call_index_s:
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
goto is_found
}
undefineVar("reg")
undefineVar("st")
undefineVar("ist")
undefineVar("res")
undefineVar("buf")
undefineVar("symbol")
undefineVar("letter")
undefineVar("digit")
undefineVar("command")
undefineVar("func_name")
systemStack = append(systemStack, -1)
goto func_call_index_ret
goto is_found
is_found:
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
print("func_call_index ERROR\n")
}else{
print("")
}
setVar("res",toInt(getVar("buf")))
setVar("letter",unicode.IsLetter([]rune(fmt.Sprintf("%v",getVar("symbol")))[0]))
setVar("digit",unicode.IsDigit([]rune(fmt.Sprintf("%v",getVar("symbol")))[0]))
if !toBool(toBool(toBool(toBool(getVar("letter"))||toBool(getVar("digit")))||toBool(isEqual(ValueFoldInterface("_"), ValueFoldInterface(getVar("symbol")))))||toBool(isEqual(ValueFoldInterface("$"), ValueFoldInterface(getVar("symbol"))))){
print("")
}else{
goto not_func
}
systemStack = append(systemStack, getVar("res"))
undefineVar("reg")
undefineVar("st")
undefineVar("ist")
undefineVar("res")
undefineVar("buf")
undefineVar("symbol")
undefineVar("letter")
undefineVar("digit")
undefineVar("command")
undefineVar("func_name")
goto func_call_index_ret
goto not_func
not_func:
goto func_call_index_s
goto end_func
end_func:
fmt.Print("")
goto in_stack
in_stack:
defineVar("s")
setVar("s", []interface{}{"end"})
defineVar("el")
setVar("el", "")
defineVar("buf")
setVar("buf", "")
defineVar("res")
setVar("res", false)
setVar("s", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("el", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
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
systemStack = append(systemStack, getVar("res"))
undefineVar("s")
undefineVar("el")
undefineVar("buf")
undefineVar("res")
if "#again_begin" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto again_begin
}
if "#begin" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto begin
}
if "#end_type" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto end_type
}
if "#set_void_flag_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto set_void_flag_end
}
if "#not_define_void_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_define_void_end
}
if "#func_args" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto func_args
}
if "#arg_type_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto arg_type_end
}
if "#br_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto br_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto func_end
}
if "#args_to_pass" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto args_to_pass
}
if "#not_send_pass_args" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_pass_args
}
if "#not_send_fig_br" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_fig_br
}
if "#not_send_fig_br_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_fig_br_end
}
if "#pop_start" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto pop_start
}
if "#not_send_pop" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_pop
}
if "#pop_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto pop_end
}
if "#check_begin" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto check_begin
}
if "#check_next" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto check_next
}
if "#check_again" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto check_again
}
if "#check_undefine_end0" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto check_undefine_end0
}
if "#inside_args" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto inside_args
}
if "#inside_arg_type_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto inside_arg_type_end
}
if "#undef_ret" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto undef_ret
}
if "#no_args" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto no_args
}
if "#not_send_return" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_return
}
if "#pop_args_start0" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto pop_args_start0
}
if "#not_send0" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send0
}
if "#pop_args_end0" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto pop_args_end0
}
if "#undefine_start0" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto undefine_start0
}
if "#undefine_end0" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto undefine_end0
}
if "#not_send_return_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_return_end
}
if "#inside_args_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto inside_args_end
}
if "#pop_args_start" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto pop_args_start
}
if "#not_send" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send
}
if "#pop_args_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto pop_args_end
}
if "#undefine_start" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto undefine_start
}
if "#undefine_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto undefine_end
}
if "#define_return_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto define_return_end
}
if "#start_entry" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto start_entry
}
if "#func_call_index_ret" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto func_call_index_ret
}
if "#call_start" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto call_start
}
if "#comma_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto comma_end
}
if "#aend" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto aend
}
if "#askip" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto askip
}
if "#not_send_push" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_push
}
if "#buf_mod_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto buf_mod_end
}
if "#call_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto call_end
}
if "#gl_fe" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto gl_fe
}
if "#nothing" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto nothing
}
if "#set_buf_nop_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto set_buf_nop_end
}
if "#not_send_void_pop_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_send_void_pop_end
}
if "#send_void_mark_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto send_void_mark_end
}
if "#mod_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto mod_end
}
if "#end_entry" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto end_entry
}
if "#command_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto command_end
}
if "#again_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto again_end
}
if "#not_is_first" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_is_first
}
if "#del_file" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto del_file
}
if "#end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto reverse_e
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto look_behind_e
}
if "#look_ret" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto look_ret
}
if "#indexes" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto indexes_s
}
if "#iend" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto iend
}
if "#iskip" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto iskip
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto indexes_e
}
if "#iret" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto iret
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto func_call_index
}
if "#look_behind_ret" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto look_behind_ret
}
if "#indexes_ret" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto indexes_ret
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto func_call_index_s
}
if "#is_found" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto is_found
}
if "#not_func" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto not_func
}
if "#end_func" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto end_func
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _in_stack_e
}
if "#ops" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto ops
}
if "#_idx" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _idx
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _op_nums_e
}
if "#_rvs" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _rvs
}
if "#func_end_func" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto func_end_func
}
if "#_fe" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _fe
}
if "#_fe2" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _fe2
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _braces_s
}
if "#_fe3" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _fe3
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _o_plus_end
}
if "#_fe4" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _fe4
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$in_stack_ret")){
goto _braces_e
}
goto ops
ops:
defineVar("command")
setVar("command", "")
defineVar("op")
setVar("op", "")
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("op", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
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
systemStack = append(systemStack, getVar("op"))
systemStack = append(systemStack, getVar("command"))
setVar("indexes_ret","#_idx")
goto indexes
goto _idx
_idx:
setVar("op_nums", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
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
if toBool(toFloat(getVar("op_num"))>toFloat(getVar("num1")))&&toBool(toFloat(getVar("op_num"))<toFloat(getVar("num2"))){
print("")
}else{
goto _is_op_end
}
setVar("to_add","False")
goto _push_op_end
goto _is_op_end
_is_op_end:
goto _quotes_s
goto _quotes_e
_quotes_e:
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
setVar("$reverse_ret","#_rvs")
systemStack = append(systemStack, getVar("res"))
goto reverse
setVar("res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
goto _rvs
_rvs:
systemStack = append(systemStack, getVar("res"))
undefineVar("command")
undefineVar("op")
undefineVar("quotes")
undefineVar("these_quotes")
undefineVar("op_nums")
undefineVar("buf")
undefineVar("res")
undefineVar("num1")
undefineVar("num2")
undefineVar("op_num")
undefineVar("was_quote")
undefineVar("to_add")
if "#again_begin" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto again_begin
}
if "#begin" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto begin
}
if "#end_type" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto end_type
}
if "#set_void_flag_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto set_void_flag_end
}
if "#not_define_void_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_define_void_end
}
if "#func_args" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto func_args
}
if "#arg_type_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto arg_type_end
}
if "#br_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto br_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto func_end
}
if "#args_to_pass" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto args_to_pass
}
if "#not_send_pass_args" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_pass_args
}
if "#not_send_fig_br" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_fig_br
}
if "#not_send_fig_br_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_fig_br_end
}
if "#pop_start" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto pop_start
}
if "#not_send_pop" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_pop
}
if "#pop_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto pop_end
}
if "#check_begin" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto check_begin
}
if "#check_next" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto check_next
}
if "#check_again" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto check_again
}
if "#check_undefine_end0" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto check_undefine_end0
}
if "#inside_args" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto inside_args
}
if "#inside_arg_type_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto inside_arg_type_end
}
if "#undef_ret" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto undef_ret
}
if "#no_args" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto no_args
}
if "#not_send_return" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_return
}
if "#pop_args_start0" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto pop_args_start0
}
if "#not_send0" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send0
}
if "#pop_args_end0" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto pop_args_end0
}
if "#undefine_start0" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto undefine_start0
}
if "#undefine_end0" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto undefine_end0
}
if "#not_send_return_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_return_end
}
if "#inside_args_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto inside_args_end
}
if "#pop_args_start" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto pop_args_start
}
if "#not_send" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send
}
if "#pop_args_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto pop_args_end
}
if "#undefine_start" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto undefine_start
}
if "#undefine_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto undefine_end
}
if "#define_return_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto define_return_end
}
if "#start_entry" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto start_entry
}
if "#func_call_index_ret" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto func_call_index_ret
}
if "#call_start" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto call_start
}
if "#comma_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto comma_end
}
if "#aend" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto aend
}
if "#askip" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto askip
}
if "#not_send_push" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_push
}
if "#buf_mod_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto buf_mod_end
}
if "#call_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto call_end
}
if "#gl_fe" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto gl_fe
}
if "#nothing" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto nothing
}
if "#set_buf_nop_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto set_buf_nop_end
}
if "#not_send_void_pop_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_send_void_pop_end
}
if "#send_void_mark_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto send_void_mark_end
}
if "#mod_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto mod_end
}
if "#end_entry" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto end_entry
}
if "#command_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto command_end
}
if "#again_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto again_end
}
if "#not_is_first" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_is_first
}
if "#del_file" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto del_file
}
if "#end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto reverse_e
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto look_behind_e
}
if "#look_ret" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto look_ret
}
if "#indexes" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto indexes_s
}
if "#iend" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto iend
}
if "#iskip" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto iskip
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto indexes_e
}
if "#iret" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto iret
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto func_call_index
}
if "#look_behind_ret" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto look_behind_ret
}
if "#indexes_ret" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto indexes_ret
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto func_call_index_s
}
if "#is_found" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto is_found
}
if "#not_func" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto not_func
}
if "#end_func" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto end_func
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _in_stack_e
}
if "#ops" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto ops
}
if "#_idx" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _idx
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _op_nums_e
}
if "#_rvs" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _rvs
}
if "#func_end_func" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto func_end_func
}
if "#_fe" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _fe
}
if "#_fe2" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _fe2
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _braces_s
}
if "#_fe3" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _fe3
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _o_plus_end
}
if "#_fe4" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _fe4
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$ops_ret")){
goto _braces_e
}
goto func_end_func
func_end_func:
defineVar("command")
setVar("command", "")
defineVar("func_begin")
setVar("func_begin", 0)
setVar("command", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("func_begin", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
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
defineVar("stack_res")
setVar("stack_res", false)
setVar("command_len",len(fmt.Sprintf("%v",getVar("command"))))
setVar("obrace","(")
setVar("cbrace",")")
setVar("o_sum","1")
setVar("c_sum","0")
setVar("pos",sum(getVar("func_begin"), 1))
setVar("$ops_ret","#_fe")
systemStack = append(systemStack, getVar("obrace"))
systemStack = append(systemStack, getVar("command"))
goto ops
goto _fe
_fe:
setVar("obraces", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
setVar("$ops_ret","#_fe2")
systemStack = append(systemStack, getVar("cbrace"))
systemStack = append(systemStack, getVar("command"))
goto ops
goto _fe2
_fe2:
setVar("cbraces", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
goto _braces_s
_braces_s:
if toFloat(getVar("pos"))<toFloat(getVar("command_len")){
print("")
}else{
goto _braces_e
}
setVar("spos",getVar("pos"))
systemStack = append(systemStack, getVar("spos"))
systemStack = append(systemStack, getVar("obraces"))
setVar("$in_stack_ret","#_fe3")
goto in_stack
goto _fe3
_fe3:
setVar("stack_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
if toBool(getVar("stack_res")){
print("")
}else{
goto _o_plus_end
}
setVar("o_sum",sum(getVar("o_sum"), 1))
goto _o_plus_end
_o_plus_end:
systemStack = append(systemStack, getVar("spos"))
systemStack = append(systemStack, getVar("cbraces"))
setVar("$in_stack_ret","#_fe4")
goto in_stack
goto _fe4
_fe4:
setVar("stack_res", systemStack[len(systemStack)-1])
if "end" != systemStack[len(systemStack)-1] {
systemStack = systemStack[:len(systemStack)-1]
}
if toBool(getVar("stack_res")){
print("")
}else{
goto _c_plus_end
}
setVar("c_sum",sum(getVar("c_sum"), 1))
goto _c_plus_end
_c_plus_end:
if isEqual(ValueFoldInterface(getVar("o_sum")), ValueFoldInterface(getVar("c_sum"))){
goto _braces_e
}else{
print("")
}
setVar("pos",sum(getVar("pos"), 1))
goto _braces_s
goto _braces_e
_braces_e:
systemStack = append(systemStack, getVar("pos"))
undefineVar("command")
undefineVar("func_begin")
undefineVar("obraces")
undefineVar("cbraces")
undefineVar("obrace")
undefineVar("cbrace")
undefineVar("symbol")
undefineVar("o_sum")
undefineVar("c_sum")
undefineVar("pos")
undefineVar("spos")
undefineVar("command_len")
undefineVar("stack_res")
if "#again_begin" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto again_begin
}
if "#begin" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto begin
}
if "#end_type" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto end_type
}
if "#set_void_flag_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto set_void_flag_end
}
if "#not_define_void_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_define_void_end
}
if "#func_args" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto func_args
}
if "#arg_type_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto arg_type_end
}
if "#br_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto br_end
}
if "#func_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto func_end
}
if "#args_to_pass" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto args_to_pass
}
if "#not_send_pass_args" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_pass_args
}
if "#not_send_fig_br" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_fig_br
}
if "#not_send_fig_br_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_fig_br_end
}
if "#pop_start" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto pop_start
}
if "#not_send_pop" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_pop
}
if "#pop_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto pop_end
}
if "#check_begin" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto check_begin
}
if "#check_next" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto check_next
}
if "#check_again" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto check_again
}
if "#check_undefine_end0" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto check_undefine_end0
}
if "#inside_args" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto inside_args
}
if "#inside_arg_type_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto inside_arg_type_end
}
if "#undef_ret" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto undef_ret
}
if "#no_args" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto no_args
}
if "#not_send_return" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_return
}
if "#pop_args_start0" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto pop_args_start0
}
if "#not_send0" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send0
}
if "#pop_args_end0" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto pop_args_end0
}
if "#undefine_start0" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto undefine_start0
}
if "#undefine_end0" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto undefine_end0
}
if "#not_send_return_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_return_end
}
if "#inside_args_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto inside_args_end
}
if "#pop_args_start" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto pop_args_start
}
if "#not_send" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send
}
if "#pop_args_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto pop_args_end
}
if "#undefine_start" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto undefine_start
}
if "#undefine_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto undefine_end
}
if "#define_return_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto define_return_end
}
if "#start_entry" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto start_entry
}
if "#func_call_index_ret" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto func_call_index_ret
}
if "#call_start" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto call_start
}
if "#comma_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto comma_end
}
if "#aend" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto aend
}
if "#askip" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto askip
}
if "#not_send_push" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_push
}
if "#buf_mod_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto buf_mod_end
}
if "#call_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto call_end
}
if "#gl_fe" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto gl_fe
}
if "#nothing" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto nothing
}
if "#set_buf_nop_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto set_buf_nop_end
}
if "#not_send_void_pop_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_send_void_pop_end
}
if "#send_void_mark_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto send_void_mark_end
}
if "#mod_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto mod_end
}
if "#end_entry" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto end_entry
}
if "#command_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto command_end
}
if "#again_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto again_end
}
if "#not_is_first" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_is_first
}
if "#del_file" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto del_file
}
if "#end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto end
}
if "#reverse" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto reverse
}
if "#reverse_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto reverse_s
}
if "#reverse_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto reverse_e
}
if "#look_behind" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto look_behind
}
if "#look_behind_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto look_behind_s
}
if "#is_not_zero" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto is_not_zero
}
if "#look_behind_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto look_behind_e
}
if "#look_ret" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto look_ret
}
if "#indexes" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto indexes
}
if "#indexes_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto indexes_s
}
if "#iend" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto iend
}
if "#iskip" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto iskip
}
if "#indexes_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto indexes_e
}
if "#iret" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto iret
}
if "#func_call_index" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto func_call_index
}
if "#look_behind_ret" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto look_behind_ret
}
if "#indexes_ret" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto indexes_ret
}
if "#func_call_index_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto func_call_index_s
}
if "#is_found" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto is_found
}
if "#not_func" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto not_func
}
if "#end_func" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto end_func
}
if "#in_stack" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto in_stack
}
if "#_in_stack_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _in_stack_s
}
if "#_no" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _no
}
if "#_in_stack_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _in_stack_e
}
if "#ops" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto ops
}
if "#_idx" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _idx
}
if "#_op_nums_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _op_nums_s
}
if "#_quotes_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _quotes_s
}
if "#_these_quotes_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _these_quotes_s
}
if "#_these_quotes_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _these_quotes_e
}
if "#_is_op_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _is_op_end
}
if "#_quotes_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _quotes_e
}
if "#_push_op_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _push_op_end
}
if "#_op_nums_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _op_nums_e
}
if "#_rvs" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _rvs
}
if "#func_end_func" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto func_end_func
}
if "#_fe" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _fe
}
if "#_fe2" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _fe2
}
if "#_braces_s" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _braces_s
}
if "#_fe3" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _fe3
}
if "#_o_plus_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _o_plus_end
}
if "#_fe4" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _fe4
}
if "#_c_plus_end" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _c_plus_end
}
if "#_braces_e" == fmt.Sprintf("%v", getVar("$func_end_res")){
goto _braces_e
}
}
