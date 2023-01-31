.data
$enter:
.ascii "\n"
$heapSize:
.quad 0
$heapPointer:
.quad 0
$buf:
.space 256, 0
$buf2:
.space 256, 0
s:
.space 200, 0
msg1:
.ascii "Hello world!\n"
len = . - msg1

.text

print:
 mov (%rsi), %al	
 cmp $0, %al	
 jz  ex			
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall		    
 inc %rsi		  
 dec %r8		    
 jnz print
ex:
 ret

set: #set strings
 # %edi = %esi
 mov (%esi), %al
 mov %al, (%edi)
 inc %esi
 inc %edi
 cmp $0, (%esi)
 jnz set
 ret 

toStr:
 # число в %rax 
 # подготовка преобразования числа в строку
  movq $0, ($buf2)
  mov $10, %r8    # делитель
  mov $$buf, %rsi  # адрес начала буфера 
  xor %rdi, %rdi  # обнуляем счетчик
# преобразуем путем последовательного деления на 10
lo:
  xor %rdx, %rdx  # число в rdx:rax
  div %r8         # делим rdx:rax на r8
  add $48, %dl    # цифру в символ цифры
  mov %dl, (%rsi) # в буфер
  inc %rsi        # на следующую позицию в буфере
  inc %rdi        # счетчик увеличиваем на 1
  cmp $0, %rax    # проверка конца алгоритма
  jnz lo          # продолжим цикл?
# число записано в обратном порядке,
# вернем правильный, перенеся в другой буфер 
  mov $$buf2, %rbx # начало нового буфера
  mov $$buf, %rcx  # старый буфер
  add %rdi, %rcx  # в конец
  dec %rcx        # старого буфера
  mov %rdi, %rdx  # длина буфера
# перенос из одного буфера в другой
exc:
  mov (%rcx), %al # из старого буфера
  mov %al, (%rbx) # в новый
  dec %rcx        # в обратном порядке  
  inc %rbx        # продвигаемся в новом буфере
  dec %rdx        # а в старом в обратном порядке
  jnz exc         # проверка конца алгоритма
  ret 

newMem:
# адрес начала выделяемой памяти в %r8
# получить адрес начала области для выделения памяти
 mov $12, %rax
 xor %rdi, %rdi
 syscall
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8 
# выделить динамическую память
 mov $4096, %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz _stop
 movq ($heapSize), %rax
 mov $4096, %rbx
 call sum 
 mov %rax, ($heapSize) 
 ret
 
sum:
# операнды в rax и в rbx
# результат в rax

 cmp $0, %rbx
 je end_sum
 dec %rbx
 inc %rax
 jmp sum
 end_sum:
 ret

defineVar:
 movq ($heapSize), %rax
 movq ($heapPointer), %rbx
 call sum
 cmp ($heapPointer), %rax 
 jl defOk
 call newMem
 movq %r8, ($heapPointer)
 defOk:
 movq ($heapPointer), %rax
 movq $72, %rbx
 call sum
 movq %rax, ($heapPointer)
 ret 

.globl _start
_start:

mov $1, %r9

#loop:
call defineVar
#dec %r9

movq ($heapPointer), %rax
call toStr
mov $$buf2, %rsi
call print
mov $$enter, %rsi
call print

#cmp $0, %r9
#jne loop

movq ($heapSize), %rax
call toStr
mov $$buf2, %rsi
call print
 
_stop:
mov $60, %rax
xor %rdi, %rdi
syscall

