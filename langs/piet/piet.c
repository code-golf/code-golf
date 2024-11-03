#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main() {
    char buffer[4096];
    size_t nbytes;
    FILE* iof;

    if (!(iof = fopen("/tmp/code.txt", "w")))
        return -1;
    
    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), stdin)))
        if (fwrite(buffer, sizeof(char), nbytes, iof) != nbytes)
            return -1;

    if (fclose(iof))
        return -1;

    execl("/usr/bin/python", "python", "/usr/bin/ascii2piet.py", "/tmp/code.txt", "/tmp/code.png", NULL);
    perror("execl");

    return 0;
}
