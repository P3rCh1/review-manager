CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_id, is_active);

CREATE INDEX IF NOT EXISTS idx_pr_reviewers_pr ON pr_reviewers(pr_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_pr_reviewers_pr_user ON pr_reviewers(pr_id, user_id);

CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user ON pr_reviewers(user_id);