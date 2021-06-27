CREATE EXTENSION citext;

CREATE TYPE hole AS ENUM (
    '12-days-of-christmas', '99-bottles-of-beer', 'abundant-numbers',
    'arabic-to-roman', 'brainfuck', 'christmas-trees', 'css-colors', 'cubes',
    'diamonds', 'divisors', 'emirp-numbers', 'emojify', 'evil-numbers',
    'fibonacci', 'fizz-buzz', 'happy-numbers', 'intersection',
    'kolakoski-constant', 'kolakoski-sequence', 'leap-years',
    'levenshtein-distance', 'leyland-numbers', 'look-and-say',
    'lucky-tickets', 'morse-decoder', 'morse-encoder', 'niven-numbers',
    'odious-numbers', 'ordinal-numbers', 'pangram-grep', 'pascals-triangle',
    'pernicious-numbers', 'poker', 'prime-numbers', 'quine', 'recamán',
    'rock-paper-scissors-spock-lizard', 'roman-to-arabic', 'rule-110',
    'seven-segment', 'sierpiński-triangle', 'spelling-numbers',
    'star-wars-opening-crawl', 'sudoku', 'sudoku-v2', 'ten-pin-bowling',
    'tongue-twisters', 'united-states', 'vampire-numbers', 'van-eck-sequence',
    'λ', 'π', 'τ', 'φ', '√2', '𝑒'
);

CREATE TYPE keymap AS ENUM ('default', 'vim');

CREATE TYPE lang AS ENUM (
    'assembly', 'bash', 'brainfuck', 'c', 'c-sharp', 'cobol', 'crystal',
    'f-sharp', 'fish', 'fortran', 'go', 'haskell', 'hexagony', 'j', 'java',
    'javascript', 'julia', 'lisp', 'lua', 'nim', 'perl', 'php', 'powershell',
    'python', 'raku', 'ruby', 'rust', 'sql', 'swift', 'v', 'zig'
);

CREATE TYPE medal AS ENUM ('diamond', 'gold', 'silver', 'bronze');

CREATE TYPE scoring AS ENUM ('bytes', 'chars');

CREATE TYPE theme AS ENUM ('auto', 'dark', 'light');

CREATE TYPE cheevo AS ENUM (
    'assembly-required', 'bakers-dozen', 'bullseye', 'caffeinated', 'cobowl',
    'different-strokes', 'dont-panic', 'elephpant-in-the-room',
    'fish-n-chips', 'forty-winks', 'happy-birthday-code-golf', 'hello-world',
    'inception', 'independence-day', 'interview-ready', 'its-over-9000',
    'may-the-4ᵗʰ-be-with-you', 'my-god-its-full-of-stars', 'ouroboros',
    'patches-welcome', 'pi-day', 'polyglot', 'rtfm', 'slowcoach',
    'solve-quine', 'the-watering-hole', 'tim-toady', 'tl-dr', 'twelvetide',
    'up-to-eleven', 'vampire-byte'
);

CREATE TABLE discord_records (
    hole    hole NOT NULL,
    lang    lang NOT NULL,
    message text NOT NULL,
    PRIMARY KEY(hole, lang)
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
    started      timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    referrer_id  int                REFERENCES users(id) ON DELETE SET NULL,
    theme        theme     NOT NULL DEFAULT 'auto',
    CHECK (id != referrer_id)   -- Can't refer yourself
);

CREATE TABLE sessions (
    id        uuid      NOT NULL DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    last_used timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id   int       NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE solutions (
    submitted timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    bytes     int       NOT NULL,
    chars     int,
    user_id   int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hole      hole      NOT NULL,
    lang      lang      NOT NULL,
    scoring   scoring   NOT NULL,
    failing   bool      NOT NULL DEFAULT false,
    code      text      NOT NULL,
    -- Assembly can only be scored on bytes, and they are compiled bytes.
    CHECK ((lang  = 'assembly' AND chars IS NULL AND scoring = 'bytes')
        OR (lang != 'assembly' AND bytes = octet_length(code)
                               AND chars = char_length(code))),
    -- Solutions are limited to 400 KiB, TODO < 128 KiB (not <=).
    CHECK (octet_length(code) <= 409600),
    PRIMARY KEY (user_id, hole, lang, scoring)
);

CREATE TABLE trophies (
    earned  timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trophy  cheevo    NOT NULL,
    PRIMARY KEY (user_id, trophy)
);

CREATE MATERIALIZED VIEW medals AS WITH ranks AS (
    SELECT user_id, hole, lang, scoring,
           RANK() OVER (
               PARTITION BY hole, lang, scoring
                   ORDER BY CASE WHEN scoring = 'bytes'
                                 THEN bytes ELSE chars END
           )
      FROM solutions
     WHERE NOT failing
) SELECT user_id, hole, lang, scoring,
         (enum_range(NULL::medal))[rank + 1] medal
    FROM ranks
   WHERE rank < 4
   UNION ALL
  SELECT MIN(user_id) user_id, hole, lang, scoring, 'diamond'::medal
    FROM ranks
   WHERE rank = 1
GROUP BY hole, lang, scoring
  HAVING COUNT(*) = 1;

CREATE VIEW bytes_points AS WITH ranked AS (
    SELECT user_id,
           RANK()   OVER (PARTITION BY hole ORDER BY MIN(bytes)),
           COUNT(*) OVER (PARTITION BY hole)
      FROM solutions
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
     WHERE NOT failing
       AND scoring = 'chars'
  GROUP BY hole, user_id
) SELECT user_id,
         SUM(ROUND(((count - rank) + 1) * (1000.0 / count))) chars_points
    FROM ranked
GROUP BY user_id;

CREATE MATERIALIZED VIEW rankings AS WITH strokes AS (
    select hole, lang, scoring, user_id, submitted,
           case when scoring = 'bytes' then bytes else chars end strokes
      from solutions
     where not failing
), min as (
    select hole, scoring, min(strokes)::numeric Sa
      from strokes
  group by hole, scoring
), min_per_lang as (
    select hole, lang, scoring, min(strokes)::numeric S, sqrt(count(*)) N
     from strokes
  group by hole, lang, scoring
), bayesian_estimators as (
    select hole, lang, scoring, S,
           (N / (N + 1)) * S + (1 / (N + 1)) * Sa Sb
      from min
      join min_per_lang using(hole, scoring)
) select hole, lang, scoring, user_id, strokes, submitted,
         round(Sb / strokes * 1000) points,
         round(S  / strokes * 1000) points_for_lang
    from strokes
    join bayesian_estimators using(hole, lang, scoring);

-- Needed to refresh concurrently
CREATE UNIQUE INDEX   medals_key ON   medals(user_id, hole, lang, scoring, medal);
CREATE UNIQUE INDEX rankings_key ON rankings(user_id, hole, lang, scoring);

-- Used by /stats
CREATE INDEX solutions_hole_key ON solutions(hole, user_id) WHERE NOT failing;
CREATE INDEX solutions_lang_key ON solutions(lang, user_id) WHERE NOT failing;

CREATE ROLE "code-golf" WITH LOGIN;

-- Only owners can refresh.
ALTER MATERIALIZED VIEW medals   OWNER TO "code-golf";
ALTER MATERIALIZED VIEW rankings OWNER TO "code-golf";

GRANT SELECT, INSERT, UPDATE         ON TABLE    discord_records TO "code-golf";
GRANT SELECT, INSERT, TRUNCATE       ON TABLE    ideas           TO "code-golf";
GRANT SELECT                         ON TABLE    bytes_points    TO "code-golf";
GRANT SELECT                         ON TABLE    chars_points    TO "code-golf";
GRANT SELECT                         ON TABLE    rankings        TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE    sessions        TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE    solutions       TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE    trophies        TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE    users           TO "code-golf";
