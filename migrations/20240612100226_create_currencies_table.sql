-- +goose Up
-- +goose StatementBegin
CREATE TABLE currencies
(
    id           UUID PRIMARY KEY,
    name         VARCHAR NOT NULL,
    code         VARCHAR NOT NULL UNIQUE,
    type         INT     NOT NULL,
    is_available BOOL    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE currencies IF EXISTS;
-- +goose StatementEnd
