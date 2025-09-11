#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* zig = "/usr/local/zig", *code = "code.zig";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(zig, zig, "version", NULL);
        ERR_AND_EXIT("execl");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen(code, "w")))
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    int zargc = argc + 7;
    char** zargv = malloc(zargc * sizeof(char*));
    zargv[0] = (char*) zig;
    zargv[1] = "run";
    zargv[2] = "-freference-trace";
    zargv[3] = "-fstrip";
    zargv[4] = "--color";
    zargv[5] = "on";
    zargv[6] = (char*) code;
    zargv[7] = "--";
    memcpy(&zargv[8], &argv[2], (argc - 2) * sizeof(char*));
    zargv[zargc - 1] = NULL;

    execv(zig, zargv);
    ERR_AND_EXIT("execv");
}
