; nasm -f bin -o run-lang run-lang.asm && chmod +x run-lang

GID_nobody      equ 99

MNT_DETACH      equ 2

MS_BIND         equ 1 << 12
MS_REC          equ 1 << 14
MS_PRIVATE      equ 1 << 18

DEV_URANDOM     equ (1 << 8) | 9  ; makedev(1, 9)
DEV_RANDOM      equ (1 << 8) | 8  ; makedev(1, 8)

S_IFCHR         equ 0o60000

S_IRUSR         equ 0o0400
S_IRGRP         equ 0o0040
S_IROTH         equ 0o0004

SYS_write       equ 1
SYS_execve      equ 59
SYS_exit        equ 60
SYS_chdir       equ 80
SYS_chmod       equ 90
SYS_setuid      equ 105
SYS_setgid      equ 106
SYS_mknod       equ 133
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
proc        db "proc", 0
tmp         db "tmp", 0
tmpfs       db "tmpfs", 0
dev         db "dev", 0
urandom     db "dev/urandom", 0
random      db "dev/random", 0

start:
        ; mount / as private
            mov rax, SYS_mount
            ; rdi starts as 0
            mov rsi, slash
            mov r10, MS_PRIVATE|MS_REC
            ; r8 starts as 0
            syscall

            test eax, eax
            jnz exit

        ; bind mount rootfs
            mov rax, SYS_mount
            mov rdi, rootfs
            mov rsi, rdi
            ; edx starts as 0
            mov r10, MS_BIND|MS_REC
            syscall

            test eax, eax
            jnz exit

        ; pivot to rootfs
            mov rax, SYS_pivot_root
            ; rdi is still rootfs
            mov rsi, fullroot
            syscall

            test eax, eax
            jnz exit

        ; change directory to /
            mov rax, SYS_chdir
            mov rdi, slash
            syscall

            test eax, eax
            jnz exit

        ; unmount the old root
            mov rax, SYS_umount2
            mov rdi, oldroot
            mov rsi, MNT_DETACH
            syscall

            test eax, eax
            jnz exit

        ; mount /proc as proc
            mov rax, SYS_mount
            mov rdi, proc
            mov rsi, rdi
            mov rdx, rdi
            xor r10, r10
            ; r8 is still 0
            syscall

            test eax, eax
            jnz exit

        ; mount /tmp as tmpfs
            mov rax, SYS_mount
            mov rdi, tmp
            mov rsi, rdi
            mov edx, tmpfs
            ; r10 is still 0
            ; r8  is still 0
            syscall

            test eax, eax
            jnz exit

        ; mount /dev as tmpfs
            mov rax, SYS_mount
            mov rdi, dev
            mov rsi, rdi
            mov edx, tmpfs
            ; r10 is still 0
            ; r8  is still 0
            syscall

            test eax, eax
            jnz exit

        ; mount /dev/urandom as block character device
            mov rax, SYS_mknod
            mov rdi, urandom
            mov rsi, S_IFCHR
            mov edx, DEV_URANDOM
            syscall

            test eax, eax
            jnz exit

        ; make /dev/urandom universally readable
            mov rax, SYS_chmod
            mov rdi, urandom
            mov rsi, S_IRUSR|S_IRGRP|S_IROTH
            syscall

            test eax, eax
            jnz exit

        ; mount /dev/random as block character device
            mov rax, SYS_mknod
            mov rdi, random
            mov rsi, S_IFCHR
            mov edx, DEV_RANDOM
            syscall

            test eax, eax
            jnz exit

        ; make /dev/random universally readable
            mov rax, SYS_chmod
            mov rdi, random
            mov rsi, S_IRUSR|S_IRGRP|S_IROTH
            syscall

            test eax, eax
            jnz exit

        ; set the hostname
            mov rax, SYS_sethostname
            mov rdi, host
            mov rsi, hostsize
            syscall

            test eax, eax
            jnz exit

        ; set the group
            mov rax, SYS_setgid
            mov rdi, GID_nobody
            syscall

            test eax, eax
            jnz exit

        ; set the user
            mov rax, SYS_setuid
            ; rdi is still GID_nobody which is identical to UID_nobody
            syscall

            test eax, eax
            jnz exit

        ; syscall(SYS_execve, argv[0], argv, 0);
            mov rax, SYS_execve
            lea rsi, [rsp + 8] ; argv
            mov rdi, [rsi]     ; argv[0]
            xor edx, edx
            syscall

exit:
            mov rax, SYS_exit
            mov rdi, 1
            syscall

filesize    equ  $ - $$
