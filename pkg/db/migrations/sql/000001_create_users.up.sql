CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    username varchar(30) UNIQUE NOT NULL,
    passwrd TEXT NOT NULL
);