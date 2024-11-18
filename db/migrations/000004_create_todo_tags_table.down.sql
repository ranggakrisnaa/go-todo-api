BEGIN;

ALTER TABLE todo_tags DROP CONSTRAINT IF EXISTS fk_todo;
DROP INDEX IF EXISTS todo_tags_id_key;
DROP INDEX IF EXISTS todo_tags_todo_id_idx;
DROP TABLE IF EXISTS todo_tags;

COMMIT;
