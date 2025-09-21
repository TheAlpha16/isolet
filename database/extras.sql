-- ### TYPES ###

CREATE TYPE user_role AS ENUM ('admin', 'author', 'captain', 'player');
CREATE TYPE challenge_type AS ENUM ('static', 'dynamic', 'on-demand');

-- ### FUNCTIONS ###

-- function to join a team; enforces team limit
CREATE OR REPLACE FUNCTION join_team(user_id bigint, teamid bigint, user_limit integer)
RETURNS void AS $$
DECLARE
    user_count integer;
BEGIN
    -- Lock the team row to prevent concurrent modifications
    PERFORM 1 FROM teams WHERE id = teamid FOR UPDATE;

    SELECT INTO user_count COUNT(*)
    FROM users
    WHERE team_id = teamid;

    IF user_count >= user_limit THEN
        RAISE EXCEPTION 'TEAM-04';
    END IF;

    UPDATE users
    SET team_id = teamid
    WHERE id = user_id;
END;
$$ LANGUAGE plpgsql;


-- ### INDEXES ###

-- partial index for team_id in users table
CREATE INDEX idx_users_team_id ON users(team_id) WHERE team_id IS NOT NULL;

-- index for (challenge_id, team_id) in solves table
CREATE INDEX idx_solves_team_id_created_at ON solves(team_id, created_at);

-- index for (team_id, created_at) in unlocked_hints table
CREATE INDEX idx_unlocked_hints_team_id_created_at ON unlocked_hints(team_id, created_at);
