-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats(
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
-- +goose StatementEnd
