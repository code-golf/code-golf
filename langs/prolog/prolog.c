#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* prolog = "/usr/bin/swipl", *code = "/tmp/code.pl";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(prolog, argv);
        ERR_AND_EXIT("execv");
    }

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

    int pargc = argc;
    char** pargv = malloc(pargc * sizeof(char*));
    pargv[0] = (char*) prolog;
    memcpy(&pargv[1], &argv[2], (argc - 2) * sizeof(char*));
    pargv[pargc - 1] = NULL;

    execv(prolog, argv);
    ERR_AND_EXIT("execv");
}
