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
-- We need to remove the created_at and updated_at columns from the view first.
DROP VIEW premis_events_view;
CREATE VIEW public.premis_events_view
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

-- alter table premis_events drop column created_at; --
-- alter table premis_events drop column updated_at; --

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

insert into event_type_lookup (id, event_type) values 
(0, "No Event Type"),
(1, "access assignment"),
(2, "accession"),
(3, "appraisal"),
(4, "capture"),
(5, "compiling"),
(6, "compression"),
(7, "creation"),
(8, "deaccession"),
(9, "decompression"),
(10, "decryption"),
(11, "deletion"),
(12, "digital signature generation"),
(13, "digital signature validation"),
(14, "displaying"),
(15, "dissemination"),
(16, "encryption"),
(17, "execution"),
(18, "exporting"),
(19, "extraction"),
(20, "filename change"),
(21, "fixity check"),
(22, "forensic feature analysis"),
(23, "format identification"),
(24, "identifier assignment"),
(25, "imaging"),
(26, "information package creation"),
(27, "information package merging"),
(28, "information package splitting"),
(29, "ingestion"),
(30, "ingestion end"),
(31, "ingestion start"),
(32, "interpreting"),
(33, "message digest calculation"),
(34, "metadata extraction"),
(35, "metadata modification"),
(36, "migration"),
(37, "modification"),
(38, "normalization"),
(39, "packing"),
(40, "policy assignment"),
(41, "printing"),
(42, "quarantine"),
(43, "recovery"),
(44, "redaction"),
(45, "refreshment"),
(46, "rendering"),
(47, "replication"),
(48, "transfer"),
(49, "unpacking"),
(50, "unquarantine"),
(51, "validation"),
(52, "virus check");

-- add foreign key restraint to event_type to map to event_type_lookup
alter table premis_events add constraint event_type_fk FOREIGN KEY event_type REFERENCES event_type_lookup(id)

-- lookup table for object in premis

-- alter table premis_events add COLUMN object_int smallint;
 
-- IMPORTANT --
-- TO DO: If there is a value in the current premis_events table --
-- for object that is NOT a match for any values in this function, --
-- probably we need to abort and roll back. If it converts to a 0, --
-- we will lose whatever information was in there. --
/* create or replace function convert_event_objects()
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

*/

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '013_shrink_db_size';
