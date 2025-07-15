CREATE EXTENSION citext;
CREATE EXTENSION hstore;

CREATE TYPE cheevo AS ENUM (
    '0xdead', 'aged-like-fine-wine', 'alchemist', 'alphabet-soup',
    'archivist', 'assembly-required', 'bakers-dozen', 'big-brother',
    'biohazard', 'bird-is-the-word', 'black-box-testing', 'blackjack',
    'bullseye', 'busy-beaver', 'caffeinated', 'centenarian', 'cobowl',
    'cunning-linguist', 'dammit-janet', 'different-strokes',
    'disappearing-act', 'dont-panic', 'double-slit-experiment',
    'elephpant-in-the-room', 'emergency-room', 'evil-scheme', 'fish-n-chips',
    'fore', 'forty-winks', 'go-forth', 'gone-in-60-holes',
    'happy-birthday-code-golf', 'happy-go-lucky', 'hello-world',
    'hextreme-agony', 'how-about-second-pi', 'hugs-and-kisses', 'inception',
    'independence-day', 'interview-ready', 'into-space', 'its-over-9000',
    'jeweler', 'just-kidding', 'like-comment-subscribe', 'marathon-runner',
    'mary-had-a-little-lambda', 'may-the-4ᵗʰ-be-with-you',
    'my-god-its-full-of-stars', 'neunundneunzig-luftballons', 'off-the-grid',
    'omniglot', 'omniglutton', 'ouroboros', 'overflowing', 'pangramglot',
    'patches-welcome', 'phileas-fogg', 'pi-day', 'polyglot', 'polyglutton',
    'real-programmers', 'right-on', 'rm-rf', 'rtfm', 'rule-34', 's-box-360',
    'slowcoach', 'smörgåsbord', 'solve-quine', 'sounds-quite-nice', 'takeout',
    'the-watering-hole', 'tim-toady', 'tl-dr', 'twelvetide',
    'twenty-kiloleagues', 'under-pressure', 'up-to-eleven', 'vampire-byte',
    'x-factor', 'x86'
);

CREATE TYPE connection AS ENUM (
    'discord', 'github', 'gitlab', 'stack-overflow'
);

CREATE TYPE hole AS ENUM (
    '12-days-of-christmas', '24-game', '99-bottles-of-beer',
    'abundant-numbers', 'abundant-numbers-long', 'apérys-constant',
    'arabic-to-roman', 'arithmetic-numbers', 'arrows', 'ascending-primes',
    'ascii-table', 'billiards', 'brainfuck', 'calendar',
    'card-number-validation', 'catalan-numbers', 'catalans-constant',
    'christmas-trees', 'collatz', 'connect-four', 'css-colors',
    'css-colors-inverse', 'css-grid', 'cubes', 'day-of-week', 'dfa-simulator',
    'diamonds', 'divisors', 'ellipse-perimeters', 'emirp-numbers',
    'emirp-numbers-long', 'emojify', 'evil-numbers', 'evil-numbers-long',
    'factorial-factorisation', 'farey-sequence', 'fibonacci', 'fizz-buzz',
    'flags', 'floyd-steinberg-dithering', 'foo-fizz-buzz-bar',
    'forsyth-edwards-notation', 'fractions', 'game-of-life',
    'gijswijts-sequence', 'gray-code-decoder', 'gray-code-encoder',
    'happy-numbers', 'happy-numbers-long', 'hexagonal-spiral', 'hexdump',
    'highly-composite-numbers', 'hilbert-curve', 'intersection',
    'inventory-sequence', 'isbn', 'jacobi-symbol', 'kaprekar-numbers',
    'kolakoski-constant', 'kolakoski-sequence', 'leap-years',
    'levenshtein-distance', 'leyland-numbers', 'ln-2', 'look-and-say',
    'lucky-numbers', 'lucky-tickets', 'mahjong', 'mandelbrot', 'maze',
    'medal-tally', 'minesweeper', 'morse-decoder', 'morse-encoder',
    'musical-chords', 'n-queens', 'nfa-simulator', 'niven-numbers',
    'niven-numbers-long', 'number-spiral', 'odd-polyomino-tiling',
    'odious-numbers', 'odious-numbers-long', 'ordinal-numbers',
    'p-adic-expansion', 'palindromemordnilap', 'pangram-grep',
    'partition-numbers', 'pascals-triangle', 'pernicious-numbers',
    'pernicious-numbers-long', 'poker', 'polygon-triangulations',
    'polyominoes', 'prime-numbers', 'prime-numbers-long', 'proximity-grid',
    'qr-decoder', 'qr-encoder', 'quadratic-formula', 'quine', 'recamán',
    'repeating-decimals', 'reverse-polish-notation', 'reversi',
    'rijndael-s-box', 'rock-paper-scissors-spock-lizard', 'roman-to-arabic',
    'rot13', 'rule-110', 'scrambled-sort', 'semiprime-numbers', 'set',
    'seven-segment', 'si-units', 'sierpiński-triangle', 'smith-numbers',
    'snake', 'spelling-numbers', 'sphenic-numbers', 'star-wars-gpt',
    'star-wars-opening-crawl', 'sudoku', 'sudoku-fill-in', 'ten-pin-bowling',
    'tic-tac-toe', 'time-distance', 'tongue-twisters', 'topological-sort',
    'transpose-sentence', 'trinomial-triangle', 'turtle', 'tutorial',
    'united-states', 'vampire-numbers', 'van-eck-sequence',
    'zeckendorf-representation', 'zodiac-signs', 'γ', 'λ', 'π', 'τ', 'φ',
    '√2', '𝑒'
);

CREATE TYPE idea_category AS ENUM ('cheevo', 'hole', 'lang', 'other');

CREATE TYPE lang AS ENUM (
    '05ab1e', 'algol-68', 'apl', 'arturo', 'assembly', 'awk', 'bash', 'basic',
    'befunge', 'berry', 'bqn', 'brainfuck', 'c', 'c-sharp', 'civet', 'cjam',
    'clojure', 'cobol', 'coconut', 'coffeescript', 'common-lisp', 'cpp',
    'crystal', 'd', 'dart', 'egel', 'elixir', 'erlang', 'f-sharp', 'factor',
    'fennel', 'fish', 'forth', 'fortran', 'gleam', 'go', 'golfscript',
    'groovy', 'harbour', 'hare', 'haskell', 'haxe', 'hexagony', 'hush', 'hy',
    'iogii', 'j', 'janet', 'java', 'javascript', 'jq', 'julia', 'k', 'kotlin',
    'lua', 'nim', 'ocaml', 'odin', 'pascal', 'perl', 'php', 'picat',
    'powershell', 'prolog', 'python', 'qore', 'r', 'racket', 'raku', 'rebol',
    'rexx', 'rockstar', 'ruby', 'rust', 'scala', 'scheme', 'sed', 'sql',
    'squirrel', 'stax', 'swift', 'tcl', 'tex', 'uiua', 'v', 'vala', 'viml',
    'vyxal', 'wren', 'zig'
);

CREATE TYPE medal AS ENUM ('unicorn', 'diamond', 'gold', 'silver', 'bronze');

CREATE TYPE pronouns AS ENUM ('he/him', 'he/they', 'she/her', 'she/they', 'they/them');

CREATE TYPE scoring AS ENUM ('bytes', 'chars');

CREATE TYPE theme AS ENUM ('auto', 'dark', 'light');

CREATE TABLE discord_records (
    hole    hole NOT NULL,
    lang    lang NOT NULL,
    message text NOT NULL,
    channel text NOT NULL,
    PRIMARY KEY (hole, lang)
);

CREATE TABLE discord_state (
    key   text NOT NULL PRIMARY KEY,
    value text NOT NULL
);

CREATE UNLOGGED TABLE ideas (
    id          int           NOT NULL PRIMARY KEY,
    thumbs_down int           NOT NULL,
    thumbs_up   int           NOT NULL,
    title       text          NOT NULL UNIQUE,
    category    idea_category NOT NULL DEFAULT 'other'
);

CREATE TABLE users (
    id           int       NOT NULL PRIMARY KEY,
    admin        bool      NOT NULL DEFAULT false,
    sponsor      bool      NOT NULL DEFAULT false,
    login        citext    NOT NULL UNIQUE,
    time_zone    text,
    delete       timestamp,
    country      char(2),
    show_country bool      NOT NULL DEFAULT false,
    started      timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    referrer_id  int                REFERENCES users(id) ON DELETE SET NULL,
    theme        theme     NOT NULL DEFAULT 'auto',
    pronouns     pronouns,
    settings     jsonb     NOT NULL DEFAULT '{}'::jsonb,
    about        text      NOT NULL DEFAULT '',
    -- TODO Make country_flag VIRTUAL not STORED when PostgreSQL supports it.
    country_flag char(2)            GENERATED ALWAYS AS
        (CASE WHEN show_country THEN country END) STORED,
    CHECK (country IS NULL OR country ~ '^[A-Z]{2}$'),
    CHECK (id != referrer_id),              -- Can't refer yourself!
    CHECK (login ~ '^[A-Za-z0-9_-]{1,42}$') -- 1 - 42 ASCII word/hyphen chars.
);

CREATE TABLE authors (
    hole    hole NOT NULL,
    user_id int  REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (hole, user_id)
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

CREATE TABLE follows (
    follower_id int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed    timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    PRIMARY KEY (follower_id, followee_id),
    CHECK (follower_id != followee_id)  -- Can't follow yourself!
);

-- config/data/holes.toml is the canonical source of truth for hole data.
-- This table is a shadow copy, updated on startup, used in DB queries.
-- TODO Move category here, remove config.HoleCategoryHstore.
CREATE UNLOGGED TABLE holes (
    id         hole NOT NULL PRIMARY KEY,
    experiment int  NOT NULL
);

-- Ditto for config/data/langs.toml.
CREATE UNLOGGED TABLE langs (
    id           lang   NOT NULL PRIMARY KEY,
    experiment   int    NOT NULL,
    digest_trunc bytea  NOT NULL,
    name         citext NOT NULL
);

CREATE TABLE notes (
    user_id int  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hole    hole NOT NULL,
    lang    lang NOT NULL,
    note    text NOT NULL,
    CHECK (octet_length(note) <= 128 * 1024),
    PRIMARY KEY (user_id, hole, lang)
);

CREATE TABLE sessions (
    id        uuid      NOT NULL DEFAULT GEN_RANDOM_UUID() PRIMARY KEY,
    last_used timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id   int       NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE solutions (
    submitted   timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    bytes       int       NOT NULL,
    chars       int,
    user_id     int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hole        hole      NOT NULL,
    lang        lang      NOT NULL,
    scoring     scoring   NOT NULL,
    failing     bool      NOT NULL DEFAULT false,
    code        text      NOT NULL,
    tested      timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    lang_digest bytea     NOT NULL,
    -- Assembly can only be scored on bytes, and they are compiled bytes.
    CHECK ((lang  = 'assembly' AND chars IS NULL AND scoring = 'bytes')
        OR (lang != 'assembly' AND bytes = octet_length(code)
                               AND chars = char_length(code))),
    -- Solutions are limited to 400 KiB, TODO < 128 KiB (not <=).
    CHECK (octet_length(code) <= 409600),
    PRIMARY KEY (user_id, hole, lang, scoring)
);

CREATE TABLE solutions_log (
    submitted timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    bytes     int       NOT NULL,
    chars     int,
    user_id   int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hole      hole      NOT NULL,
    lang      lang      NOT NULL,
    scoring   scoring   NOT NULL
);

CREATE UNLOGGED TABLE wiki (
    slug    text   NOT NULL PRIMARY KEY,
    section text,
    name    citext NOT NULL,
    html    text   NOT NULL
);

CREATE TABLE trophies (
    earned  timestamp NOT NULL DEFAULT TIMEZONE('UTC', NOW()),
    user_id int       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trophy  cheevo    NOT NULL,
    PRIMARY KEY (user_id, trophy)
);

-- Hole/lang isn't experimental and solution isn't failing.
CREATE VIEW stable_passing_solutions AS
     SELECT solutions.*
       FROM solutions
       JOIN holes ON hole = holes.id
       JOIN langs ON lang = langs.id
      WHERE holes.experiment = 0 AND langs.experiment = 0 AND NOT failing;

CREATE MATERIALIZED VIEW medals AS WITH ranks AS (
    SELECT user_id, hole, lang, scoring, submitted,
           COUNT(*) OVER (PARTITION BY hole, lang, scoring),
           RANK() OVER (
               PARTITION BY hole, lang, scoring
                   ORDER BY CASE WHEN scoring = 'bytes'
                                 THEN bytes ELSE chars END
           )
      FROM stable_passing_solutions
) SELECT user_id, hole, lang, scoring, submitted,
         (enum_range(NULL::medal))[rank + 2] medal
    FROM ranks
   WHERE rank < 4
   UNION ALL
  SELECT MIN(user_id) user_id, hole, lang, scoring, MIN(submitted), 'diamond'::medal
    FROM ranks
   WHERE rank = 1
GROUP BY hole, lang, scoring
  HAVING COUNT(*) = 1
   UNION ALL
  SELECT user_id, hole, lang, scoring, submitted, 'unicorn'::medal
    FROM ranks
   WHERE count = 1;

CREATE MATERIALIZED VIEW rankings AS WITH strokes AS (
    select hole, lang, scoring, user_id, submitted,
           holes.experiment != 0                          experimental_hole,
           langs.experiment != 0                          experimental_lang,
           holes.experiment != 0 OR langs.experiment != 0 experimental,
           case when scoring = 'bytes' then bytes else chars end strokes,
           case when scoring = 'bytes' then chars else bytes end other_strokes
      from solutions
      join holes ON hole = holes.id
      join langs ON lang = langs.id
     where not failing
), min as (
    select hole, scoring, min(strokes)::numeric Sa
      from strokes
     where not experimental
  group by hole, scoring
), min_per_lang as (
    select hole, lang, scoring, min(strokes)::numeric S, sqrt(count(*)) N
     from strokes
     where not experimental
  group by hole, lang, scoring
), bayesian_estimators as (
    select hole, lang, scoring, S,
           ((N + 2) / (N + 3)) * S + (1 / (N + 3)) * Sa Sb
      from min
      join min_per_lang using(hole, scoring)
), points as (
    select hole, lang, scoring, user_id, strokes, other_strokes, submitted,
           experimental_hole, experimental_lang, experimental,
           coalesce(round(Sb / strokes * 1000), 0) points,
           coalesce(round(S  / strokes * 1000), 0) points_for_lang
      from strokes
 left join bayesian_estimators using(hole, lang, scoring)
), ranks as (
    select *,
           count(*)     over (partition by hole, lang, scoring) golfers,
           rank()       over (partition by hole, lang, scoring
                                  order by points_for_lang desc, strokes),
           row_number() over (partition by hole, lang, scoring
                                  order by points_for_lang desc, strokes, submitted) row,
           count(*)     over (partition by hole, scoring) golfers_overall,
           rank()       over (partition by hole, scoring
                                  order by points desc, strokes) rank_overall,
           row_number() over (partition by hole, scoring
                                  order by points desc, strokes, submitted) row_overall
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

-- Materialized views. Only owners can refresh.
ALTER MATERIALIZED VIEW medals   OWNER TO "code-golf";
ALTER MATERIALIZED VIEW points   OWNER TO "code-golf";
ALTER MATERIALIZED VIEW rankings OWNER TO "code-golf";

-- Views.
GRANT SELECT ON stable_passing_solutions TO "code-golf";

-- Tables.
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE authors         TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE connections     TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE discord_records TO "code-golf";
GRANT SELECT, INSERT, UPDATE         ON TABLE discord_state   TO "code-golf";
GRANT SELECT, INSERT,         DELETE ON TABLE follows         TO "code-golf";
GRANT SELECT, INSERT, TRUNCATE       ON TABLE holes           TO "code-golf";
GRANT SELECT, INSERT, TRUNCATE       ON TABLE ideas           TO "code-golf";
GRANT SELECT, INSERT, TRUNCATE       ON TABLE langs           TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE notes           TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE sessions        TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE solutions       TO "code-golf";
GRANT SELECT, INSERT                 ON TABLE solutions_log   TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE trophies        TO "code-golf";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE users           TO "code-golf";
GRANT SELECT, INSERT, TRUNCATE       ON TABLE wiki            TO "code-golf";
