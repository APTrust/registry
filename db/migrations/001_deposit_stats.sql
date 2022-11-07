-- migration_001.sql
--
-- This migration improves support for reporting and deposit stats
-- by created cached reports and a materialized view that can be 
-- updated offline. It also makes changes to the existing 
-- "slow count" materialized views, to avoid updating them unnecessarily.

-- historical_deposit_stats contains a snapshot of deposits for each
-- month, going back to 2014, when APTrust first launched. The end_date
-- column contains a date on the first of a given month. E.g. If end_date
-- is August 1, 2022, the stats in that row show deposits through
-- July 31, 2022 at 11:59:59 pm. That is, the starts for for the end of
-- the prior month.
--
-- We insert into this table once a month using the function
-- populate_historical_deposit_stats().
create table if not exists historical_deposit_stats (
	institution_id     bigint,
	institution_name   varchar(80),
	storage_option     varchar(40),
	object_count       bigint,
	file_count         bigint,
	total_bytes        bigint,
	total_gb           double precision,
	total_tb           double precision,
	cost_gb_per_month  double precision,
	monthly_cost       double precision,
	end_date           timestamp
);		


-- This function populates the historical_deposit_stats table with 
-- numbers up to the end of the prior month. E.g. If end_date
-- is August 1, 2022, the stats in that row show deposits through
-- July 31, 2022 at 11:59:59 pm.
create or replace function populate_historical_deposit_stats(stop_date timestamp) 
	returns int
as 
$BODY$
	begin
		if not exists (select 1 from historical_deposit_stats where end_date = stop_date) then 
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
-- refresh materialized view current_deposit_stats
create materialized view if not exists current_deposit_stats as
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
  now() as end_date
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
