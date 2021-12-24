#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

const char* fortran = "/usr/local/bin/gfortran";
const char* code = "/tmp/code.f90";
const char* bin = "/tmp/code";

int main(int argc, char* argv[]) {
    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        execv(fortran, argv);
        perror("execv");
        return 0;
    }

    FILE* fp = fopen(code, "w");

    if (!fp)
        return 1;

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != nbytes)
            return 2;

    fclose(fp);

    pid_t pid = fork();
    if (!pid) {
        dup2(2, 1);
        execl(fortran, fortran, "-fbackslash", "-o", bin, code, NULL);
        perror("execl");
        return 3;
    }

    int status;
    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 4;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if(remove(code)) {
        perror("Error deleting file");
        return 5;
    }

    int fargc = argc;
    char** fargv = malloc(fargc * sizeof(char*));
    fargv[0] = (char*)bin;
    memcpy(&fargv[1], &argv[2], (argc - 2) * sizeof(char*));
    fargv[fargc - 1] = NULL;

    execv(bin, fargv);
    perror("execv");
    return 6;
}
