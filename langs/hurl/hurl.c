#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* hurl = "/usr/local/bin/hurl", *code = "code.hurl", *input = "argv.txt";

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

    for (int i = 2; i < argc; i++) {
        for (const char* c = argv[i]; *c; c++)
            if (!snprintf(&args[strlen(args)], sizeof(char*), "%c", *c))
                ERR_AND_EXIT("snprintf");

        strcat(args, "\n");
    }

    char buffer[4096];

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    if (!(fp = fopen(input, "w")))
        ERR_AND_EXIT("fopen");

    if (!fputs(args, fp)) {
        if (fclose(fp))
            ERR_AND_EXIT("fclose");

        ERR_AND_EXIT("fputs");
    }

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    execl(hurl, hurl, "run", code, NULL);
    ERR_AND_EXIT("execl");

    free(args);
}
