CREATE TABLE user_messages (
    id SERIAL NOT NULL PRIMARY KEY,
    receiver INTEGER NOT NULL REFERENCES users(id),
    sender INTEGER NOT NULL,
    tittle VARCHAR NOT NULL,
    detail VARCHAR NOT NULL
);