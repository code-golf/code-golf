CREATE FUNCTION save_solution(code text, hole hole, lang lang, user_id int)
RETURNS void AS $$
#variable_conflict use_variable
DECLARE
    bytes   int;
    chars   int;
    code_id int;
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

    -- TODO Port trophies to here.
END;
$$ LANGUAGE plpgsql;
