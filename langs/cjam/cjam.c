#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* java = "/opt/java/bin/java";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        return 0;

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen("code.cjam", "w")))
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
    cargv[0] = (char*) java;
    cargv[1] = "-jar";
    cargv[2] = "/cjam.jar";
    cargv[3] = "code.cjam";
    memcpy(&cargv[4], &argv[2], (argc - 2) * sizeof(char*));
    cargv[cargc - 1] = NULL;

    execv(java, cargv);
    ERR_AND_EXIT("execv");
}
