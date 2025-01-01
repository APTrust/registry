CREATE USER dev_user WITH PASSWORD 'password';
CREATE DATABASE apt_registry_development OWNER dev_user;
grant pg_read_server_files to dev_user;