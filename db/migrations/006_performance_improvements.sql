-- 006_performance_improvements.sql
--
-- Improve performance of three problematic operations.
--
-- 1. Queries related to fixity checking do large table scans on generic_files.
-- 2. Multiple simultaneous processes of update_count() result in deadlock.
-- 3. Multiple simultaneous processes of update_current_deposit_stats() result in deadlock.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('006_performance_improvements', now())
on conflict ("version") do update set started_at = now();

-- Production DB logs show that queries to queue and process
-- fixity checks take a long time and are run frequently.
-- This index should help.
create index if not exists ix_generic_files_state_opt_fixity 
    on generic_files("state", storage_option, last_fixity_check);




-- We're updating this to return immediately if "update counts is running".
-- This should avoid deadlock caused by two processes running simultaneously.
drop function if exists public.update_counts;
CREATE OR REPLACE FUNCTION public.update_counts()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
  begin
    -- Don't start running this if it's already running. You'll get a long deadlock.
    if exists (select 1 from ar_internal_metadata where "key"='update counts is running' and "value" = 'true') then 
        return 0;
    end if;

    if exists (select 1 from work_item_counts where updated_at < (current_timestamp - interval '60 minutes')) or not exists (select * from work_item_counts where institution_id is not null limit 1) then

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
	return 0;
  end;
$function$
;


-- Production DB logs show this update is sometimes run simultaneously by 
-- two different instances of Registry. Here, as above, we add a short-circuit
-- to return immediately if these updates are currently in process.
drop function if exists public.update_current_deposit_stats;
CREATE OR REPLACE FUNCTION public.update_current_deposit_stats()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
  begin
    -- Don't start running this if it's already running. You'll get a long deadlock.
    if exists (select 1 from ar_internal_metadata where "key"='current_deposit_stats is running' and "value" = 'true') then 
        return 0;
    end if;

    if exists (select 1 from current_deposit_stats where end_date < (current_timestamp - interval '60 minutes')) or not exists (select * from current_deposit_stats where institution_id is not null limit 1) then

        insert into ar_internal_metadata ("key", "value", created_at, updated_at) 
        values ('current_deposit_stats is running', 'true', now(), now())
   		    on conflict("key") do update set "value" = 'true';    

    	refresh materialized view concurrently current_deposit_stats;

   		update ar_internal_metadata set "value" = 'false', updated_at = now() where "key" = 'current_deposit_stats is running';

	    return 1;
	end if;
	return 0;
  end;
$function$
;


-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '006_performance_improvements';
