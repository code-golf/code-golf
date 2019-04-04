#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

// Adapted from https://github.com/TryItOnline/brainfuck

int main(int argc, char *argv[]) {
    FILE *fp = fopen("/tmp/code.c", "w");

    if (fp == 0)
        return 1;

    fputs(
        "#include <stdint.h>\n"
        "#include <stdio.h>\n"
        "int main(){\n"
        "uint8_t t[65536] = {0};\n"
        "uint16_t p = 0;\n"
        "int i;",
        fp);

    int ch;
    while ((ch = getchar()) != EOF)
        switch (ch) {
            case '>':
                fputs("++p;", fp);
                break;
            case '<':
                fputs("--p;", fp);
                break;
            case '+':
                fputs("++t[p];", fp);
                break;
            case '-':
                fputs("--t[p];", fp);
                break;
            case ',':
                // Not yet implemented
                break;
            case '.':
                fputs("putchar(t[p]);", fp);
                break;
            case '[':
                fputs("while(t[p]){", fp);
                break;
            case ']':
                fputs("}", fp);
                break;
        }

    fputs("return 0;}", fp);

    fclose(fp);

    char *tcc_argv[] = {"/usr/bin/tcc", "-run", "/tmp/code.c", NULL};
    execv("/usr/bin/tcc", tcc_argv);
    perror("execv");
}