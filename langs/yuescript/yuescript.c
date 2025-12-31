#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* yuescript = "/usr/local/bin/yue", *lua = "/usr/bin/lua", *code = "code.yue";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(yuescript, yuescript, "-v", NULL);
        ERR_AND_EXIT("execl");
    }

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
        if (close(*fd))
            ERR_AND_EXIT("close");

        if (!dup2(fd[1], STDERR_FILENO) || !dup2(fd[1], STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        if (close(fd[1]))
            ERR_AND_EXIT("close");

        execl(yuescript, yuescript, code, NULL);
        ERR_AND_EXIT("execl");
    }

    if (close(fd[1]))
        ERR_AND_EXIT("close");

    if (!(fp = fdopen(*fd, "r")))
        ERR_AND_EXIT("fdopen");

    while (fgets(buffer, sizeof(buffer), fp)) {
        if (strstr(buffer, code))
            continue;

        if (!fputs(buffer, stderr))
            ERR_AND_EXIT("fputs");
    }

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if (remove(code))
        ERR_AND_EXIT("remove");

    int yargc = argc + 1;
    char** yargv = malloc(yargc * sizeof(char*));
    yargv[0] = (char*) lua;
    yargv[1] = "code.lua";
    memcpy(&yargv[2], &argv[2], (argc - 2) * sizeof(char*));
    yargv[yargc - 1] = NULL;

    execv(lua, yargv);
    ERR_AND_EXIT("execv");
}
