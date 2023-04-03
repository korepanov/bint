.data
pageSize:
.quad 4096 
varNameSize:
.quad 32
varSize:
.quad 128
typeSize:
.quad 32 
valSize:
.quad 64 
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
strBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrBegin = . - strBegin
strPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrPointer = . - strPointer
strMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrMax = . - strMax 
isNeg:
.byte 0 
lenIsNeg = . - isNeg
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


 t0: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT0 = . - t0
 t1: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT1 = . - t1
 t2: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT2 = . - t2
 t3: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT3 = . - t3
 t4: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT4 = . - t4
 t5: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT5 = . - t5
 t6: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT6 = . - t6
 t7: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT7 = . - t7
 t8: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT8 = . - t8
 t9: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT9 = . - t9
varName0:
.ascii "a"
lenVarName0 = . - varName0
varName1:
.ascii "b"
lenVarName1 = . - varName1
varName2:
.ascii "c"
lenVarName2 = . - varName2
data0:
.ascii "3"
.space 1, 0
lenData0 = . - data0
data1:
.ascii "7"
.space 1, 0
lenData1 = . - data1
data2:
.ascii "10.5"
.space 1, 0
lenData2 = . - data2
data3:
.ascii "1"
.space 1, 0
lenData3 = . - data3
data4:
.ascii "2"
.space 1, 0
lenData4 = . - data4
data5:
.ascii "3"
.space 1, 0
lenData5 = . - data5
data6:
.ascii "4"
.space 1, 0
lenData6 = . - data6
data7:
.ascii "5"
.space 1, 0
lenData7 = . - data7
data8:
.ascii "6"
.space 1, 0
lenData8 = . - data8
data9:
.ascii "7"
.space 1, 0
lenData9 = . - data9
data10:
.ascii "8"
.space 1, 0
lenData10 = . - data10
data11:
.ascii "9"
.space 1, 0
lenData11 = . - data11
data12:
.ascii "10"
.space 1, 0
lenData12 = . - data12
data13:
.ascii "11"
.space 1, 0
lenData13 = . - data13
data14:
.ascii "12"
.space 1, 0
lenData14 = . - data14
varName3:
.ascii "s"
lenVarName3 = . - varName3
data15:
.ascii "\n"
.space 1, 0
lenData15 = . - data15
