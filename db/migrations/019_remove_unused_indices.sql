-- 019_remove_unused_indices.sql
--
-- Optimizing the database by removing indices that don't seem to be used much

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('019_remove_unused_indices', now())
on conflict ("version") do update set started_at = now();

-- None of these have dependencies according to Claude

-- These are too small and risky to justify dropping, and the storage records one is being used
-- drop index index_generic_files_on_institution_id_and_state; -- 1.3 GB, moderately used but probably only benefits staff
-- drop index index_generic_files_on_institution_id_and_updated_at; -- 1.8 GB, barely used in past week
-- drop index index_generic_files_on_updated_at; -- 1.2 GB, not used in past 2 weeks
-- drop index index_storage_records_on_url; -- 9.3 GB - Claude says this is a dead index, but it is being used in sys stats

-- Can drop but will make sysadmin search of events by institution very slow.
drop index index_premis_events_on_institution_id; -- 4 GB, not used, probably only used by staff rarely
-- Low risk according to Claude
drop index index_generic_files_on_created_at; -- 1.2 GB, not used in past 2 weeks
drop index index_premis_events_on_event_type_and_outcome; -- 4.7 GB, barely used
drop index index_premis_events_on_outcome; -- 4.4 GB, not used

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '019_remove_unused_indices';
