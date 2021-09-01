VACUUM ANALYZE;

REFRESH MATERIALIZED VIEW medals;
REFRESH MATERIALIZED VIEW rankings;
REFRESH MATERIALIZED VIEW points;

-- Check the tables are structured optimally.
-- https://www.2ndquadrant.com/en/blog/on-rocks-and-sand/
  SELECT c.relname, a.attname, t.typname, t.typalign, t.typlen
    FROM pg_attribute a
    JOIN pg_class     c ON c.oid = a.attrelid
    JOIN pg_namespace n ON n.oid = c.relnamespace
    JOIN pg_type      t ON t.oid = a.atttypid
   WHERE a.attnum >= 0
     AND c.relkind IN ('m', 'r')
     AND n.nspname = 'public'
ORDER BY c.relname, t.typlen DESC, t.typname, a.attname;

-- Table sizes.
WITH tables AS (
    SELECT relname                                            "table",
           reltuples                                          "rows",
           pg_indexes_size(c.oid)                             index_size,
           COALESCE(pg_total_relation_size(reltoastrelid), 0) toast_size,
           pg_total_relation_size(c.oid)                      total_size
      FROM pg_class     c
      JOIN pg_namespace n ON relnamespace = n.oid
     WHERE nspname = 'public' AND relkind IN ('m', 'r')
) SELECT "table", "rows",
         pg_size_pretty(total_size - index_size - toast_size) table_size,
         pg_size_pretty(index_size)                           index_size,
         pg_size_pretty(toast_size)                           toast_size,
         pg_size_pretty(total_size)                           total_size
    FROM tables
ORDER BY "table";
