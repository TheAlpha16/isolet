-- Create types
CREATE TYPE chall_type AS ENUM ('static', 'dynamic', 'on-demand');
CREATE TYPE deployment_type AS ENUM ('ssh', 'nc', 'http');

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
    score integer DEFAULT 0,
    rank integer DEFAULT 3,
    password VARCHAR(100) NOT NULL,
    teamid bigint DEFAULT -1
);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams(
    teamid bigserial PRIMARY KEY,
    teamname text NOT NULL UNIQUE,
    captain bigint NOT NULL REFERENCES users(userid),
    members bigint[] NOT NULL DEFAULT '{}',
    password VARCHAR(100) NOT NULL,
    solved text[] DEFAULT ARRAY[]::text[],
    uhints int[] DEFAULT '{}'
);

-- Create trigger function to add captain to members array
CREATE OR REPLACE FUNCTION add_captain_to_members()
RETURNS TRIGGER AS $$
BEGIN
    NEW.members := array_append(NEW.members, NEW.captain);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to call add_captain_to_members function before insert
CREATE TRIGGER add_captain_to_members_trigger
BEFORE INSERT ON teams
FOR EACH ROW
EXECUTE FUNCTION add_captain_to_members();

-- Create flags table
CREATE TABLE IF NOT EXISTS flags(
    flagid bigserial,
    teamid bigint NOT NULL REFERENCES teams(teamid),
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    password text,
    flag text NOT NULL,
    port integer,
    verified boolean DEFAULT false,
    hostname text,
    deadline bigint DEFAULT 2526249600,
    extended integer DEFAULT 1
);

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
    requirements text[] DEFAULT ARRAY[]::text[],
    hints bigint[] NOT NULL DEFAULT '{}',
    solves integer DEFAULT 0,
    author text DEFAULT 'anonymous',
    visible boolean DEFAULT false,
    tags text[] DEFAULT ARRAY[]::text[],
    links text[] DEFAULT ARRAY[]::text[]
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
    subtime timestamp NOT NULL DEFAULT NOW()
);

-- Hints table
CREATE TABLE IF NOT EXISTS hints(
    hid bigserial PRIMARY KEY,
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    hint text NOT NULL,
    cost integer NOT NULL DEFAULT 0,
    visible boolean DEFAULT false
);

-- Function to uppdate hints in challenges table
CREATE OR REPLACE FUNCTION update_hints() RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE challenges SET hints = array_append(hints, NEW.hid) WHERE chall_id = NEW.chall_id;
    RETURN NEW;
END;
$$;

-- Trigger to update hints in challenges table
CREATE TRIGGER update_hints_trigger
AFTER INSERT ON hints
FOR EACH ROW EXECUTE PROCEDURE update_hints();

-- Table to store deployment data for on-demand and dynamic challenges
CREATE TABLE IF NOT EXISTS images(
    iid bigserial PRIMARY KEY,
    chall_id integer NOT NULL REFERENCES challenges(chall_id),
    registry text NOT NULL,
    image text NOT NULL,
    deployment deployment_type NOT NULL DEFAULT 'http',
    port integer DEFAULT 80,
    subd text DEFAULT 'localhost',
    cpu integer DEFAULT 5,
    mem integer DEFAULT 10
);

-- Function to enforce instance count
CREATE OR REPLACE FUNCTION enforce_instance_count() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    max_instance_count INTEGER := 5; -- Replace with your value from config
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
        WHERE userid = NEW.userid;

        IF instance_count >= max_instance_count THEN
            RAISE EXCEPTION 'Cannot start more instances for the user.';
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

-- Trigger to enforce instance count
CREATE OR REPLACE TRIGGER enforce_instance_count_trigger
BEFORE INSERT ON running
FOR EACH ROW EXECUTE PROCEDURE enforce_instance_count();
