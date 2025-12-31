#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* crap = "/usr/local/bin/crap", *C = "/usr/bin/c", *code = "code.crap";

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

    int fd[2];

    if (pipe(fd))
        ERR_AND_EXIT("pipe");

    pid_t pid;

    if (!(pid = fork())) {
        close(fd[0]);

        if (!dup2(fd[1], STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        close(fd[1]);

        execlp(crap, crap, code, NULL);
        ERR_AND_EXIT("execlp");
    }

    if (dup2(fd[0], STDIN_FILENO))
        ERR_AND_EXIT("dup2");

    if (close(fd[0]) || close(fd[1]))
        ERR_AND_EXIT("close");

    int cargc = argc + 2;
    char** cargv = malloc(cargc * sizeof(char*));
    cargv[0] = (char*) C;
    cargv[1] = "-run";
    cargv[2] = "-";
    memcpy(&cargv[3], &argv[2], (argc - 2) * sizeof(char*));
    cargv[cargc - 1] = NULL;

    execvp(C, cargv);
    ERR_AND_EXIT("execvp");
}
