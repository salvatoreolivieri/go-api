CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
    roles (name, description, level)
VALUES
    ('user', 'A user can create post and comments', 1);

INSERT INTO
    roles (name, description, level)
VALUES
    ('moderator', 'A moderator can update other users posts', 2);

INSERT INTO
    roles (name, description, level)
VALUES
    ('Admin', 'An admin can update and delete other users post', 3);