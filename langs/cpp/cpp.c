#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* cpp = "/usr/bin/clang", *code = "code.cpp";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(cpp, argv);
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
        execl(
            cpp, cpp, "-std=c++23", "--target=x86_64-alpine-linux-musl",
            "-O2", "-lstdc++", "-fcolor-diagnostics", "-o", "code", code,
            "/unbuffered.cpp", NULL
        );
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

    int cargc = argc;
    char** cargv = malloc(cargc * sizeof(char*));
    cargv[0] = "code";
    memcpy(&cargv[1], &argv[2], (argc - 2) * sizeof(char*));
    cargv[cargc - 1] = NULL;

    execv("code", cargv);
    ERR_AND_EXIT("execv");
}
