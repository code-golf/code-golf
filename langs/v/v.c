#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* V = "/usr/local/v", *code = "code.v";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(V, V, "version", NULL);
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

    int vargc = argc + 3;
    char** vargv = malloc(vargc * sizeof(char*));
    vargv[0] = (char*) V;
    vargv[1] = "-color";
    vargv[2] = "run";
    vargv[3] = (char*) code;
    memcpy(&vargv[4], &argv[2], (argc - 2) * sizeof(char*));
    vargv[vargc - 1] = NULL;

    execv(V, vargv);
    ERR_AND_EXIT("execv");
}
