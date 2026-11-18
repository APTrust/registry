-- 013_content_transfer_requests.sql
--
-- This migration adds tables and constraints that pertain to
-- content transfer requests. Users can now request the transfer
-- of intellectual objects to be sent to different storage tiers
-- (i.e. storage classes and/or geolocations), and the database
-- will store who requests the transfers, who confirms or cancels
-- them, and the timeline that these actions took place.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('013_content_transfer_requests', now())
on conflict ("version") do update set started_at = now();

-- public.transfer_requests definition

CREATE TABLE public.transfer_requests (
	id bigserial NOT NULL,
	institution_id int4 NOT NULL,
	requested_by_id int4 NOT NULL,
	requested_at timestamp NOT NULL,
	encrypted_confirmation_token varchar NOT NULL,
	confirmed_by_id int4 NULL,
	confirmed_at timestamp NULL,
	cancelled_by_id int4 NULL,
	cancelled_at timestamp NULL,
	CONSTRAINT transfer_requests_pkey PRIMARY KEY (id)
);
CREATE INDEX index_transfer_requests_institution_id ON public.transfer_requests USING btree (institution_id);

-- public.transfer_requests_intellectual_objects definition

CREATE TABLE public.transfer_requests_intellectual_objects (
     transfer_request_id int4 NOT NULL,
	intellectual_object_id int4 NOT NULL
);
CREATE UNIQUE INDEX index_trio_unique ON public.transfer_requests_intellectual_objects USING btree (transfer_request_id, intellectual_object_id);

-- public.transfer_requests_view source

CREATE OR REPLACE VIEW public.transfer_requests_view
AS SELECT tr.id,
    tr.institution_id,
    i.name AS institution_name,
    i.identifier AS institution_identifier,
    tr.requested_by_id,
    req.name AS requested_by_name,
    req.email AS requested_by_email,
    tr.requested_at,
    tr.confirmed_by_id,
    conf.name AS confirmed_by_name,
    conf.email AS confirmed_by_email,
    tr.confirmed_at,
    tr.cancelled_by_id,
    can.name AS cancelled_by_name,
    can.email AS cancelled_by_email,
    tr.cancelled_at,
    ( SELECT count(*) AS count
           FROM transfer_requests_intellectual_objects trio
          WHERE trio.transfer_request_id = tr.id) AS object_count
   FROM transfer_requests tr
     LEFT JOIN institutions i ON tr.institution_id = i.id
     LEFT JOIN users req ON tr.requested_by_id = req.id
     LEFT JOIN users conf ON tr.confirmed_by_id = conf.id
     LEFT JOIN users can ON tr.confirmed_by_id = can.id;

-- Add constraints, work items, alerts

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '013_content_transfer_requests';
