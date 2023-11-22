.data
starSymbol:
.ascii "*"
endSymbol:
.ascii ";"
deltaSize:
.quad 1024
pageSize:
.quad 2048
shiftSize:
.quad 4096
varNameSize:
.quad 32
varSize:
.quad 128
typeSize:
.quad 32 
valSize:
.quad 64 
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
varName0:
.ascii "$ret"
.space 1, 0
lenVarName0 = . - varName0
varName1:
.ascii "$find_descr_return_var"
.space 1, 0
lenVarName1 = . - varName1
varName2:
.ascii "$find_descr_res"
.space 1, 0
lenVarName2 = . - varName2
label0:
 .quad .find_descr
labelName0:
.ascii ".find_descr"
.space 1,0
data0:
.ascii ""
.space 1, 0
lenData0 = . - data0
varName3:
.ascii "c"
.space 1, 0
lenVarName3 = . - varName3
varName4:
.ascii "b"
.space 1, 0
lenVarName4 = . - varName4
varName5:
.ascii "a"
.space 1, 0
lenVarName5 = . - varName5
data1:
.ascii ""
.space 1, 0
lenData1 = . - data1
data2:
.ascii ""
.space 1, 0
lenData2 = . - data2
data3:
.ascii "2.0"
.space 1, 0
lenData3 = . - data3
data4:
.ascii "4.0"
.space 1, 0
lenData4 = . - data4
label1:
 .quad .find_descr_end
labelName1:
.ascii ".find_descr_end"
.space 1,0
data5:
.ascii ""
.space 1, 0
lenData5 = . - data5
varName6:
.ascii "$solve_return_var"
.space 1, 0
lenVarName6 = . - varName6
label2:
 .quad .solve
labelName2:
.ascii ".solve"
.space 1,0
data6:
.ascii ""
.space 1, 0
lenData6 = . - data6
varName7:
.ascii "d"
.space 1, 0
lenVarName7 = . - varName7
varName8:
.ascii "c"
.space 1, 0
lenVarName8 = . - varName8
varName9:
.ascii "b"
.space 1, 0
lenVarName9 = . - varName9
varName10:
.ascii "a"
.space 1, 0
lenVarName10 = . - varName10
data7:
.ascii ""
.space 1, 0
lenData7 = . - data7
data8:
.ascii ""
.space 1, 0
lenData8 = . - data8
varName11:
.ascii "x1"
.space 1, 0
lenVarName11 = . - varName11
varName12:
.ascii "x2"
.space 1, 0
lenVarName12 = . - varName12
data9:
.ascii "0"
.space 1, 0
lenData9 = . - data9
data10:
.ascii ""
.space 1, 0
lenData10 = . - data10
varName13:
.ascii "$print_arg0"
.space 1, 0
lenVarName13 = . - varName13
data11:
.ascii "Нет решений\n"
.space 1, 0
lenData11 = . - data11
label3:
 .quad ._cond0_end
labelName3:
.ascii "._cond0_end"
.space 1,0
data12:
.ascii ""
.space 1, 0
lenData12 = . - data12
data13:
.ascii "-1.0"
.space 1, 0
lenData13 = . - data13
data14:
.ascii "0.5"
.space 1, 0
lenData14 = . - data14
data15:
.ascii "2.0"
.space 1, 0
lenData15 = . - data15
data16:
.ascii ""
.space 1, 0
lenData16 = . - data16
data17:
.ascii "-1.0"
.space 1, 0
lenData17 = . - data17
data18:
.ascii "0.5"
.space 1, 0
lenData18 = . - data18
data19:
.ascii "2.0"
.space 1, 0
lenData19 = . - data19
varName14:
.ascii "$print_arg0"
.space 1, 0
lenVarName14 = . - varName14
data20:
.ascii "Решение:\n"
.space 1, 0
lenData20 = . - data20
varName15:
.ascii "$s0"
.space 1, 0
lenVarName15 = . - varName15
varName16:
.ascii "$print_arg0"
.space 1, 0
lenVarName16 = . - varName16
data21:
.ascii "\n"
.space 1, 0
lenData21 = . - data21
varName17:
.ascii "$s0"
.space 1, 0
lenVarName17 = . - varName17
varName18:
.ascii "$print_arg0"
.space 1, 0
lenVarName18 = . - varName18
data22:
.ascii "\n"
.space 1, 0
lenData22 = . - data22
label4:
 .quad ._cond_exit0
labelName4:
.ascii "._cond_exit0"
.space 1,0
data23:
.ascii ""
.space 1, 0
lenData23 = . - data23
data24:
.ascii ""
.space 1, 0
lenData24 = . - data24
label5:
 .quad .solve_end
labelName5:
.ascii ".solve_end"
.space 1,0
data25:
.ascii ""
.space 1, 0
lenData25 = . - data25
varName19:
.ascii "$main_return_var"
.space 1, 0
lenVarName19 = . - varName19
label6:
 .quad .main
labelName6:
.ascii ".main"
.space 1,0
data26:
.ascii ""
.space 1, 0
lenData26 = . - data26
data27:
.ascii ""
.space 1, 0
lenData27 = . - data27
varName20:
.ascii "$print_arg0"
.space 1, 0
lenVarName20 = . - varName20
data28:
.ascii "Решение квадратных уравнений вида ax^2+bx+c=0\n"
.space 1, 0
lenData28 = . - data28
varName21:
.ascii "a"
.space 1, 0
lenVarName21 = . - varName21
varName22:
.ascii "b"
.space 1, 0
lenVarName22 = . - varName22
varName23:
.ascii "c"
.space 1, 0
lenVarName23 = . - varName23
varName24:
.ascii "d"
.space 1, 0
lenVarName24 = . - varName24
data29:
.ascii ""
.space 1, 0
lenData29 = . - data29
data30:
.ascii ""
.space 1, 0
lenData30 = . - data30
varName25:
.ascii "$for0"
.space 1, 0
lenVarName25 = . - varName25
label7:
 .quad ._for0
labelName7:
.ascii "._for0"
.space 1,0
data31:
.ascii ""
.space 1, 0
lenData31 = . - data31
data32:
.ascii ""
.space 1, 0
lenData32 = . - data32
varName26:
.ascii "$print_arg0"
.space 1, 0
lenVarName26 = . - varName26
data33:
.ascii "Введите a:\n"
.space 1, 0
lenData33 = . - data33
varName27:
.ascii "s"
.space 1, 0
lenVarName27 = . - varName27
data34:
.ascii ""
.space 1, 0
lenData34 = . - data34
data35:
.ascii "0"
.space 1, 0
lenData35 = . - data35
data36:
.ascii ""
.space 1, 0
lenData36 = . - data36
data37:
.ascii ""
.space 1, 0
lenData37 = . - data37
data38:
.ascii ""
.space 1, 0
lenData38 = . - data38
label8:
 .quad ._cond3_end
labelName8:
.ascii "._cond3_end"
.space 1,0
data39:
.ascii ""
.space 1, 0
lenData39 = . - data39
varName28:
.ascii "$F0"
.space 1, 0
lenVarName28 = . - varName28
label9:
 .quad ._attempt0
labelName9:
.ascii "._attempt0"
.space 1,0
data40:
.ascii ""
.space 1, 0
lenData40 = . - data40
data41:
.ascii "1"
.space 1, 0
lenData41 = . - data41
label10:
 .quad ._cond2_end
labelName10:
.ascii "._cond2_end"
.space 1,0
data42:
.ascii ""
.space 1, 0
lenData42 = . - data42
data43:
.ascii ""
.space 1, 0
lenData43 = . - data43
data44:
.ascii ""
.space 1, 0
lenData44 = . - data44
data45:
.ascii ""
.space 1, 0
lenData45 = . - data45
varName29:
.ascii "$print_arg0"
.space 1, 0
lenVarName29 = . - varName29
data46:
.ascii "Ошибка! Неверный формат числа\n"
.space 1, 0
lenData46 = . - data46
varName30:
.ascii "$print_arg0"
.space 1, 0
lenVarName30 = . - varName30
data47:
.ascii "\n"
.space 1, 0
lenData47 = . - data47
label11:
 .quad ._cond4_end
labelName11:
.ascii "._cond4_end"
.space 1,0
data48:
.ascii ""
.space 1, 0
lenData48 = . - data48
data49:
.ascii ""
.space 1, 0
lenData49 = . - data49
data50:
.ascii ""
.space 1, 0
lenData50 = . - data50
data51:
.ascii ""
.space 1, 0
lenData51 = . - data51
data52:
.ascii ""
.space 1, 0
lenData52 = . - data52
data53:
.ascii "1"
.space 1, 0
lenData53 = . - data53
label12:
 .quad ._cond5_end
labelName12:
.ascii "._cond5_end"
.space 1,0
data54:
.ascii ""
.space 1, 0
lenData54 = . - data54
data55:
.ascii ""
.space 1, 0
lenData55 = . - data55
label13:
 .quad ._undef_for0
labelName13:
.ascii "._undef_for0"
.space 1,0
data56:
.ascii ""
.space 1, 0
lenData56 = . - data56
label14:
 .quad ._cond1_end
labelName14:
.ascii "._cond1_end"
.space 1,0
data57:
.ascii ""
.space 1, 0
lenData57 = . - data57
data58:
.ascii ""
.space 1, 0
lenData58 = . - data58
label15:
 .quad ._cond6_end
labelName15:
.ascii "._cond6_end"
.space 1,0
data59:
.ascii ""
.space 1, 0
lenData59 = . - data59
data60:
.ascii ""
.space 1, 0
lenData60 = . - data60
label16:
 .quad ._cond7_end
labelName16:
.ascii "._cond7_end"
.space 1,0
data61:
.ascii ""
.space 1, 0
lenData61 = . - data61
label17:
 .quad ._for0_end
labelName17:
.ascii "._for0_end"
.space 1,0
data62:
.ascii ""
.space 1, 0
lenData62 = . - data62
data63:
.ascii ""
.space 1, 0
lenData63 = . - data63
data64:
.ascii ""
.space 1, 0
lenData64 = . - data64
data65:
.ascii ""
.space 1, 0
lenData65 = . - data65
varName31:
.ascii "$for1"
.space 1, 0
lenVarName31 = . - varName31
label18:
 .quad ._for1
labelName18:
.ascii "._for1"
.space 1,0
data66:
.ascii ""
.space 1, 0
lenData66 = . - data66
data67:
.ascii ""
.space 1, 0
lenData67 = . - data67
varName32:
.ascii "$print_arg0"
.space 1, 0
lenVarName32 = . - varName32
data68:
.ascii "Введите b:\n"
.space 1, 0
lenData68 = . - data68
varName33:
.ascii "s"
.space 1, 0
lenVarName33 = . - varName33
data69:
.ascii ""
.space 1, 0
lenData69 = . - data69
data70:
.ascii "0"
.space 1, 0
lenData70 = . - data70
data71:
.ascii ""
.space 1, 0
lenData71 = . - data71
data72:
.ascii ""
.space 1, 0
lenData72 = . - data72
data73:
.ascii ""
.space 1, 0
lenData73 = . - data73
label19:
 .quad ._cond10_end
labelName19:
.ascii "._cond10_end"
.space 1,0
data74:
.ascii ""
.space 1, 0
lenData74 = . - data74
varName34:
.ascii "$F0"
.space 1, 0
lenVarName34 = . - varName34
label20:
 .quad ._attempt1
labelName20:
.ascii "._attempt1"
.space 1,0
data75:
.ascii ""
.space 1, 0
lenData75 = . - data75
data76:
.ascii "1"
.space 1, 0
lenData76 = . - data76
label21:
 .quad ._cond9_end
labelName21:
.ascii "._cond9_end"
.space 1,0
data77:
.ascii ""
.space 1, 0
lenData77 = . - data77
data78:
.ascii ""
.space 1, 0
lenData78 = . - data78
data79:
.ascii ""
.space 1, 0
lenData79 = . - data79
data80:
.ascii ""
.space 1, 0
lenData80 = . - data80
varName35:
.ascii "$print_arg0"
.space 1, 0
lenVarName35 = . - varName35
data81:
.ascii "Ошибка! Неверный формат числа\n"
.space 1, 0
lenData81 = . - data81
varName36:
.ascii "$print_arg0"
.space 1, 0
lenVarName36 = . - varName36
data82:
.ascii "\n"
.space 1, 0
lenData82 = . - data82
label22:
 .quad ._cond11_end
labelName22:
.ascii "._cond11_end"
.space 1,0
data83:
.ascii ""
.space 1, 0
lenData83 = . - data83
data84:
.ascii ""
.space 1, 0
lenData84 = . - data84
data85:
.ascii ""
.space 1, 0
lenData85 = . - data85
data86:
.ascii ""
.space 1, 0
lenData86 = . - data86
data87:
.ascii ""
.space 1, 0
lenData87 = . - data87
data88:
.ascii "1"
.space 1, 0
lenData88 = . - data88
label23:
 .quad ._cond12_end
labelName23:
.ascii "._cond12_end"
.space 1,0
data89:
.ascii ""
.space 1, 0
lenData89 = . - data89
data90:
.ascii ""
.space 1, 0
lenData90 = . - data90
label24:
 .quad ._undef_for1
labelName24:
.ascii "._undef_for1"
.space 1,0
data91:
.ascii ""
.space 1, 0
lenData91 = . - data91
label25:
 .quad ._cond8_end
labelName25:
.ascii "._cond8_end"
.space 1,0
data92:
.ascii ""
.space 1, 0
lenData92 = . - data92
data93:
.ascii ""
.space 1, 0
lenData93 = . - data93
label26:
 .quad ._cond13_end
labelName26:
.ascii "._cond13_end"
.space 1,0
data94:
.ascii ""
.space 1, 0
lenData94 = . - data94
data95:
.ascii ""
.space 1, 0
lenData95 = . - data95
label27:
 .quad ._cond14_end
labelName27:
.ascii "._cond14_end"
.space 1,0
data96:
.ascii ""
.space 1, 0
lenData96 = . - data96
label28:
 .quad ._for1_end
labelName28:
.ascii "._for1_end"
.space 1,0
data97:
.ascii ""
.space 1, 0
lenData97 = . - data97
data98:
.ascii ""
.space 1, 0
lenData98 = . - data98
data99:
.ascii ""
.space 1, 0
lenData99 = . - data99
data100:
.ascii ""
.space 1, 0
lenData100 = . - data100
varName37:
.ascii "$for2"
.space 1, 0
lenVarName37 = . - varName37
label29:
 .quad ._for2
labelName29:
.ascii "._for2"
.space 1,0
data101:
.ascii ""
.space 1, 0
lenData101 = . - data101
data102:
.ascii ""
.space 1, 0
lenData102 = . - data102
varName38:
.ascii "$print_arg0"
.space 1, 0
lenVarName38 = . - varName38
data103:
.ascii "Введите c:\n"
.space 1, 0
lenData103 = . - data103
varName39:
.ascii "s"
.space 1, 0
lenVarName39 = . - varName39
data104:
.ascii ""
.space 1, 0
lenData104 = . - data104
data105:
.ascii "0"
.space 1, 0
lenData105 = . - data105
data106:
.ascii ""
.space 1, 0
lenData106 = . - data106
data107:
.ascii ""
.space 1, 0
lenData107 = . - data107
data108:
.ascii ""
.space 1, 0
lenData108 = . - data108
label30:
 .quad ._cond17_end
labelName30:
.ascii "._cond17_end"
.space 1,0
data109:
.ascii ""
.space 1, 0
lenData109 = . - data109
varName40:
.ascii "$F0"
.space 1, 0
lenVarName40 = . - varName40
label31:
 .quad ._attempt2
labelName31:
.ascii "._attempt2"
.space 1,0
data110:
.ascii ""
.space 1, 0
lenData110 = . - data110
data111:
.ascii "1"
.space 1, 0
lenData111 = . - data111
label32:
 .quad ._cond16_end
labelName32:
.ascii "._cond16_end"
.space 1,0
data112:
.ascii ""
.space 1, 0
lenData112 = . - data112
data113:
.ascii ""
.space 1, 0
lenData113 = . - data113
data114:
.ascii ""
.space 1, 0
lenData114 = . - data114
data115:
.ascii ""
.space 1, 0
lenData115 = . - data115
varName41:
.ascii "$print_arg0"
.space 1, 0
lenVarName41 = . - varName41
data116:
.ascii "Ошибка! Неверный формат числа\n"
.space 1, 0
lenData116 = . - data116
varName42:
.ascii "$print_arg0"
.space 1, 0
lenVarName42 = . - varName42
data117:
.ascii "\n"
.space 1, 0
lenData117 = . - data117
label33:
 .quad ._cond18_end
labelName33:
.ascii "._cond18_end"
.space 1,0
data118:
.ascii ""
.space 1, 0
lenData118 = . - data118
data119:
.ascii ""
.space 1, 0
lenData119 = . - data119
data120:
.ascii ""
.space 1, 0
lenData120 = . - data120
data121:
.ascii ""
.space 1, 0
lenData121 = . - data121
data122:
.ascii ""
.space 1, 0
lenData122 = . - data122
data123:
.ascii "1"
.space 1, 0
lenData123 = . - data123
label34:
 .quad ._cond19_end
labelName34:
.ascii "._cond19_end"
.space 1,0
data124:
.ascii ""
.space 1, 0
lenData124 = . - data124
data125:
.ascii ""
.space 1, 0
lenData125 = . - data125
label35:
 .quad ._undef_for2
labelName35:
.ascii "._undef_for2"
.space 1,0
data126:
.ascii ""
.space 1, 0
lenData126 = . - data126
label36:
 .quad ._cond15_end
labelName36:
.ascii "._cond15_end"
.space 1,0
data127:
.ascii ""
.space 1, 0
lenData127 = . - data127
data128:
.ascii ""
.space 1, 0
lenData128 = . - data128
label37:
 .quad ._cond20_end
labelName37:
.ascii "._cond20_end"
.space 1,0
data129:
.ascii ""
.space 1, 0
lenData129 = . - data129
data130:
.ascii ""
.space 1, 0
lenData130 = . - data130
label38:
 .quad ._cond21_end
labelName38:
.ascii "._cond21_end"
.space 1,0
data131:
.ascii ""
.space 1, 0
lenData131 = . - data131
label39:
 .quad ._for2_end
labelName39:
.ascii "._for2_end"
.space 1,0
data132:
.ascii ""
.space 1, 0
lenData132 = . - data132
data133:
.ascii ""
.space 1, 0
lenData133 = . - data133
varName43:
.ascii "$find_descr_res0"
.space 1, 0
lenVarName43 = . - varName43
data134:
.ascii "#find_descr_res0"
.space 1, 0
lenData134 = . - data134
label40:
 .quad .find_descr_res0
labelName40:
.ascii ".find_descr_res0"
.space 1,0
data135:
.ascii "#solve_res0"
.space 1, 0
lenData135 = . - data135
label41:
 .quad .solve_res0
labelName41:
.ascii ".solve_res0"
.space 1,0
data136:
.ascii ""
.space 1, 0
lenData136 = . - data136
data137:
.ascii ""
.space 1, 0
lenData137 = . - data137
label42:
 .quad .main_end
labelName42:
.ascii ".main_end"
.space 1,0
data138:
.ascii ""
.space 1, 0
lenData138 = . - data138
data139:
.ascii "#main_res0"
.space 1, 0
lenData139 = . - data139
label43:
 .quad .main_res0
labelName43:
.ascii ".main_res0"
.space 1,0
data140:
.ascii ""
.space 1, 0
lenData140 = . - data140
data141:
.ascii ""
.space 1, 0
lenData141 = . - data141
