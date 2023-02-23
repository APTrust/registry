-- This pertains to Trello issue https://trello.com/c/mlyJpN34
--
-- Newer restoration work items have action "Restore File" or
-- "Restore Object". Legacy work items from Pharos have 
-- action "Restore". Those legacy items won't show up when 
-- users filter for "Restore File" or "Restore Object". In fact,
-- they're invisible in the web UI.
--
-- This migration sets the legacy "Restore" action to either 
-- "Restore File" or "Restore Object", depending on whether the
-- WorkItem includes a generic file id. If it does, it's a file
-- restoration. If not, it's an object restoration. 
--
-- As of 2022-12-09, we have 13715 object restorations and 
-- 22 file restorations.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('002_work_item_restorations', now())
on conflict ("version") do update set started_at = now();


update work_items set action='Restore Object' where action='Restore' and generic_file_id is null;
update work_items set action='Restore File' where action='Restore' and generic_file_id is not null;


-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '002_work_item_restorations';
