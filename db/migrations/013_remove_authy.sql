-- 013_posix_metadata.sql
--
-- This migration removes Authy-related columns and indices from the database.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('013_remove_authy', now())
on conflict ("version") do update set started_at = now();

alter table public.users	drop column authy_id;
alter table public.users drop column last_sign_in_with_authy;
alter table public.users drop column authy_status varchar NULL;
alter table public.users add column mfa_status varchar NULL;

drop index index_users_on_authy_id on public.users;

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
    u.mfa_status,
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

-- Now mark the migration as completed.
update schema_migrations set finished_at = now() where "version" = '013_remove_authy';
