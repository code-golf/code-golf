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

%include "/tmp/code.asm"

filesize    equ  $ - $$