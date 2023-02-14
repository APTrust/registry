-- 004_concurrent_deposit_stats.sql
--
-- This migration makes only one change:
-- we refresh the current_deposit_stats materialized view CONCURRENTLY
-- so that the DB allows selects even as the view is refreshing.
--
-- Without this, the Registry's dashboard page often times out
-- while waiting for this view to refresh.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('004_concurrent_deposit_stats', now())
on conflict ("version") do update set started_at = now();


-- In order to use concurrent refresh, we need to create a 
-- unique index on the materialized view, so we're adding that here.
create unique index ix_current_deposits_inst_id_storage_option on current_deposit_stats(institution_id, storage_option);

-- Update deposit stats at most once per hour.
-- See note for update_counts() for info on the second if condition.
CREATE OR REPLACE FUNCTION update_current_deposit_stats()
  RETURNS integer
AS
$BODY$
  begin
    if exists (select 1 from current_deposit_stats where end_date < (current_timestamp - interval '60 minutes')) or not exists (select * from current_deposit_stats where institution_id is not null limit 1) then
    	refresh materialized view concurrently current_deposit_stats;
	    return 1;
	end if;
	return 0;
  end;
$BODY$
LANGUAGE plpgsql VOLATILE;

-- Refresh the materialized view to make sure this works.
select update_current_deposit_stats();

-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '004_concurrent_deposit_stats';
