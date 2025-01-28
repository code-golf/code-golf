#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* pony = "/usr/local/bin/ponyc", *code = "code.pony";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(pony, argv);
        ERR_AND_EXIT("execv");
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

    pid_t pid;

    if (!(pid = fork())) {
        execl(pony, pony, "--verbose", "0", NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if (remove(code))
        ERR_AND_EXIT("remove");

    int pargc = argc;
    char** pargv = malloc(pargc * sizeof(char*));
    pargv[0] = "tmp";
    memcpy(&pargv[1], &argv[2], (argc - 2) * sizeof(char*));
    pargv[pargc - 1] = NULL;

    execv("tmp", pargv);
    ERR_AND_EXIT("execv");
}
