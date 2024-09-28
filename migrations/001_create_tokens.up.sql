CREATE TABLE tokens (
    user_id UUID NOT NULL,
    token_id UUID NOT NULL PRIMARY KEY,
    valid VARCHAR(3) NOT NULL DEFAULT 'yes' CHECK (valid IN ('yes', 'no')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
