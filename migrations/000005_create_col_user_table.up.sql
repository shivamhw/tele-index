CREATE TABLE IF NOT EXISTS col_users (
    user_id VARCHAR(256), 
    col_id VARCHAR(256), 
    permission INT DEFAULT 1,
    PRIMARY KEY (user_id, col_id), 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (col_id) REFERENCES collections(id) ON DELETE CASCADE
);