.text	

__throughError:
 mov $fatalError, %rsi
 call __print 
 mov $60, %rax
 mov $1, %rdi
 syscall

 __throughUserError:
 # %rsi - адрес, по которому лежит сообщение об ошибке 
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

# посчитать количество символов до 0
# адрес начала в %rsi
# результат в %rax  
 __len:
 xor %rax, %rax 
 __lenLocal:
 mov (%rsi), %dl	
 cmp $0, %dl	
 jz  __lenEx				    
 inc %rsi		  	
 inc %rax 	    
 jmp __lenLocal  
__lenEx:
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
 cmp (strMax), %r8 
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
  cmp $0, %rax 
  jg __toStrPos
  mov $0, %rdx 
  sub $1, %rdx 
  imul %rdx  
  movb $'-', (buf2)
  jmp __toStrNeg 
  __toStrPos:
  movq $0, (buf2)
  __toStrNeg:
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
  mov (%rbx), %dl 
  cmp $'-', %dl
  jnz __toStrEmpty
  inc %rbx 
  __toStrEmpty:
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
  movb $0, (%rbx)
  mov $buf2, %rbx 
  mov (%rbx), %dl 
  cmp $'-', %dl
  jnz __toStrEnd 
  inc %rbx 
  mov (%rbx), %dl 
  cmp $'0', %dl
  jnz __toStrEnd 
  inc %rbx 
  mov (%rbx), %dl 
  cmp $0, %dl 
  jnz __toStrEnd
  mov $buf2, %rbx 
  movb $'0', (%rbx) 
  inc %rbx 
  movb $0, (%rbx) 
__toStrEnd:
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
 cmp $0, %rsi
 jz __setClearEnd
 movb $'*', (%rdx)
 dec %rsi
 inc %rdx  
 jmp __setClear
 __setClearEnd:
 dec %rdx 
 movb $0, (%rdx)

 mov %r8, %rdx  
 mov %r9, %rsi  
 
 __setLocal:
 cmp $0, %rax 
 jz __setLocalEnd
 cmp $2, %rsi
 jz __setLocalEnd
 mov (%rdi), %r11b 
 movb %r11b, (%rdx)
 inc %rdx
 inc %rdi
 dec %rax  
 dec %rsi 
 jmp __setLocal
 __setLocalEnd:
 dec %rdx 
 mov (%rdx), %rax 
 cmp $'*', %rax 
 jz __star
 inc %rdx
 __star: 
 movb $0, (%rdx)
 ret 

__concatinate:
 # входные параметры 
 # r8 - длина буфера первого операнда 
 # r9 - адрес буфера первого операнда
 # r11 - адрес буфера второго операнда 
 # выход
 # userData   
 call __clearUserData
 mov $lenUserData, %rsi # присваиваем в userData первый операнд  
 mov $userData, %rdx 
 mov %r8, %rax 
 mov %r9, %rdi 

 mov %r8, %rbx 
 mov %r11, %rcx
 call __set
 mov %rbx, %r8 
 mov %rcx, %r11 

 mov $userData, %r8
 __concNext:
 mov (%r8), %dl 
 cmp $0, %dl 
 jz __concLocal 
 inc %r8   
 jmp __concNext 
  
 __concLocal:
 
 mov (%r11), %dl 
 movb %dl, (%r8)
 inc %r8 
 inc %r11
 #dec %r10  
 cmp $0, %dl 
 jnz __concLocal
 #movb $0, (%r8)
 
 ret 

__toNumber:
 # вход: buf 
 # выход:  %rax 
 mov $buf, %rdx # our string
 movzx (%rdx), %rcx 
 cmp $'-', %rcx 
 jnz __toNumberAtoi
 inc %rdx 
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
 mov $buf, %rdx 
 movzx (%rdx), %rcx 
 cmp $'-', %rcx 
 jnz __toNumberIsPos
 mov $0, %rdx 
 sub $1, %rdx 
 //mov %rdx, (buf)
 imul %rdx 
 __toNumberIsPos:
 ret


__defineVar:
 # адрес имени переменной в $varName
 # адрес типа переменой в $varType
 
 mov $varName, %rcx 
 mov $varType, %rdx 
 cmp %r14, %r15
 jg __defOk 
 #mov %r15, %r8
 mov (strMax), %r8 
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
 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
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
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 movb $'0',(%r8)
 inc %r8 
 movb $0, (%r8)
 jmp __defEnd
 __defFloat:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 movb $'0',(%r8)
 inc %r8 
 movb $'.', (%r8)
 inc %r8 
 movb $'0', (%r8)
 inc %r8 
 movb $0, (%r8)
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8 
 jmp __defEnd
 __defBool:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 movb $'0',(%r8)
 inc %r8 
 movb $0, (%r8) 
 jmp __defEnd
 __defString:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 mov (strPointer), %rax 
 movb $0, (%rax)
 mov %r8, %r12 
 call __toStr 
 mov %r12, %r8 
 mov $buf2, %rax
 __defAddr:
 mov (%rax), %dl 
 cmp $0, %dl  
 jz __defStrEnd 
 mov %dl, (%r8)
 inc %rax 
 inc %r8 
 jmp __defAddr
 __defStrEnd:
 movb $0, (%r8)
 mov (strPointer), %rax 
 add (valSize), %rax 
 #cmp (strMax), %rax 
 #jg __defStrNewMem
 mov %rax, (strPointer)
 #jmp __defEnd 
 #__defStrNewMem:
 #mov %rax, (strPointer)
 #call __newStrMem
 __defEnd:

 add (varSize), %r14
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

 __firstStrMem:
 # адрес начала области для выделения памяти
 mov %r15, %rax
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, (strBegin)
 mov %rax, (strPointer)
 mov %rax, %r9 
 add (pageSize), %r9
 mov %r9, (strMax)
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
 __firstStrMemLo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (pageSize), %rbx
 jz  __firstStrMemEx
 jmp __firstStrMemLo
 __firstStrMemEx:
 ret 


 __newMem:
 # адрес начала выделяемой памяти в  %r8 
# запомнить адрес начала выделяемой памяти
 #mov %r8, %r14 
 mov %r8, %r9 
 add (pageSize), %r9 
 #mov %r9, %r15
 mov %r15, (oldHeapMax)
 #mov (strPointer), %r15 
 add (pageSize), %r15 
# выделить динамическую память
 mov (pageSize), %rdi
 #add %rax, %rdi
 add %r8, %rdi 
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
 mov %r15, %r8 
 call __newStrMem
 call __shiftStr
 ret  

 __newStrMem:
 # адрес начала выделяемой памяти в  %r8 
# запомнить адрес начала выделяемой памяти
 #mov %r8, %r14
 mov %r8, %r9 
 #add (pageSize), %r9
 
 mov (strPageNumber), %rax
 mov (pageSize), %rdi
 __newStrMemOkBegin: 
 cmp $0, %rax 
 jz __newStrMemOk
 dec %rax 
 add (pageSize), %rdi  
 jmp __newStrMemOkBegin
 __newStrMemOk:
 mov %rdi, (memorySize)
 add (memorySize), %r9 
 mov (strPageNumber), %rax 
 inc %rax 
 mov %rax, (strPageNumber)
 #mov %r9, (strMax)
# выделить динамическую память
 mov (memorySize), %rdi
 #add %rax, %rdi
 add %r8, %rdi 
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $'*', %dl
 mov $0, %rbx
 __newStrMemlo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (memorySize), %rbx
 jz  __newStrMemEx
 jmp __newStrMemlo
 __newStrMemEx:
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
 
 __renewStr:
 # адрес начала кучи 
 mov %r13, %r12
 # старый адрес конца кучи 
 mov (oldHeapMax), %r11 
 add (varNameSize), %r12  
 
 __renewFindStr:
 call __read
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set
 call __compare 
 cmp $1, %rax 
 jz __renewVal
 add (varSize), %r12 
 jmp __renewFindStr

 __renewVal:
 add (typeSize), %r12 
 call __read 
 call __toNumber
 __renewValLocal:
 mov (%rax), %r10b 
 cmp $0, %r10b 
 jz __renewValEnd
 movb $'*', (%rax)
 inc %rax 
 jmp __renewValLocal 
 __renewValEnd: 
 movb $'*', (%rax)
 __renewAddr:
 call __read 
 mov $buf, %rsi
 call __toNumber 
 #mov (pageSize), %rax 
 #call __toStr
 add (pageSize), %rax 
 call __toStr # в buf2 новый адрес строки 

 mov %r12, %rsi 
 __renewAddrLocal:
 mov (%rsi), %r10b 
 cmp $'*', %r10b 
 jz __renewAddrEnd
 movb $'*', (%rsi)
 inc %rsi 
 jmp __renewAddrLocal
 __renewAddrEnd:
 mov %r12, %rsi 
 mov $buf2, %rdx 
 // запись нового адреса
 __renewStrAddr: 
 mov (%rdx), %r10b 
 cmp $0, %r10b 
 jz __renewStrAddrEnd
 movb %r10b, (%rsi)
 inc %rsi 
 inc %rdx 
 jmp __renewStrAddr
 __renewStrAddrEnd:
 movb $0, (%rsi)
 ret 

 __shiftStr:
 # формируем в %r10 адрес нового начала
 mov (strBegin), %r10 
 add (pageSize), %r10
 #mov %r10, %r12 
 #add (pageSize), %r12 
 #mov %r12, (strMax)
 # адрес старого начала
 mov (strBegin), %r11
 __shiftMake: 
 mov (strMax), %r9
 cmp %r9, %r11   
 jz __shiftMakeEnd
 movb (%r11), %r12b 
 movb %r12b, (%r10)
 inc %r10
 inc %r11 
 jmp __shiftMake
 __shiftMakeEnd:
 mov (strPointer), %r10 
 add (pageSize), %r10 
 mov %r10, (strPointer)

 mov (strBegin), %r10 
 add (pageSize), %r10 
 mov %r10, (strBegin)

 mov (strMax), %r10 
 add (pageSize), %r10 
 mov %r10, (strMax)
 
 call __renewStr
 ret 

__compare:
 # сравнить строки по адресу $buf и $buf2
 # если длины строк не равны, то строки не равны 
 mov $buf, %rsi 
 call __len 
 mov %rax, %rbx 
 mov $buf2, %rsi 
 call __len 
 cmp %rax, %rbx 
 jnz __notEqual

 mov $buf, %rax 
 mov $buf2, %rbx 
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
 # вход: 
 # имя переменной по адресу $varName 
 # данные по указателю в (userData) 
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
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __setVarSearch
 

 add (varNameSize), %rbx 
 mov %rbx, %rax 
 mov %rbx, %r12 
 call __read  
 add (typeSize), %rbx 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set
 mov %rbx, %r12 
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jnz __setVarStr
 
 #add (varNameSize), %rbx 
 #add (typeSize), %rbx
 mov %rbx, %r8 
 mov %rbx, %rax 
 __setVarClear:
 add (valSize), %rax 
 __setVarClearLocal: 
 cmp %rax, %rbx 
 jz __setVarClearEnd
 movb $'*', (%rbx) 
 inc %rbx 
 jmp __setVarClearLocal
 __setVarClearEnd:
 mov %r8, %rbx 
 jmp __setVarIsNotStr

 __setVarStr:
 mov %rbx, %r12 
 call __read 
 call __toNumber
 mov %rax, %rbx  

 __setVarIsNotStr:
 
 #call __throughError
 mov %rbx, %r10 # сохраняем значение %rbx  
 mov (userData), %rax 
 xor %rdi, %rdi # счетчик количества реально записанных байт 
 __setNow:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __setVarRet 
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 inc %rdi 
 jmp __setNow 

 __setVarRet:
 movb $0, (%rbx) 
 ret

 __getVar:
 # вход: имя переменной по адресу $varName 
 # выход: указатель на данные в (userData)
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
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __getVarSearch
 
 add (varNameSize), %rbx 
 mov %rbx, %r12 
 call __read 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set
 call __compare 
 cmp $1, %rax 
 jnz __getVarNotStr
 mov $1, %r11
 jmp __getVarNotStrEnd
 __getVarNotStr:
 mov $0, %r11 
 __getVarNotStrEnd:
 mov %r12, %rbx 
 add (typeSize), %rbx  
 
 __getNow:
 #call __toNumber
 cmp $1, %r11
 jz __getVarGetStr
 mov %rbx, (userData)
 ret 
 __getVarGetStr:
 mov %rbx, %r12 
 call __read
 call __toNumber
 mov %rax, (userData) 
 ret 
 /*mov $userData, %rsi
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
 ret*/

__clearBuf:
mov $buf, %rsi 
mov $lenBuf, %rdi
__clearBufLocal: 
cmp $1, %rdi 
jz __clearBufEnd
movb $'*', (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal

__clearBufEnd:
movb $0, (%rsi)
ret

__clearBuf2:
mov $buf2, %rsi 
mov $lenBuf2, %rdi
__clearBufLocal2: 
cmp $1, %rdi 
jz __clearBufEnd2
movb $'*', (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal2

__clearBufEnd2:
movb $0, (%rsi)
ret

__clearUserData:
mov $userData, %rsi 
mov $lenUserData, %rdi
__clearUserDataLocal: 
cmp $1, %rdi 
jz __clearUserDataEnd
movb $'*', (%rsi)
inc %rsi 
dec %rdi 
jmp __clearUserDataLocal

__clearUserDataEnd:
movb $0, (%rsi)
ret

__clearBuf3:
mov $buf3, %rsi 
mov $lenBuf3, %rdi
__clearBufLocal3: 
cmp $1, %rdi 
jz __clearBufEnd3
movb $'*', (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal3

__clearBufEnd3:
movb $0, (%rsi)
ret

__clearBuf4:
mov $buf4, %rsi 
mov $lenBuf4, %rdi
__clearBufLocal4: 
cmp $1, %rdi 
jz __clearBufEnd4
movb $'*', (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal4

__clearBufEnd4:
movb $0, (%rsi)
ret

__add:
 # вход: buf и buf2
 # %rax - тип операции 
 # 0 - целочисленное сложение 
 # 1 - сложение вещественных чисел   
 # выход: userData 
 call __clearUserData
 cmp $0, %rax 
 jz __addInt 
 cmp $1, %rax 
 jz __addFloat

 call __throughError

 __addInt:
 call __toNumber
 mov %rax, %rbx 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber
 add %rbx, %rax
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 ret 
 __addFloat:
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, (buf)
 fld (buf)
 movss %xmm1, (buf)
 fadd (buf)
 fstp (buf)
 call __floatToStr

 ret 


 __sub:
 # вход: buf и buf2
 # %rax - тип операции 
 # 0 - целочисленное сложение 
 # 1 - сложение вещественных чисел   
 # выход: userData 
 call __clearUserData
 cmp $0, %rax 
 jz __subInt 
 cmp $1, %rax 
 jz __subFloat

 call __throughError

 __subInt:
 call __toNumber
 mov %rax, %rbx 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber
 sub %rax, %rbx
 mov %rbx, %rax 
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 ret 
 __subFloat:
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm1, (buf)
 fld (buf)
 movss %xmm0, (buf)
 fsub (buf)
 fstp (buf)
 call __floatToStr

 ret 

__mul:
 # вход: buf и buf2
 # %rax - тип операции 
 # 0 - целочисленное сложение 
 # 1 - сложение вещественных чисел   
 # выход: userData 
 call __clearUserData
 cmp $0, %rax 
 jz __mulInt 
 cmp $1, %rax 
 jz __mulFloat

 call __throughError

 __mulInt:
 call __toNumber
 mov %rax, %rbx 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber
 imul %rbx, %rax
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 ret 
 __mulFloat:
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, (buf)
 fld (buf)
 movss %xmm1, (buf)
 fmul (buf)
 fstp (buf)
 call __floatToStr

 ret 

__div:
 # вход: buf и buf2
 # только для вещественных чисел!   
 # выход: userData 
 call __clearUserData
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 # проверка деления на нуль
 movss (zero), %xmm2 
 cmpss $0, %xmm0, %xmm2
 pextrb $3, %xmm2, %rax
 cmp $0, %rax 
 jnz __divIsZero 
 movss %xmm1, (buf)
 fld (buf)
 movss %xmm0, (buf)
 fdiv (buf)
 fstp (buf)
 call __floatToStr
 ret 
 __divIsZero:
 mov $divZeroError, %rsi 
 call __throughUserError

__divI:
 # вход: buf и buf2
 # только для неотрицательных целых чисел!   
 # выход: userData 
 call __clearUserData
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber 
 cmp $0, %rax 
 jl __divINeg
 mov %rax, %r10 # первый операнд сохранен в %r10
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __toNumber # второй операнд в %rax 
 # проверка деления на нуль  
 cmp $0, %rax 
 jz __divIsZeroI 
 jl __divINeg
 mov %rax, %rcx # запоминаем второй операнд в %rcx 
 mov %r10, %rax # первый операнд в %rax
 xor %rdx, %rdx  
 div %rcx 
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 ret 
 __divIsZeroI:
 mov $divZeroError, %rsi 
 call __throughUserError
__divINeg:
 mov $divINegError, %rsi 
 call __throughUserError

__pow:
 # вход: buf и buf2
 # только для вещественных чисел!   
 # выход: userData 
 movb $0, (isExpNeg) # признак неотрицательного результата
 call __clearUserData
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm1, (buf)
 movss %xmm0, (buf4)
 # основание - нуль? 
 movss (zero), %xmm2 
 movss %xmm1, %xmm3 
 cmpss $0, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jnz __powZeroBase 

 movss (zero), %xmm2 
 movss (buf), %xmm3 
 cmpss $1, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jz __powIsPos
 movss (buf4), %xmm2 
 roundps $3, %xmm2, %xmm2
 movss %xmm2, (buf3)
 fld (buf4)
 fsub (buf3)
 fstp (buf3)

 movss (zero), %xmm2 
 movss (buf3), %xmm3 
 cmpss $0, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jz __powNotInt
 cvtss2si (buf4), %rax 
 mov $2, %rbx 
 xor %rdx, %rdx 
 div %rbx 
 cmp $0, %dl # число четное?
 jnz __powNotOdd
 movb $0, (isExpNeg) # признак неотрицательного результата
 jmp __powNotOddEnd
 __powNotOdd:   
 movb $1, (isExpNeg) # признак отрицательного результата
 __powNotOddEnd:
 jmp __powInt 
 __powNotInt:
 mov $powNegError, %rsi 
 call __throughUserError
 __powInt:
 __powIsPos: 

 fldln2
 fld (buf)
 fabs 
 movss %xmm0, (buf)
 fyl2x
 fmul (buf)
 fstp (buf)
 # возводим e^buf 
 fldl2e
 fmul (buf)
 fstp (buf)
 fld (buf) # смешанное число 
 frndint
 fstp (buf2) # целое число
 fld (buf)
 fsub (buf2)
 f2xm1
 fstp (buf)
 fld1 
 fadd (buf)
 fstp (buf)
 fld (buf2)
 fld1 
 fscale 
 fmul (buf)
 fstp (buf)
 
 movb (isExpNeg), %al  
 cmp $0, %al    
 jz __powEnd 
 # результат отрицательный 
 fld (zero)
 fsub (buf)
 fstp (buf)
 __powEnd:
 call __floatToStr
 ret 
 __powZeroBase:
 movss (zero), %xmm2 
 movss %xmm0, %xmm3 
 cmpss $0, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jz __powZeroExpEnd  
 mov $powZeroZeroError, %rsi 
 call __throughUserError
 __powZeroExpEnd: 
 movss (zero), %xmm2 
 movss %xmm0, %xmm3 
 cmpss $1, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax
 jz __powNegExpEnd
 mov $powZeroNegError, %rsi 
 call __throughUserError
 __powNegExpEnd:
 call __floatToStr
 ret 

__floatToStr:
# вход: buf
# выход userData
# проверяем, является ли число отрицательным
movss (zero), %xmm0 
movss (buf), %xmm1 
cmpss $1, %xmm0, %xmm1
pextrb $3, %xmm1, %rax
cmp $0, %rax 
jz __floatToStrIsPos
# меняем знак 
fld (zero) 
fsub (one)
fmul (buf)
fstp (buf)
movb $1, (isNeg) # признак отрицательного числа 
jmp __floatToStrIsNeg
__floatToStrIsPos:
movb $0, (isNeg)
__floatToStrIsNeg:
fld (buf)
movss (buf), %xmm0 
roundps $3, %xmm0, %xmm0
movss %xmm0, (buf)
cvtss2si (buf), %r12
fsub (buf) # вычитаем из значения целое значение, получаем дробное 
fstp (buf)


mov $6, %r10 # 6 знаков после запятой 

__floatToStrLocal:
fld (buf)
cmp $0, %r10
jz __floatToStrOk
dec %r10 
movss (ten), %xmm0
movss %xmm0, (buf)
fmul (buf)
fstp (buf)
jmp __floatToStrLocal
__floatToStrOk:
cvtss2si (buf), %rax # здесь содержится дробное значение 
call __toStr

call __clearBuf
mov $buf, %rax 
movb $48, (%rax)
inc %rax 
movb $0, (%rax)

__floatToStrZeros:
mov $buf2, %rsi 
call __len 
cmp $6, %rax 
jz __floatToStrEndZeros

mov $lenBuf, %r8 
mov $buf, %r9 
mov $buf2, %r11
call __concatinate

mov $lenBuf2, %rsi 
mov $buf2, %rdx 
mov $lenUserData, %rax 
mov $userData, %rdi 
call __set

jmp __floatToStrZeros

__floatToStrEndZeros:

call __clearBuf3 # в buf3 содержится дробное значение в виде строки 
mov $lenBuf3, %rsi 
mov $buf3, %rdx 
mov $lenBuf2, %rax 
mov $buf2, %rdi 
call __set 

mov %r12, %rax 
call __toStr # в buf2 содержится целое значение в виде строки
call __clearBuf
mov $buf, %rax 
movb $'.', (%rax)
inc %rax 
movb $0, (%rax) 

mov $lenBuf2, %r8 
mov $buf2, %r9 
mov $buf, %r11 
call __concatinate



call __clearBuf 
mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenUserData, %rax 
mov $userData, %rdi 
call __set


mov $lenBuf, %r8
mov $buf, %r9 
mov $buf3, %r11 
call __concatinate 

# число отрицательное?
mov (isNeg), %dl 
cmp $1, %dl  
jnz __floatToStrEnd 
call __clearBuf
mov $buf, %rax 
movb $'-', (%rax)
inc %rax 
movb $0, (%rax)

mov $lenBuf3, %rsi 
mov $buf3, %rdx 
mov $lenUserData, %rax 
mov $userData, %rdi 
call __set 

mov $lenBuf, %r8 
mov $buf, %r9 
mov $buf3, %r11 
call __concatinate

__floatToStrEnd:
ret 

__parseFloat:
# $buf - источник (строка)
# %xmm0 - результат
mov $buf, %rax 
call __clearBuf2
call __clearBuf3
mov $buf2, %rbx # здесь будет содержаться целая часть 
mov $buf3, %rcx # здесь будет содержаться дробная часть
mov (%rax), %dl 
cmp $'-', %dl 
jnz isPos
mov $1, %r12 # признак отрицательного числа 
inc %rax 
jmp __parseFloatLocal 
isPos:
mov $0, %r12 
__parseFloatLocal: 
mov (%rax), %dl 
cmp $'.', %dl
jz __point
mov %dl, (%rbx)
inc %rax 
inc %rbx 
jmp __parseFloatLocal
__point:
movb $0, (%rbx)
mov %rax, %rbx 
mov $buf2, %rsi 
call __len 
cmp $8, %rax # целое число - не более 7 цифр 
jl __parseFloatZ
mov $buf2, %rbx 
movb $48, (%rbx)
inc %rbx 
movb $0, (%rbx) 
mov $buf3, %rbx 
movb $48, (%rbx)
inc %rbx 
movb $0, (%rbx)
jmp __parseNow 
__parseFloatZ: 
mov %rbx, %rax 

__pointLocal:             
inc %rax
mov (%rax), %dl 
cmp $0, %dl 
jz __parseNow
mov (%rax), %dl 
mov %dl, (%rcx) 
inc %rcx 
jmp __pointLocal   
__parseNow:
movb $0, (%rcx)

call __clearBuf
mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenBuf2, %rax 
mov $buf2, %rdi 
call __set 
call __toNumber
mov %rax, %r10 # целая часть числа в %r10

call __clearBuf
mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenBuf3, %rax 
mov $buf3, %rdi 
call __set

mov $buf, %rsi 
call __len 
mov %rax, %rbx # длина дробной части числа в %rbx 
cmp $6, %rbx # дробная часть не более шести знаков 
jl __parseFloatCut
mov $buf, %rsi 
add $6, %rsi 
movb $0, (%rsi)
mov $6, %rbx 
__parseFloatCut: 
call __toNumber
mov %rax, (buf)
cvtsi2ss (buf), %xmm0  
movss %xmm0, (buf)

__floatLocal:
fld (buf)
cmp $0, %rbx 
jz __floatOk
dec %rbx 
movss (ten), %xmm0
movss %xmm0, (buf)
fdiv (buf)
fstp (buf)
jmp __floatLocal
__floatOk:
mov %r10, (buf)
cvtsi2ss (buf), %xmm0    
movss %xmm0, (buf) # целая часть числа 
fadd (buf)
fstp (buf)
movss (buf), %xmm0  
cmp $1, %r12 
jnz __pos
mov (zero), %rax   
mov %rax, (buf)
fld (buf)
movss %xmm0, (buf)
fsub (buf)
fstp (buf)
movss (buf), %xmm0 
__pos:
ret 

.globl _start
_start:
 call __firstMem
 call __firstStrMem

mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName0, %rax 
 mov $varName0, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName0, %rax 
 mov $varName0, %rdi 
call __set

 mov $data0, %rax  
 mov %rax, (userData)
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName3, %rax 
 mov $varName3, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName5, %rax 
 mov $varName5, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName7, %rax 
 mov $varName7, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName8, %rax 
 mov $varName8, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName9, %rax 
 mov $varName9, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName10, %rax 
 mov $varName10, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName11, %rax 
 mov $varName11, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName12, %rax 
 mov $varName12, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
 mov $varName13, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName14, %rax 
 mov $varName14, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName15, %rax 
 mov $varName15, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
 mov $varName16, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName17, %rax 
 mov $varName17, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName18, %rax 
 mov $varName18, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName19, %rax 
 mov $varName19, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName20, %rax 
 mov $varName20, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName21, %rax 
 mov $varName21, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName22, %rax 
 mov $varName22, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName23, %rax 
 mov $varName23, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName25, %rax 
 mov $varName25, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName26, %rax 
 mov $varName26, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName27, %rax 
 mov $varName27, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName28, %rax 
 mov $varName28, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName29, %rax 
 mov $varName29, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName30, %rax 
 mov $varName30, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName31, %rax 
 mov $varName31, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName32, %rax 
 mov $varName32, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName33, %rax 
 mov $varName33, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName34, %rax 
 mov $varName34, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName35, %rax 
 mov $varName35, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName36, %rax 
 mov $varName36, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName37, %rax 
 mov $varName37, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName38, %rax 
 mov $varName38, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName39, %rax 
 mov $varName39, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName40, %rax 
 mov $varName40, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName42, %rax 
 mov $varName42, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName43, %rax 
 mov $varName43, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName44, %rax 
 mov $varName44, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName45, %rax 
 mov $varName45, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName46, %rax 
 mov $varName46, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName47, %rax 
 mov $varName47, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName48, %rax 
 mov $varName48, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName49, %rax 
 mov $varName49, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName50, %rax 
 mov $varName50, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName51, %rax 
 mov $varName51, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName52, %rax 
 mov $varName52, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName53, %rax 
 mov $varName53, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName54, %rax 
 mov $varName54, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName55, %rax 
 mov $varName55, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName56, %rax 
 mov $varName56, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName57, %rax 
 mov $varName57, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName58, %rax 
 mov $varName58, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName59, %rax 
 mov $varName59, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName60, %rax 
 mov $varName60, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName61, %rax 
 mov $varName61, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName62, %rax 
 mov $varName62, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName63, %rax 
 mov $varName63, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName64, %rax 
 mov $varName64, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName65, %rax 
 mov $varName65, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName66, %rax 
 mov $varName66, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName67, %rax 
 mov $varName67, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName68, %rax 
 mov $varName68, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName69, %rax 
 mov $varName69, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName70, %rax 
 mov $varName70, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName71, %rax 
 mov $varName71, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName72, %rax 
 mov $varName72, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName73, %rax 
 mov $varName73, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName74, %rax 
 mov $varName74, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName75, %rax 
 mov $varName75, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName76, %rax 
 mov $varName76, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName77, %rax 
 mov $varName77, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName78, %rax 
 mov $varName78, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName79, %rax 
 mov $varName79, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName80, %rax 
 mov $varName80, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName81, %rax 
 mov $varName81, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName82, %rax 
 mov $varName82, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName83, %rax 
 mov $varName83, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName84, %rax 
 mov $varName84, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName85, %rax 
 mov $varName85, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName86, %rax 
 mov $varName86, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName87, %rax 
 mov $varName87, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName88, %rax 
 mov $varName88, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName89, %rax 
 mov $varName89, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName90, %rax 
 mov $varName90, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName91, %rax 
 mov $varName91, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName92, %rax 
 mov $varName92, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName93, %rax 
 mov $varName93, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName94, %rax 
 mov $varName94, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName95, %rax 
 mov $varName95, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName96, %rax 
 mov $varName96, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName97, %rax 
 mov $varName97, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName98, %rax 
 mov $varName98, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName99, %rax 
 mov $varName99, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName100, %rax 
 mov $varName100, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName101, %rax 
 mov $varName101, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName102, %rax 
 mov $varName102, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName103, %rax 
 mov $varName103, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName104, %rax 
 mov $varName104, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName105, %rax 
 mov $varName105, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName106, %rax 
 mov $varName106, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName107, %rax 
 mov $varName107, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName108, %rax 
 mov $varName108, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName109, %rax 
 mov $varName109, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName110, %rax 
 mov $varName110, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName111, %rax 
 mov $varName111, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName112, %rax 
 mov $varName112, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName113, %rax 
 mov $varName113, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName114, %rax 
 mov $varName114, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName115, %rax 
 mov $varName115, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName116, %rax 
 mov $varName116, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName117, %rax 
 mov $varName117, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName118, %rax 
 mov $varName118, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName119, %rax 
 mov $varName119, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName120, %rax 
 mov $varName120, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName121, %rax 
 mov $varName121, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName122, %rax 
 mov $varName122, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName123, %rax 
 mov $varName123, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName124, %rax 
 mov $varName124, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName125, %rax 
 mov $varName125, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName126, %rax 
 mov $varName126, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName127, %rax 
 mov $varName127, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName128, %rax 
 mov $varName128, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName129, %rax 
 mov $varName129, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName130, %rax 
 mov $varName130, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName131, %rax 
 mov $varName131, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName132, %rax 
 mov $varName132, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName133, %rax 
 mov $varName133, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName134, %rax 
 mov $varName134, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName135, %rax 
 mov $varName135, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName136, %rax 
 mov $varName136, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName137, %rax 
 mov $varName137, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName138, %rax 
 mov $varName138, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName139, %rax 
 mov $varName139, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName140, %rax 
 mov $varName140, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName141, %rax 
 mov $varName141, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName142, %rax 
 mov $varName142, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName143, %rax 
 mov $varName143, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName144, %rax 
 mov $varName144, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName145, %rax 
 mov $varName145, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName146, %rax 
 mov $varName146, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName147, %rax 
 mov $varName147, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName148, %rax 
 mov $varName148, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName149, %rax 
 mov $varName149, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName150, %rax 
 mov $varName150, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName151, %rax 
 mov $varName151, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName152, %rax 
 mov $varName152, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName153, %rax 
 mov $varName153, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName154, %rax 
 mov $varName154, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName155, %rax 
 mov $varName155, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName156, %rax 
 mov $varName156, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName157, %rax 
 mov $varName157, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName158, %rax 
 mov $varName158, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName159, %rax 
 mov $varName159, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName160, %rax 
 mov $varName160, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName161, %rax 
 mov $varName161, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName162, %rax 
 mov $varName162, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName163, %rax 
 mov $varName163, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName164, %rax 
 mov $varName164, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName165, %rax 
 mov $varName165, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName166, %rax 
 mov $varName166, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName167, %rax 
 mov $varName167, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName168, %rax 
 mov $varName168, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName169, %rax 
 mov $varName169, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName170, %rax 
 mov $varName170, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName171, %rax 
 mov $varName171, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName172, %rax 
 mov $varName172, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName173, %rax 
 mov $varName173, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName174, %rax 
 mov $varName174, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName175, %rax 
 mov $varName175, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName176, %rax 
 mov $varName176, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName177, %rax 
 mov $varName177, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName178, %rax 
 mov $varName178, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName179, %rax 
 mov $varName179, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName180, %rax 
 mov $varName180, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName181, %rax 
 mov $varName181, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName182, %rax 
 mov $varName182, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName183, %rax 
 mov $varName183, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName184, %rax 
 mov $varName184, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName185, %rax 
 mov $varName185, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName186, %rax 
 mov $varName186, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName187, %rax 
 mov $varName187, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName188, %rax 
 mov $varName188, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName189, %rax 
 mov $varName189, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName190, %rax 
 mov $varName190, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName191, %rax 
 mov $varName191, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName192, %rax 
 mov $varName192, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName193, %rax 
 mov $varName193, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName194, %rax 
 mov $varName194, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName195, %rax 
 mov $varName195, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName196, %rax 
 mov $varName196, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName197, %rax 
 mov $varName197, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName198, %rax 
 mov $varName198, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName199, %rax 
 mov $varName199, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName200, %rax 
 mov $varName200, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName201, %rax 
 mov $varName201, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName202, %rax 
 mov $varName202, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName203, %rax 
 mov $varName203, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName204, %rax 
 mov $varName204, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName205, %rax 
 mov $varName205, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName206, %rax 
 mov $varName206, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName207, %rax 
 mov $varName207, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName208, %rax 
 mov $varName208, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName209, %rax 
 mov $varName209, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName210, %rax 
 mov $varName210, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName211, %rax 
 mov $varName211, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName212, %rax 
 mov $varName212, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName213, %rax 
 mov $varName213, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName214, %rax 
 mov $varName214, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName215, %rax 
 mov $varName215, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName216, %rax 
 mov $varName216, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName217, %rax 
 mov $varName217, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName218, %rax 
 mov $varName218, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName219, %rax 
 mov $varName219, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName220, %rax 
 mov $varName220, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName221, %rax 
 mov $varName221, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName222, %rax 
 mov $varName222, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName223, %rax 
 mov $varName223, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName224, %rax 
 mov $varName224, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName225, %rax 
 mov $varName225, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName226, %rax 
 mov $varName226, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName227, %rax 
 mov $varName227, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName228, %rax 
 mov $varName228, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName229, %rax 
 mov $varName229, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName230, %rax 
 mov $varName230, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName231, %rax 
 mov $varName231, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName232, %rax 
 mov $varName232, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName233, %rax 
 mov $varName233, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName234, %rax 
 mov $varName234, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName235, %rax 
 mov $varName235, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName236, %rax 
 mov $varName236, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName237, %rax 
 mov $varName237, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName238, %rax 
 mov $varName238, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName239, %rax 
 mov $varName239, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName240, %rax 
 mov $varName240, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName241, %rax 
 mov $varName241, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName242, %rax 
 mov $varName242, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName243, %rax 
 mov $varName243, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName244, %rax 
 mov $varName244, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName245, %rax 
 mov $varName245, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName246, %rax 
 mov $varName246, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName247, %rax 
 mov $varName247, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName248, %rax 
 mov $varName248, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName249, %rax 
 mov $varName249, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName250, %rax 
 mov $varName250, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName251, %rax 
 mov $varName251, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName252, %rax 
 mov $varName252, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName253, %rax 
 mov $varName253, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName254, %rax 
 mov $varName254, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName255, %rax 
 mov $varName255, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName256, %rax 
 mov $varName256, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName257, %rax 
 mov $varName257, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName258, %rax 
 mov $varName258, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName259, %rax 
 mov $varName259, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName260, %rax 
 mov $varName260, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName261, %rax 
 mov $varName261, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName262, %rax 
 mov $varName262, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName263, %rax 
 mov $varName263, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName264, %rax 
 mov $varName264, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName265, %rax 
 mov $varName265, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName266, %rax 
 mov $varName266, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName267, %rax 
 mov $varName267, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName268, %rax 
 mov $varName268, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName269, %rax 
 mov $varName269, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName270, %rax 
 mov $varName270, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName271, %rax 
 mov $varName271, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName272, %rax 
 mov $varName272, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName273, %rax 
 mov $varName273, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName274, %rax 
 mov $varName274, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName275, %rax 
 mov $varName275, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName276, %rax 
 mov $varName276, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName277, %rax 
 mov $varName277, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName278, %rax 
 mov $varName278, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName279, %rax 
 mov $varName279, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName280, %rax 
 mov $varName280, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName281, %rax 
 mov $varName281, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName282, %rax 
 mov $varName282, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName283, %rax 
 mov $varName283, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName284, %rax 
 mov $varName284, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName285, %rax 
 mov $varName285, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName286, %rax 
 mov $varName286, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName287, %rax 
 mov $varName287, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName288, %rax 
 mov $varName288, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName289, %rax 
 mov $varName289, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName290, %rax 
 mov $varName290, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName291, %rax 
 mov $varName291, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName292, %rax 
 mov $varName292, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName293, %rax 
 mov $varName293, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName294, %rax 
 mov $varName294, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName295, %rax 
 mov $varName295, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName296, %rax 
 mov $varName296, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName297, %rax 
 mov $varName297, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName298, %rax 
 mov $varName298, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName299, %rax 
 mov $varName299, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName300, %rax 
 mov $varName300, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName301, %rax 
 mov $varName301, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName302, %rax 
 mov $varName302, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName303, %rax 
 mov $varName303, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName304, %rax 
 mov $varName304, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName305, %rax 
 mov $varName305, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName306, %rax 
 mov $varName306, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName307, %rax 
 mov $varName307, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName308, %rax 
 mov $varName308, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName309, %rax 
 mov $varName309, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName310, %rax 
 mov $varName310, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName311, %rax 
 mov $varName311, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName312, %rax 
 mov $varName312, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName313, %rax 
 mov $varName313, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName314, %rax 
 mov $varName314, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName315, %rax 
 mov $varName315, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName316, %rax 
 mov $varName316, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName317, %rax 
 mov $varName317, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName318, %rax 
 mov $varName318, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName319, %rax 
 mov $varName319, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName320, %rax 
 mov $varName320, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName321, %rax 
 mov $varName321, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName322, %rax 
 mov $varName322, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName323, %rax 
 mov $varName323, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName324, %rax 
 mov $varName324, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName325, %rax 
 mov $varName325, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName326, %rax 
 mov $varName326, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName327, %rax 
 mov $varName327, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName328, %rax 
 mov $varName328, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName329, %rax 
 mov $varName329, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName330, %rax 
 mov $varName330, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName331, %rax 
 mov $varName331, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName332, %rax 
 mov $varName332, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName333, %rax 
 mov $varName333, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName334, %rax 
 mov $varName334, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName335, %rax 
 mov $varName335, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName336, %rax 
 mov $varName336, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName337, %rax 
 mov $varName337, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName338, %rax 
 mov $varName338, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName339, %rax 
 mov $varName339, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName340, %rax 
 mov $varName340, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName341, %rax 
 mov $varName341, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName342, %rax 
 mov $varName342, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName343, %rax 
 mov $varName343, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName344, %rax 
 mov $varName344, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName345, %rax 
 mov $varName345, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName346, %rax 
 mov $varName346, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName347, %rax 
 mov $varName347, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName348, %rax 
 mov $varName348, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName349, %rax 
 mov $varName349, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName350, %rax 
 mov $varName350, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName351, %rax 
 mov $varName351, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName352, %rax 
 mov $varName352, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName353, %rax 
 mov $varName353, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName354, %rax 
 mov $varName354, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName355, %rax 
 mov $varName355, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName356, %rax 
 mov $varName356, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName357, %rax 
 mov $varName357, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName358, %rax 
 mov $varName358, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName359, %rax 
 mov $varName359, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName360, %rax 
 mov $varName360, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName361, %rax 
 mov $varName361, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName362, %rax 
 mov $varName362, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName363, %rax 
 mov $varName363, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName364, %rax 
 mov $varName364, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName365, %rax 
 mov $varName365, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName366, %rax 
 mov $varName366, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName367, %rax 
 mov $varName367, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName368, %rax 
 mov $varName368, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName369, %rax 
 mov $varName369, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName370, %rax 
 mov $varName370, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName371, %rax 
 mov $varName371, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName372, %rax 
 mov $varName372, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName373, %rax 
 mov $varName373, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName374, %rax 
 mov $varName374, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName375, %rax 
 mov $varName375, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName376, %rax 
 mov $varName376, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName377, %rax 
 mov $varName377, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName378, %rax 
 mov $varName378, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName379, %rax 
 mov $varName379, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName380, %rax 
 mov $varName380, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName381, %rax 
 mov $varName381, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName382, %rax 
 mov $varName382, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName383, %rax 
 mov $varName383, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName384, %rax 
 mov $varName384, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName385, %rax 
 mov $varName385, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName386, %rax 
 mov $varName386, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName387, %rax 
 mov $varName387, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName388, %rax 
 mov $varName388, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName389, %rax 
 mov $varName389, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName390, %rax 
 mov $varName390, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName391, %rax 
 mov $varName391, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName392, %rax 
 mov $varName392, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName393, %rax 
 mov $varName393, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName394, %rax 
 mov $varName394, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName395, %rax 
 mov $varName395, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName396, %rax 
 mov $varName396, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName397, %rax 
 mov $varName397, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName398, %rax 
 mov $varName398, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName399, %rax 
 mov $varName399, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName400, %rax 
 mov $varName400, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName401, %rax 
 mov $varName401, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName402, %rax 
 mov $varName402, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName403, %rax 
 mov $varName403, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName404, %rax 
 mov $varName404, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName405, %rax 
 mov $varName405, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName406, %rax 
 mov $varName406, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName407, %rax 
 mov $varName407, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName408, %rax 
 mov $varName408, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName409, %rax 
 mov $varName409, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName410, %rax 
 mov $varName410, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName411, %rax 
 mov $varName411, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName412, %rax 
 mov $varName412, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName413, %rax 
 mov $varName413, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName414, %rax 
 mov $varName414, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName415, %rax 
 mov $varName415, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName416, %rax 
 mov $varName416, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName417, %rax 
 mov $varName417, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName418, %rax 
 mov $varName418, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName419, %rax 
 mov $varName419, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName420, %rax 
 mov $varName420, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName421, %rax 
 mov $varName421, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName422, %rax 
 mov $varName422, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName423, %rax 
 mov $varName423, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName424, %rax 
 mov $varName424, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName425, %rax 
 mov $varName425, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName426, %rax 
 mov $varName426, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName427, %rax 
 mov $varName427, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName428, %rax 
 mov $varName428, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName429, %rax 
 mov $varName429, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName430, %rax 
 mov $varName430, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName431, %rax 
 mov $varName431, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName432, %rax 
 mov $varName432, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName433, %rax 
 mov $varName433, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName434, %rax 
 mov $varName434, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName435, %rax 
 mov $varName435, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName436, %rax 
 mov $varName436, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName437, %rax 
 mov $varName437, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName438, %rax 
 mov $varName438, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName439, %rax 
 mov $varName439, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName440, %rax 
 mov $varName440, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName441, %rax 
 mov $varName441, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName442, %rax 
 mov $varName442, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName443, %rax 
 mov $varName443, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName444, %rax 
 mov $varName444, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName445, %rax 
 mov $varName445, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName446, %rax 
 mov $varName446, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName447, %rax 
 mov $varName447, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName448, %rax 
 mov $varName448, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName449, %rax 
 mov $varName449, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName450, %rax 
 mov $varName450, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName451, %rax 
 mov $varName451, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName452, %rax 
 mov $varName452, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName453, %rax 
 mov $varName453, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName454, %rax 
 mov $varName454, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName455, %rax 
 mov $varName455, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName456, %rax 
 mov $varName456, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName457, %rax 
 mov $varName457, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName458, %rax 
 mov $varName458, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName459, %rax 
 mov $varName459, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName460, %rax 
 mov $varName460, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName461, %rax 
 mov $varName461, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName462, %rax 
 mov $varName462, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName463, %rax 
 mov $varName463, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName464, %rax 
 mov $varName464, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName465, %rax 
 mov $varName465, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName466, %rax 
 mov $varName466, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName467, %rax 
 mov $varName467, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName468, %rax 
 mov $varName468, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName469, %rax 
 mov $varName469, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName470, %rax 
 mov $varName470, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName471, %rax 
 mov $varName471, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName472, %rax 
 mov $varName472, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName473, %rax 
 mov $varName473, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName474, %rax 
 mov $varName474, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName475, %rax 
 mov $varName475, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName476, %rax 
 mov $varName476, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName477, %rax 
 mov $varName477, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName478, %rax 
 mov $varName478, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName479, %rax 
 mov $varName479, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName480, %rax 
 mov $varName480, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName481, %rax 
 mov $varName481, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName482, %rax 
 mov $varName482, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName483, %rax 
 mov $varName483, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName484, %rax 
 mov $varName484, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName485, %rax 
 mov $varName485, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName486, %rax 
 mov $varName486, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName487, %rax 
 mov $varName487, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName488, %rax 
 mov $varName488, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName489, %rax 
 mov $varName489, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName490, %rax 
 mov $varName490, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName491, %rax 
 mov $varName491, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName492, %rax 
 mov $varName492, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName493, %rax 
 mov $varName493, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName494, %rax 
 mov $varName494, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName495, %rax 
 mov $varName495, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName496, %rax 
 mov $varName496, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName497, %rax 
 mov $varName497, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName498, %rax 
 mov $varName498, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName499, %rax 
 mov $varName499, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName500, %rax 
 mov $varName500, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName501, %rax 
 mov $varName501, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName502, %rax 
 mov $varName502, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName503, %rax 
 mov $varName503, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName504, %rax 
 mov $varName504, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName505, %rax 
 mov $varName505, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName506, %rax 
 mov $varName506, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName507, %rax 
 mov $varName507, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName508, %rax 
 mov $varName508, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName509, %rax 
 mov $varName509, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName510, %rax 
 mov $varName510, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName511, %rax 
 mov $varName511, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName512, %rax 
 mov $varName512, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName513, %rax 
 mov $varName513, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName514, %rax 
 mov $varName514, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName515, %rax 
 mov $varName515, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName516, %rax 
 mov $varName516, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName517, %rax 
 mov $varName517, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName518, %rax 
 mov $varName518, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName519, %rax 
 mov $varName519, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName520, %rax 
 mov $varName520, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName521, %rax 
 mov $varName521, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName522, %rax 
 mov $varName522, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName523, %rax 
 mov $varName523, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName524, %rax 
 mov $varName524, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName525, %rax 
 mov $varName525, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName526, %rax 
 mov $varName526, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName527, %rax 
 mov $varName527, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName528, %rax 
 mov $varName528, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName529, %rax 
 mov $varName529, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName530, %rax 
 mov $varName530, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName531, %rax 
 mov $varName531, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName532, %rax 
 mov $varName532, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName533, %rax 
 mov $varName533, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName534, %rax 
 mov $varName534, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName535, %rax 
 mov $varName535, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName536, %rax 
 mov $varName536, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName537, %rax 
 mov $varName537, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName538, %rax 
 mov $varName538, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName539, %rax 
 mov $varName539, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName540, %rax 
 mov $varName540, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName541, %rax 
 mov $varName541, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName542, %rax 
 mov $varName542, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName543, %rax 
 mov $varName543, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName544, %rax 
 mov $varName544, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName545, %rax 
 mov $varName545, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName546, %rax 
 mov $varName546, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName547, %rax 
 mov $varName547, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName548, %rax 
 mov $varName548, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName549, %rax 
 mov $varName549, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName550, %rax 
 mov $varName550, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName551, %rax 
 mov $varName551, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName552, %rax 
 mov $varName552, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName553, %rax 
 mov $varName553, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName554, %rax 
 mov $varName554, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName555, %rax 
 mov $varName555, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName556, %rax 
 mov $varName556, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName557, %rax 
 mov $varName557, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName558, %rax 
 mov $varName558, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName559, %rax 
 mov $varName559, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName560, %rax 
 mov $varName560, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName561, %rax 
 mov $varName561, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName562, %rax 
 mov $varName562, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName563, %rax 
 mov $varName563, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName564, %rax 
 mov $varName564, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName565, %rax 
 mov $varName565, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName566, %rax 
 mov $varName566, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName567, %rax 
 mov $varName567, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName568, %rax 
 mov $varName568, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName569, %rax 
 mov $varName569, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName570, %rax 
 mov $varName570, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName571, %rax 
 mov $varName571, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName572, %rax 
 mov $varName572, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName573, %rax 
 mov $varName573, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName574, %rax 
 mov $varName574, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName575, %rax 
 mov $varName575, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName576, %rax 
 mov $varName576, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName577, %rax 
 mov $varName577, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName578, %rax 
 mov $varName578, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName579, %rax 
 mov $varName579, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName580, %rax 
 mov $varName580, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName581, %rax 
 mov $varName581, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName582, %rax 
 mov $varName582, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName583, %rax 
 mov $varName583, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName584, %rax 
 mov $varName584, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName585, %rax 
 mov $varName585, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName586, %rax 
 mov $varName586, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName587, %rax 
 mov $varName587, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName588, %rax 
 mov $varName588, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName589, %rax 
 mov $varName589, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName590, %rax 
 mov $varName590, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName591, %rax 
 mov $varName591, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName592, %rax 
 mov $varName592, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName593, %rax 
 mov $varName593, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName594, %rax 
 mov $varName594, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName595, %rax 
 mov $varName595, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName596, %rax 
 mov $varName596, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName597, %rax 
 mov $varName597, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName598, %rax 
 mov $varName598, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName599, %rax 
 mov $varName599, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName600, %rax 
 mov $varName600, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName601, %rax 
 mov $varName601, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName602, %rax 
 mov $varName602, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName603, %rax 
 mov $varName603, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName604, %rax 
 mov $varName604, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName605, %rax 
 mov $varName605, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName606, %rax 
 mov $varName606, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName607, %rax 
 mov $varName607, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName608, %rax 
 mov $varName608, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName609, %rax 
 mov $varName609, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName610, %rax 
 mov $varName610, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName611, %rax 
 mov $varName611, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName612, %rax 
 mov $varName612, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName613, %rax 
 mov $varName613, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName614, %rax 
 mov $varName614, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName615, %rax 
 mov $varName615, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName616, %rax 
 mov $varName616, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName617, %rax 
 mov $varName617, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName618, %rax 
 mov $varName618, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName619, %rax 
 mov $varName619, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName620, %rax 
 mov $varName620, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName621, %rax 
 mov $varName621, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName622, %rax 
 mov $varName622, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName623, %rax 
 mov $varName623, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName624, %rax 
 mov $varName624, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName625, %rax 
 mov $varName625, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName626, %rax 
 mov $varName626, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName627, %rax 
 mov $varName627, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName628, %rax 
 mov $varName628, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName629, %rax 
 mov $varName629, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName630, %rax 
 mov $varName630, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName631, %rax 
 mov $varName631, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName632, %rax 
 mov $varName632, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName633, %rax 
 mov $varName633, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName634, %rax 
 mov $varName634, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName635, %rax 
 mov $varName635, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName636, %rax 
 mov $varName636, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName637, %rax 
 mov $varName637, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName638, %rax 
 mov $varName638, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName639, %rax 
 mov $varName639, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName640, %rax 
 mov $varName640, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName641, %rax 
 mov $varName641, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName642, %rax 
 mov $varName642, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName643, %rax 
 mov $varName643, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName644, %rax 
 mov $varName644, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName645, %rax 
 mov $varName645, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName646, %rax 
 mov $varName646, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName647, %rax 
 mov $varName647, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName648, %rax 
 mov $varName648, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName649, %rax 
 mov $varName649, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName650, %rax 
 mov $varName650, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName651, %rax 
 mov $varName651, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName652, %rax 
 mov $varName652, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName653, %rax 
 mov $varName653, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName654, %rax 
 mov $varName654, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName655, %rax 
 mov $varName655, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName656, %rax 
 mov $varName656, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName657, %rax 
 mov $varName657, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName658, %rax 
 mov $varName658, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName659, %rax 
 mov $varName659, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName660, %rax 
 mov $varName660, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName661, %rax 
 mov $varName661, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName662, %rax 
 mov $varName662, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName663, %rax 
 mov $varName663, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName664, %rax 
 mov $varName664, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName665, %rax 
 mov $varName665, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName666, %rax 
 mov $varName666, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName667, %rax 
 mov $varName667, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName668, %rax 
 mov $varName668, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName669, %rax 
 mov $varName669, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName670, %rax 
 mov $varName670, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName671, %rax 
 mov $varName671, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName672, %rax 
 mov $varName672, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName673, %rax 
 mov $varName673, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName674, %rax 
 mov $varName674, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName675, %rax 
 mov $varName675, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName676, %rax 
 mov $varName676, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName677, %rax 
 mov $varName677, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName678, %rax 
 mov $varName678, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName679, %rax 
 mov $varName679, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName680, %rax 
 mov $varName680, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName681, %rax 
 mov $varName681, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName682, %rax 
 mov $varName682, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName683, %rax 
 mov $varName683, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName684, %rax 
 mov $varName684, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName685, %rax 
 mov $varName685, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName686, %rax 
 mov $varName686, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName687, %rax 
 mov $varName687, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName688, %rax 
 mov $varName688, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName689, %rax 
 mov $varName689, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName690, %rax 
 mov $varName690, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName691, %rax 
 mov $varName691, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName692, %rax 
 mov $varName692, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName693, %rax 
 mov $varName693, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName694, %rax 
 mov $varName694, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName695, %rax 
 mov $varName695, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName696, %rax 
 mov $varName696, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName697, %rax 
 mov $varName697, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName698, %rax 
 mov $varName698, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName699, %rax 
 mov $varName699, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName700, %rax 
 mov $varName700, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName701, %rax 
 mov $varName701, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName702, %rax 
 mov $varName702, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName703, %rax 
 mov $varName703, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName704, %rax 
 mov $varName704, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName705, %rax 
 mov $varName705, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName706, %rax 
 mov $varName706, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName707, %rax 
 mov $varName707, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName708, %rax 
 mov $varName708, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName709, %rax 
 mov $varName709, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName710, %rax 
 mov $varName710, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName711, %rax 
 mov $varName711, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName712, %rax 
 mov $varName712, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName713, %rax 
 mov $varName713, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName714, %rax 
 mov $varName714, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName715, %rax 
 mov $varName715, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName716, %rax 
 mov $varName716, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName717, %rax 
 mov $varName717, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName718, %rax 
 mov $varName718, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName719, %rax 
 mov $varName719, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName720, %rax 
 mov $varName720, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName721, %rax 
 mov $varName721, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName722, %rax 
 mov $varName722, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName723, %rax 
 mov $varName723, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName724, %rax 
 mov $varName724, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName725, %rax 
 mov $varName725, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName726, %rax 
 mov $varName726, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName727, %rax 
 mov $varName727, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName728, %rax 
 mov $varName728, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName729, %rax 
 mov $varName729, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName730, %rax 
 mov $varName730, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName731, %rax 
 mov $varName731, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName732, %rax 
 mov $varName732, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName733, %rax 
 mov $varName733, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName734, %rax 
 mov $varName734, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName735, %rax 
 mov $varName735, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName736, %rax 
 mov $varName736, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName737, %rax 
 mov $varName737, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName738, %rax 
 mov $varName738, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName739, %rax 
 mov $varName739, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName740, %rax 
 mov $varName740, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName741, %rax 
 mov $varName741, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName742, %rax 
 mov $varName742, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName743, %rax 
 mov $varName743, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName744, %rax 
 mov $varName744, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName745, %rax 
 mov $varName745, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName746, %rax 
 mov $varName746, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName747, %rax 
 mov $varName747, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName748, %rax 
 mov $varName748, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName749, %rax 
 mov $varName749, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName750, %rax 
 mov $varName750, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName751, %rax 
 mov $varName751, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName752, %rax 
 mov $varName752, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName753, %rax 
 mov $varName753, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName754, %rax 
 mov $varName754, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName755, %rax 
 mov $varName755, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName756, %rax 
 mov $varName756, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName757, %rax 
 mov $varName757, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName758, %rax 
 mov $varName758, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName759, %rax 
 mov $varName759, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName760, %rax 
 mov $varName760, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName761, %rax 
 mov $varName761, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName762, %rax 
 mov $varName762, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName763, %rax 
 mov $varName763, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName764, %rax 
 mov $varName764, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName765, %rax 
 mov $varName765, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName766, %rax 
 mov $varName766, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName767, %rax 
 mov $varName767, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName768, %rax 
 mov $varName768, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName769, %rax 
 mov $varName769, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName770, %rax 
 mov $varName770, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName771, %rax 
 mov $varName771, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName772, %rax 
 mov $varName772, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName773, %rax 
 mov $varName773, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName774, %rax 
 mov $varName774, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName775, %rax 
 mov $varName775, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName776, %rax 
 mov $varName776, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName777, %rax 
 mov $varName777, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName778, %rax 
 mov $varName778, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName779, %rax 
 mov $varName779, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName780, %rax 
 mov $varName780, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName781, %rax 
 mov $varName781, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName782, %rax 
 mov $varName782, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName783, %rax 
 mov $varName783, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName784, %rax 
 mov $varName784, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName785, %rax 
 mov $varName785, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName786, %rax 
 mov $varName786, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName787, %rax 
 mov $varName787, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName788, %rax 
 mov $varName788, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName789, %rax 
 mov $varName789, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName790, %rax 
 mov $varName790, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName791, %rax 
 mov $varName791, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName792, %rax 
 mov $varName792, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName793, %rax 
 mov $varName793, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName794, %rax 
 mov $varName794, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName795, %rax 
 mov $varName795, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName796, %rax 
 mov $varName796, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName797, %rax 
 mov $varName797, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName798, %rax 
 mov $varName798, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName799, %rax 
 mov $varName799, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName800, %rax 
 mov $varName800, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName801, %rax 
 mov $varName801, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName802, %rax 
 mov $varName802, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName803, %rax 
 mov $varName803, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName804, %rax 
 mov $varName804, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName805, %rax 
 mov $varName805, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName806, %rax 
 mov $varName806, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName807, %rax 
 mov $varName807, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName808, %rax 
 mov $varName808, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName809, %rax 
 mov $varName809, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName810, %rax 
 mov $varName810, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName811, %rax 
 mov $varName811, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName812, %rax 
 mov $varName812, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName813, %rax 
 mov $varName813, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName814, %rax 
 mov $varName814, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName815, %rax 
 mov $varName815, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName816, %rax 
 mov $varName816, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName817, %rax 
 mov $varName817, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName818, %rax 
 mov $varName818, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName819, %rax 
 mov $varName819, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName820, %rax 
 mov $varName820, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName821, %rax 
 mov $varName821, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName822, %rax 
 mov $varName822, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName823, %rax 
 mov $varName823, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName824, %rax 
 mov $varName824, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName825, %rax 
 mov $varName825, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName826, %rax 
 mov $varName826, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName827, %rax 
 mov $varName827, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName828, %rax 
 mov $varName828, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName829, %rax 
 mov $varName829, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName830, %rax 
 mov $varName830, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName831, %rax 
 mov $varName831, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName832, %rax 
 mov $varName832, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName833, %rax 
 mov $varName833, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName834, %rax 
 mov $varName834, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName835, %rax 
 mov $varName835, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName836, %rax 
 mov $varName836, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName837, %rax 
 mov $varName837, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName838, %rax 
 mov $varName838, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName839, %rax 
 mov $varName839, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName840, %rax 
 mov $varName840, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName841, %rax 
 mov $varName841, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName842, %rax 
 mov $varName842, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName843, %rax 
 mov $varName843, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName844, %rax 
 mov $varName844, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName845, %rax 
 mov $varName845, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName846, %rax 
 mov $varName846, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName847, %rax 
 mov $varName847, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName848, %rax 
 mov $varName848, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName849, %rax 
 mov $varName849, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName850, %rax 
 mov $varName850, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName851, %rax 
 mov $varName851, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName852, %rax 
 mov $varName852, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName853, %rax 
 mov $varName853, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName854, %rax 
 mov $varName854, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName855, %rax 
 mov $varName855, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName856, %rax 
 mov $varName856, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName857, %rax 
 mov $varName857, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName858, %rax 
 mov $varName858, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName859, %rax 
 mov $varName859, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName860, %rax 
 mov $varName860, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName861, %rax 
 mov $varName861, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName862, %rax 
 mov $varName862, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName863, %rax 
 mov $varName863, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName864, %rax 
 mov $varName864, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName865, %rax 
 mov $varName865, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName866, %rax 
 mov $varName866, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName867, %rax 
 mov $varName867, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName868, %rax 
 mov $varName868, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName869, %rax 
 mov $varName869, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName870, %rax 
 mov $varName870, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName871, %rax 
 mov $varName871, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName872, %rax 
 mov $varName872, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName873, %rax 
 mov $varName873, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName874, %rax 
 mov $varName874, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName875, %rax 
 mov $varName875, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName876, %rax 
 mov $varName876, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName877, %rax 
 mov $varName877, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName878, %rax 
 mov $varName878, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName879, %rax 
 mov $varName879, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName880, %rax 
 mov $varName880, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName881, %rax 
 mov $varName881, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName882, %rax 
 mov $varName882, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName883, %rax 
 mov $varName883, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName884, %rax 
 mov $varName884, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName885, %rax 
 mov $varName885, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName886, %rax 
 mov $varName886, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName887, %rax 
 mov $varName887, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName888, %rax 
 mov $varName888, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName889, %rax 
 mov $varName889, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName890, %rax 
 mov $varName890, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName891, %rax 
 mov $varName891, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName892, %rax 
 mov $varName892, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName893, %rax 
 mov $varName893, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName894, %rax 
 mov $varName894, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName895, %rax 
 mov $varName895, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName896, %rax 
 mov $varName896, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName897, %rax 
 mov $varName897, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName898, %rax 
 mov $varName898, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName899, %rax 
 mov $varName899, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName900, %rax 
 mov $varName900, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName901, %rax 
 mov $varName901, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName902, %rax 
 mov $varName902, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName903, %rax 
 mov $varName903, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName904, %rax 
 mov $varName904, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName905, %rax 
 mov $varName905, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName906, %rax 
 mov $varName906, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName907, %rax 
 mov $varName907, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName908, %rax 
 mov $varName908, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName909, %rax 
 mov $varName909, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName910, %rax 
 mov $varName910, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName911, %rax 
 mov $varName911, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName912, %rax 
 mov $varName912, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName913, %rax 
 mov $varName913, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName914, %rax 
 mov $varName914, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName915, %rax 
 mov $varName915, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName916, %rax 
 mov $varName916, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName917, %rax 
 mov $varName917, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName918, %rax 
 mov $varName918, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName919, %rax 
 mov $varName919, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName920, %rax 
 mov $varName920, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName921, %rax 
 mov $varName921, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName922, %rax 
 mov $varName922, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName923, %rax 
 mov $varName923, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName924, %rax 
 mov $varName924, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName925, %rax 
 mov $varName925, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName926, %rax 
 mov $varName926, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName927, %rax 
 mov $varName927, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName928, %rax 
 mov $varName928, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName929, %rax 
 mov $varName929, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName930, %rax 
 mov $varName930, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName931, %rax 
 mov $varName931, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName932, %rax 
 mov $varName932, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName933, %rax 
 mov $varName933, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName934, %rax 
 mov $varName934, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName935, %rax 
 mov $varName935, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName936, %rax 
 mov $varName936, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName937, %rax 
 mov $varName937, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName938, %rax 
 mov $varName938, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName939, %rax 
 mov $varName939, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName940, %rax 
 mov $varName940, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName941, %rax 
 mov $varName941, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName942, %rax 
 mov $varName942, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName943, %rax 
 mov $varName943, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName944, %rax 
 mov $varName944, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName945, %rax 
 mov $varName945, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName946, %rax 
 mov $varName946, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName947, %rax 
 mov $varName947, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName948, %rax 
 mov $varName948, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName949, %rax 
 mov $varName949, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName950, %rax 
 mov $varName950, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName951, %rax 
 mov $varName951, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName952, %rax 
 mov $varName952, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName953, %rax 
 mov $varName953, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName954, %rax 
 mov $varName954, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName955, %rax 
 mov $varName955, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName956, %rax 
 mov $varName956, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName957, %rax 
 mov $varName957, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName958, %rax 
 mov $varName958, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName959, %rax 
 mov $varName959, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName960, %rax 
 mov $varName960, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName961, %rax 
 mov $varName961, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName962, %rax 
 mov $varName962, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName963, %rax 
 mov $varName963, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName964, %rax 
 mov $varName964, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName965, %rax 
 mov $varName965, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName966, %rax 
 mov $varName966, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName967, %rax 
 mov $varName967, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName968, %rax 
 mov $varName968, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName969, %rax 
 mov $varName969, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName970, %rax 
 mov $varName970, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName971, %rax 
 mov $varName971, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName972, %rax 
 mov $varName972, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName973, %rax 
 mov $varName973, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName974, %rax 
 mov $varName974, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName975, %rax 
 mov $varName975, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName976, %rax 
 mov $varName976, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName977, %rax 
 mov $varName977, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName978, %rax 
 mov $varName978, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName979, %rax 
 mov $varName979, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName980, %rax 
 mov $varName980, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName981, %rax 
 mov $varName981, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName982, %rax 
 mov $varName982, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName983, %rax 
 mov $varName983, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName984, %rax 
 mov $varName984, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName985, %rax 
 mov $varName985, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName986, %rax 
 mov $varName986, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName987, %rax 
 mov $varName987, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName988, %rax 
 mov $varName988, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName989, %rax 
 mov $varName989, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName990, %rax 
 mov $varName990, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName991, %rax 
 mov $varName991, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName992, %rax 
 mov $varName992, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName993, %rax 
 mov $varName993, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName994, %rax 
 mov $varName994, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName995, %rax 
 mov $varName995, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName996, %rax 
 mov $varName996, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName997, %rax 
 mov $varName997, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName998, %rax 
 mov $varName998, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName999, %rax 
 mov $varName999, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1000, %rax 
 mov $varName1000, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1001, %rax 
 mov $varName1001, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1002, %rax 
 mov $varName1002, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1003, %rax 
 mov $varName1003, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1004, %rax 
 mov $varName1004, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1005, %rax 
 mov $varName1005, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1006, %rax 
 mov $varName1006, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1007, %rax 
 mov $varName1007, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1008, %rax 
 mov $varName1008, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1009, %rax 
 mov $varName1009, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1010, %rax 
 mov $varName1010, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1011, %rax 
 mov $varName1011, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1012, %rax 
 mov $varName1012, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1013, %rax 
 mov $varName1013, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1014, %rax 
 mov $varName1014, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1015, %rax 
 mov $varName1015, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1016, %rax 
 mov $varName1016, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1017, %rax 
 mov $varName1017, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1018, %rax 
 mov $varName1018, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1019, %rax 
 mov $varName1019, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1020, %rax 
 mov $varName1020, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1021, %rax 
 mov $varName1021, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1022, %rax 
 mov $varName1022, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1023, %rax 
 mov $varName1023, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1024, %rax 
 mov $varName1024, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1025, %rax 
 mov $varName1025, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1026, %rax 
 mov $varName1026, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1027, %rax 
 mov $varName1027, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1028, %rax 
 mov $varName1028, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1029, %rax 
 mov $varName1029, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1030, %rax 
 mov $varName1030, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1031, %rax 
 mov $varName1031, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1032, %rax 
 mov $varName1032, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1033, %rax 
 mov $varName1033, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1034, %rax 
 mov $varName1034, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1035, %rax 
 mov $varName1035, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1036, %rax 
 mov $varName1036, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1037, %rax 
 mov $varName1037, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1038, %rax 
 mov $varName1038, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1039, %rax 
 mov $varName1039, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1040, %rax 
 mov $varName1040, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1041, %rax 
 mov $varName1041, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1042, %rax 
 mov $varName1042, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1043, %rax 
 mov $varName1043, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1044, %rax 
 mov $varName1044, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1045, %rax 
 mov $varName1045, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1046, %rax 
 mov $varName1046, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1047, %rax 
 mov $varName1047, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1048, %rax 
 mov $varName1048, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1049, %rax 
 mov $varName1049, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1050, %rax 
 mov $varName1050, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1051, %rax 
 mov $varName1051, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1052, %rax 
 mov $varName1052, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1053, %rax 
 mov $varName1053, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1054, %rax 
 mov $varName1054, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1055, %rax 
 mov $varName1055, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1056, %rax 
 mov $varName1056, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1057, %rax 
 mov $varName1057, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1058, %rax 
 mov $varName1058, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1059, %rax 
 mov $varName1059, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1060, %rax 
 mov $varName1060, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1061, %rax 
 mov $varName1061, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1062, %rax 
 mov $varName1062, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1063, %rax 
 mov $varName1063, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1064, %rax 
 mov $varName1064, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1065, %rax 
 mov $varName1065, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1066, %rax 
 mov $varName1066, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1067, %rax 
 mov $varName1067, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1068, %rax 
 mov $varName1068, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1069, %rax 
 mov $varName1069, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1070, %rax 
 mov $varName1070, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1071, %rax 
 mov $varName1071, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1072, %rax 
 mov $varName1072, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1073, %rax 
 mov $varName1073, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1074, %rax 
 mov $varName1074, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1075, %rax 
 mov $varName1075, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1076, %rax 
 mov $varName1076, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1077, %rax 
 mov $varName1077, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1078, %rax 
 mov $varName1078, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1079, %rax 
 mov $varName1079, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1080, %rax 
 mov $varName1080, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1081, %rax 
 mov $varName1081, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1082, %rax 
 mov $varName1082, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1083, %rax 
 mov $varName1083, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1084, %rax 
 mov $varName1084, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1085, %rax 
 mov $varName1085, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1086, %rax 
 mov $varName1086, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1087, %rax 
 mov $varName1087, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1088, %rax 
 mov $varName1088, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1089, %rax 
 mov $varName1089, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1090, %rax 
 mov $varName1090, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1091, %rax 
 mov $varName1091, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1092, %rax 
 mov $varName1092, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1093, %rax 
 mov $varName1093, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1094, %rax 
 mov $varName1094, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1095, %rax 
 mov $varName1095, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1096, %rax 
 mov $varName1096, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1097, %rax 
 mov $varName1097, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1098, %rax 
 mov $varName1098, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1099, %rax 
 mov $varName1099, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1100, %rax 
 mov $varName1100, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1101, %rax 
 mov $varName1101, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1102, %rax 
 mov $varName1102, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1103, %rax 
 mov $varName1103, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1104, %rax 
 mov $varName1104, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1105, %rax 
 mov $varName1105, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1106, %rax 
 mov $varName1106, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1107, %rax 
 mov $varName1107, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1108, %rax 
 mov $varName1108, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1109, %rax 
 mov $varName1109, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1110, %rax 
 mov $varName1110, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1111, %rax 
 mov $varName1111, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1112, %rax 
 mov $varName1112, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1113, %rax 
 mov $varName1113, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1114, %rax 
 mov $varName1114, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1115, %rax 
 mov $varName1115, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1116, %rax 
 mov $varName1116, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1117, %rax 
 mov $varName1117, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1118, %rax 
 mov $varName1118, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1119, %rax 
 mov $varName1119, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1120, %rax 
 mov $varName1120, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName0, %rax 
 mov $varName0, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $60,  %rax
xor %rdi, %rdi
syscall
