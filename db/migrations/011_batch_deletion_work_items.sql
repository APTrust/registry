-- 011_batch_deletion_work_items.sql
-- 
-- Allow a single deletion request to map to multiple WorkItems.
-- This supports batch deletions.
-- 

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('011_batch_deletion_work_items', now())
on conflict ("version") do update set started_at = now();


do
$$
begin
	if exists(
		select 1 from information_schema.columns
		where table_schema = 'public'
		and table_name = 'deletion_requests'
		and column_name = 'work_item_id')
	then

        -- We need to rebuild this view, removing all references
        -- to the WorkItems table.
		drop view if exists deletion_requests_view; 

		create or replace view deletion_requests_view
		AS SELECT dr.id,
		    dr.institution_id,
		    i.name AS institution_name,
		    i.identifier AS institution_identifier,
		    dr.requested_by_id,
		    req.name AS requested_by_name,
		    req.email AS requested_by_email,
		    dr.requested_at,
		    dr.confirmed_by_id,
		    conf.name AS confirmed_by_name,
		    conf.email AS confirmed_by_email,
		    dr.confirmed_at,
		    dr.cancelled_by_id,
		    can.name AS cancelled_by_name,
		    can.email AS cancelled_by_email,
		    dr.cancelled_at,
		    ( SELECT count(*) AS count
		           FROM deletion_requests_generic_files drgf
		          WHERE drgf.deletion_request_id = dr.id) AS file_count,
		    ( SELECT count(*) AS count
		           FROM deletion_requests_intellectual_objects drio
		          WHERE drio.deletion_request_id = dr.id) AS object_count
		   FROM deletion_requests dr
		     LEFT JOIN institutions i ON dr.institution_id = i.id
		     LEFT JOIN users req ON dr.requested_by_id = req.id
		     LEFT JOIN users conf ON dr.confirmed_by_id = conf.id
		     LEFT JOIN users can ON dr.confirmed_by_id = can.id;
	
	
	
		-- Add deletion_request_id to work_items as a nullable
		-- foreign key to deletion_requests.
		alter table work_items add column deletion_request_id bigint null;
		alter table work_items add constraint fk_work_items_deletion_request_id 
            foreign key (deletion_request_id) references deletion_requests (id);	

        -- Copy deletion request ids from legacy requests into the work_items table.
        update work_items  
        set deletion_request_id = dr.id 
        from work_items wi inner join deletion_requests dr on dr.work_item_id = wi.id
        where dr.work_item_id = wi.id;
       
        -- Now remove the work_item_id column from deletion requests
		alter table deletion_requests drop column work_item_id;


		drop view if exists work_items_view; 

		CREATE OR REPLACE VIEW public.work_items_view
		AS SELECT wi.id,
			wi.institution_id,
			i.name AS institution_name,
			i.identifier AS institution_identifier,
			wi.intellectual_object_id,
			io.identifier AS object_identifier,
			io.alt_identifier,
			io.bag_group_identifier,
			io.storage_option,
			io.bagit_profile_identifier,
			io.source_organization,
			io.internal_sender_identifier,
			wi.generic_file_id,
			gf.identifier AS generic_file_identifier,
			wi.name,
			wi.etag,
			wi.bucket,
			wi."user",
			wi.note,
			wi.action,
			wi.stage,
			wi.status,
			wi.outcome,
			wi.bag_date,
			wi.date_processed,
			wi.retry,
			wi.node,
			wi.pid,
			wi.needs_admin_review,
			wi.size,
			wi.queued_at,
			wi.stage_started_at,
			wi.aptrust_approver,
			wi.inst_approver,
			wi.deletion_request_id,
			wi.created_at,
			wi.updated_at
		FROM work_items wi
			LEFT JOIN institutions i ON wi.institution_id = i.id
			LEFT JOIN intellectual_objects io ON wi.intellectual_object_id = io.id
			LEFT JOIN generic_files gf ON wi.generic_file_id = gf.id;


	end if;
end
$$;


-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '011_batch_deletion_work_items';
