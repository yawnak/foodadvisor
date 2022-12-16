CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    username VARCHAR(30) UNIQUE NOT NULL,
    passwrd TEXT NOT NULL,
    expiration INTERVAL DAY DEFAULT interval '0 days'
);