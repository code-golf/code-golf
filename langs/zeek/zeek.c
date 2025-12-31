#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* zeek = "/usr/local/bin/zeek", *code = "code.zeek";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(zeek, argv);
        ERR_AND_EXIT("execv");
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

    int zargc = argc + 2;
    char** zargv = malloc(zargc * sizeof(char*));
    zargv[0] = (char*) zeek;
    zargv[1] = "--";
    zargv[2] = (char*) code;
    memcpy(&zargv[3], &argv[2], (argc - 2) * sizeof(char*));
    zargv[zargc - 1] = NULL;

    execv(zeek, zargv);
    ERR_AND_EXIT("execv");
}
