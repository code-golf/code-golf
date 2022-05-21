CREATE EXTENSION citext;

CREATE TYPE cheevo AS ENUM (
    'assembly-required', 'bakers-dozen', 'black-box-testing', 'blackjack',
    'bullseye', 'caffeinated', 'cobowl', 'cunning-linguist',
    'different-strokes', 'dont-panic', 'elephpant-in-the-room',
    'fish-n-chips', 'forty-winks', 'gone-in-60-holes',
    'happy-birthday-code-golf', 'hello-world', 'inception',
    'independence-day', 'interview-ready', 'its-over-9000', 'just-kidding',
    'marathon-runner', 'may-the-4ᵗʰ-be-with-you', 'my-god-its-full-of-stars',
    'off-the-grid', 'omniglot', 'ouroboros', 'pangramglot', 'patches-welcome',
    'pi-day', 'polyglot', 'polyglutton', 'rtfm', 'rule-34', 'slowcoach',
    'solve-quine', 'sounds-quite-nice', 'takeout', 'the-watering-hole',
    'tim-toady', 'tl-dr', 'twelvetide', 'twenty-kiloleagues',
    'under-pressure', 'up-to-eleven', 'vampire-byte'
);

CREATE TYPE connection AS ENUM (
    'discord', 'github', 'gitlab', 'stack-overflow'
);

CREATE TYPE hole AS ENUM (
    '12-days-of-christmas', '99-bottles-of-beer', 'abundant-numbers',
    'arabic-to-roman', 'arrows', 'ascii-table', 'brainfuck',
    'catalan-numbers', 'catalans-constant', 'christmas-trees', 'collatz',
    'css-colors', 'cubes', 'diamonds', 'divisors', 'emirp-numbers', 'emojify',
    'evil-numbers', 'fibonacci', 'fizz-buzz', 'foo-fizz-buzz-bar',
    'fractions', 'happy-numbers', 'happy-numbers-long', 'hexdump',
    'intersection', 'isbn', 'jacobi-symbol', 'kolakoski-constant',
    'kolakoski-sequence', 'leap-years', 'levenshtein-distance',
    'leyland-numbers', 'look-and-say', 'lucky-numbers', 'lucky-tickets',
    'morse-decoder', 'morse-encoder', 'musical-chords', 'niven-numbers',
    'number-spiral', 'odious-numbers', 'ordinal-numbers', 'pangram-grep',
    'pascals-triangle', 'pernicious-numbers', 'poker', 'prime-numbers',
    'prime-numbers-long', 'qr-decoder', 'quine', 'recamán',
    'reverse-polish-notation', 'rock-paper-scissors-spock-lizard',
    'roman-to-arabic', 'rule-110', 'seven-segment', 'sierpiński-triangle',
    'smith-numbers', 'spelling-numbers', 'star-wars-opening-crawl', 'sudoku',
    'sudoku-v2', 'ten-pin-bowling', 'time-distance', 'tongue-twisters',
    'united-states', 'vampire-numbers', 'van-eck-sequence', 'λ', 'π', 'τ',
    'φ', '√2', '𝑒'
);

CREATE TYPE keymap AS ENUM ('default', 'vim');

CREATE TYPE lang AS ENUM (
    'assembly', 'bash', 'brainfuck', 'c', 'c-sharp', 'cpp', 'cobol',
    'crystal', 'd', 'elixir', 'f-sharp', 'fish', 'fortran', 'go', 'haskell',
    'hexagony', 'j', 'java', 'javascript', 'julia', 'k', 'lisp', 'lua', 'nim',
    'pascal', 'perl', 'php', 'powershell', 'prolog', 'python', 'raku', 'ruby',
    'rust', 'sed', 'sql', 'swift', 'v', 'viml', 'zig'
);

CREATE TYPE medal AS ENUM ('diamond', 'gold', 'silver', 'bronze');

CREATE TYPE scoring AS ENUM ('bytes', 'chars');

CREATE TYPE theme AS ENUM ('auto', 'dark', 'light');

CREATE TABLE discord_records (
    hole    hole NOT NULL,
    lang    lang NOT NULL,
    message text NOT NULL,
    PRIMARY KEY (hole, lang)
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
    -- TODO Make country_flag VIRTUAL not STORED when PostgreSQL supports it.
    country_flag char(2)   NOT NULL GENERATED ALWAYS AS
        (COALESCE(CASE WHEN show_country THEN country END, '')) STORED,
    CHECK (country IS NULL OR country = UPPER(country)),
    CHECK (id != referrer_id)   -- Can't refer yourself
);

CREATE TABLE authors (
    hole    hole NOT NULL,
    user_id int  REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (hole, user_id)
);

CREATE TABLE follows (
    follower_id int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed    timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    PRIMARY KEY (follower_id, followee_id),
    CHECK (follower_id != followee_id)  -- Can't follow yourself!
);

CREATE TABLE connections (
    id            bigint     NOT NULL,
    connection    connection NOT NULL,
    user_id       int        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    discriminator smallint,
    public        bool       NOT NULL DEFAULT false,
    username      text       NOT NULL,
    PRIMARY KEY (connection, id),
    UNIQUE (connection, user_id)
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

CREATE MATERIALIZED VIEW rankings AS WITH strokes AS (
    select hole, lang, scoring, user_id, submitted,
           case when scoring = 'bytes' then bytes else chars end strokes,
           case when scoring = 'bytes' then chars else bytes end other_strokes
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
           ((N + 2) / (N + 3)) * S + (1 / (N + 3)) * Sa Sb
      from min
      join min_per_lang using(hole, scoring)
), points as (
    select hole, lang, scoring, user_id, strokes, other_strokes, submitted,
           round(Sb / strokes * 1000) points,
           round(S  / strokes * 1000) points_for_lang
      from strokes
      join bayesian_estimators using(hole, lang, scoring)
), ranks as (
    select *, rank() over (
               partition by hole, lang, scoring
                   order by points_for_lang desc, strokes
           )
      from points
), tie_count as (
    select hole, lang, scoring, strokes, count(*) tie_count
      from strokes
  group by hole, lang, scoring, strokes
) select * from ranks join tie_count using (hole, lang, scoring, strokes);

CREATE MATERIALIZED VIEW points AS WITH max_points_per_hole AS (
    SELECT DISTINCT ON (user_id, hole, scoring) user_id, scoring, points
      FROM rankings
  ORDER BY user_id, hole, scoring, points DESC
) SELECT user_id, scoring, SUM(points) points
    FROM max_points_per_hole
GROUP BY user_id, scoring;

-- Needed to refresh concurrently
CREATE UNIQUE INDEX   medals_key ON   medals(user_id, hole, lang, scoring, medal);
CREATE UNIQUE INDEX   points_key ON   points(user_id, scoring);
CREATE UNIQUE INDEX rankings_key ON rankings(user_id, hole, lang, scoring);

-- Used by /stats
CREATE INDEX solutions_hole_key ON solutions(hole, user_id) WHERE NOT failing;
CREATE INDEX solutions_lang_key ON solutions(lang, user_id) WHERE NOT failing;

CREATE ROLE "code-golf" WITH LOGIN;

-- Only owners can refresh.
ALTER MATERIALIZED VIEW medals   OWNER TO "code-golf";
ALTER MATERIALIZED VIEW points   OWNER TO "code-golf";
ALTER MATERIALIZED VIEW rankings OWNER TO "code-golf";

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE authors         TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE connections     TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE discord_records TO "code-golf";
GRANT SELECT, INSERT,         DELETE ON TABLE follows         TO "code-golf";
GRANT SELECT, INSERT, TRUNCATE       ON TABLE ideas           TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE sessions        TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE solutions       TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE trophies        TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE users           TO "code-golf";
