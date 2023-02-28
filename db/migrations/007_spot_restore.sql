-- 007_spot_restore.sql
-- 
-- This migration does changes institutions.enable_spot_restore
-- to institutions.spot_restore_frequency and adds 
-- institutions.last_spot_restore_work_item_id

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('007_spot_restore', now())
on conflict ("version") do update set started_at = now();


do
$$
begin
	if exists(
		select 1 from information_schema.columns
		where table_schema = 'public'
		and table_name = 'institutions'
		and column_name = 'enable_spot_restore')
	then

        -- We need to drop this, because it refers to the institutions.enable_spot_restore
        -- column that we're about to remove.
        drop view if exists institutions_view;

        -- Add the new column with default zero
        alter table institutions add column spot_restore_frequency int not null default 0;

        -- If any institutions had spot restore enabled, 
        -- put them on a 90-day schedule. Those without spot restore
        -- enabled will stay at zero, which is the same as never/disabled.
        update institutions set spot_restore_frequency = 90 where enable_spot_restore = true;

        -- Add a column to track the WorkItem ID of the last spot restore.
        alter table institutions add column last_spot_restore_work_item_id bigint null;
        alter table institutions add constraint fk_institutions_last_spot_restore 
            foreign key (last_spot_restore_work_item_id) references work_items (id);

        -- Now get rid of the obsolete column.
        alter table institutions drop column enable_spot_restore;

        -- Now recreate the view, referring to the new spot restore columns.
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

	end if;
end
$$;


-- Now bring the old ActiveRecord table ar_internal_metadata into line with
-- our other tables by giving it a numeric serial id column. We have to 
-- get rid of the old primary key to do that, but we still want entries in
-- the "key" column to be unique, so we add a unique index at the end.
alter table ar_internal_metadata drop constraint if exists ar_internal_metadata_pkey;
alter table ar_internal_metadata add column if not exists id serial primary key;
create unique index if not exists ix_ar_internal_metadata_uniq_key on ar_internal_metadata("key");

-- Create records we'll need to track spot restorations
insert into ar_internal_metadata ("key", "value", created_at, updated_at) 
values ('spot restore is running', 'false', now(), now())
on conflict do nothing;

insert into ar_internal_metadata ("key", "value", created_at, updated_at) 
values ('spot restore last run', '2000-01-01', now(), now())
on conflict do nothing;


-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '007_spot_restore';