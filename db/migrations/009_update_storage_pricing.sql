-- 009_update_storage_pricing.sql
-- 
-- Update the prices for different storage options
-- based on APTrust's public pricing sheet at 
-- tinyurl.com/bdz6vtrf 
--
-- Trello ticket: https://trello.com/c/LXFEcjcP
--

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('009_update_storage_pricing', now())
on conflict ("version") do update set started_at = now();

-- Update the prices
update storage_options set cost_gb_per_month = 0.41016, updated_at = now() where "name" = 'Standard';
update storage_options set cost_gb_per_month = 0.05859, updated_at = now() where "service" = 'Glacier';
update storage_options set cost_gb_per_month = 0.01953, updated_at = now() where "service" = 'Glacier-Deep';

-- Now rebuild the historical deposit stats table
delete from historical_deposit_stats ;
select populate_all_historical_deposit_stats();


-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '009_update_storage_pricing';
