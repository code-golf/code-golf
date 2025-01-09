#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* rexx = "/usr/bin/regina", *code = "code.rexx";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(rexx, argv);
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

    pid_t pid;

    if (!(pid = fork())) {
        execl(rexx, rexx, "--compile", "code", code, NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    int rargc = argc + 3;
    char** rargv = malloc(rargc * sizeof(char*));
    rargv[0] = (char*) rexx;
    rargv[1] = "--execute";
    rargv[2] = "--args";
    rargv[3] = "code";
    memcpy(&rargv[4], &argv[2], (argc - 2) * sizeof(char*));
    rargv[rargc - 1] = NULL;

    execv(rexx, rargv);
    ERR_AND_EXIT("execv");
}
