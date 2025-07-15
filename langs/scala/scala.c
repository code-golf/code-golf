#include <stdio.h>
#include <stdlib.h>
#include <string.h>
// #include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* scala = "/opt/java/bin/java", *scalaRunner = "/usr/libexec/scala-cli.jar", *code = "code.scala";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

    // if (!strcmp(argv[1], "--version")) {
    //     execl(scala, scala, "-jar", scalaRunner, argv[1], NULL);
    //     ERR_AND_EXIT("execl");
    // }

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

    // pid_t pid;

    // if (!(pid = fork())) {
    //     execl(scala, scala, "-jar", scalaRunner, "--power", "bloop", NULL);
    //     ERR_AND_EXIT("execl");
    // }

    // int status;

    // waitpid(pid, &status, 0);

    // if (!WIFEXITED(status))
    //     exit(EXIT_FAILURE);

    // if (WEXITSTATUS(status))
    //     return WEXITSTATUS(status);

    int sargc = argc + 5;
    char** sargv = malloc(sargc * sizeof(char*));
    sargv[0] = (char*) scala;
    sargv[1] = "-jar";
    sargv[2] = (char*) scalaRunner;
    sargv[3] = "run";
    sargv[4] = (char*) code;
    sargv[5] = "--";
    memcpy(&sargv[6], &argv[2], (argc - 2) * sizeof(char*));
    sargv[sargc - 1] = NULL;

    execv(scala, sargv);
    ERR_AND_EXIT("execv");
}
