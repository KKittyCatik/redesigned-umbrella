DROP INDEX IF EXISTS idx_pull_request_reviewers_reviewer_id;
DROP INDEX IF EXISTS idx_pull_requests_author_id;
DROP INDEX IF EXISTS idx_users_team_name;

DROP TABLE IF EXISTS pull_request_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;