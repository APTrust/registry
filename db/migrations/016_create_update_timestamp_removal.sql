-- 016_create_update_timestamp_removal
--
-- Removes created_at and updated_at where redundant/unneeded.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('016_create_update_timestamp_removal', now())
on conflict ("version") do update set started_at = now();

-- We need to recreate affected views
drop view if exists public.premis_events_view;
drop view if exists public.checksums_view;
drop view if exists public.generic_files_view;

-- Now we can remove the redundant data
alter table premis_events drop column if exists created_at;
alter table premis_events drop column if exists updated_at;
alter table checksums drop column if exists updated_at;

-- Recreate views
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
    pe.agent
   FROM premis_events pe
     LEFT JOIN institutions i ON pe.institution_id = i.id
     LEFT JOIN intellectual_objects io ON pe.intellectual_object_id = io.id
     LEFT JOIN generic_files gf ON pe.generic_file_id = gf.id;

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

-- Note the official file checksum is now ordered by datetime, which is a nullable column
-- Since the sort order is descending, this is ok, because the latest checksum with a
-- datetime will be at the top of the list.
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
update schema_migrations set finished_at = now() where "version" = '016_create_update_timestamp_removal';
