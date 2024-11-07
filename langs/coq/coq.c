#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

int main(int argc, char* argv[]) {
    if (argc && !strcmp(argv[1], "--version")) {
        execv("/usr/bin/coqc", argv);
        perror("execv");
        return 0;
    }

    if (chdir("/tmp"))
        return -1;

    FILE* iof = fopen("code.v", "w");

    if (!iof)
        return -1;
    
    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)))
        if (fwrite(buffer, sizeof(char), nbytes, iof) != (size_t) nbytes)
            return -1;
    
    if (fclose(iof))
        return -1;
    
    //if (execl("/usr/bin/coqc", "coqc", "-o", "code.vo", "code.v", NULL)) {
    //    perror("An error occurred when compiling the source file.");
    //    return -1;
    //}
    
    execl("/usr/bin/coqc", "coqc", "-output-directory", "/tmp", "code.v", NULL);
    perror("execl");
    return 0;
}
