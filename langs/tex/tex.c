#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#undef RAND_MAX

#define RAND_MAX __UINT16_MAX__
#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

char* each_join(char* arr[], int cnt, const char* sep);
const char* sanitize_arg(const char* str);

const char* tex = "/usr/local/bin/tex";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(tex, argv);
        ERR_AND_EXIT("execv");
    }

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    char* code = argv[1];
    
    argc -= 2;
    
    // Now `argv[0]` is the executable name, `argv[1]` is the user code,
    // and `argv[2]` to `argv[2 + (argv - 1)]` are the expected args. 

    char** sanitized_argv;
    if (!(sanitized_argv = malloc(argc * sizeof(char**))))
        ERR_AND_EXIT("malloc");

    for (int i = 0; i < argc; i++) {
        sanitized_argv[i] = sanitize_arg(argv[i + 2]);
    }

    char init[RAND_MAX+1], body[RAND_MAX+1];

    if (!snprintf(init, sizeof(init), "\\def\\argc{%d}", argc))
        ERR_AND_EXIT("snprintf");

    if (argc == 0)
        strcat(init, "\\global\\def\\argv#1{}");
    else {
        // To pass in values, we need to escape some special characters.
        // All the special characters \{}%&#^_%~ are set to catcode 12 ("other").
        // Also set the whitespace characters to catcode 12, unlike \dospecials.
        // But we still need access to some to finish the definition, so we use
        //   \x01 = begin group (previously {)
        //   \x02 = end group (previously })
        //   \x03 = superscript (previously ^)
        //   \x04 = space
        //   \x05 = escape (previously \)
        //   \x06 = parameter (previously #)
        //   \x07 = comment (previously %)
        // The end of the group resets all the catcodes for future tokenizing, but it does not change the catcodes of the tokens already created.
        strcat(init, "{\\catcode1=1\\catcode2=2\\catcode3=7\\catcode5=0\\catcode6=6\\catcode4=10\\catcode7=14\\catcode`$=12\\catcode`&=12\\catcode`^=12\\catcode`_=12\\catcode37=12"
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
        const char* first = sanitized_argv[0];
        const char* rest = each_join(&sanitized_argv[1], argc - 1, "\05or\x04");
        //                         \global   \def   \argv   #1   {   \ifnum   #1=0 [first] \else   \ifcase   #1   \or [rest]  \fi   \fi   }   }
        const char* template = "\x05global\x05def\x05argv\x061\x01\x05ifnum\x061=0\x04%s\x05else\x05ifcase\x061\x05or\x04%s\x05fi\x05fi\x02\x02";
        if (!snprintf(body, sizeof(body), template, first, rest))
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

// Join the length-`cnt` array `arr` of strings with separator `sep`.
char* each_join(char* arr[], int cnt, const char* sep) {
    size_t len = strlen(sep) * (cnt - 1);

    for (int i = 0; i < cnt; i++)
        len += strlen(arr[i]);

    char* rt;

    if (!(rt = malloc(len + 1)))
        ERR_AND_EXIT("malloc");

    for (int i = 0; i < cnt - 1; i++) {
        strcat(rt, arr[i]);
        strcat(rt, sep);
    }
    strcat(rt, arr[cnt - 1]);

    return rt;
}

int count_newlines(const char* str) {
    const char* from_tmp = str;
    const char* next_newline;
    int count;
    for (count = 0; (next_newline = strchr(from_tmp, '\n')); ++count) {
        from_tmp = next_newline + 1;
    }
    return count;
}

const char* replace_newline_with(const char* str, const char* replacement) {
    char from = '\n';

    int count = count_newlines(str);

    if (count == 0) {
        // Argument doesn't contain newline. Just return the original.
        return str;
    }

    int len_replacement = strlen(replacement);
    int final_len = strlen(str) + (len_replacement - 1) * count;
    
    char* result;
    if (!(result = malloc(final_len + 1)))
        ERR_AND_EXIT("malloc");

    char* to_tmp = result;
	const char* from_tmp = str;
    const char* next_newline;
    
    while((next_newline = strchr(from_tmp, from))) {
        int len_before = next_newline - from_tmp;
        to_tmp = strncpy(to_tmp, from_tmp, len_before) + len_before;
        to_tmp = strcpy(to_tmp, replacement) + len_replacement;
        from_tmp = next_newline + 1;
    }
    strcpy(to_tmp, from_tmp);

    return result;
}

const char* sanitize_arg(const char* str) {
    // Most sanitization is already handled by the \catcode stuff.
    // We just need to replace newlines with ^^M to avoid a weird
    // quirk where trailing spaces get removed when they precede a newline.
    // Note we use \x03 instead of ^ since ^ is catcode'd to be "Other".
    return replace_newline_with(str, "\x03\x03M")
}
