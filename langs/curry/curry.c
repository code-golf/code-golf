#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

int main(int argc, char* argv[]) {
    if (argc && !strcmp(argv[1], "--version")) {
        argv[1] = "--numeric-version";
        execv("/out/bin/curry", argv);
        perror("execv");
        return 0;
    }

    FILE* iof = fopen("/tmp/code.curry", "w");

    if (!iof)
        return -1;

    char buffer[4096];
    size_t nbytes;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)))
        if (fwrite(buffer, sizeof(char), nbytes, iof) != nbytes)
            return -1;
    
    if (fclose(iof))
        return -1;
    
    execv("/usr/bin/runcurry", argv);
    perror("execv");
    return 0;
}
