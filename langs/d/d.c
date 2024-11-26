#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* dmd = "/usr/bin/dmd";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(dmd, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen("code.d", "w")))
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    int dargc = argc + 3;
    char** dargv = malloc(dargc * sizeof(char*));
    dargv[0] = (char*) dmd;
    dargv[1] = "-color=on";
    dargv[2] = "-run";
    dargv[3] = "code.d";
    memcpy(&dargv[4], &argv[2], (argc - 2) * sizeof(char*));
    dargv[dargc - 1] = NULL;

    execv(dmd, dargv);
    ERR_AND_EXIT("execv");
}
