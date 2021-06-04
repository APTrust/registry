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
create unique index if not exists index_institutions_identifier on institutions using btree(identifier);
create unique index if not exists index_institutions_receiving_bucket on institutions using btree(receiving_bucket);
create unique index if not exists index_institutions_restore_bucket on institutions using btree(restore_bucket);

-- Allow institutions to turn spot restoration tests on or off.
-- Default is off.
alter table institutions add column if not exists enable_spot_restore boolean not null default false;

-- ********************************************************************
--
-- Add tables for alerts and deletions that did not exist in Pharos
--
-- ********************************************************************

-- deletion_requests track requests for file and object deletions,
-- who initiated them and who approved them
create table if not exists public.deletion_requests (
	id bigserial primary key,
	institution_id integer not null references public.institutions(id),
	requested_by_id integer not null references public.users(id),
	requested_at timestamp not null,
	encrypted_confirmation_token varchar not null,
	confirmed_by_id integer references public.users(id),
	confirmed_at timestamp,
	cancelled_by_id integer references public.users(id),
	cancelled_at timestamp,
    work_item_id integer null references public.work_items(id)
);
create index if not exists index_deletion_requests_institution_id ON public.deletion_requests (institution_id);


-- deletion_requests_generic_files records which files belong to a deletion request
create table if not exists public.deletion_requests_generic_files (
	deletion_request_id integer not null references public.deletion_requests(id),
	generic_file_id integer not null references public.generic_files(id)
);
create index if not exists index_drgf_deletion_request_id ON public.deletion_requests_generic_files (deletion_request_id);


-- deletion_requests_intellectual_objects records which objects belong to a deletion request
create table if not exists public.deletion_requests_intellectual_objects (
	deletion_request_id integer not null references public.deletion_requests(id),
	intellectual_object_id integer not null references public.intellectual_objects(id)
);
create index if not exists index_drio_deletion_request_id ON public.deletion_requests_intellectual_objects (deletion_request_id);


-- alerts stores the content of alert messages. These messages appear in the web UI
-- and may also be emailed to users, depending on the alert type.
-- Column deletion_request_id will typically be null.
create table if not exists public.alerts (
	id bigserial primary key,
	institution_id integer references public.institutions(id),
	"type" varchar not null,
	"content" text not null,
	deletion_request_id integer references public.deletion_requests(id),
	created_at timestamp not null
);
create index if not exists index_alerts_institution_id ON public.alerts (institution_id);
create index if not exists index_alerts_type ON public.alerts ("type");


-- alerts_users tracks which users should see which alerts. Depending on alert.type,
-- the message may be emailed to the user (for example, a deletion approval alert),
-- or it may simply be displayed in the web UI.
create table if not exists public.alerts_users (
	alert_id integer not null references public.alerts(id),
	user_id integer not null references public.users(id),
	sent_at timestamp default null,
	read_at timestamp default null
);
create index if not exists index_alerts_users_alert_id ON public.alerts_users (alert_id);
create index if not exists index_alerts_users_user_id ON public.alerts_users (user_id);
create unique index if not exists index_alerts_users_unique ON public.alerts_users (alert_id, user_id);


-- alerts_work_items link an alerts to one or more work items.
create table if not exists public.alerts_work_items (
	alert_id integer not null references public.alerts(id),
	work_item_id integer not null references public.work_items(id)
);
create index if not exists index_alerts_work_items_alert_id ON public.alerts_work_items (alert_id);
create unique index if not exists index_alerts_work_items_unique ON public.alerts_work_items (alert_id, work_item_id);


-- alerts_premis_events link alerts to one or more premis_events.
create table if not exists public.alerts_premis_events (
	alert_id integer not null references public.alerts(id),
	premis_event_id integer not null references public.premis_events(id)
);
create index if not exists index_alerts_premis_events_alert_id ON public.alerts_premis_events(alert_id);
create unique index if not exists index_alerts_premis_events_unique ON public.alerts_premis_events(alert_id, premis_event_id);
