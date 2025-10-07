#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* mojo = "/usr/bin/pixi", *code = "code.mojo";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(mojo, mojo, "run", "mojo", argv[1], NULL);
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

    int margc = argc + 3;
    char** margv = malloc(margc * sizeof(char*));
    margv[0] = (char*) mojo;
    margv[1] = "run";
    margv[2] = "mojo";
    margv[3] = (char*) code;
    memcpy(&margv[4], &argv[2], (argc - 2) * sizeof(char*));
    margv[margc - 1] = NULL;

    execv(mojo, margv);
    ERR_AND_EXIT("execv");
}
