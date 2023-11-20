CREATE TABLE city (
                    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                    gmt FLOAT DEFAULT '0.00',
                    city_id INT DEFAULT 0,
                    iata_code varchar(255),
                    country_iso2 varchar(255),
                    geoname_id INT DEFAULT 0,
                    latitude float8 DEFAULT '0.00',
                    longitude float8 DEFAULT '0.00',
                    city_name varchar(255),
                    timezone varchar(255),
                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW ()
);
CREATE TABLE country (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       country_name varchar(255),
                       country_iso2 varchar(255),
                       country_iso3 varchar(255),
                       country_iso_numeric INT DEFAULT 0,
                       population INT DEFAULT 0,
                       capital varchar(255),
                       continent varchar (255),
                       currency_name varchar(255),
                       currency_code varchar(255),
                       fips_code varchar(255),
                       phone_prefix varchar(255),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW ()
);
