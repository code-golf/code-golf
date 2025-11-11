#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* squirrel = "/usr/local/bin/sq_static", *code = "code.nut";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(squirrel, squirrel, "-v", NULL);
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

    int sargc = argc + 1;
    char** sargv = malloc(sargc * sizeof(char*));
    sargv[0] = (char*) squirrel;
    sargv[1] = (char*) code;
    memcpy(&sargv[2], &argv[2], (argc - 2) * sizeof(char*));
    sargv[sargc - 1] = NULL;

    execv(squirrel, sargv);
    ERR_AND_EXIT("execv");
}
