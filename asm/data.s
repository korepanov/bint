.data
starSymbol:
.ascii "*"
endSymbol:
.ascii ";"
deltaSize:
.quad 0
pageSize:
.quad 262144 # 4096
shiftSize:
.quad 262144 # 4096
varNameSize:
.quad 32
varSize:
.quad 128
typeSize:
.quad 32 
valSize:
.quad 64 
strValSize:
.quad 32768 # 1024
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
 t10: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT10 = . - t10
 t11: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT11 = . - t11
 t12: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT12 = . - t12
 t13: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT13 = . - t13
 t14: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT14 = . - t14
 t15: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT15 = . - t15
 t16: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT16 = . - t16
 t17: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT17 = . - t17
 t18: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT18 = . - t18
 t19: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT19 = . - t19
 t20: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT20 = . - t20
 t21: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT21 = . - t21
 t22: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT22 = . - t22
 t23: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT23 = . - t23
 t24: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT24 = . - t24
 t25: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT25 = . - t25
 t26: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT26 = . - t26
 t27: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT27 = . - t27
 t28: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT28 = . - t28
 t29: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT29 = . - t29
 t30: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT30 = . - t30
 t31: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT31 = . - t31
 t32: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT32 = . - t32
 t33: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT33 = . - t33
 t34: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT34 = . - t34
 t35: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT35 = . - t35
 t36: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT36 = . - t36
 t37: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT37 = . - t37
 t38: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT38 = . - t38
 t39: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT39 = . - t39
 t40: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT40 = . - t40
 t41: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT41 = . - t41
 t42: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT42 = . - t42
 t43: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT43 = . - t43
 t44: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT44 = . - t44
 t45: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT45 = . - t45
 t46: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT46 = . - t46
 t47: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT47 = . - t47
 t48: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT48 = . - t48
 t49: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT49 = . - t49
 t50: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT50 = . - t50
 t51: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT51 = . - t51
 t52: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT52 = . - t52
 t53: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT53 = . - t53
 t54: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT54 = . - t54
 t55: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT55 = . - t55
 t56: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT56 = . - t56
 t57: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT57 = . - t57
 t58: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT58 = . - t58
 t59: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT59 = . - t59
 t60: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT60 = . - t60
 t61: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT61 = . - t61
 t62: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT62 = . - t62
 t63: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT63 = . - t63
 t64: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT64 = . - t64
 t65: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT65 = . - t65
 t66: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT66 = . - t66
 t67: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT67 = . - t67
 t68: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT68 = . - t68
 t69: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT69 = . - t69
 t70: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT70 = . - t70
 t71: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT71 = . - t71
 t72: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT72 = . - t72
 t73: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT73 = . - t73
 t74: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT74 = . - t74
 t75: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT75 = . - t75
 t76: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT76 = . - t76
 t77: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT77 = . - t77
 t78: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT78 = . - t78
 t79: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT79 = . - t79
 t80: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT80 = . - t80
 t81: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT81 = . - t81
 t82: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT82 = . - t82
 t83: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT83 = . - t83
 t84: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT84 = . - t84
 t85: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT85 = . - t85
 t86: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT86 = . - t86
 t87: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT87 = . - t87
 t88: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT88 = . - t88
 t89: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT89 = . - t89
 t90: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT90 = . - t90
 t91: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT91 = . - t91
 t92: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT92 = . - t92
 t93: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT93 = . - t93
 t94: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT94 = . - t94
 t95: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT95 = . - t95
 t96: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT96 = . - t96
 t97: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT97 = . - t97
 t98: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT98 = . - t98
 t99: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT99 = . - t99
 t100: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT100 = . - t100
 t101: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT101 = . - t101
 t102: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT102 = . - t102
 t103: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT103 = . - t103
 t104: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT104 = . - t104
 t105: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT105 = . - t105
 t106: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT106 = . - t106
 t107: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT107 = . - t107
 t108: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT108 = . - t108
 t109: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT109 = . - t109
 t110: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT110 = . - t110
 t111: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT111 = . - t111
 t112: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT112 = . - t112
 t113: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT113 = . - t113
 t114: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT114 = . - t114
 t115: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT115 = . - t115
 t116: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT116 = . - t116
 t117: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT117 = . - t117
 t118: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT118 = . - t118
 t119: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT119 = . - t119
 t120: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT120 = . - t120
 t121: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT121 = . - t121
 t122: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT122 = . - t122
 t123: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT123 = . - t123
 t124: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT124 = . - t124
 t125: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT125 = . - t125
 t126: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT126 = . - t126
 t127: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenT127 = . - t127
systemVarName0:
.ascii "^systemVar0"
.space 1, 0
lenSystemVarName0 = . - systemVarName0
systemVarName1:
.ascii "^systemVar1"
.space 1, 0
lenSystemVarName1 = . - systemVarName1
systemVarName2:
.ascii "^systemVar2"
.space 1, 0
lenSystemVarName2 = . - systemVarName2
systemVarName3:
.ascii "^systemVar3"
.space 1, 0
lenSystemVarName3 = . - systemVarName3
systemVarName4:
.ascii "^systemVar4"
.space 1, 0
lenSystemVarName4 = . - systemVarName4
systemVarName5:
.ascii "^systemVar5"
.space 1, 0
lenSystemVarName5 = . - systemVarName5
systemVarName6:
.ascii "^systemVar6"
.space 1, 0
lenSystemVarName6 = . - systemVarName6
systemVarName7:
.ascii "^systemVar7"
.space 1, 0
lenSystemVarName7 = . - systemVarName7
systemVarName8:
.ascii "^systemVar8"
.space 1, 0
lenSystemVarName8 = . - systemVarName8
systemVarName9:
.ascii "^systemVar9"
.space 1, 0
lenSystemVarName9 = . - systemVarName9
systemVarName10:
.ascii "^systemVar10"
.space 1, 0
lenSystemVarName10 = . - systemVarName10
systemVarName11:
.ascii "^systemVar11"
.space 1, 0
lenSystemVarName11 = . - systemVarName11
systemVarName12:
.ascii "^systemVar12"
.space 1, 0
lenSystemVarName12 = . - systemVarName12
systemVarName13:
.ascii "^systemVar13"
.space 1, 0
lenSystemVarName13 = . - systemVarName13
systemVarName14:
.ascii "^systemVar14"
.space 1, 0
lenSystemVarName14 = . - systemVarName14
systemVarName15:
.ascii "^systemVar15"
.space 1, 0
lenSystemVarName15 = . - systemVarName15
systemVarName16:
.ascii "^systemVar16"
.space 1, 0
lenSystemVarName16 = . - systemVarName16
systemVarName17:
.ascii "^systemVar17"
.space 1, 0
lenSystemVarName17 = . - systemVarName17
systemVarName18:
.ascii "^systemVar18"
.space 1, 0
lenSystemVarName18 = . - systemVarName18
systemVarName19:
.ascii "^systemVar19"
.space 1, 0
lenSystemVarName19 = . - systemVarName19
systemVarName20:
.ascii "^systemVar20"
.space 1, 0
lenSystemVarName20 = . - systemVarName20
systemVarName21:
.ascii "^systemVar21"
.space 1, 0
lenSystemVarName21 = . - systemVarName21
systemVarName22:
.ascii "^systemVar22"
.space 1, 0
lenSystemVarName22 = . - systemVarName22
systemVarName23:
.ascii "^systemVar23"
.space 1, 0
lenSystemVarName23 = . - systemVarName23
systemVarName24:
.ascii "^systemVar24"
.space 1, 0
lenSystemVarName24 = . - systemVarName24
systemVarName25:
.ascii "^systemVar25"
.space 1, 0
lenSystemVarName25 = . - systemVarName25
systemVarName26:
.ascii "^systemVar26"
.space 1, 0
lenSystemVarName26 = . - systemVarName26
systemVarName27:
.ascii "^systemVar27"
.space 1, 0
lenSystemVarName27 = . - systemVarName27
systemVarName28:
.ascii "^systemVar28"
.space 1, 0
lenSystemVarName28 = . - systemVarName28
systemVarName29:
.ascii "^systemVar29"
.space 1, 0
lenSystemVarName29 = . - systemVarName29
systemVarName30:
.ascii "^systemVar30"
.space 1, 0
lenSystemVarName30 = . - systemVarName30
systemVarName31:
.ascii "^systemVar31"
.space 1, 0
lenSystemVarName31 = . - systemVarName31
systemVarName32:
.ascii "^systemVar32"
.space 1, 0
lenSystemVarName32 = . - systemVarName32
systemVarName33:
.ascii "^systemVar33"
.space 1, 0
lenSystemVarName33 = . - systemVarName33
systemVarName34:
.ascii "^systemVar34"
.space 1, 0
lenSystemVarName34 = . - systemVarName34
systemVarName35:
.ascii "^systemVar35"
.space 1, 0
lenSystemVarName35 = . - systemVarName35
systemVarName36:
.ascii "^systemVar36"
.space 1, 0
lenSystemVarName36 = . - systemVarName36
systemVarName37:
.ascii "^systemVar37"
.space 1, 0
lenSystemVarName37 = . - systemVarName37
systemVarName38:
.ascii "^systemVar38"
.space 1, 0
lenSystemVarName38 = . - systemVarName38
systemVarName39:
.ascii "^systemVar39"
.space 1, 0
lenSystemVarName39 = . - systemVarName39
systemVarName40:
.ascii "^systemVar40"
.space 1, 0
lenSystemVarName40 = . - systemVarName40
systemVarName41:
.ascii "^systemVar41"
.space 1, 0
lenSystemVarName41 = . - systemVarName41
systemVarName42:
.ascii "^systemVar42"
.space 1, 0
lenSystemVarName42 = . - systemVarName42
systemVarName43:
.ascii "^systemVar43"
.space 1, 0
lenSystemVarName43 = . - systemVarName43
systemVarName44:
.ascii "^systemVar44"
.space 1, 0
lenSystemVarName44 = . - systemVarName44
systemVarName45:
.ascii "^systemVar45"
.space 1, 0
lenSystemVarName45 = . - systemVarName45
systemVarName46:
.ascii "^systemVar46"
.space 1, 0
lenSystemVarName46 = . - systemVarName46
systemVarName47:
.ascii "^systemVar47"
.space 1, 0
lenSystemVarName47 = . - systemVarName47
systemVarName48:
.ascii "^systemVar48"
.space 1, 0
lenSystemVarName48 = . - systemVarName48
systemVarName49:
.ascii "^systemVar49"
.space 1, 0
lenSystemVarName49 = . - systemVarName49
systemVarName50:
.ascii "^systemVar50"
.space 1, 0
lenSystemVarName50 = . - systemVarName50
systemVarName51:
.ascii "^systemVar51"
.space 1, 0
lenSystemVarName51 = . - systemVarName51
systemVarName52:
.ascii "^systemVar52"
.space 1, 0
lenSystemVarName52 = . - systemVarName52
systemVarName53:
.ascii "^systemVar53"
.space 1, 0
lenSystemVarName53 = . - systemVarName53
systemVarName54:
.ascii "^systemVar54"
.space 1, 0
lenSystemVarName54 = . - systemVarName54
systemVarName55:
.ascii "^systemVar55"
.space 1, 0
lenSystemVarName55 = . - systemVarName55
systemVarName56:
.ascii "^systemVar56"
.space 1, 0
lenSystemVarName56 = . - systemVarName56
systemVarName57:
.ascii "^systemVar57"
.space 1, 0
lenSystemVarName57 = . - systemVarName57
systemVarName58:
.ascii "^systemVar58"
.space 1, 0
lenSystemVarName58 = . - systemVarName58
systemVarName59:
.ascii "^systemVar59"
.space 1, 0
lenSystemVarName59 = . - systemVarName59
systemVarName60:
.ascii "^systemVar60"
.space 1, 0
lenSystemVarName60 = . - systemVarName60
systemVarName61:
.ascii "^systemVar61"
.space 1, 0
lenSystemVarName61 = . - systemVarName61
systemVarName62:
.ascii "^systemVar62"
.space 1, 0
lenSystemVarName62 = . - systemVarName62
systemVarName63:
.ascii "^systemVar63"
.space 1, 0
lenSystemVarName63 = . - systemVarName63
systemVarName64:
.ascii "^systemVar64"
.space 1, 0
lenSystemVarName64 = . - systemVarName64
systemVarName65:
.ascii "^systemVar65"
.space 1, 0
lenSystemVarName65 = . - systemVarName65
systemVarName66:
.ascii "^systemVar66"
.space 1, 0
lenSystemVarName66 = . - systemVarName66
systemVarName67:
.ascii "^systemVar67"
.space 1, 0
lenSystemVarName67 = . - systemVarName67
systemVarName68:
.ascii "^systemVar68"
.space 1, 0
lenSystemVarName68 = . - systemVarName68
systemVarName69:
.ascii "^systemVar69"
.space 1, 0
lenSystemVarName69 = . - systemVarName69
systemVarName70:
.ascii "^systemVar70"
.space 1, 0
lenSystemVarName70 = . - systemVarName70
systemVarName71:
.ascii "^systemVar71"
.space 1, 0
lenSystemVarName71 = . - systemVarName71
systemVarName72:
.ascii "^systemVar72"
.space 1, 0
lenSystemVarName72 = . - systemVarName72
systemVarName73:
.ascii "^systemVar73"
.space 1, 0
lenSystemVarName73 = . - systemVarName73
systemVarName74:
.ascii "^systemVar74"
.space 1, 0
lenSystemVarName74 = . - systemVarName74
systemVarName75:
.ascii "^systemVar75"
.space 1, 0
lenSystemVarName75 = . - systemVarName75
systemVarName76:
.ascii "^systemVar76"
.space 1, 0
lenSystemVarName76 = . - systemVarName76
systemVarName77:
.ascii "^systemVar77"
.space 1, 0
lenSystemVarName77 = . - systemVarName77
systemVarName78:
.ascii "^systemVar78"
.space 1, 0
lenSystemVarName78 = . - systemVarName78
systemVarName79:
.ascii "^systemVar79"
.space 1, 0
lenSystemVarName79 = . - systemVarName79
systemVarName80:
.ascii "^systemVar80"
.space 1, 0
lenSystemVarName80 = . - systemVarName80
systemVarName81:
.ascii "^systemVar81"
.space 1, 0
lenSystemVarName81 = . - systemVarName81
systemVarName82:
.ascii "^systemVar82"
.space 1, 0
lenSystemVarName82 = . - systemVarName82
systemVarName83:
.ascii "^systemVar83"
.space 1, 0
lenSystemVarName83 = . - systemVarName83
systemVarName84:
.ascii "^systemVar84"
.space 1, 0
lenSystemVarName84 = . - systemVarName84
systemVarName85:
.ascii "^systemVar85"
.space 1, 0
lenSystemVarName85 = . - systemVarName85
systemVarName86:
.ascii "^systemVar86"
.space 1, 0
lenSystemVarName86 = . - systemVarName86
systemVarName87:
.ascii "^systemVar87"
.space 1, 0
lenSystemVarName87 = . - systemVarName87
systemVarName88:
.ascii "^systemVar88"
.space 1, 0
lenSystemVarName88 = . - systemVarName88
systemVarName89:
.ascii "^systemVar89"
.space 1, 0
lenSystemVarName89 = . - systemVarName89
systemVarName90:
.ascii "^systemVar90"
.space 1, 0
lenSystemVarName90 = . - systemVarName90
systemVarName91:
.ascii "^systemVar91"
.space 1, 0
lenSystemVarName91 = . - systemVarName91
systemVarName92:
.ascii "^systemVar92"
.space 1, 0
lenSystemVarName92 = . - systemVarName92
systemVarName93:
.ascii "^systemVar93"
.space 1, 0
lenSystemVarName93 = . - systemVarName93
systemVarName94:
.ascii "^systemVar94"
.space 1, 0
lenSystemVarName94 = . - systemVarName94
systemVarName95:
.ascii "^systemVar95"
.space 1, 0
lenSystemVarName95 = . - systemVarName95
systemVarName96:
.ascii "^systemVar96"
.space 1, 0
lenSystemVarName96 = . - systemVarName96
systemVarName97:
.ascii "^systemVar97"
.space 1, 0
lenSystemVarName97 = . - systemVarName97
systemVarName98:
.ascii "^systemVar98"
.space 1, 0
lenSystemVarName98 = . - systemVarName98
systemVarName99:
.ascii "^systemVar99"
.space 1, 0
lenSystemVarName99 = . - systemVarName99
systemVarName100:
.ascii "^systemVar100"
.space 1, 0
lenSystemVarName100 = . - systemVarName100
systemVarName101:
.ascii "^systemVar101"
.space 1, 0
lenSystemVarName101 = . - systemVarName101
systemVarName102:
.ascii "^systemVar102"
.space 1, 0
lenSystemVarName102 = . - systemVarName102
systemVarName103:
.ascii "^systemVar103"
.space 1, 0
lenSystemVarName103 = . - systemVarName103
systemVarName104:
.ascii "^systemVar104"
.space 1, 0
lenSystemVarName104 = . - systemVarName104
systemVarName105:
.ascii "^systemVar105"
.space 1, 0
lenSystemVarName105 = . - systemVarName105
systemVarName106:
.ascii "^systemVar106"
.space 1, 0
lenSystemVarName106 = . - systemVarName106
systemVarName107:
.ascii "^systemVar107"
.space 1, 0
lenSystemVarName107 = . - systemVarName107
systemVarName108:
.ascii "^systemVar108"
.space 1, 0
lenSystemVarName108 = . - systemVarName108
systemVarName109:
.ascii "^systemVar109"
.space 1, 0
lenSystemVarName109 = . - systemVarName109
systemVarName110:
.ascii "^systemVar110"
.space 1, 0
lenSystemVarName110 = . - systemVarName110
systemVarName111:
.ascii "^systemVar111"
.space 1, 0
lenSystemVarName111 = . - systemVarName111
systemVarName112:
.ascii "^systemVar112"
.space 1, 0
lenSystemVarName112 = . - systemVarName112
systemVarName113:
.ascii "^systemVar113"
.space 1, 0
lenSystemVarName113 = . - systemVarName113
systemVarName114:
.ascii "^systemVar114"
.space 1, 0
lenSystemVarName114 = . - systemVarName114
systemVarName115:
.ascii "^systemVar115"
.space 1, 0
lenSystemVarName115 = . - systemVarName115
systemVarName116:
.ascii "^systemVar116"
.space 1, 0
lenSystemVarName116 = . - systemVarName116
systemVarName117:
.ascii "^systemVar117"
.space 1, 0
lenSystemVarName117 = . - systemVarName117
systemVarName118:
.ascii "^systemVar118"
.space 1, 0
lenSystemVarName118 = . - systemVarName118
systemVarName119:
.ascii "^systemVar119"
.space 1, 0
lenSystemVarName119 = . - systemVarName119
systemVarName120:
.ascii "^systemVar120"
.space 1, 0
lenSystemVarName120 = . - systemVarName120
systemVarName121:
.ascii "^systemVar121"
.space 1, 0
lenSystemVarName121 = . - systemVarName121
systemVarName122:
.ascii "^systemVar122"
.space 1, 0
lenSystemVarName122 = . - systemVarName122
systemVarName123:
.ascii "^systemVar123"
.space 1, 0
lenSystemVarName123 = . - systemVarName123
systemVarName124:
.ascii "^systemVar124"
.space 1, 0
lenSystemVarName124 = . - systemVarName124
systemVarName125:
.ascii "^systemVar125"
.space 1, 0
lenSystemVarName125 = . - systemVarName125
systemVarName126:
.ascii "^systemVar126"
.space 1, 0
lenSystemVarName126 = . - systemVarName126
systemVarName127:
.ascii "^systemVar127"
.space 1, 0
lenSystemVarName127 = . - systemVarName127
 st0: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt0 = . - st0
 st1: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt1 = . - st1
 st2: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt2 = . - st2
 st3: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt3 = . - st3
 st4: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt4 = . - st4
 st5: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt5 = . - st5
 st6: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt6 = . - st6
 st7: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt7 = . - st7
 st8: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt8 = . - st8
 st9: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt9 = . - st9
 st10: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt10 = . - st10
 st11: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt11 = . - st11
 st12: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt12 = . - st12
 st13: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt13 = . - st13
 st14: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt14 = . - st14
 st15: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt15 = . - st15
 st16: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt16 = . - st16
 st17: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt17 = . - st17
 st18: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt18 = . - st18
 st19: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt19 = . - st19
 st20: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt20 = . - st20
 st21: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt21 = . - st21
 st22: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt22 = . - st22
 st23: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt23 = . - st23
 st24: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt24 = . - st24
 st25: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt25 = . - st25
 st26: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt26 = . - st26
 st27: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt27 = . - st27
 st28: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt28 = . - st28
 st29: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt29 = . - st29
 st30: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt30 = . - st30
 st31: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt31 = . - st31
 st32: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt32 = . - st32
 st33: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt33 = . - st33
 st34: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt34 = . - st34
 st35: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt35 = . - st35
 st36: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt36 = . - st36
 st37: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt37 = . - st37
 st38: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt38 = . - st38
 st39: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt39 = . - st39
 st40: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt40 = . - st40
 st41: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt41 = . - st41
 st42: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt42 = . - st42
 st43: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt43 = . - st43
 st44: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt44 = . - st44
 st45: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt45 = . - st45
 st46: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt46 = . - st46
 st47: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt47 = . - st47
 st48: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt48 = . - st48
 st49: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt49 = . - st49
 st50: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt50 = . - st50
 st51: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt51 = . - st51
 st52: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt52 = . - st52
 st53: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt53 = . - st53
 st54: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt54 = . - st54
 st55: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt55 = . - st55
 st56: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt56 = . - st56
 st57: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt57 = . - st57
 st58: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt58 = . - st58
 st59: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt59 = . - st59
 st60: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt60 = . - st60
 st61: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt61 = . - st61
 st62: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt62 = . - st62
 st63: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt63 = . - st63
 st64: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt64 = . - st64
 st65: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt65 = . - st65
 st66: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt66 = . - st66
 st67: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt67 = . - st67
 st68: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt68 = . - st68
 st69: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt69 = . - st69
 st70: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt70 = . - st70
 st71: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt71 = . - st71
 st72: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt72 = . - st72
 st73: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt73 = . - st73
 st74: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt74 = . - st74
 st75: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt75 = . - st75
 st76: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt76 = . - st76
 st77: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt77 = . - st77
 st78: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt78 = . - st78
 st79: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt79 = . - st79
 st80: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt80 = . - st80
 st81: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt81 = . - st81
 st82: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt82 = . - st82
 st83: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt83 = . - st83
 st84: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt84 = . - st84
 st85: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt85 = . - st85
 st86: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt86 = . - st86
 st87: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt87 = . - st87
 st88: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt88 = . - st88
 st89: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt89 = . - st89
 st90: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt90 = . - st90
 st91: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt91 = . - st91
 st92: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt92 = . - st92
 st93: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt93 = . - st93
 st94: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt94 = . - st94
 st95: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt95 = . - st95
 st96: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt96 = . - st96
 st97: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt97 = . - st97
 st98: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt98 = . - st98
 st99: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt99 = . - st99
 st100: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt100 = . - st100
 st101: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt101 = . - st101
 st102: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt102 = . - st102
 st103: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt103 = . - st103
 st104: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt104 = . - st104
 st105: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt105 = . - st105
 st106: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt106 = . - st106
 st107: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt107 = . - st107
 st108: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt108 = . - st108
 st109: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt109 = . - st109
 st110: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt110 = . - st110
 st111: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt111 = . - st111
 st112: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt112 = . - st112
 st113: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt113 = . - st113
 st114: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt114 = . - st114
 st115: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt115 = . - st115
 st116: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt116 = . - st116
 st117: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt117 = . - st117
 st118: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt118 = . - st118
 st119: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt119 = . - st119
 st120: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt120 = . - st120
 st121: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt121 = . - st121
 st122: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt122 = . - st122
 st123: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt123 = . - st123
 st124: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt124 = . - st124
 st125: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt125 = . - st125
 st126: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt126 = . - st126
 st127: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt127 = . - st127
varName0:
.ascii "file"
.space 1, 0
lenVarName0 = . - varName0
varName1:
.ascii "mode"
.space 1, 0
lenVarName1 = . - varName1
varName2:
.ascii "bytes"
.space 1, 0
lenVarName2 = . - varName2
varName3:
.ascii "f"
.space 1, 0
lenVarName3 = . - varName3
varName4:
.ascii "size"
.space 1, 0
lenVarName4 = . - varName4
varName5:
.ascii "s"
.space 1, 0
lenVarName5 = . - varName5
varName6:
.ascii "res"
.space 1, 0
lenVarName6 = . - varName6
data0:
.ascii "4096"
.space 1, 0
lenData0 = . - data0
data1:
.ascii "/home/slava/books/cookbook.txt"
.space 1, 0
lenData1 = . - data1
data2:
.ascii "0"
.space 1, 0
lenData2 = . - data2
label0:
 .quad .again
labelName0:
.ascii ".again"
.space 1,0
data3:
.ascii ""
.space 1, 0
lenData3 = . - data3
data4:
.ascii ""
.space 1, 0
lenData4 = . - data4
label1:
 .quad .end
labelName1:
.ascii ".end"
.space 1,0
data5:
.ascii ""
.space 1, 0
lenData5 = . - data5
