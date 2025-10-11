CREATE TABLE IF NOT EXISTS col_items (
    item_id VARCHAR(256), 
    col_id VARCHAR(256), 
    PRIMARY KEY (item_id, col_id), 
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
    FOREIGN KEY (col_id) REFERENCES collections(id) ON DELETE CASCADE
);