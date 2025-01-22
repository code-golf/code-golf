#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* scala = "/usr/local/bin/scala", *code = "code.scala";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(scala, scala, "-version", NULL);
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

    int sargc = argc + 1;
    char** sargv = malloc(sargc * sizeof(char*));
    sargv[0] = (char*) scala;
    sargv[1] = (char*) code;
    memcpy(&sargv[2], &argv[2], (argc - 2) * sizeof(char*));
    sargv[sargc - 1] = NULL;

    execv(scala, sargv);
    ERR_AND_EXIT("execv");
}
