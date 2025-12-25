#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* gleam = "/usr/local/bin/gleam", *code = "src/main.gleam";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(gleam, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/usr/local"))
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
        if (!dup2(STDERR_FILENO, STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        execl(gleam, gleam, "build", "--no-print-progress", NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    int gargc = argc + 3;
    char** gargv = malloc(gargc * sizeof(char*));
    gargv[0] = (char*) gleam;
    gargv[1] = "run";
    gargv[2] = "--no-print-progress";
    gargv[3] = "--";
    memcpy(&gargv[4], &argv[2], (argc - 2) * sizeof(char*));
    gargv[gargc - 1] = NULL;

    execv(gleam, gargv);
    ERR_AND_EXIT("execv");
}
