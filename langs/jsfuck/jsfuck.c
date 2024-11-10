#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

int main(int argc, char* argv[]) {
    // FIXME Avoid keyboard interruption from never-terminating containers
    if (argc > 1 && !strcmp(argv[1], "--version"))
        exit(EXIT_SUCCESS);

    FILE* fp;

    if (!(fp = fopen("/tmp/code.js", "w")))
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");
    
    if (fclose(fp))
        ERR_AND_EXIT("fclose");
    
    int nargc = argc + 1;
    char** nargv = malloc(nargc * sizeof(char*));
    nargv[0] = "/usr/bin/node";
    nargv[1] = "/tmp/code.js";
    memcpy(&nargv[2], &argv[2], (argc - 2) * sizeof(char*));
    nargv[nargc - 1] = NULL;
    
    execv("/usr/bin/node", nargv);
    ERR_AND_EXIT("execv");
}
