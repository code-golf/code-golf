#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

int main(int argc, char* argv[]) {
    if (setenv("PATH", "/usr/lib/fpc/3.2.2:/usr/bin", 1))
        ERR_AND_EXIT("setenv");

    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        execl("/usr/bin/fpc", "fpc", "-iV", NULL);
        ERR_AND_EXIT("execv");
    }

    if(chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp = fopen("code.pas", "w");
    if (!fp)
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t)nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp) < 0)
        ERR_AND_EXIT("fclose");

    pid_t pid = fork();
    if (!pid) {
        // Move the compiler output to STDERR.
        if (dup2(STDERR_FILENO, STDOUT_FILENO) < 0)
            ERR_AND_EXIT("dup2");

        execl("/usr/bin/fpc", "fpc", "code.pas", NULL);
        ERR_AND_EXIT("execl");
    }

    int status;
    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 1;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if(remove("code.pas"))
        ERR_AND_EXIT("remove");

    int pargc = argc;
    char** pargv = malloc(pargc * sizeof(char*));
    pargv[0] = "code";
    memcpy(&pargv[1], &argv[2], (argc - 2) * sizeof(char*));
    pargv[pargc - 1] = NULL;

    execv("code", pargv);
    ERR_AND_EXIT("execv");
}
