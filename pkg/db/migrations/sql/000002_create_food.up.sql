CREATE TABLE food (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    cooktime INTERVAL MINUTE NOT NULL,
    price INT NOT NULL,
    isbreakfast BOOLEAN NOT NULL,
    isdinner BOOLEAN NOT NULL,
    issupper BOOLEAN NOT NULL
);

CREATE TABLE foodtouser (
    userid INT REFERENCES users(id),
    foodid INT REFERENCES food(id),
    lasteaten DATE,
    CONSTRAINT pk_foodtouser PRIMARY KEY (userid, foodid)
);