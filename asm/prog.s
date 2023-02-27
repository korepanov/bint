.data
pageSize:
.quad 4096 
varNameSize:
.quad 32
varSize:
.quad 160 
typeSize:
.quad 32 
metaSize:
.quad 32 
valSize:
.quad 64 
buf:
.quad 0, 0, 0, 0
lenBuf = . - buf 
buf2:
.quad 0, 0, 0, 0
lenBuf2 = . - buf2 
varType:
.quad 0, 0, 0, 0
lenVarType = . - varType 
varName:
.quad 0, 0, 0, 0
lenVarName = . - varName 
userData:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserData = . - userData
varName0:
.ascii "iVar"
lenVarName0 = . - varName0
varType0:
.ascii "int"
lenVarType0 = . - varType0
varName1:
.ascii "sVar"
lenVarName1 = . - varName1
varType1:
.ascii "string"
lenVarType1 = . - varType1
varName2:
.ascii "fVar"
lenVarName2 = . - varName2
varType2:
.ascii "float"
lenVarType2 = . - varType2
initInt:
.ascii "0"
.space 1, 0
lenInitInt = . - initInt
initFloat:
.ascii "0.0"
.space 1, 0
lenInitFloat = . - initFloat 
initBool:
.ascii "0"
.space 1, 0
lenIniBool = . - initBool 
initString:
.ascii ""
.space 1, 0
lenInitString = . - initString 
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
data0:
.ascii "25"
.space 1, 0
lenData0 = . - data0 
data1:
.ascii "Привет, мир!"
.space 1, 0
lenData1 = . - data1 

fatalError:
.ascii "fatal error: internal error\n"
.space 1, 0 

.text	

__throughError:
 mov $fatalError, %rsi
 call __print 
 mov $60, %rax
 mov $1, %rdi
 syscall

__print:
 mov (%rsi), %al	
 cmp $0, %al	
 jz  __printEx			
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall		    
 inc %rsi		  		    
 jnz __print
__printEx:
 ret

 __printSymbol:
 mov (%rsi), %al				
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall		    
 ret

__printHeap:  
 mov %r13, %r8  
 __printHeapLoop:
 cmp %r15, %r8 
 jz __printHeapEx
 mov %r8, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 __printHeapEx:
 ret 


__toStr:
 # число в %rax 
 # подготовка преобразования числа в строку
  movq $0, (buf2)
  mov $10, %r8    # делитель
  mov $buf, %rsi  # адрес начала буфера 
  xor %rdi, %rdi  # обнуляем счетчик
# преобразуем путем последовательного деления на 10
__toStrlo:
  xor %rdx, %rdx  # число в rdx:rax
  div %r8         # делим rdx:rax на r8
  add $48, %dl    # цифру в символ цифры
  mov %dl, (%rsi) # в буфер
  inc %rsi        # на следующую позицию в буфере
  inc %rdi        # счетчик увеличиваем на 1
  cmp $0, %rax    # проверка конца алгоритма
  jnz __toStrlo          # продолжим цикл?
# число записано в обратном порядке,
# вернем правильный, перенеся в другой буфер 
  mov $buf2, %rbx # начало нового буфера
  mov $buf, %rcx  # старый буфер
  add %rdi, %rcx  # в конец
  dec %rcx        # старого буфера
  mov %rdi, %rdx  # длина буфера
# перенос из одного буфера в другой
__toStrexc:
  mov (%rcx), %al # из старого буфера
  mov %al, (%rbx) # в новый
  dec %rcx        # в обратном порядке  
  inc %rbx        # продвигаемся в новом буфере
  dec %rdx        # а в старом в обратном порядке
  jnz __toStrexc         # проверка конца алгоритма
  ret

__set: #set strings
 # входные параметры 
 # rsi - длина буфера назначения 
 # rdx - адрес буфера назначения
 # rax - длина буфера источника 
 # rdi - адрес буфера источника 
 mov %rdx, %r8 
 mov %rsi, %r9
  
 __setClear:
 movb $'*', (%rdx)
 dec %rsi
 inc %rdx
 cmp $0, %rsi  
 jnz __setClear
 dec %rdx  
 movb $0, (%rdx)

 mov %r8, %rdx   
 
 __setLocal:
 mov (%rdi), %r11b
 movb %r11b, (%rdx)
 inc %rdx
 inc %rdi
 dec %rax  
 cmp $0, %rax
 jnz __setLocal
 movb $0, (%rdx)
 ret 

__toNumber:
 # вход: buf 
 # выход:  %rax 
 mov $buf, %rdx # our string
 __toNumberAtoi:
 xor %rax, %rax # zero a "result so far"
 __toNumberTop:
 movzx (%rdx), %rcx # get a character
 inc %rdx # ready for next one
 cmp $0, %rcx # end?
 jz __toNumberDone
 sub $48, %rcx # "convert" character to number
 imul $10, %rax # multiply "result so far" by ten
 add %rcx, %rax # add in current digit
 jmp __toNumberTop # until done
 __toNumberDone:
 ret


__defineVar:
 # адрес имени переменной в $varName
 # адрес типа переменой в $varType
 
 mov $varName, %rcx 
 mov $varType, %rdx 
 mov %r14, %rax

 cmp %rax, %r15
 jg __defOk 
 mov %r15, %r8
 call __newMem
 mov $varName, %rcx 
 mov $varType, %rdx
 __defOk:
 mov %r14, %r8 
 __defOkLocal:
 movb (%rcx), %r11b 
 cmp $'*', %r11b
 jz __defOkLocalEx
 movb %r11b, (%r8)
 inc %rcx 
 inc %r8 
 jmp __defOkLocal
 __defOkLocalEx:
 mov %r14, %r8 
 add (varNameSize), %r8 
 __defOkTypeLocal:
 movb (%rdx), %r11b
 cmp $'*', %r11b 
 jz __defOkTypeLocalEx
 movb %r11b, (%r8)
 inc %rdx
 inc %r8 
 jmp __defOkTypeLocal
 __defOkTypeLocalEx:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 
 # rsi - длина буфера назначения 
 # rdx - адрес буфера назначения
 # rax - длина буфера источника 
 # rdi - адрес буфера источника 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarType, %rax 
 mov $varType, %rdi 
 call __set 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
  
 call __compare 
 cmp $1, %rax 
 jz __defInt 

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi 
 call __set 
 call __compare
 cmp $1, %rax 
 jz __defFloat 

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi 
 call __set 
 call __compare
 cmp $1, %rax  
 jz __defBool 

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set 
 call __compare
 cmp $1, %rax 
 jz __defString 
 call __throughError

 __defInt:
 mov $intType, %rsi
 call __print  
 jmp __defEnd
 __defFloat:
 mov $floatType, %rsi
 call __print 
 jmp __defEnd
 __defBool:
 mov $boolType, %rsi
 call __print 
 jmp __defEnd
 __defString:
 mov $stringType, %rsi
 call __print 

 __defEnd:

 mov %r14, %rax 
 add (varSize), %rax 
 mov %rax, %r14
 ret 

# r12 - pointer (общего назначения)
# r13 - heapBegin 
# r14 - heapPointer 
# r15 - heapMax 

__firstMem:
 # получить адрес начала области для выделения памяти
 mov $12, %rax
 xor %rdi, %rdi
 syscall
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, %r13
 mov %rax, %r14
 mov %rax, %r9 
 add (pageSize), %r9
 mov %r9, %r15 
# выделить динамическую память
 mov (pageSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $'*', %dl
 mov $0, %rbx
 __lo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (pageSize), %rbx
 jz  __ex
 jmp __lo
 __ex:
 ret 

 __newMem:
 # адрес начала выделяемой памяти в  %r8 
# запомнить адрес начала выделяемой памяти
 #mov %r8, %r14
 mov %r8, %r9 
 add (pageSize), %r9 
 mov %r9, %r15
# выделить динамическую память
 mov (pageSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $'*', %dl
 mov $0, %rbx
 __newMemlo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (pageSize), %rbx
 jz  __newMemEx
 jmp __newMemlo
 __newMemEx:
 ret 

__read: 
 # считать в buf по указателю в %r12 
 mov %r12, %r8
 mov $buf, %r10 
 mov $lenBuf, %rsi 
__readClear:
 movb $0, (%r10)
 dec %rsi
 inc %r10
 cmp $0, %rsi  
 jnz __readClear
 mov $buf, %r10
 __readLocal: 
 movb (%r8), %r9b
 cmp $'*', %r9b  
 jz __readEx
 mov %r9b, (%r10)
 inc %r10 
 inc %r8 
 jmp __readLocal
 __readEx:
 mov $buf, %r10 
 cmp $0, (%r10)
 jnz __readOk
 movb $'*', (%r10)
 __readOk:
 ret 

__compare:
 # сравнить строки по адресу $buf и $varName  
 mov $buf, %rax 
 mov $varName, %rbx 
 __compareLocal:
 movb (%rax), %dl 
 cmp $0, %dl 
 jz __equal
 movb (%rax), %dl
 cmp %dl, (%rbx) 
 jnz __notEqual
 inc %rax 
 inc %rbx 
 jmp __compareLocal

 __notEqual:
 mov $0, %rax 
 ret 
 __equal:
 mov $1, %rax  
 ret 

__setVar:
 # имя переменной по адресу $varName 
 # данные в $userData 
 mov %r13, %rbx
 __setVarLocal:
 cmp %r15, %rbx
 jg __setVarEnd

 mov %rbx, %r12 
 call __read 
 cmp $'*', (buf)
 jz __setVarEnd 
 
 add (varSize), %rbx 
 jmp __setVarLocal  
  
 __setVarEnd:
 __setVarSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $'*', (buf)
 jz __throughError
 mov $buf, %rsi 
 mov %rbx, %r12 
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __setVarSearch
 
 add (varNameSize), %rbx 
 add (typeSize), %rbx
 mov %rbx, %r10 # сохраняем значение %rbx  
 mov $userData, %rax 
 xor %rdi, %rdi # счетчик количества реально записанных байт 
 __setNow:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __setMeta  
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 inc %rdi 
 jmp __setNow 
 __setMeta: 
 mov %r10, %rbx 
 add (valSize), %rbx
 mov %rbx, %r10 
 mov %rdi, %rax 
 call __toStr 
 mov $buf2, %rax
 mov %r10, %rbx 
 
 __setMetaLocal: 
 mov (%rax), %dl 
 mov %dl, (%rbx)
 cmp $0, %dl 
 jz __setVarRet 
 inc %rbx 
 inc %rax 

 jmp __setMetaLocal

 __setVarRet:
 ret

 __getVar:
 # вход: имя переменной по адресу $varName 
 # выход: $userData 
 mov %r13, %rbx
 __getVarLocal:
 cmp %r15, %rbx
 jg __getVarEnd

 mov %rbx, %r12 
 call __read 
 cmp $'*', (buf)
 jz __getVarEnd 
 
 add (varSize), %rbx 
 jmp __getVarLocal  
  
 __getVarEnd:
 __getVarSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $'*', (buf)
 jz __throughError
 mov $buf, %rsi 
 mov %rbx, %r12 
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __getVarSearch
 
 add (varNameSize), %rbx 
 add (typeSize), %rbx 
 mov $userData, %rax 
 mov $lenUserData, %r12 

 mov $userData, %r10 
 mov $lenUserData, %rsi 
 __getClear:
 movb $0, (%r10)
 dec %rsi
 inc %r10
 cmp $0, %rsi  
 jnz __getClear

 mov $buf, %r10 
 mov $lenBuf, %rsi 
 __getClearBuf:
 movb $0, (%r10)
 dec %rsi
 inc %r10
 cmp $0, %rsi  
 jnz __getClearBuf 
 
 __getMeta:

 mov $buf, %rsi 
 mov %rbx, %r10 
 add (valSize), %rbx 
 mov $1, %dl  
 __getMetaLocal:
 cmp $'*', %dl 
 jz __getNow 
 mov (%rbx), %dl
 mov %dl, (%rsi)  
 inc %rbx 
 inc %rsi  
 jmp __getMetaLocal
 
 __getNow:
 call __toNumber
 mov %r10, %rbx
 mov $userData, %rsi
 __getNowLocal:  
 cmp $0, %rax 
 jz __getVarRet
 mov (%rbx), %dl 
 mov %dl, (%rsi)
 inc %rsi 
 inc %rbx  
 dec %rax 
 jmp __getNowLocal 
 __getVarRet:
 ret




.globl _start
_start:
 call __firstMem

 #iVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName0, %rax 
 mov $varName0, %rdi 
 call __set 

 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenVarType0, %rax 
 mov $varType0, %rdi 
 call __set 

 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 
 #sVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi 
 call __set 

 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenVarType1, %rax 
 mov $varType1, %rdi 
 call __set

 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 
#fVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi 
 call __set 

 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenVarType2, %rax 
 mov $varType2, %rdi 
 call __set

 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi 
 call __set 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenData1, %rax 
 mov $data1, %rdi 
 call __set
 call __setVar 

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName0, %rax 
 mov $varName0, %rdi 
 call __set 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenData0, %rax 
 mov $data0, %rdi 
 call __set
 call __setVar

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi 
 call __set
 #call __getVar 
 #mov $userData, %rsi 
 #call __print 
 #call __printHeap

__stop:
 mov $60,  %rax      # номер системного вызова exit
 xor %rdi, %rdi      # код возврата (0 - выход без ошибок)
 syscall   
