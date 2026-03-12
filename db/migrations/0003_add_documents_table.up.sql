CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS documents (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id       INT NOT NULL REFERENCES users(id),
    file_name     VARCHAR(255) NOT NULL,
    file_path     VARCHAR(512) NOT NULL,
    file_size     BIGINT NOT NULL,
    content_type  VARCHAR(100) NOT NULL,
    uploaded_at   TIMESTAMP DEFAULT NOW()
);
