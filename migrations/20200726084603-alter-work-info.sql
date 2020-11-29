
-- +migrate Up
ALTER TABLE work_info DROP COLUMN `start`;
ALTER TABLE work_info DROP COLUMN `end`;
ALTER TABLE work_info DROP COLUMN `rest`;
ALTER TABLE work_info ADD COLUMN `start_hour` VARCHAR(2) DEFAULT NULL AFTER `date`;
ALTER TABLE work_info ADD COLUMN `start_minute` VARCHAR(2) DEFAULT NULL AFTER `start_hour`;
ALTER TABLE work_info ADD COLUMN `end_hour` VARCHAR(2) DEFAULT NULL AFTER `start_minute`;
ALTER TABLE work_info ADD COLUMN `end_minute` VARCHAR(2) DEFAULT NULL AFTER `end_hour`;
ALTER TABLE work_info ADD COLUMN `rest_hour` VARCHAR(2) DEFAULT NULL AFTER `end_minute`;
ALTER TABLE work_info ADD COLUMN `rest_minute` VARCHAR(2) DEFAULT NULL AFTER `rest_hour`;

-- +migrate Down
ALTER TABLE work_info DROP COLUMN `start_hour`;
ALTER TABLE work_info DROP COLUMN `start_minute`;
ALTER TABLE work_info DROP COLUMN `end_hour`;
ALTER TABLE work_info DROP COLUMN `end_minute`;
ALTER TABLE work_info DROP COLUMN `rest_hour`;
ALTER TABLE work_info DROP COLUMN `rest_minute`;
ALTER TABLE work_info ADD COLUMN `start` VARCHAR(5) DEFAULT NULL AFTER `date`;
ALTER TABLE work_info ADD COLUMN `end` VARCHAR(5) DEFAULT NULL AFTER `start`;
ALTER TABLE work_info ADD COLUMN `rest` VARCHAR(5) DEFAULT NULL AFTER `end`;
