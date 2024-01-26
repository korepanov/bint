.data
clearSymbol:
.ascii "$clear"
.space 1, 0
lenClearSymbol = . - clearSymbol 
starSymbol:
.ascii "*"
endSymbol:
.ascii ";"
deltaSize:
.quad 0
pageSize:
.quad 4096 
shiftSize:
.quad 8192 
varNameSize:
.quad 32
varSize:
.quad 128
typeSize:
.quad 32 
valSize:
.quad 64 
strValSize:
.quad 1024 
labelSize:
.quad 128 
labelsMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
buf:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf = . - buf 
buf2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf2 = . - buf2 
buf3:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf3 = . - buf3
buf4:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf4 = . - buf4
inputBuf:
.byte 0
lenInputBuf = . - inputBuf 
readBuf:
.byte 0, 0
lenReadBuf = . - readBuf
readBufSum:
.space 4096
lenReadBufSum = . - readBufSum
numberOfReadBytes:
.quad 0
userMem:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserMem = . - userMem 
userMem2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserMem2 = . - userMem2 
mem:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem = . - mem 
mem2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem2 = . - mem2 
mem3:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem3 = . - mem3 
mem4:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem4 = . - mem4 
mem5:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem5 = . - mem5
mem6:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem6 = . - mem6
mem7:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem7 = . - mem7
mem8:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem8 = . - mem8
mem9:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem9 = . - mem9
mem10:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem10 = . - mem10 
mem11:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem11 = . - mem11 
mem12:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem12 = . - mem12 
mem13:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem13 = . - mem13 
mem14:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem14 = . - mem14 
mem15:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem15 = . - mem15 
mem16:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem16 = . - mem16 
mem17:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem17 = . - mem17 
mem18:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem18 = . - mem18 
mem19:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem19 = . - mem19 
mem20:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem20 = . - mem20  
strBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrBegin = . - strBegin
oldHeapMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenOldHeapMax = . - oldHeapMax
strPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrPointer = . - strPointer
strMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrMax = . - strMax 
isNeg:
.byte 0 
lenIsNeg = . - isNeg
userConcatinateFlag:
.byte 0
lenUserConcatinateFlag = . - userConcatinateFlag
userConcatinateFlag2:
.byte 0
lenUserConcatinateFlag2 = . - userConcatinateFlag2
isExpNeg:
.byte 0
lenIsExpNeg = . - isExpNeg
varType:
.quad 0, 0, 0, 0
lenVarType = . - varType 
varName:
.quad 0, 0, 0, 0
lenVarName = . - varName 
userData:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserData = . - userData
lastUserData:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLastUserData = . - lastUserData
strPageNumber:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrNumber = . - strPageNumber
memorySize:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMemorySize = . - memorySize
memoryBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMemoryBegin = . - memoryBegin
labelsEnd:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLabelsEnd = . - labelsEnd
labelsPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLabelsPointer = . - labelsPointer 
varNamePanic:
.ascii "^toPanic"
.space 1, 0
lenVarNamePanic = . - varNamePanic 
varNameError:
.ascii "error"
.space 1, 0
lenVarNameError = . - varNameError 
systemVarName:
.ascii "^systemVar"
.space 1, 0
lenSystemVarName = . - systemVarName
popVarName:
.ascii "^popVar"
.space 1, 0
lenPopVarName = . - popVarName

intType:
.ascii "int"
.space 1, 0
lenIntType = . - intType
floatType:
.ascii "float"
.space 1, 0
lenFloatType = . - floatType
boolType:
.ascii "bool"
.space 1, 0
lenBoolType = . - boolType
stringType:
.ascii "string"
.space 1, 0
lenStringType = . - stringType
enter:
.ascii "\n"
.space 1, 0
lenEnter = . - enter 
bigTrueVal:
.ascii "True"
.space 1, 0
lenBigTrueVal = . - bigTrueVal
trueVal:
.ascii "true"
.space 1, 0
lenTrueVal = . - trueVal 
bigFalseVal:
.ascii "False"
.space 1, 0
lenBigFalseVal = . - bigFalseVal 
falseVal:
.ascii "false"
.space 1, 0
lenFalseVal = . - falseVal 
oneVal:
.ascii "1"
.space 1, 0
lenOneVal = . - oneVal 
zeroVal:
.ascii "0"
.space 1, 0
lenZeroVal = . - zeroVal 

ten:
.float 10.0 
one:
.float 1.0 
zero:
.float 0.0 
floatTail:
.ascii ".0"
.space 1, 0 
lenFloatTail =  . - floatTail

fatalError:
.ascii "fatal error: internal error\n"
.space 1, 0 
divZeroError:
.ascii "runtime error: dividing by zero"
.space 1, 0
divINegError:
.ascii "runtime error: @ is not defined for negative numbers"
.space 1, 0
powNegError:
.ascii "runtime error: ^ is not defined for negative base and fractional exponent"
.space 1, 0
powZeroNegError:
.ascii "runtime error: ^ is not defined for zero base and negative exponent"
.space 1, 0
powZeroZeroError:
.ascii "runtime error: ^ is not defined for zero base and zero exponent"
.space 1, 0
noSuchMarkError:
.ascii "runtime error: no such mark: "
.space 1, 0
concError:
.ascii "could not concatinate not string arguments"
.space 1, 0 
strError:
.ascii "the type of the variable to which you want to assign the result of string concatenation is not a string"
.space 1, 0 
parseNumberError:
.ascii "could not parse number: invalid number format"
.space 1, 0  
parseBoolError:
.ascii "could not parse bool: invalid bool format"
.space 1, 0 
sliceBoundError:
.ascii "slice index is out of bounds"
.space 1, 0 
isLetterError:
.ascii "is_letter argument length error"
.space 1, 0 
isDigitError:
.ascii "is_digit argument length error"
.space 1, 0 
openFileError:
.ascii "could not open file"
.space 1, 0
readFromFileError:
.ascii "could not read file"
.space 1, 0
