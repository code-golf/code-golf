use t;

like post-solution(:lang<assembly>)<Err>, rx[^
    'Signal: segmentation violation after line 1 (%rip was ' <xdigit> ** 16 ')'
 \n 'Register dump:'
 \n '    %rax = ' <xdigit> ** 16 '    %r8  = ' <xdigit> ** 16
 \n '    %rbx = ' <xdigit> ** 16 '    %r9  = ' <xdigit> ** 16
 \n '    %rcx = ' <xdigit> ** 16 '    %r10 = ' <xdigit> ** 16
 \n '    %rdx = ' <xdigit> ** 16 '    %r11 = ' <xdigit> ** 16
 \n '    %rsi = ' <xdigit> ** 16 '    %r12 = ' <xdigit> ** 16
 \n '    %rdi = ' <xdigit> ** 16 '    %r13 = ' <xdigit> ** 16
 \n '    %rsp = ' <xdigit> ** 16 '    %r14 = ' <xdigit> ** 16
 \n '    %rbp = ' <xdigit> ** 16 '    %r15 = ' <xdigit> ** 16
$];

done-testing;
