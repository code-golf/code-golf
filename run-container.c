// 1152 bytes

long _call(int, void*, void*, void*, void*);

#define call1(n,a)       _call(n,(void*)(a),         0,         0,         0)
#define call2(n,a,b)     _call(n,(void*)(a),(void*)(b),         0,         0)
#define call3(n,a,b,c)   _call(n,(void*)(a),(void*)(b),(void*)(c),         0)
#define call4(n,a,b,c,d) _call(n,(void*)(a),(void*)(b),(void*)(c),(void*)(d))

#define EXIT_FAILURE 1

#define GID_nobody 99
#define UID_nobody 99

#define MNT_DETACH 2

#define MS_BIND    1 << 12
#define MS_REC     1 << 14
#define MS_PRIVATE 1 << 18

#define STDERR_FILENO 2

#define SYS_write       1
#define SYS_execve      59
#define SYS_exit        60
#define SYS_chdir       80
#define SYS_setuid      105
#define SYS_setgid      106
#define SYS_pivot_root  155
#define SYS_mount       165
#define SYS_umount2     166
#define SYS_sethostname 170

unsigned long strlen(const char* s) {
    unsigned long len = 0;
    while (*s++) ++len;
    return len;
}

void die(const char* const s) {
    call3(SYS_write, STDERR_FILENO, s, strlen(s));
    call3(SYS_write, STDERR_FILENO, "\n", 1);

    call1(SYS_exit, EXIT_FAILURE);
}

int main(int argc, char const* const argv[]) {
    if (call4(SYS_mount, 0, "/", 0, MS_PRIVATE|MS_REC) < 0)
        die("private");

    if (call4(SYS_mount, "rootfs", "rootfs", 0, MS_BIND|MS_REC) < 0)
        die("bind");

    if (call2(SYS_pivot_root, "rootfs", "rootfs/old-root") < 0)
        die("pivot");

    if (call1(SYS_chdir, "/") < 0)
        die("chdir");

    if (call2(SYS_umount2, "/old-root", MNT_DETACH) < 0)
        die("umount");

    if (call2(SYS_sethostname, argv[1], strlen(argv[1])) < 0)
        die("host");

    if (call1(SYS_setgid, GID_nobody) < 0)
        die("setgid");

    if (call1(SYS_setuid, UID_nobody) < 0)
        die("setuid");

    call2(SYS_execve, argv[0], argv + 1);
    die("exec");

    return 0;
}
