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

subtest 'io_uring' => {
    is run( 'perl', 'syscall 425; exit $!' )<exit_code>, 38,
        'io_uring_setup syscall returns ENOSY';

    is run( 'perl', 'syscall 426, -1; exit $!' )<exit_code>, 38,
        'io_uring_enter syscall returns ENOSY';

    is run( 'perl', 'syscall 427, -1; exit $!' )<exit_code>, 38,
        'io_uring_register syscall returns ENOSY';
};

subtest 'vsock' => {
    is run( 'python', 'import socket; socket.socket(socket.AF_UNIX)' )<exit_code>,
        0, 'socket AF_UNIX is allowed';

    is run( 'python', 'import socket; socket.socket(socket.AF_VSOCK)' )<stderr>,
        'signal: bad system call (core dumped)', 'socket AF_VSOCK is killed';

    is run( 'python', 'import socket; socket.socketpair(socket.AF_UNIX)' )<exit_code>,
        0, 'socketpair AF_UNIX is allowed';

    is run( 'python', 'import socket; socket.socketpair(socket.AF_VSOCK)' )<stderr>,
        'signal: bad system call (core dumped)', 'socketpair AF_VSOCK is killed';
};

sub run($lang, $code) { post-solution( :$code :$lang )<runs>[0] }

done-testing;
