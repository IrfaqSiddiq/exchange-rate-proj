CREATE TABLE items_info(
    id serial PRIMARY KEY ,
    item varchar(1000),
    amount double precision,
    purchase_date timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);