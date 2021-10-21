create table Users (
    id SERIAL PRIMARY key,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(50) NOT NULL,
    PASSWORD VARCHAR(1000) NOT NULL
);

create table Products(
    id SERIAL not NULL, 
    name VARCHAR(100) NOT NULL,
    categories VARCHAR(200)[] NOT NULL,
    brand VARCHAR(50) NOT NULL,
    price real not null,
    size VARCHAR(50)[] not null,
    color VARCHAR(50)[] not null,
    quantity int not null,
    description VARCHAR(1000),
    created_date date DEFAULT current_date,
    PRIMARY KEY(id)
);

create table ProductUser(
    productid int not null REFERENCES Products(id),
    userid int not null REFERENCES Users(id),
    PRIMARY KEY(productid, userid)
)

-- create table Order(
--     id SERIAL PRIMARY KEY,
--     productid int not null REFERENCES Products(id),
--     quantity int not null,
--     color
-- )
