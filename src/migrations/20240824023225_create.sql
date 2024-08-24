-- +goose Up
-- +goose StatementBegin
CREATE TABLE invalid_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) COMMENT 'トークン',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    expired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '期限日時'
) COMMENT '無効トークンテーブル'
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invalid_tokens;
-- +goose StatementEnd
