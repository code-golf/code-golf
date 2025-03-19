#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* rockstar = "/opt/java/bin/java", *code = "code.rock";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl(rockstar, rockstar, "-jar", "rocky.jar", "-h", NULL);
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

    int rargc = argc + 6;
    char** rargv = malloc(rargc * sizeof(char*));
    rargv[0] = (char*) rockstar;
    rargv[1] = "-jar";
    rargv[2] = "/rocky.jar";
    rargv[3] = "run";
    rargv[4] = "--infinite-loops";
    rargv[5] = "--rocky";
    rargv[6] = (char*) code;
    memcpy(&rargv[7], &argv[2], (argc - 2) * sizeof(char*));
    rargv[rargc - 1] = NULL;

    execv(rockstar, rargv);
    ERR_AND_EXIT("execv");
}
