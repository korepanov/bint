__throughError:
 mov $fatalError, %rsi
 call __print 
 mov $60, %rax
 mov $1, %rdi
 syscall

 __throughUserError:
 # %rsi - адрес, по которому лежит сообщение об ошибке 
 # puts error message in the error variable
 mov %rsi, (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi  
 call __set

 xor %rax, %rax
 call __setVar

 mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarNamePanic, %rax 
 mov $varNamePanic, %rdi
 call __set
 call __getVar
 mov (userData), %rax 
 mov (%rax), %dl 
 cmp $'0', %dl 
 jz __throughUserErrorEnd 
 
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi
 call __set
 call __getVar
 
 mov (userData), %rsi 
 call __print
 mov $enter, %rsi 
 call __print 
 
 mov $60, %rax
 mov $1, %rdi
 syscall

 __throughUserErrorEnd:
 ret

__print:
 push %rsi 
 push %rdi 
 push %rdx 
 push %rax
 
 __printBegin:
 mov (%rsi), %al	
 cmp $0, %al	
 jz  __printEx			
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall		    
 inc %rsi		  		    
 jnz __printBegin
__printEx:

 pop %rax 
 pop %rdx 
 pop %rdi 
 pop %rsi 
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
 mov (memoryBegin), %r8  
 __printHeapLoop:
 cmp (strMax), %r8 
 jz __printHeapEx
 mov (%r8), %dl 
 /*cmp $0, %dl 
 jnz __printHeapNotZero
 mov $endSymbol, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 __printHeapNotZero: */
 cmp $2, %dl 
 jnz __printHeapNotTwo
 mov $endSymbol, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 __printHeapNotTwo: 

 cmp $1, %dl 
 jnz printHeapNopEnd
 mov $starSymbol, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 printHeapNopEnd:
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

 cmp %rdi, %rdx 
 jz __setEnd

 mov %rdx, %r8 
 mov %rsi, %r9
  
 __setClear:
 cmp $0, %rsi
 jz __setClearEnd
 movb $1, (%rdx)
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
 cmp $1, %rax 
 jz __star
 inc %rdx
 __star: 
 movb $0, (%rdx)

 __setEnd:
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
 
 __userConcatinate:
 # входные параметры 
 # r8 - адрес начала первой строки (либо адрес имени переменной)
 # r9 - адрес начала второй строки (либо адрес имени переменной)
 # rax = 1 - строка слева лежит в переменной
 # rax = 0 - строка слева в статической памяти
 # rbx = 1 - строка справа лежит в переменной 
 # rbx = 0 - строка справа лежит в статичесой памяти 
 # $varName - адрес имени переменной, куда положить результат
 movb $0, (userConcatinateFlag)
 movb $0, (userConcatinateFlag2)
 
 push %r8 
 push %r9 
 push %rax 
 push %rbx 


 mov %r13, %rbx
 __userConcatinateLocal:
 cmp %r15, %rbx
 jg __userConcatinateEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __userConcatinateEnd 
 
 add (varSize), %rbx 
 jmp __userConcatinateLocal
 
 __userConcatinateEnd:
  __userConcatinateSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
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
 jz __userConcatinateSearch
 

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
 jnz __userConcatinateErrEnd
 mov $strError, %rsi 
 call __throughUserError
 ret 

 __userConcatinateErrEnd:

 pop %rbx 
 pop %rax 
 pop %r9 
 pop %r8 



 cmp $0, %rax 
 jz __userConcatinateLeftZero
 cmp $0, %rbx 
 jz __userConcatinateRightZero 
 // слева и справа динамические данные 
 
 push %r8 
 push %r9  
 push %r12 

 call __getVar 
 mov (userData), %rax 
 pop %r12 
 pop %r9 
 pop %r8
 
 push %rax 
 push %r8 
 push %r9 
 push %r12 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov %r8, %rdi
 call __set
 call __getVar 
 mov (userData), %rax 
 
 pop %r12 
 pop %r9 
 pop %r8
 push %rax 
 push %r8 
 push %r9 
 push %r12 
 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov %r9, %rdi
 call __set
 call __getVar 
 mov (userData), %rax 
 pop %r12 # for internalShiftStr 
 pop %r9 
 pop %r8

 push %rax

 pop %rcx # address of the second string 
 pop %rbx # address os the first string 
 pop %rax # address of the result 

 cmp %rax, %rbx 
 jz __userConcatinateTwoOnesTheSame1
 cmp %rax, %rcx 
 jz __userConcatinateTwoOnesTheSame3 
 cmp %rbx, %rcx 
 jz __userConcatinateTwoOnesTheSame4
 // all three variables are different 
 __userConcatinateTwoOnesStart:
 cmp %rax, %rbx 
 jng __userConcatinateTwoOnesFirstFlagEnd
 movb $1, (userConcatinateFlag) # to shift the address of the first string  
 __userConcatinateTwoOnesFirstFlagEnd:
 cmp %rax, %rcx 
 jng __userConcatinateTwoOnesSecondFlagEnd
 movb $1, (userConcatinateFlag2) # to shift the address of the second string 
 __userConcatinateTwoOnesSecondFlagEnd:
 push %rax
 push %rbx 
 push %rcx 
 push %r12 
 
 mov %rax, %rsi 
 call __clear 
 mov %rbx, %rsi 
 call __len 
 push %rax 
 mov %rcx, %rsi 
 call __len 
 mov %rax, %rdx 
 pop %rax 
 add %rax, %rdx 

 pop %r12 
 pop %rcx
 pop %rbx 
 pop %rax 
 
 push %rax  

  __userConcatinateTwoOnesPrepare:
 cmp $0, %rdx  
 jz __userConcatinateTwoOnesPrepareEnd
 mov (%rax), %dil 
 cmp $2, %dil 
 jnz __userConcatinateTwoOnesMoreMemEnd

 mov (userConcatinateFlag), %dil
 cmp $1, %dil 
 jnz __userConcatinateTwoOnesAddEnd
 mov (strValSize), %r8  
 add %r8, %rbx
 __userConcatinateTwoOnesAddEnd:

 mov (userConcatinateFlag2), %dil
 cmp $1, %dil 
 jnz __userConcatinateTwoOnesAddEnd2
 mov (strValSize), %r8  
 add %r8, %rcx
 __userConcatinateTwoOnesAddEnd2:

 push %rax 
 push %rbx 
 push %rcx 
 push %rdx 
 push %r12 

 call __internalShiftStr
 
 pop %r12
 pop %rdx 
 pop %rcx  
 pop %rbx 
 pop %rax 
 
 __userConcatinateTwoOnesMoreMemEnd:
 inc %rax 
 dec %rdx 
 jmp __userConcatinateTwoOnesPrepare
 __userConcatinateTwoOnesPrepareEnd:
 
 pop %rax 

 __userConcatinateTwoOnesNow:
 mov (%rbx), %dil 
 cmp $0, %dil 
 jz __userConcatinateTwoOnesNowEnd
 mov %dil, (%rax)
 inc %rbx 
 inc %rax 
 jmp __userConcatinateTwoOnesNow 
 __userConcatinateTwoOnesNowEnd:
 
 __userConcatinateTwoOnesNow2:
 mov (%rcx), %dil 
 cmp $0, %dil 
 jz __userConcatinateTwoOnesNowEnd2 
 mov %dil, (%rax)
 inc %rcx 
 inc %rax
 jmp __userConcatinateTwoOnesNow2  
 __userConcatinateTwoOnesNowEnd2:
 movb $0, (%rax)

 ret 
 __userConcatinateTwoOnesTheSame4:
 // the first and the second variables are the same
 jmp __userConcatinateTwoOnesStart
 ret 
 __userConcatinateTwoOnesTheSame3:
 // result and the second variable are the same 
 cmp %rax, %rbx 
 jng __userConcatinateTwoOnesTheSameRightFlagEnd
 movb $1, (userConcatinateFlag) # to shift the address of the first string  
 __userConcatinateTwoOnesTheSameRightFlagEnd:
 
 push %rax
 push %rbx 
 push %rcx 
 push %r12 
 
 mov %rbx, %rsi 
 call __len 
 mov %rax, %rdx 
 push %rdx 

 mov %rcx, %rsi 
 call __len 
 mov %rax, %rdi   

 pop %rdx 
 pop %r12 
 pop %rcx
 pop %rbx 
 pop %rax 
 
 push %rax 
 push %rdx 
 push %rdi 

 add %rdi, %rax 

  __userConcatinateTwoOnesTheSameRightPrepare:
 cmp $0, %rdx  
 jz __userConcatinateTwoOnesTheSameRightPrepareEnd
 mov (%rax), %dil 
 cmp $2, %dil 
 jnz __userConcatinateTwoOnesMoreMemTheSameRightEnd

 mov (userConcatinateFlag), %dil
 cmp $1, %dil 
 jnz __userConcatinateTwoOnesTheSameRightAddEnd
 mov (strValSize), %r8  
 add %r8, %rbx
 __userConcatinateTwoOnesTheSameRightAddEnd:

 push %rax 
 push %rbx 
 push %rcx 
 push %rdx 
 push %r12 

 call __internalShiftStr
 
 pop %r12
 pop %rdx 
 pop %rcx  
 pop %rbx 
 pop %rax 
 
 __userConcatinateTwoOnesMoreMemTheSameRightEnd:
 inc %rax 
 dec %rdx 
 jmp __userConcatinateTwoOnesTheSameRightPrepare
 __userConcatinateTwoOnesTheSameRightPrepareEnd:
 
 mov %rbx, %rsi # address of the first string in the variable 
 pop %rax # length of the string in the variable
 pop %rdx # length of the static value 
 pop %rbx # address of the string begin in the variable 

 push %rbx 
 add %rax, %rbx # the end if the string in the variable
 __userConcatinateTwoOnesShiftTheSameRight:
 cmp $0, %rax 
 jl __userConcatinateTwoOnesShiftTheSameRightEnd
 mov (%rbx), %dil
 movb $1, (%rbx)
 add %rdx, %rbx 
 mov %dil, (%rbx)
 sub %rdx, %rbx 
 dec %rbx 
 dec %rax 
 jmp __userConcatinateTwoOnesShiftTheSameRight 
 __userConcatinateTwoOnesShiftTheSameRightEnd:
 pop %rbx 

 __userConcatinateTwoOnesTheSameRightNow:
 mov (%rsi), %dil 
 cmp $0, %rdx 
 jz __userConcatinateTwoOnesTheSameRightNowEnd 
 mov %dil, (%rbx)
 inc %rsi 
 inc %rbx
 dec %rdx  
 jmp __userConcatinateTwoOnesTheSameRightNow 
 __userConcatinateTwoOnesTheSameRightNowEnd:

 ret 
 __userConcatinateTwoOnesTheSame1:
 cmp %rax, %rcx 
 jz __userConcatinateTwoOnesTheSame2
 // result and the first variable are the same 
 cmp %rax, %rcx 
 jng __userConcatinateTwoOnesTheSameLeftFlagEnd
 movb $1, (userConcatinateFlag) # to shift the address of the second string  
 __userConcatinateTwoOnesTheSameLeftFlagEnd:



 push %rax 
 push %rbx 
 push %rcx
 push %r12  

 mov %rax, %rsi 
 call __len 
 mov %rax, %rdx  
 
 pop %r12 
 pop %rcx 
 pop %rbx 
 pop %rax 
 
 push %rbx 
 push %rcx 
 push %r12 

 add %rdx, %rax 
 push %rax 
 mov %rcx, %rsi 
 call __len 
 mov %rax, %rdx 
 
 pop %rax 
 pop %r12 
 pop %rcx 
 pop %rbx 

 push %rax 
 push %rbx 
 push %r12 


 __userConcatinateTwoOnesTheSameLeftPrepare:
 cmp $0, %rdx  
 jz __userConcatinateTwoOnesTheSameLeftPrepareEnd
 mov (%rax), %dil 
 cmp $2, %dil 
 jnz __userConcatinateTwoOnesMoreMemTheSameLeftEnd

 mov (userConcatinateFlag), %dil
 cmp $1, %dil 
 jnz __userConcatinateTwoOnesTheSameLeftAddEnd
 mov (strValSize), %r8  
 add %r8, %rcx
 __userConcatinateTwoOnesTheSameLeftAddEnd:

 push %rax 
 push %rbx 
 push %rcx 
 push %rdx 
 push %r12 

 call __internalShiftStr
 
 pop %r12
 pop %rdx 
 pop %rcx  
 pop %rbx 
 pop %rax 
 
 __userConcatinateTwoOnesMoreMemTheSameLeftEnd:
 inc %rax 
 dec %rdx 
 jmp __userConcatinateTwoOnesTheSameLeftPrepare
 __userConcatinateTwoOnesTheSameLeftPrepareEnd:

 pop %r12 
 pop %rbx 
 pop %rax

 __userConcatinateTwoOnesTheSameLeftNow:
 mov (%rcx), %dil 
 cmp $0, %dil 
 jz __userConcatinateTwoOnesTheSameLeftNowEnd 
 mov %dil, (%rax)
 inc %rcx 
 inc %rax 
 jmp __userConcatinateTwoOnesTheSameLeftNow
 __userConcatinateTwoOnesTheSameLeftNowEnd:
 movb $0, (%rax)

 ret  
 __userConcatinateTwoOnesTheSame2:
 // all three variables are the same 
 push %rax 
 mov %rax, %rsi 
 call __len 
 mov %rax, %rbx # the length 
 pop %rax #from here to read 
 mov %rax, %rcx 
 add %rbx, %rcx # here to write  


 push %rax 
 push %rbx 
 push %rcx 
 push %r12 

 mov %rbx, %rdx 
 add %rbx, %rax 

 __userConcatinateTwoOnesTheSamePrepare2:
 cmp $0, %rdx  
 jz __userConcatinateTwoOnesTheSamePrepareEnd2
 mov (%rax), %dil 
 cmp $2, %dil 
 jnz __userConcatinateTwoOnesMoreMemTheSameEnd2

 push %rax 
 push %rbx 
 push %rcx 
 push %rdx 
 push %r12 

 call __internalShiftStr
 
 pop %r12
 pop %rdx 
 pop %rcx  
 pop %rbx 
 pop %rax 
 
 __userConcatinateTwoOnesMoreMemTheSameEnd2:
 inc %rax 
 dec %rdx 
 jmp __userConcatinateTwoOnesTheSamePrepare2
 __userConcatinateTwoOnesTheSamePrepareEnd2:

 pop %r12 
 pop %rcx 
 pop %rbx 
 pop %rax

 __userConcatinateTwoOnesTheSameNow2:
 cmp $0, %rbx 
 jz __userConcatinateTwoOnesTheSameNowEnd2
 mov (%rax), %dil 
 mov %dil, (%rcx)
 inc %rax 
 inc %rcx 
 dec %rbx 
 jmp __userConcatinateTwoOnesTheSameNow2 
 __userConcatinateTwoOnesTheSameNowEnd2:
 movb $0, (%rcx) 

 ret 

 __userConcatinateRightZero:
 // слева динамические данные, а справа статические 
 push %r8 
 push %r9 
 push %r12 

 call __getVar 
 mov (userData), %rbx 

 pop %r12  
 pop %r9 
 pop %r8 


 push %r8 
 push %r9 
 push %rbx
 push %r12 

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov %r8, %rdi
 call __set
 call __getVar 


 mov (userData), %rax 
 
 pop %r12 
 pop %rbx 
 pop %r9 
 pop %r8 

 push %r8 
 push %r9 
 push %rax 
 push %rbx 
 push %r12 

 # %rbx - dest
 # %rax - source 
 cmp %rbx, %rax 
 jz __userConcatinateRightZeroTheSame 
 jg __userConcatinateRightZeroShift 
 jmp __userConcatinateRightZeroShiftEnd 

 __userConcatinateRightZeroTheSame:
 # the same variable 
 pop %r12 
 pop %rbx 
 pop %rax 
 pop %r9 
 pop %r8 
 push %r12 

 mov (userData), %rsi 
 call __len
 pop %r12  

 add %rax, %rbx 
 push %rbx

 mov %r9, %rsi 
 call __len 

 pop %rbx 
 
 push %rbx 
 __userConcatinateRightZeroTheSamePrepare:
 cmp $0, %rax 
 jz __userConcatinateRightZeroTheSamePrepareEnd
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateRightZeroTheSameMoreMemEnd

 push %rax 
 push %rbx 
 push %r12 
 push %r9 
 call __internalShiftStr
 pop %r9 
 pop %r12 
 pop %rbx 
 pop %rax 
 
 __userConcatinateRightZeroTheSameMoreMemEnd:
 inc %rbx 
 dec %rax 
 jmp __userConcatinateRightZeroTheSamePrepare
 __userConcatinateRightZeroTheSamePrepareEnd:
 
 pop %rbx 
 mov %rbx, %rcx 
 mov %r9, %rbx 

 __userConcatinateRightZeroTheSameNow:
 mov (%rbx), %dil 
 cmp $0, %dil 
 jz __userConcatinateRightZeroTheSameEnd 
 mov %dil, (%rcx)
 inc %rbx 
 inc %rcx  
 jmp __userConcatinateRightZeroTheSameNow
 __userConcatinateRightZeroTheSameEnd:
 movb $0, (%rcx)
 ret  
 __userConcatinateRightZeroShift:
 movb $1, (userConcatinateFlag) # add size of the shift to the source 
 __userConcatinateRightZeroShiftEnd:

 mov %rax, %rsi 
 call __len 
 mov %rax, %rcx 
 
 pop %r12 
 pop %rbx 
 pop %rax 
 pop %r9 
 pop %r8

 push %r8 
 push %r9 
 push %rax 
 push %rbx
 push %rcx 
 push %r12 

 mov %rbx, %rsi 
 call __clear 

 mov %r9, %rsi 
 call __len 
  
 pop %r12 
 pop %rcx 
 add %rax, %rcx # length of the result 
 inc %rcx # 0 byte 
 pop %rbx
 pop %rax  
 push %rbx 
 push %r12 
 pop %r12 

 __userConcatinateRightZeroPrepare:
 cmp $0, %rcx 
 jz __userConcatinateRightZeroPrepareEnd 
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateRightZeroMoreMemEnd

 mov (userConcatinateFlag), %dil
 cmp $1, %dil 
 jnz __userConcatinateRightZeroAddEnd
 mov (strValSize), %rdx   
 add %rdx, %rax  
 __userConcatinateRightZeroAddEnd:
 push %r8 
 push %r9 
 push %rax 
 push %rbx
 push %rcx 
 push %r12 
 call __internalShiftStr
 pop %r12 
 pop %rcx
 pop %rbx
 pop %rax 
 pop %r9 
 pop %r8  
 __userConcatinateRightZeroMoreMemEnd:
 inc %rbx 
 dec %rcx 
 jmp __userConcatinateRightZeroPrepare
 __userConcatinateRightZeroPrepareEnd:
 
 pop %rbx 

 __userConcatinateRightZeroFirst:
 mov (%rax), %dil 
 cmp $0, %dil 
 jz __userConcatinateRightZeroFirstEnd
 mov %dil, (%rbx)
 inc %rax 
 inc %rbx 
 jmp __userConcatinateRightZeroFirst
 __userConcatinateRightZeroFirstEnd: 
 pop %r9
 pop %r8 

 __userConcatinateRightZeroSecond:
 mov (%r9), %dil 
 cmp $0, %dil 
 jz __userConcatinateRightZeroSecondEnd
 mov %dil, (%rbx)
 inc %r9 
 inc %rbx 
 jmp __userConcatinateRightZeroSecond 
 __userConcatinateRightZeroSecondEnd:
 movb $0, (%rbx)

 ret 
 __userConcatinateLeftZero:
 cmp $0, %rbx 
 jz __userConcatinateTwoZeros
 // слева статические данные, а справа динамические 
 
 push %r8 
 push %r9 
 push %r12 

 call __getVar 
 mov (userData), %rbx 

 pop %r12  
 pop %r9 
 pop %r8 


 push %r8 
 push %r9 
 push %rbx
 push %r12 

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov %r9, %rdi
 call __set
 call __getVar 


 mov (userData), %rax 
 mov %rax, %rsi 
 
 pop %r12 
 pop %rbx 
 pop %r9 
 pop %r8 

 push %r8 
 push %r9 
 push %rax 
 push %rbx 
 push %r12 

 # %rbx - dest
 # %rax - source 
 cmp %rbx, %rax 
 jz __userConcatinateLeftZeroTheSame 
 jg __userConcatinateLeftZeroShift 
 jmp __userConcatinateLeftZeroShiftEnd 

__userConcatinateLeftZeroTheSame:

# the same variable 
 pop %r12 
 pop %rbx 
 pop %rax 
 pop %r9 
 pop %r8 
 push %r12 

 mov (userData), %rsi 
 
 call __len
 
 pop %r12  
 push %rax 
 
 push %rbx 
 add %rax, %rbx 
 push %rbx

 mov %r8, %rsi 
 call __len 
 pop %rbx 
 push %rax
 
 __userConcatinateLeftZeroTheSamePrepare:
 cmp $0, %rax 
 jz __userConcatinateLeftZeroTheSamePrepareEnd
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateLeftZeroTheSameMoreMemEnd

 push %rax 
 push %rbx 
 push %r12 
 push %r8 
 push %r9 
 call __internalShiftStr
 pop %r9 
 pop %r8 
 pop %r12 
 pop %rbx 
 pop %rax 
 
 __userConcatinateLeftZeroTheSameMoreMemEnd:
 inc %rbx 
 dec %rax 
 jmp __userConcatinateLeftZeroTheSamePrepare
 __userConcatinateLeftZeroTheSamePrepareEnd:
 pop %rdx # length of the static value 
 pop %rbx # address of the string begin in the variable 
 pop %rax # length of the string in the variable

 push %rbx 
 add %rax, %rbx # the end if the string in the variable
 __userConcatinateLeftZeroShiftTheSame:
 cmp $0, %rax 
 jl __userConcatinateLeftZeroShiftTheSameEnd
 mov (%rbx), %dil
 movb $1, (%rbx)
 add %rdx, %rbx 
 mov %dil, (%rbx)
 sub %rdx, %rbx 
 dec %rbx 
 dec %rax 
 jmp __userConcatinateLeftZeroShiftTheSame
 __userConcatinateLeftZeroShiftTheSameEnd:
 pop %rbx 

 __userConcatinateLeftZeroTheSameNow:
 cmp $0, %rdx 
 jz __userConcatinateLeftZeroTheSameNowEnd 
 mov (%r8), %dil 
 mov %dil, (%rbx) 
 inc %r8
 inc %rbx 
 dec %rdx  
 jmp __userConcatinateLeftZeroTheSameNow
 __userConcatinateLeftZeroTheSameNowEnd:

 ret 

__userConcatinateLeftZeroShift:
movb $1, (userConcatinateFlag)
__userConcatinateLeftZeroShiftEnd:
 
 mov %rax, %rsi
 call __len 
 mov %rax, %rcx 
 
 pop %r12 
 pop %rbx 
 pop %rax 
 pop %r9 
 pop %r8

 push %r8 
 push %r9 
 push %rax 
 push %rbx
 push %rcx 
 push %r12 

 mov %rbx, %rsi 
 call __clear 

 mov %r8, %rsi 
 call __len 
  
 pop %r12 
 pop %rcx 
 add %rax, %rcx # length of the result 
 inc %rcx # 0 byte 
 pop %rbx
 pop %rax  
 push %rbx 
 push %r12 
 pop %r12 


 __userConcatinateLeftZeroPrepare:
 cmp $0, %rcx 
 jz __userConcatinateLeftZeroPrepareEnd 
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateLeftZeroMoreMemEnd

 mov (userConcatinateFlag), %dil
 cmp $1, %dil 
 jnz __userConcatinateLeftZeroAddEnd
 mov (strValSize), %rdx   
 add %rdx, %rax  
 __userConcatinateLeftZeroAddEnd: 
 push %r8 
 push %r9 
 push %rax 
 push %rbx
 push %rcx 
 push %r12 
 call __internalShiftStr
 pop %r12 
 pop %rcx
 pop %rbx
 pop %rax 
 pop %r9 
 pop %r8  
 __userConcatinateLeftZeroMoreMemEnd:
 inc %rbx 
 dec %rcx 
 jmp __userConcatinateLeftZeroPrepare
 __userConcatinateLeftZeroPrepareEnd:
 
 pop %rbx 

 pop %r9
 pop %r8 

 __userConcatinateLeftZeroFirst:
 mov (%r8), %dil 
 cmp $0, %dil 
 jz __userConcatinateLeftZeroFirstEnd
 mov %dil, (%rbx)
 inc %r8 
 inc %rbx 
 jmp __userConcatinateLeftZeroFirst 
 __userConcatinateLeftZeroFirstEnd:
 movb $0, (%rbx)
 
 __userConcatinateLeftZeroSecond:
 mov (%rax), %dil 
 cmp $0, %dil 
 jz __userConcatinateLeftZeroSecondEnd
 mov %dil, (%rbx)
 inc %rax 
 inc %rbx 
 jmp __userConcatinateLeftZeroSecond 
 __userConcatinateLeftZeroSecondEnd: 
 movb $0, (%rbx)

 ret 
 __userConcatinateTwoZeros:
 // слева и справа статические данные 
 push %r8 
 push %r9 
 push %r12 
 call __getVar 
 mov (userData), %rsi 
 call __clear
 pop %r12  
 pop %r9 
 pop %r8 
 
 push %r8 
 push %r9 
 push %r12 

 mov %r8, %rsi
 call __len 

 pop %r12 
 pop %r9
 pop %r8 

 push %r8 
 push %r9  
 push %rax 
 push %r12 

 mov %r9, %rsi 
 call __len 

 pop %r12 
 pop %rbx # old %rax  

 add %rbx, %rax 
 inc %rax # 0 byte 
 push %rax 
 push %r12 

 call __getVar

 pop %r12  
 pop %rax 

 mov (userData), %rbx 
 push %rbx 

 __userConcatinateTwoZerosPrepare:
 cmp $0, %rax 
 jz __userConcatinateTwoZerosNow
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateTwoZerosMoreMemEnd
 __userConcatinateTwoZerosMoreMem: 
 push %rax 
 push %rbx 
 push %r12 
 call __internalShiftStr
 pop %r12 
 pop %rbx 
 pop %rax 
 __userConcatinateTwoZerosMoreMemEnd:
 dec %rax 
 inc %rbx 
 jmp __userConcatinateTwoZerosPrepare

 __userConcatinateTwoZerosNow:
 
 pop %rbx 
 pop %r9 
 pop %r8 
 
 __userConcatinateTwoZerosFirst:
 mov (%r8), %dil 
 cmp $0, %dil 
 jz __userConcatinateTwoZerosFirstEnd
 mov %dil, (%rbx)
 inc %r8 
 inc %rbx 
 jmp __userConcatinateTwoZerosFirst
 __userConcatinateTwoZerosFirstEnd:
 
 __userConcatinateTwoZerosSecond:
 mov (%r9), %dil 
 cmp $0, %dil 
 jz __userConcatinateTwoZerosSecondEnd 
 mov %dil, (%rbx)
 inc %r9 
 inc %rbx 
 jmp __userConcatinateTwoZerosSecond
 __userConcatinateTwoZerosSecondEnd:
 movb $0, (%rbx)

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

 __userToNumber:
 # вход: buf 
 # выход:  %rax 
 mov $buf, %rdx # our string
 movzx (%rdx), %rcx 
 cmp $'-', %rcx 
 jnz __userToNumberAtoi
 inc %rdx 
 __userToNumberAtoi:
 xor %rax, %rax # zero a "result so far"
 __userToNumberTop:
 movzx (%rdx), %rcx # get a character
 inc %rdx # ready for next one
 cmp $0, %rcx # end?
 jz __userToNumberDone
 cmp $'.', %rcx # дробная точка?
 jz __userToNumberDone  
 cmp $48, %rcx 
 jl __userToNumberException
 cmp $57, %rcx 
 jg __userToNumberException
 sub $48, %rcx # "convert" character to number
 imul $10, %rax # multiply "result so far" by ten
 add %rcx, %rax # add in current digit
 jmp __userToNumberTop # until done
 __userToNumberDone:
 mov $buf, %rdx 
 movzx (%rdx), %rcx 
 cmp $'-', %rcx 
 jnz __userToNumberIsPos
 mov $0, %rdx 
 sub $1, %rdx 
 //mov %rdx, (buf)
 imul %rdx 
 __userToNumberIsPos:
 ret
 __userToNumberException:
 mov $parseNumberError, %rsi 
 call __throughUserError 
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
 cmp $1, %r11b
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
 cmp $1, %r11b 
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
 add (valSize), %rax
 dec %rax  
 movb $2, (%rax) # признак конца поля для строки 
 inc %rax 
 sub (valSize), %rax 
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


  __clearVars:
 mov %r13, %rbx
 __clearVarsLocal:
 cmp %r15, %rbx
 jg __clearVarsEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __clearVarsEnd 
 
 push %rbx 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenClearSymbol, %rax 
 mov $clearSymbol, %rdi 
 call __set
 call __compare
 pop %rbx  
 cmp $1, %rax 
 jz __clearThis 

 add (varSize), %rbx 
 jmp __clearVarsLocal  
  

 __clearThis:
 push %rbx 
 
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
 jz __clearVarsStr 

 pop %rbx  
 mov %rbx, %rax 
 add (varSize), %rbx 

 __clearVarsCompress: 
 cmp %rbx, %r15 
 jl __clearVarsCompressEnd
 mov (%rbx), %dil 
 mov %dil, (%rax)
 inc %rax 
 inc %rbx 

 jmp __clearVarsCompress
 __clearVarsCompressEnd:
 
 sub (varSize), %r14 
 
 jmp __clearVars  

 __clearVarsStr:
 
 pop %rbx 
 add (varNameSize), %rbx 
 add (typeSize), %rbx 

 push %rbx 
 mov %rbx, %r12 
 call __read 
 call __toNumber 

 push %rax 
 __clearVarsStrNow:
 mov (%rax), %dil 
 cmp $2, %dil 
 jz __clearVarsStrNowEnd
 movb $'!', (%rax)
 inc %rax  
 jmp __clearVarsStrNow 
 __clearVarsStrNowEnd: 
 movb $'!', (%rax)
 inc %rax 
 
 pop %rbx #begin 
 
 mov (strMax), %rcx 

 mov %rax, %rdx 
 sub %rbx, %rdx # size to shift 
 
 __clearVarsStrCompress:
 cmp %rax, %rcx 
 jl __clearVarsStrCompressEnd
 mov (%rax), %dil 
 mov %dil, (%rbx)
 inc %rax 
 inc %rbx 
 jmp __clearVarsStrCompress
 __clearVarsStrCompressEnd:
 
 pop %rbx # place in the variables table
 push %rdx 

 mov %rbx, %rax  
 sub (typeSize), %rax 
 sub (varNameSize), %rax # begin  
 push %rax 

 mov %rbx, %rcx 
 add (valSize), %rcx 

 __clearVarsStrVarTable:
 cmp %rbx, %rcx 
 jz __clearVarsStrVarTableEnd
 movb $'!', (%rbx)
 inc %rbx 
 jmp __clearVarsStrVarTable
 __clearVarsStrVarTableEnd:


 __clearVarsCompress2: 
 cmp %rbx, %r15 
 jl __clearVarsCompressEnd2
 mov (%rbx), %dil 
 mov %dil, (%rax)
 inc %rax 
 inc %rbx 

 jmp __clearVarsCompress2
 __clearVarsCompressEnd2:

 sub (varSize), %r14 

 mov (strPointer), %rbx 
 sub %rdx, %rbx 
 mov %rbx, (strPointer)

 pop %rbx # from here change next addresses for the strings 
 pop %rdx # size to change address 
   
 add (varNameSize), %rbx   

 __clearVarsChangeAddr: 
 mov (%rbx), %dil 
 cmp $1, %dil
 jz __clearVarsChangeAddrEnd
 mov %rbx, %r12 
 call __read 

 push %rbx
 push %rdx 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set 
 call __compare
 pop %rdx 
 pop %rbx
   
 cmp $1, %rax
 jnz __clearVarsChangeAddrNowEnd
 __clearVarsChangeAddrNow:
 add (typeSize), %rbx 
 mov %rbx, %r12 
 call __read  
 push %rbx 
 push %rdx 
 call __toNumber
 pop %rdx 
 pop %rbx 

 sub %rdx, %rax
 
 push %rbx 
 push %rdx 
 call __toStr
 pop %rdx 
 pop %rbx 

 mov %rbx, %rsi 
 call __clear 
 mov $buf2, %rax 
 
 push %rbx 
 __clearVarsSetAddrNow:
 mov (%rax), %dil 
 cmp $0, %dil 
 jz __clearVarsSetAddrNowEnd 
 mov %dil, (%rbx)
 inc %rax 
 inc %rbx 
 jmp __clearVarsSetAddrNow
 __clearVarsSetAddrNowEnd:
 movb $0, (%rbx)
 pop %rbx 
 sub (typeSize), %rbx 

 __clearVarsChangeAddrNowEnd:
 
 add (varSize), %rbx 
 jmp __clearVarsChangeAddr 
 __clearVarsChangeAddrEnd: 
 
 jmp __clearVars 

 __clearVarsEnd:
 ret 



__undefineVar:
 # вход: имя переменной по адресу $varName 
 mov %r13, %rbx
 __undefVarLocal:
 cmp %r15, %rbx
 jg __undefVarEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __undefVarEnd 
 
 add (varSize), %rbx 
 jmp __undefVarLocal  
  
 __undefVarEnd:
 __undefVarSearch:
 sub (varSize), %rbx 
 cmp %rbx, %r13 
 jg __undefEnd2 

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __undefEnd
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
 jz __undefVarSearch
 
 push %rbx 
 mov %rbx, %r12 
 add (varNameSize), %r12

 push %rbx 

 __undefName: 
 cmp %rbx, %r12 
 jz __undefNameEx
 movb $1, (%rbx)
 inc %rbx 
 jmp __undefName 
 __undefNameEx:
 
 pop %rbx 
 # from here set the mark for clearVars
 mov $clearSymbol, %rax 

 __undefSetMark:
 mov (%rax), %dil 
 cmp $0, %dil   
 jz __undefSetMarkEnd 
 movb %dil, (%rbx)
 inc %rbx 
 inc %rax 
 jmp __undefSetMark  
 __undefSetMarkEnd:
 movb $0, (%rbx)

 __undefEnd:

 pop %rax # begin
 
 __undefEnd2:
 ret 

# r12 - pointer (общего назначения)
# r13 - heapBegin 
# r14 - heapPointer 
# r15 - heapMax 

__firstMem:
 # в %rax адрес начала выделяемой памяти 
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
 mov $1, %dl
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

 __newLabelMem:
 # в %rax адрес начала выделяемой памяти 
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, %r9 
 add (labelSize), %r9
# выделить динамическую память
 mov (labelSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __newLabelMemlo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (labelSize), %rbx
 jz  __newLabelex
 jmp __newLabelMemlo
 __newLabelex:
 ret

 __firstStrMem:
 # адрес начала области для выделения памяти
 mov %r15, %rax
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, (strBegin)
 mov %rax, (strPointer)
 mov %rax, %r9 
 add (shiftSize), %r9
 mov %r9, (strMax)
# выделить динамическую память
 mov (shiftSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __firstStrMemLo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (shiftSize), %rbx
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
 mov $1, %dl
 mov $0, %rbx
 __newMemlo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (pageSize), %rbx
 jz  __newMemEx
 jmp __newMemlo
 __newMemEx:
 mov (strMax), %r8 
 call __newStrMem
 call __shiftStr

  mov (strValSize), %r9
 add (strValSize), %r9 
 
 mov (strValSizeMax), %r8 
 cmp %r9, %r8 
 jl __newMemNoMoreStrValSize 
 mov %r9, (strValSize)
 __newMemNoMoreStrValSize:

 mov (pageSize), %r9 
 add (pageSize), %r9 

 mov (pageSizeMax), %r8 
 cmp %r9, %r8 
 jl __newMemNoMorePageSize
 mov %r9, (pageSize)
 __newMemNoMorePageSize:

 mov (pageSize), %r9 
 add (pageSize), %r9 

 mov (shiftSizeMax), %r8
 cmp %r9, %r8 
 jl __newMemNoMoreShiftSize
 mov %r9, (shiftSize)
 __newMemNoMoreShiftSize:
 ret  

 __newStrMem:
 # адрес начала выделяемой памяти в  %r8 
# запомнить адрес начала выделяемой памяти
 #mov %r8, %r14
 mov %r8, %r9 
 #add (pageSize), %r9
 
 mov (strPageNumber), %rax
 add $1, %rax 
 mov (shiftSize), %rdi
 __newStrMemOkBegin: 
 cmp $0, %rax 
 jz __newStrMemOk
 dec %rax 
 add (shiftSize), %rdi  
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
 mov $1, %dl
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
 push %r8 
 push %r10 
 push %rsi 
 push %r9

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
 cmp $1, %r9b  
 jz __readEx
 cmp $0, %r9b  
 jz __readEx
 mov %r9b, (%r10)
 inc %r10 
 inc %r8 
 jmp __readLocal
 __readEx:
 mov $buf, %r10 
 cmp $0, (%r10)
 jnz __readOk
 movb $1, (%r10)
 __readOk:

 pop %r9
 pop %rsi 
 pop %r10 
 pop %r8 
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
 cmp %r11, %r12  
 jge __renewStrEnd
 jmp __renewFindStr

 __renewVal:
 add (typeSize), %r12 
 call __read 
 call __toNumber
 __renewValLocal:
 mov (%rax), %r10b 
 cmp $2, %r10b 
 jz __renewValEnd
 movb $1, (%rax)
 inc %rax 
 jmp __renewValLocal 
 __renewValEnd: 
 movb $1, (%rax)
 __renewAddr:
 call __read 
 mov $buf, %rsi
 call __toNumber 
 #mov (pageSize), %rax 
 #call __toStr
 add (shiftSize), %rax 
 call __toStr # в buf2 новый адрес строки 

 mov %r12, %rsi 
 __renewAddrLocal:
 mov (%rsi), %r10b 
 cmp $1, %r10b 
 jz __renewAddrEnd
 movb $1, (%rsi)
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
 sub (typeSize), %r12 
 add (varSize), %r12 
 jmp __renewFindStr
 __renewStrEnd:
  
 ret 

 __shiftStr:
 # формируем в %r10 адрес нового начала
 mov (strBegin), %r10 
 add (shiftSize), %r10
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
 add (shiftSize), %r10 
 mov %r10, (strPointer)

 mov (strBegin), %r10 
 add (shiftSize), %r10 
 mov %r10, (strBegin)

 mov (strMax), %r10 
 add (shiftSize), %r10 
 mov %r10, (strMax)
 
 call __renewStr 
 
 mov (strMax), %r10 
 add (deltaSize), %r10 
 mov %r10, (strMax)

 mov (shiftSize), %r10 
 add (deltaSize), %r10 
 mov %r10, (shiftSize)

 mov (strMax), %r10 
 sub (strBegin), %r10 

 xor %rdx, %rdx 
 mov %r10, %rax 
 mov $2, %rbx
 div %rbx 

 mov %rax, (deltaSize)
 
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

__internalMakeShiftStr:
# %r12 - адрес внутри таблицы строк, начиная с которого нужно сделать сдвиг 
 mov %r12, %rax # начиная с этого адреса нужно сдвинуть непосредственно строки на valSize   
 mov %rax, (mem11)
 # формируем адрес нового конца 
 mov (strPointer), %r10 
 add (strValSize), %r10
 mov %r10, (mem12)
 cmp (strMax), %r10 
 jl __internalMakeShiftStrOk
 mov (strMax), %r8 
 call __newStrMem
 mov (pageSize), %r8 
 add %r8, (strMax) 
 __internalMakeShiftStrOk:
 # адрес старого конца
 mov (strPointer), %r11
 mov (mem12), %r10 
 mov (mem11), %rax  
 mov %rax, %r9
 __internalShiftMake: 

 cmp %r9, %r11   
 jl __internalShiftMakeEnd
 movb (%r11), %al 
 movb %al, (%r10)
 movb $1, (%r11)

 dec %r10
 dec %r11 
 jmp __internalShiftMake
 __internalShiftMakeEnd:
 
 ret 


__internalShiftStr:
# %r12 - место внутри таблицы переменных, после которого нужно сделать сдвиг  
mov (strValSize), %rsi 
add %rsi, (strPointer) 

mov %r12, (mem9)
call __read

call __toNumber
mov (mem9), %r12 

mov %rax, (mem9) # с этого адреса нужно сделать сдвиг в таблице строк 

mov %r12, %rsi 
sub (typeSize), %rsi 
add (varSize), %rsi

__internalShiftStrLocal: 
cmp %r14, %rsi  
 
jge __internalShiftStrEnd

mov %rsi, (mem)
mov %rax, (mem2)

mov %rsi, %rdi  
call __len 
mov $lenBuf2, %rsi 
mov $buf2, %rdx 
mov $lenVarType, %rax 
call __set 

mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenStringType, %rax 
mov $stringType, %rdi 
call __set 

call __compare 
cmp $1, %rax 

jnz __internalShiftStrChangeNext

mov (mem), %rsi 
mov (mem2), %rax

mov (mem), %rdi 
add (typeSize), %rsi 
add (typeSize), %rdi 
add (valSize), %rdi 
mov %rdi, (mem6)
mov %rsi, (mem7)
mov %rsi, %r12 
call __read 
call __toNumber
mov %rax, (mem8)
add (strValSize), %rax 

call __toStr 
mov (mem6), %rdi
mov (mem7), %rsi   
__internalShiftStrClear:
cmp %rdi, %rsi 
jge __internalShiftStrClearEnd
movb $1, (%rsi)
inc %rsi 
jmp __internalShiftStrClear

__internalShiftStrClearEnd:
dec %rsi 
#movb $0, (%rsi) 

mov (mem), %rsi 
add (typeSize), %rsi 
 
mov $buf2, %rdi 
 
__internalShiftStrSet:
mov (%rdi), %al 
cmp $0, %al 
jz __internalShiftStrChangeEnd
mov (%rdi), %al 
mov %al, (%rsi)
inc %rdi 
inc %rsi 
jmp __internalShiftStrSet

__internalShiftStrChangeEnd:
movb $0, (%rsi)

__internalShiftStrChangeNext:
mov (mem), %rsi 
mov (mem2), %rax
add (varSize), %rsi 

mov %rsi, (mem)
mov %rax, (mem2)

mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenBuf2, %rax 
mov $buf2, %rdi 
call __set 

call __toNumber
cmp (strMax), %rax 
mov (mem), %rsi 
mov (mem2), %rax 
jl __internalShiftStrOk 
mov (strMax), %r8 
call __newStrMem
mov (pageSize), %r8
add %r8, (strMax)
mov (mem), %rsi 
mov (mem2), %rax 
 
__internalShiftStrOk:

jmp  __internalShiftStrLocal

__internalShiftStrEnd:
mov (mem9), %r12 

__internalMakeShiftStrNowLoop:
mov (%r12), %dl 
cmp $2, %dl 
jz __internalMakeShiftStrNow
inc %r12 

jmp __internalMakeShiftStrNowLoop
__internalMakeShiftStrNow:


call __internalMakeShiftStr
ret 


__setVar:
 # вход: 
 # имя переменной по адресу $varName 
 # данные по указателю в (userData) - статика 
 # адрес имени переменной в (userData) - данные в переменной 
 # rax = 0 - статические данные 
 # rax = 1 - данные в переменной
 cmp $1, %rax 
 jnz __setVarStat

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf, %rax 
 mov (userData), %rdi 
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare 

 cmp $1, %rax 
 jnz __setVarStat0 
 ret # присвоение переменной в саму себя 
 __setVarStat0:
 mov $1, %rax 
 __setVarStat:
 mov %rax, (mem4) 

 mov %r13, %rbx
 __setVarLocal:
 cmp %r15, %rbx
 jg __setVarEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __setVarEnd 
 
 add (varSize), %rbx 
 jmp __setVarLocal  
  
 __setVarEnd:
 __setVarSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
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
 movb $1, (%rbx) 
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
 mov %rbx, %r10 

 __setVarClearStr:
 mov (%rbx), %dil 
 cmp $2, %dil 
 jz __setVarClearStrEnd
 movb $1, (%rbx)
 inc %rbx 
 jmp __setVarClearStr
 __setVarClearStrEnd:
 dec %rbx 
 movb $0, (%rbx)
 mov %r10, %rbx 
 mov %rbx, (mem5)
 mov %r12, (mem10)

 __setVarIsNotStr:

 cmp $1, (mem4) # переменная? 
 jnz __setVarNotVar
 
 mov %rbx, (mem19) # сохраним %rbx 

 mov $lenMem18, %rsi # сохранить varName 
 mov $mem18, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenUserData, %rax 
 mov (userData), %rdi 
 call __set
 
 call __getVar
 

 mov (userData), %rsi 
 call __len  
 mov (mem5), %rbx 
 mov (mem10), %r12 

 __setVarPrepare:
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __setVarMoreMemEnd0
 mov %rax, (mem4)
 mov %rbx, (mem5) 
 mov %r12, (mem10)
 call __internalShiftStr
 mov (mem10), %r12 
 mov (mem4), %rax
 mov (mem5), %rbx
 __setVarMoreMemEnd0:
 cmp $0, %rax 
 jz __setVarPrepareEnd
 inc %rbx 
 dec %rax 
 jmp __setVarPrepare
 __setVarPrepareEnd:
 call __getVar # обновим указатель в (userData) 
 mov (mem19), %rbx # восстановим первоначальный %rbx 

 __setVarNotStr0:
 mov (userData), %rax 
 __setNow0:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __setVarRet0
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 jmp __setNow0

 __setVarRet0:
 movb $0, (%rbx)
 mov $lenVarName, %rsi # восстановим varName 
 mov $varName, %rdx 
 mov $lenMem18, %rax 
 mov $mem18, %rdi 
 call __set
 ret
 __setVarNotVar: 
 mov (userData), %rax 
 __setNow:
   
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __setVarMoreMemEnd
 mov %rax, (mem4)
 mov %rbx, (mem5) 
 mov %r12, (mem10)
 call __internalShiftStr
 mov (mem10), %r12 
 mov (mem4), %rax
 mov (mem5), %rbx 
 
 __setVarMoreMemEnd:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __setVarRet 
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 
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
 cmp $1, (buf)
 jz __getVarEnd 
 
 add (varSize), %rbx 
 jmp __getVarLocal  
  
 __getVarEnd:
 __getVarSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
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
movb $1, (%rsi)
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
movb $1, (%rsi)
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
movb $1, (%rsi)
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
movb $1, (%rsi)
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
movb $1, (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal4

__clearBufEnd4:
movb $0, (%rsi)
ret

__clear:
# %rsi - adress from witch to clear everything until 0 byte  
__clearLocal: 
mov (%rsi), %al
cmp $0, %al  
jz __clearEnd
movb $1, (%rsi)
inc %rsi 
jmp __clearLocal

__clearEnd:
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
 ret 

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
 ret
__divINeg:
 mov $divINegError, %rsi 
 call __throughUserError
 ret 

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
 ret 
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
 fstp (buf2)
 
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
 ret 
 __powZeroExpEnd: 
 movss (zero), %xmm2 
 movss %xmm0, %xmm3 
 cmpss $1, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax
 jz __powNegExpEnd
 mov $powZeroNegError, %rsi 
 call __throughUserError
 ret
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
fstp (buf)
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

__userParseFloat:
# $buf - источник (строка)
# %xmm0 - результат
mov $buf, %rax 

mov (%rax), %dl 
cmp $'.', %dl 
jz __userParseFloatException

__userParseFloatCheckPoint:
cmp $0, %dl 
jz __userParseFloatCheckPointNotOk
inc %rax 
mov (%rax), %dl 
cmp $'.', %dl 
jz __userParseFloatCheckPointOk
jmp __userParseFloatCheckPoint

__userParseFloatCheckPointNotOk:
movb $'.', (%rax)
inc %rax 
movb $'0', (%rax)
inc %rax 
movb $0, (%rax)

__userParseFloatCheckPointOk:
mov $buf, %rax 

call __clearBuf2
call __clearBuf3
mov $buf2, %rbx # здесь будет содержаться целая часть 
mov $buf3, %rcx # здесь будет содержаться дробная часть
mov (%rax), %dl 
cmp $'-', %dl 
jnz __userIsPos
mov $1, %r12 # признак отрицательного числа 
inc %rax 
jmp __userParseFloatLocal 
__userIsPos:
mov $0, %r12 
__userParseFloatLocal: 
mov (%rax), %dl 
cmp $'.', %dl
jz __userPoint
mov %dl, (%rbx)
inc %rax 
inc %rbx 
jmp __userParseFloatLocal
__userPoint:
movb $0, (%rbx)
mov %rax, %rbx 
mov $buf2, %rsi 
call __len 
cmp $8, %rax # целое число - не более 7 цифр 
jl __userParseFloatZ
mov $buf2, %rbx 
movb $48, (%rbx)
inc %rbx 
movb $0, (%rbx) 
mov $buf3, %rbx 
movb $48, (%rbx)
inc %rbx 
movb $0, (%rbx)
jmp __userParseNow 
__userParseFloatZ: 
mov %rbx, %rax 

__userPointLocal:             
inc %rax
mov (%rax), %dl 
cmp $0, %dl 
jz __userParseNow
mov (%rax), %dl 
mov %dl, (%rcx) 
inc %rcx 
jmp __userPointLocal   
__userParseNow:
movb $0, (%rcx)

call __clearBuf
mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenBuf2, %rax 
mov $buf2, %rdi 
call __set 
call __userToNumber
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
jl __userParseFloatCut
mov $buf, %rsi 
add $6, %rsi 
movb $0, (%rsi)
mov $6, %rbx 
__userParseFloatCut: 
call __userToNumber

mov %rax, (buf)
cvtsi2ss (buf), %xmm0  
movss %xmm0, (buf)

__userFloatLocal:
fld (buf)
cmp $0, %rbx 
jz __userFloatOk
dec %rbx 
movss (ten), %xmm0
movss %xmm0, (buf)
fdiv (buf)
fstp (buf)
jmp __userFloatLocal
__userFloatOk:
mov %r10, (buf)
cvtsi2ss (buf), %xmm0    
movss %xmm0, (buf) # целая часть числа 
fadd (buf)
fstp (buf)
movss (buf), %xmm0  
cmp $1, %r12 
jnz __userPos
mov (zero), %rax   
mov %rax, (buf)
fld (buf)
movss %xmm0, (buf)
fsub (buf)
fstp (buf)
movss (buf), %xmm0 
__userPos:
ret
__userParseFloatException:
mov $parseNumberError, %rsi 
call __throughUserError
ret 


 __goto:
 # адрес имени метки, по которой нужно прыгнуть, в %rdi 
 mov %rdi, %rsi 
 call __len 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 call __set 

 mov (memoryBegin), %rax
  __gotoSearch: 
 mov %rax, (labelsPointer)
 cmp %rax, (labelsMax)
 jl __gotoEnd
 mov %rax, %r12 
 call __read  
 
 call __compare 
 cmp $1, %rax 
 jz __goNow  
 mov (labelsPointer), %rax 
 add (labelSize), %rax 
  
 jmp __gotoSearch
 __goNow:
 mov (labelsPointer), %rax 
 add (valSize), %rax 
 mov %rax, %r8 #  cохраняем %rax 
 mov %rax, %rsi 
 call __len 
 
  
 mov $lenBuf, %rsi 
 mov $buf, %rdx  
 mov %r8, %rdi 
 call __set 
 
 call __toNumber
 jmp *%rax 

 __gotoEnd: 
 mov $noSuchMarkError, %rsi 
 call __len 
 mov %rax, %r8 
 mov $noSuchMarkError, %r9 
 mov $buf2, %r11 
 call __concatinate

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set 
 
 mov $lenUserData, %r9 
 mov $buf, %r9 
 mov $enter, %r11 
 call __concatinate

 mov $userData, %rsi 
 call __throughUserError
 ret 

 __less:
 # вход: buf и buf2 
 # %rax - тип операции 
 # 0 - целочисленный
 # 1 - вещественный 
 # выход: userData 
 cmp $0, %rax 
 jnz __lessFloat
 call __toNumber
 mov %rax, %r12 # сохраняем %rax 
 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set 
 call __toNumber 
 
 cmp %r12, %rax 
 jg __isLess 
 call __clearUserData
 movb $'0', (userData)
 ret 
 __isLess:
 call __clearUserData
 movb $'1', (userData)
 ret 
 __lessFloat:
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __parseFloat

 cmpss $2, %xmm1, %xmm0
 pextrb $3, %xmm0, %rax
 cmp $0, %rax 
 jz __isLess 
 call __clearUserData
 movb $'0', (userData)
 ret 

  __lessOrEqual:
 # вход: buf и buf2 
 # %rax - тип операции 
 # 0 - целочисленный
 # 1 - вещественный 
 # выход: userData 
 cmp $0, %rax 
 jnz __lessOrEqualFloat
 call __toNumber
 mov %rax, %r12 # сохраняем %rax 
 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set 
 call __toNumber 
 
 cmp %r12, %rax 
 jge __isLessOrEqual 
 call __clearUserData
 movb $'0', (userData)
 ret 
 __isLessOrEqual:
 call __clearUserData
 movb $'1', (userData)
 ret 
 __lessOrEqualFloat:
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __parseFloat

 cmpss $1, %xmm1, %xmm0
 pextrb $3, %xmm0, %rax
 cmp $0, %rax 
 jz __isLessOrEqual 
 call __clearUserData
 movb $'0', (userData)
 ret 

 __more:
 call __lessOrEqual
 xor %rax, %rax 
 mov (userData), %al
 cmp $'0', %al 
 jz __isMore
 movb $'0', (userData)
 ret 
 __isMore:
 movb $'1', (userData)
 ret 

 __moreOrEqual:
 call __less
 xor %rax, %rax 
 mov (userData), %al
 cmp $'0', %al 
 jz __isMoreOrEqual
 movb $'0', (userData)
 ret 
 __isMoreOrEqual:
 movb $'1', (userData)
 ret

 __eq:
 # вход: buf и buf2 
 # %rax - тип операции 
 # 0 - целочисленный
 # 1 - вещественный 
 # выход: userData 
 cmp $0, %rax 
 jnz __equalFloat
 call __toNumber
 mov %rax, %r12 # сохраняем %rax 
 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set 
 call __toNumber 
 
 cmp %r12, %rax 
 jz __isEqual  
 call __clearUserData
 movb $'0', (userData)
 ret 
 __isEqual:
 call __clearUserData
 movb $'1', (userData)
 ret 
 __equalFloat:
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __parseFloat

 cmpss $4, %xmm1, %xmm0
 pextrb $3, %xmm0, %rax
 cmp $0, %rax 
 jz __isEqual  
 call __clearUserData
 movb $'0', (userData)
 ret 

  __eqString:
 # (buf) - адрес первой строки 
 # (buf2) - адрес второй строки 
 # выход: (userData) = 1 - равны, (userData) = 0 - не равны 
 mov (buf), %rsi 
 call __len 
 mov %rax, %rbx 
 mov (buf2), %rsi 
 call __len 
 cmp %rax, %rbx 
 jnz __eqStringNotEqual

 mov (buf), %rax 
 mov (buf2), %rbx 
 __eqStringCompareLocal:
 movb (%rax), %dl 
 cmp $0, %dl 
 jz __eqStringEqual
 movb (%rax), %dl
 cmp %dl, (%rbx) 
 jnz __eqStringNotEqual
 inc %rax 
 inc %rbx 
 jmp __eqStringCompareLocal

 __eqStringNotEqual: 
 movb $'0', (userData)
 ret 
 __eqStringEqual:
 movb $'1', (userData)
 ret


 __parseBool:
 # buf - источник (строка)
 # %rax - результат

 xor %rax, %rax 
 mov (buf), %al  
 cmp $'1', %al 
 jnz __parseFalse
 mov $1, %rax 
 ret  
 __parseFalse:
 mov $0, %rax 
 ret 

 __userParseBool:
 # buf - источник (строка)
 # %rax - результат
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBigTrueVal, %rax 
 mov $bigTrueVal, %rdi 
 call __set
 call __compare
 cmp $1, %rax 
 jz __userParseBoolTrue 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBigFalseVal, %rax 
 mov $bigFalseVal, %rdi 
 call __set
 call __compare
 cmp $1, %rax 
 jz __userParseBoolFalse 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenTrueVal, %rax 
 mov $trueVal, %rdi 
 call __set
 call __compare
 cmp $1, %rax 
 jz __userParseBoolTrue

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenFalseVal, %rax 
 mov $falseVal, %rdi 
 call __set
 call __compare
 cmp $1, %rax 
 jz __userParseBoolFalse 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenOneVal, %rax 
 mov $oneVal, %rdi 
 call __set
 call __compare
 cmp $1, %rax 
 jz __userParseBoolTrue 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenZeroVal, %rax 
 mov $zeroVal, %rdi 
 call __set
 call __compare
 cmp $1, %rax 
 jz __userParseBoolFalse  

 jmp __userParseBoolException 
 __userParseBoolTrue:
 mov $1, %rax 
 ret
 __userParseBoolFalse:
 mov $0, %rax 
 ret 
 __userParseBoolException:
 mov $parseBoolError, %rsi 
 call __throughUserError
 ret 

__boolToStr:
 # вход: buf
 # выход: userData
 call __clearUserData
 mov (buf), %al 
 cmp $1, %al 
 jnz __boolToStrEndTrue
 mov $userData, %rax 
 movb $'1', (%rax)
 inc %rax 
 movb $0, (%rax)
 ret 
 __boolToStrEndTrue:
 mov $userData, %rax 
 movb $'0', (%rax)
 inc %rax 
 movb $0, (%rax)
 ret

 __and:
 # вход: buf и buf2 в виде строк 
 # выход: userData в виде строки 
 call __clearUserData
 call __parseBool 
 mov %rax, (userData)
 mov (buf2), %rax 
 mov %rax, (buf)
 call __parseBool
 mov %rax, (buf2)
 mov (userData), %rax 
 mov %rax, (buf)


 mov (buf), %rax 
 and (buf2), %rax  

 cmp $1, %rax 
 jz __andTrue 
 movb $'0', (userData)
 ret 
 __andTrue:
 movb $'1', (userData)
 
 ret 

 __or:
 # вход: buf и buf2 в виде строк 
 # выход: userData в виде строки 
 call __clearUserData
 call __parseBool 
 mov %rax, (userData)
 mov (buf2), %rax 
 mov %rax, (buf)
 call __parseBool
 mov %rax, (buf2)
 mov (userData), %rax 
 mov %rax, (buf)


 mov (buf), %rax 
 or (buf2), %rax  

 cmp $1, %rax 
 jz __orTrue 
 movb $'0', (userData)
 ret 
 __orTrue:
 movb $'1', (userData)
 
 ret 

 __xor:
 # вход: buf и buf2 в виде строк 
 # выход: userData в виде строки 
 call __clearUserData
 call __parseBool 
 mov %rax, (userData)
 mov (buf2), %rax 
 mov %rax, (buf)
 call __parseBool
 mov %rax, (buf2)
 mov (userData), %rax 
 mov %rax, (buf)


 mov (buf), %rax 
 xor (buf2), %rax  

 cmp $1, %rax 
 jz __xorTrue 
 movb $'0', (userData)
 ret 
 __xorTrue:
 movb $'1', (userData)
 
 ret 

__not:
 # вход: buf в виде строки 
 # выход: userData в виде строки 
 call __clearUserData

 mov (buf), %al 
 cmp $'1', %al 
 jz __notTrue
 movb $'1', (userData)
 ret 
 __notTrue:
 movb $'0', (userData)
 ret 

 __exists:
# %rdi - адрес строки с именем файла 
# открыть файл
  mov $0,  %rsi   # открываем для чтения
  mov $2,  %rax   # номер системного вызова
  syscall         # вызов функции "открыть файл"
  cmp $0, %rax    # нет ли ошибки при открытии  
  jl __existsNot
  # закрыть файл
  mov  %rax, %rdi  # дескриптор файла
  mov  $3, %rax   # номер системного вызова
  syscall 
  mov $1, %rax 
  ret 
  __existsNot:
  xor %rax, %rax 
  ret
 
 __index:
# вход:
# %rsi - адрес строки
# %rdi - адрес подстроки 
# выход: %rax
mov %rsi, %r8 
call __len 
mov %r8, %rsi 
mov %rax, %r8 
mov %rsi, %r9 
mov %rdi, %rsi 
call __len 
mov %r9, %rsi

mov %rax, %r9 
mov $-1,  %rax
cmp %r8, %r9 
jg __indexEnd

mov %rdi, %r8 # save %rdi  

xor %rax, %rax 
__indexCompare:
mov $1, %dl # is substring in this index?
mov %rsi, %r9 # save %rsi 
__indexCompareLoop:
mov (%rdi), %bl
cmp $0, %bl
jz __indexCompareEnd  
mov (%rsi), %bl
cmp $0, %bl
jz __indexNo   
mov (%rdi), %bl 
mov (%rsi), %cl 
cmp %bl, %cl
jnz __indexChange
inc %rdi
inc %rsi  
jmp __indexCompareLoop   
__indexChange:
mov $0, %dl 
__indexCompareEnd:
cmp $1, %dl 
jz __indexEnd 
mov %r9, %rsi # restore %rsi 
inc %rsi 
mov (%rsi), %bl 
cmp $0, %bl 
jnz __indexNoEnd
__indexNo: 
mov $-1, %rax  
jmp __indexEnd 
__indexNoEnd: 
inc %rax 
mov %r8, %rdi # restore %rdi 
jmp __indexCompare 

__indexEnd:

ret 

__singleSlice:
# input: %rax - address of the string 
# %rbx - number
# output: ^systemVar 
push %rax 
push %rbx

mov $lenVarName, %rsi 
mov $varName, %rdx 
mov $lenSystemVarName, %rax 
mov $systemVarName, %rdi 
call __set 
call __getVar

mov (userData), %rsi 
call __clear

pop %rbx
pop %rax

push %rax 
push %rbx 

mov %rax, %rsi 
call __len  
mov %rax, %rsi 

pop %rbx 
pop %rax 

# check if we are out of bounds 
cmp $0, %rbx 
jl __singleSliceException
cmp %rbx, %rsi 
jle __sliceException

mov (userData), %rsi 
add %rbx, %rax 
# make slice 
mov (%rax), %bl 
mov %bl, (%rsi)
inc %rsi 
movb $0, (%rsi)
ret 
__singleSliceException:
mov $sliceBoundError, %rsi 
call __throughUserError 
ret 


__slice:
# input: %rax - address of the string 
# %rbx - left number 
# %rcx - right number 
# output: ^systemVar
push %rax
push %rbx
push %rcx

mov $lenVarName, %rsi 
mov $varName, %rdx 
mov $lenSystemVarName, %rax 
mov $systemVarName, %rdi 
call __set 
call __getVar

mov (userData), %rsi 
call __clear 
 
pop %rcx
pop %rbx 
pop %rax 

push %rax 
push %rbx 
push %rcx 

mov %rax, %rsi 
call __len  
mov %rax, %rsi 

pop %rcx 
pop %rbx 
pop %rax 

# check if we are out of bounds 
cmp $0, %rbx 
jl __sliceException
cmp $0, %rcx 
jl __sliceException
cmp %rbx, %rsi 
jle __sliceException
cmp %rcx, %rsi 
jl __sliceException 
cmp %rbx, %rcx 
jl __sliceException

# make slice 
mov (userData), %rdi 
mov %rax, %rsi 
add %rbx, %rax # left bound 
add %rcx, %rsi # right bound 

__sliceMake:
cmp %rax, %rsi 
jz __sliceMakeEnd  
mov (%rax), %dl 
mov %dl, (%rdi)
inc %rax 
inc %rdi 
jmp __sliceMake 
__sliceMakeEnd:
movb $0, (%rdi)
ret 
__sliceException:  
mov $sliceBoundError, %rsi 
call __throughUserError 
ret 

__isLetter:
# only for English language!
# input: %rax - address of the string 
# output: 
# %rbx = 0 - is not letter 
# %rbx = 1 - is letter
push %rax

mov %rax, %rsi 
call __len 
cmp $1, %rax 
jnz __isLetterException

pop %rax
mov (%rax), %bl 

cmp $65, %bl 
jl __isLetterNo  
cmp $91, %bl 
jl __isLetterYes
cmp $97, %bl
jl __isLetterNo 
cmp $123, %bl 
jl __isLetterYes 

__isLetterNo:
xor %rbx, %rbx 
ret 
__isLetterYes:
mov $1, %rbx 
ret 
__isLetterException:
mov $isLetterError, %rsi 
call __throughUserError
ret

__isDigit:
# input: %rax - address of the string 
# output: 
# %rbx = 0 - is not digit 
# %rbx = 1 - is digit 
push %rax

mov %rax, %rsi 
call __len 
cmp $1, %rax 
jnz __isDigitException

pop %rax
mov (%rax), %bl 

cmp $48, %bl 
jl __isDigitNo  
cmp $58, %bl 
jl __isDigitYes

__isDigitNo:
xor %rbx, %rbx 
ret 
__isDigitYes:
mov $1, %rbx 
ret 
__isDigitException:
mov $isDigitError, %rsi 
call __throughUserError
ret

__input:
# вход: имя переменной по адресу $varName  
call __getVar
mov (userData), %rsi  
call __clear
mov (userData), %rsi 
movb $0, (%rsi) 
__inputLoop:

mov $lenInputBuf, %rdx 
mov $inputBuf, %rsi 
mov $0, %rdi
mov $0, %rax
syscall        

mov $inputBuf, %rsi 
mov (%rsi), %al 
cmp $'\n', %al
jz __inputEnd 

mov $varName, %r8 
mov $inputBuf, %r9 
mov $1, %rax 
mov $0, %rbx 
call __userConcatinate


jmp __inputLoop 

__inputEnd:
ret  

__openFile:
# input:
# %rax - address of the string with file name
# %rbx - file opening mode
# 0 - read
# 1 - write 
# 2 - append
# output:
# %rax - file descriptor number 
# открыть файл
cmp $0, %rbx 
jz __openFileRead 
cmp $1, %rbx 
jz __openFileWrite
cmp $2, %rbx 
jz __openFileAppend 
jmp __openFileException

__openFileRead:
  mov %rax, %rdi   # адрес строки с именем файла
  mov $0,  %rsi   # открываем для чтения
  mov $2,  %rax   # номер системного вызова
  syscall         # вызов функции "открыть файл"
  cmp $0, %rax    # нет ли ошибки при открытии
  jl  __openFileException          # перейти к концу программы
  ret  
__openFileWrite:
  mov %rax, %rdi    # адрес строки с именем файла
  mov $1,  %rsi    # открываем для записи (O_WRONLY)
  or  $64, %rsi   # создать, если файла нет (O_CREAT)
  mov $0777, %rdx  # права доступа создаваемого файла
  mov $2,  %rax    # номер системного вызова
  syscall          # вызов функции "открыть файл"
  cmp $0, %rax    # нет ли ошибки при открытии
  jl  __openFileException
  ret
__openFileAppend:
  mov %rax, %rdi    # адрес строки с именем файла
  mov $0x441,  %rsi    # O_CREAT| O_WRONLY | O_APPEND
  mov $0777, %rdx  # права доступа создаваемого файла
  mov $2,  %rax    # номер системного вызова
  syscall          # вызов функции "открыть файл"
  cmp $0, %rax    # нет ли ошибки при открытии
  jl  __openFileException
  ret 
__openFileException:
  mov $openFileError, %rsi 
  call __throughUserError
  ret 

__closeFile:
# закрыть файл
# %rdi - descriptor file number 
mov  $3, %rax   # номер системного вызова
syscall
ret 

__delFile:
 # input:
 # %rdi - адрес строки с путем до файла 
 mov $87,  %rax   
 syscall      
 ret 

__writeToFile:
# input:
# %r10 - file descriptor number 
# %r8 - variable name address 
push %r10 

mov $lenVarName, %rsi 
mov $varName, %rdx 
mov $lenVarName, %rax 
mov %r8, %rdi
call __set 
call __getVar
mov (userData), %rsi 
call __len 

pop %r10 
mov (userData), %rsi 
mov %r10, %rdi 
mov %rax, %rdx 
mov $1, %rax 
syscall 

ret

__readFromFile:
# input:
# %r10 - file descriptor number 
# %rbx - size to read 
# %r8 - variable name address
# output:
# %rax - number of bytes was read
movq $0, (numberOfReadBytes) 
xor %r9, %r9 
cmp $0, %rbx 
jle __readFromFileEnd

push %r10 
push %rbx 
push %r8 

mov $lenVarName, %rsi 
mov $varName, %rdx 
mov $lenVarName, %rax 
mov %r8, %rdi
call __set 
call __getVar
mov (userData), %rsi 
call __clear 
mov (userData), %rsi 
movb $0, (%rsi)
mov $readBuf, %rsi 
call __clear

pop %r8 
pop %rbx 
pop %r10 

# читать из файла
xor %r9, %r9 # sum of the read bytes 
mov $readBuf,  %rsi  # адрес буфера
mov $lenReadBuf, %rdx  # длина буфера
dec %rdx 

mov $readBufSum, %r12 
__readFromFileLo:
mov %r10,  %rdi  # номер дескриптора
mov $0,   %rax  # номер системной функции чтения
syscall         # системный вызов
cmp $0, %rax 
jl __readFromFileException

add %rax, %r9

mov (numberOfReadBytes), %rcx 
add %rax, %rcx
mov %rcx, (numberOfReadBytes) 

inc %rax 
mov %rax, %rdi 

push %rbx 

mov $readBuf, %rbx 
add %rax, %rbx
dec %rbx  
movb $0, (%rbx)

mov  %rax, %rdi

 

mov (readBuf), %al 
mov %al, (%r12)
inc %r12  

pop %rbx 


mov (numberOfReadBytes),%rcx 
cmp %rcx, %rbx 
jz __readFromFileEnd

inc %r9 # 0 byte at the end  


cmp $lenReadBufSum, %r9 
jge __readFromFileUnion 

dec %r9 
cmp  $lenReadBuf, %rdi
jz   __readFromFileLo         

__readFromFileEnd:
mov %r9, %rax 
cmp $0, %rax 

jl __readFromFileException 

push %r8 
movb $0, (%r12)

mov $lenVarName, %rsi 
mov $varName, %rdx 
mov $lenVarName, %rax 
mov %r8, %rdi
call __set 

pop %r8 
mov $readBufSum, %r9
mov $1, %rax 
mov $0, %rbx  
call __userConcatinate

mov $readBufSum, %rsi 
call __clear 
mov $readBufSum, %r12 

mov (numberOfReadBytes), %rax 
ret 

__readFromFileUnion:

push %rsi 
push %rdi 
push %rdx 
push %r8
push %rbx 
push %r10

push %r8 
movb $0, (%r12)

mov $lenVarName, %rsi 
mov $varName, %rdx 
mov $lenVarName, %rax 
mov %r8, %rdi
call __set 

pop %r8 
mov $readBufSum, %r9
mov $1, %rax 
mov $0, %rbx  
call __userConcatinate

mov $readBufSum, %rsi 
call __clear 
mov $readBufSum, %r12 

xor %r9, %r9 
pop %r10 
pop %rbx 
pop %r8
pop %rdx 
pop %rdi 
pop %rsi 

jmp __readFromFileLo

__readFromFileException:
mov $readFromFileError, %rsi 
call __throughUserError
ret 
 

.globl _start
_start:
 call __initLabels
 call __firstMem
 call __firstStrMem

 # toPanic variable
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNamePanic, %rax 
 mov $varNamePanic, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi
 call __set 
 call __defineVar

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNamePanic, %rax 
 mov $varNamePanic, %rdi
 call __set 
 mov $oneVal, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar 

 # errorVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar

 # ^systemVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName, %rax 
 mov $systemVarName, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar 
