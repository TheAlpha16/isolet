package database

import (
	"fmt"
	"log"

	"github.com/CyberLabs-Infosec/isolet/goapi/config"
)

func CreateTables() error {
	var err error

	_, err = DB.Query(`
	CREATE TABLE IF NOT EXISTS users(
		userid bigserial PRIMARY KEY,
		email text NOT NULL UNIQUE,
		username text NOT NULL UNIQUE,
		score integer DEFAULT 0,
		rank integer DEFAULT 3,
		password VARCHAR(100) NOT NULL
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = DB.Query(`
	CREATE TABLE IF NOT EXISTS flags(
		flagid bigserial,
		userid bigint NOT NULL REFERENCES users(userid),
		level integer,
		password text,
		flag text NOT NULL,
		port integer NOT NULL,
		verified bool DEFAULT false,
		hostname text
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = DB.Query(`
	CREATE TABLE IF NOT EXISTS running(
		runid bigserial,
		userid bigint NOT NULL REFERENCES users(userid),
		level integer
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = DB.Query(`
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

	_, err = DB.Query(`
	CREATE TABLE IF NOT EXISTS challenges(
		chall_id serial PRIMARY KEY,
		level integer NOT NULL UNIQUE,
		chall_name text NOT NULL,
		prompt text,
		solves integer DEFAULT 0,
		tags text[]
	)`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = DB.Query(fmt.Sprintf(`
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

	_, err = DB.Query(`
	CREATE OR REPLACE TRIGGER toverify_delete_old_rows_trigger
    	BEFORE INSERT ON toverify
    	EXECUTE PROCEDURE toverify_delete_old_rows();
	`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = DB.Query(fmt.Sprintf(`
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

	_, err = DB.Query(`
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