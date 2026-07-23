-- 019_remove_unused_indices.sql
--
-- Optimizing the database by removing indices that don't seem to be used much

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('019_remove_unused_indices', now())
on conflict ("version") do update set started_at = now();

-- Todo check the utilization of these
drop index index_premis_events_on_event_type_and_outcome;
drop index index_premis_events_on_institution_id;
drop index index_premis_events_on_outcome;
drop index index_generic_files_on_created_at;
drop index index_generic_files_on_institution_id_and_state;
drop index index_generic_files_on_institution_id_and_updated_at;
drop index index_generic_files_on_updated_at;
drop index index_storage_records_on_url;

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '019_remove_unused_indices';
