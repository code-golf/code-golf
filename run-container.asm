GID_nobody      equ 99
UID_nobody      equ 99

MNT_DETACH      equ 2

MS_BIND         equ 1 << 12
MS_REC          equ 1 << 14
MS_PRIVATE      equ 1 << 18

SYS_write       equ 1
SYS_execve      equ 59
SYS_exit        equ 60
SYS_chdir       equ 80
SYS_setuid      equ 105
SYS_setgid      equ 106
SYS_pivot_root  equ 155
SYS_mount       equ 165
SYS_umount2     equ 166
SYS_sethostname equ 170

BITS 64
            org 0x400000

ehdr:                         ; Elf64_Ehdr
            db 0x7f, "ELF", 2 ;  e_ident
    times 2 db 1              ;  e_ident cont.
    times 9 db 0              ;  e_ident cont.
            dw 2              ;  e_type
            dw 0x3e           ;  e_machine
            dd 1              ;  e_version
            dq start          ;  e_entry
            dq phdr - $$      ;  e_phoff
            dq 0              ;  e_shoff
            dd 0              ;  e_flags
            dw ehdrsize       ;  e_ehsize
            dw phdrsize       ;  e_phentsize
            dw 1              ;  e_phnum
    times 3 dw 0              ;  e_shentsize, e_shnum, e_shstrndx

ehdrsize    equ $ - ehdr

phdr:                         ; Elf64_Phdr
            dd 1              ;  p_type
            dd 5              ;  p_flags
            dq 0              ;  p_offset
            dq $$             ;  p_vaddr
            dq $$             ;  p_paddr
    times 2 dq filesize       ;  p_filesz, p_memsz
            dq 0x1000         ;  p_align

phdrsize    equ $ - phdr

host        db  "code-golf"
hostsize    equ $ - host

fullroot    db "rootfs/old-root", 0
oldroot     db "old-root", 0
rootfs      db "rootfs", 0
slash       db "/", 0
slash_proc  db "/proc", 0
slash_tmp   db "/tmp", 0
proc        db "proc", 0
tmpfs       db "tmpfs", 0
stage       db "ABCDEFGHIJK"

start:
        ; [A] mount / as private
            mov rax, SYS_mount
            ; rdi starts as 0
            mov rsi, slash
            mov r10, MS_PRIVATE|MS_REC
            ; r8 starts as 0
            syscall
            mov r12, stage

            test eax, eax
            jnz exit

        ; [B] bind mount rootfs
            mov rax, SYS_mount
            mov rdi, rootfs
            mov rsi, rdi
            xor edx, edx
            mov r10, MS_BIND|MS_REC
            syscall
            inc r12

            test eax, eax
            jnz exit

        ; [C] pivot to rootfs
            mov rax, SYS_pivot_root
            ; rdi is still rootfs
            mov rsi, fullroot
            syscall
            inc r12

            test eax, eax
            jnz exit

        ; [D] change directory to /
            mov rax, SYS_chdir
            mov rdi, slash
            syscall
            inc r12

            test eax, eax
            jnz exit

        ; [E] unmount the old root
            mov rax, SYS_umount2
            mov rdi, oldroot
            mov rsi, MNT_DETACH
            syscall
            inc r12

            test eax, eax
            jnz exit

        ; [F] mount /proc as proc
            mov rax, SYS_mount
            mov rdi, slash_proc
            mov rsi, rdi
            mov edx, proc
            xor r10, r10
            xor r8, r8
            syscall
            mov r12, stage

            test eax, eax
            jnz exit

        ; [G] mount /tmp as tmpfs
            mov rax, SYS_mount
            mov rdi, slash_tmp
            mov rsi, rdi
            mov edx, tmpfs
            xor r10, r10
            xor r8, r8
            syscall
            mov r12, stage

            test eax, eax
            jnz exit

        ; [H] set the hostname
            mov rax, SYS_sethostname
            mov rdi, host
            mov rsi, hostsize
            syscall
            inc r12

            test eax, eax
            jnz exit

        ; [I] set the group
            mov rax, SYS_setgid
            mov rdi, GID_nobody
            syscall
            inc r12

            test eax, eax
            jnz exit

        ; [J] set the user
            mov rax, SYS_setuid
            mov rdi, UID_nobody
            syscall
            inc r12

            test eax, eax
            jnz exit

        ; [K] syscall(SYS_execve, argv[0], argv, 0);
            mov rax, SYS_execve
            lea rsi, [rsp + 8] ; argv
            mov rdi, [rsi]     ; argv[0]
            xor edx, edx
            syscall
            inc r12

exit:
            mov rax, SYS_write
            mov rdi, 1
            mov rsi, r12
            mov rdx, 1
            syscall

            mov rax, SYS_exit
            mov rdi, 1
            syscall

filesize    equ  $ - $$
