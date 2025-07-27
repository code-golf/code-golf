#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#undef RAND_MAX

#define RAND_MAX __UINT16_MAX__
#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

char* each_join(char* arr[], int cnt, const char* sep);

const char* tex = "/usr/local/bin/tex";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(tex, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    char* code = argv[1];

    memmove(&argv[1], &argv[2], argc-- * sizeof(char*));

    char init[RAND_MAX+1], body[RAND_MAX+1];

    if (!snprintf(init, sizeof(init), "\\def\\argc{%d}", --argc))
        ERR_AND_EXIT("snprintf");

    if (argc == 0)
        strcat(init, "\\global\\def\\argv#1{}");
    else {
        // To pass in values, we need to escape some special characters.
        // All the special characters \{}%&#^_%~ are set to catcode 12 ("other").
        // Also set the whitespace characters to catcode 12, unlike \dospecials.
        // But we still need access to some to finish the definition, so we use
        //   byte 1 = begin group (previously {)
        //   byte 2 = end group (previously })
        //   byte 4 = space
        //   byte 5 = escape (previously \)
        //   byte 6 = parameter (previously #)
        //   byte 7 = comment (previously %)
        // The end of the group resets all the catcodes for future tokenizing, but it does not change the catcodes of the tokens already created.
        strcat(init, "{\\catcode1=1\\catcode2=2\\catcode5=0\\catcode6=6\\catcode4=10\\catcode7=14\\catcode`$=12\\catcode`&=12\\catcode`^=12\\catcode`_=12\\catcode37=12"
            "\\catcode`~=12\\catcode`#=12\\catcode`{=12\\catcode`}=12\\catcode9=12\\catcode32=12\\catcode10=12\\catcode12=12\\catcode13=13\\catcode92=12");
        // The following line reads something like
        //   \global\def\argv#1{
        //     \ifnum#1=0
        //       arg0
        //     \else
        //       \ifcase#1 \or arg1 \or arg2 \or ... \or argN \fi
        //     \fi
        //   }
        // we handle 0 separately to avoid it starting with a space:
        //   \ifcase#1arg0\or...    doesn't read the index right
        //   \ifcase#1 arg0\or...   starts with an space
        //   \ifcase#1{}arg0\or...  starts with an empty group (invisible but brutal)
        //   \ifcase#1\relax arg0\or...  starts with a \relax (invisible but brutal)
        // no easy way directly in the \ifcase, so use the extra \ifnum wrapper
        //   \ifnum#1=0 arg0\else\ifcase\or arg1\or arg2\fi\fi
        // ensures that arg0 is preceded by a digit (which removes the space), or a macro which doesn't get expanded (which removes the space).
        //
        // Note: If you ever come back to this, make sure it works for argv{\the\i}, argv\i, and argv{\count0}, and hole args starting with digits.
        if (!snprintf(body, sizeof(body), "globaldefargv1ifnum1=0%selseifcase1or%sfifi", argv[1], each_join(&argv[1], argc, "or")))
            ERR_AND_EXIT("snprintf");

        strcat(init, body);
    }
    // \octet enables the octet font.
    // \footline={} disables page numbers.
    // \parindent=0pt prevents per-paragraph indentation.
    // \hsize and \vsize set the page dimensions. I set them a bit less than the maximum legal dimension which is less than 16384pt.
    // \bye closes the document (TeX doesn't handle EOF how you might expect).
    if (!snprintf(body, sizeof(body), "\\octet\\footline={}\\parindent=0pt\\hsize=16000pt\\vsize=16000pt\\relax\n%s\n%s\n\\bye", init, code))
        ERR_AND_EXIT("snprintf");

    srand((unsigned) time(NULL));

    char file[64];
    // The randoms might not actually prevent reading the file, but they do certainly help avoid grepping out the wrong line.
    if (!snprintf(file, sizeof(file), "solution_%u_%u_%u_%u", rand() % RAND_MAX, rand() % RAND_MAX, rand() % RAND_MAX, rand() % RAND_MAX))
        ERR_AND_EXIT("snprintf");

    char src[128];

    if (!snprintf(src, sizeof(src), "%s.tex", file))
        ERR_AND_EXIT("snprintf");

    FILE* fp;

    if (!(fp = fopen(src, "w")))
        ERR_AND_EXIT("fopen");

    if (fputs(body, fp)) {
        if (fclose(fp))
            ERR_AND_EXIT("fclose");

        ERR_AND_EXIT("fputs");
    }

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    char cmd[256];

    if (!snprintf(cmd, sizeof(cmd), "%s %s > /dev/null", tex, src))
        ERR_AND_EXIT("snprintf");

    system(cmd); // Here we omit `ERR_AND_EXIT` to prevent the program from terminating as soon as /usr/local/bin/tex returns an error.

    char log[128];

    if (!snprintf(log, sizeof(log), "%s.log", file))
        ERR_AND_EXIT("snprintf");

    if (!(fp = fopen(log, "r")))
        ERR_AND_EXIT("fopen");

    char line[4096];

    while (fgets(line, sizeof(line), fp)) {
        if (strstr(line, file) || strstr(line, "Version 3.141592653") || strstr(line, " [1] )") ||
            (strstr(line, "=\\count") && !strstr(line, " =\\count ")))
            continue;

        if (fputs(line, stderr))
            ERR_AND_EXIT("fputs");
    }

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    char dvi[128];

    if (!snprintf(dvi, sizeof(dvi), "%s.dvi", file))
        ERR_AND_EXIT("snprintf");

    if (!access(dvi, F_OK)) {
        if (!snprintf(cmd, sizeof(cmd), "dvi-to-text %s", dvi))
            ERR_AND_EXIT("snprintf");

        if (system(cmd))
            ERR_AND_EXIT("system");
    }
}

// Based on the original Bash implementation:
// https://stackoverflow.com/a/17841619/7481517
char* each_join(char* arr[], int cnt, const char* sep) {
    size_t len;

    for (int i = 0; i < cnt; i++)
        len += strlen(arr[i]) + strlen(sep);

    char* rt;

    if (!(rt = malloc(len + 1)))
        ERR_AND_EXIT("malloc");

    for (int i = 1; i < cnt; i++) {
        strcat(rt, arr[i]);
        strcat(rt, sep);
    }

    return rt;
}
