#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* odin = "/usr/local/bin/odin", *code = "code.odin";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(odin, odin, "version", NULL);
        ERR_AND_EXIT("execl");
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

    int oargc = argc + 4;
    char** oargv = malloc(oargc * sizeof(char*));
    oargv[0] = (char*) odin;
    oargv[1] = "run";
    oargv[2] = (char*) code;
    oargv[3] = "-file";
    oargv[4] = "-o:none";
    memcpy(&oargv[5], &argv[2], (argc - 2) * sizeof(char*));
    oargv[oargc - 1] = NULL;

    execv(odin, oargv);
    ERR_AND_EXIT("execv");
}
