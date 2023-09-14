package compilerVars

var CompilerVars map[string]int // карта, которая хранит соответствие имен переменных кучи и номеров имен переменных asm
var VarsCounter int             // счетчик переменных asm
var DataNumber int              // счетчик данных asm
var LabelCounter int            // счетчик пользовательских меток asm
var BranchCounter int           // число ветвлений
var TempVarsNum = 128           // количество временных переменных в asm на одно арифметическое выражение
var TempStringVarsNum = 128     // количество временных пользовательских переменных типа string на одну конкатенацию
