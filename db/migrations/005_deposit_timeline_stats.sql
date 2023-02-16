-- 005_deposit_timeline_stats.sql
--
-- Add data to historical deposit stats to explicitly record
-- when depositors had zero data in the repo.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('005_deposit_timeline_stats', now())
on conflict ("version") do update set started_at = now();


-- Change end_date from timestamp to date to simplify reports.
do
$$
begin
	if exists(
		select 1 from information_schema.columns
		where table_schema = 'public'
		and table_name = 'historical_deposit_stats'
		and column_name = 'end_date'
		and data_type ilike 'timestamp%')
	then
		alter table historical_deposit_stats alter column end_date type date;
	end if;
end
$$;

-- Create unique index on what amounts to the stats table's natural key.
create unique index if not exists ix_historical_inst_opt_date on 
historical_deposit_stats (institution_id, storage_option, end_date);

-- For timeline reports, we to know about months where a depositor had
-- nothing in storage. These are entries that explicitly show zero bytes 
-- in storage for a particular month. populate_historical_deposit stats
-- doesn't add these entries because its aggregate functions return nothing.
drop function if exists populate_empty_deposit_stats();
create or replace function populate_empty_deposit_stats() 
  returns int
as
$$
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
$$ LANGUAGE plpgsql;


-- Now we change this function to take a date parameter instead of a 
-- timestamp, and we tell it to call populate_empty_deposit_stats before 
-- it returns. This function runs once a month.
create or replace function populate_historical_deposit_stats(stop_date date) 
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

			select populate_empty_deposit_stats();
		
			return 1;
		else
			return 0;
		end if;
	end;
$BODY$
LANGUAGE plpgsql VOLATILE;


-- An early run of the function that populated historical deposit stats
-- left some incorrect object counts. Let's wipe out this table and
-- rebuild it.
delete from historical_deposit_stats ;

-- Rebuild the deposit stats table. This will take a long time on prod.
select populate_all_historical_deposit_stats();

-- Update the current stats, if necessary.
select update_current_deposit_stats();

-- Now, the point of this whole migration: Add explicit zero counts to
-- the historical_deposit_stats table for months when depositors had no
-- materials in APTrust. Our front-end charts need the data points to 
-- correctly plot the time series data.
select populate_empty_deposit_stats();

-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '005_deposit_timeline_stats';
