#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* nial = "/usr/local/bin/nial", *code = "code.ndf";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen(code, "w")))
        ERR_AND_EXIT("fopen");

    if (fputs("setwidth 0;", fp))
        ERR_AND_EXIT("fputs");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    int nargc = argc + 3;
    char** nargv = malloc(nargc * sizeof(char*));
    nargv[0] = (char*) nial;
    nargv[1] = "-defs";
    nargv[2] = (char*) code;
    nargv[3] = "--";
    memcpy(&nargv[4], &argv[2], (argc - 2) * sizeof(char*));
    nargv[nargc - 1] = NULL;

    execv(nial, nargv);
    ERR_AND_EXIT("execv");
}
