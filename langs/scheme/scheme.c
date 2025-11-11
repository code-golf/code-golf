#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* scheme = "/usr/bin/scheme", *code = "code.ss";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(scheme, argv);
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

    int sargc = argc + 2;
    char** sargv = malloc(sargc * sizeof(char*));
    sargv[0] = (char*) scheme;
    sargv[1] = "--script";
    sargv[2] = (char*) code;
    memcpy(&sargv[3], &argv[2], (argc - 2) * sizeof(char*));
    sargv[sargc - 1] = NULL;

    execv(scheme, sargv);
    ERR_AND_EXIT("execv");
}
