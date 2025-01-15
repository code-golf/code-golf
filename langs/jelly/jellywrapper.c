#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* jelly = "/usr/bin/jelly", *code = "code.jelly";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

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

    strcat(args, "[");

    for (int i = 2; i < argc; i++) {
        if (i > 2)
            strcat(args, ", ");

        strcat(args, "\"");

        for (const char* c = argv[i]; *c; c++)
            switch (*c) {
                case '\\':
                    strcat(args, "\\\\");
                    break;
                case '\"':
                    strcat(args, "\\\"");
                    break;
                case '\n':
                    strcat(args, "\\n");
                    break;
                default:
                    if (!snprintf(&args[strlen(args)], sizeof(char*), "%c", *c))
                        ERR_AND_EXIT("snprintf");
            }

        strcat(args, "\"");
    }

    strcat(args, "]");

    char buffer[4096];

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    execl(jelly, jelly, "fun", code, args, NULL);
    ERR_AND_EXIT("execl");

    free(args);
}
