-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    link VARCHAR NOT NULL,
    fakelink VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links;
-- +goose StatementEnd
