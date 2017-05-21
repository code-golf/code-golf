#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mount.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <unistd.h>

void main (int argc, char *argv[]) {
    if (mount(NULL, "/", NULL, MS_PRIVATE|MS_REC, NULL) < 0) {
        perror("mount private");
        exit(EXIT_FAILURE);
    }

    if (mount("rootfs", "rootfs", "bind", MS_BIND|MS_REC, NULL) < 0) {
        perror("mount bind");
        exit(EXIT_FAILURE);
    }

    if (syscall(SYS_pivot_root, "rootfs", "rootfs/old-root") < 0) {
        perror("pivot_root");
        exit(EXIT_FAILURE);
    }

    if (chdir("/") < 0) {
        perror("chdir");
        exit(EXIT_FAILURE);
    }

    if (umount2("/old-root", MNT_DETACH) < 0) {
        perror("umount2");
        exit(EXIT_FAILURE);
    }

    if (sethostname(argv[0], strlen(argv[0])) < 0) {
        perror("sethostname");
        exit(EXIT_FAILURE);
    }

    char *newargv[] = {argv[0], NULL};

    execv(argv[1], newargv);
    perror("execv");
    exit(EXIT_FAILURE);
}
