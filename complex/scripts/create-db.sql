-- psql postgres -f create-db.sql


DROP DATABASE IF EXISTS nobbydobby;

DROP ROLE IF EXISTS nobbydobby;

CREATE ROLE nobbydobbyapp LOGIN CREATEDB;

CREATE DATABASE nobbydobby
  WITH OWNER = nobbydobbyapp
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       LC_COLLATE = 'de_DE.UTF-8'
       LC_CTYPE = 'de_DE.UTF-8'
       CONNECTION LIMIT = -1;
       
GRANT ALL ON DATABASE nobbydobby TO nobbydobbyapp;

REVOKE ALL ON DATABASE nobbydobby FROM public;
