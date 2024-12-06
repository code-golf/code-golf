#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* vyxal = "/usr/bin/vyxal", *code = "code.vy";

int main(int argc, char* argv[]) {
    if (argc && !strcmp(argv[1], "--version")) {
        execl(vyxal, vyxal, "-", "⋎", NULL);
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

    int vargc = argc + 2;
    char** vargv = malloc(vargc * sizeof(char*));
    vargv[0] = (char*) vyxal;
    vargv[1] = (char*) code;
    vargv[2] = "MOa";
    memcpy(&vargv[3], &argv[2], (argc - 2) * sizeof(char*));
    vargv[vargc - 1] = NULL;

    execv(vyxal, vargv);
    ERR_AND_EXIT("execv");
}
