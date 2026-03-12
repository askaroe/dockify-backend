--- For openGauss

CREATE OR REPLACE FUNCTION gen_uuid()
RETURNS uuid
LANGUAGE sql
AS $$
    SELECT md5(random()::text || clock_timestamp()::text)::uuid;
$$;

CREATE TABLE IF NOT EXISTS documents (
    id            UUID PRIMARY KEY DEFAULT gen_uuid(),
    user_id       INT NOT NULL REFERENCES users(id),
    file_name     VARCHAR(255) NOT NULL,
    file_path     VARCHAR(512) NOT NULL,
    file_size     BIGINT NOT NULL,
    content_type  VARCHAR(100) NOT NULL,
    uploaded_at   TIMESTAMP DEFAULT NOW()
);