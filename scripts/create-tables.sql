--psql nobbydobby -U nobbydobbyapp -f create-tables.sql

DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
 	id serial NOT NULL primary key,
 	firstname varchar(50) NOT NULL ,
	lastname varchar(50) NOT NULL ,
 	email varchar(50) unique NOT NULL
 	);


ALTER TABLE users
  OWNER TO nobbydobbyapp;


