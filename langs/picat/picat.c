#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* picat = "/usr/local/bin/picat", *code = "code.pi";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(picat, argv);
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

    int pargc = argc + 1;
    char** pargv = malloc(pargc * sizeof(char*));
    pargv[0] = (char*) picat;
    pargv[1] = (char*) code;
    memcpy(&pargv[2], &argv[2], (argc - 2) * sizeof(char*));
    pargv[pargc - 1] = NULL;

    execv(picat, pargv);
    ERR_AND_EXIT("execv");
}
