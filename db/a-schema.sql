CREATE EXTENSION citext;
CREATE EXTENSION pgcrypto;  -- For GEN_RANDOM_UUID(), not needed under PG13.

CREATE TYPE hole AS ENUM (
    '12-days-of-christmas', '99-bottles-of-beer', 'abundant-numbers',
    'arabic-to-roman', 'brainfuck', 'christmas-trees', 'css-colors', 'cubes',
    'diamonds', 'divisors', 'emirp-numbers', 'evil-numbers', 'fibonacci',
    'fizz-buzz', 'happy-numbers', 'leap-years', 'lucky-tickets',
    'morse-decoder', 'morse-encoder', 'niven-numbers', 'odious-numbers',
    'ordinal-numbers', 'pangram-grep', 'pascals-triangle',
    'pernicious-numbers', 'poker', 'prime-numbers', 'quine',
    'rock-paper-scissors-spock-lizard', 'roman-to-arabic', 'rule-110',
    'seven-segment', 'sierpiński-triangle', 'spelling-numbers', 'sudoku',
    'ten-pin-bowling', 'united-states', 'λ', 'π', 'τ', 'φ', '√2', '𝑒'
);

CREATE TYPE lang AS ENUM (
    'bash', 'brainfuck', 'c', 'c-sharp', 'cobol', 'f-sharp', 'fortran', 'go',
    'haskell', 'j', 'java', 'javascript', 'julia', 'lisp', 'lua', 'nim',
    'perl', 'php', 'powershell', 'python', 'raku', 'ruby', 'rust', 'swift'
);

CREATE TYPE trophy AS ENUM (
    'caffeinated', 'elephpant-in-the-room', 'happy-birthday-code-golf',
    'hello-world', 'inception', 'independence-day', 'interview-ready',
    'its-over-9000', 'my-god-its-full-of-stars', 'ouroboros',
    'patches-welcome', 'pi-day', 'polyglot', 'slowcoach', 'tim-toady',
    'the-watering-hole'
);

CREATE TABLE ideas (
    id          int  NOT NULL PRIMARY KEY,
    thumbs_down int  NOT NULL,
    thumbs_up   int  NOT NULL,
    title       text NOT NULL UNIQUE
);

CREATE TABLE users (
    id      int    NOT NULL PRIMARY KEY,
    admin   bool   NOT NULL DEFAULT false,
    sponsor bool   NOT NULL DEFAULT false,
    login   citext NOT NULL UNIQUE
);

CREATE TABLE sessions (
    id        uuid      NOT NULL DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    last_used timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id   int       NOT NULL REFERENCES users(id)
);

CREATE TABLE solutions (
    submitted timestamp NOT NULL,
    user_id   int       NOT NULL REFERENCES users(id),
    hole      hole      NOT NULL,
    lang      lang      NOT NULL,
    code      text      NOT NULL,
    failing   bool      NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, hole, lang)
);

CREATE TABLE trophies (
    earned  timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id int       NOT NULL REFERENCES users(id),
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
     AND c.relname IN ('ideas', 'sessions', 'solutions', 'trophies', 'users')
ORDER BY c.relname, t.typlen DESC, t.typname, a.attname;

CREATE VIEW points AS WITH ranked AS (
    SELECT user_id,
           RANK()   OVER (PARTITION BY hole ORDER BY MIN(LENGTH(code))),
           COUNT(*) OVER (PARTITION BY hole)
      FROM solutions
     WHERE NOT failing
  GROUP BY hole, user_id
) SELECT user_id,
         SUM(ROUND(((count - rank) + 1) * (1000.0 / count))) points
    FROM ranked
GROUP BY user_id;

-- Used by /stats
CREATE INDEX solutions_hole_idx ON solutions(hole, user_id) WHERE NOT failing;
CREATE INDEX solutions_lang_idx ON solutions(lang, user_id) WHERE NOT failing;

CREATE ROLE "code-golf" WITH LOGIN;

GRANT SELECT, INSERT, TRUNCATE       ON TABLE ideas     TO "code-golf";
GRANT SELECT                         ON TABLE points    TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE sessions  TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE solutions TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE trophies  TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE users     TO "code-golf";
