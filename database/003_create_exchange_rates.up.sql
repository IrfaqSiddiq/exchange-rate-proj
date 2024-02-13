CREATE TABLE exchange_rates (
        id serial PRIMARY KEY ,
        exchange_amount double precision,
        created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
        exchange_rate_time timestamp with time zone,
        country_id INT REFERENCES supported_countries (id)
);



COMMENT ON COLUMN exchange_rates.exchange_rate_time IS 'Stores time when exchange rate was recorded by service provider in UTC';
COMMENT ON COLUMN exchange_rates.id IS 'primary key';
COMMENT ON COLUMN exchange_rates.exchange_amount IS 'equivalent amount related to usd';
COMMENT ON COLUMN exchange_rates.created_at IS 'record creation time';
COMMENT ON COLUMN exchange_rates.country_id IS 'country id mapped with supported countries id';