#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* luau = "/usr/local/bin/luau", *code = "code.luau";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

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

    int largc = argc + 4;
    char** largv = malloc(largc * sizeof(char*));
    largv[0] = (char*) luau;
    largv[1] = "-O2";
    largv[2] = "-g0";
    largv[3] = (char*) code;
    largv[4] = "-a";
    memcpy(&largv[5], &argv[2], (argc - 2) * sizeof(char*));
    largv[largc - 1] = NULL;

    execv(luau, largv);
    ERR_AND_EXIT("execv");
}
