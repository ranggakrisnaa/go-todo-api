BEGIN;

DROP TABLE IF EXISTS todos;
DROP INDEX IF EXISTS todo_tags_id_key;
DROP INDEX IF EXISTS todo_tags_todo_id_idx;
ALTER TABLE todo_tags DROP CONSTRAINT fk_todo;

COMMIT;