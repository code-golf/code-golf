#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* D = "/usr/bin/dmd", *code = "code.d";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(D, argv);
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

    int dargc = argc + 3;
    char** dargv = malloc(dargc * sizeof(char*));
    dargv[0] = (char*) D;
    dargv[1] = "-color=on";
    dargv[2] = "-run";
    dargv[3] = (char*) code;
    memcpy(&dargv[4], &argv[2], (argc - 2) * sizeof(char*));
    dargv[dargc - 1] = NULL;

    execv(D, dargv);
    ERR_AND_EXIT("execv");
}
