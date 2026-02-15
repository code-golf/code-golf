#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* iogii = "/usr/local/bin/iogii", *ghcii = "/usr/bin/ghcii", *code[2] = {"code.iog", "code.ism"}, *input = "argv.txt";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(iogii, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen(code[0], "w")))
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    pid_t pid;
    int fd;

    if (!(pid = fork())) {
        if (!dup2(fd = open(code[1], O_CREAT | O_TRUNC | O_WRONLY, 0644), STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        if (close(fd))
            ERR_AND_EXIT("close");

        execl(iogii, iogii, "-e", code[0], NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if (remove(code[0]))
        ERR_AND_EXIT("remove");

    int iargc = argc + 1;
    char** iargv = malloc(iargc * sizeof(char*));
    iargv[0] = (char*) ghcii;
    iargv[1] = (char*) code[1];
    memcpy(&iargv[2], &argv[2], (argc - 2) * sizeof(char*));
    iargv[iargc - 1] = NULL;
    execv(ghcii, iargv);
    ERR_AND_EXIT("execv");

    if (remove(code[1]))
        ERR_AND_EXIT("remove");
}
