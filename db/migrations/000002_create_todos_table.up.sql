BEGIN;

CREATE TABLE todos (
    id SERIAL NOT NULL PRIMARY KEY,
    uuid  UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT TRUE,
    due_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX todos_id_key ON todos(id);

CREATE INDEX todos_user_id_idx ON todos(user_id);

COMMIT;


