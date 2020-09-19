CREATE TABLE IF NOT EXISTS `event`
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    aggregate_id   BINARY(16),
    aggregate_name VARCHAR(255),
    event_type     VARCHAR(255),
    event_data     MEDIUMTEXT,
    created_at     DATETIME
) ENGINE = InnoDB
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci