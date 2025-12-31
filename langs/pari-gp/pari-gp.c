#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* pari = "/usr/bin/gp", *code = "code.gp";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(pari, pari, "--version-short", NULL);
        ERR_AND_EXIT("execl");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen(code, "w")))
        ERR_AND_EXIT("fopen");

    ssize_t nbytes;

    for (int i = 2; i < argc; i++)
        nbytes += strlen(argv[i]);

    char* args;

    if (!(args = malloc(nbytes * sizeof(char*))))
        ERR_AND_EXIT("malloc");

    if (!strcat(args, "["))
        ERR_AND_EXIT("strcat");

    for (int i = 2; i < argc; i++) {
        if (i > 2)
            if (!strcat(args, ", "))
                ERR_AND_EXIT("strcat");

        if (!strcat(args, "\""))
            ERR_AND_EXIT("strcat");

        for (const char* c = argv[i]; *c; c++)
            switch (*c) {
                case '\\':
                    if (!strcat(args, "\\\\"))
                        ERR_AND_EXIT("strcat");
                    break;
                case '\"':
                    if (!strcat(args, "\\\""))
                        ERR_AND_EXIT("strcat");
                    break;
                case '\n':
                    if (!strcat(args, "\\n"))
                        ERR_AND_EXIT("strcat");
                    break;
                default:
                    if (!snprintf(&args[strlen(args)], sizeof(char*), "%c", *c))
                        ERR_AND_EXIT("snprintf");
            }

        if (!strcat(args, "\""))
            ERR_AND_EXIT("strcat");
    }

    if (!strcat(args, "]"))
        ERR_AND_EXIT("strcat");

    char buffer[4096];

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    if (!(fp = fopen("args", "w")))
        ERR_AND_EXIT("fopen");

    if (fputs(args, fp)) {
        if (fclose(fp))
            ERR_AND_EXIT("fclose");

        ERR_AND_EXIT("fputs");
    }

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    execl(pari, pari, "--fast", "--quiet", code, NULL);
    ERR_AND_EXIT("execl");

    free(args);
}
