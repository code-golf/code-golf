#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* iogii = "/usr/local/bin/iogii", *ghcii = "/usr/bin/ghcii", *code[2] = {"/proc/self/fd/0", "code.ism"};

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(iogii, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    pid_t pid;

    if (!(pid = fork())) {
        int fd;

        if (!dup2(fd = open(code[1], O_CREAT | O_TRUNC | O_WRONLY, 0644), STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        if (close(fd))
            ERR_AND_EXIT("close");

        execl(iogii, iogii, "--engine", code[0], NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    int iargc = argc + 1;
    char** iargv = malloc(iargc * sizeof(char*));
    iargv[0] = (char*) ghcii;
    iargv[1] = (char*) code[1];
    memcpy(&iargv[2], &argv[2], (argc - 2) * sizeof(char*));
    iargv[iargc - 1] = NULL;

    execv(ghcii, iargv);
    ERR_AND_EXIT("execv");
}
