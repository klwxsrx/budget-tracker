CREATE TABLE `event` (
    id BINARY(16) PRIMARY KEY,
    type VARCHAR(255),
    aggregate_id BINARY(16),
    aggregate_name VARCHAR(255),
    event_data MEDIUMTEXT,
    created_at DATETIME
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci