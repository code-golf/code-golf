#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* viml = "/usr/bin/vim";

char* make_args_list(int argc, char* argv[]) {
    size_t capacity = sizeof("let args = []");
    for (int i = 2; i < argc; i++) {
        capacity += 2 * strlen(argv[i]); // Assume worst-case escaping.
        capacity += 3; // Quotes and comma.
    }

    char* ret = malloc(capacity);
    char* out = stpcpy(ret, "let args = [");

    for (int i = 2; i < argc; i++) {
        *out++ = '"';

        for (const char* c = argv[i]; *c; c++)
            switch (*c) {
                case '\\':
                    *out++ = '\\';
                    *out++ = '\\';
                    break;
                case '\"':
                    *out++ = '\\';
                    *out++ = '"';
                    break;
                case '\n':
                    *out++ = '\\';
                    *out++ = 'n';
                    break;
                default:
                    *out++ = *c;
            }

        *out++ = '"';
        *out++ = ',';
    }

    *out++ = ']';
    *out++ = '\0';

    return ret;
}

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(viml, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen("output", "w")))
        ERR_AND_EXIT("fopen");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    pid_t pid;

    if (!(pid = fork())) {
        if (!dup2(STDERR_FILENO, STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        execl(
            viml,
            viml,
            "--clean",
            "--not-a-term",
            "--noplugin",
            "-eZ",
            "-V1",
            "-c", make_args_list(argc, argv),
            "-S", "/proc/self/fd/0",
            "output",
            NULL
        );

        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    char stdout_buf[4096];
    ssize_t nbytes;

    if (!(fp = fopen("output", "r")))
        ERR_AND_EXIT("fopen");

    while ((nbytes = fread(stdout_buf, sizeof(char), sizeof(stdout_buf), fp)))
        if (fwrite(stdout_buf, sizeof(char), nbytes, stdout) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    return WEXITSTATUS(status);
}
