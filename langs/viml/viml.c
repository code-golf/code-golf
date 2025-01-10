#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

const char* viml = "/usr/bin/vim", *code = "code.vim";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(viml, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen(code, "w")))
        ERR_AND_EXIT("fopen");

    if (!fprintf(fp, "let args = ["))
        ERR_AND_EXIT("fprintf");

    for (int i = 2; i < argc; i++) {
        if (i > 2)
            if (!fprintf(fp, ", "))
                ERR_AND_EXIT("fprintf");

        if (!fputc('"', fp))
            ERR_AND_EXIT("fputc");

        for (const char* c = argv[i]; *c; c++)
            switch (*c) {
                case '\\':
                    if (fputs("\\\\", fp))
                        ERR_AND_EXIT("fputs");
                    break;
                case '\"':
                    if (fputs("\\\"", fp))
                        ERR_AND_EXIT("fputs");
                    break;
                case '\n':
                    if (fputs("\\n", fp))
                        ERR_AND_EXIT("fputs");
                    break;
                default:
                    if (!fputc(*c, fp))
                        ERR_AND_EXIT("fputc");
            }

        if (!fputc('"', fp))
            ERR_AND_EXIT("fputc");
    }

    if (fputs("]\n", fp))
        ERR_AND_EXIT("fputs");

    char stderr_buf[4096], stdout_buf[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, stdout_buf, sizeof(stdout_buf))))
        if (fwrite(stdout_buf, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    if (!(fp = fopen("output", "w")))
        ERR_AND_EXIT("fopen");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    if (!(fp = popen("/usr/bin/vim --clean --not-a-term --noplugin -eZ -V1 +0 -S code.vim output 2>&1", "r")))
        ERR_AND_EXIT("popen");

    while (fgets(stdout_buf, sizeof(stdout_buf), fp)) {
        strncat(stderr_buf, stdout_buf, sizeof(stderr_buf) - strlen(stderr_buf) - sizeof(char));
        nbytes += strlen(stdout_buf);
    }

    int status;

    if (!(status = pclose(fp))) {
        if (!(fp = fopen("output", "r")))
            ERR_AND_EXIT("fopen");

        while ((nbytes = fread(stdout_buf, sizeof(char), sizeof(stdout_buf), fp)))
            if (fwrite(stdout_buf, sizeof(char), nbytes, stdout) != (size_t) nbytes)
                ERR_AND_EXIT("fwrite");

        if (fclose(fp))
            ERR_AND_EXIT("fclose");
    } else
        if (nbytes)
            if (fputs(stderr_buf, stderr))
                ERR_AND_EXIT("fputs");
}
