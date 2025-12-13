-- psql postgres -f create-db.sql


DROP DATABASE IF EXISTS felix;

DROP ROLE IF EXISTS felix;

CREATE ROLE felixapp LOGIN CREATEDB;

CREATE DATABASE felix
  WITH OWNER = felixapp
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       LC_COLLATE = 'de_DE.UTF-8'
       LC_CTYPE = 'de_DE.UTF-8'
       CONNECTION LIMIT = -1;
       
GRANT ALL ON DATABASE felix TO felixapp;

REVOKE ALL ON DATABASE felix FROM public;
