CREATE TABLE food (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE foodtouser (
    userid INT REFERENCES users(id),
    foodid INT REFERENCES food(id),
    lasteaten DATE,
    CONSTRAINT pk_foodtouser PRIMARY KEY (userid, foodid)
);