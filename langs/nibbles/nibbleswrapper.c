#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execl("/usr/local/bin/nibbles", "nibbles", "-v", NULL);
        ERR_AND_EXIT("execl");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    // Forgot something here. Pfff not now
    
    FILE* fp;

    if (!(fp = fopen("code.nbl", "w")))
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
        execl("/usr/local/bin/nibbles", "nibbles", "-hs", "code.nbl", NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        return 1;
    
    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);
    
    if (remove("code.nbl"))
        ERR_AND_EXIT("remove");
    
    int nargc = argc + 2;
    char** nargv = malloc(nargc * sizeof(char*));
    nargv[0] = "/usr/local/bin/runghc";
    nargv[1] = "out.hs";
    nargv[2] = "--";
    memcpy(&nargv[3], &argv[2], (argc - 2) * sizeof(char*));
    nargv[nargc - 1] = NULL;

    execv("/usr/local/bin/runghc", nargv);
    ERR_AND_EXIT("execv");
}
