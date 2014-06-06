.globl cao

cao:
    push %ebp
    mov %esp, %ebp

    push %eax
    push %ebx

    mov 1, %eax
    mov 2, %ebx
    testl %eax, %ebx

    pop %ebx
    pop %eax

    leave
    ret
