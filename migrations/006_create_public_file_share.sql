-- +goose Up

CREATE TABLE public_file_shares(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,

    CONSTRAINT fk_public_file_shares_file_owner
        FOREIGN KEY(file_id, owner_id)
        REFERENCES files(id,user_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_public_file_share_token
    ON public_file_shares(token);

-- +goose Down
DROP TABLE public_file_shares;