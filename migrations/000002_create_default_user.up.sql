INSERT INTO teams (name, created_at) 
VALUES ('admin-team', NOW())
ON CONFLICT (name) DO NOTHING;

INSERT INTO users (id, username, team_name, is_active, created_at) 
VALUES (
    'admin-id', 
    'admin', 
    'admin-team', 
    true, 
    NOW()
)
ON CONFLICT (id) DO NOTHING;
