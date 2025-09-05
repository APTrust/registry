-- 013_shrink_db_size.sql
--
-- This migration contains several optimizations that will reduce the size of the database.
-- They include:
-- Removing columns that are no longer used
-- Converting certain enumerated string fields to integer and adding lookup tables

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('013_shrink_db_size', now())
on conflict ("version") do update set started_at = now();

--
alter table premis_events drop column created_at;
alter table premis_events drop column updated_at;

alter table premis_events drop column old_uuid;
alter table premis_events event_type;

alter table premis_events add COLUMN event_type_int smallint;

create or replace function convert_event_types()
returns void as $$
begin
    update premis_events
    set event_type_int = case event_type
        when 'access assignment' then 1
        when 'creation' then 2
        when 'deletion' then 3
        when 'message digest calculation' then 4
        when 'fixity check' then 5
        when 'identifier assignment' then 6
        when 'ingestion' then 7
        when 'replication' then 8
        when 'validation' then 9
        else 0  -- default
    end;
end;
$$ language plpgsql;

select convert_event_types();

-- if exists
alter table premis_events drop column event_type;
alter table premis_events rename column event_type_int TO event_type;
-- create table event_type_lookup
-- Most of these, we are not using at the moment
-- accession
-- appraisal
-- capture
-- compiling
-- compression
-- creation
-- deaccession
-- decompression
-- decryption
-- deletion
-- digital signature generation
-- digital signature validation
-- displaying
-- dissemination
-- encryption
-- execution
-- exporting
-- extraction
-- filename change
-- fixity check
-- forensic feature analysis
-- format identification
-- imaging
-- information package creation
-- information package merging
-- information package splitting
-- ingestion
-- ingestion end
-- ingestion start
-- interpreting
-- message digest calculation
-- metadata extraction
-- metadata modification
-- migration
-- modification
-- normalization
-- packing
-- policy assignment
-- printing
-- quarantine
-- recovery
-- redaction
-- refreshment
-- rendering
-- replication
-- transfer
-- unpacking
-- unquarantine
-- validation
-- virus check

-- add foreign key restraint to event_type to map to event_type_lookup

-- lookup table for object in premis

-- lookup table for agent in premis

alter table checksums drop column created_at;
alter table checksums drop column updated_at;

-- storage records url

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '013_shrink_db_size';
