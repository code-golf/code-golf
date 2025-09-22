#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* haskell = "/usr/local/bin/ghc", *code = "/tmp/code.hs";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(haskell, argv);
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

    int hargc = argc + 5;
    char** hargv = malloc(hargc * sizeof(char*));
    hargv[0] = (char*) haskell;
    hargv[1] = "--run";
    hargv[2] = (char*) code;
    hargv[3] = "-fdiagnostics-color=always";
    hargv[4] = "--";
    hargv[5] = "";
    memcpy(&hargv[6], &argv[2], (argc - 2) * sizeof(char*));
    hargv[hargc - 1] = NULL;

    execv(haskell, hargv);
    ERR_AND_EXIT("execv");
}
