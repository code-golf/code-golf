#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* cjam = "/opt/java/bin/java", *code = "code.cjam";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

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

    int cargc = argc + 3;
    char** cargv = malloc(cargc * sizeof(char*));
    cargv[0] = (char*) cjam;
    cargv[1] = "-jar";
    cargv[2] = "/cjam.jar";
    cargv[3] = (char*) code;
    memcpy(&cargv[4], &argv[2], (argc - 2) * sizeof(char*));
    cargv[cargc - 1] = NULL;

    execv(cjam, cargv);
    ERR_AND_EXIT("execv");
}
