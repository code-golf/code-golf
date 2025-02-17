#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* fennel = "/usr/local/bin/fennel", *code = "code.fnl";

int main (int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(fennel, argv);
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

    int fargc = argc + 1;
    char** fargv = malloc(fargc * sizeof(char*));
    fargv[0] = (char*) fennel;
    fargv[1] = (char*) code;
    memcpy(&fargv[2], &argv[2], (argc - 2) * sizeof(char*));
    fargv[fargc - 1] = NULL;

    execv(fennel, fargv);
    ERR_AND_EXIT("execv");
}
