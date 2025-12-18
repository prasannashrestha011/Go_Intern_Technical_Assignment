-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "order" RENAME COLUMN create_at to created_at 
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "order" RENAME COLUMN created_at to create_at 
-- +goose StatementEnd
