-- 014_add_encrypted_passkey_session.sql
--

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('014_add_encrypted_passkey_session', now())
on conflict ("version") do update set started_at = now();

-- Add new POSIX metadata column to generic_files table.
alter table users add column if not exists encrypted_passkey_session varchar null;

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '014_add_encrypted_passkey_session';
