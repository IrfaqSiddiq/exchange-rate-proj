CREATE TABLE supported_countries(
    id serial PRIMARY KEY,
    country_name character varying(20) UNIQUE,
    country_code character varying(2) UNIQUE,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    currency_code varchar(3)
);


INSERT INTO supported_countries(country_name,country_code,currency_code)VALUES('Zambia','ZM','ZMW');
INSERT INTO exchange_rates(exchange_amount,exchange_rate_time,country_id)VALUES(20.0,'2024-01-27',1);
insert into admin_settings (profit_perc )values(10);