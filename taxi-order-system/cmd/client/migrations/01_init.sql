CREATE SCHEMA IF NOT EXISTS taxi;

CREATE TABLE if not exists taxi.Clients (
    id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY, 
    iemai text
    )