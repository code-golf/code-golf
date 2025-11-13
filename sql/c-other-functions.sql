CREATE FUNCTION following(int, int) RETURNS int[] AS $$
  WITH follows AS (
    SELECT followee_id
      FROM follows
     WHERE follower_id = $1
  ORDER BY followed
     LIMIT $2
  ) SELECT array_append(array(SELECT * FROM follows), $1)
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

CREATE TYPE hole_best_except_user_ret AS (strokes int, golfer_count int, user_id int);

CREATE FUNCTION hole_best_except_user(hole hole, lang lang, scoring scoring, user_id int)
RETURNS SETOF hole_best_except_user_ret AS $$
BEGIN
    RETURN QUERY EXECUTE FORMAT(
        'WITH ranks AS (
            SELECT %I, RANK() OVER (ORDER BY %I), user_id
              FROM solutions
             WHERE NOT failing AND hole = $1 AND lang = $2 AND scoring = $3 AND user_id != $4
        ) SELECT %I,
                 (SELECT COUNT(*)::int FROM ranks r WHERE r.rank = ranks.rank),
                 user_id
            FROM ranks
        ORDER BY rank
           LIMIT 1',
        scoring, scoring, scoring
    ) USING hole, lang, scoring, user_id;
END;
$$ LANGUAGE plpgsql;

-- TODO Simplify this function by using the langs table.
CREATE FUNCTION letters(langs lang[]) RETURNS text[] IMMUTABLE RETURN (
    WITH letters AS (
        SELECT DISTINCT unnest(regexp_split_to_array(nullif(regexp_replace(
               lang::text, '-sharp$|pp$|^fish$|[^a-z]', '', 'g'), ''), ''))
          FROM unnest(langs) lang
    ) SELECT array_agg(upper(unnest)) FROM letters
);

CREATE FUNCTION pangramglot(langs lang[]) RETURNS int IMMUTABLE RETURN (
    SELECT coalesce(cardinality(letters(langs)), 0)
);

CREATE TYPE save_solution_ret AS (
    earned                         cheevo[],
    failing_bytes                  int,
    failing_chars                  int,
    new_bytes                      int,
    new_bytes_joint                bool,
    new_bytes_rank                 int,
    new_bytes_solution_count       int,
    new_chars                      int,
    new_chars_joint                bool,
    new_chars_rank                 int,
    new_chars_solution_count       int,
    old_bytes                      int,
    old_bytes_joint                bool,
    old_bytes_rank                 int,
    old_chars                      int,
    old_chars_joint                bool,
    old_chars_rank                 int,
    old_best_bytes                 int,
    old_best_bytes_submitted       timestamp,
    old_best_bytes_first_golfer_id int,
    old_best_bytes_golfer_count    int,
    old_best_bytes_golfer_id       int,
    old_best_chars                 int,
    old_best_chars_submitted       timestamp,
    old_best_chars_first_golfer_id int,
    old_best_chars_golfer_count    int,
    old_best_chars_golfer_id       int
);

CREATE FUNCTION save_solution(
    bytes int, chars int, code text, hole hole, lang lang, time_ms smallint, user_id int
) RETURNS save_solution_ret AS $$
#variable_conflict use_variable
DECLARE
    lang_digest bytea;
    old_best    hole_best_except_user_ret;
    old_bytes   int;
    old_chars   int;
    old_strokes int;
    rank        hole_rank_ret;
    ret         save_solution_ret;
    scoring     scoring;
    strokes     int;
BEGIN
    -- Ensure we're the only one messing with solutions.
    LOCK TABLE solutions IN EXCLUSIVE MODE;

    rank                := hole_rank(hole, lang, 'bytes', user_id);
    ret.old_bytes       := rank.strokes;
    ret.old_bytes_joint := rank.joint;
    ret.old_bytes_rank  := rank.rank;

    SELECT digest_trunc INTO lang_digest FROM langs WHERE id = lang;

    SELECT solutions.bytes, solutions.submitted, solutions.user_id
      INTO ret.old_best_bytes, ret.old_best_bytes_submitted, ret.old_best_bytes_first_golfer_id
      FROM solutions
     WHERE solutions.failing = false
       AND solutions.hole    = hole
       AND solutions.lang    = lang
       AND solutions.scoring = 'bytes'
  ORDER BY solutions.bytes, solutions.submitted
     LIMIT 1;

    -- If the user previously had a failing solution, get the number of strokes.
    SELECT solutions.bytes
      INTO ret.failing_bytes
      FROM solutions
     WHERE solutions.failing = true
       AND solutions.hole    = hole
       AND solutions.lang    = lang
       AND solutions.scoring = 'bytes'
       AND solutions.user_id = user_id;

    IF bytes <= ret.old_best_bytes THEN
        old_best := hole_best_except_user(hole, lang, 'bytes', user_id);
        IF old_best.strokes = ret.old_best_bytes THEN
            ret.old_best_bytes_golfer_count := old_best.golfer_count;
            IF old_best.golfer_count = 1 THEN
                ret.old_best_bytes_golfer_id := old_best.user_id;
            END IF;
        END IF;
    END IF;

    IF chars IS NOT NULL THEN
        rank                := hole_rank(hole, lang, 'chars', user_id);
        ret.old_chars       := rank.strokes;
        ret.old_chars_joint := rank.joint;
        ret.old_chars_rank  := rank.rank;

        SELECT solutions.chars, solutions.submitted, solutions.user_id
          INTO ret.old_best_chars, ret.old_best_chars_submitted, ret.old_best_chars_first_golfer_id
          FROM solutions
         WHERE solutions.failing = false
           AND solutions.hole    = hole
           AND solutions.lang    = lang
           AND solutions.scoring = 'chars'
      ORDER BY solutions.chars, solutions.submitted
         LIMIT 1;

        SELECT solutions.chars
          INTO ret.failing_chars
          FROM solutions
         WHERE solutions.failing = true
           AND solutions.hole    = hole
           AND solutions.lang    = lang
           AND solutions.scoring = 'chars'
           AND solutions.user_id = user_id;

        IF chars <= ret.old_best_chars THEN
            old_best := hole_best_except_user(hole, lang, 'chars', user_id);
            IF old_best.strokes = ret.old_best_chars THEN
                ret.old_best_chars_golfer_count := old_best.golfer_count;
                IF old_best.golfer_count = 1 THEN
                    ret.old_best_chars_golfer_id := old_best.user_id;
                END IF;
            END IF;
        END IF;
    END IF;

    -- Update solutions. First bytes, then chars.
    FOREACH scoring IN ARRAY '{bytes, chars}'::scoring[] LOOP
        strokes := bytes;
        IF scoring = 'chars' THEN strokes := chars; END IF;

        -- Not all langs support all scorings. e.g. Asm has no char scoring.
        IF strokes IS NULL THEN CONTINUE; END IF;

        -- Select information about the current (non-failing) solution.
        SELECT sol.bytes, sol.chars
          INTO old_bytes, old_chars
          FROM solutions sol
         WHERE sol.failing = false
           AND sol.hole    = hole
           AND sol.lang    = lang
           AND sol.scoring = scoring
           AND sol.user_id = user_id;

        old_strokes := old_bytes;
        IF scoring = 'chars' THEN old_strokes := old_chars; END IF;

        -- No existing solution, or it was failing, or the new one is shorter.
        -- Insert or update everything. Also add a history entry.
        IF NOT FOUND OR strokes < old_strokes THEN
            INSERT INTO solutions (bytes, chars, code, hole, lang,
                                   lang_digest, scoring, time_ms, user_id)
                 VALUES           (bytes, chars, code, hole, lang,
                                   lang_digest, scoring, time_ms, user_id)
            ON CONFLICT ON CONSTRAINT solutions_pkey
              DO UPDATE SET bytes       = excluded.bytes,
                            chars       = excluded.chars,
                            code        = excluded.code,
                            failing     = false,
                            lang_digest = excluded.lang_digest,
                            submitted   = excluded.submitted,
                            tested      = excluded.tested,
                            time_ms     = excluded.time_ms;

            INSERT INTO solutions_log (bytes, chars, hole, lang, scoring, user_id)
                 VALUES               (bytes, chars, hole, lang, scoring, user_id);

        -- The new solution is the same length. Keep old submitted, this stops
        -- a user moving down the leaderboard by matching their personal best.
        -- We keep the lowest runtime iff the lang digest hasn't changed.
        ELSIF strokes = old_strokes THEN
            UPDATE solutions
               SET bytes       = bytes,
                   chars       = chars,
                   code        = code,
                   lang_digest = lang_digest,
                   tested      = DEFAULT,
                   time_ms     = CASE WHEN lang_digest = solutions.lang_digest
                                      THEN LEAST(time_ms, solutions.time_ms)
                                      ELSE time_ms
                                      END
             WHERE solutions.hole    = hole
               AND solutions.lang    = lang
               AND solutions.scoring = scoring
               AND solutions.user_id = user_id;

        -- Else, the solution is bigger so don't save it.
        END IF;
    END LOOP;

    rank                := hole_rank(hole, lang, 'bytes', user_id);
    ret.new_bytes       := rank.strokes;
    ret.new_bytes_joint := rank.joint;
    ret.new_bytes_rank  := rank.rank;

    SELECT COUNT(*)
      INTO ret.new_bytes_solution_count
      FROM solutions
     WHERE solutions.failing = false
       AND solutions.hole    = hole
       AND solutions.lang    = lang
       AND solutions.scoring = 'bytes';

    IF chars IS NOT NULL THEN
        rank                := hole_rank(hole, lang, 'chars', user_id);
        ret.new_chars       := rank.strokes;
        ret.new_chars_joint := rank.joint;
        ret.new_chars_rank  := rank.rank;

        SELECT COUNT(*)
          INTO ret.new_chars_solution_count
          FROM solutions
         WHERE solutions.failing = false
           AND solutions.hole    = hole
           AND solutions.lang    = lang
           AND solutions.scoring = 'chars';
    END IF;

    -- Only earn cheevos if the hole and lang aren't experimental.
    SELECT experiment = 0 INTO found FROM holes WHERE id = hole;
    IF found THEN
        SELECT experiment = 0 INTO found FROM langs WHERE id = lang;
        IF found THEN
            ret.earned := earn_cheevos(hole, lang, user_id);
        END IF;
    END IF;

    RETURN ret;
END;
$$ LANGUAGE plpgsql;
