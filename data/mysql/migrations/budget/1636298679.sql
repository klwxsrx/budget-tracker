# Creates test budget cause I don't wanna write code for budget creation :)

SET @event_id = UUID();
SET @budget_id = UUID();

INSERT INTO `event` (surrogate_id, id, aggregate_id, aggregate_name, event_type, event_data, created_at)
VALUES
    (DEFAULT, UUID_TO_BIN(@event_id), UUID_TO_BIN(@budget_id), 'budget', 'budget.created', CONCAT('{"id": "', @budget_id, '","name": "budget","type": "budget.created","title": "Personal budget","currency": "rub"}'), NOW());

INSERT INTO `unsent_event` (id) VALUES (UUID_TO_BIN(@event_id));