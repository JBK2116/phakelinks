-- +goose Up
-- +goose StatementBegin
CREATE TABLE explanations (
    id             BIGSERIAL PRIMARY KEY,
    url_mapping_id BIGINT NOT NULL REFERENCES url_mappings(id) ON DELETE CASCADE,
    explanation    TEXT NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE explanations;
-- +goose StatementEnd
