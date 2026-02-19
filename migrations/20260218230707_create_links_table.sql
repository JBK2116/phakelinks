-- +goose Up
-- +goose StatementBegin
CREATE TABLE url_mappings (
    id           BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    fake_url     TEXT NOT NULL,
    technique    TEXT NOT NULL,
    mode         TEXT NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links;
-- +goose StatementEnd
