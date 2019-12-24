CREATE TABLE IF NOT EXISTS sessions (    
    session_id    text PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    created_at    timestamp NOT NULL,
    expires_at    timestamp NOT NULL
)
