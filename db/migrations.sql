-- migrations.sql
--
-- This file contains ALL alterations that should be applied to the
-- existing Pharos DB schema (as it exists in the feature/storage-option
-- branch) to make it match schema.sql.
--
-- All operations in this file must be idempotent, so we can run it
-- any number of times and always know that it will leave the DB in a
-- consistent and known state that matches schema.sql.
--
-- NOTE: When migrating old Pharos DB, we will also need to create the
--       views in the schema.sql file.
-------------------------------------------------------------------------------

-- We need to fix the user role structure. Pharos allows a user to have
-- multiple roles at a single institution, though our business rules disallow
-- that, and no user has ever had more than one role. To simplify the DB
-- and our queries, we need to do the following:
--
-- 1. Create a role column in the users table.
-- 2. Populate that column with the value with each user's role from
--    user -> roles_user -> roles.
-- 3. Drop the roles_users table.
-- 4. Drop the roles table.

do $$
begin
  if not exists (select 1 from information_schema.columns where table_schema='public' AND table_name='users' AND column_name='role') then
 	alter table users add column "role" varchar(50) not null default 'none';
 	update users u set "role" = coalesce((select r.name from "roles" r inner join roles_users ru on ru.role_id = r.id where ru.user_id = u.id), 'none');
    drop table if exists roles_users;
    drop table if exists roles;
  end if;
end
$$;

-- The ingest_state columns were part of a proposed architecture we never
-- implemented. They have never been used, and we don't need them.
alter table intellectual_objects drop column if exists ingest_state;
alter table generic_files drop column if exists ingest_state;

-- Remove object_identifier and generic_file_identifier from work_items.
-- We can use a view to join the files & objects tables, avoiding the
-- duplicate data.
alter table work_items drop column if exists object_identifier;
alter table work_items drop column if exists generic_file_identifier;

-- The work_items.date column actually refers to the datetime on which
-- the item was last processed by one of our Go workers. The name "date"
-- is too vague and ambiguous, so let's call it what it is.
do $$
begin
  if exists (select 1 from information_schema.columns
    where table_schema='public' AND table_name='work_items' AND column_name='date') then
    alter table work_items rename column "date" to "date_processed";
  end if;
end
$$;

-- Update the indexes on work_items to reflect the change from date to
-- date_processed.

drop index if exists index_work_items_on_date;
drop index if exists index_work_items_on_institution_id_and_date;

create index if not exists index_work_items_on_date_processed on work_items(date_processed);
create index if not exists index_work_items_on_inst_id_and_date_processed on work_items(institution_id, date_processed);

-- These columns are unnecessary in the premis_events table.
-- We can get them by joining other tables in a view.
-- This will have a large impact on storage space, since these
-- two fields together average ~100 bytes per record, and as of
-- Jan. 2021, we have around 120M records. Dropping these columns
-- also drops the indexes on these two columns, which are huge and
-- are rarely used.
alter table premis_events drop column if exists intellectual_object_identifier;
alter table premis_events drop column if exists generic_file_identifier;

-- Get rid of useless indexes. We can add these back if they actually turn
-- out to be useful.
--
-- generic_files
drop index if exists index_files_on_inst_state_and_format;
drop index if exists index_files_on_inst_state_and_updated;
drop index if exists index_generic_files_on_file_format;
drop index if exists index_generic_files_on_file_format_and_state;
drop index if exists index_generic_files_on_institution_id_and_size_and_state;
drop index if exists index_generic_files_on_intellectual_object_id_and_file_format;
drop index if exists index_generic_files_on_intellectual_object_id_and_state;
drop index if exists index_generic_files_on_size;
drop index if exists index_generic_files_on_size_and_state;
drop index if exists index_generic_files_on_state;
drop index if exists index_generic_files_on_state_and_updated_at;

-- intellectual_objects
drop index if exists index_intellectual_objects_on_access;
drop index if exists index_intellectual_objects_on_institution_id_and_state;
drop index if exists index_intellectual_objects_on_state;

-- premis_events
drop index if exists index_premis_events_on_generic_file_id_and_event_type;
drop index if exists index_premis_events_on_identifier_and_institution_id;

-- institutions
create unique index if not exists index_institutions_identifier on public.institutions using btree(identifier);
create unique index if not exists index_institutions_receiving_bucket on public.institutions using btree(receiving_bucket);
create unique index if not exists index_institutions_restore_bucket on public.institutions using btree(restore_bucket);
