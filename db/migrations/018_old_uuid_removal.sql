-- 018_old_uuid_removal.sql
--
-- Removes unneeded old_uuid column from PremisEvents.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('018_old_uuid_removal', now())
on conflict ("version") do update set started_at = now();

-- First we have to remove dependent objects temporarily
drop index index_premis_events_on_event_type;
drop index index_premis_events_on_event_type_and_outcome;
drop index ix_premis_event_counts;
drop materialized view public.premis_event_counts;
drop view public.premis_events_view;

-- Drop the column.
alter table premis_events drop column old_uuid;

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
    pe.created_at,
    pe.updated_at
   FROM premis_events pe
     LEFT JOIN institutions i ON pe.institution_id = i.id
     LEFT JOIN intellectual_objects io ON pe.intellectual_object_id = io.id
     LEFT JOIN generic_files gf ON pe.generic_file_id = gf.id;

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '018_old_uuid_removal';
