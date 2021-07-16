#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main (int argc, char *argv[]) {
    char buffer[4096];
    ssize_t nbytes;
    FILE *fp = fopen("/tmp/code.jl", "w");

    if (fp == 0)
        return 1;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != nbytes)
            return 2;

    fclose(fp);

    putenv("HOME=/");

    execv("/usr/bin/julia", argv);
    perror("execv");
}
