CREATE EXTENSION citext;
CREATE EXTENSION pgcrypto;  -- For GEN_RANDOM_UUID(), not needed under PG13.

CREATE TYPE hole AS ENUM (
    '12-days-of-christmas', '99-bottles-of-beer', 'abundant-numbers',
    'arabic-to-roman', 'brainfuck', 'christmas-trees', 'css-colors', 'cubes',
    'diamonds', 'divisors', 'emirp-numbers', 'emojify', 'evil-numbers',
    'fibonacci', 'fizz-buzz', 'happy-numbers', 'intersection', 'leap-years',
    'levenshtein-distance', 'leyland-numbers', 'look-and-say',
    'lucky-tickets', 'morse-decoder', 'morse-encoder', 'niven-numbers',
    'odious-numbers', 'ordinal-numbers', 'pangram-grep', 'pascals-triangle',
    'pernicious-numbers', 'poker', 'prime-numbers', 'quine', 'recamán',
    'rock-paper-scissors-spock-lizard', 'roman-to-arabic', 'rule-110',
    'seven-segment', 'sierpiński-triangle', 'spelling-numbers', 'sudoku',
    'ten-pin-bowling', 'tongue-twisters', 'united-states', 'vampire-numbers',
    'λ', 'π', 'τ', 'φ', '√2', '𝑒'
);

CREATE TYPE keymap AS ENUM ('default', 'vim');

CREATE TYPE lang AS ENUM (
    'bash', 'brainfuck', 'c', 'c-sharp', 'cobol', 'f-sharp', 'fish',
    'fortran', 'go', 'haskell', 'j', 'java', 'javascript', 'julia', 'lisp',
    'lua', 'nim', 'perl', 'php', 'powershell', 'python', 'raku', 'ruby',
    'rust', 'sql', 'swift', 'v', 'zig'
);

CREATE TYPE scoring AS ENUM ('bytes', 'chars');

-- TODO Fix 'tim-toady' & 'the-watering-hole' order when renamed to cheevos.
CREATE TYPE trophy AS ENUM (
    'bakers-dozen', 'bullseye', 'caffeinated', 'cobowl', 'dont-panic',
    'elephpant-in-the-room', 'forty-winks', 'happy-birthday-code-golf',
    'hello-world', 'inception', 'independence-day', 'interview-ready',
    'its-over-9000', 'my-god-its-full-of-stars', 'ouroboros',
    'patches-welcome', 'pi-day', 'polyglot', 'rtfm', 'slowcoach', 'tim-toady',
    'the-watering-hole', 'tl-dr', 'twelvetide', 'up-to-eleven', 'vampire-byte'
);

CREATE TABLE code (
    id    int    NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    bytes int    NOT NULL GENERATED ALWAYS AS (octet_length(code)) STORED,
    chars int    NOT NULL GENERATED ALWAYS AS  (char_length(code)) STORED,
    code  text   NOT NULL,
    CHECK (bytes <= 409600), -- 400 KiB
    EXCLUDE USING hash(code WITH =)
);

CREATE TABLE ideas (
    id          int  NOT NULL PRIMARY KEY,
    thumbs_down int  NOT NULL,
    thumbs_up   int  NOT NULL,
    title       text NOT NULL UNIQUE
);

CREATE TABLE users (
    id           int       NOT NULL PRIMARY KEY,
    admin        bool      NOT NULL DEFAULT false,
    sponsor      bool      NOT NULL DEFAULT false,
    login        citext    NOT NULL UNIQUE,
    time_zone    text,
    delete       timestamp,
    keymap       keymap    NOT NULL DEFAULT 'default',
    country      char(2),
    show_country bool      NOT NULL DEFAULT false,
    started      timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW())
);

CREATE TABLE sessions (
    id        uuid      NOT NULL DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    last_used timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id   int       NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE solutions (
    submitted timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    code_id   int       NOT NULL REFERENCES code(id),
    user_id   int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hole      hole      NOT NULL,
    lang      lang      NOT NULL,
    scoring   scoring   NOT NULL,
    failing   bool      NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, hole, lang, scoring)
);

CREATE TABLE trophies (
    earned  timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trophy  trophy    NOT NULL,
    PRIMARY KEY (user_id, trophy)
);

-- Check the tables are structured optimally.
-- https://www.2ndquadrant.com/en/blog/on-rocks-and-sand/
  SELECT c.relname, a.attname, t.typname, t.typalign, t.typlen
    FROM pg_attribute a
    JOIN pg_class     c ON a.attrelid = c.oid
    JOIN pg_type      t ON a.atttypid = t.oid
   WHERE a.attnum >= 0
     AND c.relname IN ('code', 'ideas', 'sessions', 'solutions', 'trophies', 'users')
ORDER BY c.relname, t.typlen DESC, t.typname, a.attname;

CREATE MATERIALIZED VIEW medals AS WITH ranks AS (
    SELECT user_id,
           RANK() OVER (PARTITION BY hole, lang ORDER BY bytes)
      FROM solutions
      JOIN code ON code_id = id
     WHERE NOT failing AND scoring = 'bytes'
     UNION ALL
    SELECT user_id,
           RANK() OVER (PARTITION BY hole, lang ORDER BY chars)
      FROM solutions
      JOIN code ON code_id = id
     WHERE NOT failing AND scoring = 'chars'
), medals AS (
    SELECT user_id,
           COUNT(*) FILTER (WHERE rank = 1) gold,
           COUNT(*) FILTER (WHERE rank = 2) silver,
           COUNT(*) FILTER (WHERE rank = 3) bronze
      FROM ranks
  GROUP BY user_id
) SELECT * FROM medals WHERE gold > 0 OR silver > 0 OR bronze > 0;

CREATE VIEW bytes_points AS WITH ranked AS (
    SELECT user_id,
           RANK()   OVER (PARTITION BY hole ORDER BY MIN(bytes)),
           COUNT(*) OVER (PARTITION BY hole)
      FROM solutions
      JOIN code ON code_id = id
     WHERE NOT failing
       AND scoring = 'bytes'
  GROUP BY hole, user_id
) SELECT user_id,
         SUM(ROUND(((count - rank) + 1) * (1000.0 / count))) bytes_points
    FROM ranked
GROUP BY user_id;

CREATE VIEW chars_points AS WITH ranked AS (
    SELECT user_id,
           RANK()   OVER (PARTITION BY hole ORDER BY MIN(chars)),
           COUNT(*) OVER (PARTITION BY hole)
      FROM solutions
      JOIN code ON code_id = id
     WHERE NOT failing
       AND scoring = 'chars'
  GROUP BY hole, user_id
) SELECT user_id,
         SUM(ROUND(((count - rank) + 1) * (1000.0 / count))) chars_points
    FROM ranked
GROUP BY user_id;

-- Used by delete_orphaned_code()
CREATE INDEX solutions_code_id_key ON solutions(code_id);

-- Used by /golfers
CREATE UNIQUE INDEX medals_user_id_key ON medals(user_id);

-- Used by /stats
CREATE INDEX solutions_hole_key ON solutions(hole, user_id) WHERE NOT failing;
CREATE INDEX solutions_lang_key ON solutions(lang, user_id) WHERE NOT failing;

CREATE ROLE "code-golf" WITH LOGIN;

-- Only owners can refresh.
ALTER MATERIALIZED VIEW medals OWNER TO "code-golf";

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE    code         TO "code-golf";
GRANT SELECT                         ON SEQUENCE code_id_seq  TO "code-golf";
GRANT SELECT, INSERT, TRUNCATE       ON TABLE    ideas        TO "code-golf";
GRANT SELECT                         ON TABLE    bytes_points TO "code-golf";
GRANT SELECT                         ON TABLE    chars_points TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE    sessions     TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE    solutions    TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE    trophies     TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE    users        TO "code-golf";
