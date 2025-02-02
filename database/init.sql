-- Create types
CREATE TYPE chall_type AS ENUM ('static', 'dynamic', 'on-demand');
CREATE TYPE deployment_type AS ENUM ('ssh', 'nc', 'http');
CREATE TYPE token_type AS ENUM ('password_reset', 'invite_token');

-- Create categories table
CREATE TABLE IF NOT EXISTS categories(
    category_id serial PRIMARY KEY,
    category_name text NOT NULL UNIQUE
);

-- Create users table
CREATE TABLE IF NOT EXISTS users(
    userid bigserial PRIMARY KEY,
    email text NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    rank integer DEFAULT 3,
    password VARCHAR(100) NOT NULL,
    teamid bigint DEFAULT -1
);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams(
    teamid bigserial PRIMARY KEY,
    teamname text NOT NULL UNIQUE,
    captain bigint NOT NULL REFERENCES users(userid),
    password VARCHAR(100) NOT NULL,
    cost int DEFAULT 0,
    last_submission bigint DEFAULT EXTRACT(EPOCH FROM NOW())
);

-- Create trigger function to rank up captains
CREATE OR REPLACE FUNCTION rank_up_captain()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users SET rank = 2, teamid = NEW.teamid WHERE userid = NEW.captain;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update captain's rank
CREATE OR REPLACE TRIGGER rank_up_captain_trigger
AFTER INSERT ON teams
FOR EACH ROW
EXECUTE FUNCTION rank_up_captain();

-- Create toverify table
CREATE TABLE IF NOT EXISTS toverify(
    vid bigserial PRIMARY KEY,
    email text NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    timestamp timestamp NOT NULL DEFAULT NOW()
);

-- Function to delete old toverify entries
CREATE OR REPLACE FUNCTION toverify_delete_old_rows() RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM toverify WHERE timestamp < NOW() - INTERVAL '30 minutes';
    RETURN NEW;
END;
$$;

-- Trigger to delete old toverify entries
CREATE OR REPLACE TRIGGER toverify_delete_old_rows_trigger
BEFORE INSERT ON toverify
EXECUTE PROCEDURE toverify_delete_old_rows();

-- Create a table for tokens
CREATE TABLE IF NOT EXISTS tokens (
    tid SERIAL PRIMARY KEY,
    token UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    type token_type NOT NULL,
    userid BIGINT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
    expiry TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '30 minutes')
);

-- Function to delete old tokens
CREATE OR REPLACE FUNCTION tokens_delete_old_rows() RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM tokens WHERE expiry < NOW();
    RETURN NEW;
END;
$$;

-- Trigger to delete old tokens
CREATE OR REPLACE TRIGGER tokens_delete_old_rows_trigger
BEFORE INSERT ON tokens
EXECUTE PROCEDURE tokens_delete_old_rows();

-- Create challenges table
CREATE TABLE IF NOT EXISTS challenges(
    chall_id serial PRIMARY KEY,
    chall_name text NOT NULL UNIQUE,
    category_id integer NOT NULL REFERENCES categories(category_id),
    prompt text,
    flag text,
    type chall_type NOT NULL DEFAULT 'static',
    points integer NOT NULL DEFAULT 100,
    files text[] DEFAULT ARRAY[]::text[],
    requirements int[] DEFAULT '{}',
    hints int[] NOT NULL DEFAULT '{}',
    solves integer DEFAULT 0,
    author text DEFAULT 'anonymous',
    visible boolean DEFAULT false,
    tags text[] DEFAULT ARRAY[]::text[],
    links text[] DEFAULT ARRAY[]::text[]
);

-- Create flags table
CREATE TABLE IF NOT EXISTS flags(
    flagid bigserial,
    teamid bigint NOT NULL REFERENCES teams(teamid),
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    password text,
    flag text NOT NULL,
    port integer,
    hostname text,
    deadline bigint DEFAULT 2526249600,
    extended integer DEFAULT 1
);

-- Create sublogs table
CREATE TABLE IF NOT EXISTS sublogs(
    sid bigserial PRIMARY KEY,
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    userid bigint NOT NULL REFERENCES users(userid),
    teamid bigint NOT NULL REFERENCES teams(teamid),
    flag text NOT NULL,
    correct boolean NOT NULL,
    ip inet NOT NULL,
    timestamp timestamp NOT NULL DEFAULT NOW()
);

-- Hints table
CREATE TABLE IF NOT EXISTS hints(
    hid serial PRIMARY KEY,
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    hint text NOT NULL,
    cost integer NOT NULL DEFAULT 0,
    visible boolean DEFAULT false
);

-- Function to update hints in challenges table
CREATE OR REPLACE FUNCTION update_hints() RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE challenges SET hints = array_append(hints, NEW.hid) WHERE chall_id = NEW.chall_id;
    RETURN NEW;
END;
$$;

-- Trigger to update hints in challenges table
CREATE OR REPLACE TRIGGER update_hints_trigger
AFTER INSERT ON hints
FOR EACH ROW EXECUTE PROCEDURE update_hints();

-- Create solves table
CREATE TABLE IF NOT EXISTS solves(
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    teamid bigint NOT NULL REFERENCES teams(teamid),
    timestamp timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (teamid, chall_id)
);

-- Create unlocked hints table
CREATE TABLE IF NOT EXISTS uhints(
    hid integer NOT NULL REFERENCES hints(hid),
    teamid bigint NOT NULL REFERENCES teams(teamid),
    timestamp timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (teamid, hid)
);

-- Function to update cost in teams when a new hint is unlocked
CREATE OR REPLACE FUNCTION update_team_cost()
RETURNS TRIGGER AS $$
DECLARE
    hint_cost int;
BEGIN
    SELECT cost INTO hint_cost
    FROM hints
    WHERE hid = NEW.hid;

    UPDATE teams
    SET cost = cost + hint_cost
    WHERE teamid = NEW.teamid;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to initiate hint cost update function
CREATE OR REPLACE TRIGGER update_team_cost_trigger
AFTER INSERT
ON uhints
FOR EACH ROW
EXECUTE FUNCTION update_team_cost();

-- Table to store deployment data for on-demand and dynamic challenges
CREATE TABLE IF NOT EXISTS images(
    iid serial PRIMARY KEY,
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    registry text DEFAULT '',
    image text NOT NULL,
    deployment deployment_type NOT NULL DEFAULT 'http',
    port integer DEFAULT 80,
    subd text DEFAULT '',
    cpu integer DEFAULT 0,
    mem integer DEFAULT 0
);

-- Table to store buffer for running on-demand challenges
CREATE TABLE IF NOT EXISTS running(
    runid bigserial,
    teamid bigint NOT NULL REFERENCES teams(teamid),
    chall_id integer
);

-- Function to enforce instance count
CREATE OR REPLACE FUNCTION enforce_instance_count() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    max_instance_count INTEGER := 2; 
    instance_count INTEGER := 0;
    must_check BOOLEAN := false;
BEGIN
    IF TG_OP = 'INSERT' THEN
        must_check := true;
    END IF;

    IF must_check THEN
        -- prevent concurrent inserts from multiple transactions
        LOCK TABLE running IN EXCLUSIVE MODE;

        SELECT INTO instance_count COUNT(*) 
        FROM running 
        WHERE teamid = NEW.teamid;

        IF instance_count >= max_instance_count THEN
            RAISE EXCEPTION 'Cannot start more instances for the team.';
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

-- Trigger to enforce instance count
CREATE OR REPLACE TRIGGER enforce_instance_count_trigger
BEFORE INSERT ON running
FOR EACH ROW EXECUTE PROCEDURE enforce_instance_count();

-- Function to add entry to solves table on correct submission
CREATE OR REPLACE FUNCTION add_solve_entry()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.correct = TRUE THEN
        INSERT INTO solves (chall_id, teamid, timestamp)
        VALUES (NEW.chall_id, NEW.teamid, NOW());

        UPDATE challenges
        SET solves = solves + 1
        WHERE chall_id = NEW.chall_id;

        UPDATE teams
        SET last_submission = EXTRACT(EPOCH FROM NOW())
        WHERE teamid = NEW.teamid;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to add entry to solves table on correct submission
CREATE OR REPLACE TRIGGER add_solve_entry_trigger
AFTER INSERT ON sublogs
FOR EACH ROW
EXECUTE FUNCTION add_solve_entry();

-- Create a GIN index on the requirements column
CREATE INDEX IF NOT EXISTS idx_requirements_gin ON challenges USING gin (requirements);

-- Create a function to retrieve the challenge data for a team
CREATE OR REPLACE FUNCTION get_challenges(team_id bigint)
RETURNS TABLE (
    chall_id integer,
    chall_name text,
    prompt text,
    type chall_type,
    points integer,
    files text[],
    hints json,
    solves integer,
    author text,
    tags text[],
    links text[],
    category_name text,
    deployment deployment_type,
    port integer,
    subd text,
    done boolean
) AS $$
BEGIN
    RETURN QUERY 
    WITH solved_challenges AS (
        SELECT ARRAY_AGG(solves.chall_id) AS solved_array
        FROM solves
        WHERE teamid = team_id
    )
    SELECT 
        ch.chall_id,
        ch.chall_name,
        ch.prompt,
        ch.type,
        ch.points,
        ch.files,
        COALESCE((
            SELECT json_agg(
                jsonb_build_object(
                    'hid', h.hid,
                    'hint', CASE WHEN uh.hid IS NOT NULL THEN h.hint ELSE '' END,
                    'cost', h.cost,
                    'unlocked', uh.hid IS NOT NULL
                )
            ) 
            FROM hints h
            LEFT JOIN uhints uh ON uh.teamid = team_id AND uh.hid = h.hid 
            WHERE h.visible = true 
            AND h.hid = any(ch.hints)
        ), '[]'::json) AS hints,
        ch.solves,
        ch.author,
        ch.tags,
        ch.links,
        cat.category_name,
        COALESCE(img.deployment, 'http') AS deployment,
        COALESCE(img.port, 0) AS port,
        COALESCE(img.subd, '') AS subd,
        ch.chall_id = any(solved_array) AS done
    FROM challenges ch
    JOIN categories cat 
        ON ch.category_id = cat.category_id
    LEFT JOIN images img 
        ON img.chall_id = ch.chall_id
    CROSS JOIN solved_challenges
    WHERE ch.visible = true
    AND (
        ch.requirements = '{}' 
        OR ch.requirements <@ solved_array
    );
END;
$$ LANGUAGE plpgsql;

-- Create a function to calculate score for a team
CREATE OR REPLACE FUNCTION calculate_score(team_id bigint)
RETURNS integer AS $$
DECLARE
    score integer := 0;
BEGIN
    SELECT INTO score COALESCE(SUM(ch.points), 0) - t.cost
    FROM teams t
    LEFT JOIN solves s
        ON s.teamid = t.teamid
    LEFT JOIN challenges ch
        ON ch.chall_id = s.chall_id
    WHERE t.teamid = team_id
    GROUP BY t.teamid;

    RETURN score;
END;
$$ LANGUAGE plpgsql;

-- Create a function to unlock a hint for a team
CREATE OR REPLACE FUNCTION unlock_hint(team_id bigint, hint_id integer)
RETURNS text AS $$
DECLARE
    hint_cost integer;
    hint_hint text;
    challid integer;
    challname text;
    team_score integer;
BEGIN
    IF EXISTS (
        SELECT 1
        FROM uhints
        WHERE teamid = team_id AND hid = hint_id
    ) THEN
        RAISE EXCEPTION 'hint already unlocked';
    END IF;

    SELECT cost, hint, chall_id 
    INTO hint_cost, hint_hint, challid
    FROM hints
    WHERE hid = hint_id AND visible = true;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'hint does not exist';
    END IF;

    WITH solved_challenges AS (
        SELECT ARRAY_AGG(solves.chall_id) AS solved_array
        FROM solves
        WHERE teamid = team_id
    )

    SELECT challenges.chall_name 
    INTO challname
    FROM challenges
    CROSS JOIN solved_challenges
    WHERE challenges.chall_id = challid
    AND challenges.visible = true
    AND (
        challenges.requirements = '{}' 
        OR challenges.requirements <@ solved_array
    );

    IF challname IS NULL THEN
        RAISE EXCEPTION 'hint does not exist';
    END IF;

    SELECT calculate_score(team_id) INTO team_score;

    IF team_score < hint_cost THEN
        RAISE EXCEPTION 'insufficient points';
    END IF;

    INSERT INTO uhints (hid, teamid)
    VALUES (hint_id, team_id);

    RETURN hint_hint;
END;
$$ LANGUAGE plpgsql;

-- Create a function to retrieve scoreboard data
CREATE OR REPLACE FUNCTION get_scoreboard(perPage integer, pageOffset integer)
RETURNS TABLE (
    teamid bigint,
    teamname text,
    rank bigint,
    score bigint
) AS $$
BEGIN
    RETURN QUERY 
    SELECT 
        teams.teamid AS teamid,
        teams.teamname AS teamname,
        RANK() OVER (ORDER BY COALESCE(SUM(challenges.points), 0) - teams.cost DESC, teams.last_submission ASC) AS rank,
        COALESCE(SUM(challenges.points), 0) - teams.cost AS score
    FROM teams
    LEFT JOIN solves
        ON solves.teamid = teams.teamid
    LEFT JOIN challenges
        ON challenges.chall_id = solves.chall_id
    GROUP BY teams.teamid, teams.teamname, teams.cost, teams.last_submission
    ORDER BY rank ASC
    LIMIT perPage
    OFFSET pageOffset;
END;
$$ LANGUAGE plpgsql;

-- Create a function to retrieve top teams submissions
CREATE OR REPLACE FUNCTION get_top_teams_submissions()
RETURNS TABLE (
    teamid bigint,
    teamname text,
    rank bigint,
    submissions json
) AS $$
BEGIN
    RETURN QUERY
    WITH top_teams AS (
        SELECT 
            get_scoreboard.teamid,
            get_scoreboard.teamname,
            get_scoreboard.rank
        FROM get_scoreboard(10, 0)
    ),
    combined_events AS (
        SELECT 
            t.teamid,
            t.teamname,
            t.rank,
            jsonb_build_object(
                'points', c.points,
                'timestamp', s.timestamp
            ) AS event
        FROM top_teams t
        LEFT JOIN solves s ON s.teamid = t.teamid
        LEFT JOIN challenges c ON c.chall_id = s.chall_id
        WHERE c.points IS NOT NULL AND s.timestamp IS NOT NULL

        UNION ALL

        SELECT 
            t.teamid,
            t.teamname,
            t.rank,
            jsonb_build_object(
                'points', -h.cost,
                'timestamp', uh.timestamp
            ) AS event
        FROM top_teams t
        LEFT JOIN uhints uh ON uh.teamid = t.teamid
        LEFT JOIN hints h ON h.hid = uh.hid
        WHERE h.cost IS NOT NULL AND uh.timestamp IS NOT NULL
    )
    SELECT 
        combined_events.teamid,
        combined_events.teamname,
        combined_events.rank,
        COALESCE(json_agg(event), '[]'::json) AS submissions
    FROM combined_events
    GROUP BY combined_events.teamid, combined_events.teamname, combined_events.rank
    ORDER BY combined_events.rank ASC;
END;
$$ LANGUAGE plpgsql;

-- Create a function to join a team
CREATE OR REPLACE FUNCTION join_team(team_id bigint, user_id bigint, user_limit integer)
RETURNS void AS $$
DECLARE
    user_count integer;
BEGIN
    SELECT INTO user_count COUNT(*)
    FROM users
    WHERE teamid = team_id;

    IF user_count >= user_limit THEN
        RAISE EXCEPTION 'team is full';
    END IF;

    UPDATE users
    SET teamid = team_id
    WHERE userid = user_id;
END;
$$ LANGUAGE plpgsql;
