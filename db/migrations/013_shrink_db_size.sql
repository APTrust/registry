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

-- alter table premis_events event_type; --

alter table premis_events add COLUMN event_type_int smallint;
 
-- IMPORTANT --
-- TO DO: If there is a value in the current premis_events table --
-- for eventType that is NOT a match for any values in this function, --
-- probably we need to abort and roll back. If it converts to a 0, --
-- we will lose whatever information was in there. Same for object and agent fields --
create or replace function convert_event_types()
returns void as $$
begin
    update premis_events
    set event_type_int = case event_type
        when 'access assignment' then 1
        when 'accession' then 2
        when 'appraisal' then 3
        when 'capture' then 4
        when 'compiling' then 5
        when 'compression' then 6
        when 'creation' then 7
        when 'deaccession' then 8
        when 'decompression' then 9
        when 'decryption' then 10
        when 'deletion' then 11
        when 'digital signature generation' then 12
        when 'digital signature validation' then 13
        when 'displaying' then 14
        when 'dissemination' then 15
        when 'encryption' then 16
        when 'execution' then 17
        when 'exporting' then 18
        when 'extraction' then 19
        when 'filename change' then 20
        when 'fixity check' then 21
        when 'forensic feature analysis' then 22
        when 'format identification' then 23
        when 'identifier assignment' then 24
        when 'imaging' then 25
        when 'information package creation' then 26
        when 'information package merging' then 27
        when 'information package splitting' then 28
        when 'ingestion' then 29
        when 'ingestion end' then 30
        when 'ingestion start' then 31
        when 'interpreting' then 32
        when 'message digest calculation' then 33
        when 'metadata extraction' then 34
        when 'metadata modification' then 35
        when 'migration' then 36
        when 'modification' then 37
        when 'normalization' then 38
        when 'packing' then 39
        when 'policy assignment' then 40
        when 'printing' then 41
        when 'quarantine' then 42
        when 'recovery' then 43
        when 'redaction' then 44
        when 'refreshment' then 45
        when 'rendering' then 46
        when 'replication' then 47
        when 'transfer' then 48
        when 'unpacking' then 49
        when 'unquarantine' then 50
        when 'validation' then 51
        when 'virus check' then 52
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
create table if not exists event_type_lookup (
     id int primary key,
     event_type varchar not null
);

-- add foreign key restraint to event_type to map to event_type_lookup
alter table premis_events add constraint event_type_fk FOREIGN KEY event_type REFERENCES event_type_lookup(id)

-- lookup table for object in premis

alter table premis_events add COLUMN object_int smallint;
 
-- IMPORTANT --
-- TO DO: If there is a value in the current premis_events table --
-- for object that is NOT a match for any values in this function, --
-- probably we need to abort and roll back. If it converts to a 0, --
-- we will lose whatever information was in there. --
create or replace function convert_event_objects()
returns void as $$
begin
    update premis_events
    set object_int = case "object"
        when 'object' then 1

        else 0  -- default
    end;
end;
$$ language plpgsql;

select convert_event_objects();

-- if exists
alter table premis_events drop column "object";
alter table premis_events rename column object_int TO "object";


create table if not exists object_lookup (
     id int primary key,
     "object" varchar not null
);

alter table premis_events add constraint object_fk FOREIGN KEY "object" REFERENCES object_lookup(id)


-- lookup table for agent in premis

alter table premis_events add COLUMN agent_int smallint;
 
-- IMPORTANT --
-- TO DO: If there is a value in the current premis_events table --
-- for agent that is NOT a match for any values in this function, --
-- probably we need to abort and roll back. If it converts to a 0, --
-- we will lose whatever information was in there. --
create or replace function convert_event_agents()
returns void as $$
begin
    update premis_events
    set agent_int = case agent
        when 'agent' then 1
        else 0  -- default
    end;
end;
$$ language plpgsql;

select convert_event_agents();

-- if exists
alter table premis_events drop column agent;
alter table premis_events rename column agent_int TO agent;


create table if not exists agent_lookup (
     id int primary key,
     agent varchar not null
);

alter table premis_events add constraint agent_fk FOREIGN KEY agent REFERENCES agent_lookup(id)


-- checksums

alter table checksums drop column created_at;
alter table checksums drop column updated_at;

-- storage records url



-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '013_shrink_db_size';
