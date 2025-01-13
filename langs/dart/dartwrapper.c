#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* dart = "/usr/bin/dart", *code = "/tmp/code.dart";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(dart, argv);
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

    int dargc = argc + 1;
    char** dargv = malloc(dargc * sizeof(char*));
    dargv[0] = (char*) dart;
    dargv[1] = (char*) code;
    memcpy(&dargv[2], &argv[2], (argc - 2) * sizeof(char*));
    dargv[dargc - 1] = NULL;

    execv(dart, dargv);
    ERR_AND_EXIT("execv");
}
