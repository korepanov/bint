package main 

import(
"strings"
"fmt"
"io"
"os"
"strconv"
"unicode"
"errors"
)
var systemStack = []interface{}{"end"}
var vars = make(map[string][]interface{})

func defineVar(varName string){
vars[varName] = append(vars[varName], interface{}(nil))
}

func undefineVar(varName string){
vars[varName] = vars[varName][:len(vars[varName]) - 1]
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
if "True" == s.(string){
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

func isEqual(a interface{}, b interface{}) bool{
if "\"" == string(fmt.Sprintf("%v", a)[0]) && "\"" == string(fmt.Sprintf("%v", b)[0]) {
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
if "\"" == string(fmt.Sprintf("%v", a)[0]) && "\"" == string(fmt.Sprintf("%v", b)[0]) {
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

func main(){

_ = strings.Index("", "")
fmt.Printf("")
f, _ := os.Create("temp.b")
f2, _ := os.Create("temp2.b")
io.Copy(f, f2)
os.Remove(f.Name())
os.Remove(f2.Name());
_, _ = strconv.ParseFloat("64", 64)
_ = unicode.IsLetter('a')
	
	
