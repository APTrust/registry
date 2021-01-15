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
