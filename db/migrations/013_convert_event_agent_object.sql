-- 013_convert_event_agent_object.sql
-- Creates lookup tables for agent and object fields of premis_events.
-- This allows us to save space in the database. Currently these fields are of type varchar.
-- But because of repetition in the data, we can convert these columns to type smallint and add lookup tables.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('013_convert_event_agent_object', now())
on conflict ("version") do update set started_at = now();

-- We'll need to drop and recreate the premis events view
drop view if exists public.premis_events_view;

-- Add the new columns.
alter table premis_events add COLUMN if not exists agent_int smallint;
alter table premis_events add COLUMN if not exists object_int smallint;

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
(5, 'https://github.com/APTrust/preservation-services'),
(6, 'APTrust preservation services'),
(7, 'http://github.com/google/uuid'),
(8, 'http://golang.org/pkg/crypto/sha256/'),
(9, 'http://golang.org/pkg/crypto/md5/'),
(10, 'Registry Unit Test'),
(11, 'Maxwell Smart'),
(12, 'https://github.com/APTrust/exchange'),
(13, 'http://github.com/satori/go.uuid'),
(14, 'http://github.com/nu7hatch/gouuid'),
(15, 'https://github.com/APTrust/bagman'),
(16, 'https://github.com/aws/aws-sdk-go'),
(17, 'https://github.com/crowdmob/goamz'),
(18, 'https://github.com/marcel/aws-s3/tree/master'),
(19, 'https://github.com/satori/go.uuid'),
(20, 'https://launchpad.net/goamz'),
(21, 'https://github.com/APTrust/auditing/blob/1.0/cleanup_001.py');

insert into event_object_lookup (id, event_object) values 
(0, 'unknown event object'),
(1, 'APTrust preservation services'),
(2, 'Minio S3 client'),
(3, 'Minio S3 library'),
(4, 'preservation-services + Minio S3 client'),
(5, 'Go uuid library + Minio S3 library')
(6, 'Go language crypto/sha256'),
(7, 'Go language crypto/md5'),
(8, 'scissors'),
(9, 'APTrust exchange/ingest processor'),
(10, 'Fake event object'),
(11, 'APTrust Go Exchange + Amazon S3 client'),
(12, 'SHA-256 thingy'),
(13, 'Exchange ingest code'),
(14, 'Deleterbot code'),
(15, 'APTrust bagman'),
(16, 'APTrust bag processor'),
(17, 'APTrust exchange'),
(18, 'APTrust Exchange apt_delete service'),
(19, 'APTrust Exchange ingest services'),
(20, 'APTrust exchange using Satori go.uuid'),
(21, 'AWS Go SDK S3 client'),
(22, 'AWS Go SDK S3 Library'),
(23, 'bagman + goamz s3 client'),
(24, 'exchange + AWS Go SDK S3 client'),
(25, 'exchange + goamz S3 client'),
(26, 'goamz S3 client'),
(27, 'Goamz S3 Client'),
(28, 'Go language cryptohash'),
(29, 'Go crypto/md5'),
(30, 'Go uuid library + APTrust DPN services'),
(31, 'Go uuid library + AWS Go SDK S3 library'),
(32, 'Go uuid library + goamz S3 library'),
(33, 'Ruby aws-s3 gem'),
(34, 'APTrust audit and cleanup scripts for audit_001');

-- IMPORTANT - Rollback if any agents or objects appear as 0
create or replace function convert_event_agents_and_objects()
returns void as $$
begin
    update premis_events set agent_int = case
        when agent='https://github.com/minio/minio-go v4' then 1
        when agent='https://github.com/minio/minio-go v5' then 2
        when agent='https://github.com/minio/minio-go v6' then 3
        when agent='https://github.com/minio/minio-go v7' then 4
        when agent='https://github.com/APTrust/preservation-services' then 5
        when agent='APTrust preservation services' then 6
        when agent='http://github.com/google/uuid' then 7
        when agent='http://golang.org/pkg/crypto/sha256/' then 8
        when agent='http://golang.org/pkg/crypto/md5/' then 9
        when agent='Registry Unit Test' then 10
        when agent='Maxwell Smart' then 11
        when agent='https://github.com/APTrust/exchange' then 12
        when agent='http://github.com/satori/go.uuid' then 13
        when agent='http://github.com/nu7hatch/gouuid' then 14
        when agent='https://github.com/APTrust/bagman' then 15
        when agent='https://github.com/aws/aws-sdk-go' then 16
        when agent='https://github.com/crowdmob/goamz' then 17
        when agent='https://github.com/marcel/aws-s3/tree/master' then 18
        when agent='https://github.com/satori/go.uuid' then 19
        when agent='https://launchpad.net/goamz' then 20
        when agent='https://github.com/APTrust/auditing/blob/1.0/cleanup_001.py' then 21
        else 0
    end,
    object_int = case
        when "object"='APTrust preservation services' then 1
        when "object"='Minio S3 client' then 2
        when "object"='Minio S3 library' then 3
        when "object"='preservation-services + Minio S3 client' then 4
        when "object"='Go uuid library + Minio S3 library' then 5
        when "object"='Go language crypto/sha256' then 6
        when "object"='Go language crypto/md5' then 7
        when "object"='scissors' then 8
        when "object"='APTrust exchange/ingest processor' then 9
        when "object"='Fake event object' then 10
        when "object"='APTrust Go Exchange + Amazon S3 client' then 11
        when "object"='SHA-256 thingy' then 12
        when "object"='Exchange ingest code' then 13
        when "object"='Deleterbot code' then 14
        when "object"='APTrust bagman' then 15
        when "object"='APTrust bag processor' then 16
        when "object"='APTrust exchange' then 17
        when "object"='APTrust Exchange apt_delete service' then 18
        when "object"='APTrust Exchange ingest services' then 19
        when "object"='APTrust exchange using Satori go.uuid' then 20
        when "object"='AWS Go SDK S3 client' then 21
        when "object"='AWS Go SDK S3 Library' then 22
        when "object"='bagman + goamz s3 client' then 23
        when "object"='exchange + AWS Go SDK S3 client' then 24
        when "object"='exchange + goamz S3 client' then 25
        when "object"='goamz S3 client' then 26
        when "object"='Goamz S3 Client' then 27
        when "object"='Go language cryptohash' then 28
        when "object"='Go crypto/md5' then 29
        when "object"='Go uuid library + APTrust DPN services' then 30
        when "object"='Go uuid library + AWS Go SDK S3 library' then 31
        when "object"='Go uuid library + goamz S3 library' then 32
        when "object"='Ruby aws-s3 gem' then 33
        when "object"='APTrust audit and cleanup scripts for audit_001' then 34
        else 0
    end;
end;
$$ language plpgsql;

-- Call function
select convert_event_agents_and_objects();

-- Drop original data and now rename the columns
alter table premis_events drop column if exists agent;
alter table premis_events drop column if exists "object";
alter table premis_events rename column event_agent_int TO agent;
alter table premis_events rename column event_object_int TO "object";

-- Add foreign key restraint to map columns to lookup tables
alter table premis_events add constraint event_agent_fk FOREIGN KEY (agent) REFERENCES event_agent_lookup(id);
alter table premis_events add constraint event_object_fk FOREIGN KEY ("object") REFERENCES event_object_lookup(id);

-- Recreate view
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
    pe.updated_at,
    pe.old_uuid
   FROM premis_events pe
     LEFT JOIN institutions i ON pe.institution_id = i.id
     LEFT JOIN intellectual_objects io ON pe.intellectual_object_id = io.id
     LEFT JOIN generic_files gf ON pe.generic_file_id = gf.id;

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '013_convert_event_agent_object';
