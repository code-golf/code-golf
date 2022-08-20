#include "sqlite3.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]) {
    if (argc > 1 && strcmp(argv[1], "-v") == 0) {
        printf("SQLite %s\n", sqlite3_libversion());
        return 0;
    }

    sqlite3 *db;
    if (sqlite3_open(":memory:", &db) != SQLITE_OK) {
        fprintf(stderr, "%s\n", sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    char *err;
    if (sqlite3_exec(db, "CREATE TABLE argv(arg TEXT)", 0, 0, &err) != SQLITE_OK ) {
        fprintf(stderr, "%s\n", err);
        sqlite3_free(err);
        sqlite3_close(db);
        return 1;
    }

    sqlite3_stmt *res;
    for (int i = 2; i < argc; i++) {
        int rc = sqlite3_prepare_v2(db, "INSERT INTO argv VALUES(?)", -1, &res, 0);
        sqlite3_bind_text(res, 1, argv[i], -1, 0);
        sqlite3_step(res);
        sqlite3_finalize(res);
    }

    char *sql = 0;
    size_t len = 0, size = 128;
    while (!feof(stdin)){
        sql = realloc(sql, size *= 2);
        len += fread(&sql[len], 1, size - len - 1, stdin);
    }
    sql[len] = '\0';

    const char* statement = sql;

    while(*statement) {
        if (sqlite3_prepare_v2(db, statement, -1, &res, &statement) != SQLITE_OK) {
            fprintf(stderr, "%s\n", sqlite3_errmsg(db));
            sqlite3_close(db);
            return 1;
        }

        while (sqlite3_step(res) == SQLITE_ROW)
            if (sqlite3_column_type(res, 0) != SQLITE_NULL)
                printf("%s\n", sqlite3_column_text(res, 0));

        sqlite3_finalize(res);
    }

    sqlite3_close(db);
    return 0;
}
