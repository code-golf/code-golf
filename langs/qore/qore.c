#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* qore = "/usr/local/bin/qore", *code = "code.q";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(qore, argv);
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

    int qargc = argc + 3;
    char** qargv = malloc(qargc * sizeof(char*));
    qargv[0] = (char*) qore;
    qargv[1] = "--allow-bare-refs";
    qargv[2] = "--no-external-access";
    qargv[3] = (char*) code;
    memcpy(&qargv[4], &argv[2], (argc - 2) * sizeof(char*));
    qargv[qargc - 1] = NULL;

    execv(qore, qargv);
    ERR_AND_EXIT("execv");
}
