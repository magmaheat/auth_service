CREATE TABLE tokens (
                        user_id UUID NOT NULL,
                        token_id UUID NOT NULL PRIMARY KEY,
                        valid VARCHAR(3) NOT NULL DEFAULT 'yes' CHECK (valid IN ('yes', 'no')),
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION get_and_invalidate_token(token UUID)
RETURNS TABLE(valid VARCHAR(3)) AS $$
BEGIN

RETURN QUERY
SELECT valid
FROM tokens
WHERE token_id = $1;

UPDATE tokens
SET valid = 'no'
WHERE token_id = $1 AND valid = 'yes';
END;
$$ LANGUAGE plpgsql;
