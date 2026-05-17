-- Create enum type "challenge_type"
CREATE TYPE "challenge_type" AS ENUM ('static', 'dynamic', 'on-demand');
-- Create enum type "user_role"
CREATE TYPE "user_role" AS ENUM ('admin', 'author', 'captain', 'player');
-- Create enum type "protocol_type"
CREATE TYPE "protocol_type" AS ENUM ('http', 'https', 'nc', 'ssh');
-- Create enum type "resource_name_type"
CREATE TYPE "resource_name_type" AS ENUM ('cpu', 'memory', 'storage', 'ephemeral-storage');
-- Create enum type "resource_type"
CREATE TYPE "resource_type" AS ENUM ('request', 'limit');
-- Create "categories" table
CREATE TABLE "categories" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "name" text NOT NULL,
  "is_visible" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_categories_name" to table: "categories"
CREATE UNIQUE INDEX "idx_categories_name" ON "categories" ("name");
-- Create "challenges" table
CREATE TABLE "challenges" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "name" text NOT NULL,
  "prompt" text NOT NULL,
  "category_id" bigint NOT NULL,
  "flag" text NULL,
  "type" "challenge_type" NOT NULL DEFAULT 'static',
  "points" bigint NOT NULL,
  "requirements" bigint[] NULL DEFAULT '{}',
  "files" text[] NULL DEFAULT '{}',
  "author" text NOT NULL DEFAULT 'anonymous',
  "tags" text[] NULL DEFAULT '{}',
  "links" text[] NULL DEFAULT '{}',
  "is_visible" boolean NOT NULL DEFAULT false,
  "max_attempts" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
-- Create index "idx_challenges_name" to table: "challenges"
CREATE UNIQUE INDEX "idx_challenges_name" ON "challenges" ("name");
-- Create "config_vars" table
CREATE TABLE "config_vars" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_config_vars_key" to table: "config_vars"
CREATE UNIQUE INDEX "idx_config_vars_key" ON "config_vars" ("key");
-- Create "endpoint_specs" table
CREATE TABLE "endpoint_specs" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "name" text NOT NULL,
  "protocol" "protocol_type" NOT NULL,
  "target_port" integer NOT NULL,
  "manifest_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_endpoint_specs_manifest_id" to table: "endpoint_specs"
CREATE INDEX "idx_endpoint_specs_manifest_id" ON "endpoint_specs" ("manifest_id");
-- Create "endpoints" table
CREATE TABLE "endpoints" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "instance_id" bigint NOT NULL,
  "team_id" bigint NULL,
  "name" text NOT NULL,
  "protocol" "protocol_type" NOT NULL,
  "target_port" integer NOT NULL,
  "hostname" text NULL,
  "port" integer NULL,
  "ready" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_endpoints_instance_id" to table: "endpoints"
CREATE INDEX "idx_endpoints_instance_id" ON "endpoints" ("instance_id");
-- Create "hints" table
CREATE TABLE "hints" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "text" text NOT NULL,
  "cost" bigint NOT NULL DEFAULT 0,
  "is_visible" boolean NOT NULL DEFAULT false,
  "challenge_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_hints_challenge_id" to table: "hints"
CREATE INDEX "idx_hints_challenge_id" ON "hints" ("challenge_id");
-- Create "instances" table
CREATE TABLE "instances" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "team_id" bigint NULL,
  "challenge_id" bigint NOT NULL,
  "flag" text NULL,
  "expires_at" bigint NULL,
  "available_at" bigint NULL,
  "allow_extension" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create "manifests" table
CREATE TABLE "manifests" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "challenge_id" bigint NOT NULL,
  "image" text NOT NULL,
  "flag_template" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_manifests_challenge_id" to table: "manifests"
CREATE UNIQUE INDEX "idx_manifests_challenge_id" ON "manifests" ("challenge_id");
-- Create "resources" table
CREATE TABLE "resources" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "name" "resource_name_type" NOT NULL,
  "value" text NOT NULL,
  "type" "resource_type" NOT NULL,
  "manifest_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_resources_manifest_id" to table: "resources"
CREATE INDEX "idx_resources_manifest_id" ON "resources" ("manifest_id");
-- Create "solves" table
CREATE TABLE "solves" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "challenge_id" bigint NOT NULL,
  "team_id" bigint NOT NULL,
  "submission_id" bigint NOT NULL,
  "points" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_solves_challenge_team" to table: "solves"
CREATE UNIQUE INDEX "idx_solves_challenge_team" ON "solves" ("challenge_id", "team_id");
-- Create index "idx_solves_team_id" to table: "solves"
CREATE INDEX "idx_solves_team_id" ON "solves" ("team_id");
-- Create "submissions" table
CREATE TABLE "submissions" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "team_id" bigint NOT NULL,
  "challenge_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "flag" text NOT NULL,
  "is_correct" boolean NOT NULL,
  "ip_address" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_submissions_team_challenge" to table: "submissions"
CREATE INDEX "idx_submissions_team_challenge" ON "submissions" ("team_id", "challenge_id", "is_correct");
-- Create "teams" table
CREATE TABLE "teams" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "name" text NOT NULL,
  "captain_id" bigint NOT NULL,
  "password" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_teams_name" to table: "teams"
CREATE UNIQUE INDEX "idx_teams_name" ON "teams" ("name");
-- Create "unlocked_hints" table
CREATE TABLE "unlocked_hints" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "team_id" bigint NOT NULL,
  "hint_id" bigint NOT NULL,
  "cost" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_unlocked_hints_team" to table: "unlocked_hints"
CREATE UNIQUE INDEX "idx_unlocked_hints_team" ON "unlocked_hints" ("team_id", "hint_id");
-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "username" text NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "team_id" bigint NULL,
  "role" "user_role" NOT NULL DEFAULT 'player',
  "is_banned" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- Modify "challenges" table
ALTER TABLE "challenges" ADD CONSTRAINT "fk_challenges_category" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "endpoint_specs" table
ALTER TABLE "endpoint_specs" ADD CONSTRAINT "fk_manifests_endpoint_specs" FOREIGN KEY ("manifest_id") REFERENCES "manifests" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "endpoints" table
ALTER TABLE "endpoints" ADD CONSTRAINT "fk_endpoints_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_instances_endpoints" FOREIGN KEY ("instance_id") REFERENCES "instances" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "hints" table
ALTER TABLE "hints" ADD CONSTRAINT "fk_challenges_hints" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "instances" table
ALTER TABLE "instances" ADD CONSTRAINT "fk_instances_challenge" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_instances_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "manifests" table
ALTER TABLE "manifests" ADD CONSTRAINT "fk_manifests_challenge" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "resources" table
ALTER TABLE "resources" ADD CONSTRAINT "fk_manifests_limits" FOREIGN KEY ("manifest_id") REFERENCES "manifests" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_manifests_requests" FOREIGN KEY ("manifest_id") REFERENCES "manifests" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "solves" table
ALTER TABLE "solves" ADD CONSTRAINT "fk_solves_challenge" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_solves_submission" FOREIGN KEY ("submission_id") REFERENCES "submissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_solves_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "submissions" table
ALTER TABLE "submissions" ADD CONSTRAINT "fk_submissions_challenge" FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_submissions_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_submissions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "teams" table
ALTER TABLE "teams" ADD CONSTRAINT "fk_teams_captain" FOREIGN KEY ("captain_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "unlocked_hints" table
ALTER TABLE "unlocked_hints" ADD CONSTRAINT "fk_unlocked_hints_hint" FOREIGN KEY ("hint_id") REFERENCES "hints" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_unlocked_hints_team" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "users" table
ALTER TABLE "users" ADD CONSTRAINT "fk_teams_members" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
