-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title VARCHAR(255),
    url VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    feed_id UUID NOT NULL,
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT posts_feed_id_fkey
        FOREIGN KEY (feed_id) REFERENCES feed(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;