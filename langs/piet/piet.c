#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc __attribute__((unused)), char* argv[] __attribute__((unused))) {
    if (chdir("/tmp"))
        return -1;
    
    FILE* iof = fopen("code.txt", "w");

    if (!iof)
        return -1;
    
    char buffer[4096];
    ssize_t nbytes;
    
    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)))
        if (fwrite(buffer, sizeof(char), nbytes, iof) != (size_t) nbytes)
            return -1;

    if (fclose(iof))
        return -1;

    if (execl("/usr/bin/python", "python", "/ascii2piet.py", "code.txt", "code.png", NULL)) {
        perror("An error occurred when generating the PNG image.");
        return -1;
    }
    
    execl("/usr/bin/npiet", "npiet", "-d", "code.png", NULL); // npiet takes as input a portable pixmap (a ppm file) and since v0.3a png and gif files too. (Source: https://www.bertnase.de/npiet/)
    perror("execl");
    return 0;
}
