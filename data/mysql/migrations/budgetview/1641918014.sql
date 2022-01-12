CREATE TABLE IF NOT EXISTS `budget`
(
    id         BINARY(16),
    title      VARCHAR(1000),
    currency   CHAR(3),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
) ENGINE = InnoDB
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci