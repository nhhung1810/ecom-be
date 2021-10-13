create table Users
(
    id SERIAL PRIMARY key,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(50) NOT NULL,
    PASSWORD VARCHAR(1000) NOT NULL
);