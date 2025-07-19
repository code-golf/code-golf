#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* umka = "/usr/local/bin/umka", *code = "code.um";

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

    int uargc = argc + 2;
    char** uargv = malloc(uargc * sizeof(char*));
    uargv[0] = (char*) umka;
    uargv[1] = "-warn";
    uargv[2] = (char*) code;
    memcpy(&uargv[3], &argv[2], (argc - 2) * sizeof(char*));
    uargv[uargc - 1] = NULL;

    execv(umka, uargv);
    ERR_AND_EXIT("execv");
}
