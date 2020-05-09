
-- +migrate Up
CREATE TABLE IF NOT EXISTS projects (
    id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    project_name VARCHAR(255) NOT NULL,
    reward_type TINYINT(9) NOT NULL,
    lower_limit_time INT(3) DEFAULT NULL,
    limit_time INT(3) DEFAULT NULL,
    unit_price INT(11) DEFAULT NULL,
    over_unit_price INT(11) DEFAULT NULL,
    deduction_unit_price INT(11) DEFAULT NULL,
    hourly_wage INT(11) DEFAULT NULL,
    work_time_per_day VARCHAR (4) DEFAULT NULL,
    rest_time VARCHAR (4) DEFAULT NULL
);
-- +migrate Down

DROP TABLE IF EXISTS projects;
