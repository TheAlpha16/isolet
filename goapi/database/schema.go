package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/goapi/config"
)

func CreateTables() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// challenge type
	_, err = DB.QueryContext(ctx, 
		`CREATE TYPE chall_type AS ENUM ('static', 'dynamic', 'on-demand'
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// challenge categories
	_, err = DB.QueryContext(ctx, 
		`CREATE TABLE IF NOT EXISTS categories(
			category_id serial PRIMARY KEY,
			category_name text NOT NULL UNIQUE
		)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// users
	_, err = DB.QueryContext(ctx, `
	CREATE TABLE IF NOT EXISTS users(
		userid bigserial PRIMARY KEY,
		email text NOT NULL UNIQUE,
		username text NOT NULL UNIQUE,
		score integer DEFAULT 0,
		rank integer DEFAULT 3,
		password VARCHAR(100) NOT NULL,
		teamid bigint DEFAULT -1
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// teams
	_, err = DB.QueryContext(ctx, `
	CREATE TABLE IF NOT EXISTS teams(
		teamid bigserial PRIMARY KEY,
		teamname text NOT NULL UNIQUE,
		captain bigint NOT NULL REFERENCES users(userid),
		members bigint[] NOT NULL DEFAULT '{}',
		password VARCHAR(100) NOT NULL,
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// flags
	_, err = DB.QueryContext(ctx, `
	CREATE TABLE IF NOT EXISTS flags(
		flagid bigserial,
		userid bigint NOT NULL REFERENCES users(userid),
		level integer,
		password text,
		flag text NOT NULL,
		port integer NOT NULL,
		verified bool DEFAULT false,
		hostname text,
		deadline bigint DEFAULT 2526249600,
		extended integer DEFAULT 1
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// submssion logs
	_, err = DB.QueryContext(ctx, `
	CREATE TABLE IF NOT EXISTS sublogs(
		sid bigserial PRIMARY KEY,
		chall_id integer NOT NULL REFERENCES challenges(chall_id),
		userid bigint NOT NULL REFERENCES users(userid),
		teamid bigint NOT NULL REFERENCES teams(teamid),
		flag text NOT NULL,
		correct boolean NOT NULL,
		ip inet NOT NULL,
		subtime timestamp NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// running instances
	_, err = DB.QueryContext(ctx, `
	CREATE TABLE IF NOT EXISTS running(
		runid bigserial,
		userid bigint NOT NULL REFERENCES users(userid),
		level integer
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// to verify users
	_, err = DB.QueryContext(ctx, `
	CREATE TABLE IF NOT EXISTS toverify(
		vid bigserial PRIMARY KEY,
		email text NOT NULL UNIQUE,
		username text NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		timestamp timestamp NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// challenges
	_, err = DB.QueryContext(ctx, `
	CREATE TABLE IF NOT EXISTS challenges(
		chall_id serial PRIMARY KEY,
		level integer NOT NULL UNIQUE,
		chall_name text NOT NULL,
		category_id integer NOT NULL REFERENCES categories(category_id),
		prompt text,
		flag text,
		type chall_type NOT NULL DEFAULT 'static',
		points integer NOT NULL DEFAULT 100,
		files text[] DEFAULT ARRAY[]::text[],
		hints text[],
		solves integer DEFAULT 0,
		author text default 'anonymous',
		visible boolean DEFAULT false,
		tags text[],
		port integer DEFAULT 0,
		subd text DEFAULT 'localhost',
		cpu integer DEFAULT 5,
		mem integer DEFAULT 10
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	// clear old toverify entries
	_, err = DB.QueryContext(ctx, fmt.Sprintf(`
	CREATE OR REPLACE FUNCTION toverify_delete_old_rows() RETURNS trigger
	LANGUAGE plpgsql
	AS $$
		BEGIN
			DELETE FROM toverify WHERE timestamp < NOW() - INTERVAL '%d minutes';
			RETURN NEW;
		END;
	$$;
	`, config.TOKEN_EXP))
	if err != nil {
		log.Println(err)
		return err
	}

	// enforce instance count
	_, err = DB.QueryContext(ctx, fmt.Sprintf(`
	CREATE OR REPLACE FUNCTION enforce_instance_count() RETURNS trigger
	LANGUAGE plpgsql
	AS $$
	DECLARE
		max_instance_count INTEGER := %d;
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
	`, config.CONCURRENT_INSTANCES))
	if err != nil {
		log.Println(err)
		return err
	}
	
	// trigger to add captain to members
	_, err = DB.QueryContext(ctx, `
	CREATE OR REPLACE FUNCTION add_captain_to_members()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.members := array_append(NEW.members, NEW.captain);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		log.Println(err)
		return err
	}

	// trigger to delete old toverify entries
	_, err = DB.QueryContext(ctx, `
	CREATE OR REPLACE TRIGGER toverify_delete_old_rows_trigger
	BEFORE INSERT ON toverify
	EXECUTE PROCEDURE toverify_delete_old_rows();
	`)
	if err != nil {
		log.Println(err)
		return err
	}

	// trigger to enforce instance count
	_, err = DB.QueryContext(ctx, `
	CREATE TRIGGER add_captain_to_members_trigger
	BEFORE INSERT ON teams
	FOR EACH ROW
	EXECUTE FUNCTION add_captain_to_members();
	`)
	if err != nil {
		log.Println(err)
		return err
	}

	// trigger to enforce instance count
	_, err = DB.QueryContext(ctx, `
	CREATE OR REPLACE TRIGGER enforce_instance_count_trigger
	BEFORE INSERT ON running
	FOR EACH ROW EXECUTE PROCEDURE enforce_instance_count();
	`)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
