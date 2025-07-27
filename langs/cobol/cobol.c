#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* cobol = "/usr/local/bin/cobc", *C = "/usr/bin/c", *code = "code.c";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(cobol, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    pid_t pid;

    if (!(pid = fork())) {
        execl(cobol, cobol, "-CFxo", code, "-", NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0); 

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    int cargc = argc + 3;
    char** cargv = malloc(cargc * sizeof(char*));
    cargv[0] = (char*) C;
    cargv[1] = "-lcob";
    cargv[2] = "-run";
    cargv[3] = (char*) code;
    memcpy(&cargv[4], &argv[2], (argc - 2) * sizeof(char*));
    cargv[cargc - 1] = NULL;

    execv(C, cargv);
    ERR_AND_EXIT("execv");
}
