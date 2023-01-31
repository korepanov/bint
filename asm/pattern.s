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
toStr:
mov $10, %r8
mov $$buf, %rsi
xor %rdi, %rdi
lo:
xor %rdx, %rdx
div %r8
add $48, %dl
mov %dl, (%rsi)
inc %rsi
inc %rdi
cmp $0, %rax
jnz lo
mov $$buf2, %rbx
mov $$buf, %rcx
add %rdi, %rcx
dec %rcx
mov %rdi, %rdx
exc:
mov (%rcx), %al
mov %al, (%rbx)
dec %rcx
inc %rbx
dec %rdx
jnz exc
ret
.globl _start
_start:

