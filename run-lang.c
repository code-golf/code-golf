//go:build none

#define _GNU_SOURCE
#include <linux/filter.h>
#include <linux/sched.h>
#include <linux/seccomp.h>
#include <sched.h>
#include <spawn.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/mount.h>
#include <sys/prctl.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <sys/sysmacros.h>
#include <sys/wait.h>
#include <unistd.h>

#define NOBODY 65534

#define ALLOW(name) \
    BPF_JUMP(BPF_JMP+BPF_JEQ+BPF_K, __NR_##name, 0, 1), \
    BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_ALLOW)

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)
#define STR_WITH_LEN(str) str, sizeof(str) - 1

int main(__attribute__((unused)) int argc, char *argv[]) {
    struct clone_args args = {
        .exit_signal = SIGCHLD,
        .flags       = CLONE_NEWIPC | CLONE_NEWNET | CLONE_NEWNS
                     | CLONE_NEWPID | CLONE_NEWUTS,
    };

    pid_t pid = syscall(__NR_clone3, &args, sizeof(struct clone_args));
    if (pid < 0)
        ERR_AND_EXIT("clone");

    if (pid) {
        int status;
        do {
            if (waitpid(pid, &status, 0) < 0)
                ERR_AND_EXIT("waitpid");
        } while (!WIFEXITED(status) && !WIFSIGNALED(status));

        puts("parent end");

        return WEXITSTATUS(status);
    }

    if (mount(NULL, "/", NULL, MS_PRIVATE|MS_REC, NULL) < 0)
        ERR_AND_EXIT("mount private");

    if (mount("rootfs", "rootfs", "bind", MS_BIND|MS_REC, NULL) < 0)
        ERR_AND_EXIT("mount bind");

    if (syscall(SYS_pivot_root, "rootfs", "rootfs") < 0)
        ERR_AND_EXIT("pivot_root");

    if (chdir("/") < 0)
        ERR_AND_EXIT("chdir");

    if (umount2("/", MNT_DETACH) < 0)
        ERR_AND_EXIT("umount2");

    if (mount("tmpfs", "/dev", "tmpfs", MS_NOSUID, NULL) < 0)
        ERR_AND_EXIT("mount dev");

    if (mknod("/dev/null", S_IFCHR|0666, makedev(1, 3)) < 0)
        ERR_AND_EXIT("mknod");

    // FIXME This shouldn't be needed, 0666 should suffice, but without it Zig
    //       fails with permission denied when opening /dev/null as O_RDWR.
    if (chown("/dev/null", NOBODY, NOBODY) < 0)
        ERR_AND_EXIT("chown /dev/null");

    if (mknod("/dev/urandom", S_IFCHR|0444, makedev(1, 9)) < 0)
        ERR_AND_EXIT("mknod");

    if (mount("proc", "/proc", "proc", MS_NODEV|MS_NOEXEC|MS_NOSUID|MS_RDONLY, NULL) < 0)
        ERR_AND_EXIT("mount proc");

    if (mount("tmpfs", "/tmp", "tmpfs", MS_NODEV|MS_NOSUID, NULL) < 0)
        ERR_AND_EXIT("mount tmp");

    if (sethostname(STR_WITH_LEN("code-golf")) < 0)
        ERR_AND_EXIT("sethostname");

    // Allow /proc/self/fd/0 to be read by the lang after we change user.
    if (chown("/proc/self/fd/0", NOBODY, NOBODY) < 0)
        ERR_AND_EXIT("chown /proc/self/fd/0");

    if (setgid(NOBODY) < 0)
        ERR_AND_EXIT("setgid");

    if (setuid(NOBODY) < 0)
        ERR_AND_EXIT("setuid");

    // sudo journalctl -f _AUDIT_TYPE_NAME=SECCOMP
    // ... SECCOMP ... syscall=xxx ...
    struct sock_filter filter[] = {
        BPF_STMT(BPF_LD+BPF_W+BPF_ABS, offsetof(struct seccomp_data, nr)),

        // FIXME Julia attempts this wildly high syscall :-S
        // SECCOMP auid=x uid=65534 gid=65534 ses=x pid=x comm="julia" exe="/usr/bin/julia" sig=31 arch=c000003e syscall=1008 compat=0 ip=x code=0x0
        #define __NR_julia 1008
        ALLOW(julia),

        /*************\
        | File System |
        \*************/

        // File Operations
        ALLOW(close),             // 3
        ALLOW(close_range),       // 436
        ALLOW(creat),             // 85
        ALLOW(fallocate),         // 285
        ALLOW(ftruncate),         // 77
        ALLOW(memfd_create),      // 319
        ALLOW(mknod),             // 133
        ALLOW(mknodat),           // 259
        ALLOW(name_to_handle_at), // 303
        ALLOW(open),              // 2
        ALLOW(open_by_handle_at), // 304
        ALLOW(openat),            // 257
        ALLOW(openat2),           // 437
        ALLOW(rename),            // 82
        ALLOW(renameat),          // 264
        ALLOW(renameat2),         // 316
        ALLOW(truncate),          // 76
        ALLOW(userfaultfd),       // 323

        // Directory Operations
        ALLOW(chdir),          // 80
        ALLOW(chroot),         // 161
        ALLOW(fchdir),         // 81
        ALLOW(getcwd),         // 79
        ALLOW(getdents64),     // 217
        ALLOW(getdents),       // 78
        ALLOW(lookup_dcookie), // 212
        ALLOW(mkdir),          // 83
        ALLOW(mkdirat),        // 258
        ALLOW(rmdir),          // 84

        // Link Operations
        ALLOW(link),       // 86
        ALLOW(linkat),     // 265
        ALLOW(readlink),   // 89
        ALLOW(readlinkat), // 267
        ALLOW(symlink),    // 88
        ALLOW(symlinkat),  // 266
        ALLOW(unlink),     // 87
        ALLOW(unlinkat),   // 263

        // Basic File Attributes
        ALLOW(access),     // 21
        ALLOW(chmod),      // 90
        ALLOW(chown),      // 92
        ALLOW(faccessat),  // 269
        ALLOW(faccessat2), // 439
        ALLOW(fchmod),     // 91
        ALLOW(fchmodat),   // 268
        ALLOW(fchown),     // 93
        ALLOW(fchownat),   // 260
        ALLOW(fstat),      // 5
        ALLOW(futimesat),  // 261
        ALLOW(lchown),     // 94
        ALLOW(lstat),      // 6
        ALLOW(newfstatat), // 262
        ALLOW(stat),       // 4
        ALLOW(statx),      // 332
        ALLOW(umask),      // 95
        ALLOW(utime),      // 132
        ALLOW(utimensat),  // 280
        ALLOW(utimes),     // 235

        // Extended File Attributes
        ALLOW(fgetxattr),    // 193
        ALLOW(flistxattr),   // 196
        ALLOW(fremovexattr), // 199
        ALLOW(fsetxattr),    // 190
        ALLOW(getxattr),     // 191
        ALLOW(lgetxattr),    // 192
        ALLOW(listxattr),    // 194
        ALLOW(llistxattr),   // 195
        ALLOW(lremovexattr), // 198
        ALLOW(lsetxattr),    // 189
        ALLOW(removexattr),  // 197
        ALLOW(setxattr),     // 188

        // File Descriptor Manipulations
        ALLOW(dup2),  // 33
        ALLOW(dup),   // 32
        ALLOW(dup3),  // 292
        ALLOW(fcntl), // 72
        ALLOW(flock), // 73
        ALLOW(ioctl), // 16

        // Read/Write
        ALLOW(copy_file_range), // 326
        ALLOW(lseek),           // 8
        ALLOW(pread64),         // 17
        ALLOW(preadv2),         // 327
        ALLOW(preadv),          // 295
        ALLOW(pwrite64),        // 18
        ALLOW(pwritev2),        // 328
        ALLOW(pwritev),         // 296
        ALLOW(read),            // 0
        ALLOW(readv),           // 19
        ALLOW(sendfile),        // 40
        ALLOW(write),           // 1
        ALLOW(writev),          // 20

        // Synchronized I/O
        ALLOW(fdatasync),       // 75
        ALLOW(fsync),           // 74
        ALLOW(msync),           // 26
        ALLOW(sync),            // 162
        ALLOW(sync_file_range), // 277
        ALLOW(syncfs),          // 306

        // Asynchronous I/O
        ALLOW(io_pgetevents),        // 333
        ALLOW(io_cancel),            // 210
        ALLOW(io_destroy),           // 207
        ALLOW(io_getevents),         // 208
        ALLOW(io_setup),             // 206
        ALLOW(io_submit),            // 209
        // ALLOW(io_uring_enter),    // 426
        // ALLOW(io_uring_register), // 427
        // ALLOW(io_uring_setup),    // 425

        // Multiplexed I/O
        ALLOW(epoll_create),  // 213
        ALLOW(epoll_create1), // 291
        ALLOW(epoll_ctl),     // 233
        ALLOW(epoll_pwait),   // 281
        ALLOW(epoll_wait),    // 232
        ALLOW(poll),          // 7
        ALLOW(ppoll),         // 271
        ALLOW(pselect6),      // 270
        ALLOW(select),        // 23

        // Monitoring File Events
        // ALLOW(fanotify_init),     // 300
        // ALLOW(fanotify_mark),     // 301
        // ALLOW(inotify_add_watch), // 254
        // ALLOW(inotify_init1),     // 294
        // ALLOW(inotify_init),      // 253
        // ALLOW(inotify_rm_watch),  // 255

        // Miscellaneous
        ALLOW(fadvise64), // 221
        ALLOW(getrandom), // 318
        ALLOW(readahead), // 187

        /*********\
        | Network |
        \*********/

        // Socket Operations
        ALLOW(accept),      // 43
        ALLOW(accept4),     // 288
        ALLOW(bind),        // 49
        ALLOW(connect),     // 42
        ALLOW(getpeername), // 52
        ALLOW(getsockname), // 51
        ALLOW(getsockopt),  // 55
        ALLOW(listen),      // 50
        ALLOW(setsockopt),  // 54
        ALLOW(shutdown),    // 48
        ALLOW(socket),      // 41
        ALLOW(socketpair),  // 53

        // Send/Receive
        ALLOW(recvfrom), // 45
        ALLOW(recvmmsg), // 299
        ALLOW(recvmsg),  // 47
        ALLOW(sendmmsg), // 307
        ALLOW(sendmsg),  // 46
        ALLOW(sendto),   // 44

        // Naming
        // ALLOW(setdomainname), // 171
        // ALLOW(sethostname),   // 170

        // Packet Filtering
        // ALLOW(bpf), // 321

        /******\
        | Time |
        \******/

        // Current Time of Day
        ALLOW(gettimeofday),    // 96
        // ALLOW(settimeofday), // 164
        ALLOW(time),            // 201

        // POSIX Clocks
        ALLOW(clock_adjtime),   // 305
        ALLOW(clock_getres),    // 229
        ALLOW(clock_gettime),   // 228
        ALLOW(clock_nanosleep), // 230
        ALLOW(clock_settime),   // 227

        // Clocks Based Timers
        ALLOW(timer_create),     // 222
        ALLOW(timer_delete),     // 226
        ALLOW(timer_getoverrun), // 225
        ALLOW(timer_gettime),    // 224
        ALLOW(timer_settime),    // 223

        // Timers
        ALLOW(alarm),     // 37
        ALLOW(getitimer), // 36
        ALLOW(setitimer), // 38

        // File Descriptor Based Timers
        ALLOW(timerfd_create),  // 283
        ALLOW(timerfd_gettime), // 287
        ALLOW(timerfd_settime), // 286

        // Miscellaneous
        // ALLOW(adjtimex), // 159
        ALLOW(nanosleep),   // 35
        ALLOW(times),       // 100

        /***********\
        | Processes |
        \***********/

        // Creation and Termination
        ALLOW(clone),      // 56
        ALLOW(clone3),     // 435
        ALLOW(execve),     // 59
        ALLOW(execveat),   // 322
        ALLOW(exit),       // 60
        ALLOW(exit_group), // 231
        ALLOW(fork),       // 57
        ALLOW(vfork),      // 58
        ALLOW(wait4),      // 61
        ALLOW(waitid),     // 247

        // Process ID
        ALLOW(getpid),      // 39
        ALLOW(getppid),     // 110
        ALLOW(gettid),      // 186
        ALLOW(pidfd_getfd), // 438
        ALLOW(pidfd_open),  // 434

        // Session ID
        ALLOW(getsid), // 124
        ALLOW(setsid), // 112

        // Process Group ID
        ALLOW(getpgid), // 121
        ALLOW(getpgrp), // 111
        ALLOW(setpgid), // 109

        // Users and Groups
        ALLOW(getegid),   // 108
        ALLOW(geteuid),   // 107
        ALLOW(getgid),    // 104
        ALLOW(getgroups), // 115
        ALLOW(getresgid), // 120
        ALLOW(getresuid), // 118
        ALLOW(getuid),    // 102
        ALLOW(setfsgid),  // 123
        ALLOW(setfsuid),  // 122
        ALLOW(setgid),    // 106
        ALLOW(setgroups), // 116
        ALLOW(setregid),  // 114
        ALLOW(setresgid), // 119
        ALLOW(setresuid), // 117
        ALLOW(setreuid),  // 113
        ALLOW(setuid),    // 105

        // Namespaces
        // ALLOW(setns), // 308

        // Resource Limits
        ALLOW(getrlimit), // 97
        ALLOW(getrusage), // 98
        ALLOW(prlimit64), // 302
        ALLOW(setrlimit), // 160

        // Process Scheduling
        ALLOW(getpriority),            // 140
        ALLOW(ioprio_get),             // 252
        ALLOW(ioprio_set),             // 251
        ALLOW(sched_getaffinity),      // 204
        ALLOW(sched_getattr),          // 315
        ALLOW(sched_getparam),         // 143
        ALLOW(sched_get_priority_max), // 146
        ALLOW(sched_get_priority_min), // 147
        ALLOW(sched_getscheduler),     // 145
        ALLOW(sched_rr_get_interval),  // 148
        ALLOW(sched_setaffinity),      // 203
        ALLOW(sched_setattr),          // 314
        ALLOW(sched_setparam),         // 142
        ALLOW(sched_setscheduler),     // 144
        ALLOW(sched_yield),            // 24
        ALLOW(setpriority),            // 141

        // Virtual Memory
        ALLOW(brk),           // 12
        ALLOW(madvise),       // 28
        ALLOW(membarrier),    // 324
        ALLOW(mincore),       // 27
        ALLOW(mlock),         // 149
        ALLOW(mlock2),        // 325
        ALLOW(mlockall),      // 151
        ALLOW(mmap),          // 9
        ALLOW(modify_ldt),    // 154
        ALLOW(mprotect),      // 10
        ALLOW(mremap),        // 25
        ALLOW(munlock),       // 150
        ALLOW(munlockall),    // 152
        ALLOW(munmap),        // 11
        ALLOW(pkey_alloc),    // 330
        ALLOW(pkey_free),     // 331
        ALLOW(pkey_mprotect), // 329

        // Threads
        ALLOW(arch_prctl),      // 158
        ALLOW(capget),          // 125
        ALLOW(capset),          // 126
        ALLOW(get_thread_area), // 211
        ALLOW(set_thread_area), // 205
        ALLOW(set_tid_address), // 218

        // Miscellaneous
        // ALLOW(kcmp),              // 312
        ALLOW(prctl),                // 157
        // ALLOW(process_vm_readv),  // 310
        // ALLOW(process_vm_writev), // 311
        ALLOW(ptrace),               // 101
        // ALLOW(seccomp),           // 317
        // ALLOW(unshare),           // 272
        // ALLOW(uselib),            // 134

        /*********\
        | Signals |
        \*********/

        // Standard Signals
        ALLOW(kill),   // 62
        ALLOW(pause),  // 34
        ALLOW(tgkill), // 234
        ALLOW(tkill),  // 200

        // Real-time Signals
        ALLOW(rt_sigaction),      // 13
        ALLOW(rt_sigpending),     // 127
        ALLOW(rt_sigprocmask),    // 14
        ALLOW(rt_sigqueueinfo),   // 129
        ALLOW(rt_sigreturn),      // 15
        ALLOW(rt_sigsuspend),     // 130
        ALLOW(rt_sigtimedwait),   // 128
        ALLOW(rt_tgsigqueueinfo), // 297
        ALLOW(sigaltstack),       // 131

        // File Descriptor Based Signals
        ALLOW(eventfd),           // 284
        ALLOW(eventfd2),          // 290
        ALLOW(pidfd_send_signal), // 424
        ALLOW(signalfd),          // 282
        ALLOW(signalfd4),         // 289

        // Miscellaneous
        ALLOW(restart_syscall), // 219

        /*****\
        | IPC |
        \*****/

        // Pipe
        ALLOW(pipe),     // 22
        ALLOW(pipe2),    // 293
        ALLOW(splice),   // 275
        ALLOW(tee),      // 276
        ALLOW(vmsplice), // 278

        // Shared Memory
        ALLOW(shmat),  // 30
        ALLOW(shmctl), // 31
        ALLOW(shmdt),  // 67
        ALLOW(shmget), // 29

        // Semaphores
        // ALLOW(semctl),     // 66
        // ALLOW(semget),     // 64
        // ALLOW(semop),      // 65
        // ALLOW(semtimedop), // 220

        // Futexes
        ALLOW(futex),           // 202
        ALLOW(get_robust_list), // 274
        ALLOW(set_robust_list), // 273

        // System V Message Queue
        // ALLOW(msgctl), // 71
        // ALLOW(msgget), // 68
        // ALLOW(msgrcv), // 70
        // ALLOW(msgsnd), // 69

        // POSIX Message Queue
        // ALLOW(mq_getsetattr),   // 245
        // ALLOW(mq_notify),       // 244
        // ALLOW(mq_open),         // 240
        // ALLOW(mq_timedreceive), // 243
        // ALLOW(mq_timedsend),    // 242
        // ALLOW(mq_unlink),       // 241

        /******\
        | NUMA |
        \******/

        // CPU Node
        // ALLOW(getcpu), // 309

        // Memory Node
        ALLOW(get_mempolicy), // 239
        ALLOW(mbind),         // 237
        ALLOW(migrate_pages), // 256
        ALLOW(move_pages),    // 279
        ALLOW(set_mempolicy), // 238

        /****************\
        | Key Management |
        \****************/

        // ALLOW(add_key),     // 248
        // ALLOW(keyctl),      // 250
        // ALLOW(request_key), // 249

        /*************\
        | System-Wide |
        \*************/

        // Loadable Modules
        // ALLOW(create_module),   // 174
        // ALLOW(delete_module),   // 176
        // ALLOW(finit_module),    // 313
        // ALLOW(get_kernel_syms), // 177
        // ALLOW(init_module),     // 175
        // ALLOW(query_module),    // 178

        // Accounting and Quota
        // ALLOW(acct),     // 163
        // ALLOW(quotactl), // 179

        // Filesystem (privileged)
        // ALLOW(fsconfig),   // 431
        // ALLOW(fsmount),    // 432
        // ALLOW(fsopen),     // 430
        // ALLOW(fspick),     // 433
        // ALLOW(mount),      // 165
        // ALLOW(move_mount), // 429
        // ALLOW(nfsservctl), // 180
        // ALLOW(open_tree),  // 428
        // ALLOW(pivot_root), // 155
        // ALLOW(swapoff),    // 168
        // ALLOW(swapon),     // 167
        // ALLOW(umount2),    // 166

        // Filesystem (unprivileged)
        ALLOW(fstatfs), // 138
        ALLOW(statfs),  // 137
        ALLOW(sysfs),   // 139
        ALLOW(ustat),   // 136

        // Miscellaneous (privileged)
        ALLOW(ioperm),             // 173 (FreeBASIC)
        // ALLOW(iopl),            // 172
        // ALLOW(kexec_file_load), // 320
        // ALLOW(kexec_load),      // 246
        // ALLOW(perf_event_open), // 298
        // ALLOW(personality),     // 135
        // ALLOW(reboot),          // 169
        // ALLOW(_sysctl),         // 156
        // ALLOW(syslog),          // 103
        // ALLOW(vhangup),         // 153

        // Miscellaneous (unprivileged)
        ALLOW(rseq),    // 334
        ALLOW(sysinfo), // 99
        ALLOW(uname),   // 63

        /************\
        | Deprecated |
        \************/

        // ALLOW(remap_file_pages), // 216

        /***************\
        | Unimplemented |
        \***************/

        // ALLOW(afs_syscall),     // 183
        // ALLOW(create_module),   // 174
        // ALLOW(epoll_ctl_old),   // 214
        // ALLOW(epoll_wait_old),  // 215
        // ALLOW(get_kernel_syms), // 177
        // ALLOW(getpmsg),         // 181
        // ALLOW(nfsservctl),      // 180
        // ALLOW(putpmsg),         // 182
        // ALLOW(query_module),    // 178
        // ALLOW(security),        // 185
        // ALLOW(tuxcall),         // 184
        // ALLOW(vserver),         // 236

        BPF_STMT(BPF_RET+BPF_K, SECCOMP_RET_KILL),
    };

    struct sock_fprog fprog = {
        (unsigned short) (sizeof(filter) / sizeof(filter[0])),
        filter,
    };

    if (prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0) < 0)
        ERR_AND_EXIT("prctl(NO_NEW_PRIVS)");

    if (prctl(PR_SET_SECCOMP, SECCOMP_MODE_FILTER, &fprog) < 0)
        ERR_AND_EXIT("prctl(SECCOMP)");

    execvp(argv[0], argv);
    ERR_AND_EXIT("execvp");
}
