#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

int main (int argc, char *argv[]) {
    pid_t pid = fork();
    if (!pid) {
        execl("/usr/bin/cobc", "/usr/bin/cobc", "-CFxo", "/tmp/code.c", "-", NULL);
        perror("execl");
        return 1;
    }

    int status;
    waitpid(pid, &status, 0); 

    if (!WIFEXITED(status))
        return 2;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    int ntcc = argc + 3;
    char **tcc = malloc(ntcc * sizeof(char*));
    tcc[0] = "/usr/bin/tcc";
    tcc[1] = "-lcob";
    tcc[2] = "-run";
    tcc[3] = "/tmp/code.c";
    memcpy(&tcc[4], &argv[2], (argc - 2) * sizeof(char*));
    tcc[ntcc - 1] = NULL;

    execv("/usr/bin/tcc", tcc);
    perror("execv");
    return 3;
}
