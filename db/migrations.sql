-- migrations.sql
--
-- This file contains ALL alterations that should be applied to the
-- existing Pharos DB schema to make it match schema.sql.
--
-- All operations in this file must be idempotent, so we can run it
-- any number of times and always know that it will leave the DB in a
-- consistent and known state that matches schema.sql.
--
-------------------------------------------------------------------------------

-- First off, make sure role names are unique.
-- I have no idea why this constraint did not exist in Pharos.
-- Note that the name of the constraint is created by postgres
-- as 'roles_names_key' when we define roles.name as varchar unique.
select create_constraint_if_not_exists('roles', 'roles_name_key', 'unique("name");');


-- We need to fix the user role structure. Pharos allows a user to have
-- multiple roles at a single institution, though our business rules disallow
-- that, and no user has ever had more than one role. To simplify the DB
-- and our queries, we need to do the following:
--
-- 1. Create a role column in the users table.
-- 2. Populate that column with the value with each user's role from
--    user -> roles_user -> roles.
-- 3. Drop the roles_users table.
-- 4. Drop the roles table.

do $$
begin
  if not exists (select 1 from information_schema.columns where table_schema='public' AND table_name='users' AND column_name='role') then
 	alter table users add column "role" varchar(50) not null default 'none';
 	update users u set "role" = coalesce((select r.name from "roles" r inner join roles_users ru on ru.role_id = r.id where ru.user_id = u.id), 'none');
    drop table if exists roles_users;
    drop table if exists roles;
  end if;
end
$$
