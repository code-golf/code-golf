#include <dirent.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* reason = "/usr/local/bin/dune", *code = "code.re";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen("dune", "w")))
        ERR_AND_EXIT("fopen");

    if (fputs("(executable (name code))", fp))
        ERR_AND_EXIT("fputs");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    if (!(fp = fopen("dune-project", "w")))
        ERR_AND_EXIT("fopen");

    if (fputs("(lang dune 3.20)", fp))
        ERR_AND_EXIT("fputs");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    if (!(fp = fopen(code, "w")))
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    int rargc = argc + 3;
    char** rargv = malloc(rargc * sizeof(char*));
    rargv[0] = (char*) reason;
    rargv[1] = "exec";
    rargv[2] = "./code.exe";
    rargv[3] = "--";
    memcpy(&rargv[4], &argv[2], (argc - 2) * sizeof(char*));
    rargv[rargc - 1] = NULL;

    execv(reason, rargv);
    ERR_AND_EXIT("execv");
}
