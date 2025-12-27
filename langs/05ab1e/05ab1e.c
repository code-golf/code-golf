#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* osabie = "/usr/bin/escript", *code = "code.abe", *input = "argv.txt";

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

    const char* table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

    for (int i = 2; i < argc; i++) {
        int len = strlen(argv[i]);

        char* arg, *str;

        if (!(arg = str = malloc((len + 2) / 3 * 4 + 1)))
            ERR_AND_EXIT("malloc");

        int j, chunk;

        typedef unsigned char UCHAR;

        for (j = 0; j + 2 < len; j += 3) {
            chunk = ((UCHAR*) argv[i])[j] << 16 | ((UCHAR*) argv[i])[j+1] << 8 | ((UCHAR*) argv[i])[j+2];

            *str++ = table[chunk >> 18 & 63];
            *str++ = table[chunk >> 12 & 63];
            *str++ = table[chunk >> 6 & 63];
            *str++ = table[chunk & 63];
        }

        if (j < len) {
            chunk = ((UCHAR*) argv[i])[j] << 16 | (j + 1 < len) * ((UCHAR*) argv[i])[j+1] << 8;

            *str++ = table[chunk >> 18 & 63];
            *str++ = table[chunk >> 12 & 63];
            *str++ = j + 1 < len ? table[chunk >> 6 & 63] : '=';
            *str++ = '=';
        }

        *str = 0;

        if (!write(fd, arg, strlen(arg)) || !write(fd, "\n", sizeof(char)))
            ERR_AND_EXIT("write");

        free(arg);
    }

    if (close(fd))
        ERR_AND_EXIT("close");

    pid_t pid;

    if (!(pid = fork())) {
        if (dup2(open(input, O_RDONLY), STDIN_FILENO))
            ERR_AND_EXIT("dup2");

        execl(osabie, osabie, "/usr/local/bin/osabie", code, NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if (remove(code))
        ERR_AND_EXIT("remove");
}
