-- 017_create_update_timestamp_removal
--
-- Removes created_at and updated_at where redundant/unneeded.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('017_create_update_timestamp_removal', now())
on conflict ("version") do update set started_at = now();

-- We need to recreate all premis event views and indices
drop index index_premis_events_on_event_type;
drop index index_premis_events_on_event_type_and_outcome;
drop index ix_premis_event_counts;
drop materialized view public.premis_event_counts;
drop view public.premis_events_view;
drop index index_checksums_on_generic_file_id;
drop view public.checksums_view;
drop view public.generic_files_view;

-- Now we can remove the redundant data
alter table premis_events drop column created_at;
alter table premis_events drop column updated_at;
alter table checksums drop column created_at;
alter table checksums drop column updated_at;

-- Recreate indices and views that use this and reindex - may take some time before indexing is complete
CREATE INDEX index_premis_events_on_event_type ON public.premis_events USING btree (event_type);
CREATE INDEX index_premis_events_on_event_type_and_outcome ON public.premis_events USING btree (event_type, outcome);

CREATE MATERIALIZED VIEW public.premis_event_counts
TABLESPACE pg_default
AS SELECT premis_events.institution_id,
    count(premis_events.id) AS row_count,
    premis_events.event_type,
    premis_events.outcome,
    CURRENT_TIMESTAMP AS updated_at
   FROM premis_events
  GROUP BY CUBE(premis_events.institution_id, premis_events.event_type, premis_events.outcome)
  ORDER BY premis_events.institution_id, premis_events.event_type, premis_events.outcome
WITH DATA;

CREATE UNIQUE INDEX ix_premis_event_counts ON public.premis_event_counts USING btree (institution_id, event_type, outcome);

CREATE OR REPLACE VIEW public.premis_events_view
AS SELECT pe.id,
    pe.identifier,
    pe.institution_id,
    i.name AS institution_name,
    pe.intellectual_object_id,
    io.identifier AS intellectual_object_identifier,
    pe.generic_file_id,
    gf.identifier AS generic_file_identifier,
    pe.event_type,
    pe.date_time,
    pe.detail,
    pe.outcome,
    pe.outcome_detail,
    pe.outcome_information,
    pe.object,
    pe.agent,
    pe.old_uuid
   FROM premis_events pe
     LEFT JOIN institutions i ON pe.institution_id = i.id
     LEFT JOIN intellectual_objects io ON pe.intellectual_object_id = io.id
     LEFT JOIN generic_files gf ON pe.generic_file_id = gf.id;

CREATE INDEX index_checksums_on_generic_file_id ON public.checksums USING btree (generic_file_id);

CREATE OR REPLACE VIEW public.checksums_view
AS SELECT cs.id,
    cs.algorithm,
    cs.datetime,
    cs.digest,
    gf.state,
    gf.identifier AS generic_file_identifier,
    cs.generic_file_id,
    gf.intellectual_object_id,
    gf.institution_id
   FROM checksums cs
     LEFT JOIN generic_files gf ON cs.generic_file_id = gf.id;

CREATE OR REPLACE VIEW public.generic_files_view
AS SELECT gf.id,
    gf.file_format,
    gf.size,
    gf.identifier,
    gf.intellectual_object_id,
    io.identifier AS object_identifier,
    io.access,
    gf.state,
    gf.last_fixity_check,
    gf.institution_id,
    i.name AS institution_name,
    i.identifier AS institution_identifier,
    gf.storage_option,
    gf.uuid,
    gf.mod_time,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'md5'::text
          ORDER BY checksums.datetime DESC
         LIMIT 1) AS md5,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'sha1'::text
          ORDER BY checksums.datetime DESC
         LIMIT 1) AS sha1,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'sha256'::text
          ORDER BY checksums.datetime DESC
         LIMIT 1) AS sha256,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'sha512'::text
          ORDER BY checksums.datetime DESC
         LIMIT 1) AS sha512,
    gf.created_at,
    gf.updated_at
   FROM generic_files gf
     LEFT JOIN intellectual_objects io ON io.id = gf.intellectual_object_id
     LEFT JOIN institutions i ON i.id = gf.institution_id;


-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '017_create_update_timestamp_removal';
