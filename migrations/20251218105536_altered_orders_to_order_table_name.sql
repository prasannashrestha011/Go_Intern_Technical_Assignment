-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER table orders RENAME to "order";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER table "order" RENAME to orders;
-- +goose StatementEnd
