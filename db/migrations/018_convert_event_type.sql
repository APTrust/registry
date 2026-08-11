-- 018_convert_event_type.sql
--
-- This migration helps to shrink the size of our database by optimizing premis events.
-- By converting event_type from a string to an int, and indexing these int values,
-- we will save on space to store each event.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('018_convert_event_type', now())
on conflict ("version") do update set started_at = now();

-- First we need to drop related indices and views temporarily.
-- We will recreate them at the end.
drop index if exists index_premis_events_on_event_type;
drop index if exists index_premis_events_on_event_type_and_outcome;
drop index if exists ix_premis_event_counts;
drop materialized view public.premis_event_counts;
drop view if exists public.premis_events_view;

-- Add a new temporary column for the event_type as int
alter table premis_events add COLUMN if not exists event_type_int smallint;

-- Create a lookup table for these new event_type ints.
-- Each int is linked to a string describing the event type.
drop table if exists event_type_lookup;
create table if not exists event_type_lookup (
     id int primary key,
     event_type varchar not null
);

insert into event_type_lookup (id, event_type) values 
(0, 'unknown event'),
(1, 'access assignment'),
(2, 'accession'),
(3, 'appraisal'),
(4, 'capture'),
(5, 'compiling'),
(6, 'compression'),
(7, 'creation'),
(8, 'deaccession'),
(9, 'decompression'),
(10, 'decryption'),
(11, 'deletion'),
(12, 'digital signature generation'),
(13, 'digital signature validation'),
(14, 'displaying'),
(15, 'dissemination'),
(16, 'encryption'),
(17, 'execution'),
(18, 'exporting'),
(19, 'extraction'),
(20, 'filename change'),
(21, 'fixity check'),
(22, 'forensic feature analysis'),
(23, 'format identification'),
(24, 'identifier assignment'),
(25, 'imaging'),
(26, 'information package creation'),
(27, 'information package merging'),
(28, 'information package splitting'),
(29, 'ingestion'),
(30, 'ingestion end'),
(31, 'ingestion start'),
(32, 'interpreting'),
(33, 'message digest calculation'),
(34, 'metadata extraction'),
(35, 'metadata modification'),
(36, 'migration'),
(37, 'modification'),
(38, 'normalization'),
(39, 'packing'),
(40, 'policy assignment'),
(41, 'printing'),
(42, 'quarantine'),
(43, 'recovery'),
(44, 'redaction'),
(45, 'refreshment'),
(46, 'rendering'),
(47, 'replication'),
(48, 'transfer'),
(49, 'unpacking'),
(50, 'unquarantine'),
(51, 'validation'),
(52, 'virus check');
 
-- Function that performs the actual conversion.
create or replace function convert_event_types()
returns void as $$
begin
    update premis_events set event_type_int = case
        when event_type='access assignment' then 1
        when event_type='accession' then 2
        when event_type='appraisal' then 3
        when event_type='capture' then 4
        when event_type='compiling' then 5
        when event_type='compression' then 6
        when event_type='creation' then 7
        when event_type='deaccession' then 8
        when event_type='decompression' then 9
        when event_type='decryption' then 10
        when event_type='deletion' then 11
        when event_type='digital signature generation' then 12
        when event_type='digital signature validation' then 13
        when event_type='displaying' then 14
        when event_type='dissemination' then 15
        when event_type='encryption' then 16
        when event_type='execution' then 17
        when event_type='exporting' then 18
        when event_type='extraction' then 19
        when event_type='filename change' then 20
        when event_type='fixity check' then 21
        when event_type='forensic feature analysis' then 22
        when event_type='format identification' then 23
        when event_type='identifier assignment' then 24
        when event_type='imaging' then 25
        when event_type='information package creation' then 26
        when event_type='information package merging' then 27
        when event_type='information package splitting' then 28
        when event_type='ingestion' then 29
        when event_type='ingestion end' then 30
        when event_type='ingestion start' then 31
        when event_type='interpreting' then 32
        when event_type='message digest calculation' then 33
        when event_type='metadata extraction' then 34
        when event_type='metadata modification' then 35
        when event_type='migration' then 36
        when event_type='modification' then 37
        when event_type='normalization' then 38
        when event_type='packing' then 39
        when event_type='policy assignment' then 40
        when event_type='printing' then 41
        when event_type='quarantine' then 42
        when event_type='recovery' then 43
        when event_type='redaction' then 44
        when event_type='refreshment' then 45
        when event_type='rendering' then 46
        when event_type='replication' then 47
        when event_type='transfer' then 48
        when event_type='unpacking' then 49
        when event_type='unquarantine' then 50
        when event_type='validation' then 51
        when event_type='virus check' then 52
        else 0  -- default
    end;
end;
$$ language plpgsql;

-- Run the function.
select convert_event_types();

-- Now that we have converted every event, we can drop the old event_type column.
alter table premis_events drop column if exists event_type;

-- Rename the event_type_int column to event_type.
alter table premis_events rename column if exists event_type_int TO event_type;

-- Add foreign key restraint to event_type to map to event_type_lookup
alter table premis_events add constraint event_type_fk FOREIGN KEY (event_type) REFERENCES event_type_lookup(id);

-- Recreate indices and views that use this and reindex - may take some time before indexing is complete
CREATE INDEX if not exists index_premis_events_on_event_type ON public.premis_events USING btree (event_type);
CREATE INDEX if not exists index_premis_events_on_event_type_and_outcome ON public.premis_events USING btree (event_type, outcome);

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

CREATE if not exists UNIQUE INDEX ix_premis_event_counts ON public.premis_event_counts USING btree (institution_id, event_type, outcome);

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

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '018_convert_event_type';
