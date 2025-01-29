#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* stax = "/usr/bin/stax", *code = "code.stax", *input = "argv.txt";

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

    int fd;

    if (!(fd = open(input, O_CREAT | O_TRUNC | O_WRONLY, 0644)))
        ERR_AND_EXIT("open");

    for (int i = 2; i < argc; i++)
        if (write(fd, argv[i], strlen(argv[i])) < 0 || write(fd, "\n", sizeof(char)) < 0) {
            if (close(fd))
                ERR_AND_EXIT("close");

            ERR_AND_EXIT("write");
        }

    if (close(fd))
        ERR_AND_EXIT("close");

    pid_t pid;

    if (!(pid = fork())) {
        if (dup2(open(input, O_RDONLY), STDIN_FILENO))
            ERR_AND_EXIT("dup2");

        execl(stax, stax, "-u", code, NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);
}
