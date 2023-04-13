-- 010_alerts_view_email.sql
-- 
-- Add user email address to alerts_view. 
--

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('010_alerts_view_email', now())
on conflict ("version") do update set started_at = now();


drop view public.alerts_view;

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
    u.email as user_email,
    au.sent_at,
    au.read_at
   FROM alerts a
     LEFT JOIN alerts_users au ON a.id = au.alert_id
     LEFT JOIN users u ON au.user_id = u.id
     LEFT JOIN institutions i ON a.institution_id = i.id;

-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '010_alerts_view_email';
