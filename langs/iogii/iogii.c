#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* iogii = "/usr/local/bin/iogii", *code = "code.iog";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(iogii, argv);
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

    int iargc = argc + 2;
    char** iargv = malloc(iargc * sizeof(char*));
    iargv[0] = (char*) iogii;
    iargv[1] = (char*) code;
    iargv[2] = "--";
    memcpy(&iargv[3], &argv[2], (argc - 2) * sizeof(char*));
    iargv[iargc - 1] = NULL;

    execv(iogii, iargv);
    ERR_AND_EXIT("execv");
}
