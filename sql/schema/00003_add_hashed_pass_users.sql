-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN hashed_password TEXT DEFAULT 'unset';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN hashed_password;
-- +goose StatementEnd
