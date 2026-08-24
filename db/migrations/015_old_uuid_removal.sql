-- 015_old_uuid_removal.sql
--
-- Removes unneeded old_uuid column from PremisEvents.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('015_old_uuid_removal', now())
on conflict ("version") do update set started_at = now();

-- First we have to remove dependent objects temporarily
drop view public.premis_events_view;

-- Drop the column.
alter table premis_events drop column if exists old_uuid;

-- Recreate the Premis Events view
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
update schema_migrations set finished_at = now() where "version" = '015_old_uuid_removal';
