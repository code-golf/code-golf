#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

int main(int argc, char* argv[]) {
    if (argc > 1 && strcmp(argv[1], "-h") == 0) {
        execl("/opt/java/openjdk/bin/java", "java", "-jar", "/rocky.jar", "-h", NULL);
        ERR_AND_EXIT("execl");
    }

    if(chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp = fopen("code.rock", "w");
    if (!fp)
        ERR_AND_EXIT("fopen");

    // Copy STDIN into code.rock
    char buffer[4096];
    ssize_t nbytes;
    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t)nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp) < 0)
        ERR_AND_EXIT("fclose");

    execl("/opt/java/openjdk/bin/java", "java", "-jar", "/rocky.jar",
        "run", "--infinite-loops", "--rocky", "code.rock", NULL);
    ERR_AND_EXIT("execl");
}
