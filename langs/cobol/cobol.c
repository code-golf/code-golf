#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

int main (int argc, char *argv[]) {
    char buffer[4096];
    ssize_t nbytes;
    FILE *fp = fopen("/tmp/code.cbl", "w");

    if (!fp)
        return 1;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != nbytes)
            return 2;

    fclose(fp);

    pid_t pid = fork();
    if (!pid) {
        execl("/usr/local/bin/cobc", "/usr/local/bin/cobc", "-FCx",
              "/tmp/code.cbl", "-o", "/tmp/code.c", NULL);
        perror("execl");
        return 3;
    }

    int status;             
    waitpid(pid, &status, 0); 

    if (!WIFEXITED(status))
        return 4;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if (remove("/tmp/code.cbl")) {
        perror("remove");
        return 5;
    }

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
    return 6;
}
