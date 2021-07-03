#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

int main (int argc, char *argv[]) {
    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        execv("/usr/bin/cobc", argv);
        perror("execv");
        return 0;
    }

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

    char buffer[4096];
    ssize_t nbytes;
    FILE *fp1 = fopen("/tmp/code.c", "r");
    FILE *fp2 = fopen("/tmp/code2.c", "w");

    if (!fp1 || !fp2)
        return 3;

    const char *header = "#include <math.h>\n";

    if (fwrite(header, sizeof(char), strlen(header), fp2) != strlen(header))
        return 4;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), fp1)) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp2) != nbytes)
            return 5;

    fclose(fp1);
    fclose(fp2);

    int ntcc = argc + 3;
    char **tcc = malloc(ntcc * sizeof(char*));
    tcc[0] = "/usr/bin/tcc";
    tcc[1] = "-lcob";
    tcc[2] = "-run";
    tcc[3] = "/tmp/code2.c";
    memcpy(&tcc[4], &argv[2], (argc - 2) * sizeof(char*));
    tcc[ntcc - 1] = NULL;

    execv("/usr/bin/tcc", tcc);
    perror("execv");
    return 6;
}
