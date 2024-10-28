#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

const char* java = "/opt/jdk/bin/java";
const char* javac = "/opt/jdk/bin/javac";
const char* code = "/tmp/code.java";

int main(int argc, char* argv[]) {
    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        execv(java, argv);
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
        execl(javac, javac, code, NULL);
        perror("execl");
        return 3;
    }

    int status;
    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 4;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if(remove(code)) {
        perror("Error deleting file");
        return 5;
    }

    // Find the class name.
    DIR* dir = opendir("/tmp/");

    if (!dir)
        return 6;

    char class[256];
    memset(class, 0, sizeof(class));

    errno = 0;

    struct dirent* d;

    while (d = readdir(dir)) {
        int len = strlen(d->d_name) - 6;
        if (len > 0 && strcmp(d->d_name + len, ".class") == 0) {
            if (!class[0] || strcmp(d->d_name, class) < 0)
                memcpy(class, d->d_name, len);
        }
    }

    if (errno)
        return 7;

    if (closedir(dir))
        return 8;

    if(chdir("/tmp")) {
        perror("Error changing directory");
        return 9;
    }

    int jargc = argc + 1;
    char** jargv = malloc(jargc * sizeof(char*));
    jargv[0] = (char*)java;
    jargv[1] = class;
    memcpy(&jargv[2], &argv[2], (argc - 2) * sizeof(char*));
    jargv[jargc - 1] = NULL;

    execv(java, jargv);
    perror("execv");
    return 10;
}
