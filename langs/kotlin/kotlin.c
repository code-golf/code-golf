#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

const char* kotlinc = "/kotlinc/bin/kotlinc";
const char* kotlin = "/kotlinc/bin/kotlin";
const char* code = "/tmp/code.kt";

int main(int argc, char* argv[]) {
    if (setenv("LC_ALL", "C.UTF-8", 1) ||
        setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/jdk/openjdk/bin", 1)) {
        perror("Error setting environment variable");
        return 1;
    }

    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        argv[1] = "-version";
        execv(kotlin, argv);
        perror("execv");
        return 0;
    }

    FILE* fp = fopen(code, "w");

    if (!fp)
        return 2;

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != nbytes)
            return 3;

    fclose(fp);


    if(chdir("/tmp")) {
        perror("Error changing directory");
        return 10;
    }

    pid_t pid = fork();
    if (!pid) {
        execl(kotlinc, kotlinc, code, NULL);
        perror("execl");
        return 4;
    }

    int status;
    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 5;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if(remove(code)) {
        perror("Error deleting file");
        return 6;
    }

    // Determine name of the main class file.
    DIR* dir = opendir("/tmp/");

    if (!dir)
        return 7;

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
        return 8;

    if (closedir(dir))
        return 9;

    int kargc = argc + 2;
    char** kargv = malloc(kargc * sizeof(char*));
    kargv[0] = kotlin;
    kargv[1] = class;
    memcpy(&kargv[2], &argv[2], (argc - 2) * sizeof(char*));
    kargv[kargc - 1] = NULL;

    execv(kotlin, kargv);
    perror("execv");
    return 11;
}
