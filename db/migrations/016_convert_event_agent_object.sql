-- 016_convert_event_agent_object.sql
--
-- Creates lookup tables for agent and object fields of premis_events.
-- This allows us to save space in the database. Currently these fields are of type varchar.
-- But because of repetition in the data, we can convert these columns to type smallint and add lookup tables.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('016_convert_event_agent_object', now())
on conflict ("version") do update set started_at = now();

-- We'll need to drop and recreate premis_event dependent objects
drop index index_premis_events_on_event_type;
drop index index_premis_events_on_event_type_and_outcome;
drop index ix_premis_event_counts;
drop materialized view public.premis_event_counts;
drop view public.premis_events_view;

-- Add the new columns.
alter table premis_events add COLUMN agent_int smallint;
alter table premis_events add COLUMN object_int smallint;

-- Create lookup tables.
drop table if exists event_agent_lookup;
create table event_agent_lookup (
     id int primary key,
     event_agent varchar not null
);
drop table if exists event_object_lookup;
create table event_object_lookup (
     id int primary key,
     event_object varchar not null
);

-- Add possible values for agent and object on events.
insert into event_agent_lookup (id, event_agent) values 
(0, 'unknown event agent'),
(1, 'https://github.com/minio/minio-go v4'),
(2, 'https://github.com/minio/minio-go v5'),
(3, 'https://github.com/minio/minio-go v6'),
(4, 'https://github.com/minio/minio-go v7'),
(5, 'Registry Unit Test'),
(6, 'APTrust preservation services');

insert into event_object_lookup (id, event_object) values 
(0, 'unknown event object'),
(1, 'preservation-services + Minio S3 client'),
(2, 'Minio S3 client'),
(3, 'APTrust preservation services'),
(4, 'Go uuid library + Minio S3 library')
(5, 'Go language crypto/sha256'),
(6, 'Minio S3 library');
 
-- IMPORTANT - Rollback if any agents or objects appear as 0
create or replace function convert_event_agents()
returns void as $$
begin
    update premis_events set event_agent_int = case
        when agent='https://github.com/minio/minio-go v4' then 1
        when agent='https://github.com/minio/minio-go v5' then 2
        when agent='https://github.com/minio/minio-go v6' then 3
        when agent='https://github.com/minio/minio-go v7' then 4
        when agent='Registry Unit Test' then 5
        when agent='APTrust preservation services' then 6
        else 0  -- default
    end;
end;
$$ language plpgsql;

create or replace function convert_event_objects()
returns void as $$
begin
    update premis_events set event_object_int = case
        when "object"='preservation-services + Minio S3 client' then 1
        when "object"='Minio S3 client' then 2
        when "object"='APTrust preservation services' then 3
        when "object"='Go uuid library + Minio S3 library' then 4
        when "object"='Go language crypto/sha256' then 5
        when "object"='Minio S3 library' then 6
        else 0  -- default
    end;
end;
$$ language plpgsql;

-- Call functions
select convert_event_agents();
select convert_event_objects();

-- if exists
alter table premis_events drop column agent;
alter table premis_events drop column "object";
alter table premis_events rename column event_agent_int TO agent;
alter table premis_events rename column event_object_int TO "object";

-- add foreign key restraint to map columns to lookup tables
alter table premis_events add constraint event_agent_fk FOREIGN KEY (agent) REFERENCES event_agent_lookup(id);
alter table premis_events add constraint event_object_fk FOREIGN KEY ("object") REFERENCES event_object_lookup(id);

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
update schema_migrations set finished_at = now() where "version" = '016_convert_event_agent_object';
