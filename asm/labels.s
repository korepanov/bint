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
mov $labelName5, %rbx
 __initLabelsName5: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx5
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName5
 __initLabelsNameEx5:
 movb $0, (%rdi)

 mov (label5), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr5:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx5
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr5
 __initLabelsAddrEx5:
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
mov $labelName6, %rbx
 __initLabelsName6: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx6
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName6
 __initLabelsNameEx6:
 movb $0, (%rdi)

 mov (label6), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr6:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx6
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr6
 __initLabelsAddrEx6:
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
mov $labelName7, %rbx
 __initLabelsName7: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx7
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName7
 __initLabelsNameEx7:
 movb $0, (%rdi)

 mov (label7), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr7:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx7
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr7
 __initLabelsAddrEx7:
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
mov $labelName8, %rbx
 __initLabelsName8: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx8
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName8
 __initLabelsNameEx8:
 movb $0, (%rdi)

 mov (label8), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr8:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx8
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr8
 __initLabelsAddrEx8:
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
mov $labelName9, %rbx
 __initLabelsName9: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx9
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName9
 __initLabelsNameEx9:
 movb $0, (%rdi)

 mov (label9), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr9:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx9
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr9
 __initLabelsAddrEx9:
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
mov $labelName10, %rbx
 __initLabelsName10: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx10
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName10
 __initLabelsNameEx10:
 movb $0, (%rdi)

 mov (label10), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr10:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx10
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr10
 __initLabelsAddrEx10:
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
mov $labelName11, %rbx
 __initLabelsName11: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx11
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName11
 __initLabelsNameEx11:
 movb $0, (%rdi)

 mov (label11), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr11:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx11
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr11
 __initLabelsAddrEx11:
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
mov $labelName12, %rbx
 __initLabelsName12: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx12
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName12
 __initLabelsNameEx12:
 movb $0, (%rdi)

 mov (label12), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr12:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx12
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr12
 __initLabelsAddrEx12:
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
mov $labelName13, %rbx
 __initLabelsName13: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx13
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName13
 __initLabelsNameEx13:
 movb $0, (%rdi)

 mov (label13), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr13:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx13
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr13
 __initLabelsAddrEx13:
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
mov $labelName14, %rbx
 __initLabelsName14: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx14
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName14
 __initLabelsNameEx14:
 movb $0, (%rdi)

 mov (label14), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr14:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx14
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr14
 __initLabelsAddrEx14:
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
mov $labelName15, %rbx
 __initLabelsName15: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx15
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName15
 __initLabelsNameEx15:
 movb $0, (%rdi)

 mov (label15), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr15:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx15
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr15
 __initLabelsAddrEx15:
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
mov $labelName16, %rbx
 __initLabelsName16: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx16
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName16
 __initLabelsNameEx16:
 movb $0, (%rdi)

 mov (label16), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr16:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx16
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr16
 __initLabelsAddrEx16:
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
mov $labelName17, %rbx
 __initLabelsName17: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx17
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName17
 __initLabelsNameEx17:
 movb $0, (%rdi)

 mov (label17), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr17:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx17
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr17
 __initLabelsAddrEx17:
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
mov $labelName18, %rbx
 __initLabelsName18: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx18
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName18
 __initLabelsNameEx18:
 movb $0, (%rdi)

 mov (label18), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr18:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx18
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr18
 __initLabelsAddrEx18:
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
mov $labelName19, %rbx
 __initLabelsName19: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx19
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName19
 __initLabelsNameEx19:
 movb $0, (%rdi)

 mov (label19), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr19:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx19
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr19
 __initLabelsAddrEx19:
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
mov $labelName20, %rbx
 __initLabelsName20: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx20
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName20
 __initLabelsNameEx20:
 movb $0, (%rdi)

 mov (label20), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr20:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx20
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr20
 __initLabelsAddrEx20:
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
mov $labelName21, %rbx
 __initLabelsName21: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx21
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName21
 __initLabelsNameEx21:
 movb $0, (%rdi)

 mov (label21), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr21:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx21
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr21
 __initLabelsAddrEx21:
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
mov $labelName22, %rbx
 __initLabelsName22: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx22
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName22
 __initLabelsNameEx22:
 movb $0, (%rdi)

 mov (label22), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr22:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx22
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr22
 __initLabelsAddrEx22:
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
mov $labelName23, %rbx
 __initLabelsName23: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx23
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName23
 __initLabelsNameEx23:
 movb $0, (%rdi)

 mov (label23), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr23:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx23
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr23
 __initLabelsAddrEx23:
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
mov $labelName24, %rbx
 __initLabelsName24: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx24
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName24
 __initLabelsNameEx24:
 movb $0, (%rdi)

 mov (label24), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr24:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx24
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr24
 __initLabelsAddrEx24:
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
mov $labelName25, %rbx
 __initLabelsName25: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx25
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName25
 __initLabelsNameEx25:
 movb $0, (%rdi)

 mov (label25), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr25:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx25
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr25
 __initLabelsAddrEx25:
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
mov $labelName26, %rbx
 __initLabelsName26: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx26
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName26
 __initLabelsNameEx26:
 movb $0, (%rdi)

 mov (label26), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr26:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx26
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr26
 __initLabelsAddrEx26:
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
mov $labelName27, %rbx
 __initLabelsName27: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx27
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName27
 __initLabelsNameEx27:
 movb $0, (%rdi)

 mov (label27), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr27:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx27
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr27
 __initLabelsAddrEx27:
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
mov $labelName28, %rbx
 __initLabelsName28: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx28
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName28
 __initLabelsNameEx28:
 movb $0, (%rdi)

 mov (label28), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr28:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx28
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr28
 __initLabelsAddrEx28:
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
mov $labelName29, %rbx
 __initLabelsName29: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx29
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName29
 __initLabelsNameEx29:
 movb $0, (%rdi)

 mov (label29), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr29:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx29
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr29
 __initLabelsAddrEx29:
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
mov $labelName30, %rbx
 __initLabelsName30: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx30
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName30
 __initLabelsNameEx30:
 movb $0, (%rdi)

 mov (label30), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr30:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx30
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr30
 __initLabelsAddrEx30:
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
mov $labelName31, %rbx
 __initLabelsName31: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx31
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName31
 __initLabelsNameEx31:
 movb $0, (%rdi)

 mov (label31), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr31:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx31
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr31
 __initLabelsAddrEx31:
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
mov $labelName32, %rbx
 __initLabelsName32: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx32
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName32
 __initLabelsNameEx32:
 movb $0, (%rdi)

 mov (label32), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr32:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx32
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr32
 __initLabelsAddrEx32:
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
mov $labelName33, %rbx
 __initLabelsName33: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx33
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName33
 __initLabelsNameEx33:
 movb $0, (%rdi)

 mov (label33), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr33:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx33
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr33
 __initLabelsAddrEx33:
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
mov $labelName34, %rbx
 __initLabelsName34: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx34
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName34
 __initLabelsNameEx34:
 movb $0, (%rdi)

 mov (label34), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr34:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx34
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr34
 __initLabelsAddrEx34:
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
mov $labelName35, %rbx
 __initLabelsName35: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx35
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName35
 __initLabelsNameEx35:
 movb $0, (%rdi)

 mov (label35), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr35:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx35
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr35
 __initLabelsAddrEx35:
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
mov $labelName36, %rbx
 __initLabelsName36: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx36
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName36
 __initLabelsNameEx36:
 movb $0, (%rdi)

 mov (label36), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr36:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx36
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr36
 __initLabelsAddrEx36:
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
mov $labelName37, %rbx
 __initLabelsName37: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx37
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName37
 __initLabelsNameEx37:
 movb $0, (%rdi)

 mov (label37), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr37:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx37
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr37
 __initLabelsAddrEx37:
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
mov $labelName38, %rbx
 __initLabelsName38: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx38
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName38
 __initLabelsNameEx38:
 movb $0, (%rdi)

 mov (label38), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr38:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx38
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr38
 __initLabelsAddrEx38:
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
mov $labelName39, %rbx
 __initLabelsName39: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx39
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName39
 __initLabelsNameEx39:
 movb $0, (%rdi)

 mov (label39), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr39:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx39
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr39
 __initLabelsAddrEx39:
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
mov $labelName40, %rbx
 __initLabelsName40: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx40
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName40
 __initLabelsNameEx40:
 movb $0, (%rdi)

 mov (label40), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr40:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx40
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr40
 __initLabelsAddrEx40:
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
mov $labelName41, %rbx
 __initLabelsName41: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx41
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName41
 __initLabelsNameEx41:
 movb $0, (%rdi)

 mov (label41), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr41:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx41
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr41
 __initLabelsAddrEx41:
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
mov $labelName42, %rbx
 __initLabelsName42: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx42
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName42
 __initLabelsNameEx42:
 movb $0, (%rdi)

 mov (label42), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr42:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx42
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr42
 __initLabelsAddrEx42:
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
mov $labelName43, %rbx
 __initLabelsName43: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx43
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName43
 __initLabelsNameEx43:
 movb $0, (%rdi)

 mov (label43), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr43:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx43
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr43
 __initLabelsAddrEx43:
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
