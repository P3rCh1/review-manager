CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE teams (
    id UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE
);

CREATE TABLE pull_requests (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    status pr_status NOT NULL DEFAULT 'OPEN',
    
    author_id TEXT NOT NULL REFERENCES users(id),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    merged_at TIMESTAMPTZ NULL
);

CREATE TABLE pr_reviewers (
    pr_id TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,  

    PRIMARY KEY (pr_id, user_id)
);
