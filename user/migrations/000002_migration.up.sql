CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    login varchar(256) unique NOT NULL,
    name varchar(256) NOT NULL,
    email varchar(256) unique NOT NULL,
    password varchar(256) NOT NULL
);
