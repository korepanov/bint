.data
pageSize:
.quad 4096
varSize:
.quad 200 
enter:
.ascii "\n"
heapBegin:
.quad 0 
heapSize:
.quad 0
heapMax:
.quad 0
heapPointer:
.quad 0
buf:
.space 256, 0
buf2:
.space 256, 0
s:
.space 200, 0
msg1:
.ascii "Hello world!\n"
len = . - msg1

.text

__print:
 mov (%rsi), %al	
 cmp $0, %al	
 jz  __ex			
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall		    
 inc %rsi		  
 dec %r8		    
 jnz __print
__ex:
 ret

__set: #set strings
 # %edi = %esi
 mov (%esi), %al
 mov %al, (%edi)
 inc %esi
 inc %edi
 cmp $0, (%esi)
 jnz __set
 ret 

__toStr:
 # число в %rax 
 # подготовка преобразования числа в строку
  movq $0, (buf2)
  mov $10, %r8    # делитель
  mov $buf, %rsi  # адрес начала буфера 
  xor %rdi, %rdi  # обнуляем счетчик
# преобразуем путем последовательного деления на 10
__lo:
  xor %rdx, %rdx  # число в rdx:rax
  div %r8         # делим rdx:rax на r8
  add $48, %dl    # цифру в символ цифры
  mov %dl, (%rsi) # в буфер
  inc %rsi        # на следующую позицию в буфере
  inc %rdi        # счетчик увеличиваем на 1
  cmp $0, %rax    # проверка конца алгоритма
  jnz __lo          # продолжим цикл?
# число записано в обратном порядке,
# вернем правильный, перенеся в другой буфер 
  mov $buf2, %rbx # начало нового буфера
  mov $buf, %rcx  # старый буфер
  add %rdi, %rcx  # в конец
  dec %rcx        # старого буфера
  mov %rdi, %rdx  # длина буфера
# перенос из одного буфера в другой
__exc:
  mov (%rcx), %al # из старого буфера
  mov %al, (%rbx) # в новый
  dec %rcx        # в обратном порядке  
  inc %rbx        # продвигаемся в новом буфере
  dec %rdx        # а в старом в обратном порядке
  jnz __exc         # проверка конца алгоритма
  ret 

__newMem:
# адрес начала выделяемой памяти в %r8
# получить адрес начала области для выделения памяти
 mov $12, %rax
 xor %rdi, %rdi
 syscall
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8 
 mov %r8, (heapPointer)
# выделить динамическую память
 mov (pageSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __stop
 movq (heapSize), %rax
 mov (pageSize), %rbx
 call __sum 
 mov %rax, (heapSize) 
 mov (pageSize), %rax
 mov (heapPointer), %rbx
 call __sum 
 mov %rax, (heapMax) 
 ret
 
__sum:
# операнды в rax и в rbx
# результат в rax

 cmp $0, %rbx
 je __end_sum
 dec %rbx
 inc %rax
 jmp __sum
 __end_sum:
 ret

__defineVar:
 movq (heapMax), %rax
 cmp (heapPointer), %rax 
 jg defOk
 call __newMem
 movq %r8, (heapPointer)
 defOk:
 movq (heapPointer), %rax
 movq (varSize), %rbx
 call __sum
 movq %rax, (heapPointer)
 ret 

.globl _start
_start:
call __newMem
mov (heapPointer), %rax 
mov %rax, (heapBegin)

movq (heapBegin), %rax
call __toStr
mov $buf2, %rsi
call __print
mov $enter, %rsi
call __print

mov $43, %r9

loop:
call __defineVar
dec %r9

movq (heapPointer), %rax
call __toStr
mov $buf2, %rsi
call __print
mov $enter, %rsi
call __print

cmp $0, %r9
jne loop

movq (heapSize), %rax
call __toStr
mov $buf2, %rsi
call __print
 
__stop:
mov $60, %rax
xor %rdi, %rdi
syscall

