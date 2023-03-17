-- 008_fix_update_counts.sql
-- 
-- This migration fixes a problem with update_counts that
-- sometimes caused two calls to the function to overlap,
-- leading to deadlock. This is a long-running call that
-- runs in the background every hour or so.


-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('008_fix_update_counts', now())
on conflict ("version") do update set started_at = now();

drop function if exists lock_update_counts();
drop function if exists unlock_update_counts();

create unique index if not exists ix_generic_file_counts on generic_file_counts (institution_id, "state");
create unique index if not exists ix_intellectual_object_counts on intellectual_object_counts (institution_id, "state");
create unique index if not exists ix_premis_event_counts on premis_event_counts (institution_id, event_type, outcome);
create unique index if not exists ix_work_item_counts on work_item_counts (institution_id, "action");


CREATE OR REPLACE FUNCTION public.update_counts()
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
  begin
    -- Don't start running this if it's already running. You'll get a long deadlock.
    if exists (select 1 from ar_internal_metadata where "key"='update counts is running' and "value" = 'true') then 
    	raise notice 'update_counts is running in another process (has value true)';
        return 0;
    end if;

    -- Another hint that this function is already running is
    -- a lock on "update counts" row in the metadata table.
    -- That update isn't committed until the entire function 
    -- completes, which had been causing deadlocks. 
    -- This is the key addition in migration 008_fix_update_counts.
	if exists (SELECT id FROM ar_internal_metadata aim where "key" = 'update counts is running') and not exists (SELECT id FROM ar_internal_metadata aim where "key" = 'update counts is running' FOR UPDATE SKIP locked) then 
    	raise notice 'update_counts is running in another process (metadata row is locked)';
		return 0;
	end if;
      
    if exists (select 1 from work_item_counts where updated_at < (current_timestamp - interval '60 minutes')) or not exists (select * from work_item_counts where institution_id is not null limit 1) then

		-- Use ar_internal_metadata to track whether this function is running.
		-- These are some long-running queries, especially for premis events.
		-- we want to avoid the case where this function gets kicked off while
		-- a previous iteration is still in progress.
   		insert into ar_internal_metadata ("key", "value", created_at, updated_at) values ('update counts is running', 'true', now(), now())
   		on conflict("key") do update set "value" = 'true';
   	
    	refresh materialized view concurrently premis_event_counts;
   		refresh materialized view concurrently intellectual_object_counts;
   		refresh materialized view concurrently generic_file_counts;
   		refresh materialized view concurrently work_item_counts;    

   		update ar_internal_metadata set "value" = 'false', updated_at = now() where "key" = 'update counts is running';
   		return 1;
	end if;
	raise notice 'update_counts ran recently and does not need to be re-run now';
	return 0;
  end;
$function$
;

-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '008_fix_update_counts';