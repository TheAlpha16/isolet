-- team_id queries (GetByTeam)
CREATE INDEX "idx_instances_team_id" ON "instances" ("team_id");

-- challenge_id + team_id IS NOT NULL queries (GetByRefs/DeleteByRefs for on-demand instances)
CREATE INDEX "idx_instances_challenge_team_notnull" ON "instances" ("challenge_id", "team_id") WHERE "team_id" IS NOT NULL;

-- challenge_id with team_id IS NULL queries (GetByRefs/DeleteByRefs for dynamic instances)
CREATE INDEX "idx_instances_challenge_null_team" ON "instances" ("challenge_id") WHERE "team_id" IS NULL;

-- expires_at queries (DeleteExpired)
CREATE INDEX "idx_instances_expires_at" ON "instances" ("expires_at") WHERE "expires_at" IS NOT NULL;
