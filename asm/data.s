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

divZeroError:
.ascii "fatal error: dividing by zero\n"
.space 1, 0

fatalError:
.ascii "fatal error: internal error\n"
.space 1, 0 

divINegError:
.ascii "@ is not defined for negative numbers\n"
.space 1, 0

powNegError:
.ascii "runtime error: ^ is not defined for negative base and fractional exponent\n"
.space 1, 0

powZeroNegError:
.ascii "runtime error: ^ is not defined for zero base and negative exponent\n"
.space 1, 0
powZeroZeroError:
.ascii "runtime error: ^ is not defined for zero base and zero exponent\n"
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
varName0:
.ascii "s"
lenVarName0 = . - varName0
data0:
.ascii "some cool string\n"
.space 1, 0
lenData0 = . - data0
varName1:
.ascii "s2"
lenVarName1 = . - varName1
varName2:
.ascii "s2"
lenVarName2 = . - varName2
varName3:
.ascii "s2"
lenVarName3 = . - varName3
varName4:
.ascii "s2"
lenVarName4 = . - varName4
varName5:
.ascii "s2"
lenVarName5 = . - varName5
varName6:
.ascii "s2"
lenVarName6 = . - varName6
varName7:
.ascii "s2"
lenVarName7 = . - varName7
varName8:
.ascii "s2"
lenVarName8 = . - varName8
varName9:
.ascii "s2"
lenVarName9 = . - varName9
varName10:
.ascii "s2"
lenVarName10 = . - varName10
varName11:
.ascii "s2"
lenVarName11 = . - varName11
varName12:
.ascii "s2"
lenVarName12 = . - varName12
varName13:
.ascii "s2"
lenVarName13 = . - varName13
varName14:
.ascii "s2"
lenVarName14 = . - varName14
varName15:
.ascii "s2"
lenVarName15 = . - varName15
varName16:
.ascii "s2"
lenVarName16 = . - varName16
varName17:
.ascii "s2"
lenVarName17 = . - varName17
varName18:
.ascii "s2"
lenVarName18 = . - varName18
varName19:
.ascii "s2"
lenVarName19 = . - varName19
varName20:
.ascii "s2"
lenVarName20 = . - varName20
varName21:
.ascii "s2"
lenVarName21 = . - varName21
varName22:
.ascii "s2"
lenVarName22 = . - varName22
varName23:
.ascii "s2"
lenVarName23 = . - varName23
varName24:
.ascii "s2"
lenVarName24 = . - varName24
varName25:
.ascii "s2"
lenVarName25 = . - varName25
varName26:
.ascii "s2"
lenVarName26 = . - varName26
varName27:
.ascii "s2"
lenVarName27 = . - varName27
varName28:
.ascii "s2"
lenVarName28 = . - varName28
varName29:
.ascii "s2"
lenVarName29 = . - varName29
varName30:
.ascii "s2"
lenVarName30 = . - varName30
varName31:
.ascii "s2"
lenVarName31 = . - varName31
varName32:
.ascii "s2"
lenVarName32 = . - varName32
varName33:
.ascii "s2"
lenVarName33 = . - varName33
varName34:
.ascii "s2"
lenVarName34 = . - varName34
varName35:
.ascii "s2"
lenVarName35 = . - varName35
varName36:
.ascii "s2"
lenVarName36 = . - varName36
varName37:
.ascii "s2"
lenVarName37 = . - varName37
varName38:
.ascii "s2"
lenVarName38 = . - varName38
varName39:
.ascii "s2"
lenVarName39 = . - varName39
varName40:
.ascii "s2"
lenVarName40 = . - varName40
varName41:
.ascii "s2"
lenVarName41 = . - varName41
varName42:
.ascii "s2"
lenVarName42 = . - varName42
varName43:
.ascii "s2"
lenVarName43 = . - varName43
varName44:
.ascii "s2"
lenVarName44 = . - varName44
varName45:
.ascii "s2"
lenVarName45 = . - varName45
varName46:
.ascii "s2"
lenVarName46 = . - varName46
varName47:
.ascii "s2"
lenVarName47 = . - varName47
varName48:
.ascii "s2"
lenVarName48 = . - varName48
varName49:
.ascii "s2"
lenVarName49 = . - varName49
varName50:
.ascii "s2"
lenVarName50 = . - varName50
varName51:
.ascii "s2"
lenVarName51 = . - varName51
varName52:
.ascii "s2"
lenVarName52 = . - varName52
varName53:
.ascii "s2"
lenVarName53 = . - varName53
varName54:
.ascii "s2"
lenVarName54 = . - varName54
varName55:
.ascii "s2"
lenVarName55 = . - varName55
varName56:
.ascii "s2"
lenVarName56 = . - varName56
varName57:
.ascii "s2"
lenVarName57 = . - varName57
varName58:
.ascii "s2"
lenVarName58 = . - varName58
varName59:
.ascii "s2"
lenVarName59 = . - varName59
varName60:
.ascii "s2"
lenVarName60 = . - varName60
varName61:
.ascii "s2"
lenVarName61 = . - varName61
varName62:
.ascii "s2"
lenVarName62 = . - varName62
varName63:
.ascii "s2"
lenVarName63 = . - varName63
varName64:
.ascii "s2"
lenVarName64 = . - varName64
varName65:
.ascii "s2"
lenVarName65 = . - varName65
varName66:
.ascii "s2"
lenVarName66 = . - varName66
varName67:
.ascii "s2"
lenVarName67 = . - varName67
varName68:
.ascii "s2"
lenVarName68 = . - varName68
varName69:
.ascii "s2"
lenVarName69 = . - varName69
varName70:
.ascii "s2"
lenVarName70 = . - varName70
varName71:
.ascii "s2"
lenVarName71 = . - varName71
varName72:
.ascii "s2"
lenVarName72 = . - varName72
varName73:
.ascii "s2"
lenVarName73 = . - varName73
varName74:
.ascii "s2"
lenVarName74 = . - varName74
varName75:
.ascii "s2"
lenVarName75 = . - varName75
varName76:
.ascii "s2"
lenVarName76 = . - varName76
varName77:
.ascii "s2"
lenVarName77 = . - varName77
varName78:
.ascii "s2"
lenVarName78 = . - varName78
varName79:
.ascii "s2"
lenVarName79 = . - varName79
varName80:
.ascii "s2"
lenVarName80 = . - varName80
