-- +goose Up

CREATE TABLE file_shares(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    shared_with_user_id UUID NOT NULL,

    permission VARCHAR(20) NOT NULL DEFAULT 'read',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_file_shares_file_owner
        FOREIGN KEY (file_id,owner_id)
        REFERENCES files(id, user_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_file_shares_shared_user
        FOREIGN KEY (shared_with_user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_file_shares_permission
        CHECK(permission='read'),

    CONSTRAINT chk_file_shares_not_self
        CHECK(owner_id <> shared_with_user_id),

    CONSTRAINT uq_file_shares_file_user
        UNIQUE(file_id,shared_with_user_id)
);

CREATE INDEX idx_file_shares_shared_with_user
    ON file_shares(shared_with_user_id);

-- +goose Down

DROP TABLE file_shares;