-- +goose Up
CREATE INDEX idx_files_user_id
ON files(user_id);

CREATE INDEX idx_files_created_at
ON files(created_at);

-- +goose Down

DROP INDEX idx_files_user_id;
DROP INDEX idx_files_creaed_at;