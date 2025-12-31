#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* ox = "/usr/local/bin/oxl", *code = "code.ox";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

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
        if (!dup2(STDERR_FILENO, STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        execl(ox, ox, "-b", "-r-", code, NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    int oargc = argc + 2;
    char** oargv = malloc(oargc * sizeof(char*));
    oargv[0] = (char*) ox;
    oargv[1] = "-b";
    oargv[2] = (char*) code;
    memcpy(&oargv[3], &argv[2], (argc - 2) * sizeof(char*));
    oargv[oargc - 1] = NULL;

    execv(ox, oargv);
    ERR_AND_EXIT("execv");
}
