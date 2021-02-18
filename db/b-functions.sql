CREATE FUNCTION earn(INOUT earned trophy[], trophy trophy, user_id int) AS $$
BEGIN
    INSERT INTO trophies VALUES (DEFAULT, user_id, trophy)
             ON CONFLICT DO NOTHING;

    IF found THEN
        earned := array_append(earned, trophy);
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE hole_rank_ret AS (strokes int, rank int, joint bool);

CREATE FUNCTION hole_rank(hole hole, lang lang, scoring scoring, user_id int)
RETURNS SETOF hole_rank_ret AS $$
BEGIN
    RETURN QUERY EXECUTE FORMAT(
        'WITH ranks AS (
            SELECT %I, RANK() OVER (ORDER BY %I), user_id
              FROM solutions
              JOIN code ON id = code_id
             WHERE NOT failing AND hole = $1 AND lang = $2 AND scoring = $3
        ) SELECT %I, rank::int,
                 (SELECT COUNT(*) != 1 FROM ranks r WHERE r.rank = ranks.rank)
            FROM ranks WHERE user_id = $4',
        scoring, scoring, scoring
    ) USING hole, lang, scoring, user_id;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE save_solution_ret AS (
    earned          trophy[],
    new_bytes       int,
    new_bytes_joint bool,
    new_bytes_rank  int,
    new_chars       int,
    new_chars_joint bool,
    new_chars_rank  int,
    old_bytes       int,
    old_bytes_joint bool,
    old_bytes_rank  int,
    old_chars       int,
    old_chars_joint bool,
    old_chars_rank  int
);

CREATE FUNCTION save_solution(code text, hole hole, lang lang, user_id int)
RETURNS save_solution_ret AS $$
#variable_conflict use_variable
DECLARE
    bytes   int;
    chars   int;
    code_id int;
    earned  trophy[] := '{}'::trophy[];
    holes   int;
    rank    hole_rank_ret;
    ret     save_solution_ret;
BEGIN
    bytes := octet_length(code);
    chars :=  char_length(code);

    -- Ensure we're the only one messing with code.
    LOCK TABLE code IN EXCLUSIVE MODE;

    rank                := hole_rank(hole, lang, 'bytes', user_id);
    ret.old_bytes       := rank.strokes;
    ret.old_bytes_joint := rank.joint;
    ret.old_bytes_rank  := rank.rank;

    rank                := hole_rank(hole, lang, 'chars', user_id);
    ret.old_chars       := rank.strokes;
    ret.old_chars_joint := rank.joint;
    ret.old_chars_rank  := rank.rank;

    -- Lookup the code ID, creating a new record if necessary.
    SELECT id INTO code_id FROM code WHERE code.code = code;
    IF NOT FOUND THEN
        INSERT INTO code (code) VALUES (code) RETURNING id INTO code_id;
    END IF;

    -- Update the code if it's the same length or less, but only update the
    -- submitted time if the solution is shorter. This avoids a user moving
    -- down the leaderboard by matching their personal best.
    INSERT INTO solutions (code_id, hole, lang, scoring, user_id)
         VALUES           (code_id, hole, lang, 'bytes', user_id)
    ON CONFLICT ON CONSTRAINT solutions_pkey
    DO UPDATE SET failing = false,
                submitted = CASE
                    WHEN solutions.failing
                        OR bytes
                        < (SELECT code.bytes FROM code WHERE id = solutions.code_id)
                    THEN excluded.submitted
                    ELSE solutions.submitted
                END,
                    code_id = CASE
                    WHEN solutions.failing
                        OR bytes
                        <= (SELECT code.bytes FROM code WHERE id = solutions.code_id)
                    THEN excluded.code_id
                    ELSE solutions.code_id
                END;

    INSERT INTO solutions (code_id, hole, lang, scoring, user_id)
         VALUES           (code_id, hole, lang, 'chars', user_id)
    ON CONFLICT ON CONSTRAINT solutions_pkey
    DO UPDATE SET failing = false,
                submitted = CASE
                    WHEN solutions.failing
                        OR chars
                        < (SELECT code.chars FROM code WHERE id = solutions.code_id)
                    THEN excluded.submitted
                    ELSE solutions.submitted
                END,
                    code_id = CASE
                    WHEN solutions.failing
                        OR chars
                        <= (SELECT code.chars FROM code WHERE id = solutions.code_id)
                    THEN excluded.code_id
                    ELSE solutions.code_id
                END;

    rank                := hole_rank(hole, lang, 'bytes', user_id);
    ret.new_bytes       := rank.strokes;
    ret.new_bytes_joint := rank.joint;
    ret.new_bytes_rank  := rank.rank;

    rank                := hole_rank(hole, lang, 'chars', user_id);
    ret.new_chars       := rank.strokes;
    ret.new_chars_joint := rank.joint;
    ret.new_chars_rank  := rank.rank;

    -- Remove any orphaned code.
    DELETE FROM code WHERE NOT EXISTS (SELECT FROM solutions WHERE id = solutions.code_id);

    -- Earn trophies.
    SELECT COUNT(DISTINCT solutions.hole) INTO holes
      FROM solutions WHERE NOT failing AND solutions.user_id = user_id;

    IF holes >= 1  THEN earned := earn(earned, 'hello-world',       user_id); END IF;
    IF holes >= 11 THEN earned := earn(earned, 'up-to-eleven',      user_id); END IF;
    IF holes >= 13 THEN earned := earn(earned, 'bakers-dozen',      user_id); END IF;
    IF holes >= 19 THEN earned := earn(earned, 'the-watering-hole', user_id); END IF;
    IF holes >= 40 THEN earned := earn(earned, 'forty-winks',       user_id); END IF;
    IF holes >= 42 THEN earned := earn(earned, 'dont-panic',        user_id); END IF;
    if holes >= 50 THEN earned := earn(earned, 'bullseye',          user_id); END IF;

    IF hole = 'brainfuck' AND lang = 'brainfuck' THEN
        earned := earn(earned, 'inception', user_id);
    END IF;

    IF hole = 'fizz-buzz' THEN
        earned := earn(earned, 'interview-ready', user_id);
    END IF;

    IF hole = 'quine' AND lang = 'python' THEN
        earned := earn(earned, 'ouroboros', user_id);
    END IF;

    IF hole = 'ten-pin-bowling' AND lang = 'cobol' THEN
        earned := earn(earned, 'cobowl', user_id);
    END IF;

    IF lang = 'php' THEN
        earned := earn(earned, 'elephpant-in-the-room', user_id);
    END IF;

    IF (SELECT COUNT(DISTINCT solutions.code_id) > 1 FROM solutions
        WHERE   solutions.user_id = user_id
        AND     solutions.hole = hole
        AND     solutions.lang = lang) THEN
        earned := earn(earned, 'different-strokes', user_id);
    END IF;

    ret.earned := earned;

    RETURN ret;
END;
$$ LANGUAGE plpgsql;
