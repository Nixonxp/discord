CREATE TABLE IF NOT EXISTS user_friends (
     user_id uuid NOT NULL,
     friend_id uuid NOT NULL,
     UNIQUE (user_id, friend_id)
);
