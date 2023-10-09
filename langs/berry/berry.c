#include "berry.h"
#include "be_vm.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static int handle_result(bvm *vm, int res) {
    switch (res) {
    case BE_OK: /* everything is OK */
        return be_pcall(vm, 0);
    case BE_EXCEPTION: /* uncatched exception */
        be_dumpexcept(vm);
        return 1;
    case BE_EXIT: /* return exit code */
        return be_toindex(vm, -1);
    case BE_IO_ERROR:
        be_writestring("error: "); 
        be_writestring(be_tostring(vm, -1));
        be_writenewline();
        return -2;
    case BE_MALLOC_FAIL:
        be_writestring("error: memory allocation failed.\n");
        return -1;
    default: /* unkonw result */
        return 2;
    }
}

static void push_args(bvm *vm, int argc, char *argv[]) {
    be_newobject(vm, "list");
    while (argc--) {
        be_pushstring(vm, *argv++);
        be_data_push(vm, -2);
        be_pop(vm, 1);
    }
    be_pop(vm, 1);
    be_setglobal(vm, "_argv");
    be_pop(vm, 1);
}

int main(int argc, char *argv[]) {
    setvbuf(stdin, NULL, _IONBF, 0);
    setvbuf(stdout, NULL, _IONBF, 0);

    if (argc > 1 && strcmp(argv[1], "-v") == 0) {
        be_writestring(BERRY_VERSION "\n");
        return 0;
    }

    int res;
    bvm *vm = be_vm_new();
    be_module_path_set(vm, "/usr/local/lib/berry/packages");
    push_args(vm, argc - 1, argv + 1);

    // Slurp stdin.
    char *code = 0;
    size_t len = 0, size = 128;
    while (!feof(stdin)){
        code = realloc(code, size *= 2);
        len += fread(&code[len], 1, size - len - 1, stdin);
    }
    code[len] = '\0';

    res = be_loadstring(vm, code);
    res = handle_result(vm, res);

    be_vm_delete(vm);
    return res;
}
