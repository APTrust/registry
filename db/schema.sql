-- Registry Schema - 2023-03-17

-- public.ar_internal_metadata definition

-- Drop table

-- DROP TABLE ar_internal_metadata;

CREATE TABLE ar_internal_metadata (
	"key" varchar NOT NULL,
	value varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	id serial4 NOT NULL,
	CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX ix_ar_internal_metadata_uniq_key ON public.ar_internal_metadata USING btree (key);


-- public.bulk_delete_jobs definition

-- Drop table

-- DROP TABLE bulk_delete_jobs;

CREATE TABLE bulk_delete_jobs (
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

-- Drop table

-- DROP TABLE bulk_delete_jobs_emails;

CREATE TABLE bulk_delete_jobs_emails (
	bulk_delete_job_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_emails_on_bulk_delete_job_id ON public.bulk_delete_jobs_emails USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_emails_on_email_id ON public.bulk_delete_jobs_emails USING btree (email_id);


-- public.bulk_delete_jobs_generic_files definition

-- Drop table

-- DROP TABLE bulk_delete_jobs_generic_files;

CREATE TABLE bulk_delete_jobs_generic_files (
	bulk_delete_job_id int8 NULL,
	generic_file_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_generic_files_on_bulk_delete_job_id ON public.bulk_delete_jobs_generic_files USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_generic_files_on_generic_file_id ON public.bulk_delete_jobs_generic_files USING btree (generic_file_id);


-- public.bulk_delete_jobs_institutions definition

-- Drop table

-- DROP TABLE bulk_delete_jobs_institutions;

CREATE TABLE bulk_delete_jobs_institutions (
	bulk_delete_job_id int8 NULL,
	institution_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_institutions_on_bulk_delete_job_id ON public.bulk_delete_jobs_institutions USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_institutions_on_institution_id ON public.bulk_delete_jobs_institutions USING btree (institution_id);


-- public.bulk_delete_jobs_intellectual_objects definition

-- Drop table

-- DROP TABLE bulk_delete_jobs_intellectual_objects;

CREATE TABLE bulk_delete_jobs_intellectual_objects (
	bulk_delete_job_id int8 NULL,
	intellectual_object_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_intellectual_objects_on_bulk_job_id ON public.bulk_delete_jobs_intellectual_objects USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_intellectual_objects_on_object_id ON public.bulk_delete_jobs_intellectual_objects USING btree (intellectual_object_id);


-- public.confirmation_tokens definition

-- Drop table

-- DROP TABLE confirmation_tokens;

CREATE TABLE confirmation_tokens (
	id bigserial NOT NULL,
	"token" varchar NULL,
	intellectual_object_id int4 NULL,
	generic_file_id int4 NULL,
	institution_id int4 NULL,
	user_id int4 NULL,
	CONSTRAINT confirmation_tokens_pkey PRIMARY KEY (id)
);


-- public.emails definition

-- Drop table

-- DROP TABLE emails;

CREATE TABLE emails (
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

-- Drop table

-- DROP TABLE emails_generic_files;

CREATE TABLE emails_generic_files (
	generic_file_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_generic_files_on_email_id ON public.emails_generic_files USING btree (email_id);
CREATE INDEX index_emails_generic_files_on_generic_file_id ON public.emails_generic_files USING btree (generic_file_id);


-- public.emails_intellectual_objects definition

-- Drop table

-- DROP TABLE emails_intellectual_objects;

CREATE TABLE emails_intellectual_objects (
	intellectual_object_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_intellectual_objects_on_email_id ON public.emails_intellectual_objects USING btree (email_id);
CREATE INDEX index_emails_intellectual_objects_on_intellectual_object_id ON public.emails_intellectual_objects USING btree (intellectual_object_id);


-- public.emails_premis_events definition

-- Drop table

-- DROP TABLE emails_premis_events;

CREATE TABLE emails_premis_events (
	premis_event_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_premis_events_on_email_id ON public.emails_premis_events USING btree (email_id);
CREATE INDEX index_emails_premis_events_on_premis_event_id ON public.emails_premis_events USING btree (premis_event_id);


-- public.emails_work_items definition

-- Drop table

-- DROP TABLE emails_work_items;

CREATE TABLE emails_work_items (
	work_item_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_work_items_on_email_id ON public.emails_work_items USING btree (email_id);
CREATE INDEX index_emails_work_items_on_work_item_id ON public.emails_work_items USING btree (work_item_id);


-- public.generic_files definition

-- Drop table

-- DROP TABLE generic_files;

CREATE TABLE generic_files (
	id serial4 NOT NULL,
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
CREATE INDEX ix_generic_files_state_opt_fixity ON public.generic_files USING btree (state, storage_option, last_fixity_check);
CREATE INDEX ix_gf_last_fixity_check ON public.generic_files USING btree (last_fixity_check);


-- public.historical_deposit_stats definition

-- Drop table

-- DROP TABLE historical_deposit_stats;

CREATE TABLE historical_deposit_stats (
	institution_id int8 NULL,
	institution_name varchar(80) NULL,
	storage_option varchar(40) NULL,
	object_count int8 NULL,
	file_count int8 NULL,
	total_bytes int8 NULL,
	total_gb float8 NULL,
	total_tb float8 NULL,
	cost_gb_per_month float8 NULL,
	monthly_cost float8 NULL,
	end_date date NULL,
	member_institution_id int4 NULL,
	primary_sort varchar NULL,
	secondary_sort varchar NULL
);
CREATE INDEX ix_historical_deposit_stats_end_date ON public.historical_deposit_stats USING btree (end_date);
CREATE INDEX ix_historical_deposit_stats_inst_id ON public.historical_deposit_stats USING btree (institution_id);
CREATE INDEX ix_historical_deposit_stats_storage_option ON public.historical_deposit_stats USING btree (storage_option);
CREATE UNIQUE INDEX ix_historical_inst_opt_date ON public.historical_deposit_stats USING btree (institution_id, storage_option, end_date);


-- public.intellectual_objects definition

-- Drop table

-- DROP TABLE intellectual_objects;

CREATE TABLE intellectual_objects (
	id serial4 NOT NULL,
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
	bag_group_identifier varchar NULL,
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

-- Drop table

-- DROP TABLE old_passwords;

CREATE TABLE old_passwords (
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

-- Drop table

-- DROP TABLE premis_events;

CREATE TABLE premis_events (
	id serial4 NOT NULL,
	identifier varchar NULL,
	event_type varchar NULL,
	date_time timestamp NULL,
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

-- Drop table

-- DROP TABLE schema_migrations;

CREATE TABLE schema_migrations (
	"version" varchar NOT NULL,
	started_at timestamp NULL,
	finished_at timestamp NULL,
	CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);


-- public.snapshots definition

-- Drop table

-- DROP TABLE snapshots;

CREATE TABLE snapshots (
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


-- public.storage_options definition

-- Drop table

-- DROP TABLE storage_options;

CREATE TABLE storage_options (
	id bigserial NOT NULL,
	provider varchar NOT NULL,
	service varchar NOT NULL,
	region varchar NOT NULL,
	"name" varchar NOT NULL,
	cost_gb_per_month numeric(12, 8) NOT NULL,
	"comment" varchar NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT storage_options_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_storage_options_name ON public.storage_options USING btree (name);


-- public.usage_samples definition

-- Drop table

-- DROP TABLE usage_samples;

CREATE TABLE usage_samples (
	id serial4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	institution_id varchar NULL,
	"data" text NULL,
	CONSTRAINT usage_samples_pkey PRIMARY KEY (id)
);


-- public.work_items definition

-- Drop table

-- DROP TABLE work_items;

CREATE TABLE work_items (
	id serial4 NOT NULL,
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
CREATE INDEX index_work_items_etag_instid_and_name ON public.work_items USING btree (etag, institution_id, name);
CREATE INDEX index_work_items_on_action ON public.work_items USING btree (action);
CREATE INDEX index_work_items_on_date_processed ON public.work_items USING btree (date_processed);
CREATE INDEX index_work_items_on_etag_and_name ON public.work_items USING btree (etag, name);
CREATE INDEX index_work_items_on_generic_file_id ON public.work_items USING btree (generic_file_id);
CREATE INDEX index_work_items_on_inst_id_and_date_processed ON public.work_items USING btree (institution_id, date_processed);
CREATE INDEX index_work_items_on_institution_id ON public.work_items USING btree (institution_id);
CREATE INDEX index_work_items_on_intellectual_object_id ON public.work_items USING btree (intellectual_object_id);
CREATE INDEX index_work_items_on_stage ON public.work_items USING btree (stage);
CREATE INDEX index_work_items_on_status ON public.work_items USING btree (status);


-- public.checksums definition

-- Drop table

-- DROP TABLE checksums;

CREATE TABLE checksums (
	id serial4 NOT NULL,
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


-- public.institutions definition

-- Drop table

-- DROP TABLE institutions;

CREATE TABLE institutions (
	id serial4 NOT NULL,
	"name" varchar NULL,
	identifier varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	state varchar NULL,
	"type" varchar NULL,
	member_institution_id int4 NULL,
	deactivated_at timestamp NULL,
	otp_enabled bool NULL,
	receiving_bucket varchar NOT NULL,
	restore_bucket varchar NOT NULL,
	spot_restore_frequency int4 NOT NULL DEFAULT 0,
	last_spot_restore_work_item_id int8 NULL,
	CONSTRAINT institutions_pkey PRIMARY KEY (id),
	CONSTRAINT fk_institutions_last_spot_restore FOREIGN KEY (last_spot_restore_work_item_id) REFERENCES work_items(id)
);
CREATE UNIQUE INDEX index_institutions_identifier ON public.institutions USING btree (identifier);
CREATE INDEX index_institutions_on_name ON public.institutions USING btree (name);
CREATE UNIQUE INDEX index_institutions_receiving_bucket ON public.institutions USING btree (receiving_bucket);
CREATE UNIQUE INDEX index_institutions_restore_bucket ON public.institutions USING btree (restore_bucket);


-- public.storage_records definition

-- Drop table

-- DROP TABLE storage_records;

CREATE TABLE storage_records (
	id bigserial NOT NULL,
	generic_file_id int4 NULL,
	url varchar NULL,
	CONSTRAINT storage_records_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_a126ea6adc FOREIGN KEY (generic_file_id) REFERENCES generic_files(id)
);
CREATE INDEX index_storage_records_on_generic_file_id ON public.storage_records USING btree (generic_file_id);


-- public.users definition

-- Drop table

-- DROP TABLE users;

CREATE TABLE users (
	id serial4 NOT NULL,
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
	encrypted_otp_sent_at timestamp NULL,
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
	awaiting_second_factor bool NOT NULL DEFAULT false,
	"role" varchar(50) NOT NULL DEFAULT 'none'::character varying,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_7fcf39ca13 FOREIGN KEY (institution_id) REFERENCES institutions(id)
);
CREATE INDEX index_users_on_authy_id ON public.users USING btree (authy_id);
CREATE UNIQUE INDEX index_users_on_email ON public.users USING btree (email);
CREATE INDEX index_users_on_institution_id ON public.users USING btree (institution_id);
CREATE INDEX index_users_on_password_changed_at ON public.users USING btree (password_changed_at);
CREATE UNIQUE INDEX index_users_on_reset_password_token ON public.users USING btree (reset_password_token);


-- public.deletion_requests definition

-- Drop table

-- DROP TABLE deletion_requests;

CREATE TABLE deletion_requests (
	id bigserial NOT NULL,
	institution_id int4 NOT NULL,
	requested_by_id int4 NOT NULL,
	requested_at timestamp NOT NULL,
	encrypted_confirmation_token varchar NOT NULL,
	confirmed_by_id int4 NULL,
	confirmed_at timestamp NULL,
	cancelled_by_id int4 NULL,
	cancelled_at timestamp NULL,
	work_item_id int4 NULL,
	CONSTRAINT deletion_requests_pkey PRIMARY KEY (id),
	CONSTRAINT deletion_requests_cancelled_by_id_fkey FOREIGN KEY (cancelled_by_id) REFERENCES users(id),
	CONSTRAINT deletion_requests_confirmed_by_id_fkey FOREIGN KEY (confirmed_by_id) REFERENCES users(id),
	CONSTRAINT deletion_requests_institution_id_fkey FOREIGN KEY (institution_id) REFERENCES institutions(id),
	CONSTRAINT deletion_requests_requested_by_id_fkey FOREIGN KEY (requested_by_id) REFERENCES users(id),
	CONSTRAINT deletion_requests_work_item_id_fkey FOREIGN KEY (work_item_id) REFERENCES work_items(id)
);
CREATE INDEX index_deletion_requests_institution_id ON public.deletion_requests USING btree (institution_id);


-- public.deletion_requests_generic_files definition

-- Drop table

-- DROP TABLE deletion_requests_generic_files;

CREATE TABLE deletion_requests_generic_files (
	deletion_request_id int4 NOT NULL,
	generic_file_id int4 NOT NULL,
	CONSTRAINT deletion_requests_generic_files_deletion_request_id_fkey FOREIGN KEY (deletion_request_id) REFERENCES deletion_requests(id),
	CONSTRAINT deletion_requests_generic_files_generic_file_id_fkey FOREIGN KEY (generic_file_id) REFERENCES generic_files(id)
);
CREATE UNIQUE INDEX index_drgf_unique ON public.deletion_requests_generic_files USING btree (deletion_request_id, generic_file_id);


-- public.deletion_requests_intellectual_objects definition

-- Drop table

-- DROP TABLE deletion_requests_intellectual_objects;

CREATE TABLE deletion_requests_intellectual_objects (
	deletion_request_id int4 NOT NULL,
	intellectual_object_id int4 NOT NULL,
	CONSTRAINT deletion_requests_intellectual_obje_intellectual_object_id_fkey FOREIGN KEY (intellectual_object_id) REFERENCES intellectual_objects(id),
	CONSTRAINT deletion_requests_intellectual_objects_deletion_request_id_fkey FOREIGN KEY (deletion_request_id) REFERENCES deletion_requests(id)
);
CREATE UNIQUE INDEX index_drio_unique ON public.deletion_requests_intellectual_objects USING btree (deletion_request_id, intellectual_object_id);


-- public.alerts definition

-- Drop table

-- DROP TABLE alerts;

CREATE TABLE alerts (
	id bigserial NOT NULL,
	institution_id int4 NULL,
	"type" varchar NOT NULL,
	subject varchar NOT NULL,
	"content" text NOT NULL,
	deletion_request_id int4 NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT alerts_pkey PRIMARY KEY (id),
	CONSTRAINT alerts_deletion_request_id_fkey FOREIGN KEY (deletion_request_id) REFERENCES deletion_requests(id),
	CONSTRAINT alerts_institution_id_fkey FOREIGN KEY (institution_id) REFERENCES institutions(id)
);
CREATE INDEX index_alerts_institution_id ON public.alerts USING btree (institution_id);
CREATE INDEX index_alerts_type ON public.alerts USING btree (type);


-- public.alerts_premis_events definition

-- Drop table

-- DROP TABLE alerts_premis_events;

CREATE TABLE alerts_premis_events (
	alert_id int4 NOT NULL,
	premis_event_id int4 NOT NULL,
	CONSTRAINT alerts_premis_events_alert_id_fkey FOREIGN KEY (alert_id) REFERENCES alerts(id),
	CONSTRAINT alerts_premis_events_premis_event_id_fkey FOREIGN KEY (premis_event_id) REFERENCES premis_events(id)
);
CREATE INDEX index_alerts_premis_events_alert_id ON public.alerts_premis_events USING btree (alert_id);
CREATE UNIQUE INDEX index_alerts_premis_events_unique ON public.alerts_premis_events USING btree (alert_id, premis_event_id);


-- public.alerts_users definition

-- Drop table

-- DROP TABLE alerts_users;

CREATE TABLE alerts_users (
	alert_id int4 NOT NULL,
	user_id int4 NOT NULL,
	sent_at timestamp NULL,
	read_at timestamp NULL,
	CONSTRAINT alerts_users_alert_id_fkey FOREIGN KEY (alert_id) REFERENCES alerts(id),
	CONSTRAINT alerts_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX index_alerts_users_alert_id ON public.alerts_users USING btree (alert_id);
CREATE UNIQUE INDEX index_alerts_users_unique ON public.alerts_users USING btree (alert_id, user_id);
CREATE INDEX index_alerts_users_user_id ON public.alerts_users USING btree (user_id);


-- public.alerts_work_items definition

-- Drop table

-- DROP TABLE alerts_work_items;

CREATE TABLE alerts_work_items (
	alert_id int4 NOT NULL,
	work_item_id int4 NOT NULL,
	CONSTRAINT alerts_work_items_alert_id_fkey FOREIGN KEY (alert_id) REFERENCES alerts(id),
	CONSTRAINT alerts_work_items_work_item_id_fkey FOREIGN KEY (work_item_id) REFERENCES work_items(id)
);
CREATE INDEX index_alerts_work_items_alert_id ON public.alerts_work_items USING btree (alert_id);
CREATE UNIQUE INDEX index_alerts_work_items_unique ON public.alerts_work_items USING btree (alert_id, work_item_id);


-- public.alerts_view source

CREATE OR REPLACE VIEW public.alerts_view
AS SELECT a.id,
    a.institution_id,
    i.name AS institution_name,
    i.identifier AS institution_identifier,
    a.type,
    a.subject,
    a.content,
    a.deletion_request_id,
    a.created_at,
    au.user_id,
    u.name AS user_name,
    au.sent_at,
    au.read_at
   FROM alerts a
     LEFT JOIN alerts_users au ON a.id = au.alert_id
     LEFT JOIN users u ON au.user_id = u.id
     LEFT JOIN institutions i ON a.institution_id = i.id;


-- public.checksums_view source

CREATE OR REPLACE VIEW public.checksums_view
AS SELECT cs.id,
    cs.algorithm,
    cs.datetime,
    cs.digest,
    gf.state,
    gf.identifier AS generic_file_identifier,
    cs.generic_file_id,
    gf.intellectual_object_id,
    gf.institution_id,
    cs.created_at,
    cs.updated_at
   FROM checksums cs
     LEFT JOIN generic_files gf ON cs.generic_file_id = gf.id;


-- public.current_deposit_stats source

CREATE MATERIALIZED VIEW public.current_deposit_stats
TABLESPACE pg_default
AS SELECT i2.id AS institution_id,
    i2.member_institution_id,
    COALESCE(stats.institution_name, 'All Institutions'::character varying) AS institution_name,
    COALESCE(stats.storage_option, 'Total'::character varying) AS storage_option,
    stats.file_count,
    stats.object_count,
    stats.total_bytes,
    stats.total_bytes / 1073741824::numeric AS total_gb,
    stats.total_bytes / '1099511627776'::bigint::numeric AS total_tb,
    so.cost_gb_per_month,
    stats.total_bytes / 1073741824::numeric * so.cost_gb_per_month AS monthly_cost,
    now() AS end_date,
    COALESCE(stats.institution_name, 'zzz'::character varying) AS primary_sort,
    COALESCE(stats.storage_option, 'zzz'::character varying) AS secondary_sort
   FROM ( SELECT i.name AS institution_name,
            count(gf.id) AS file_count,
            count(DISTINCT gf.intellectual_object_id) AS object_count,
            sum(gf.size) AS total_bytes,
            gf.storage_option
           FROM generic_files gf
             LEFT JOIN institutions i ON i.id = gf.institution_id
          WHERE gf.state::text = 'A'::text
          GROUP BY CUBE(i.name, gf.storage_option)) stats
     LEFT JOIN storage_options so ON so.name::text = stats.storage_option::text
     LEFT JOIN institutions i2 ON i2.name::text = stats.institution_name::text
  ORDER BY stats.institution_name, stats.storage_option
WITH DATA;

-- View indexes:
CREATE UNIQUE INDEX ix_current_deposits_inst_id_storage_option ON public.current_deposit_stats USING btree (institution_id, storage_option);


-- public.deletion_requests_view source

CREATE OR REPLACE VIEW public.deletion_requests_view
AS SELECT dr.id,
    dr.institution_id,
    i.name AS institution_name,
    i.identifier AS institution_identifier,
    dr.requested_by_id,
    req.name AS requested_by_name,
    req.email AS requested_by_email,
    dr.requested_at,
    dr.confirmed_by_id,
    conf.name AS confirmed_by_name,
    conf.email AS confirmed_by_email,
    dr.confirmed_at,
    dr.cancelled_by_id,
    can.name AS cancelled_by_name,
    can.email AS cancelled_by_email,
    dr.cancelled_at,
    ( SELECT count(*) AS count
           FROM deletion_requests_generic_files drgf
          WHERE drgf.deletion_request_id = dr.id) AS file_count,
    ( SELECT count(*) AS count
           FROM deletion_requests_intellectual_objects drio
          WHERE drio.deletion_request_id = dr.id) AS object_count,
    dr.work_item_id,
    wi.stage,
    wi.status,
    wi.date_processed,
    wi.size,
    wi.note
   FROM deletion_requests dr
     LEFT JOIN institutions i ON dr.institution_id = i.id
     LEFT JOIN users req ON dr.requested_by_id = req.id
     LEFT JOIN users conf ON dr.confirmed_by_id = conf.id
     LEFT JOIN users can ON dr.confirmed_by_id = can.id
     LEFT JOIN work_items wi ON dr.work_item_id = wi.id;


-- public.generic_file_counts source

CREATE MATERIALIZED VIEW public.generic_file_counts
TABLESPACE pg_default
AS SELECT generic_files.institution_id,
    count(generic_files.id) AS row_count,
    generic_files.state,
    CURRENT_TIMESTAMP AS updated_at
   FROM generic_files
  GROUP BY CUBE(generic_files.institution_id, generic_files.state)
  ORDER BY generic_files.institution_id, generic_files.state
WITH DATA;

-- View indexes:
CREATE UNIQUE INDEX ix_generic_file_counts ON public.generic_file_counts USING btree (institution_id, state);


-- public.generic_files_view source

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


-- public.institutions_view source

CREATE OR REPLACE VIEW public.institutions_view
AS SELECT i.id,
    i.name,
    i.identifier,
    i.state,
    i.type,
    i.deactivated_at,
    i.otp_enabled,
    i.receiving_bucket,
    i.restore_bucket,
    i.spot_restore_frequency,
    i.last_spot_restore_work_item_id,
    i.created_at,
    i.updated_at,
    i.member_institution_id AS parent_id,
    parent.name AS parent_name,
    parent.identifier AS parent_identifier,
    parent.state AS parent_state,
    parent.deactivated_at AS parent_deactivated_at
   FROM institutions i
     LEFT JOIN institutions parent ON i.member_institution_id = parent.id;


-- public.intellectual_object_counts source

CREATE MATERIALIZED VIEW public.intellectual_object_counts
TABLESPACE pg_default
AS SELECT intellectual_objects.institution_id,
    count(intellectual_objects.id) AS row_count,
    intellectual_objects.state,
    CURRENT_TIMESTAMP AS updated_at
   FROM intellectual_objects
  GROUP BY CUBE(intellectual_objects.institution_id, intellectual_objects.state)
  ORDER BY intellectual_objects.institution_id, intellectual_objects.state
WITH DATA;

-- View indexes:
CREATE UNIQUE INDEX ix_intellectual_object_counts ON public.intellectual_object_counts USING btree (institution_id, state);


-- public.intellectual_objects_view source

CREATE OR REPLACE VIEW public.intellectual_objects_view
AS SELECT io.id,
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
    i.name AS institution_name,
    i.identifier AS institution_identifier,
    i.type AS institution_type,
    i.member_institution_id AS institution_parent_id,
    ( SELECT count(*) AS count
           FROM generic_files gf
          WHERE gf.intellectual_object_id = io.id AND gf.state::text = 'A'::text) AS file_count,
    ( SELECT sum(gf.size) AS sum
           FROM generic_files gf
          WHERE gf.intellectual_object_id = io.id AND gf.state::text = 'A'::text) AS size,
    ( SELECT count(*) AS count
           FROM generic_files gf
          WHERE gf.intellectual_object_id = io.id AND gf.state::text = 'A'::text AND gf.identifier::text ~~ concat(io.identifier, '/data/%')) AS payload_file_count,
    ( SELECT sum(gf.size) AS sum
           FROM generic_files gf
          WHERE gf.intellectual_object_id = io.id AND gf.state::text = 'A'::text AND gf.identifier::text ~~ concat(io.identifier, '/data/%')) AS payload_size
   FROM intellectual_objects io
     LEFT JOIN institutions i ON io.institution_id = i.id;


-- public.premis_event_counts source

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

-- View indexes:
CREATE UNIQUE INDEX ix_premis_event_counts ON public.premis_event_counts USING btree (institution_id, event_type, outcome);


-- public.premis_events_view source

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


-- public.storage_option_stats source

CREATE OR REPLACE VIEW public.storage_option_stats
AS WITH stats AS (
         SELECT sum(gf.size) AS total_bytes,
            count(*) AS file_count,
            gf.institution_id,
            gf.storage_option
           FROM generic_files gf
          WHERE gf.state::text = 'A'::text
          GROUP BY ROLLUP(gf.institution_id, gf.storage_option)
        )
 SELECT s.total_bytes,
    s.file_count,
    s.institution_id,
    s.storage_option,
    i.name AS institution_name,
    i.identifier AS institution_identifier
   FROM stats s
     LEFT JOIN institutions i ON s.institution_id = i.id;


-- public.users_view source

CREATE OR REPLACE VIEW public.users_view
AS SELECT u.id,
    u.name,
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
    u.role,
    i.name AS institution_name,
    i.identifier AS institution_identifier,
    i.state AS institution_state,
    i.type AS institution_type,
    i.member_institution_id,
    i2.name AS member_institution_name,
    i2.identifier AS member_institution_identifier,
    i2.state AS member_institution_state,
    i.otp_enabled,
    i.receiving_bucket,
    i.restore_bucket
   FROM users u
     LEFT JOIN institutions i ON u.institution_id = i.id
     LEFT JOIN institutions i2 ON i.member_institution_id = i2.id;


-- public.work_item_counts source

CREATE MATERIALIZED VIEW public.work_item_counts
TABLESPACE pg_default
AS SELECT work_items.institution_id,
    count(work_items.id) AS row_count,
    work_items.action,
    CURRENT_TIMESTAMP AS updated_at
   FROM work_items
  GROUP BY CUBE(work_items.institution_id, work_items.action)
  ORDER BY work_items.institution_id, work_items.action
WITH DATA;

-- View indexes:
CREATE UNIQUE INDEX ix_work_item_counts ON public.work_item_counts USING btree (institution_id, action);


-- public.work_items_view source

CREATE OR REPLACE VIEW public.work_items_view
AS SELECT wi.id,
    wi.institution_id,
    i.name AS institution_name,
    i.identifier AS institution_identifier,
    wi.intellectual_object_id,
    io.identifier AS object_identifier,
    io.alt_identifier,
    io.bag_group_identifier,
    io.storage_option,
    io.bagit_profile_identifier,
    io.source_organization,
    io.internal_sender_identifier,
    wi.generic_file_id,
    gf.identifier AS generic_file_identifier,
    wi.name,
    wi.etag,
    wi.bucket,
    wi."user",
    wi.note,
    wi.action,
    wi.stage,
    wi.status,
    wi.outcome,
    wi.bag_date,
    wi.date_processed,
    wi.retry,
    wi.node,
    wi.pid,
    wi.needs_admin_review,
    wi.size,
    wi.queued_at,
    wi.stage_started_at,
    wi.aptrust_approver,
    wi.inst_approver,
    wi.created_at,
    wi.updated_at
   FROM work_items wi
     LEFT JOIN institutions i ON wi.institution_id = i.id
     LEFT JOIN intellectual_objects io ON wi.intellectual_object_id = io.id
     LEFT JOIN generic_files gf ON wi.generic_file_id = gf.id;



CREATE OR REPLACE FUNCTION public.create_constraint_if_not_exists(t_name text, c_name text, constraint_sql text)
 RETURNS void
 LANGUAGE plpgsql
AS $function$
  begin
    -- Look for our constraint
    if not exists (select constraint_name
                   from information_schema.constraint_column_usage
                   where table_name = t_name  and constraint_name = c_name) then
        execute 'ALTER TABLE ' || t_name || ' ADD CONSTRAINT ' || c_name || ' ' || constraint_sql;
    end if;
end;
$function$
;

CREATE OR REPLACE FUNCTION public.populate_all_historical_deposit_stats()
 RETURNS void
 LANGUAGE plpgsql
AS $function$
DECLARE
   current_year    INTEGER := date_part('year', now());
   current_month   INTEGER := date_part('month', now());
   start_year      INTEGER := 2014;
   start_month     INTEGER := 1;
   already_populating VARCHAR;
BEGIN 
	select "value" into already_populating from ar_internal_metadata where "key" = 'historical deposit stats is running';
	raise notice '%', already_populating;
	if (already_populating is null or already_populating != 'true') then	
		-- Set a flag in ar_internal_metadata so know this process is running.
		-- We do this because multiple Registry containers may call this function 
		-- while it's already running (on the first of the month). This is a long-running
		-- select/insert query, and we don't want to overtax the DB, nor do we want
		-- to end up with duplicate rows in the historical_deposit_stats table.
		insert into ar_internal_metadata ("key", "value", created_at, updated_at) 
		values ('historical deposit stats is running', 'true', now(), now())
		on conflict ("key") do 
		update set "value" = 'true', updated_at = now();
		
		for year in start_year..current_year loop
   			for month in 1..12 loop
	   			if make_timestamp(year, month,1,0,0,0) < now() then
	   				perform populate_historical_deposit_stats(make_timestamp(year, month,1,0,0,0));
	    		end if;
   			end loop;
   		end loop;
   	
   		-- Now clear the metadata flag.
   		update ar_internal_metadata set "value" = 'false' where key = 'historical deposit stats is running';
   	end if;
end; 
$function$
;

CREATE OR REPLACE FUNCTION public.populate_current_deposit_stats()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
	declare report_date timestamp;
	begin
		report_date := now();
		insert into historical_deposit_stats
		select
		  i2.id as institution_id,
		  coalesce(stats.institution_name, 'Total') as institution_name,
		  coalesce(stats.storage_option, 'Total') as storage_option,
		  stats.file_count,
		  stats.object_count,
		  stats.total_bytes,
		  (stats.total_bytes / 1073741824) as total_gb,
		  (stats.total_bytes / 1099511627776) as total_tb,
		  so.cost_gb_per_month,
		  ((stats.total_bytes / 1073741824) * so.cost_gb_per_month) as monthly_cost,
		  report_date as end_date
		from
		  (select
			i."name" as institution_name,
			count(gf.id) as file_count,
			count(distinct(gf.intellectual_object_id)) as object_count,
			sum(gf.size) as total_bytes,
			gf.storage_option
		  from generic_files gf
		  left join institutions i on i.id = gf.institution_id
		  where gf.state = 'A'
		  group by cube (i."name", gf.storage_option)) stats
		left join storage_options so on so."name" = stats.storage_option
		left join institutions i2 on i2."name" = stats.institution_name;		
	
		return 1;
	end;
$function$
;

CREATE OR REPLACE FUNCTION public.populate_deposit_stats(stop_date date)
 RETURNS void
 LANGUAGE plpgsql
AS $function$
	begin
		if not exists (select 1 from deposit_stats where end_date = stop_date) then 
			insert into deposit_stats
			select
			  i2.id as institution_id,
			  coalesce(stats.institution_name, 'Total') as institution_name,
			  coalesce(stats.storage_option, 'Total') as storage_option,
			  stats.file_count,
			  stats.object_count,
			  stats.total_bytes,
			  (stats.total_bytes / 1073741824) as total_gb,
			  (stats.total_bytes / 1099511627776) as total_tb,
			  so.cost_gb_per_month,
			  ((stats.total_bytes / 1073741824) * so.cost_gb_per_month) as monthly_cost,
			  stop_date as end_date
			from
			  (select
				i."name" as institution_name,
				count(gf.id) as file_count,
				count(distinct(gf.intellectual_object_id)) as object_count,
				sum(gf.size) as total_bytes,
				gf.storage_option
			  from generic_files gf
			  left join institutions i on i.id = gf.institution_id
			  where gf.state = 'A'
			  and gf.created_at < stop_date
			  group by cube (i."name", gf.storage_option)) stats
			left join storage_options so on so."name" = stats.storage_option
			left join institutions i2 on i2."name" = stats.institution_name;
		end if;
	end;
$function$
;

CREATE OR REPLACE FUNCTION public.populate_empty_deposit_stats()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
declare
	inst_id int8;
	end_dt date;
	storage_opt varchar;
begin 
	for inst_id in select distinct(institution_id) from historical_deposit_stats 
	loop
		for end_dt in select distinct(end_date) from historical_deposit_stats
		loop 
			for storage_opt in select distinct(storage_option) from historical_deposit_stats
			loop 

				if inst_id is null and not exists (select * from historical_deposit_stats where institution_id is null and storage_option = storage_opt and end_date = end_dt) then 
					insert into historical_deposit_stats (institution_id, institution_name, storage_option, object_count, 
						file_count, total_bytes, total_gb, total_tb, cost_gb_per_month, monthly_cost, end_date, 
						member_institution_id, primary_sort, secondary_sort)
					values (null, 'All Institutions', storage_opt,0,0,0,0,0,0,0, end_dt, 0, 'zzz', storage_opt);
				end if;

				if not exists (select * from historical_deposit_stats where institution_id = inst_id and storage_option = storage_opt and end_date = end_dt) then 
					insert into historical_deposit_stats (institution_id, institution_name, storage_option, object_count, 
						file_count, total_bytes, total_gb, total_tb, cost_gb_per_month, monthly_cost, end_date, 
						member_institution_id, primary_sort, secondary_sort)
					select i.id, i.name, storage_opt, 0,0,0,0,0,0,0,end_dt, i.member_institution_id, i.name, storage_opt from institutions i where i.id = inst_id;
				end if;
			end loop;
		end loop;
	end loop;
	update historical_deposit_stats set secondary_sort = 'zzz' where secondary_sort = 'Total';
    return 1;
end;
$function$
;

CREATE OR REPLACE FUNCTION public.populate_historical_deposit_stats(stop_date date)
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
	begin
		if not exists (select 1 from historical_deposit_stats where end_date = stop_date) then 
			insert into historical_deposit_stats (
			  institution_id,
              member_institution_id,
			  institution_name,
			  storage_option,
			  file_count,
			  object_count,
			  total_bytes,
			  total_gb,
			  total_tb,
			  cost_gb_per_month,
			  monthly_cost,
			  end_date, 
              primary_sort,
              secondary_sort
            )
			select
			  i2.id as institution_id,
              i2.member_institution_id as member_institution_id,
			  coalesce(stats.institution_name, 'All Institutions') as institution_name,
			  coalesce(stats.storage_option, 'Total') as storage_option,
			  coalesce(stats.file_count, 0) as file_count,
			  coalesce(stats.object_count, 0) as object_count,
			  coalesce(stats.total_bytes, 0) as total_bytes,
			  coalesce((stats.total_bytes / 1073741824), 0) as total_gb,
			  coalesce((stats.total_bytes / 1099511627776), 0) as total_tb,
			  coalesce(so.cost_gb_per_month, 0) as cost_gb_per_month,
			  coalesce(((stats.total_bytes / 1073741824) * so.cost_gb_per_month), 0) as monthly_cost,
			  stop_date as end_date,
			  coalesce(stats.institution_name, 'zzz') as primary_sort,
			  coalesce(stats.storage_option, 'zzz') as secondary_sort
			from
			  (select
				i."name" as institution_name,
				count(gf.id) as file_count,
				count(distinct(gf.intellectual_object_id)) as object_count,
				sum(gf.size) as total_bytes,
				gf.storage_option
			  from generic_files gf
			  left join institutions i on i.id = gf.institution_id
			  where gf.state = 'A'
			  and gf.created_at < stop_date
			  group by cube (i."name", gf.storage_option)) stats
			left join storage_options so on so."name" = stats.storage_option
			left join institutions i2 on i2."name" = stats.institution_name;

			select populate_empty_deposit_stats();
		
			return 1;
		else
			return 0;
		end if;
	end;
$function$
;

CREATE OR REPLACE FUNCTION public.populate_historical_deposit_stats(stop_date timestamp without time zone)
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
	begin
		if not exists (select 1 from historical_deposit_stats where end_date = stop_date) then 
			insert into historical_deposit_stats (
			  institution_id,
              member_institution_id,
			  institution_name,
			  storage_option,
			  file_count,
			  object_count,
			  total_bytes,
			  total_gb,
			  total_tb,
			  cost_gb_per_month,
			  monthly_cost,
			  end_date, 
              primary_sort,
              secondary_sort
            )
			select
			  i2.id as institution_id,
              i2.member_institution_id as member_institution_id,
			  coalesce(stats.institution_name, 'All Institutions') as institution_name,
			  coalesce(stats.storage_option, 'Total') as storage_option,
			  coalesce(stats.file_count, 0) as file_count,
			  coalesce(stats.object_count, 0) as object_count,
			  coalesce(stats.total_bytes, 0) as total_bytes,
			  coalesce((stats.total_bytes / 1073741824), 0) as total_gb,
			  coalesce((stats.total_bytes / 1099511627776), 0) as total_tb,
			  coalesce(so.cost_gb_per_month, 0) as cost_gb_per_month,
			  coalesce(((stats.total_bytes / 1073741824) * so.cost_gb_per_month), 0) as monthly_cost,
			  stop_date as end_date,
			  coalesce(stats.institution_name, 'zzz') as primary_sort,
			  coalesce(stats.storage_option, 'zzz') as secondary_sort
			from
			  (select
				i."name" as institution_name,
				count(gf.id) as file_count,
				count(distinct(gf.intellectual_object_id)) as object_count,
				sum(gf.size) as total_bytes,
				gf.storage_option
			  from generic_files gf
			  left join institutions i on i.id = gf.institution_id
			  where gf.state = 'A'
			  and gf.created_at < stop_date
			  group by cube (i."name", gf.storage_option)) stats
			left join storage_options so on so."name" = stats.storage_option
			left join institutions i2 on i2."name" = stats.institution_name;
		
			return 1;
		else
			return 0;
		end if;
	end;
$function$
;

CREATE OR REPLACE FUNCTION public.update_counts()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
  begin
    -- Don't start running this if it's already running. You'll get a long deadlock.
    if exists (select 1 from ar_internal_metadata where "key"='update counts is running' and "value" = 'true') then 
    	raise notice 'update_counts is running in another process (has value true)';
        return 0;
    end if;

    -- Another hint that this function is already running is
    -- a lock on "update counts" row in the metadata table.
    -- That update isn't committed until the entire function 
    -- completes, which had been causing deadlocks. 
    -- This is the key addition in migration 008_fix_update_counts.
	if exists (SELECT id FROM ar_internal_metadata aim where "key" = 'update counts is running') and not exists (SELECT id FROM ar_internal_metadata aim where "key" = 'update counts is running' FOR UPDATE SKIP locked) then 
    	raise notice 'update_counts is running in another process (metadata row is locked)';
		return 0;
	end if;
      
    if exists (select 1 from work_item_counts where updated_at < (current_timestamp - interval '60 minutes')) or not exists (select * from work_item_counts where institution_id is not null limit 1) then

		-- Use ar_internal_metadata to track whether this function is running.
		-- These are some long-running queries, especially for premis events.
		-- we want to avoid the case where this function gets kicked off while
		-- a previous iteration is still in progress.
   		insert into ar_internal_metadata ("key", "value", created_at, updated_at) values ('update counts is running', 'true', now(), now())
   		on conflict("key") do update set "value" = 'true';
   	
    	refresh materialized view concurrently premis_event_counts;
   		refresh materialized view concurrently intellectual_object_counts;
   		refresh materialized view concurrently generic_file_counts;
   		refresh materialized view concurrently work_item_counts;    

   		update ar_internal_metadata set "value" = 'false', updated_at = now() where "key" = 'update counts is running';
   		return 1;
	end if;
	raise notice 'update_counts ran recently and does not need to be re-run now';
	return 0;
  end;
$function$
;

CREATE OR REPLACE FUNCTION public.update_current_deposit_stats()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
  begin
    -- Don't start running this if it's already running. You'll get a long deadlock.
    if exists (select 1 from ar_internal_metadata where "key"='current_deposit_stats is running' and "value" = 'true') then 
        return 0;
    end if;

    if exists (select 1 from current_deposit_stats where end_date < (current_timestamp - interval '60 minutes')) or not exists (select * from current_deposit_stats where institution_id is not null limit 1) then

        insert into ar_internal_metadata ("key", "value", created_at, updated_at) 
        values ('current_deposit_stats is running', 'true', now(), now())
   		    on conflict("key") do update set "value" = 'true';    

    	refresh materialized view concurrently current_deposit_stats;

   		update ar_internal_metadata set "value" = 'false', updated_at = now() where "key" = 'current_deposit_stats is running';

	    return 1;
	end if;
	return 0;
  end;
$function$
;
