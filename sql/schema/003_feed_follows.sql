







-- +goose Up
CREATE TABLE feeds_follows(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE ON UPDATE CASCADE,
    unique(user_id, feed_id)
);

-- +goose Down
DROP TABLE feeds_follows;