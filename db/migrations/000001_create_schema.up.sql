create table Users (
    id SERIAL PRIMARY key,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(50) NOT NULL,
    PASSWORD VARCHAR(1000) NOT NULL
);

create table Products(
    id SERIAL not NULL,
    name VARCHAR(100) NOT NULL,
    categories VARCHAR(200) [] NOT NULL,
    brand VARCHAR(50) NOT NULL,
    price real not null,
    size VARCHAR(50) [] not null,
    color VARCHAR(50) [] not null,
    quantity int not null,
    description VARCHAR(1000),
    created_date date DEFAULT current_date,
    PRIMARY KEY(id)
);

create table ProductUser(
    productid int not null REFERENCES Products(id),
    userid int not null REFERENCES Users(id),
    PRIMARY KEY(productid, userid)
);

-- create table Orders(
--     id SERIAL PRIMARY KEY,
--     userid int not null REFERENCES Users(id),
--     status int not NULL DEFAULT 0
-- );

create table ProductsOrder(
    orderid SERIAL PRIMARY KEY,
    userid int not null REFERENCES Users(id),
    status VARCHAR(10) not null DEFAULT 'Pending',
    productid int not null REFERENCES Products(id),
    quantity int not null,
    price real not null,
    color VARCHAR(50) not null,
    size VARCHAR(50) not null
    created_date date DEFAULT current_date,
);


ALTER TABLE ProductsOrder 
ADD CONSTRAINT status_domain 
CHECK (status in ('Completed', 'Pending', 'Cancel'))

ALTER TABLE Products
ADD COLUMN archive int default 0;



-- INSERT INTO ProductsOrder(userid, productid, quantity, price, color, size)
-- VALUES ($1, $2, $3, $4, $5, $6) RETURNING orderid


-- EXCUTE EACH OF THIS IN ISOLATION AND IN ORDER
-- FOR UNKNOW REASON, IT WORK PERFECTLY BUT THE 
-- PARSER DON'T UNDERSTAND IT
CREATE OR REPLACE FUNCTION find_remain(integer) RETURNS integer
    AS '
		select 
            (max(p.quantity) - sum(ps.quantity)) as remaining 
        from products as p
		join productsorder as ps on p.id = ps.productid
		WHERE p.id = $1
		GROUP BY p.id;
	'
    LANGUAGE SQL
    IMMUTABLE
    RETURNS NULL ON NULL INPUT;

CREATE OR REPLACE FUNCTION check_insert() 
   RETURNS TRIGGER 
   LANGUAGE PLPGSQL
AS $$
BEGIN
   if NEW.quantity > find_remain(NEW.productid) THEN
   		RAISE EXCEPTION 'Quantity limit exceed';
   END IF;
   
   RETURN NEW;
END;
$$
-- 

CREATE TRIGGER check_quantity
BEFORE INSERT OR UPDATE
ON productsorder
FOR EACH ROW 
EXECUTE PROCEDURE check_insert()