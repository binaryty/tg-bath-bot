CREATE TABLE articles
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    thumb_url VARCHAR(255),
    published_at VARCHAR(255),
    loaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);