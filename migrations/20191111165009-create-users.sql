
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(16) NOT NULL,
    password VARCHAR(255) NOT NULL,
    UNIQUE KEY uq_user_id(user_id)
);
-- +migrate Down
DROP TABLE IF EXISTS users;
