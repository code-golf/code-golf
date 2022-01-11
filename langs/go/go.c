#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

const char* go = "/usr/local/go/bin/go";
const char* code = "/tmp/code.go";
const char* object = "/tmp/code.o";
const char* bin = "/tmp/code";

int main(int argc, char* argv[]) {
    if (argc > 1 && strcmp(argv[1], "version") == 0) {
        execv(go, argv);
        perror("execv");
        return 0;
    }

    FILE* fp = fopen(code, "w");

    if (!fp)
        return 1;

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != nbytes)
            return 2;

    fclose(fp);

    pid_t pid = fork();
    if (!pid) {
        dup2(2, 1);
        execl(go, go, "tool", "compile", "-o", object, code, NULL);
        perror("execl");
        return 3;
    }

    int status;
    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 4;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    pid = fork();
    if (!pid) {
        dup2(2, 1);
        execl(go, go, "tool", "link", "-o", bin, object, NULL);
        perror("execl");
        return 5;
    }

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 6;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if(remove(code)) {
        perror("Error deleting file");
        return 7;
    }

    int gargc = argc;
    char** gargv = malloc(gargc * sizeof(char*));
    gargv[0] = (char*)bin;
    memcpy(&gargv[1], &argv[2], (argc - 2) * sizeof(char*));
    gargv[gargc - 1] = NULL;

    execv(bin, gargv);
    perror("execv");
    return 8;
}
