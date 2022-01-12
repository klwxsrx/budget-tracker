CREATE TABLE IF NOT EXISTS `account`
(
    id              BINARY(16),
    budget_id       BINARY(16),
    status          TINYINT,
    title           VARCHAR(1000),
    initial_balance INT,
    current_balance INT,
    PRIMARY KEY (id),
    FOREIGN KEY (budget_id) REFERENCES budget (id) ON DELETE CASCADE
) ENGINE = InnoDB
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci