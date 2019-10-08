#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/wait.h>

#define NASM "/usr/bin/nasm"
#define SOURCE "/tmp/code.asm"
#define BINARY "/tmp/code"

int main (int argc, char *argv[]) {
    char buffer[4096];
    ssize_t nbytes;
    FILE *fp = fopen(SOURCE, "w");

    if (fp == 0) {
        perror("Error opening " SOURCE);
        return 1;
    }

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)) > 0) {
        if (fwrite(buffer, sizeof(char), nbytes, fp) != nbytes) {
            perror("Error writing to " SOURCE);
            return 1;
        }
    }

    fclose(fp);

    int pid = fork();
    if (pid == 0) {
        char *nasm_argv[] = {NASM, "-f", "bin", "-o", BINARY, "/header.asm", NULL};
        execv(NASM, nasm_argv);
        perror("Error executing " NASM);
    } else if (pid > 0) {
        int status;
        // wait for forked nasm
        if (wait(&status) == -1) {
            perror("Error waiting for child process");
            return 1;
        }
        if (WIFEXITED(status)) {
            if (WEXITSTATUS(status) != 0) {
                // nasm exited with non-zero
                return WEXITSTATUS(status);
            }
        } else {
            // nasm exited in some other way
            return 1;
        }
        chmod(BINARY, 0777);
        execv(BINARY, argv);
        perror("Error executing " BINARY);
    } else {
        perror("Error forking");
    }

    return 1;
}
