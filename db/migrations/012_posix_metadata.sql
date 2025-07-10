-- 012_posix_metadata.sql
--
-- This migration adds eight columns to the generic_files
-- table to track POSIX metadata related to files we ingest.
-- All of these fields are nullable, so as not to substantially
-- increase the size of the 40+ million rows already in the DB.
--
-- Bags ingested after this migration may or may not include
-- POSIX metadata, depending on whether that data is stored in
-- the tarred bag upload at the time it was created.
--

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('012_posix_metadata', now())
on conflict ("version") do update set started_at = now();

-- Add new POSIX metadata columns to generic_files table.
alter table generic_files add column if not exists access_time timestamp null;
alter table generic_files add column if not exists change_time timestamp null;
alter table generic_files add column if not exists mod_time timestamp null;
alter table generic_files add column if not exists gid int8 null;
alter table generic_files add column if not exists gname varchar null;
alter table generic_files add column if not exists "uid" int8 null;
alter table generic_files add column if not exists uname varchar null;
alter table generic_files add column if not exists "mode" int8 null;


-- Now add POSIX metadata columns to generic_files_view.
drop view if exists public.generic_files_view;

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
    gf.access_time,
    gf.change_time,
    gf.mod_time,
    gf.gid,
    gf.gname,
    gf.uid,
    gf.uname,
    gf.mode,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'md5'::text
          ORDER BY checksums.created_at DESC
         LIMIT 1) AS md5,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'sha1'::text
          ORDER BY checksums.created_at DESC
         LIMIT 1) AS sha1,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'sha256'::text
          ORDER BY checksums.created_at DESC
         LIMIT 1) AS sha256,
    ( SELECT checksums.digest
           FROM checksums
          WHERE checksums.generic_file_id = gf.id AND checksums.algorithm::text = 'sha512'::text
          ORDER BY checksums.created_at DESC
         LIMIT 1) AS sha512,
    gf.created_at,
    gf.updated_at
   FROM generic_files gf
     LEFT JOIN intellectual_objects io ON io.id = gf.intellectual_object_id
     LEFT JOIN institutions i ON i.id = gf.institution_id;


-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '012_posix_metadata';
