#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* io = "/usr/local/bin/io", *code = "code.io";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(io, argv);
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

    int iargc = argc + 1;
    char** iargv = malloc(iargc * sizeof(char*));
    iargv[0] = (char*) io;
    iargv[1] = (char*) code;
    memcpy(&iargv[2], &argv[2], (argc - 2) * sizeof(char*));
    iargv[iargc - 1] = NULL;

    execv(io, iargv);
    ERR_AND_EXIT("execv");
}
