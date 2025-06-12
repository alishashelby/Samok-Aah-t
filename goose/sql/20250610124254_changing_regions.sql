-- +goose Up
-- +goose StatementBegin
ALTER TABLE regions
ALTER COLUMN name TYPE VARCHAR(41);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regions
ALTER COLUMN name TYPE VARCHAR(40);
-- +goose StatementEnd
