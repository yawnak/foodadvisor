CREATE TABLE ingridients (
    foodid INT REFERENCES food(id),
    name TEXT,
    amount INT,
    units VARCHAR(10),
    CONSTRAINT pk_ingridients PRIMARY KEY (foodid, name)
);