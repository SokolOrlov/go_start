CREATE SCHEMA IF NOT EXISTS taxi;


CREATE TABLE if not exists taxi.Drivers (
    id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY, 
    driverid text,
    tripid integer NULL,
	CONSTRAINT uq_tripId
        UNIQUE (tripid)
    );

 INSERT INTO taxi.drivers(driverid)
 SELECT 'driver1' WHERE NOT EXISTS (SELECT * FROM taxi.drivers WHERE driverid = 'driver1');

insert into taxi.drivers(driverid)
 SELECT 'driver2' WHERE NOT EXISTS (SELECT * FROM taxi.drivers WHERE driverid = 'driver2');