#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl("/usr/bin/yue", "yue", "-v", NULL);
        ERR_AND_EXIT("execl");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");
    
    FILE* fp;

    if (!(fp = fopen("code.yue", "w")))
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
        if (!dup2(STDERR_FILENO, STDOUT_FILENO))
            ERR_AND_EXIT("dup2");
        
        execl("/usr/bin/yue", "yue", "code.yue", NULL);
        ERR_AND_EXIT("execl");
    }
    
    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 1;
    
    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);
    
    if (remove("code.yue"))
        ERR_AND_EXIT("remove");
    
    int yargc = argc + 1;
    char** yargv = malloc(yargc * sizeof(char*));
    yargv[0] = "/usr/bin/lua";
    yargv[1] = "code.lua";
    memcpy(&yargv[2], &argv[2], (argc - 2) * sizeof(char*));
    yargv[yargc - 1] = NULL;

    execv("/usr/bin/lua", yargv);
    ERR_AND_EXIT("execv");
}
