#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

// Due to compiling with `-Werror` have both params been marked as "unused", for now.
int main(int argc __attribute__((unused)), char* argv[] __attribute__((unused))) {
    if (chdir("/tmp")) {
        perror("An error occurred when changing directory.");
        return -1;
    }
    
    FILE* iof = fopen("code.txt", "w");

    if (!iof) {
        perror("An error occurred when opening the source file for writing.");
        return -1;
    }
    
    char buffer[4096];
    ssize_t nbytes;
    
    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)))
        if (fwrite(buffer, sizeof(char), nbytes, iof) != (size_t) nbytes)
            return -1;

    if (fclose(iof)) {
        perror("An error occurred when closing the source file.");
        return -1;
    }

    if (execl("/usr/bin/python", "python", "/ascii2piet.py", "code.txt", "code.png", NULL)) {
        perror("An error occurred when saving the output file.");
        return -1;
    }
    
    if (execl("/usr/bin/npiet", "npiet", "code.png", NULL)) {
        perror("An error occurred when reading the input file.");
        return -1;
    }

    return 0;
}
