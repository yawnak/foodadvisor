CREATE TABLE ingridients (
    foodid INT REFERENCES food(id),
    name TEXT,
    amount INT NOT NULL,
    units VARCHAR(10) NOT NULL,
    CONSTRAINT pk_ingridients PRIMARY KEY (foodid, name)
);