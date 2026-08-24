-- 014_remove_unused_indices.sql
--
-- Optimizing the database by removing indices that don't seem to be used much.
-- We can always rebuild them later if we need to.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('014_remove_unused_indices', now())
on conflict ("version") do update set started_at = now();

-- Drop indices and reclaim total space of approximately 14.3 GB
drop index if exists index_premis_events_on_institution_id; -- 4 GB, not used recently
drop index if exists index_generic_files_on_created_at; -- 1.2 GB, not used recently
drop index if exists index_premis_events_on_event_type_and_outcome; -- 4.7 GB, barely used
drop index if exists index_premis_events_on_outcome; -- 4.4 GB, not used recently

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '014_remove_unused_indices';
