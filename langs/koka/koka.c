#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* koka = "/usr/local/bin/koka", *code = "code.kk";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(koka, argv);
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

    int kargc = argc + 5;
    char** kargv = malloc(kargc * sizeof(char*));
    kargv[0] = (char*) koka;
    kargv[1] = "--target=js";
    kargv[2] = "-v0";
    kargv[3] = "-e";
    kargv[4] = (char*) code;
    kargv[5] = "--";
    memcpy(&kargv[6], &argv[2], (argc - 2) * sizeof(char*));
    kargv[kargc - 1] = NULL;

    execv(koka, kargv);
    ERR_AND_EXIT("execv");
}
