#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

const char* clang = "/usr/bin/clang";
const char* buffer_fix = "/tmp/buffer_fix.cpp";
const char* code = "/tmp/code.cpp";
const char* bin = "/tmp/code";

int main(int argc, char* argv[]) {
    if (argc > 1 && strcmp(argv[1], "--version") == 0) {
        execv(clang, argv);
        perror("execv");
        return 0;
    }

    FILE* fp = fopen(code, "w");

    if (!fp)
        return 1;

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)) > 0)
        if (fwrite(buffer, sizeof(char), nbytes, fp) != nbytes)
            return 2;

    fclose(fp);

    // Before main runs, set stdout and stderr to unbuffered
    fp = fopen(buffer_fix, "w");
    fputs(
        "#include <stdio.h>\n"
        "static int success_stdout = setvbuf(stdout, NULL, _IONBF, 0);\n"
        "static int success_stderr = setvbuf(stderr, NULL, _IONBF, 0);\n",
        fp);
    fclose(fp);

    pid_t pid = fork();
    if (!pid) {
        // See https://clang.llvm.org/cxx_status.html for valid -std values.
        execl(clang, clang, "-std=c++2b", "-target", "x86_64-alpine-linux-musl", "-O2", "-lstdc++",
            "-fcolor-diagnostics", "-I/usr/include/c++/10.3.1/",
            "-I/usr/include/c++/10.3.1/x86_64-alpine-linux-musl/",
            "-I/usr/include/c++/10.3.1/backward/", "-o", bin, code, buffer_fix, NULL);
        perror("execl");
        return 3;
    }

    int status;
    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 4;

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    if (remove(code) || remove(buffer_fix)) {
        perror("Error deleting file");
        return 5;
    }

    int cargc = argc;
    char** cargv = malloc(cargc * sizeof(char*));
    cargv[0] = (char*)bin;
    memcpy(&cargv[1], &argv[2], (argc - 2) * sizeof(char*));
    cargv[cargc - 1] = NULL;

    execv(bin, cargv);
    perror("execv");
    return 6;
}
