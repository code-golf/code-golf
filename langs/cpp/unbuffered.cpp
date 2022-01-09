#include <stdio.h>

// Before main runs, set stdout and stderr to unbuffered
static int success_stdout = setvbuf(stdout, NULL, _IONBF, 0);
static int success_stderr = setvbuf(stderr, NULL, _IONBF, 0);
