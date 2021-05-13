--
-- schema.sql
--
-- This file represents the full registry schema.
--
-- The migrations.sql file updates the old Pharos schema, as it existed
-- in the feature/storage_record branch, so that it matches this schema.
--
-- https://github.com/APTrust/pharos/blob/feature/storage-record/db/schema.rb
--


-- public.ar_internal_metadata definition
--
-- This is a legacy table used by ActiveRecord.
-- We can drop this once production is stable.
--
CREATE TABLE public.ar_internal_metadata (
	"key" varchar NOT NULL,
	value varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key)
);


-- public.bulk_delete_jobs definition
--
-- This is the legacy bulk delete table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.bulk_delete_jobs (
	id bigserial NOT NULL,
	requested_by varchar NULL,
	institutional_approver varchar NULL,
	aptrust_approver varchar NULL,
	institutional_approval_at timestamp NULL,
	aptrust_approval_at timestamp NULL,
	note text NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	institution_id int4 NOT NULL,
	CONSTRAINT bulk_delete_jobs_pkey PRIMARY KEY (id)
);


-- public.bulk_delete_jobs_emails definition
--
-- Related to the legacy bulk delete table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.bulk_delete_jobs_emails (
	bulk_delete_job_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_emails_on_bulk_delete_job_id ON public.bulk_delete_jobs_emails USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_emails_on_email_id ON public.bulk_delete_jobs_emails USING btree (email_id);


-- public.bulk_delete_jobs_generic_files definition
--
-- Related to the legacy bulk delete table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.bulk_delete_jobs_generic_files (
	bulk_delete_job_id int8 NULL,
	generic_file_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_generic_files_on_bulk_delete_job_id ON public.bulk_delete_jobs_generic_files USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_generic_files_on_generic_file_id ON public.bulk_delete_jobs_generic_files USING btree (generic_file_id);


-- public.bulk_delete_jobs_institutions definition
--
-- Related to the legacy bulk delete table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.bulk_delete_jobs_institutions (
	bulk_delete_job_id int8 NULL,
	institution_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_institutions_on_bulk_delete_job_id ON public.bulk_delete_jobs_institutions USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_institutions_on_institution_id ON public.bulk_delete_jobs_institutions USING btree (institution_id);


-- public.bulk_delete_jobs_intellectual_objects definition
--
-- Related to the legacy bulk delete table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.bulk_delete_jobs_intellectual_objects (
	bulk_delete_job_id int8 NULL,
	intellectual_object_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_intellectual_objects_on_bulk_job_id ON public.bulk_delete_jobs_intellectual_objects USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_intellectual_objects_on_object_id ON public.bulk_delete_jobs_intellectual_objects USING btree (intellectual_object_id);


-- public.confirmation_tokens definition
--
-- This is a legacy table. We can likely delete it
-- after production is stable.
--
CREATE TABLE public.confirmation_tokens (
	id bigserial NOT NULL,
	"token" varchar NULL,
	intellectual_object_id int4 NULL,
	generic_file_id int4 NULL,
	institution_id int4 NULL,
	user_id int4 NULL,
	CONSTRAINT confirmation_tokens_pkey PRIMARY KEY (id)
);


-- public.emails definition
--
-- This legacy table will be replaced by the new alerts table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.emails (
	id bigserial NOT NULL,
	email_type varchar NULL,
	event_identifier varchar NULL,
	item_id int4 NULL,
	email_text text NULL,
	user_list text NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	intellectual_object_id int4 NULL,
	generic_file_id int4 NULL,
	institution_id int4 NULL,
	CONSTRAINT emails_pkey PRIMARY KEY (id)
);


-- public.emails_generic_files definition
--
-- Related to the legacy emails table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.emails_generic_files (
	generic_file_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_generic_files_on_email_id ON public.emails_generic_files USING btree (email_id);
CREATE INDEX index_emails_generic_files_on_generic_file_id ON public.emails_generic_files USING btree (generic_file_id);


-- public.emails_intellectual_objects definition
--
-- Related to the legacy emails table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.emails_intellectual_objects (
	intellectual_object_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_intellectual_objects_on_email_id ON public.emails_intellectual_objects USING btree (email_id);
CREATE INDEX index_emails_intellectual_objects_on_intellectual_object_id ON public.emails_intellectual_objects USING btree (intellectual_object_id);


-- public.emails_premis_events definition
--
-- Related to the legacy emails table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.emails_premis_events (
	premis_event_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_premis_events_on_email_id ON public.emails_premis_events USING btree (email_id);
CREATE INDEX index_emails_premis_events_on_premis_event_id ON public.emails_premis_events USING btree (premis_event_id);


-- public.emails_work_items definition
--
-- Related to the legacy emails table.
-- We should eventually archive this, but keep the
-- data somewhere for historical and audit purposes.
--
CREATE TABLE public.emails_work_items (
	work_item_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_work_items_on_email_id ON public.emails_work_items USING btree (email_id);
CREATE INDEX index_emails_work_items_on_work_item_id ON public.emails_work_items USING btree (work_item_id);


-- public.generic_files definition
CREATE TABLE public.generic_files (
	id serial NOT NULL,
	file_format varchar NULL,
	"size" int8 NULL,
	identifier varchar NULL,
	intellectual_object_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	state varchar NULL,
	last_fixity_check timestamp NOT NULL DEFAULT '2000-01-01 00:00:00'::timestamp without time zone,
	institution_id int4 NOT NULL,
	storage_option varchar NOT NULL DEFAULT 'Standard'::character varying,
	uuid varchar NOT NULL,
	CONSTRAINT generic_files_pkey PRIMARY KEY (id)
);
CREATE INDEX index_generic_files_on_created_at ON public.generic_files USING btree (created_at);
CREATE UNIQUE INDEX index_generic_files_on_identifier ON public.generic_files USING btree (identifier);
CREATE INDEX index_generic_files_on_institution_id ON public.generic_files USING btree (institution_id);
CREATE INDEX index_generic_files_on_institution_id_and_state ON public.generic_files USING btree (institution_id, state);
CREATE INDEX index_generic_files_on_institution_id_and_updated_at ON public.generic_files USING btree (institution_id, updated_at);
CREATE INDEX index_generic_files_on_intellectual_object_id ON public.generic_files USING btree (intellectual_object_id);
CREATE INDEX index_generic_files_on_updated_at ON public.generic_files USING btree (updated_at);
CREATE UNIQUE INDEX index_generic_files_on_uuid ON public.generic_files USING btree (uuid);
CREATE INDEX ix_gf_last_fixity_check ON public.generic_files USING btree (last_fixity_check);


-- public.institutions definition
CREATE TABLE public.institutions (
	id serial NOT NULL,
	"name" varchar NULL,
	identifier varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	state varchar NULL,
	"type" varchar NULL,
	member_institution_id int4 NULL,
	deactivated_at timestamp NULL,
	otp_enabled bool NULL,
    enable_spot_restore bool not null default false,
	receiving_bucket varchar NOT NULL,
	restore_bucket varchar NOT NULL,
	CONSTRAINT institutions_pkey PRIMARY KEY (id)
);
CREATE INDEX index_institutions_on_name ON public.institutions USING btree (name);


-- public.intellectual_objects definition
CREATE TABLE public.intellectual_objects (
	id serial NOT NULL,
	title varchar NULL,
	description text NULL,
	identifier varchar NULL,
	alt_identifier varchar NULL,
	"access" varchar NULL,
	bag_name varchar NULL,
	institution_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	state varchar NULL,
	etag varchar NULL,
	bag_group_identifier varchar NOT NULL DEFAULT ''::character varying,
	storage_option varchar NOT NULL DEFAULT 'Standard'::character varying,
	bagit_profile_identifier varchar NULL,
	source_organization varchar NULL,
	internal_sender_identifier varchar NULL,
	internal_sender_description text NULL,
	CONSTRAINT intellectual_objects_pkey PRIMARY KEY (id)
);
CREATE INDEX index_intellectual_objects_on_bag_name ON public.intellectual_objects USING btree (bag_name);
CREATE INDEX index_intellectual_objects_on_created_at ON public.intellectual_objects USING btree (created_at);
CREATE UNIQUE INDEX index_intellectual_objects_on_identifier ON public.intellectual_objects USING btree (identifier);
CREATE INDEX index_intellectual_objects_on_institution_id ON public.intellectual_objects USING btree (institution_id);
CREATE INDEX index_intellectual_objects_on_updated_at ON public.intellectual_objects USING btree (updated_at);


-- public.old_passwords definition
--
-- This is a legacy table. We can get rid of this if we
-- abandon mandatory password changes. Wait until that's
-- decided.
--
CREATE TABLE public.old_passwords (
	id bigserial NOT NULL,
	encrypted_password varchar NOT NULL,
	password_salt varchar NULL,
	password_archivable_type varchar NOT NULL,
	password_archivable_id int4 NOT NULL,
	created_at timestamp NULL,
	CONSTRAINT old_passwords_pkey PRIMARY KEY (id)
);
CREATE INDEX index_password_archivable ON public.old_passwords USING btree (password_archivable_type, password_archivable_id);


-- public.premis_events definition
CREATE TABLE public.premis_events (
	id serial NOT NULL,
	identifier varchar NULL,
	event_type varchar NULL,
	date_time varchar NULL,
	outcome_detail varchar NULL,
	detail varchar NULL,
	outcome_information varchar NULL,
	"object" varchar NULL,
	agent varchar NULL,
	intellectual_object_id int4 NULL,
	generic_file_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	outcome varchar NULL,
	institution_id int4 NULL,
	old_uuid varchar NULL,
	CONSTRAINT premis_events_pkey PRIMARY KEY (id)
);
CREATE INDEX index_premis_events_date_time_desc ON public.premis_events USING btree (date_time DESC);
CREATE INDEX index_premis_events_on_event_type ON public.premis_events USING btree (event_type);
CREATE INDEX index_premis_events_on_event_type_and_outcome ON public.premis_events USING btree (event_type, outcome);
CREATE INDEX index_premis_events_on_generic_file_id ON public.premis_events USING btree (generic_file_id);
CREATE UNIQUE INDEX index_premis_events_on_identifier ON public.premis_events USING btree (identifier);
CREATE INDEX index_premis_events_on_institution_id ON public.premis_events USING btree (institution_id);
CREATE INDEX index_premis_events_on_intellectual_object_id ON public.premis_events USING btree (intellectual_object_id);
CREATE INDEX index_premis_events_on_outcome ON public.premis_events USING btree (outcome);


-- public.schema_migrations definition
--
-- This is a legacy table used by ActiveRecord
-- to track which migrations have been applied.
--
-- We can delete this when production is stable.
--
CREATE TABLE public.schema_migrations (
	"version" varchar NOT NULL,
	CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);


-- public.snapshots definition
--
-- This table stores point-in-time snapshots
-- of how much data each depositor has in the repo.
--
-- We may eventually delete this, but keep for now.
--
CREATE TABLE public.snapshots (
	id bigserial NOT NULL,
	audit_date timestamp NULL,
	institution_id int4 NULL,
	apt_bytes int8 NULL,
	"cost" numeric NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	snapshot_type varchar NULL,
	cs_bytes int8 NULL,
	go_bytes int8 NULL,
	CONSTRAINT snapshots_pkey PRIMARY KEY (id)
);


-- public.usage_samples definition
--
-- No idea what this was supposed to be, but the
-- table never had any data in the production DB.
-- Consider deleting this after we're sure it's not used.
--
CREATE TABLE public.usage_samples (
	id serial NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	institution_id varchar NULL,
	"data" text NULL,
	CONSTRAINT usage_samples_pkey PRIMARY KEY (id)
);


-- public.work_items definition
CREATE TABLE public.work_items (
	id serial NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	intellectual_object_id int4 NULL,
	generic_file_id int4 NULL,
	"name" varchar NULL,
	etag varchar NULL,
	bucket varchar NULL,
	"user" varchar NULL,
	note text NULL,
	"action" varchar NULL,
	stage varchar NULL,
	status varchar NULL,
	outcome text NULL,
	bag_date timestamp NULL,
	date_processed timestamp NULL,
	retry bool NOT NULL DEFAULT false,
	node varchar(255) NULL,
	pid int4 NULL DEFAULT 0,
	needs_admin_review bool NOT NULL DEFAULT false,
	institution_id int4 NULL,
	queued_at timestamp NULL,
	"size" int8 NULL,
	stage_started_at timestamp NULL,
	aptrust_approver varchar NULL,
	inst_approver varchar NULL,
	CONSTRAINT work_items_pkey PRIMARY KEY (id)
);
CREATE INDEX index_work_items_on_action ON public.work_items USING btree (action);
CREATE INDEX index_work_items_on_date ON public.work_items USING btree (date_processed);
CREATE INDEX index_work_items_on_etag_and_name ON public.work_items USING btree (etag, name);
CREATE INDEX index_work_items_on_generic_file_id ON public.work_items USING btree (generic_file_id);
CREATE INDEX index_work_items_on_institution_id ON public.work_items USING btree (institution_id);
CREATE INDEX index_work_items_on_institution_id_and_date ON public.work_items USING btree (institution_id, date_processed);
CREATE INDEX index_work_items_on_intellectual_object_id ON public.work_items USING btree (intellectual_object_id);
CREATE INDEX index_work_items_on_stage ON public.work_items USING btree (stage);
CREATE INDEX index_work_items_on_status ON public.work_items USING btree (status);


-- public.checksums definition
CREATE TABLE public.checksums (
	id serial NOT NULL,
	algorithm varchar NULL,
	datetime timestamp NULL,
	digest varchar NULL,
	generic_file_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT checksums_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_89bb0866e7 FOREIGN KEY (generic_file_id) REFERENCES generic_files(id)
);
CREATE INDEX index_checksums_on_generic_file_id ON public.checksums USING btree (generic_file_id);


-- public.storage_records definition
CREATE TABLE public.storage_records (
	id bigserial NOT NULL,
	generic_file_id int4 NULL,
	url varchar NULL,
	CONSTRAINT storage_records_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_a126ea6adc FOREIGN KEY (generic_file_id) REFERENCES generic_files(id)
);
CREATE INDEX index_storage_records_on_generic_file_id ON public.storage_records USING btree (generic_file_id);


-- public.users definition
CREATE TABLE public.users (
	id serial NOT NULL,
	"name" varchar NULL,
	email varchar NULL,
	phone_number varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	encrypted_password varchar NOT NULL DEFAULT ''::character varying,
	reset_password_token varchar NULL,
	reset_password_sent_at timestamp NULL,
	remember_created_at timestamp NULL,
	sign_in_count int4 NOT NULL DEFAULT 0,
	current_sign_in_at timestamp NULL,
	last_sign_in_at timestamp NULL,
	current_sign_in_ip varchar NULL,
	last_sign_in_ip varchar NULL,
	institution_id int4 NULL,
	encrypted_api_secret_key text NULL,
	password_changed_at timestamp NULL,
	encrypted_otp_secret varchar NULL,
	encrypted_otp_secret_iv varchar NULL,
	encrypted_otp_secret_salt varchar NULL,
	consumed_timestep int4 NULL,
	otp_required_for_login bool NULL,
	deactivated_at timestamp NULL,
	enabled_two_factor bool NULL DEFAULT false,
	confirmed_two_factor bool NULL DEFAULT false,
	otp_backup_codes _varchar NULL,
	authy_id varchar NULL,
	last_sign_in_with_authy timestamp NULL,
	authy_status varchar NULL,
	email_verified bool NULL DEFAULT false,
	initial_password_updated bool NULL DEFAULT false,
	force_password_update bool NULL DEFAULT false,
	account_confirmed bool NULL DEFAULT true,
	grace_period timestamp NULL,
    "role" varchar(50) NOT NULL default 'none',
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_7fcf39ca13 FOREIGN KEY (institution_id) REFERENCES institutions(id)
);
CREATE INDEX index_users_on_authy_id ON public.users USING btree (authy_id);
CREATE UNIQUE INDEX index_users_on_email ON public.users USING btree (email);
CREATE INDEX index_users_on_institution_id ON public.users USING btree (institution_id);
CREATE INDEX index_users_on_password_changed_at ON public.users USING btree (password_changed_at);
CREATE UNIQUE INDEX index_users_on_reset_password_token ON public.users USING btree (reset_password_token);

-------------------------------------------------------------------------------
-- New tables that were not part of Pharos.
-- alerts              + join tables
-- deletion_requests   + join tables
-------------------------------------------------------------------------------

-- deletion_requests track requests for file and object deletions,
-- who initiated them and who approved them
create table public.deletion_requests (
	id bigserial primary key,
	institution_id integer not null references public.institutions(id),
	requested_by_id integer not null references public.users(id),
	requested_at timestamp not null,
	encrypted_confirmation_token varchar not null,
	encrypted_cancellation_token varchar not null,
	confirmed_by_id integer references public.users(id),
	confirmed_at timestamp,
	cancelled_by_id integer references public.users(id),
	cancelled_at timestamp
);
create index index_deletion_requests_institution_id ON public.deletion_requests (institution_id);


-- deletion_requests_generic_files records which files belong to a deletion request
create table public.deletion_requests_generic_files (
	deletion_request_id integer not null references public.deletion_requests(id),
	generic_file_id integer not null references public.generic_files(id)
);
create index index_drgf_deletion_request_id ON public.deletion_requests_generic_files (deletion_request_id);


-- deletion_requests_intellectual_objects records which objects belong to a deletion request
create table public.deletion_requests_intellectual_objects (
	deletion_request_id integer not null references public.deletion_requests(id),
	intellectual_object_id integer not null references public.intellectual_objects(id)
);
create index index_drio_deletion_request_id ON public.deletion_requests_intellectual_objects (deletion_request_id);


-- alerts stores the content of alert messages. These messages appear in the web UI
-- and may also be emailed to users, depending on the alert type.
-- Column deletion_request_id will typically be null.
create table public.alerts (
	id bigserial primary key,
	institution_id integer references public.institutions(id),
	"type" varchar not null,
	"content" text not null,
	deletion_request_id integer references public.deletion_requests(id),
	created_at timestamp not null
);
create index index_alerts_institution_id ON public.alerts (institution_id);
create index index_alerts_type ON public.alerts ("type");


-- alerts_users tracks which users should see which alerts. Depending on alert.type,
-- the message may be emailed to the user (for example, a deletion approval alert),
-- or it may simply be displayed in the web UI.
create table public.alerts_users (
	alert_id integer not null references public.alerts(id),
	user_id integer not null references public.users(id),
	sent_at timestamp default null,
	read_at timestamp default null
);
create index index_alerts_users_alert_id ON public.alerts_users (alert_id);
create index index_alerts_users_user_id ON public.alerts_users (user_id);


-- alerts_work_items link an alerts to one or more work items.
create table public.alerts_work_items (
	alert_id integer not null references public.alerts(id),
	work_item_id integer not null references public.work_items(id)
);
create index index_alerts_work_items_alert_id ON public.alerts_work_items (alert_id);


-- alerts_premis_events link alerts to one or more premis_events.
create table public.alerts_premis_events (
	alert_id integer not null references public.alerts(id),
	premis_event_id integer not null references public.premis_events(id)
);
create index index_alerts_premis_events_alert_id ON public.alerts_premis_events(alert_id);




-------------------------------------------------------------------------------
-- Views
-------------------------------------------------------------------------------

-- Create a view to pull in columns that used to be in the WorkItems table.
-- These include intellectual_object_identifier and generic_file_identifier.
create or replace view work_items_view as
select
	wi.id,
	wi.institution_id,
	i."name" as institution_name,
	i.identifier as institution_identifier,
	wi.intellectual_object_id,
	io.identifier as "object_identifier",
	io.alt_identifier,
	io.bag_group_identifier,
	io.storage_option,
	io.bagit_profile_identifier,
	io.source_organization,
	io.internal_sender_identifier,
	wi.generic_file_id,
	gf.identifier as "generic_file_identifier",
	wi."name",
	wi.etag,
	wi.bucket,
	wi."user",
	wi.note,
	wi."action",
	wi.stage,
	wi.status,
	wi.outcome,
	wi.bag_date,
	wi.date_processed,
	wi.retry,
	wi.node,
	wi.pid,
	wi.needs_admin_review,
	wi."size",
	wi.queued_at,
	wi.stage_started_at,
	wi.aptrust_approver,
	wi.inst_approver,
	wi.created_at,
	wi.updated_at
from work_items wi
left join institutions i on wi.institution_id = i.id
left join intellectual_objects io on wi.intellectual_object_id = io.id
left join generic_files gf on wi.generic_file_id = gf.id;


-- Create a view to pull in columns that used to be in the PremisEvents table.
-- These include intellectual_object_identifier and generic_file_identifier.
create or replace view premis_events_view as
select
	pe.id,
	pe.identifier,
	pe.institution_id,
	i."name" as institution_name,
	pe.intellectual_object_id,
	io.identifier as intellectual_object_identifier,
	pe.generic_file_id,
	gf.identifier as generic_file_identifier,
	pe.event_type,
	pe.date_time,
	pe.detail,
	pe.outcome,
	pe.outcome_detail,
	pe.outcome_information,
	pe."object",
	pe.agent,
	pe.created_at,
	pe.updated_at,
	pe.old_uuid
from premis_events pe
left join institutions i on pe.institution_id = i.id
left join intellectual_objects io on pe.intellectual_object_id = io.id
left join generic_files gf on pe.generic_file_id = gf.id;


-- users_view makes it easier to list and search on institution-related
-- attributes for users.
create or replace view users_view as
select
	u.id,
	u."name",
	u.email,
	u.phone_number,
	u.created_at,
	u.updated_at,
	u.reset_password_sent_at,
	u.remember_created_at,
	u.sign_in_count,
	u.current_sign_in_at,
	u.last_sign_in_at,
	u.current_sign_in_ip,
	u.last_sign_in_ip,
	u.institution_id,
	u.password_changed_at,
	u.consumed_timestep,
	u.otp_required_for_login,
	u.deactivated_at,
	u.enabled_two_factor,
	u.confirmed_two_factor,
	u.authy_id,
	u.last_sign_in_with_authy,
	u.authy_status,
	u.email_verified,
	u.initial_password_updated,
	u.force_password_update,
	u.account_confirmed,
	u.grace_period,
	u."role",
	i."name" as institution_name,
	i.identifier as institution_identifier,
	i.state as institution_state,
	i."type" as institution_type,
	i.member_institution_id,
	i2."name" as member_institution_name,
	i2.identifier as member_institution_identifier,
	i2.state as member_institution_state,
	i.otp_enabled,
	i.receiving_bucket,
	i.restore_bucket
from users u
left join institutions i on u.institution_id = i.id
left join institutions i2 on i.member_institution_id = i2.id;


-- institutions_view shows an institution along with essential
-- information about its parent, if it has a parent.

create or replace view institutions_view as
select
	i.id,
	i."name",
	i.identifier,
	i.state,
	i."type",
	i.deactivated_at,
	i.otp_enabled,
	i.enable_spot_restore,
	i.receiving_bucket,
	i.restore_bucket,
	i.created_at,
	i.updated_at,
	i.member_institution_id as "parent_id",
	parent."name" as "parent_name",
	parent.identifier as "parent_identifier",
	parent.state as "parent_state",
	parent.deactivated_at as "parent_deactivated_at"
from institutions i
left join institutions parent on i.member_institution_id = parent.id;

-- intellectual_objects_view

create or replace view intellectual_objects_view as
select
	io.id,
	io.title,
	io.description,
	io.identifier,
	io.alt_identifier,
	io.access,
	io.bag_name,
	io.institution_id,
	io.state,
	io.etag,
	io.bag_group_identifier,
	io.storage_option,
	io.bagit_profile_identifier,
	io.source_organization,
	io.internal_sender_identifier,
	io.internal_sender_description,
	io.created_at,
	io.updated_at,
	i."name" as institution_name,
	i.identifier as institution_identifier,
	i."type" as institution_type,
	i.member_institution_id as institution_parent_id,
	(select count(*) from generic_files gf where gf.intellectual_object_id = io.id and gf.state = 'A') as "file_count",
	(select sum(gf."size") from generic_files gf where gf.intellectual_object_id = io.id and gf.state = 'A') as "size"
from intellectual_objects io
left join institutions i on io.institution_id = i.id;


-------------------------------------------------------------------------------
-- Functions
-------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION create_constraint_if_not_exists (t_name text, c_name text, constraint_sql text)
  RETURNS void
AS
$BODY$
  begin
    -- Look for our constraint
    if not exists (select constraint_name
                   from information_schema.constraint_column_usage
                   where table_name = t_name  and constraint_name = c_name) then
        execute 'ALTER TABLE ' || t_name || ' ADD CONSTRAINT ' || c_name || ' ' || constraint_sql;
    end if;
end;
$BODY$
LANGUAGE plpgsql VOLATILE;
