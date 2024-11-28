BEGIN;

CREATE TABLE todo_tags (
    id SERIAL NOT NULL PRIMARY KEY,
    uuid  UUID NOT NULL DEFAULT gen_random_uuid(),
    todo_id INT NOT NULL,
    tag_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_todo FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
    CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE INDEX todo_tags_id_key ON todo_tags(id);

CREATE INDEX todo_tags_todo_id_idx ON todo_tags(todo_id);

COMMIT;