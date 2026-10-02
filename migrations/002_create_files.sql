-- +goose Up

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    storage_key TEXT NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_files_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE  
);

-- +goose Down
DROP TABLE files;