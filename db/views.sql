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

-- deletion_requests_view

create or replace view deletion_requests_view as
select dr.id,
       dr.institution_id,
       i."name" as institution_name,
       i.identifier as institution_identifier,
       dr.requested_by_id,
       req."name" as requested_by_name,
       req.email as requested_by_email,
       dr.requested_at,
       dr.confirmed_by_id,
       conf."name" as confirmed_by_name,
       conf.email as confirmed_by_email,
       dr.confirmed_at,
       dr.cancelled_by_id,
       can."name" as cancelled_by_name,
       can.email as cancelled_by_email,
       dr.cancelled_at,
       (select count(*) from deletion_requests_generic_files drgf where drgf.deletion_request_id = dr.id) as file_count,
       (select count(*) from deletion_requests_intellectual_objects drio where drio.deletion_request_id = dr.id) as object_count,
       dr.work_item_id,
       wi.stage,
       wi.status,
       wi.date_processed,
       wi."size",
       wi.note
from deletion_requests dr
left join institutions i on dr.institution_id = i.id
left join users req on dr.requested_by_id = req.id
left join users conf on dr.confirmed_by_id = conf.id
left join users can on dr.confirmed_by_id = can.id
left join work_items wi on dr.work_item_id = wi.id;

-- alerts view

create or replace view alerts_view as
select
	a.id,
	a.institution_id,
	i."name" as institution_name,
	i.identifier as institution_identifier,
	a."type",
	a.subject,
	a."content",
	a.deletion_request_id,
	a.created_at,
	au.user_id,
	u."name" as user_name,
	au.sent_at,
	au.read_at
from alerts a
left join alerts_users au on a.id = au.alert_id
left join users u on au.user_id = u.id
left join institutions i on a.institution_id = i.id;

-- storage_option_stats
create or replace view storage_option_stats as
with stats as (
	select
		sum("gf"."size") as total_bytes,
		count(*) as file_count,
		gf.institution_id,
		gf.storage_option
	from generic_files gf
	where gf.state = 'A'
	group by rollup (gf.institution_id, gf.storage_option)
	)
select
	s.total_bytes,
	s.file_count,
	s.institution_id,
	s.storage_option,
	i."name" as institution_name,
	i.identifier as institution_identifier
from stats s
left join institutions i on s.institution_id = i.id;

create or replace view generic_files_view as
select
	gf.id,
	gf.file_format,
	gf."size",
	gf.identifier,
	gf.intellectual_object_id,
	io.identifier as object_identifier,
	io."access",
	gf.state,
	gf.last_fixity_check,
	gf.institution_id,
	i."name" as institution_name,
	i.identifier as institution_identifier,
	gf.storage_option,
	gf.uuid,
	(select digest from checksums where generic_file_id = gf.id and algorithm='md5' order by created_at desc limit 1) as "md5",
	(select digest from checksums where generic_file_id = gf.id and algorithm='sha1' order by created_at desc limit 1) as "sha1",
	(select digest from checksums where generic_file_id = gf.id and algorithm='sha256' order by created_at desc limit 1) as "sha256",
	(select digest from checksums where generic_file_id = gf.id and algorithm='sha512' order by created_at desc limit 1) as "sha512",
	gf.created_at,
	gf.updated_at
from generic_files gf
left join intellectual_objects io on io.id = gf.intellectual_object_id
left join institutions i on i.id = gf.institution_id;

-- Checksums view
create or replace view checksums_view as
select cs.id,
       cs.algorithm,
       cs.datetime,
       cs.digest,
       gf.state,
       gf.identifier as "generic_file_identifier",
       cs.generic_file_id,
       gf.intellectual_object_id,
       gf.institution_id,
       cs.created_at,
       cs.updated_at
from checksums cs
left join generic_files gf on cs.generic_file_id = gf.id;
