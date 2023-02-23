# Legacy Schema and Migrations

This folder contains legacy SQL files, including the initial
launch schema for Registry, and some older migrations.

We moved them here because, during unit, integration and end-to-end tests,
the the `runMigrations()` function in db/testutil.go sets up the intial
DB schema and then runs all the migrations on it before running tests.

It may do this multiple times in a single test run, and that gets to be
time-consuming.

In addition, a number of the migrations define, re-define, and alter the
same set of functions and views. This gets confusing. Developers should
know they can ignore anything in this legacy directory, as it's not 
current. 

The current schema is the sum of `db/schema.sql` plus the migrations in 
`db/migrations`.

Items in this directory should be considered archival and exist purely for 
the historical record.

