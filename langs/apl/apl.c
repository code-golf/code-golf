#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* apl = "/usr/bin/dyalogscript", *code = "code.dyalog";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(apl, apl, "apl.dws", NULL);
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

    int aargc = argc + 1;
    char** aargv = malloc(aargc * sizeof(char*));
    aargv[0] = (char*) apl;
    aargv[1] = (char*) code;
    memcpy(&aargv[2], &argv[2], (argc - 2) * sizeof(char*));
    aargv[aargc - 1] = NULL;

    execv(apl, aargv);
    ERR_AND_EXIT("execv");
}
