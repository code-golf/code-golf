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
        "size_t i;\n",
        fp);

    // char in[] = "arg2\0arg3\0argN\0";
    fputs("char in[] = \"", fp);
    for (int i = 2; i < argc; i++) {
        // Write each character of argv[i] using \xNN escape
        for (char *a = argv[i]; *a; a++)
            fprintf(fp, "\\x%02x", *a & 255);
        fputs("\\0",fp);
    }
    fputs("\";\n", fp);

    int ch;
    while ((ch = getchar()) != EOF)
        switch (ch) {
            case '>':
                fputs("++p;\n", fp);
                break;
            case '<':
                fputs("--p;\n", fp);
                break;
            case '+':
                fputs("++t[p];\n", fp);
                break;
            case '-':
                fputs("--t[p];\n", fp);
                break;
            case ',':
                // Cell is unchanged when trying to read past the end of stdin
                fputs("if(i<sizeof(in)-1)t[p]=in[i++];\n", fp);
                break;
            case '.':
                fputs("putchar(t[p]);\n", fp);
                break;
            case '[':
                fputs("while(t[p]){\n", fp);
                break;
            case ']':
                fputs("}\n", fp);
                break;
        }

    fputs("return 0;}", fp);

    fclose(fp);

    char *tcc_argv[] = {"/usr/bin/tcc", "-run", "/tmp/code.c", NULL};
    execv("/usr/bin/tcc", tcc_argv);
    perror("execv");
}