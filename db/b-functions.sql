CREATE FUNCTION earn(INOUT earned trophy[], trophy trophy, user_id int) AS $$
BEGIN
    INSERT INTO trophies VALUES (DEFAULT, user_id, trophy)
             ON CONFLICT DO NOTHING;

    IF found THEN
        earned := array_append(earned, trophy);
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION save_solution(code text, hole hole, lang lang, user_id int)
RETURNS trophy[] AS $$
#variable_conflict use_variable
DECLARE
    bytes   int;
    chars   int;
    code_id int;
    earned  trophy[] := '{}'::trophy[];
    holes   int;
BEGIN
    bytes := octet_length(code);
    chars :=  char_length(code);

    -- Ensure we're the only one messing with code.
    LOCK TABLE code IN EXCLUSIVE MODE;

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

    RETURN earned;
END;
$$ LANGUAGE plpgsql;
