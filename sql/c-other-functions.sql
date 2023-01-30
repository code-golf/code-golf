
CREATE FUNCTION following(int, int) RETURNS int[] AS $$
    SELECT array_append(array_agg(followee_id), $1)
      FROM follows
     WHERE follower_id = $1
     LIMIT $2
$$ LANGUAGE SQL STABLE;

CREATE TYPE hole_rank_ret AS (strokes int, rank int, joint bool);

CREATE FUNCTION hole_rank(hole hole, lang lang, scoring scoring, user_id int)
RETURNS SETOF hole_rank_ret AS $$
BEGIN
    RETURN QUERY EXECUTE FORMAT(
        'WITH ranks AS (
            SELECT %I, RANK() OVER (ORDER BY %I), user_id
              FROM solutions
             WHERE NOT failing AND hole = $1 AND lang = $2 AND scoring = $3
        ) SELECT %I, rank::int,
                 (SELECT COUNT(*) != 1 FROM ranks r WHERE r.rank = ranks.rank)
            FROM ranks WHERE user_id = $4',
        scoring, scoring, scoring
    ) USING hole, lang, scoring, user_id;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE hole_best_ret AS (strokes int, rank int, joint bool, user_id int);

CREATE FUNCTION hole_best(hole hole, lang lang, scoring scoring)
RETURNS SETOF hole_best_ret AS $$
BEGIN
    RETURN QUERY EXECUTE FORMAT(
        'WITH ranks AS (
            SELECT %I, RANK() OVER (ORDER BY %I), user_id
              FROM solutions
             WHERE NOT failing AND hole = $1 AND lang = $2 AND scoring = $3
        ) SELECT %I, rank::int,
                 (SELECT COUNT(*) != 1 FROM ranks r WHERE r.rank = ranks.rank),
                 user_id
            FROM ranks
        ORDER BY rank
           LIMIT 1',
        scoring, scoring, scoring
    ) USING hole, lang, scoring;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION pangramglot(langs lang[]) RETURNS int AS $$
    WITH letters AS (
        SELECT DISTINCT unnest(regexp_split_to_array(nullif(regexp_replace(
               lang::text, '-sharp$|pp$|^fish$|[^a-z]', '', 'g'), ''), ''))
          FROM unnest(langs) lang
    ) SELECT COUNT(*) FROM letters
$$ LANGUAGE SQL STABLE;

CREATE TYPE save_solution_ret AS (
    beat_bytes           int,
    beat_chars           int,
    earned               cheevo[],
    new_bytes            int,
    new_bytes_joint      bool,
    new_bytes_rank       int,
    new_chars            int,
    new_chars_joint      bool,
    new_chars_rank       int,
    old_bytes            int,
    old_bytes_joint      bool,
    old_bytes_rank       int,
    old_chars            int,
    old_chars_joint      bool,
    old_chars_rank       int,
    old_best_bytes       int,
    old_best_bytes_joint bool,
    old_best_chars       int,
    old_best_chars_joint bool
);

CREATE FUNCTION save_solution(
    bytes int, chars int, code text, hole hole, lang lang, user_id int
) RETURNS save_solution_ret AS $$
#variable_conflict use_variable
DECLARE
    old_best hole_best_ret;
    rank     hole_rank_ret;
    ret      save_solution_ret;
BEGIN
    -- Ensure we're the only one messing with solutions.
    LOCK TABLE solutions IN EXCLUSIVE MODE;

    rank                := hole_rank(hole, lang, 'bytes', user_id);
    ret.old_bytes       := rank.strokes;
    ret.old_bytes_joint := rank.joint;
    ret.old_bytes_rank  := rank.rank;

    old_best                 := hole_best(hole, lang, 'bytes');
    ret.old_best_bytes       := old_best.strokes;
    ret.old_best_bytes_joint := old_best.joint;

    IF chars IS NOT NULL THEN
        rank                := hole_rank(hole, lang, 'chars', user_id);
        ret.old_chars       := rank.strokes;
        ret.old_chars_joint := rank.joint;
        ret.old_chars_rank  := rank.rank;

        old_best                 := hole_best(hole, lang, 'chars');
        ret.old_best_chars       := old_best.strokes;
        ret.old_best_chars_joint := old_best.joint;
    END IF;

    -- Update the code if it's the same length or less, but only update the
    -- submitted time if the solution is shorter. This avoids a user moving
    -- down the leaderboard by matching their personal best.
    INSERT INTO solutions (bytes, chars, code, hole, lang, scoring, user_id)
         VALUES           (bytes, chars, code, hole, lang, 'bytes', user_id)
    ON CONFLICT ON CONSTRAINT solutions_pkey
    DO UPDATE SET failing = false,
                    bytes = CASE
                    WHEN solutions.failing OR excluded.bytes <= solutions.bytes
                    THEN excluded.bytes ELSE solutions.bytes END,
                    chars = CASE
                    WHEN solutions.failing OR excluded.bytes <= solutions.bytes
                    THEN excluded.chars ELSE solutions.chars END,
                     code = CASE
                    WHEN solutions.failing OR excluded.bytes <= solutions.bytes
                    THEN excluded.code ELSE solutions.code END,
                submitted = CASE
                    WHEN solutions.failing OR excluded.bytes < solutions.bytes
                    THEN excluded.submitted ELSE solutions.submitted END;

    IF chars IS NOT NULL THEN
        INSERT INTO solutions (bytes, chars, code, hole, lang, scoring, user_id)
             VALUES           (bytes, chars, code, hole, lang, 'chars', user_id)
        ON CONFLICT ON CONSTRAINT solutions_pkey
        DO UPDATE SET failing = false,
                        bytes = CASE
                        WHEN solutions.failing OR excluded.chars <= solutions.chars
                        THEN excluded.bytes ELSE solutions.bytes END,
                        chars = CASE
                        WHEN solutions.failing OR excluded.chars <= solutions.chars
                        THEN excluded.chars ELSE solutions.chars END,
                         code = CASE
                        WHEN solutions.failing OR excluded.chars <= solutions.chars
                        THEN excluded.code ELSE solutions.code END,
                    submitted = CASE
                        WHEN solutions.failing OR excluded.chars < solutions.chars
                        THEN excluded.submitted ELSE solutions.submitted END;
    END IF;

    rank                := hole_rank(hole, lang, 'bytes', user_id);
    ret.new_bytes       := rank.strokes;
    ret.new_bytes_joint := rank.joint;
    ret.new_bytes_rank  := rank.rank;

    IF chars IS NOT NULL THEN
        rank                := hole_rank(hole, lang, 'chars', user_id);
        ret.new_chars       := rank.strokes;
        ret.new_chars_joint := rank.joint;
        ret.new_chars_rank  := rank.rank;
    END IF;

    IF ret.new_bytes_rank = ret.old_bytes_rank THEN
        ret.beat_bytes = ret.old_bytes;
    ELSE
        SELECT MIN(solutions.bytes) INTO ret.beat_bytes
          FROM solutions
         WHERE solutions.hole  = hole
           AND solutions.lang  = lang
           AND solutions.bytes > bytes;
    END IF;

    IF chars IS NOT NULL THEN
        IF ret.new_chars_rank = ret.old_chars_rank THEN
             ret.beat_chars = ret.old_chars;
        ELSE
            SELECT MIN(solutions.chars) INTO ret.beat_chars
              FROM solutions
             WHERE solutions.hole  = hole
               AND solutions.lang  = lang
               AND solutions.chars > chars;
        END IF;
    END IF;

    ret.earned := earn_cheevos(hole, lang, user_id);

    RETURN ret;
END;
$$ LANGUAGE plpgsql;
