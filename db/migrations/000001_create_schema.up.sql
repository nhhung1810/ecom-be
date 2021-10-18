create table Users (
    id SERIAL PRIMARY key,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(50) NOT NULL,
    PASSWORD VARCHAR(1000) NOT NULL
);

create table Products(
    id SERIAL not NULL, 
    name VARCHAR(100) NOT NULL,
    categories VARCHAR(200) NOT NULL,
    brand VARCHAR(50) NOT NULL,
    price real not null,
    size VARCHAR(50) not null,
    color VARCHAR(50) not null,
    quantity int not null,
    description VARCHAR(1000),
    PRIMARY KEY(id)
);

create table Images(
    id VARCHAR(30),
    dat text,
    PRIMARY KEY(id)
);

create table ProductImages(
    productid int not null REFERENCES Products(id),
    imageid VARCHAR(30) not null REFERENCES Images(id),
    PRIMARY KEY(imageid, productid)
);

create table ProductUser(
    productid int not null REFERENCES Products(id),
    userid int not null REFERENCES Users(id),
    PRIMARY KEY(productid, userid)
)

