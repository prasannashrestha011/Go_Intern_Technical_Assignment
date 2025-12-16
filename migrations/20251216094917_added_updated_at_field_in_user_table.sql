-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "user" ADD COLUMN updated_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "user" DROP COLUMN updated_at;
-- +goose StatementEnd
