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
.globl _start
_start:
mov $data0, %esi
mov $s, %edi
mark0:
mov (%esi), %al
mov %al, (%edi)
inc %esi
inc %edi
cmp $0, (%esi)
jnz mark0
mov $s, %esi
mov $d, %edi
mark1:
mov (%esi), %al
mov %al, (%edi)
inc %esi
inc %edi
cmp $0, (%esi)
jnz mark1
mov $d, %rsi
call print
mov $msg0, %rsi
call print
mov $msg1, %rsi
call print
mov $60,   %rax
mov $1, %rdi
syscall
mov $msg2, %rsi
call print
mov $60,  %rax
xor %rdi, %rdi
syscall
