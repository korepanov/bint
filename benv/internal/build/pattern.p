package main 

import(
	"strings"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)
func nop(interface{}){
	fmt.Printf("")
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
	
	
