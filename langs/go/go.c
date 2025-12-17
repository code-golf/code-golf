#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* go = "/usr/local/go/bin/go", *code = "code.go";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(go, go, "version", NULL);
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

    int gargc = argc + 2;
    char** gargv = malloc(gargc * sizeof(char*));
    gargv[0] = (char*) go;
    gargv[1] = "run";
    gargv[2] = (char*) code;
    memcpy(&gargv[3], &argv[2], (argc - 2) * sizeof(char*));
    gargv[gargc - 1] = NULL;

    execv(go, gargv);
    ERR_AND_EXIT("execv");
}
