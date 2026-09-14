-- 013_add_auth_app_secret.sql
-- Adds a field to represent a secret from an authenticator app for use with MFA

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('013_add_auth_app_secret', now())
on conflict ("version") do update set started_at = now();

-- Add new encrypted_auth_app_secret to the users table
alter table users add column if not exists encrypted_auth_app_secret varchar null;

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '013_add_auth_app_secret';
