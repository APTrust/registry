-- 003_alter_deposit_stats.sql
--
-- This migration adds column member_institution_id to the
-- materialized views for deposit stats.
--

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('003_alter_deposit_stats', now())
on conflict ("version") do update set started_at = now();


-- Add the member_institutiton_id column to historical_deposit_stats.
alter table historical_deposit_stats add column if not exists member_institution_id int4;
alter table historical_deposit_stats add column if not exists primary_sort varchar;
alter table historical_deposit_stats add column if not exists secondary_sort varchar;


-- Populate the new member_institution_id column.
update historical_deposit_stats
set member_institution_id = i.member_institution_id
from institutions i 
where institution_id = i.id 
and institution_id is not null;

-- Update the sort columns
update historical_deposit_stats set primary_sort = institution_name;
update historical_deposit_stats set primary_sort = 'zzz' where institution_name = 'All Institutions';

update historical_deposit_stats set secondary_sort = storage_option;
update historical_deposit_stats set secondary_sort = 'zzz' where storage_option = 'Total';


-- This function populates the historical_deposit_stats table with 
-- numbers up to the end of the prior month. E.g. If end_date
-- is August 1, 2022, the stats in that row show deposits through
-- July 31, 2022 at 11:59:59 pm.
--
-- This new version includes the new member_institution_id column.
create or replace function populate_historical_deposit_stats(stop_date timestamp) 
	returns int
as 
$BODY$
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
$BODY$
LANGUAGE plpgsql VOLATILE;


-- current_deposit_stats contains current deposit stats.
-- We update this hourly. These stats take a few minutes to
-- gather, and we don't want to collect them while a user
-- is waiting for a page to load.
-- 
-- We can refresh this at any time using
-- refresh materialized view current_deposit_stats.
--
-- Note that this drop statement may generate a warning saying 
-- "materialized view "current_deposit_stats" does not exist, skipping"
--
-- That's fine. It just means it skipped this step because there
-- was nothing to do.
drop materialized view if exists current_deposit_stats;
create materialized view current_deposit_stats as
select
  i2.id as institution_id,
  i2.member_institution_id as member_institution_id,
  coalesce(stats.institution_name, 'All Institutions') as institution_name,
  coalesce(stats.storage_option, 'Total') as storage_option,
  stats.file_count,
  stats.object_count,
  stats.total_bytes,
  (stats.total_bytes / 1073741824) as total_gb,
  (stats.total_bytes / 1099511627776) as total_tb,
  so.cost_gb_per_month,
  ((stats.total_bytes / 1073741824) * so.cost_gb_per_month) as monthly_cost,
  now() as end_date,
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
  group by cube (i."name", gf.storage_option)) stats
left join storage_options so on so."name" = stats.storage_option
left join institutions i2 on i2."name" = stats.institution_name
order by stats.institution_name, stats.storage_option;



-- Now populate the materialized views.
select update_current_deposit_stats();
select populate_all_historical_deposit_stats();

-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '003_alter_deposit_stats';
