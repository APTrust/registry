# Database Changes

- [x] Add unique constraint on roles.name

## Standardize IDs

Some tables use serial (32-bit int) ids and some use bigserial (64-bit int). Foreign keys pointing back to primary ids are also a mix of int4 (32-bit) and int8 (64-bit).

The system should use int8/bigserial/64-bit throughout. Note that the method for making this change in Postgres 10+ may involve two steps, [as described here](https://stackoverflow.com/questions/52195303/postgresql-primary-key-id-datatype-from-serial-to-bigserial). We have to alter both the sequence and the ID column. Most other answers say the sequence already returns bigint, and we only have to alter the column from int to bigint.

## Fix Missing Constraints

* institutions.identifier should be NOT NULL
* institutions.identifier should be UNIQUE

## Review Indexes

* generic_files has way too many indexes, and most of them are probably not used
* intellectual_objects has too many indexes, most are probably unused
* work_items has quite a few indexes. Review and remove unnecessary ones.
* review indexes on intellectual_objects

## Normalize / Deduplicate Data

### Premis Events

* Remove premis_events.intellectual_object_identifier
* Remove premis_events.generic_file_identifier
* Use a view to join premis_events to objects and files tables to get identifiers

### Work Items

* Remove work_items.object_identifier
* Remove work_items.generic_file_identifier
* Use a view to pull in those two fields based on obj ID and gf ID.

## Remove Unused Columns

* generic_files.ingest_state
* intellectual_object.ingest_state

## Rename Ambiguous Columns

* work_items.date should be something like "date\_processed"

## Multi-Step Involved Changes

### Move Role into Users Table

Currently a user has one institution (users.institution\_id) and can have multiple roles through roles\_users to the roles table. This makes no sense. Because of our business rules, a user can have only one role at one institution. To fix this:

1. Add column users.role.
2. Copy each user's role from user->roles_users->role.name to users.role.
3. Drop table users_roles.
4. Drop table roles.

We may have to apply this change back in the feature/storage-record branch of Pharos, or if that's too difficult due to the brittleness of the old Rails code, apply it stages here. If in stages, we would do items 1 and 2 above, and save items 3 and 4 (dropping those tables) until after the registry goes into production.
