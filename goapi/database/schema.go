package database

import "github.com/TitanCrew/isolet/config"

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
		return err
	}

	_, _ = DB.Query(`
	CREATE FUNCTION toverify_delete_old_rows() RETURNS trigger
	LANGUAGE plpgsql
	AS $$
		BEGIN
			DELETE FROM toverify WHERE timestamp < NOW() - INTERVAL '$1 minutes';
			RETURN NEW;
		END;
	$$;
	`, config.TOKEN_EXP)

	_, _ = DB.Query(`
	CREATE TRIGGER toverify_delete_old_rows_trigger
    	BEFORE INSERT ON toverify
    	EXECUTE PROCEDURE toverify_delete_old_rows();
	`)

	return nil
}