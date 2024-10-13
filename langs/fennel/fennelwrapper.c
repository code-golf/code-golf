#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

int main (int argc, char *argv[]) {
    if (argc > 1 && strcmp(argv[1], "-v") == 0) {
        execv("/usr/bin/fennel", argv);
        ERR_AND_EXIT("execl");
    }

    FILE* fp = fopen("/tmp/code.fnl", "w");
    if (!fp)
        ERR_AND_EXIT("fopen");

    // Copy STDIN into code.fnl
    char buffer[4096];
    ssize_t nbytes;
    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t)nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp) < 0)
        ERR_AND_EXIT("fclose");

    execv("/usr/bin/fennel", argv);
    ERR_AND_EXIT("execl");
}
