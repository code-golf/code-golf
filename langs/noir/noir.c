#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* noir = "/usr/local/bin/nargo", *code = "src/main.nr";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(noir, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/noir"))
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

    int nargc = argc + 1;
    char** nargv = malloc(nargc * sizeof(char*));
    nargv[0] = (char*) noir;
    nargv[1] = "execute";
    memcpy(&nargv[2], &argv[2], (argc - 2) * sizeof(char*));
    nargv[nargc - 1] = NULL;

    execv(noir, nargv);
    ERR_AND_EXIT("execv");
}
