
-- +migrate Up
CREATE TABLE IF NOT EXISTS holiday_master (
    id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    date DATE NOT NULL,
    holiday_name VARCHAR(32) NOT NULL,
    UNIQUE KEY uq_date(date)
);

-- +migrate Down
DROP TABLE IF EXISTS holiday_master;
