.text
__initLabels:
 # формирует в %rax адрес начала выделяемой памяти для firstMem
 mov $12, %rax
 xor %rdi, %rdi
 syscall
 mov %rax, (memoryBegin)
 call __newLabelMem
 mov (memoryBegin), %rax 
 add (labelSize), %rax 
 mov %rax, %r12 # сохраняем rax 
 mov (memoryBegin), %rdi
 mov (memoryBegin), %r9
 
 
 mov %rdi, %r10 # сохраняем %rdi  
 mov %r9, %rsi # сохраняем %rsi 
 mov %r12, %rax 
 
 call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi # восстанавливаем  
 mov %rsi, %r9 
 

mov $labelName0, %rbx
 __initLabelsName0: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx0
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName0
 __initLabelsNameEx0:
 movb $0, (%rdi)

 mov (label0), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr0:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx0
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr0
 __initLabelsAddrEx0:
 movb $0, (%rdi)
 add (valSize), %r9 
 mov %r9, %rdi 
 mov %rdi, %r10 
 mov %r9, %rsi 
 mov %r12, %rax
call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi 
 mov %rsi, %r9
mov $labelName1, %rbx
 __initLabelsName1: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx1
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName1
 __initLabelsNameEx1:
 movb $0, (%rdi)

 mov (label1), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr1:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx1
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr1
 __initLabelsAddrEx1:
 movb $0, (%rdi)
 add (valSize), %r9 
 mov %r9, %rdi 
 mov %rdi, %r10 
 mov %r9, %rsi 
 mov %r12, %rax
call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi 
 mov %rsi, %r9
mov $labelName2, %rbx
 __initLabelsName2: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx2
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName2
 __initLabelsNameEx2:
 movb $0, (%rdi)

 mov (label2), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr2:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx2
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr2
 __initLabelsAddrEx2:
 movb $0, (%rdi)
 add (valSize), %r9 
 mov %r9, %rdi 
 mov %rdi, %r10 
 mov %r9, %rsi 
 mov %r12, %rax
call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi 
 mov %rsi, %r9
mov $labelName3, %rbx
 __initLabelsName3: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx3
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName3
 __initLabelsNameEx3:
 movb $0, (%rdi)

 mov (label3), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr3:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx3
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr3
 __initLabelsAddrEx3:
 movb $0, (%rdi)
 add (valSize), %r9 
 mov %r9, %rdi 
 mov %rdi, %r10 
 mov %r9, %rsi 
 mov %r12, %rax
call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi 
 mov %rsi, %r9
mov $labelName4, %rbx
 __initLabelsName4: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx4
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName4
 __initLabelsNameEx4:
 movb $0, (%rdi)

 mov (label4), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr4:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx4
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr4
 __initLabelsAddrEx4:
 movb $0, (%rdi)
 add (valSize), %r9 
 mov %r9, %rdi 
 mov %rdi, %r10 
 mov %r9, %rsi 
 mov %r12, %rax
call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi 
 mov %rsi, %r9
 mov %r12, %rax 
 mov %r12, (labelsMax)
 ret 
