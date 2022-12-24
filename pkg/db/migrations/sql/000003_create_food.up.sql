CREATE TABLE food(
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    cooktime INTERVAL MINUTE NOT NULL,
    price INT NOT NULL,
    mealtype meal NOT NULL,
    dishtype dish NOT NULL
);

CREATE TABLE foodtousers(
    userid INT REFERENCES users(id),
    foodid INT REFERENCES food(id),
    lasteaten DATE NOT NULL,
    CONSTRAINT pk_foodtousers PRIMARY KEY (userid, foodid)
);