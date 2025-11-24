SELECT 'CREATE ROLE myuser WITH LOGIN PASSWORD ''mypassword'''
    WHERE NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'myuser');
\gexec

SELECT 'CREATE DATABASE prdb'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'prdb');
\gexec

SELECT 'GRANT ALL PRIVILEGES ON DATABASE prdb TO myuser'
    WHERE EXISTS (SELECT FROM pg_database WHERE datname = 'prdb')
	AND EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'myuser');
\gexec

\connect prdb

CREATE TABLE teams (
                       team_name TEXT PRIMARY KEY,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
                       user_id TEXT PRIMARY KEY,
                       username TEXT NOT NULL,
                       team_name TEXT NOT NULL REFERENCES teams(team_name) ON DELETE RESTRICT,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_users_team_name ON users(team_name);
CREATE INDEX idx_users_is_active ON users(is_active);

CREATE TABLE pull_requests (
                               pull_request_id TEXT PRIMARY KEY,
                               pull_request_name TEXT NOT NULL,
                               author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
                               status TEXT NOT NULL CHECK (status IN ('OPEN','MERGED')) DEFAULT 'OPEN',
                               created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                               merged_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_pr_author ON pull_requests(author_id);
CREATE INDEX idx_pr_status ON pull_requests(status);

CREATE TABLE pr_reviewers (
                              pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
                              reviewer_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
                              assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                              PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX idx_pr_reviewers_pr ON pr_reviewers(pull_request_id);
CREATE INDEX idx_pr_reviewers_reviewer ON pr_reviewers(reviewer_id);
