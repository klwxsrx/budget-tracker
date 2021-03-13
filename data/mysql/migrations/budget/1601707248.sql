CREATE TABLE IF NOT EXISTS `last_notified_event`
(
    id       INT PRIMARY KEY,
    event_id INT NOT NULL
) ENGINE = InnoDB
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci