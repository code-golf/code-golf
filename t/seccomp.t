use t;

is run( 'perl', 'syscall 248' )<stderr>, 'signal: bad system call (core dumped)',
    'add_key syscall kills process';

# i386 syscalls, see https://github.com/code-golf/code-golf/pull/2729
subtest 'i386 syscalls' => {
    is run( 'assembly', q:to/ASM/ )<exit_code>, 31, 'i386 keyctl (288) blocked';
        mov $288, %eax          # i386 keyctl (288)
        int $0x80

        mov $60,  %eax          # x64 exit
        xor %edi, %edi
        syscall
    ASM

    is run( 'assembly', q:to/ASM/ )<exit_code>, 0, 'x64 accept4 (288) allowed';
        mov $288, %rax          # x64 accept4 (288)
        xor %rdi, %rdi          # sockfd = 0
        xor %rsi, %rsi          # addr = NULL
        xor %rdx, %rdx          # addrlen = NULL
        xor %r10, %r10          # flags = 0
        syscall

        mov $60,  %eax          # x64 exit
        xor %edi, %edi
        syscall
    ASM

    is run( 'assembly', q:to/ASM/ )<stdout>, 'Hello', 'i386 write (4) allowed';
        mov $4,   %eax          # i386 write
        mov $1,   %ebx          # stdout
        mov $buf, %ecx
        mov $len, %edx
        int $0x80

        mov $1,   %eax          # i386 exit
        xor %ebx, %ebx
        int $0x80

    .data
        buf: .string "Hello\n"
        len = . - buf
    ASM
};

sub run($lang, $code) { post-solution( :$code :$lang )<runs>[0] }

done-testing;
