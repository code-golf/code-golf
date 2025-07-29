#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* java = "/opt/java/bin/java", *javac = "/opt/java/bin/javac", *code = "code.java";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(java, argv);
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

    pid_t pid;

    if (!(pid = fork())) {
        execl(javac, javac, code, NULL);
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

    DIR* dir;

    if (!(dir = opendir(".")))
        ERR_AND_EXIT("opendir");

    char class[256];

    memset(class, errno = 0, sizeof(char));

    struct dirent* entry;
    int len;

    // Find the class name.
    while ((entry = readdir(dir)))
        if ((len = strlen(entry->d_name) - 6) && !strcmp(entry->d_name + len, ".class"))
            if (!class[0] || strcmp(entry->d_name, class))
                memcpy(class, entry->d_name, len);

    if (errno)
        exit(EXIT_FAILURE);

    if (closedir(dir))
        ERR_AND_EXIT("closedir");

    int jargc = argc + 1;
    char** jargv = malloc(jargc * sizeof(char*));
    jargv[0] = (char*) java;
    jargv[1] = class;
    memcpy(&jargv[2], &argv[2], (argc - 2) * sizeof(char*));
    jargv[jargc - 1] = NULL;

    execv(java, jargv);
    ERR_AND_EXIT("execv");
}
