#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* erlang = "/usr/bin/escript", *code = "main.erl";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen(code, "w")))
        ERR_AND_EXIT("fopen");

    if (!fputs("-module(main).\n-export([main/1]).\n", fp))
        ERR_AND_EXIT("fputs");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    pid_t pid;

    if (!(pid = fork())) {
        if (!dup2(STDERR_FILENO, STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        execl("/usr/bin/erlc", "erlc", code, NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if (remove(code))
        ERR_AND_EXIT("remove");

    int eargc = argc + 1;
    char** eargv = malloc(eargc * sizeof(char*));
    eargv[0] = (char*) erlang;
    eargv[1] = "main.beam";
    memcpy(&eargv[2], &argv[2], (argc - 2) * sizeof(char*));
    eargv[eargc - 1] = NULL;

    execv(erlang, eargv);
    ERR_AND_EXIT("execv");
}
