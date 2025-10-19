CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(256) PRIMARY KEY,
    msg_id BIGINT,
    from_id BIGINT,
    chat_id BIGINT,
    size BIGINT,
    filename TEXT,
    createdon DATETIME DEFAULT CURRENT_TIMESTAMP
);