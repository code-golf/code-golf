#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* octave = "/usr/local/bin/octave", *code = "code.m";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(octave, argv);
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

    int oargc = argc + 2;
    char** oargv = malloc(oargc * sizeof(char*));
    oargv[0] = (char*) octave;
    oargv[1] = "--quiet";
    oargv[2] = (char*) code;
    memcpy(&oargv[3], &argv[2], (argc - 2) * sizeof(char*));
    oargv[oargc - 1] = NULL;

    execv(octave, oargv);
    ERR_AND_EXIT("execv");
}
