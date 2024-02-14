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
mov $labelName44, %rbx
 __initLabelsName44: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx44
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName44
 __initLabelsNameEx44:
 movb $0, (%rdi)

 mov (label44), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr44:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx44
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr44
 __initLabelsAddrEx44:
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
mov $labelName45, %rbx
 __initLabelsName45: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx45
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName45
 __initLabelsNameEx45:
 movb $0, (%rdi)

 mov (label45), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr45:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx45
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr45
 __initLabelsAddrEx45:
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
mov $labelName46, %rbx
 __initLabelsName46: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx46
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName46
 __initLabelsNameEx46:
 movb $0, (%rdi)

 mov (label46), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr46:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx46
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr46
 __initLabelsAddrEx46:
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
mov $labelName47, %rbx
 __initLabelsName47: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx47
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName47
 __initLabelsNameEx47:
 movb $0, (%rdi)

 mov (label47), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr47:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx47
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr47
 __initLabelsAddrEx47:
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
mov $labelName48, %rbx
 __initLabelsName48: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx48
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName48
 __initLabelsNameEx48:
 movb $0, (%rdi)

 mov (label48), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr48:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx48
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr48
 __initLabelsAddrEx48:
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
mov $labelName49, %rbx
 __initLabelsName49: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx49
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName49
 __initLabelsNameEx49:
 movb $0, (%rdi)

 mov (label49), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr49:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx49
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr49
 __initLabelsAddrEx49:
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
mov $labelName50, %rbx
 __initLabelsName50: 
 mov (%rbx), %dl 
 cmp $0, %dl 
jz __initLabelsNameEx50
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi
 jmp __initLabelsName50
 __initLabelsNameEx50:
 movb $0, (%rdi)

 mov (label50), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi
 mov $buf2, %rbx 
__initLabelsAddr50:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx50
 mov %dl, (%rdi)
 inc %rbx
 inc %rdi 
 jmp __initLabelsAddr50
 __initLabelsAddrEx50:
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
 
 .data
stackEndSymbol:
.ascii "end"
.space 1, 0 
clearSymbol:
.ascii "$clear"
.space 1, 0
lenClearSymbol = . - clearSymbol 
starSymbol:
.ascii "*"
endSymbol:
.ascii ";"
deltaSize:
.quad 2048
pageSize:
.quad 4096 
pageSizeMax:
.quad 4500000000 
shiftSize:
.quad 8192
stackShiftSize:
.quad 8192 
shiftSizeMax:
.quad 9000000000 
varNameSize:
.quad 32
varSize:
.quad 128
typeSize:
.quad 32 
valSize:
.quad 64 
strValSize:
.quad 1024
strValSizeMax:
.quad 22500000  
labelSize:
.quad 128 
labelsMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
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
inputBuf:
.byte 0
lenInputBuf = . - inputBuf 
readBuf:
.byte 0, 0
lenReadBuf = . - readBuf
readBufSum:
.space 4096
lenReadBufSum = . - readBufSum
numberOfReadBytes:
.quad 0
userMem:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserMem = . - userMem 
userMem2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserMem2 = . - userMem2 
mem:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem = . - mem 
mem2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem2 = . - mem2 
mem3:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem3 = . - mem3 
mem4:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem4 = . - mem4 
mem5:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem5 = . - mem5
mem6:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem6 = . - mem6
mem7:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem7 = . - mem7
mem8:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem8 = . - mem8
mem9:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem9 = . - mem9
mem10:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem10 = . - mem10 
mem11:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem11 = . - mem11 
mem12:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem12 = . - mem12 
mem13:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem13 = . - mem13 
mem14:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem14 = . - mem14 
mem15:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem15 = . - mem15 
mem16:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem16 = . - mem16 
mem17:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem17 = . - mem17 
mem18:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem18 = . - mem18 
mem19:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem19 = . - mem19 
mem20:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem20 = . - mem20  
strBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrBegin = . - strBegin
stackBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStackBegin = . - stackBegin  
oldHeapMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenOldHeapMax = . - oldHeapMax
strPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrPointer = . - strPointer
strMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrMax = . - strMax 
stackMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStackMax = . - stackMax
stackPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStackPointer = . - stackPointer   
isNeg:
.byte 0 
lenIsNeg = . - isNeg
userConcatinateFlag:
.byte 0
lenUserConcatinateFlag = . - userConcatinateFlag
userConcatinateFlag2:
.byte 0
lenUserConcatinateFlag2 = . - userConcatinateFlag2
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
strPageNumber:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrNumber = . - strPageNumber
memorySize:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMemorySize = . - memorySize
memoryBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMemoryBegin = . - memoryBegin
labelsEnd:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLabelsEnd = . - labelsEnd
labelsPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLabelsPointer = . - labelsPointer 
varNamePanic:
.ascii "^toPanic"
.space 1, 0
lenVarNamePanic = . - varNamePanic 
varNameError:
.ascii "error"
.space 1, 0
lenVarNameError = . - varNameError 
systemVarName:
.ascii "^systemVar"
.space 1, 0
lenSystemVarName = . - systemVarName
popVarName:
.ascii "^popVar"
.space 1, 0
lenPopVarName = . - popVarName

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
stackType:
.ascii "stack"
.space 1, 0
lenStackType = . - stackType 
enter:
.ascii "\n"
.space 1, 0
lenEnter = . - enter 
bigTrueVal:
.ascii "True"
.space 1, 0
lenBigTrueVal = . - bigTrueVal
trueVal:
.ascii "true"
.space 1, 0
lenTrueVal = . - trueVal 
bigFalseVal:
.ascii "False"
.space 1, 0
lenBigFalseVal = . - bigFalseVal 
falseVal:
.ascii "false"
.space 1, 0
lenFalseVal = . - falseVal 
oneVal:
.ascii "1"
.space 1, 0
lenOneVal = . - oneVal 
zeroVal:
.ascii "0"
.space 1, 0
lenZeroVal = . - zeroVal 

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

fatalError:
.ascii "fatal error: internal error\n"
.space 1, 0 
divZeroError:
.ascii "runtime error: dividing by zero"
.space 1, 0
divINegError:
.ascii "runtime error: @ is not defined for negative numbers"
.space 1, 0
powNegError:
.ascii "runtime error: ^ is not defined for negative base and fractional exponent"
.space 1, 0
powZeroNegError:
.ascii "runtime error: ^ is not defined for zero base and negative exponent"
.space 1, 0
powZeroZeroError:
.ascii "runtime error: ^ is not defined for zero base and zero exponent"
.space 1, 0
noSuchMarkError:
.ascii "runtime error: no such mark: "
.space 1, 0
concError:
.ascii "could not concatinate not string arguments"
.space 1, 0 
strError:
.ascii "the type of the variable to which you want to assign the result of string concatenation is not a string"
.space 1, 0 
parseNumberError:
.ascii "could not parse number: invalid number format"
.space 1, 0  
parseBoolError:
.ascii "could not parse bool: invalid bool format"
.space 1, 0 
sliceBoundError:
.ascii "slice index is out of bounds"
.space 1, 0 
isLetterError:
.ascii "is_letter argument length error"
.space 1, 0 
isDigitError:
.ascii "is_digit argument length error"
.space 1, 0 
openFileError:
.ascii "could not open file"
.space 1, 0
readFromFileError:
.ascii "could not read file"
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
systemVarName0:
.ascii "^systemVar0"
.space 1, 0
lenSystemVarName0 = . - systemVarName0
systemVarName1:
.ascii "^systemVar1"
.space 1, 0
lenSystemVarName1 = . - systemVarName1
systemVarName2:
.ascii "^systemVar2"
.space 1, 0
lenSystemVarName2 = . - systemVarName2
systemVarName3:
.ascii "^systemVar3"
.space 1, 0
lenSystemVarName3 = . - systemVarName3
systemVarName4:
.ascii "^systemVar4"
.space 1, 0
lenSystemVarName4 = . - systemVarName4
systemVarName5:
.ascii "^systemVar5"
.space 1, 0
lenSystemVarName5 = . - systemVarName5
systemVarName6:
.ascii "^systemVar6"
.space 1, 0
lenSystemVarName6 = . - systemVarName6
systemVarName7:
.ascii "^systemVar7"
.space 1, 0
lenSystemVarName7 = . - systemVarName7
systemVarName8:
.ascii "^systemVar8"
.space 1, 0
lenSystemVarName8 = . - systemVarName8
systemVarName9:
.ascii "^systemVar9"
.space 1, 0
lenSystemVarName9 = . - systemVarName9
systemVarName10:
.ascii "^systemVar10"
.space 1, 0
lenSystemVarName10 = . - systemVarName10
systemVarName11:
.ascii "^systemVar11"
.space 1, 0
lenSystemVarName11 = . - systemVarName11
systemVarName12:
.ascii "^systemVar12"
.space 1, 0
lenSystemVarName12 = . - systemVarName12
systemVarName13:
.ascii "^systemVar13"
.space 1, 0
lenSystemVarName13 = . - systemVarName13
systemVarName14:
.ascii "^systemVar14"
.space 1, 0
lenSystemVarName14 = . - systemVarName14
systemVarName15:
.ascii "^systemVar15"
.space 1, 0
lenSystemVarName15 = . - systemVarName15
systemVarName16:
.ascii "^systemVar16"
.space 1, 0
lenSystemVarName16 = . - systemVarName16
systemVarName17:
.ascii "^systemVar17"
.space 1, 0
lenSystemVarName17 = . - systemVarName17
systemVarName18:
.ascii "^systemVar18"
.space 1, 0
lenSystemVarName18 = . - systemVarName18
systemVarName19:
.ascii "^systemVar19"
.space 1, 0
lenSystemVarName19 = . - systemVarName19
systemVarName20:
.ascii "^systemVar20"
.space 1, 0
lenSystemVarName20 = . - systemVarName20
systemVarName21:
.ascii "^systemVar21"
.space 1, 0
lenSystemVarName21 = . - systemVarName21
systemVarName22:
.ascii "^systemVar22"
.space 1, 0
lenSystemVarName22 = . - systemVarName22
systemVarName23:
.ascii "^systemVar23"
.space 1, 0
lenSystemVarName23 = . - systemVarName23
systemVarName24:
.ascii "^systemVar24"
.space 1, 0
lenSystemVarName24 = . - systemVarName24
systemVarName25:
.ascii "^systemVar25"
.space 1, 0
lenSystemVarName25 = . - systemVarName25
systemVarName26:
.ascii "^systemVar26"
.space 1, 0
lenSystemVarName26 = . - systemVarName26
systemVarName27:
.ascii "^systemVar27"
.space 1, 0
lenSystemVarName27 = . - systemVarName27
systemVarName28:
.ascii "^systemVar28"
.space 1, 0
lenSystemVarName28 = . - systemVarName28
systemVarName29:
.ascii "^systemVar29"
.space 1, 0
lenSystemVarName29 = . - systemVarName29
systemVarName30:
.ascii "^systemVar30"
.space 1, 0
lenSystemVarName30 = . - systemVarName30
systemVarName31:
.ascii "^systemVar31"
.space 1, 0
lenSystemVarName31 = . - systemVarName31
systemVarName32:
.ascii "^systemVar32"
.space 1, 0
lenSystemVarName32 = . - systemVarName32
systemVarName33:
.ascii "^systemVar33"
.space 1, 0
lenSystemVarName33 = . - systemVarName33
systemVarName34:
.ascii "^systemVar34"
.space 1, 0
lenSystemVarName34 = . - systemVarName34
systemVarName35:
.ascii "^systemVar35"
.space 1, 0
lenSystemVarName35 = . - systemVarName35
systemVarName36:
.ascii "^systemVar36"
.space 1, 0
lenSystemVarName36 = . - systemVarName36
systemVarName37:
.ascii "^systemVar37"
.space 1, 0
lenSystemVarName37 = . - systemVarName37
systemVarName38:
.ascii "^systemVar38"
.space 1, 0
lenSystemVarName38 = . - systemVarName38
systemVarName39:
.ascii "^systemVar39"
.space 1, 0
lenSystemVarName39 = . - systemVarName39
systemVarName40:
.ascii "^systemVar40"
.space 1, 0
lenSystemVarName40 = . - systemVarName40
systemVarName41:
.ascii "^systemVar41"
.space 1, 0
lenSystemVarName41 = . - systemVarName41
systemVarName42:
.ascii "^systemVar42"
.space 1, 0
lenSystemVarName42 = . - systemVarName42
systemVarName43:
.ascii "^systemVar43"
.space 1, 0
lenSystemVarName43 = . - systemVarName43
systemVarName44:
.ascii "^systemVar44"
.space 1, 0
lenSystemVarName44 = . - systemVarName44
systemVarName45:
.ascii "^systemVar45"
.space 1, 0
lenSystemVarName45 = . - systemVarName45
systemVarName46:
.ascii "^systemVar46"
.space 1, 0
lenSystemVarName46 = . - systemVarName46
systemVarName47:
.ascii "^systemVar47"
.space 1, 0
lenSystemVarName47 = . - systemVarName47
systemVarName48:
.ascii "^systemVar48"
.space 1, 0
lenSystemVarName48 = . - systemVarName48
systemVarName49:
.ascii "^systemVar49"
.space 1, 0
lenSystemVarName49 = . - systemVarName49
systemVarName50:
.ascii "^systemVar50"
.space 1, 0
lenSystemVarName50 = . - systemVarName50
systemVarName51:
.ascii "^systemVar51"
.space 1, 0
lenSystemVarName51 = . - systemVarName51
systemVarName52:
.ascii "^systemVar52"
.space 1, 0
lenSystemVarName52 = . - systemVarName52
systemVarName53:
.ascii "^systemVar53"
.space 1, 0
lenSystemVarName53 = . - systemVarName53
systemVarName54:
.ascii "^systemVar54"
.space 1, 0
lenSystemVarName54 = . - systemVarName54
systemVarName55:
.ascii "^systemVar55"
.space 1, 0
lenSystemVarName55 = . - systemVarName55
systemVarName56:
.ascii "^systemVar56"
.space 1, 0
lenSystemVarName56 = . - systemVarName56
systemVarName57:
.ascii "^systemVar57"
.space 1, 0
lenSystemVarName57 = . - systemVarName57
systemVarName58:
.ascii "^systemVar58"
.space 1, 0
lenSystemVarName58 = . - systemVarName58
systemVarName59:
.ascii "^systemVar59"
.space 1, 0
lenSystemVarName59 = . - systemVarName59
systemVarName60:
.ascii "^systemVar60"
.space 1, 0
lenSystemVarName60 = . - systemVarName60
systemVarName61:
.ascii "^systemVar61"
.space 1, 0
lenSystemVarName61 = . - systemVarName61
systemVarName62:
.ascii "^systemVar62"
.space 1, 0
lenSystemVarName62 = . - systemVarName62
systemVarName63:
.ascii "^systemVar63"
.space 1, 0
lenSystemVarName63 = . - systemVarName63
systemVarName64:
.ascii "^systemVar64"
.space 1, 0
lenSystemVarName64 = . - systemVarName64
systemVarName65:
.ascii "^systemVar65"
.space 1, 0
lenSystemVarName65 = . - systemVarName65
systemVarName66:
.ascii "^systemVar66"
.space 1, 0
lenSystemVarName66 = . - systemVarName66
systemVarName67:
.ascii "^systemVar67"
.space 1, 0
lenSystemVarName67 = . - systemVarName67
systemVarName68:
.ascii "^systemVar68"
.space 1, 0
lenSystemVarName68 = . - systemVarName68
systemVarName69:
.ascii "^systemVar69"
.space 1, 0
lenSystemVarName69 = . - systemVarName69
systemVarName70:
.ascii "^systemVar70"
.space 1, 0
lenSystemVarName70 = . - systemVarName70
systemVarName71:
.ascii "^systemVar71"
.space 1, 0
lenSystemVarName71 = . - systemVarName71
systemVarName72:
.ascii "^systemVar72"
.space 1, 0
lenSystemVarName72 = . - systemVarName72
systemVarName73:
.ascii "^systemVar73"
.space 1, 0
lenSystemVarName73 = . - systemVarName73
systemVarName74:
.ascii "^systemVar74"
.space 1, 0
lenSystemVarName74 = . - systemVarName74
systemVarName75:
.ascii "^systemVar75"
.space 1, 0
lenSystemVarName75 = . - systemVarName75
systemVarName76:
.ascii "^systemVar76"
.space 1, 0
lenSystemVarName76 = . - systemVarName76
systemVarName77:
.ascii "^systemVar77"
.space 1, 0
lenSystemVarName77 = . - systemVarName77
systemVarName78:
.ascii "^systemVar78"
.space 1, 0
lenSystemVarName78 = . - systemVarName78
systemVarName79:
.ascii "^systemVar79"
.space 1, 0
lenSystemVarName79 = . - systemVarName79
systemVarName80:
.ascii "^systemVar80"
.space 1, 0
lenSystemVarName80 = . - systemVarName80
systemVarName81:
.ascii "^systemVar81"
.space 1, 0
lenSystemVarName81 = . - systemVarName81
systemVarName82:
.ascii "^systemVar82"
.space 1, 0
lenSystemVarName82 = . - systemVarName82
systemVarName83:
.ascii "^systemVar83"
.space 1, 0
lenSystemVarName83 = . - systemVarName83
systemVarName84:
.ascii "^systemVar84"
.space 1, 0
lenSystemVarName84 = . - systemVarName84
systemVarName85:
.ascii "^systemVar85"
.space 1, 0
lenSystemVarName85 = . - systemVarName85
systemVarName86:
.ascii "^systemVar86"
.space 1, 0
lenSystemVarName86 = . - systemVarName86
systemVarName87:
.ascii "^systemVar87"
.space 1, 0
lenSystemVarName87 = . - systemVarName87
systemVarName88:
.ascii "^systemVar88"
.space 1, 0
lenSystemVarName88 = . - systemVarName88
systemVarName89:
.ascii "^systemVar89"
.space 1, 0
lenSystemVarName89 = . - systemVarName89
systemVarName90:
.ascii "^systemVar90"
.space 1, 0
lenSystemVarName90 = . - systemVarName90
systemVarName91:
.ascii "^systemVar91"
.space 1, 0
lenSystemVarName91 = . - systemVarName91
systemVarName92:
.ascii "^systemVar92"
.space 1, 0
lenSystemVarName92 = . - systemVarName92
systemVarName93:
.ascii "^systemVar93"
.space 1, 0
lenSystemVarName93 = . - systemVarName93
systemVarName94:
.ascii "^systemVar94"
.space 1, 0
lenSystemVarName94 = . - systemVarName94
systemVarName95:
.ascii "^systemVar95"
.space 1, 0
lenSystemVarName95 = . - systemVarName95
systemVarName96:
.ascii "^systemVar96"
.space 1, 0
lenSystemVarName96 = . - systemVarName96
systemVarName97:
.ascii "^systemVar97"
.space 1, 0
lenSystemVarName97 = . - systemVarName97
systemVarName98:
.ascii "^systemVar98"
.space 1, 0
lenSystemVarName98 = . - systemVarName98
systemVarName99:
.ascii "^systemVar99"
.space 1, 0
lenSystemVarName99 = . - systemVarName99
systemVarName100:
.ascii "^systemVar100"
.space 1, 0
lenSystemVarName100 = . - systemVarName100
systemVarName101:
.ascii "^systemVar101"
.space 1, 0
lenSystemVarName101 = . - systemVarName101
systemVarName102:
.ascii "^systemVar102"
.space 1, 0
lenSystemVarName102 = . - systemVarName102
systemVarName103:
.ascii "^systemVar103"
.space 1, 0
lenSystemVarName103 = . - systemVarName103
systemVarName104:
.ascii "^systemVar104"
.space 1, 0
lenSystemVarName104 = . - systemVarName104
systemVarName105:
.ascii "^systemVar105"
.space 1, 0
lenSystemVarName105 = . - systemVarName105
systemVarName106:
.ascii "^systemVar106"
.space 1, 0
lenSystemVarName106 = . - systemVarName106
systemVarName107:
.ascii "^systemVar107"
.space 1, 0
lenSystemVarName107 = . - systemVarName107
systemVarName108:
.ascii "^systemVar108"
.space 1, 0
lenSystemVarName108 = . - systemVarName108
systemVarName109:
.ascii "^systemVar109"
.space 1, 0
lenSystemVarName109 = . - systemVarName109
systemVarName110:
.ascii "^systemVar110"
.space 1, 0
lenSystemVarName110 = . - systemVarName110
systemVarName111:
.ascii "^systemVar111"
.space 1, 0
lenSystemVarName111 = . - systemVarName111
systemVarName112:
.ascii "^systemVar112"
.space 1, 0
lenSystemVarName112 = . - systemVarName112
systemVarName113:
.ascii "^systemVar113"
.space 1, 0
lenSystemVarName113 = . - systemVarName113
systemVarName114:
.ascii "^systemVar114"
.space 1, 0
lenSystemVarName114 = . - systemVarName114
systemVarName115:
.ascii "^systemVar115"
.space 1, 0
lenSystemVarName115 = . - systemVarName115
systemVarName116:
.ascii "^systemVar116"
.space 1, 0
lenSystemVarName116 = . - systemVarName116
systemVarName117:
.ascii "^systemVar117"
.space 1, 0
lenSystemVarName117 = . - systemVarName117
systemVarName118:
.ascii "^systemVar118"
.space 1, 0
lenSystemVarName118 = . - systemVarName118
systemVarName119:
.ascii "^systemVar119"
.space 1, 0
lenSystemVarName119 = . - systemVarName119
systemVarName120:
.ascii "^systemVar120"
.space 1, 0
lenSystemVarName120 = . - systemVarName120
systemVarName121:
.ascii "^systemVar121"
.space 1, 0
lenSystemVarName121 = . - systemVarName121
systemVarName122:
.ascii "^systemVar122"
.space 1, 0
lenSystemVarName122 = . - systemVarName122
systemVarName123:
.ascii "^systemVar123"
.space 1, 0
lenSystemVarName123 = . - systemVarName123
systemVarName124:
.ascii "^systemVar124"
.space 1, 0
lenSystemVarName124 = . - systemVarName124
systemVarName125:
.ascii "^systemVar125"
.space 1, 0
lenSystemVarName125 = . - systemVarName125
systemVarName126:
.ascii "^systemVar126"
.space 1, 0
lenSystemVarName126 = . - systemVarName126
systemVarName127:
.ascii "^systemVar127"
.space 1, 0
lenSystemVarName127 = . - systemVarName127
 st0: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt0 = . - st0
 st1: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt1 = . - st1
 st2: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt2 = . - st2
 st3: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt3 = . - st3
 st4: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt4 = . - st4
 st5: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt5 = . - st5
 st6: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt6 = . - st6
 st7: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt7 = . - st7
 st8: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt8 = . - st8
 st9: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt9 = . - st9
 st10: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt10 = . - st10
 st11: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt11 = . - st11
 st12: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt12 = . - st12
 st13: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt13 = . - st13
 st14: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt14 = . - st14
 st15: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt15 = . - st15
 st16: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt16 = . - st16
 st17: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt17 = . - st17
 st18: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt18 = . - st18
 st19: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt19 = . - st19
 st20: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt20 = . - st20
 st21: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt21 = . - st21
 st22: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt22 = . - st22
 st23: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt23 = . - st23
 st24: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt24 = . - st24
 st25: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt25 = . - st25
 st26: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt26 = . - st26
 st27: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt27 = . - st27
 st28: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt28 = . - st28
 st29: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt29 = . - st29
 st30: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt30 = . - st30
 st31: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt31 = . - st31
 st32: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt32 = . - st32
 st33: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt33 = . - st33
 st34: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt34 = . - st34
 st35: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt35 = . - st35
 st36: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt36 = . - st36
 st37: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt37 = . - st37
 st38: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt38 = . - st38
 st39: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt39 = . - st39
 st40: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt40 = . - st40
 st41: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt41 = . - st41
 st42: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt42 = . - st42
 st43: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt43 = . - st43
 st44: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt44 = . - st44
 st45: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt45 = . - st45
 st46: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt46 = . - st46
 st47: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt47 = . - st47
 st48: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt48 = . - st48
 st49: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt49 = . - st49
 st50: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt50 = . - st50
 st51: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt51 = . - st51
 st52: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt52 = . - st52
 st53: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt53 = . - st53
 st54: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt54 = . - st54
 st55: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt55 = . - st55
 st56: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt56 = . - st56
 st57: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt57 = . - st57
 st58: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt58 = . - st58
 st59: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt59 = . - st59
 st60: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt60 = . - st60
 st61: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt61 = . - st61
 st62: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt62 = . - st62
 st63: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt63 = . - st63
 st64: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt64 = . - st64
 st65: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt65 = . - st65
 st66: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt66 = . - st66
 st67: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt67 = . - st67
 st68: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt68 = . - st68
 st69: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt69 = . - st69
 st70: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt70 = . - st70
 st71: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt71 = . - st71
 st72: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt72 = . - st72
 st73: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt73 = . - st73
 st74: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt74 = . - st74
 st75: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt75 = . - st75
 st76: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt76 = . - st76
 st77: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt77 = . - st77
 st78: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt78 = . - st78
 st79: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt79 = . - st79
 st80: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt80 = . - st80
 st81: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt81 = . - st81
 st82: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt82 = . - st82
 st83: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt83 = . - st83
 st84: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt84 = . - st84
 st85: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt85 = . - st85
 st86: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt86 = . - st86
 st87: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt87 = . - st87
 st88: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt88 = . - st88
 st89: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt89 = . - st89
 st90: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt90 = . - st90
 st91: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt91 = . - st91
 st92: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt92 = . - st92
 st93: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt93 = . - st93
 st94: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt94 = . - st94
 st95: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt95 = . - st95
 st96: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt96 = . - st96
 st97: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt97 = . - st97
 st98: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt98 = . - st98
 st99: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt99 = . - st99
 st100: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt100 = . - st100
 st101: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt101 = . - st101
 st102: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt102 = . - st102
 st103: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt103 = . - st103
 st104: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt104 = . - st104
 st105: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt105 = . - st105
 st106: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt106 = . - st106
 st107: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt107 = . - st107
 st108: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt108 = . - st108
 st109: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt109 = . - st109
 st110: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt110 = . - st110
 st111: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt111 = . - st111
 st112: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt112 = . - st112
 st113: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt113 = . - st113
 st114: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt114 = . - st114
 st115: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt115 = . - st115
 st116: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt116 = . - st116
 st117: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt117 = . - st117
 st118: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt118 = . - st118
 st119: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt119 = . - st119
 st120: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt120 = . - st120
 st121: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt121 = . - st121
 st122: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt122 = . - st122
 st123: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt123 = . - st123
 st124: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt124 = . - st124
 st125: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt125 = . - st125
 st126: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt126 = . - st126
 st127: 
 .quad 0, 0, 0, 0, 0, 0, 0, 0 
 lenSt127 = . - st127
varName0:
.ascii "$ret"
.space 1, 0
lenVarName0 = . - varName0
varName1:
.ascii "$mod_return_var"
.space 1, 0
lenVarName1 = . - varName1
varName2:
.ascii "$mod_res"
.space 1, 0
lenVarName2 = . - varName2
label0:
 //.quad .mod
labelName0:
.ascii ".mod"
.space 1,0
data0:
.ascii ""
.space 1, 0
lenData0 = . - data0
varName3:
.ascii "b"
.space 1, 0
lenVarName3 = . - varName3
varName4:
.ascii "a"
.space 1, 0
lenVarName4 = . - varName4
data1:
.ascii ""
.space 1, 0
lenData1 = . - data1
data2:
.ascii ""
.space 1, 0
lenData2 = . - data2
varName5:
.ascii "fa"
.space 1, 0
lenVarName5 = . - varName5
varName6:
.ascii "fb"
.space 1, 0
lenVarName6 = . - varName6
varName7:
.ascii "$F0"
.space 1, 0
lenVarName7 = . - varName7
varName8:
.ascii "$F0"
.space 1, 0
lenVarName8 = . - varName8
varName9:
.ascii "c"
.space 1, 0
lenVarName9 = . - varName9
varName10:
.ascii "$I0"
.space 1, 0
lenVarName10 = . - varName10
label1:
 //.quad .mod_end
labelName1:
.ascii ".mod_end"
.space 1,0
data3:
.ascii ""
.space 1, 0
lenData3 = . - data3
varName11:
.ascii "$next_prime_return_var"
.space 1, 0
lenVarName11 = . - varName11
varName12:
.ascii "$next_prime_res"
.space 1, 0
lenVarName12 = . - varName12
label2:
 //.quad .next_prime
labelName2:
.ascii ".next_prime"
.space 1,0
data4:
.ascii ""
.space 1, 0
lenData4 = . - data4
varName13:
.ascii "n"
.space 1, 0
lenVarName13 = . - varName13
data5:
.ascii ""
.space 1, 0
lenData5 = . - data5
data6:
.ascii ""
.space 1, 0
lenData6 = . - data6
data7:
.ascii ""
.space 1, 0
lenData7 = . - data7
data8:
.ascii ""
.space 1, 0
lenData8 = . - data8
data9:
.ascii ""
.space 1, 0
lenData9 = . - data9
varName14:
.ascii "$for0"
.space 1, 0
lenVarName14 = . - varName14
label3:
 //.quad ._for0
labelName3:
.ascii "._for0"
.space 1,0
data10:
.ascii ""
.space 1, 0
lenData10 = . - data10
data11:
.ascii ""
.space 1, 0
lenData11 = . - data11
data12:
.ascii "1"
.space 1, 0
lenData12 = . - data12
varName15:
.ascii "is_prime"
.space 1, 0
lenVarName15 = . - varName15
data13:
.ascii "1"
.space 1, 0
lenData13 = . - data13
data14:
.ascii ""
.space 1, 0
lenData14 = . - data14
varName16:
.ascii "i"
.space 1, 0
lenVarName16 = . - varName16
data15:
.ascii "2"
.space 1, 0
lenData15 = . - data15
varName17:
.ascii "$for3"
.space 1, 0
lenVarName17 = . - varName17
label4:
 //.quad ._for3
labelName4:
.ascii "._for3"
.space 1,0
data16:
.ascii ""
.space 1, 0
lenData16 = . - data16
data17:
.ascii "2"
.space 1, 0
lenData17 = . - data17
data18:
.ascii "1"
.space 1, 0
lenData18 = . - data18
data19:
.ascii ""
.space 1, 0
lenData19 = . - data19
varName18:
.ascii "$mod_res0"
.space 1, 0
lenVarName18 = . - varName18
data20:
.ascii "#mod_res0"
.space 1, 0
lenData20 = . - data20
label5:
 //.quad .mod_res0
labelName5:
.ascii ".mod_res0"
.space 1,0
data21:
.ascii "0"
.space 1, 0
lenData21 = . - data21
data22:
.ascii ""
.space 1, 0
lenData22 = . - data22
data23:
.ascii "0"
.space 1, 0
lenData23 = . - data23
data24:
.ascii "1"
.space 1, 0
lenData24 = . - data24
label6:
 //.quad ._cond4_end
labelName6:
.ascii "._cond4_end"
.space 1,0
data25:
.ascii ""
.space 1, 0
lenData25 = . - data25
data26:
.ascii ""
.space 1, 0
lenData26 = . - data26
data27:
.ascii "1"
.space 1, 0
lenData27 = . - data27
label7:
 //.quad ._undef_for3
labelName7:
.ascii "._undef_for3"
.space 1,0
data28:
.ascii ""
.space 1, 0
lenData28 = . - data28
label8:
 //.quad ._cond3_end
labelName8:
.ascii "._cond3_end"
.space 1,0
data29:
.ascii ""
.space 1, 0
lenData29 = . - data29
data30:
.ascii ""
.space 1, 0
lenData30 = . - data30
label9:
 //.quad ._cond5_end
labelName9:
.ascii "._cond5_end"
.space 1,0
data31:
.ascii ""
.space 1, 0
lenData31 = . - data31
data32:
.ascii "2"
.space 1, 0
lenData32 = . - data32
data33:
.ascii "1"
.space 1, 0
lenData33 = . - data33
data34:
.ascii ""
.space 1, 0
lenData34 = . - data34
label10:
 //.quad ._cond6_end
labelName10:
.ascii "._cond6_end"
.space 1,0
data35:
.ascii ""
.space 1, 0
lenData35 = . - data35
label11:
 //.quad ._for3_end
labelName11:
.ascii "._for3_end"
.space 1,0
data36:
.ascii ""
.space 1, 0
lenData36 = . - data36
label12:
 //.quad ._cond2_end
labelName12:
.ascii "._cond2_end"
.space 1,0
data37:
.ascii ""
.space 1, 0
lenData37 = . - data37
data38:
.ascii ""
.space 1, 0
lenData38 = . - data38
data39:
.ascii ""
.space 1, 0
lenData39 = . - data39
label13:
 //.quad ._cond7_end
labelName13:
.ascii "._cond7_end"
.space 1,0
data40:
.ascii ""
.space 1, 0
lenData40 = . - data40
data41:
.ascii ""
.space 1, 0
lenData41 = . - data41
data42:
.ascii ""
.space 1, 0
lenData42 = . - data42
label14:
 //.quad ._undef_for0
labelName14:
.ascii "._undef_for0"
.space 1,0
data43:
.ascii ""
.space 1, 0
lenData43 = . - data43
label15:
 //.quad ._cond1_end
labelName15:
.ascii "._cond1_end"
.space 1,0
data44:
.ascii ""
.space 1, 0
lenData44 = . - data44
data45:
.ascii ""
.space 1, 0
lenData45 = . - data45
label16:
 //.quad ._cond8_end
labelName16:
.ascii "._cond8_end"
.space 1,0
data46:
.ascii ""
.space 1, 0
lenData46 = . - data46
data47:
.ascii ""
.space 1, 0
lenData47 = . - data47
label17:
 //.quad ._cond9_end
labelName17:
.ascii "._cond9_end"
.space 1,0
data48:
.ascii ""
.space 1, 0
lenData48 = . - data48
label18:
 //.quad ._for0_end
labelName18:
.ascii "._for0_end"
.space 1,0
data49:
.ascii ""
.space 1, 0
lenData49 = . - data49
label19:
 //.quad ._cond0_end
labelName19:
.ascii "._cond0_end"
.space 1,0
data50:
.ascii ""
.space 1, 0
lenData50 = . - data50
data51:
.ascii ""
.space 1, 0
lenData51 = . - data51
data52:
.ascii "-1"
.space 1, 0
lenData52 = . - data52
label20:
 //.quad .next_prime_end
labelName20:
.ascii ".next_prime_end"
.space 1,0
data53:
.ascii ""
.space 1, 0
lenData53 = . - data53
varName19:
.ascii "$main_return_var"
.space 1, 0
lenVarName19 = . - varName19
varName20:
.ascii "$main_res"
.space 1, 0
lenVarName20 = . - varName20
label21:
 //.quad .main
labelName21:
.ascii ".main"
.space 1,0
data54:
.ascii ""
.space 1, 0
lenData54 = . - data54
data55:
.ascii ""
.space 1, 0
lenData55 = . - data55
varName21:
.ascii "$print_arg0"
.space 1, 0
lenVarName21 = . - varName21
data56:
.ascii "Программа вычисления таблицы простых чисел\n"
.space 1, 0
lenData56 = . - data56
varName22:
.ascii "$print_arg0"
.space 1, 0
lenVarName22 = . - varName22
data57:
.ascii "Задайте количество чисел (не более 25)\n"
.space 1, 0
lenData57 = . - data57
varName23:
.ascii "s"
.space 1, 0
lenVarName23 = . - varName23
varName24:
.ascii "n"
.space 1, 0
lenVarName24 = . - varName24
data58:
.ascii ""
.space 1, 0
lenData58 = . - data58
data59:
.ascii ""
.space 1, 0
lenData59 = . - data59
data60:
.ascii ""
.space 1, 0
lenData60 = . - data60
varName25:
.ascii "$for1"
.space 1, 0
lenVarName25 = . - varName25
label22:
 //.quad ._for1
labelName22:
.ascii "._for1"
.space 1,0
data61:
.ascii ""
.space 1, 0
lenData61 = . - data61
data62:
.ascii ""
.space 1, 0
lenData62 = . - data62
data63:
.ascii ""
.space 1, 0
lenData63 = . - data63
data64:
.ascii "0"
.space 1, 0
lenData64 = . - data64
data65:
.ascii ""
.space 1, 0
lenData65 = . - data65
data66:
.ascii ""
.space 1, 0
lenData66 = . - data66
data67:
.ascii ""
.space 1, 0
lenData67 = . - data67
label23:
 //.quad ._cond13_end
labelName23:
.ascii "._cond13_end"
.space 1,0
data68:
.ascii ""
.space 1, 0
lenData68 = . - data68
data69:
.ascii ""
.space 1, 0
lenData69 = . - data69
data70:
.ascii ""
.space 1, 0
lenData70 = . - data70
label24:
 //.quad ._cond14_end
labelName24:
.ascii "._cond14_end"
.space 1,0
data71:
.ascii ""
.space 1, 0
lenData71 = . - data71
varName26:
.ascii "$I0"
.space 1, 0
lenVarName26 = . - varName26
label25:
// .quad ._attempt0
labelName25:
.ascii "._attempt0"
.space 1,0
data72:
.ascii ""
.space 1, 0
lenData72 = . - data72
data73:
.ascii "1"
.space 1, 0
lenData73 = . - data73
label26:
 //.quad ._cond12_end
labelName26:
.ascii "._cond12_end"
.space 1,0
data74:
.ascii ""
.space 1, 0
lenData74 = . - data74
data75:
.ascii ""
.space 1, 0
lenData75 = . - data75
data76:
.ascii ""
.space 1, 0
lenData76 = . - data76
data77:
.ascii ""
.space 1, 0
lenData77 = . - data77
varName27:
.ascii "$print_arg0"
.space 1, 0
lenVarName27 = . - varName27
data78:
.ascii "Ошибка! Неправильный формат числа\n"
.space 1, 0
lenData78 = . - data78
varName28:
.ascii "$print_arg0"
.space 1, 0
lenVarName28 = . - varName28
varName29:
.ascii "$print_arg0"
.space 1, 0
lenVarName29 = . - varName29
data79:
.ascii "\n"
.space 1, 0
lenData79 = . - data79
varName30:
.ascii "$print_arg0"
.space 1, 0
lenVarName30 = . - varName30
data80:
.ascii "Повторите ввод\n"
.space 1, 0
lenData80 = . - data80
label27:
 //.quad ._cond15_end
labelName27:
.ascii "._cond15_end"
.space 1,0
data81:
.ascii ""
.space 1, 0
lenData81 = . - data81
data82:
.ascii "0"
.space 1, 0
lenData82 = . - data82
data83:
.ascii ""
.space 1, 0
lenData83 = . - data83
data84:
.ascii "Введено отрицательное число"
.space 1, 0
lenData84 = . - data84
data85:
.ascii ""
.space 1, 0
lenData85 = . - data85
varName31:
.ascii "$print_arg0"
.space 1, 0
lenVarName31 = . - varName31
data86:
.ascii "Ошибка! "
.space 1, 0
lenData86 = . - data86
data87:
.ascii "\n"
.space 1, 0
lenData87 = . - data87
varName32:
.ascii "$print_arg0"
.space 1, 0
lenVarName32 = . - varName32
data88:
.ascii "Повторите ввод\n"
.space 1, 0
lenData88 = . - data88
label28:
 //.quad ._cond16_end
labelName28:
.ascii "._cond16_end"
.space 1,0
data89:
.ascii ""
.space 1, 0
lenData89 = . - data89
data90:
.ascii "25"
.space 1, 0
lenData90 = . - data90
data91:
.ascii ""
.space 1, 0
lenData91 = . - data91
data92:
.ascii "Слишком большое количество чисел"
.space 1, 0
lenData92 = . - data92
data93:
.ascii ""
.space 1, 0
lenData93 = . - data93
varName33:
.ascii "$print_arg0"
.space 1, 0
lenVarName33 = . - varName33
data94:
.ascii "Ошибка! "
.space 1, 0
lenData94 = . - data94
data95:
.ascii "\n"
.space 1, 0
lenData95 = . - data95
varName34:
.ascii "$print_arg0"
.space 1, 0
lenVarName34 = . - varName34
data96:
.ascii "Повторите ввод\n"
.space 1, 0
lenData96 = . - data96
label29:
 //.quad ._cond17_end
labelName29:
.ascii "._cond17_end"
.space 1,0
data97:
.ascii ""
.space 1, 0
lenData97 = . - data97
label30:
 //.quad ._cond_exit0
labelName30:
.ascii "._cond_exit0"
.space 1,0
data98:
.ascii ""
.space 1, 0
lenData98 = . - data98
data99:
.ascii ""
.space 1, 0
lenData99 = . - data99
data100:
.ascii ""
.space 1, 0
lenData100 = . - data100
data101:
.ascii ""
.space 1, 0
lenData101 = . - data101
data102:
.ascii ""
.space 1, 0
lenData102 = . - data102
data103:
.ascii "1"
.space 1, 0
lenData103 = . - data103
label31:
 //.quad ._cond18_end
labelName31:
.ascii "._cond18_end"
.space 1,0
data104:
.ascii ""
.space 1, 0
lenData104 = . - data104
data105:
.ascii ""
.space 1, 0
lenData105 = . - data105
label32:
 //.quad ._undef_for1
labelName32:
.ascii "._undef_for1"
.space 1,0
data106:
.ascii ""
.space 1, 0
lenData106 = . - data106
label33:
 //.quad ._cond11_end
labelName33:
.ascii "._cond11_end"
.space 1,0
data107:
.ascii ""
.space 1, 0
lenData107 = . - data107
data108:
.ascii ""
.space 1, 0
lenData108 = . - data108
label34:
 //.quad ._cond19_end
labelName34:
.ascii "._cond19_end"
.space 1,0
data109:
.ascii ""
.space 1, 0
lenData109 = . - data109
data110:
.ascii ""
.space 1, 0
lenData110 = . - data110
label35:
 //.quad ._cond20_end
labelName35:
.ascii "._cond20_end"
.space 1,0
data111:
.ascii ""
.space 1, 0
lenData111 = . - data111
label36:
 //.quad ._for1_end
labelName36:
.ascii "._for1_end"
.space 1,0
data112:
.ascii ""
.space 1, 0
lenData112 = . - data112
label37:
 //.quad ._cond10_end
labelName37:
.ascii "._cond10_end"
.space 1,0
data113:
.ascii ""
.space 1, 0
lenData113 = . - data113
data114:
.ascii ""
.space 1, 0
lenData114 = . - data114
data115:
.ascii "0"
.space 1, 0
lenData115 = . - data115
data116:
.ascii ""
.space 1, 0
lenData116 = . - data116
data117:
.ascii "0"
.space 1, 0
lenData117 = . - data117
label38:
 //.quad ._cond21_end
labelName38:
.ascii "._cond21_end"
.space 1,0
data118:
.ascii ""
.space 1, 0
lenData118 = . - data118
data119:
.ascii ""
.space 1, 0
lenData119 = . - data119
varName35:
.ascii "p"
.space 1, 0
lenVarName35 = . - varName35
data120:
.ascii "2"
.space 1, 0
lenData120 = . - data120
data121:
.ascii "1"
.space 1, 0
lenData121 = . - data121
data122:
.ascii ""
.space 1, 0
lenData122 = . - data122
varName36:
.ascii "$s0"
.space 1, 0
lenVarName36 = . - varName36
varName37:
.ascii "$print_arg0"
.space 1, 0
lenVarName37 = . - varName37
varName38:
.ascii "$print_arg0"
.space 1, 0
lenVarName38 = . - varName38
data123:
.ascii "\n"
.space 1, 0
lenData123 = . - data123
data124:
.ascii "0"
.space 1, 0
lenData124 = . - data124
label39:
 //.quad ._cond22_end
labelName39:
.ascii "._cond22_end"
.space 1,0
data125:
.ascii ""
.space 1, 0
lenData125 = . - data125
data126:
.ascii ""
.space 1, 0
lenData126 = . - data126
varName39:
.ascii "$s0"
.space 1, 0
lenVarName39 = . - varName39
varName40:
.ascii "$print_arg0"
.space 1, 0
lenVarName40 = . - varName40
data127:
.ascii ", "
.space 1, 0
lenData127 = . - data127
data128:
.ascii ""
.space 1, 0
lenData128 = . - data128
varName41:
.ascii "i"
.space 1, 0
lenVarName41 = . - varName41
data129:
.ascii "0"
.space 1, 0
lenData129 = . - data129
varName42:
.ascii "$for2"
.space 1, 0
lenVarName42 = . - varName42
label40:
 //.quad ._for2
labelName40:
.ascii "._for2"
.space 1,0
data130:
.ascii ""
.space 1, 0
lenData130 = . - data130
data131:
.ascii "1"
.space 1, 0
lenData131 = . - data131
data132:
.ascii ""
.space 1, 0
lenData132 = . - data132
varName43:
.ascii "$next_prime_res0"
.space 1, 0
lenVarName43 = . - varName43
data133:
.ascii "#next_prime_res0"
.space 1, 0
lenData133 = . - data133
label41:
 //.quad .next_prime_res0
labelName41:
.ascii ".next_prime_res0"
.space 1,0
varName44:
.ascii "$s0"
.space 1, 0
lenVarName44 = . - varName44
varName45:
.ascii "$print_arg0"
.space 1, 0
lenVarName45 = . - varName45
data134:
.ascii "2"
.space 1, 0
lenData134 = . - data134
data135:
.ascii ""
.space 1, 0
lenData135 = . - data135
varName46:
.ascii "$print_arg0"
.space 1, 0
lenVarName46 = . - varName46
data136:
.ascii ", "
.space 1, 0
lenData136 = . - data136
label42:
 //.quad ._cond25_end
labelName42:
.ascii "._cond25_end"
.space 1,0
data137:
.ascii ""
.space 1, 0
lenData137 = . - data137
data138:
.ascii ""
.space 1, 0
lenData138 = . - data138
data139:
.ascii "1"
.space 1, 0
lenData139 = . - data139
label43:
 //.quad ._undef_for2
labelName43:
.ascii "._undef_for2"
.space 1,0
data140:
.ascii ""
.space 1, 0
lenData140 = . - data140
label44:
 //.quad ._cond24_end
labelName44:
.ascii "._cond24_end"
.space 1,0
data141:
.ascii ""
.space 1, 0
lenData141 = . - data141
data142:
.ascii ""
.space 1, 0
lenData142 = . - data142
label45:
 //.quad ._cond26_end
labelName45:
.ascii "._cond26_end"
.space 1,0
data143:
.ascii ""
.space 1, 0
lenData143 = . - data143
data144:
.ascii "1"
.space 1, 0
lenData144 = . - data144
data145:
.ascii ""
.space 1, 0
lenData145 = . - data145
label46:
 //.quad ._cond27_end
labelName46:
.ascii "._cond27_end"
.space 1,0
data146:
.ascii ""
.space 1, 0
lenData146 = . - data146
label47:
 //.quad ._for2_end
labelName47:
.ascii "._for2_end"
.space 1,0
data147:
.ascii ""
.space 1, 0
lenData147 = . - data147
label48:
// .quad ._cond23_end
labelName48:
.ascii "._cond23_end"
.space 1,0
data148:
.ascii ""
.space 1, 0
lenData148 = . - data148
data149:
.ascii ""
.space 1, 0
lenData149 = . - data149
varName47:
.ascii "$print_arg0"
.space 1, 0
lenVarName47 = . - varName47
data150:
.ascii "\n"
.space 1, 0
lenData150 = . - data150
data151:
.ascii "0"
.space 1, 0
lenData151 = . - data151
label49:
 //.quad .main_end
labelName49:
.ascii ".main_end"
.space 1,0
data152:
.ascii ""
.space 1, 0
lenData152 = . - data152
varName48:
.ascii "res"
.space 1, 0
lenVarName48 = . - varName48
varName49:
.ascii "$main_res0"
.space 1, 0
lenVarName49 = . - varName49
varName50:
.ascii "myStack"
.space 1, 0
lenVarName50 = . - varName50 
data153:
.ascii "#main_res0"
.space 1, 0
lenData153 = . - data153
data154:
.ascii "/home/slava/Go/bint/test.txt"
.space 1, 0
lenData154 = . - data154
label50:
// .quad .main_res0
labelName50:
.ascii ".main_res0"
.space 1,0

.text 
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
 cmp (stackMax), %r8 
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
 cmp $0, %r11b 
 jz __defOkLocalEx
 cmp $1, %r11b
 jz __defOkLocalEx
 movb %r11b, (%r8)
 inc %rcx 
 inc %r8 
 jmp __defOkLocal
 __defOkLocalEx:
 movb $0, (%r8)
 mov %r14, %r8 
 add (varNameSize), %r8 
 __defOkTypeLocal:
 movb (%rdx), %r11b
 cmp $0, %r11b 
 jz __defOkTypeLocalEx
 cmp $1, %r11b 
 jz __defOkTypeLocalEx
 movb %r11b, (%r8)
 inc %rdx
 inc %r8 
 jmp __defOkTypeLocal
 __defOkTypeLocalEx:
 movb $0, (%r8)

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

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenStackType, %rax 
 mov $stackType, %rdi 
 call __set 
 call __compare
 cmp $1, %rax 
 jz __defStack  

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
 mov %rax, (strPointer)
 jmp __defEnd 
 



 __defStack:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 mov (stackPointer), %rax 
 mov $stackEndSymbol, %r12
 __defStackEndSymbol: 
 mov (%r12), %dil 
 cmp $0, %dil 
 jz __defStackEndSymbolOk
 mov %dil, (%rax)
 inc %r12 
 jmp __defStackEndSymbol
 __defStackEndSymbolOk:
 movb $0, (%rax)
 mov %r8, %r12 
 mov (stackPointer), %rax 
 call __toStr 
 mov %r12, %r8 
 mov $buf2, %rax
 __defStackAddr:
 mov (%rax), %dl 
 cmp $0, %dl  
 jz __defStackEnd 
 mov %dl, (%r8)
 inc %rax 
 inc %r8 
 jmp __defStackAddr
 __defStackEnd:
 movb $0, (%r8)
 mov (stackPointer), %rax 
 add (valSize), %rax 
 mov %rax, (stackPointer)
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

 __firstStackMem:
 # адрес начала области для выделения памяти
 mov (strMax), %rax
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, (stackBegin)
 mov %rax, (stackPointer)
 mov %rax, %r9 
 add (stackShiftSize), %r9
 mov %r9, (stackMax)
# выделить динамическую память
 mov (stackShiftSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __firstStackMemLo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (stackShiftSize), %rbx
 jz  __firstStackMemEx
 jmp __firstStackMemLo
 __firstStackMemEx:
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
 call __firstStackMem

 # toPanic variable
 /*mov $lenVarName, %rsi 
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

mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName0, %rax 
 mov $systemVarName0, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName1, %rax 
 mov $systemVarName1, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName2, %rax 
 mov $systemVarName2, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName3, %rax 
 mov $systemVarName3, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName4, %rax 
 mov $systemVarName4, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName5, %rax 
 mov $systemVarName5, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName6, %rax 
 mov $systemVarName6, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName7, %rax 
 mov $systemVarName7, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName8, %rax 
 mov $systemVarName8, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName9, %rax 
 mov $systemVarName9, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName10, %rax 
 mov $systemVarName10, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName11, %rax 
 mov $systemVarName11, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName12, %rax 
 mov $systemVarName12, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName13, %rax 
 mov $systemVarName13, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName14, %rax 
 mov $systemVarName14, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName15, %rax 
 mov $systemVarName15, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName16, %rax 
 mov $systemVarName16, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName17, %rax 
 mov $systemVarName17, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName18, %rax 
 mov $systemVarName18, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName19, %rax 
 mov $systemVarName19, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName20, %rax 
 mov $systemVarName20, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName21, %rax 
 mov $systemVarName21, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName22, %rax 
 mov $systemVarName22, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName23, %rax 
 mov $systemVarName23, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName24, %rax 
 mov $systemVarName24, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName25, %rax 
 mov $systemVarName25, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName26, %rax 
 mov $systemVarName26, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName27, %rax 
 mov $systemVarName27, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName28, %rax 
 mov $systemVarName28, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName29, %rax 
 mov $systemVarName29, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName30, %rax 
 mov $systemVarName30, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName31, %rax 
 mov $systemVarName31, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName32, %rax 
 mov $systemVarName32, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName33, %rax 
 mov $systemVarName33, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName34, %rax 
 mov $systemVarName34, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName35, %rax 
 mov $systemVarName35, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName36, %rax 
 mov $systemVarName36, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName37, %rax 
 mov $systemVarName37, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName38, %rax 
 mov $systemVarName38, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName39, %rax 
 mov $systemVarName39, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName40, %rax 
 mov $systemVarName40, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName41, %rax 
 mov $systemVarName41, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName42, %rax 
 mov $systemVarName42, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName43, %rax 
 mov $systemVarName43, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName44, %rax 
 mov $systemVarName44, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName45, %rax 
 mov $systemVarName45, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName46, %rax 
 mov $systemVarName46, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName47, %rax 
 mov $systemVarName47, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName48, %rax 
 mov $systemVarName48, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName49, %rax 
 mov $systemVarName49, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName50, %rax 
 mov $systemVarName50, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName51, %rax 
 mov $systemVarName51, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName52, %rax 
 mov $systemVarName52, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName53, %rax 
 mov $systemVarName53, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName54, %rax 
 mov $systemVarName54, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName55, %rax 
 mov $systemVarName55, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName56, %rax 
 mov $systemVarName56, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName57, %rax 
 mov $systemVarName57, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName58, %rax 
 mov $systemVarName58, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName59, %rax 
 mov $systemVarName59, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName60, %rax 
 mov $systemVarName60, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName61, %rax 
 mov $systemVarName61, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName62, %rax 
 mov $systemVarName62, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName63, %rax 
 mov $systemVarName63, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName64, %rax 
 mov $systemVarName64, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName65, %rax 
 mov $systemVarName65, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName66, %rax 
 mov $systemVarName66, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName67, %rax 
 mov $systemVarName67, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName68, %rax 
 mov $systemVarName68, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName69, %rax 
 mov $systemVarName69, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName70, %rax 
 mov $systemVarName70, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName71, %rax 
 mov $systemVarName71, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName72, %rax 
 mov $systemVarName72, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName73, %rax 
 mov $systemVarName73, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName74, %rax 
 mov $systemVarName74, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName75, %rax 
 mov $systemVarName75, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName76, %rax 
 mov $systemVarName76, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName77, %rax 
 mov $systemVarName77, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName78, %rax 
 mov $systemVarName78, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName79, %rax 
 mov $systemVarName79, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName80, %rax 
 mov $systemVarName80, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName81, %rax 
 mov $systemVarName81, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName82, %rax 
 mov $systemVarName82, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName83, %rax 
 mov $systemVarName83, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName84, %rax 
 mov $systemVarName84, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName85, %rax 
 mov $systemVarName85, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName86, %rax 
 mov $systemVarName86, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName87, %rax 
 mov $systemVarName87, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName88, %rax 
 mov $systemVarName88, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName89, %rax 
 mov $systemVarName89, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName90, %rax 
 mov $systemVarName90, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName91, %rax 
 mov $systemVarName91, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName92, %rax 
 mov $systemVarName92, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName93, %rax 
 mov $systemVarName93, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName94, %rax 
 mov $systemVarName94, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName95, %rax 
 mov $systemVarName95, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName96, %rax 
 mov $systemVarName96, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName97, %rax 
 mov $systemVarName97, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName98, %rax 
 mov $systemVarName98, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName99, %rax 
 mov $systemVarName99, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName100, %rax 
 mov $systemVarName100, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName101, %rax 
 mov $systemVarName101, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName102, %rax 
 mov $systemVarName102, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName103, %rax 
 mov $systemVarName103, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName104, %rax 
 mov $systemVarName104, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName105, %rax 
 mov $systemVarName105, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName106, %rax 
 mov $systemVarName106, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName107, %rax 
 mov $systemVarName107, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName108, %rax 
 mov $systemVarName108, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName109, %rax 
 mov $systemVarName109, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName110, %rax 
 mov $systemVarName110, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName111, %rax 
 mov $systemVarName111, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName112, %rax 
 mov $systemVarName112, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName113, %rax 
 mov $systemVarName113, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName114, %rax 
 mov $systemVarName114, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName115, %rax 
 mov $systemVarName115, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName116, %rax 
 mov $systemVarName116, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName117, %rax 
 mov $systemVarName117, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName118, %rax 
 mov $systemVarName118, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName119, %rax 
 mov $systemVarName119, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName120, %rax 
 mov $systemVarName120, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName121, %rax 
 mov $systemVarName121, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName122, %rax 
 mov $systemVarName122, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName123, %rax 
 mov $systemVarName123, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName124, %rax 
 mov $systemVarName124, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName125, %rax 
 mov $systemVarName125, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName126, %rax 
 mov $systemVarName126, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName127, %rax 
 mov $systemVarName127, %rdi
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
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stringType, %rdi
 call __set 
 call __defineVar
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
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
jmp .mod_end
.mod:

mov $data0, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName3, %rax 
 mov $varName3, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName4, %rax 
 mov $varName4, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $data1, %rsi
call __print
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
mov $lenVarType, %rsi
mov $varType, %rdx
mov $lenStringType, %rax
mov $stringType, %rdi
call __set
mov $varName, %rcx
mov $varType, %rdx
call __defineVar
pop (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenPopVarName, %rax 
 mov $popVarName, %rdi
 call __set 
 xor %rax, %rax 
 call __setVar
mov $popVarName, %rax 
mov %rax, (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName3, %rax 
 mov $varName3, %rdi
 call __set 
 mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
call __undefineVar
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
mov $lenVarType, %rsi
mov $varType, %rdx
mov $lenStringType, %rax
mov $stringType, %rdi
call __set
mov $varName, %rcx
mov $varType, %rdx
call __defineVar
pop (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenPopVarName, %rax 
 mov $popVarName, %rdi
 call __set 
 xor %rax, %rax 
 call __setVar
mov $popVarName, %rax 
mov %rax, (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
call __undefineVar
mov $data2, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName5, %rax 
 mov $varName5, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName6, %rax 
 mov $varName6, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName7, %rax 
 mov $varName7, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov (userData), %rdi
 call __set
 call __userParseFloat
 
 movss %xmm0, (buf) 
 call __floatToStr 
  
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName7, %rax 
 mov $varName7, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName7, %rax 
 mov $varName7, %rdi
 call __set 
 mov $t0, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName5, %rax 
 mov $varName5, %rdi
 call __set
mov $varName7, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName7, %rax 
mov $varName7, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName8, %rax 
 mov $varName8, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName3, %rax 
 mov $varName3, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov (userData), %rdi
 call __set
 call __userParseFloat
 
 movss %xmm0, (buf) 
 call __floatToStr 
  
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName8, %rax 
 mov $varName8, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName8, %rax 
 mov $varName8, %rdi
 call __set 
 mov $t0, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set
mov $varName8, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set
mov $varName8, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName8, %rax 
mov $varName8, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName9, %rax 
 mov $varName9, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName10, %rax 
 mov $varName10, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName5, %rax 
 mov $varName5, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 mov $1, %rax 

 call __div 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenT0, %rax 
 mov $t0, %rdi
 call __set
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set 
 call __userToNumber
 call __toStr 
  
 mov $lenT1, %rsi 
 mov $t1, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName10, %rax 
 mov $varName10, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName10, %rax 
 mov $varName10, %rdi
 call __set 
 mov $t1, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName9, %rax 
 mov $varName9, %rdi
 call __set
mov $varName10, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName10, %rax 
mov $varName10, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName3, %rax 
 mov $varName3, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __divI 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenT0, %rax 
 mov $t0, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName3, %rax 
 mov $varName3, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __mul 
 mov $lenT1, %rsi 
 mov $t1, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenT1, %rax 
 mov $t1, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __sub 
 mov $lenT2, %rsi 
 mov $t2, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenSt0, %rsi
mov $st0, %rdx
mov $lenSt0, %rsi
mov $t2, %rdi
call __set 
push $st0
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName9, %rax 
mov $varName9, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName6, %rax 
mov $varName6, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName5, %rax 
mov $varName5, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName3, %rax 
mov $varName3, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
mov $varName4, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi 
 call __set 
 call __getVar
mov (userData), %rdi 
 movb $'.', (%rdi) 
 jmp __goto
.mod_end:

mov $data3, %rsi
call __print
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
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
jmp .next_prime_end
.next_prime:

mov $data4, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName13, %rax 
 mov $varName13, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $data5, %rsi
call __print
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
mov $lenVarType, %rsi
mov $varType, %rdx
mov $lenStringType, %rax
mov $stringType, %rdi
call __set
mov $varName, %rcx
mov $varType, %rdx
call __defineVar
pop (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenPopVarName, %rax 
 mov $popVarName, %rdi
 call __set 
 xor %rax, %rax 
 call __setVar
mov $popVarName, %rax 
mov %rax, (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
 mov $varName13, %rdi
 call __set 
 mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
call __undefineVar
mov $data6, %rsi
call __print
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right0
mov $data7, %rsi
call __print
jmp __rightEnd0
 __right0:
jmp ._cond0_end
__rightEnd0:
mov $data8, %rsi
call __print
mov $data9, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName14, %rax 
 mov $varName14, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
._for0:

mov $data10, %rsi
call __print
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right1
mov $data11, %rsi
call __print
jmp __rightEnd1
 __right1:
jmp ._cond1_end
__rightEnd1:
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData12, %rax 
 mov $data12, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
 mov $varName13, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __add 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
 mov $varName13, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
 mov $varName13, %rdi
 call __set 
 mov $t0, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName15, %rax 
 mov $varName15, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName15, %rax 
 mov $varName15, %rdi 
call __set

 mov $data13, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right2
mov $data14, %rsi
call __print
jmp __rightEnd2
 __right2:
jmp ._cond2_end
__rightEnd2:
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName16, %rax 
 mov $varName16, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
 mov $varName16, %rdi 
call __set

 mov $data15, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName17, %rax 
 mov $varName17, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
._for3:

mov $data16, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData17, %rax 
 mov $data17, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
 mov $varName13, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __divI 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData18, %rax 
 mov $data18, %rdi
 call __set
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenT0, %rax 
 mov $t0, %rdi
 call __set
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __add 
 mov $lenT1, %rsi 
 mov $t1, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenT1, %rax 
 mov $t1, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
 mov $varName16, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __less 
 mov $lenT2, %rsi 
 mov $t2, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right3
mov $data19, %rsi
call __print
jmp __rightEnd3
 __right3:
jmp ._cond3_end
__rightEnd3:
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName18, %rax 
 mov $varName18, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov $varName13, %rdi
 call __set 
 call __getVar 
 push (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov $varName16, %rdi
 call __set 
 call __getVar 
 push (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi 
call __set

 mov $data20, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
jmp .mod
.mod_res0:

mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
mov $lenVarType, %rsi
mov $varType, %rdx
mov $lenStringType, %rax
mov $stringType, %rdi
call __set
mov $varName, %rcx
mov $varType, %rdx
call __defineVar
pop (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenPopVarName, %rax 
 mov $popVarName, %rdi
 call __set 
 xor %rax, %rax 
 call __setVar
mov $popVarName, %rax 
mov %rax, (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi
 call __set 
 mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName18, %rax 
 mov $varName18, %rdi
 call __set
mov $varName2, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData21, %rax 
 mov $data21, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName18, %rax 
 mov $varName18, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __eq 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right4
mov $data22, %rsi
call __print
jmp __rightEnd4
 __right4:
jmp ._cond4_end
__rightEnd4:
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName18, %rax 
mov $varName18, %rdi 
 call __set 
call __undefineVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName15, %rax 
 mov $varName15, %rdi 
call __set

 mov $data23, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName17, %rax 
 mov $varName17, %rdi 
call __set

 mov $data24, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
jmp ._undef_for3
._cond4_end:

mov $data25, %rsi
call __print
mov $data26, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData27, %rax 
 mov $data27, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
 mov $varName16, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __add 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
 mov $varName16, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
 mov $varName16, %rdi
 call __set 
 mov $t0, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
._undef_for3:

mov $data28, %rsi
call __print
._cond3_end:

mov $data29, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName17, %rax 
 mov $varName17, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 mov (%rsi), %rdx 
 mov %rdx, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right5
mov $data30, %rsi
call __print
jmp __rightEnd5
 __right5:
jmp ._cond5_end
__rightEnd5:
jmp ._for3_end
._cond5_end:

mov $data31, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData32, %rax 
 mov $data32, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
 mov $varName13, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __divI 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData33, %rax 
 mov $data33, %rdi
 call __set
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenT0, %rax 
 mov $t0, %rdi
 call __set
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __add 
 mov $lenT1, %rsi 
 mov $t1, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenT1, %rax 
 mov $t1, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
 mov $varName16, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __less 
 mov $lenT2, %rsi 
 mov $t2, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right6
mov $data34, %rsi
call __print
jmp __rightEnd6
 __right6:
jmp ._cond6_end
__rightEnd6:
 call __clearVars
jmp ._for3
._cond6_end:

mov $data35, %rsi
call __print
._for3_end:

mov $data36, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName17, %rax 
mov $varName17, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName16, %rax 
mov $varName16, %rdi 
 call __set 
call __undefineVar
._cond2_end:

mov $data37, %rsi
call __print
mov $data38, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName15, %rax 
 mov $varName15, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 mov (%rsi), %rdx 
 mov %rdx, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right7
mov $data39, %rsi
call __print
jmp __rightEnd7
 __right7:
jmp ._cond7_end
__rightEnd7:
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName14, %rax 
mov $varName14, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName15, %rax 
mov $varName15, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov $varName13, %rdi
 call __set 
 call __getVar 
 push (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
mov $varName13, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName11, %rax 
 mov $varName11, %rdi 
 call __set 
 call __getVar
mov (userData), %rdi 
 movb $'.', (%rdi) 
 jmp __goto
._cond7_end:

mov $data40, %rsi
call __print
mov $data41, %rsi
call __print
mov $data42, %rsi
call __print
._undef_for0:

mov $data43, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName15, %rax 
mov $varName15, %rdi 
 call __set 
call __undefineVar
._cond1_end:

mov $data44, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName14, %rax 
 mov $varName14, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 mov (%rsi), %rdx 
 mov %rdx, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right8
mov $data45, %rsi
call __print
jmp __rightEnd8
 __right8:
jmp ._cond8_end
__rightEnd8:
jmp ._for0_end
._cond8_end:

mov $data46, %rsi
call __print
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right9
mov $data47, %rsi
call __print
jmp __rightEnd9
 __right9:
jmp ._cond9_end
__rightEnd9:
 call __clearVars
jmp ._for0
._cond9_end:

mov $data48, %rsi
call __print
._for0_end:

mov $data49, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName14, %rax 
mov $varName14, %rdi 
 call __set 
call __undefineVar
._cond0_end:

mov $data50, %rsi
call __print
mov $data51, %rsi
call __print
push $data52
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName13, %rax 
mov $varName13, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName11, %rax 
 mov $varName11, %rdi 
 call __set 
 call __getVar
mov (userData), %rdi 
 movb $'.', (%rdi) 
 jmp __goto
.next_prime_end:

mov $data53, %rsi
call __print
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
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
jmp .main_end*/
.main:

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
 mov $lenVarName23, %rax 
 mov $varName23, %rdi 
call __set

 mov $data154, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar

mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName50, %rax 
 mov $varName50, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax
 mov $stackType, %rdi
 call __set 
 call __defineVar

call __printHeap 
call __throughError


call __getVar 

mov (userData), %rax 
mov $2, %rbx 
call __openFile 
push %rax 
mov %rax, %r10
mov $varName23, %r8
call __writeToFile

call __toStr 
mov $buf2, %rsi 
call __print 
pop %rax 
mov %rax, %rdi 
call __closeFile 


/*push %rax 

mov %rax, %r10 
mov $4096, %rbx 
mov $varName23, %r8 
call __readFromFile

mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName23, %rax 
 mov $varName23, %rdi 
call __set
call __getVar 
mov (userData), %rsi 
call __print 

pop %rdi 
call __closeFile*/
#call __printHeap 
/*call __throughError
mov $data54, %rsi
call __print
mov $data55, %rsi
call __print
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
 mov $lenVarName21, %rax 
 mov $varName21, %rdi 
call __set

 mov $data56, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName21, %rax 
 mov $varName21, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName21, %rax 
mov $varName21, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName22, %rax 
 mov $varName22, %rdi 
call __set

 mov $data57, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName22, %rax 
 mov $varName22, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName22, %rax 
mov $varName22, %rdi 
 call __set 
call __undefineVar
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
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right10
mov $data58, %rsi
call __print
jmp __rightEnd10
 __right10:
jmp ._cond10_end
__rightEnd10:
mov $data59, %rsi
call __print
mov $data60, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName25, %rax 
 mov $varName25, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
._for1:

mov $data61, %rsi
call __print
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right11
mov $data62, %rsi
call __print
jmp __rightEnd11
 __right11:
jmp ._cond11_end
__rightEnd11:
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi 
call __set

 mov $data63, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNamePanic, %rax 
 mov $varNamePanic, %rdi 
call __set

 mov $data64, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right12
mov $data65, %rsi
call __print
jmp __rightEnd12
 __right12:
jmp ._cond12_end
__rightEnd12:
 mov $data66, %rax 
 mov %rax, (buf4)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi
 call __set 
 call __getVar 
 mov (userData), %rax 
 mov %rax, (buf3)
 mov (buf3), %rax 
 mov %rax, (buf)
 mov (buf4), %rax 
 mov %rax, (buf2)
 call __eqString 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (t0), %al 
 mov %al, (buf3)
mov (buf3), %al 
 mov %al, (buf) 

 call __not 
 mov (userData), %al 
 mov %al, (t1)
mov (userData), %al  
 cmp $'0', %al 
 jz __right13
mov $data67, %rsi
call __print
jmp __rightEnd13
 __right13:
jmp ._cond13_end
__rightEnd13:
jmp ._attempt0
._cond13_end:

mov $data68, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName23, %rax 
 mov $varName23, %rdi
 call __set 
 call __input
 mov $data69, %rax 
 mov %rax, (buf4)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi
 call __set 
 call __getVar 
 mov (userData), %rax 
 mov %rax, (buf3)
 mov (buf3), %rax 
 mov %rax, (buf)
 mov (buf4), %rax 
 mov %rax, (buf2)
 call __eqString 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (t0), %al 
 mov %al, (buf3)
mov (buf3), %al 
 mov %al, (buf) 

 call __not 
 mov (userData), %al 
 mov %al, (t1)
mov (userData), %al  
 cmp $'0', %al 
 jz __right14
mov $data70, %rsi
call __print
jmp __rightEnd14
 __right14:
jmp ._cond14_end
__rightEnd14:
jmp ._attempt0
._cond14_end:

mov $data71, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName26, %rax 
 mov $varName26, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName23, %rax 
 mov $varName23, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov (userData), %rdi
 call __set
 call __userToNumber
 call __toStr 
  
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName26, %rax 
 mov $varName26, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName26, %rax 
 mov $varName26, %rdi
 call __set 
 mov $t0, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set
mov $varName26, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set
mov $varName26, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName26, %rax 
mov $varName26, %rdi 
 call __set 
call __undefineVar
._attempt0:

mov $data72, %rsi
call __print
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNamePanic, %rax 
 mov $varNamePanic, %rdi 
call __set

 mov $data73, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
._cond12_end:

mov $data74, %rsi
call __print
mov $data75, %rsi
call __print
 mov $data76, %rax 
 mov %rax, (buf4)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi
 call __set 
 call __getVar 
 mov (userData), %rax 
 mov %rax, (buf3)
 mov (buf3), %rax 
 mov %rax, (buf)
 mov (buf4), %rax 
 mov %rax, (buf2)
 call __eqString 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (t0), %al 
 mov %al, (buf3)
mov (buf3), %al 
 mov %al, (buf) 

 call __not 
 mov (userData), %al 
 mov %al, (t1)
mov (userData), %al  
 cmp $'0', %al 
 jz __right15
mov $data77, %rsi
call __print
jmp __rightEnd15
 __right15:
jmp ._cond15_end
__rightEnd15:
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
 mov $lenVarName27, %rax 
 mov $varName27, %rdi 
call __set

 mov $data78, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName27, %rax 
 mov $varName27, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName27, %rax 
mov $varName27, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName28, %rax 
 mov $varName28, %rdi
 call __set
mov $varNameError, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName28, %rax 
 mov $varName28, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName28, %rax 
mov $varName28, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName29, %rax 
 mov $varName29, %rdi 
call __set

 mov $data79, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName29, %rax 
 mov $varName29, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName29, %rax 
mov $varName29, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName30, %rax 
 mov $varName30, %rdi 
call __set

 mov $data80, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName30, %rax 
 mov $varName30, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName30, %rax 
mov $varName30, %rdi 
 call __set 
call __undefineVar
jmp ._cond_exit0
._cond15_end:

mov $data81, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData82, %rax 
 mov $data82, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __less 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right16
mov $data83, %rsi
call __print
jmp __rightEnd16
 __right16:
jmp ._cond16_end
__rightEnd16:
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi 
call __set

 mov $data84, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $data85, %rsi
call __print
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
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenData86, %rax 
 mov $data86, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName0, %rax 
 mov $systemVarName0, %rdi 
 call __set
 mov $data86, %r8 
 mov $varNameError, %r9 
 xor %rax, %rax 
 mov $1, %rbx 
 call __userConcatinate
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData87, %rax 
 mov $data87, %rdi
 call __set
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenSystemVarName0, %rax 
 mov $systemVarName0, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName1, %rax 
 mov $systemVarName1, %rdi 
 call __set
 mov $systemVarName0, %r8 
 mov $data87, %r9 
 mov $1, %rax 
 xor %rbx, %rbx 
 call __userConcatinate
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName31, %rax 
 mov $varName31, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName31, %rax 
 mov $varName31, %rdi
 call __set 
mov $systemVarName1, %rax 
 mov %rax, (userData) 
mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName31, %rax 
 mov $varName31, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName31, %rax 
mov $varName31, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName32, %rax 
 mov $varName32, %rdi 
call __set

 mov $data88, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName32, %rax 
 mov $varName32, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName32, %rax 
mov $varName32, %rdi 
 call __set 
call __undefineVar
jmp ._cond_exit0
._cond16_end:

mov $data89, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData90, %rax 
 mov $data90, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __more 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right17
mov $data91, %rsi
call __print
jmp __rightEnd17
 __right17:
jmp ._cond17_end
__rightEnd17:
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi 
call __set

 mov $data92, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $data93, %rsi
call __print
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
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenData94, %rax 
 mov $data94, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName0, %rax 
 mov $systemVarName0, %rdi 
 call __set
 mov $data94, %r8 
 mov $varNameError, %r9 
 xor %rax, %rax 
 mov $1, %rbx 
 call __userConcatinate
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData95, %rax 
 mov $data95, %rdi
 call __set
mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov $lenSystemVarName0, %rax 
 mov $systemVarName0, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName1, %rax 
 mov $systemVarName1, %rdi 
 call __set
 mov $systemVarName0, %r8 
 mov $data95, %r9 
 mov $1, %rax 
 xor %rbx, %rbx 
 call __userConcatinate
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName33, %rax 
 mov $varName33, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName33, %rax 
 mov $varName33, %rdi
 call __set 
mov $systemVarName1, %rax 
 mov %rax, (userData) 
mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName33, %rax 
 mov $varName33, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName33, %rax 
mov $varName33, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName34, %rax 
 mov $varName34, %rdi 
call __set

 mov $data96, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName34, %rax 
 mov $varName34, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName34, %rax 
mov $varName34, %rdi 
 call __set 
call __undefineVar
jmp ._cond_exit0
._cond17_end:

mov $data97, %rsi
call __print
._cond_exit0:

mov $data98, %rsi
call __print
mov $data99, %rsi
call __print
 mov $data100, %rax 
 mov %rax, (buf4)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarNameError, %rax 
 mov $varNameError, %rdi
 call __set 
 call __getVar 
 mov (userData), %rax 
 mov %rax, (buf3)
 mov (buf3), %rax 
 mov %rax, (buf)
 mov (buf4), %rax 
 mov %rax, (buf2)
 call __eqString 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (t0), %al 
 mov %al, (buf3)
mov (buf3), %al 
 mov %al, (buf) 

 call __not 
 mov (userData), %al 
 mov %al, (t1)
mov (t1), %al 
 mov %al, (buf3)
mov (buf3), %al 
 mov %al, (buf) 

 call __not 
 mov (userData), %al 
 mov %al, (t2)
mov (userData), %al  
 cmp $'0', %al 
 jz __right18
mov $data101, %rsi
call __print
jmp __rightEnd18
 __right18:
jmp ._cond18_end
__rightEnd18:
mov $data102, %rsi
call __print
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName25, %rax 
 mov $varName25, %rdi 
call __set

 mov $data103, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
jmp ._undef_for1
._cond18_end:

mov $data104, %rsi
call __print
mov $data105, %rsi
call __print
._undef_for1:

mov $data106, %rsi
call __print
._cond11_end:

mov $data107, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName25, %rax 
 mov $varName25, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 mov (%rsi), %rdx 
 mov %rdx, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right19
mov $data108, %rsi
call __print
jmp __rightEnd19
 __right19:
jmp ._cond19_end
__rightEnd19:
jmp ._for1_end
._cond19_end:

mov $data109, %rsi
call __print
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right20
mov $data110, %rsi
call __print
jmp __rightEnd20
 __right20:
jmp ._cond20_end
__rightEnd20:
 call __clearVars
jmp ._for1
._cond20_end:

mov $data111, %rsi
call __print
._for1_end:

mov $data112, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName25, %rax 
mov $varName25, %rdi 
 call __set 
call __undefineVar
._cond10_end:

mov $data113, %rsi
call __print
mov $data114, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData115, %rax 
 mov $data115, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __eq 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right21
mov $data116, %rsi
call __print
jmp __rightEnd21
 __right21:
jmp ._cond21_end
__rightEnd21:
push $data117
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
mov $varName24, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName23, %rax 
mov $varName23, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName19, %rax 
 mov $varName19, %rdi 
 call __set 
 call __getVar
mov (userData), %rdi 
 movb $'.', (%rdi) 
 jmp __goto
._cond21_end:

mov $data118, %rsi
call __print
mov $data119, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName35, %rax 
 mov $varName35, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName35, %rax 
 mov $varName35, %rdi 
call __set

 mov $data120, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData121, %rax 
 mov $data121, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __eq 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right22
mov $data122, %rsi
call __print
jmp __rightEnd22
 __right22:
jmp ._cond22_end
__rightEnd22:
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
 mov $lenVarName36, %rax 
 mov $varName36, %rdi
 call __set
mov $varName35, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
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
 mov $lenVarName37, %rax 
 mov $varName37, %rdi
 call __set
mov $varName36, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName37, %rax 
 mov $varName37, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName37, %rax 
mov $varName37, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName36, %rax 
mov $varName36, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName38, %rax 
 mov $varName38, %rdi 
call __set

 mov $data123, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName38, %rax 
 mov $varName38, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName38, %rax 
mov $varName38, %rdi 
 call __set 
call __undefineVar
push $data124
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName35, %rax 
mov $varName35, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName23, %rax 
mov $varName23, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
mov $varName24, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName19, %rax 
 mov $varName19, %rdi 
 call __set 
 call __getVar
mov (userData), %rdi 
 movb $'.', (%rdi) 
 jmp __goto
._cond22_end:

mov $data125, %rsi
call __print
mov $data126, %rsi
call __print
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
 mov $lenVarName39, %rax 
 mov $varName39, %rdi
 call __set
mov $varName35, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
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
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData127, %rax 
 mov $data127, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName0, %rax 
 mov $systemVarName0, %rdi 
 call __set
 mov $varName39, %r8 
 mov $data127, %r9 
 mov $1, %rax 
 xor %rbx, %rbx 
 call __userConcatinate
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName40, %rax 
 mov $varName40, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName40, %rax 
 mov $varName40, %rdi
 call __set 
mov $systemVarName0, %rax 
 mov %rax, (userData) 
mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName40, %rax 
 mov $varName40, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName40, %rax 
mov $varName40, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName39, %rax 
mov $varName39, %rdi 
 call __set 
call __undefineVar
 movb $1, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right23
mov $data128, %rsi
call __print
jmp __rightEnd23
 __right23:
jmp ._cond23_end
__rightEnd23:
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName41, %rax 
 mov $varName41, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi 
call __set

 mov $data129, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName42, %rax 
 mov $varName42, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
._for2:

mov $data130, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData131, %rax 
 mov $data131, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __sub 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenT0, %rax 
 mov $t0, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __less 
 mov $lenT1, %rsi 
 mov $t1, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right24
mov $data132, %rsi
call __print
jmp __rightEnd24
 __right24:
jmp ._cond24_end
__rightEnd24:
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName43, %rax 
 mov $varName43, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov $varName35, %rdi
 call __set 
 call __getVar 
 push (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName11, %rax 
 mov $varName11, %rdi 
call __set

 mov $data133, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
jmp .next_prime
.next_prime_res0:

mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
mov $lenVarType, %rsi
mov $varType, %rdx
mov $lenStringType, %rax
mov $stringType, %rdi
call __set
mov $varName, %rcx
mov $varType, %rdx
call __defineVar
pop (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenPopVarName, %rax 
 mov $popVarName, %rdi
 call __set 
 xor %rax, %rax 
 call __setVar
mov $popVarName, %rax 
mov %rax, (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName12, %rax 
 mov $varName12, %rdi
 call __set 
 mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName43, %rax 
 mov $varName43, %rdi
 call __set
mov $varName12, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName35, %rax 
 mov $varName35, %rdi
 call __set
mov $varName43, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName43, %rax 
mov $varName43, %rdi 
 call __set 
call __undefineVar
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
 mov $lenVarName44, %rax 
 mov $varName44, %rdi
 call __set
mov $varName35, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
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
 mov $lenVarName45, %rax 
 mov $varName45, %rdi
 call __set
mov $varName44, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName45, %rax 
 mov $varName45, %rdi
 call __set
mov $varName44, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName45, %rax 
 mov $varName45, %rdi
 call __set
mov $varName44, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName45, %rax 
 mov $varName45, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName45, %rax 
mov $varName45, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName44, %rax 
mov $varName44, %rdi 
 call __set 
call __undefineVar
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData134, %rax 
 mov $data134, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __sub 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenT0, %rax 
 mov $t0, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __eq 
 mov $lenT1, %rsi 
 mov $t1, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (t1), %al 
 mov %al, (buf3)
mov (buf3), %al 
 mov %al, (buf) 

 call __not 
 mov (userData), %al 
 mov %al, (t2)
mov (userData), %al  
 cmp $'0', %al 
 jz __right25
mov $data135, %rsi
call __print
jmp __rightEnd25
 __right25:
jmp ._cond25_end
__rightEnd25:
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
 mov $lenVarName46, %rax 
 mov $varName46, %rdi 
call __set

 mov $data136, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName46, %rax 
 mov $varName46, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName46, %rax 
mov $varName46, %rdi 
 call __set 
call __undefineVar
._cond25_end:

mov $data137, %rsi
call __print
mov $data138, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData139, %rax 
 mov $data139, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __add 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi 
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi
 call __set 
 mov $t0, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
._undef_for2:

mov $data140, %rsi
call __print
._cond24_end:

mov $data141, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName42, %rax 
 mov $varName42, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 mov (%rsi), %rdx 
 mov %rdx, (userData)
mov (userData), %al  
 cmp $'0', %al 
 jz __right26
mov $data142, %rsi
call __print
jmp __rightEnd26
 __right26:
jmp ._cond26_end
__rightEnd26:
jmp ._for2_end
._cond26_end:

mov $data143, %rsi
call __print
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenData144, %rax 
 mov $data144, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
 mov $varName24, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __sub 
 mov $lenT0, %rsi 
 mov $t0, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenT0, %rax 
 mov $t0, %rdi
 call __set
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
 mov $varName41, %rdi
 call __set 
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi
 call __set 
mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf3, %rax 
 mov $buf3, %rdi
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set 
 xor %rax, %rax 

 call __less 
 mov $lenT1, %rsi 
 mov $t1, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set
mov (userData), %al  
 cmp $'0', %al 
 jz __right27
mov $data145, %rsi
call __print
jmp __rightEnd27
 __right27:
jmp ._cond27_end
__rightEnd27:
 call __clearVars
jmp ._for2
._cond27_end:

mov $data146, %rsi
call __print
._for2_end:

mov $data147, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName42, %rax 
mov $varName42, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName41, %rax 
mov $varName41, %rdi 
 call __set 
call __undefineVar
._cond23_end:

mov $data148, %rsi
call __print
mov $data149, %rsi
call __print
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
 mov $lenVarName47, %rax 
 mov $varName47, %rdi 
call __set

 mov $data150, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName47, %rax 
 mov $varName47, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName47, %rax 
mov $varName47, %rdi 
 call __set 
call __undefineVar
push $data151
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName24, %rax 
mov $varName24, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName23, %rax 
mov $varName23, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName35, %rax 
mov $varName35, %rdi 
 call __set 
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName19, %rax 
 mov $varName19, %rdi 
 call __set 
 call __getVar
mov (userData), %rdi 
 movb $'.', (%rdi) 
 jmp __goto
.main_end:

mov $data152, %rsi
call __print
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName48, %rax 
 mov $varName48, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName49, %rax 
 mov $varName49, %rdi 
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
 mov $varName, %rcx 
 mov $varType, %rdx  
 call __defineVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName19, %rax 
 mov $varName19, %rdi 
call __set

 mov $data153, %rax  
 mov %rax, (userData)
 xor %rax, %rax
call __setVar
jmp .main
.main_res0:

mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
mov $lenVarType, %rsi
mov $varType, %rdx
mov $lenStringType, %rax
mov $stringType, %rdi
call __set
mov $varName, %rcx
mov $varType, %rdx
call __defineVar
pop (userData)
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenPopVarName, %rax 
 mov $popVarName, %rdi
 call __set 
 xor %rax, %rax 
 call __setVar
mov $popVarName, %rax 
mov %rax, (userData)
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName20, %rax 
 mov $varName20, %rdi
 call __set 
 mov $1, %rax 
 call __setVar
mov $lenVarName, %rsi
mov $varName, %rdx
mov $lenPopVarName, %rax
mov $popVarName, %rdi
call __set
call __undefineVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName49, %rax 
 mov $varName49, %rdi
 call __set
mov $varName20, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName48, %rax 
 mov $varName48, %rdi
 call __set
mov $varName49, %rax
mov %rax, (userData)
 mov $1, %rax
call __setVar
mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName49, %rax 
mov $varName49, %rdi 
 call __set 
call __undefineVar*/
mov $60,  %rax
xor %rdi, %rdi
syscall
