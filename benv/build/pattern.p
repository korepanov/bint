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
return fmt.Sprintf("%v", a) + fmt.Sprintf("%v", b)
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

func getRootSource(source string) string{
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
	const chunkSize = 100
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
	
	
