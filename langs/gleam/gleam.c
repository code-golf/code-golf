#include <dirent.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <unistd.h>

#define ERR_AND_EXIT(msg) do { perror(msg); exit(EXIT_FAILURE); } while (0)

void copy_file(const char* src_file, const char* dst_file);
void copy_folder(const char* src_dir, const char* dst_dir);

const char* gleam = "/usr/local/bin/gleam", *code = "src/main.gleam";

int main(int argc, char* argv[]) {
    if (!strcmp(argv[1], "--version")) {
        execv(gleam, argv);
        ERR_AND_EXIT("execv");
    }

    copy_folder("/usr/local", "/tmp");

    if (chdir("/tmp"))
        ERR_AND_EXIT("chdir");

    FILE* fp;

    if (!(fp = fopen(code, "w")))
        ERR_AND_EXIT("fopen");

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = read(STDIN_FILENO, buffer, sizeof(buffer))))
        if (fwrite(buffer, sizeof(char), nbytes, fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(fp))
        ERR_AND_EXIT("fclose");

    pid_t pid;

    if (!(pid = fork())) {
        if (!dup2(STDERR_FILENO, STDOUT_FILENO))
            ERR_AND_EXIT("dup2");

        execl(gleam, gleam, "build", "--no-print-progress", NULL);
        ERR_AND_EXIT("execl");
    }

    int status;

    waitpid(pid, &status, 0);

    if (!WIFEXITED(status))
        exit(EXIT_FAILURE);

    if (WEXITSTATUS(status))
        return WEXITSTATUS(status);

    int gargc = argc + 2;
    char** gargv = malloc(gargc * sizeof(char*));
    gargv[0] = (char*) gleam;
    gargv[1] = "run";
    gargv[2] = "--no-print-progress";
    memcpy(&gargv[3], &argv[2], (argc - 2) * sizeof(char*));
    gargv[gargc - 1] = NULL;

    execv(gleam, gargv);
    ERR_AND_EXIT("execv");
}

void copy_file(const char* src_file, const char* dst_file) {
    FILE* src_fp, *dst_fp;

    if (!(src_fp = fopen(src_file, "r")))
        ERR_AND_EXIT("fopen");

    if (!(dst_fp = fopen(dst_file, "w"))) {
        if (fclose(src_fp))
            ERR_AND_EXIT("fclose");

        ERR_AND_EXIT("fopen");
    }

    char buffer[4096];
    ssize_t nbytes;

    while ((nbytes = fread(buffer, sizeof(char), sizeof(buffer), src_fp)))
        if (fwrite(buffer, sizeof(char), nbytes, dst_fp) != (size_t) nbytes)
            ERR_AND_EXIT("fwrite");

    if (fclose(src_fp) || fclose(dst_fp))
        ERR_AND_EXIT("fclose");
}

void copy_folder(const char* src_dir, const char* dst_dir) {
    DIR* dp;

    if (!(dp = opendir(src_dir)))
        ERR_AND_EXIT("opendir");

    struct stat st;

    if (stat(dst_dir, &st))
        if (mkdir(dst_dir, 0755)) {
            if (closedir(dp))
                ERR_AND_EXIT("closedir");

            ERR_AND_EXIT("mkdir");
        }

    struct dirent* entry;
    char src_buf[4096], dst_buf[4096];

    while ((entry = readdir(dp))) {
        if (!strcmp(entry->d_name, ".") || !strcmp(entry->d_name, ".."))
            continue;

        snprintf(src_buf, sizeof(src_buf), "%s/%s", src_dir, entry->d_name);
        snprintf(dst_buf, sizeof(dst_buf), "%s/%s", dst_dir, entry->d_name);

        if (!stat(src_buf, &st) && S_ISDIR(st.st_mode))
            copy_folder(src_buf, dst_buf);
        else if (!stat(src_buf, &st) && S_ISREG(st.st_mode))
            copy_file(src_buf, dst_buf);
    }

    if (closedir(dp))
        ERR_AND_EXIT("closedir");
}
