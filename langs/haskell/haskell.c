#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

#define VERSION "9.4.4"

int main(int argc, char* argv[]) {
    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        puts(VERSION);
        return 0;
    }

    FILE* fp = fopen("/tmp/code.hs", "w");
    if (!fp)
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t)nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp) < 0)
        ERR_AND_EXIT("fclose");

    int hargc = argc + 3;
    char **hargv = malloc(hargc * sizeof(char*));
    hargv[0] = "runghc-" VERSION;
    hargv[1] = "--ghc-arg=-fdiagnostics-color=always";
    hargv[2] = "/tmp/code.hs";
    hargv[3] = "--";
    memcpy(&hargv[4], &argv[2], (argc - 2) * sizeof(char*));
    hargv[hargc - 1] = NULL;

    execv("/usr/lib/ghc-" VERSION "/bin/runghc-" VERSION, hargv);
    ERR_AND_EXIT("execv");
}
