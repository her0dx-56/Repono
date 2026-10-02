-- +goose Up

ALTER TABLE files
ADD CONSTRAINT uq_files_id_user
UNIQUE (id,user_id);

-- +goose Down

ALTER TABLE files
DROP CONSTRAINT uq_files_id_user;