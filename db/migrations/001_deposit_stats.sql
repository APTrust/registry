-- migration_001.sql
--
-- This migration improves support for reporting and deposit stats
-- by creating cached reports and a materialized view that can be 
-- updated asynchronously. It also makes changes to the existing 
-- "slow count" materialized views, to avoid updating them unnecessarily.
--
-- It also adds support for tracking migrations and tracking whether
-- some long-running queries/functions are currently running.


-- Let's start tracking our schema migrations.
alter table schema_migrations add column if not exists started_at timestamp not null;
alter table schema_migrations add column if not exists finished_at timestamp null;

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('001_deposit_stats', now())
on conflict ("version") do update set started_at = now();


-- historical_deposit_stats contains a snapshot of deposits for each
-- month, going back to 2014, when APTrust first launched. The end_date
-- column contains a date on the first of a given month. E.g. If end_date
-- is August 1, 2022, the stats in that row show deposits through
-- July 31, 2022 at 11:59:59 pm. That is, the starts for for the end of
-- the prior month.
--
-- We insert into this table once a month using the function
-- populate_historical_deposit_stats().
create table historical_deposit_stats (
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
create index ix_historical_deposit_stats_inst_id on historical_deposit_stats(institution_id);
create index ix_historical_deposit_stats_storage_option on historical_deposit_stats(storage_option);
create index ix_historical_deposit_stats_end_date on historical_deposit_stats(end_date);

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
			  coalesce(stats.institution_name, 'All Institutions') as institution_name,
			  coalesce(stats.storage_option, 'Total') as storage_option,
			  coalesce(stats.file_count, 0) as file_count,
			  coalesce(stats.object_count, 0) as object_count,
			  coalesce(stats.total_bytes, 0) as total_bytes,
			  coalesce((stats.total_bytes / 1073741824), 0) as total_gb,
			  coalesce((stats.total_bytes / 1099511627776), 0) as total_tb,
			  coalesce(so.cost_gb_per_month, 0) as cost_gb_per_month,
			  coalesce(((stats.total_bytes / 1073741824) * so.cost_gb_per_month), 0) as monthly_cost,
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


-- populate_all_historical_deposit_stats fills in historical deposit stats
-- for every month between January, 2014 and last month. We call this as a 
-- cron job from the Registry. Note that the underlying function, 
-- populate_historical_deposit_stats does no work if the stats for the
-- requested date already exist in the historical_deposit_stats table.
-- The stats queries are expensive, so we want to avoid running them when
-- they're not necessary.
--
-- Since Registry runs this as a cron job in the background, it will not
-- slow requests from users.
create or replace function populate_all_historical_deposit_stats() 
	returns void
as 
$BODY$
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
$BODY$
LANGUAGE plpgsql VOLATILE;


-- current_deposit_stats contains current deposit stats.
-- We update this hourly. These stats take a few minutes to
-- gather, and we don't want to collect them while a user
-- is waiting for a page to load.
-- 
-- We can refresh this at any time using
-- refresh materialized view current_deposit_stats
drop materialized view if exists current_deposit_stats;
create materialized view current_deposit_stats as
select
  i2.id as institution_id,
  coalesce(stats.institution_name, 'All Institutions') as institution_name,
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


-- We also want to add a timestamp to the slow counts views
-- to include a timestamp, so we know when they were last updated.
-- This helps prevent excessive refreshing of these views.
-- The refresh is slow and expensive in terms of read operations.
--
-- We need to drop these views first, then re-create them.

-- premis_event_counts
drop materialized view if exists premis_event_counts;
create materialized view premis_event_counts as
	select institution_id, count(id) as row_count, event_type, outcome, current_timestamp as updated_at
	from premis_events group by cube(institution_id, event_type, outcome)
	order by institution_id, event_type, outcome;

-- intellectual_object_counts
drop materialized view if exists intellectual_object_counts;
create materialized view intellectual_object_counts as
	select institution_id, count(id) as row_count, "state", current_timestamp as updated_at 
	from intellectual_objects group by cube(institution_id, "state")
	order by institution_id, state; 

-- generic_file_counts
drop materialized view if exists generic_file_counts;
create materialized view generic_file_counts as
	select institution_id, count(id) as row_count, "state", current_timestamp as updated_at
	from generic_files group by cube(institution_id, "state")
	order by institution_id, state; 

-- work_item_counts
drop materialized view if exists work_item_counts;
create materialized view work_item_counts as
	select institution_id, count(id) as row_count, "action", current_timestamp as updated_at
	from work_items group by cube(institution_id, "action")
	order by institution_id, "action";


-- Update counts at most once per hour.
-- The second part of the if condition tests to see if the view
-- is essentially empty, which happens after setup for local 
-- development and testing. This doesn't happen in staging, demo,
-- or production.
CREATE OR REPLACE FUNCTION update_counts()
  RETURNS integer
AS
$BODY$
  begin
    if exists (select 1 from work_item_counts where updated_at < (current_timestamp - interval '60 minutes')) or not exists (select * from work_item_counts where institution_id is not null limit 1) then
    	if not exists (select 1 from ar_internal_metadata where "key"='update counts is running' and "value" = 'true') then 

			-- Use ar_internal_metadata to track whether this function is running.
			-- These are some long-running queries, especially for premis events.
			-- we want to avoid the case where this function gets kicked off while
			-- a previous iteration is still in progress.
    		insert into ar_internal_metadata ("key", "value", created_at, updated_at) values ('update counts is running', 'true', now(), now())
    		on conflict("key") do update set "value" = 'true';
    	
	    	refresh materialized view premis_event_counts;
    		refresh materialized view intellectual_object_counts;
    		refresh materialized view generic_file_counts;
    		refresh materialized view work_item_counts;    

    		update ar_internal_metadata set "value" = 'false', updated_at = now() where "key" = 'update counts is running';

    		return 1;
	   	end if;
	end if;
	return 0;
  end;
$BODY$
LANGUAGE plpgsql VOLATILE;


-- Update deposit stats at most once per hour.
-- See note for update_counts() for info on the second if condition.
CREATE OR REPLACE FUNCTION update_current_deposit_stats()
  RETURNS integer
AS
$BODY$
  begin
    if exists (select 1 from current_deposit_stats where end_date < (current_timestamp - interval '60 minutes')) or not exists (select * from current_deposit_stats where institution_id is not null limit 1) then
    	refresh materialized view current_deposit_stats;
	    return 1;
	end if;
	return 0;
  end;
$BODY$
LANGUAGE plpgsql VOLATILE;


-- Now populate the materialized views.
-- In staging, demo and production systems, this will populate
-- the views with actual data. This will take several minutes,
-- and possible up to an hour to run on the production system.
-- 
-- In dev and test systems, the views will remain empty because 
-- there's no data until we load the fixtures, so we'll have to 
-- call these functions again from within our tests to properly 
-- populate the views.
select update_counts();
select update_current_deposit_stats();
select populate_all_historical_deposit_stats();

-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '001_deposit_stats';
