CREATE EXTENSION citext;

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
    'bash', 'brainfuck', 'c', 'c-sharp', 'f-sharp', 'fortran', 'go',
    'haskell', 'j', 'java', 'javascript', 'julia', 'lisp', 'lua', 'nim',
    'perl', 'php', 'powershell', 'python', 'raku', 'ruby', 'rust', 'swift'
);

CREATE TYPE trophy AS ENUM (
    'caffeinated', 'elephpant-in-the-room', 'happy-birthday-code-golf',
    'hello-world', 'inception', 'interview-ready', 'its-over-9000',
    'my-god-its-full-of-stars', 'ouroboros', 'patches-welcome', 'polyglot',
    'slowcoach', 'tim-toady', 'the-watering-hole'
);

CREATE TABLE ideas (
    id          integer NOT NULL PRIMARY KEY,
    thumbs_down integer NOT NULL,
    thumbs_up   integer NOT NULL,
    title       text    NOT NULL UNIQUE
);

CREATE TABLE users (
    id      integer               NOT NULL PRIMARY KEY,
    admin   boolean DEFAULT false NOT NULL,
    sponsor boolean DEFAULT false NOT NULL,
    login   citext                NOT NULL UNIQUE
);

CREATE TABLE solutions (
    submitted timestamp without time zone NOT NULL,
    user_id   integer                     NOT NULL REFERENCES users(id),
    hole      hole                        NOT NULL,
    lang      lang                        NOT NULL,
    code      text                        NOT NULL,
    failing   boolean DEFAULT false       NOT NULL,
    PRIMARY KEY (user_id, hole, lang)
);

CREATE TABLE trophies (
    earned  timestamp without time zone NOT NULL,
    user_id integer                     NOT NULL REFERENCES users(id),
    trophy  trophy                      NOT NULL,
    UNIQUE (user_id, trophy)
);

-- Check the tables are structured optimally.
-- https://www.2ndquadrant.com/en/blog/on-rocks-and-sand/
  SELECT c.relname, a.attname, t.typname, t.typalign, t.typlen
    FROM pg_attribute a
    JOIN pg_class     c ON a.attrelid = c.oid
    JOIN pg_type      t ON a.atttypid = t.oid
   WHERE a.attnum >= 0
     AND c.relname IN ('ideas', 'solutions', 'trophies', 'users')
ORDER BY c.relname, t.typlen DESC, t.typname, a.attname;

CREATE VIEW points AS WITH leaderboard AS (
    SELECT DISTINCT ON (hole, user_id)
           hole,
           length(code) strokes,
           user_id
      FROM solutions
     WHERE NOT failing
  ORDER BY hole, user_id, length(code), submitted
), scored_leaderboard AS (
    SELECT hole,
           round(
                (   (   count(*) OVER (PARTITION BY hole)
                        -
                        rank() OVER (PARTITION BY hole ORDER BY strokes)
                    ) + 1
                ) * (1000.0 / count(*) OVER (PARTITION BY hole))
           ) score,
           user_id
      FROM leaderboard l
) SELECT user_id,
         sum(score) points,
         count(*)   holes
    FROM scored_leaderboard
GROUP BY user_id;

CREATE ROLE "code-golf" WITH LOGIN;

GRANT SELECT, INSERT, TRUNCATE       ON TABLE ideas     TO "code-golf";
GRANT SELECT                         ON TABLE points    TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE solutions TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE trophies  TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE users     TO "code-golf";
