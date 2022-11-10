-- Initial schema, as it existed after migrating the old Pharos DB
-- to the new Registry structure.
--
-- Note that this schema still contains a number of Pharos legacy tables,
-- such as ar_internal_metadata, confirmation_tokens, and the bulk_delete
-- tables.
-- 
-- We may drop some of these in future, though we may need to keep the
-- bulk delete tables for auditing. And we will keep ar_internal_metadata
-- and schema_migrations because they're useful.



CREATE TABLE ar_internal_metadata (
	"key" varchar NOT NULL,
	value varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key)
);



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



CREATE TABLE bulk_delete_jobs_emails (
	bulk_delete_job_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_emails_on_bulk_delete_job_id ON public.bulk_delete_jobs_emails USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_emails_on_email_id ON public.bulk_delete_jobs_emails USING btree (email_id);



CREATE TABLE bulk_delete_jobs_generic_files (
	bulk_delete_job_id int8 NULL,
	generic_file_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_generic_files_on_bulk_delete_job_id ON public.bulk_delete_jobs_generic_files USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_generic_files_on_generic_file_id ON public.bulk_delete_jobs_generic_files USING btree (generic_file_id);



CREATE TABLE bulk_delete_jobs_institutions (
	bulk_delete_job_id int8 NULL,
	institution_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_institutions_on_bulk_delete_job_id ON public.bulk_delete_jobs_institutions USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_institutions_on_institution_id ON public.bulk_delete_jobs_institutions USING btree (institution_id);



CREATE TABLE bulk_delete_jobs_intellectual_objects (
	bulk_delete_job_id int8 NULL,
	intellectual_object_id int8 NULL
);
CREATE INDEX index_bulk_delete_jobs_intellectual_objects_on_bulk_job_id ON public.bulk_delete_jobs_intellectual_objects USING btree (bulk_delete_job_id);
CREATE INDEX index_bulk_delete_jobs_intellectual_objects_on_object_id ON public.bulk_delete_jobs_intellectual_objects USING btree (intellectual_object_id);



CREATE TABLE confirmation_tokens (
	id bigserial NOT NULL,
	"token" varchar NULL,
	intellectual_object_id int4 NULL,
	generic_file_id int4 NULL,
	institution_id int4 NULL,
	user_id int4 NULL,
	CONSTRAINT confirmation_tokens_pkey PRIMARY KEY (id)
);



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



CREATE TABLE emails_generic_files (
	generic_file_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_generic_files_on_email_id ON public.emails_generic_files USING btree (email_id);
CREATE INDEX index_emails_generic_files_on_generic_file_id ON public.emails_generic_files USING btree (generic_file_id);



CREATE TABLE emails_intellectual_objects (
	intellectual_object_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_intellectual_objects_on_email_id ON public.emails_intellectual_objects USING btree (email_id);
CREATE INDEX index_emails_intellectual_objects_on_intellectual_object_id ON public.emails_intellectual_objects USING btree (intellectual_object_id);



CREATE TABLE emails_premis_events (
	premis_event_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_premis_events_on_email_id ON public.emails_premis_events USING btree (email_id);
CREATE INDEX index_emails_premis_events_on_premis_event_id ON public.emails_premis_events USING btree (premis_event_id);



CREATE TABLE emails_work_items (
	work_item_id int8 NULL,
	email_id int8 NULL
);
CREATE INDEX index_emails_work_items_on_email_id ON public.emails_work_items USING btree (email_id);
CREATE INDEX index_emails_work_items_on_work_item_id ON public.emails_work_items USING btree (work_item_id);



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
CREATE INDEX ix_gf_last_fixity_check ON public.generic_files USING btree (last_fixity_check);



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
	enable_spot_restore bool NOT NULL DEFAULT false,
	receiving_bucket varchar NOT NULL,
	restore_bucket varchar NOT NULL,
	CONSTRAINT institutions_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_institutions_identifier ON public.institutions USING btree (identifier);
CREATE INDEX index_institutions_on_name ON public.institutions USING btree (name);
CREATE UNIQUE INDEX index_institutions_receiving_bucket ON public.institutions USING btree (receiving_bucket);
CREATE UNIQUE INDEX index_institutions_restore_bucket ON public.institutions USING btree (restore_bucket);



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



CREATE TABLE schema_migrations (
	"version" varchar NOT NULL,
	CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);



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



CREATE TABLE usage_samples (
	id serial4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	institution_id varchar NULL,
	"data" text NULL,
	CONSTRAINT usage_samples_pkey PRIMARY KEY (id)
);



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



CREATE TABLE storage_records (
	id bigserial NOT NULL,
	generic_file_id int4 NULL,
	url varchar NULL,
	CONSTRAINT storage_records_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_a126ea6adc FOREIGN KEY (generic_file_id) REFERENCES generic_files(id)
);
CREATE INDEX index_storage_records_on_generic_file_id ON public.storage_records USING btree (generic_file_id);



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



CREATE TABLE deletion_requests_generic_files (
	deletion_request_id int4 NOT NULL,
	generic_file_id int4 NOT NULL,
	CONSTRAINT deletion_requests_generic_files_deletion_request_id_fkey FOREIGN KEY (deletion_request_id) REFERENCES deletion_requests(id),
	CONSTRAINT deletion_requests_generic_files_generic_file_id_fkey FOREIGN KEY (generic_file_id) REFERENCES generic_files(id)
);
CREATE UNIQUE INDEX index_drgf_unique ON public.deletion_requests_generic_files USING btree (deletion_request_id, generic_file_id);



CREATE TABLE deletion_requests_intellectual_objects (
	deletion_request_id int4 NOT NULL,
	intellectual_object_id int4 NOT NULL,
	CONSTRAINT deletion_requests_intellectual_obje_intellectual_object_id_fkey FOREIGN KEY (intellectual_object_id) REFERENCES intellectual_objects(id),
	CONSTRAINT deletion_requests_intellectual_objects_deletion_request_id_fkey FOREIGN KEY (deletion_request_id) REFERENCES deletion_requests(id)
);
CREATE UNIQUE INDEX index_drio_unique ON public.deletion_requests_intellectual_objects USING btree (deletion_request_id, intellectual_object_id);



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



CREATE TABLE alerts_premis_events (
	alert_id int4 NOT NULL,
	premis_event_id int4 NOT NULL,
	CONSTRAINT alerts_premis_events_alert_id_fkey FOREIGN KEY (alert_id) REFERENCES alerts(id),
	CONSTRAINT alerts_premis_events_premis_event_id_fkey FOREIGN KEY (premis_event_id) REFERENCES premis_events(id)
);
CREATE INDEX index_alerts_premis_events_alert_id ON public.alerts_premis_events USING btree (alert_id);
CREATE UNIQUE INDEX index_alerts_premis_events_unique ON public.alerts_premis_events USING btree (alert_id, premis_event_id);



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



CREATE TABLE alerts_work_items (
	alert_id int4 NOT NULL,
	work_item_id int4 NOT NULL,
	CONSTRAINT alerts_work_items_alert_id_fkey FOREIGN KEY (alert_id) REFERENCES alerts(id),
	CONSTRAINT alerts_work_items_work_item_id_fkey FOREIGN KEY (work_item_id) REFERENCES work_items(id)
);
CREATE INDEX index_alerts_work_items_alert_id ON public.alerts_work_items USING btree (alert_id);
CREATE UNIQUE INDEX index_alerts_work_items_unique ON public.alerts_work_items USING btree (alert_id, work_item_id);



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



CREATE MATERIALIZED VIEW public.generic_file_counts
TABLESPACE pg_default
AS SELECT generic_files.institution_id,
    count(generic_files.id) AS row_count,
    generic_files.state
   FROM generic_files
  GROUP BY CUBE(generic_files.institution_id, generic_files.state)
  ORDER BY generic_files.institution_id, generic_files.state
WITH DATA;



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



CREATE OR REPLACE VIEW public.institutions_view
AS SELECT i.id,
    i.name,
    i.identifier,
    i.state,
    i.type,
    i.deactivated_at,
    i.otp_enabled,
    i.enable_spot_restore,
    i.receiving_bucket,
    i.restore_bucket,
    i.created_at,
    i.updated_at,
    i.member_institution_id AS parent_id,
    parent.name AS parent_name,
    parent.identifier AS parent_identifier,
    parent.state AS parent_state,
    parent.deactivated_at AS parent_deactivated_at
   FROM institutions i
     LEFT JOIN institutions parent ON i.member_institution_id = parent.id;



CREATE MATERIALIZED VIEW public.intellectual_object_counts
TABLESPACE pg_default
AS SELECT intellectual_objects.institution_id,
    count(intellectual_objects.id) AS row_count,
    intellectual_objects.state
   FROM intellectual_objects
  GROUP BY CUBE(intellectual_objects.institution_id, intellectual_objects.state)
  ORDER BY intellectual_objects.institution_id, intellectual_objects.state
WITH DATA;



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



CREATE MATERIALIZED VIEW public.premis_event_counts
TABLESPACE pg_default
AS SELECT premis_events.institution_id,
    count(premis_events.id) AS row_count,
    premis_events.event_type,
    premis_events.outcome
   FROM premis_events
  GROUP BY CUBE(premis_events.institution_id, premis_events.event_type, premis_events.outcome)
  ORDER BY premis_events.institution_id, premis_events.event_type, premis_events.outcome
WITH DATA;



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



CREATE MATERIALIZED VIEW public.work_item_counts
TABLESPACE pg_default
AS SELECT work_items.institution_id,
    count(work_items.id) AS row_count,
    work_items.action
   FROM work_items
  GROUP BY CUBE(work_items.institution_id, work_items.action)
  ORDER BY work_items.institution_id, work_items.action
WITH DATA;



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


CREATE OR REPLACE FUNCTION public.update_counts()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
  begin
    refresh materialized view premis_event_counts;
    refresh materialized view intellectual_object_counts;
    refresh materialized view generic_file_counts;
    refresh materialized view work_item_counts;    
    return 1;
  end;
$function$
;
