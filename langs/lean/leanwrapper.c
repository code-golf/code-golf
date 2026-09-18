#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* lean = "/usr/local/lean/bin/lean", *leanc = "/usr/local/lean/bin/leanc",
    *code = "code.lean", *codec = "code.c", *executable = "code";

int main(int argc, char* argv[]) {
    // guard against argv[1] being NULL
    if (argc > 1 && !strcmp(argv[1], "--version")) {
        execv(lean, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    // 1. write stdin to code.lean

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
    int status;

    // 2. lean -c code.c code.lean

    if (!(pid = fork())) {
        // Redirect stdout (1) to point to stderr (2)
        if (dup2(STDERR_FILENO, STDOUT_FILENO) < 0) {
            ERR_AND_EXIT("dup2");
        }
        execl(
            lean, lean, "-c", codec, code, NULL
        );
        ERR_AND_EXIT("execl");
    }

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    // 3. leanc -o code code.c

    if (!(pid = fork())) {
        execl(
            leanc, leanc, "-o", executable, codec, NULL
        );
        ERR_AND_EXIT("execl");
    }

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    // 4. remove code.lean and code.c

    if (remove(code))
        ERR_AND_EXIT("remove");

    if (remove(codec))
        ERR_AND_EXIT("remove");

    // 5. run executable with suitable argv
    int iargc = argc + 1;
    char** iargv = malloc(iargc * sizeof(char*));
    iargv[0] = (char *) executable;
    memcpy(&iargv[1], &argv[1], (argc - 1) * sizeof(char*));
    iargv[iargc - 1] = NULL;

    execv(executable, iargv);
    ERR_AND_EXIT("execv");
}