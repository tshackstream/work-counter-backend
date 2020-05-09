
-- +migrate Up
CREATE TABLE IF NOT EXISTS status_master (
    id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    status_name VARCHAR(32) NOT NULL
);

INSERT INTO status_master(status_name) VALUES
('出勤'),
('欠勤');

-- +migrate Down
DROP TABLE IF EXISTS status_master;
