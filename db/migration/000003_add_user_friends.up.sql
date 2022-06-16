CREATE TABLE user_friends (
    id SERIAL NOT NULL PRIMARY KEY,
    id_user INTEGER NOT NULL REFERENCES users (id),
    id_friend INTEGER NOT NULL,
    friendList json
);