--
-- PostgreSQL database dump
--

-- Dumped from database version 14.18 (Homebrew)
-- Dumped by pg_dump version 14.18 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: challenge_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.challenge_type AS ENUM (
    'static',
    'dynamic',
    'on-demand'
);


--
-- Name: user_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.user_role AS ENUM (
    'admin',
    'author',
    'captain',
    'player'
);


--
-- Name: join_team(bigint, bigint, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.join_team(user_id bigint, teamid bigint, user_limit integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
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
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    name text NOT NULL,
    is_visible boolean DEFAULT false NOT NULL
);


--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: challenges; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.challenges (
    id bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    name text NOT NULL,
    prompt text NOT NULL,
    category_id bigint NOT NULL,
    flag text,
    type public.challenge_type DEFAULT 'static'::public.challenge_type NOT NULL,
    points bigint NOT NULL,
    files text[] DEFAULT '{}'::text[],
    author text DEFAULT 'anonymous'::text NOT NULL,
    tags text[] DEFAULT '{}'::text[],
    links text[] DEFAULT '{}'::text[],
    is_visible boolean DEFAULT false NOT NULL,
    max_attempts bigint DEFAULT 0 NOT NULL,
    requirements bigint[] DEFAULT '{}'::bigint[]
);


--
-- Name: challenges_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.challenges_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: challenges_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.challenges_id_seq OWNED BY public.challenges.id;


--
-- Name: config_vars; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.config_vars (
    id bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    key text NOT NULL,
    value text NOT NULL
);


--
-- Name: config_vars_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.config_vars_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: config_vars_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.config_vars_id_seq OWNED BY public.config_vars.id;


--
-- Name: hints; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hints (
    id bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    text text NOT NULL,
    cost bigint DEFAULT 0 NOT NULL,
    is_visible boolean DEFAULT false NOT NULL,
    challenge_id bigint NOT NULL
);


--
-- Name: hints_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.hints_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: hints_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.hints_id_seq OWNED BY public.hints.id;


--
-- Name: solves; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solves (
    id bigint NOT NULL,
    created_at bigint,
    challenge_id bigint NOT NULL,
    team_id bigint NOT NULL,
    submission_id bigint NOT NULL,
    points bigint NOT NULL
);


--
-- Name: solves_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solves_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solves_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solves_id_seq OWNED BY public.solves.id;


--
-- Name: submissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions (
    id bigint NOT NULL,
    created_at bigint,
    team_id bigint NOT NULL,
    challenge_id bigint NOT NULL,
    user_id bigint NOT NULL,
    flag text NOT NULL,
    is_correct boolean NOT NULL,
    ip_address text NOT NULL
);


--
-- Name: submissions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.submissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: submissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.submissions_id_seq OWNED BY public.submissions.id;


--
-- Name: teams; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.teams (
    id bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    name text NOT NULL,
    captain_id bigint NOT NULL,
    password text NOT NULL
);


--
-- Name: teams_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.teams_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: teams_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.teams_id_seq OWNED BY public.teams.id;


--
-- Name: unlocked_hints; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.unlocked_hints (
    id bigint NOT NULL,
    created_at bigint,
    team_id bigint NOT NULL,
    hint_id bigint NOT NULL,
    cost bigint NOT NULL
);


--
-- Name: unlocked_hints_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.unlocked_hints_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: unlocked_hints_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.unlocked_hints_id_seq OWNED BY public.unlocked_hints.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    team_id bigint,
    role public.user_role DEFAULT 'player'::public.user_role NOT NULL,
    is_banned boolean DEFAULT false NOT NULL
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: challenges id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.challenges ALTER COLUMN id SET DEFAULT nextval('public.challenges_id_seq'::regclass);


--
-- Name: config_vars id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config_vars ALTER COLUMN id SET DEFAULT nextval('public.config_vars_id_seq'::regclass);


--
-- Name: hints id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hints ALTER COLUMN id SET DEFAULT nextval('public.hints_id_seq'::regclass);


--
-- Name: solves id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solves ALTER COLUMN id SET DEFAULT nextval('public.solves_id_seq'::regclass);


--
-- Name: submissions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ALTER COLUMN id SET DEFAULT nextval('public.submissions_id_seq'::regclass);


--
-- Name: teams id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.teams ALTER COLUMN id SET DEFAULT nextval('public.teams_id_seq'::regclass);


--
-- Name: unlocked_hints id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unlocked_hints ALTER COLUMN id SET DEFAULT nextval('public.unlocked_hints_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: challenges challenges_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.challenges
    ADD CONSTRAINT challenges_pkey PRIMARY KEY (id);


--
-- Name: config_vars config_vars_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config_vars
    ADD CONSTRAINT config_vars_pkey PRIMARY KEY (id);


--
-- Name: hints hints_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hints
    ADD CONSTRAINT hints_pkey PRIMARY KEY (id);


--
-- Name: solves solves_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solves
    ADD CONSTRAINT solves_pkey PRIMARY KEY (id);


--
-- Name: submissions submissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions
    ADD CONSTRAINT submissions_pkey PRIMARY KEY (id);


--
-- Name: teams teams_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.teams
    ADD CONSTRAINT teams_pkey PRIMARY KEY (id);


--
-- Name: unlocked_hints unlocked_hints_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unlocked_hints
    ADD CONSTRAINT unlocked_hints_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_categories_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_categories_name ON public.categories USING btree (name);


--
-- Name: idx_challenges_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_challenges_name ON public.challenges USING btree (name);


--
-- Name: idx_config_vars_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_config_vars_key ON public.config_vars USING btree (key);


--
-- Name: idx_hints_challenge_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_hints_challenge_id ON public.hints USING btree (challenge_id);


--
-- Name: idx_solves_challenge_team; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_solves_challenge_team ON public.solves USING btree (challenge_id, team_id);


--
-- Name: idx_solves_team_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_solves_team_id ON public.solves USING btree (team_id);


--
-- Name: idx_solves_team_id_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_solves_team_id_created_at ON public.solves USING btree (team_id, created_at);


--
-- Name: idx_submissions_team_challenge; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_submissions_team_challenge ON public.submissions USING btree (team_id, challenge_id);


--
-- Name: idx_teams_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_teams_name ON public.teams USING btree (name);


--
-- Name: idx_unlocked_hints_team; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_unlocked_hints_team ON public.unlocked_hints USING btree (team_id, hint_id);


--
-- Name: idx_unlocked_hints_team_id_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_unlocked_hints_team_id_created_at ON public.unlocked_hints USING btree (team_id, created_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_team_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_team_id ON public.users USING btree (team_id) WHERE (team_id IS NOT NULL);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: challenges fk_challenges_category; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.challenges
    ADD CONSTRAINT fk_challenges_category FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: hints fk_challenges_hints; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hints
    ADD CONSTRAINT fk_challenges_hints FOREIGN KEY (challenge_id) REFERENCES public.challenges(id) ON DELETE CASCADE;


--
-- Name: solves fk_solves_challenge; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solves
    ADD CONSTRAINT fk_solves_challenge FOREIGN KEY (challenge_id) REFERENCES public.challenges(id) ON DELETE CASCADE;


--
-- Name: solves fk_solves_submission; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solves
    ADD CONSTRAINT fk_solves_submission FOREIGN KEY (submission_id) REFERENCES public.submissions(id) ON DELETE CASCADE;


--
-- Name: solves fk_solves_team; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solves
    ADD CONSTRAINT fk_solves_team FOREIGN KEY (team_id) REFERENCES public.teams(id) ON DELETE CASCADE;


--
-- Name: submissions fk_submissions_challenge; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions
    ADD CONSTRAINT fk_submissions_challenge FOREIGN KEY (challenge_id) REFERENCES public.challenges(id) ON DELETE CASCADE;


--
-- Name: submissions fk_submissions_team; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions
    ADD CONSTRAINT fk_submissions_team FOREIGN KEY (team_id) REFERENCES public.teams(id) ON DELETE CASCADE;


--
-- Name: submissions fk_submissions_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions
    ADD CONSTRAINT fk_submissions_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: teams fk_teams_captain; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.teams
    ADD CONSTRAINT fk_teams_captain FOREIGN KEY (captain_id) REFERENCES public.users(id);


--
-- Name: users fk_teams_members; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_teams_members FOREIGN KEY (team_id) REFERENCES public.teams(id);


--
-- Name: unlocked_hints fk_unlocked_hints_hint; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unlocked_hints
    ADD CONSTRAINT fk_unlocked_hints_hint FOREIGN KEY (hint_id) REFERENCES public.hints(id) ON DELETE CASCADE;


--
-- Name: unlocked_hints fk_unlocked_hints_team; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unlocked_hints
    ADD CONSTRAINT fk_unlocked_hints_team FOREIGN KEY (team_id) REFERENCES public.teams(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

