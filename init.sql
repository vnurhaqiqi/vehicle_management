CREATE SCHEMA vehicle AUTHORIZATION myuser;

CREATE TABLE vehicle.vehicle_locations (
	id serial4 NOT NULL,
	vehicle_id varchar(20) NOT NULL,
	latitude float8 NOT NULL,
	longitude float8 NOT NULL,
	"timestamp" int8 NOT NULL,
	CONSTRAINT vehicle_locations_pkey PRIMARY KEY (id)
);

