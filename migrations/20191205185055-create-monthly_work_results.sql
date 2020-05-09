
-- +migrate Up
CREATE TABLE IF NOT EXISTS monthly_work_results (
    id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    project_id INT(11) NOT NULL,
    month DATE NOT NULL,
    business_day INT(2) NOT NULL,
    input_day INT(2) NOT NULL,
    work_time VARCHAR(16) NOT NULL,
    prospected_work_time VARCHAR(16) NOT NULL,
    prospected_decimal_work_time DECIMAL(5, 2) NOT NULL,
    prospected_reward INT(11) NOT NULL,
    UNIQUE KEY uq_project_id_date(project_id, month)
);

-- +migrate Down
DROP TABLE IF EXISTS monthly_work_results;
