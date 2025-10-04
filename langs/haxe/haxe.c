#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* haxe = "/usr/local/bin/haxe", *code = "Main.hx";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(haxe, argv);
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

    int hargc = argc + 2;
    char** hargv = malloc(hargc * sizeof(char*));
    hargv[0] = (char*) haxe;
    hargv[1] = "--run";
    hargv[2] = "Main";
    memcpy(&hargv[3], &argv[2], (argc - 2) * sizeof(char*));
    hargv[hargc - 1] = NULL;

    execv(haxe, hargv);
    ERR_AND_EXIT("execv");
}
