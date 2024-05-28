CREATE TABLE IF NOT EXISTS friend_invites (
     id uuid PRIMARY KEY,
     owner_id uuid NOT NULL,
     user_id uuid NOT NULL,
     status varchar(256) NOT NULL,
     UNIQUE (owner_id, user_id)
);
