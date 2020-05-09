
-- +migrate Up
CREATE TABLE IF NOT EXISTS work_info (
    id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    project_id INT(11) NOT NULL,
    date DATE NOT NULL,
    status INT(11) DEFAULT NULL,
    start VARCHAR(5) DEFAULT NULL,
    end VARCHAR(5) DEFAULT NULL,
    rest VARCHAR(5) DEFAULT NULL,
    total VARCHAR(5) DEFAULT NULL,
    note VARCHAR(255) DEFAULT NULL,
    UNIQUE KEY uq_project_id_date(project_id, date)
);

-- +migrate Down

DROP TABLE IF EXISTS work_info;
