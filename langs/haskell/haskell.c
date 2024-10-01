#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

int main(int argc, char* argv[]) {
    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        execv("/usr/bin/ghc", argv);
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

    char **hargv = malloc(argc * sizeof(char*));
    hargv[0] = "--ghc-arg=-fdiagnostics-color=always";
    hargv[1] = "/tmp/code.hs";
    memcpy(&hargv[2], &argv[1], argc * sizeof(char*));

    execv("/usr/bin/runghc", hargv);
    ERR_AND_EXIT("execv");
}
