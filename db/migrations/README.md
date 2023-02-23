# Database Migrations

Database migrations

* must be placed in the `db/migrations` directory
* must end in `.sql`
* must begin with a lexically sortable number, they are run in the correct
order during testing
* must be idempotent, so if they're run repeatedly, they always leave the DB 
in the same state
* must update the schema migrations tables, so we know they were run

Registry's test suite will automatically run migrations when you run 
`./registry test`, as long as the migration files are in the migrations 
directory and end with `.sql`.

See the `runMigrations()` function in db/testutil.go for the implementation. 

## Naming Migration Files

Use a name like `001_deposit_stats.sql`, which has the following components:

* `001` is a lexically sortable numeric prefix, so we know `runMigrations()` 
will run this migration before `002`, `003`, etc.
* `deposit_stats` describes what this migration is about.
* `.sql` tells `runMigrations()` to run this as a SQL file.

**Note** As of February 2023, we've already run six migrations. The prefix
for the next one should start with `007`.

## Idempotency

A migration should leave the DB in the same state, even if it's run 
repeatedly. You can accomplish this by testing conditions before executing
a change. 

For example, use `create index if not exists` rather than `create index`.
The same goes for tables and views. Use `create or replace function`
instead of `create function`. In some cases, you may need to 
`drop function if exists` before recreating the function.

For materialized views, you may need to drop them and then create them with
a new definition.

When adding or altering columns, test the current state of the DB before
running your create/alter command. For example:

```sql

-- Convert a varchar column to timestamp, but only if the
-- column is still varchar. The first time you run this,
-- it will convert the column type and values. Subsequent
-- runs will do nothing, because the column has already 
-- been converted.
--
-- This is idempotent: the command will always leave the
-- DB in the desired state, no matter what state it was in 
-- before.
do
$$
begin
	if exists(
		select 1 from information_schema.columns
		where table_schema = 'public'
		and table_name = 'premis_events'
		and column_name = 'date_time'
		and data_type = 'character varying')
	then
		alter table premis_events alter column date_time type timestamp
		using TO_TIMESTAMP(date_time, 'YYYY-MM-DDTHH24:MI:SS');
	end if;
end
$$;

-- Add a column, but only if it does not already exist.
do $$
begin
  if not exists (select 1 from information_schema.columns where table_schema='public' AND table_name='users' AND column_name='awaiting_second_factor') then
 	alter table users add column "awaiting_second_factor" boolean not null default false;
  end if;
end
$$;

```

## Updating schema_migrations

Migrations should update the `schema_migrations` table when they begin and
end, so we know that they have started and completed. See the migration 
template below.


## Migration Template


```sql

-- <number>_<description>.sql
-- 
-- This migration does whatever <description> says.
-- Briefly elaborate on <description>.

-- Note that we're starting the migration.
insert into schema_migrations ("version", started_at) values ('<number>_<description>', now())
on conflict ("version") do update set started_at = now();

-- The body of the migration goes here.
-- Remember: idempotent!

-- Now note that the migration is complete.
update schema_migrations set finished_at = now() where "version" = '<number_description>';

```

## Sample Migrations

See the `db/legacy_migrations` directory for example migrations that were 
actually run on the Registry DB.

## Migration Testing

Test migrations **heavily**, first in your local dev environment, then
on staging.

Do not run migrations on demo or production unless they've been successful
on staging.

If you have a potentially dangerous migration that alters data, create a 
copy of the entire demo or production DB in RDS and run the migration on the
copy first.

For example, the sample above that changes the `premis_events.date_time`
column from varchar to timestamp had to be run individually on a **copy**
of all our databases to ensure that none of them contained an 
un-convertable value in the `date_time` column. 

A single bad value in this column would cause the entire type conversion to 
fail. We really don't want that happening in production, so we do a test
run first.

(Don't get me started on why someone would create a varchar column to 
store timestamps and then name it date_time.)

## Running Migrations on Staging, Demo and Production

Copy the migration file to the bastion server of your target environment
like this:

```
scp db/migrations/005_deposit_timeline_stats.sql bastion-staging:
```

Log in to the bastion server of the target environment. You'll find a script
in your home directory called `db_connect.sh`. You'll also see your sql file.

If you expect the migration to take more than a few seconds to run, start
a screen session so that the migration will continue even if your ssh session
gets cut off. You can start a named `screen` session with this command:

```
screen -S migration
```

Now you can run the migration with this command:

```
./db_connect.sh < 005_deposit_timeline_stats.sql > migration_005_output.txt 2> migration_005_error.txt 
```

After you run the command, you'll be prompted to enter the database password,
which `db_connect.sh` is kind enough to tell you. If you're running a long
migration in screen, you can disconnect with `Ctrl-A Ctrl-D` to get out of 
screen. You can disconnect from the ssh session as well and check back later.

Before you disconnect, check the output of `screen -ls`. It should say there
is a session called `migration` running **and that the session is detached.**
Don't log out or disconnect unless the session is detached.

When you return, the results of the migration will be available in the file
`migration_005_output.txt`, with error messages in  `migration_005_error.txt`.

You can re-attach to the screen session using `screen -R migration`. If it has
completed, you can kill the session from within using `Ctrl-D`.

You could likely also run migrations on our target environments using a local
GUI running through an SSH tunnel. In that case, however, you lose the 
protection of screen, and network interruptions may cause your migration to
fail.

If the `db_connect.sh` script ever gets wiped out, keep in mind that you can
get all the DB parameters, including endpoint URL, username and password
from AWS Parameter Store.

All that's really in the `db_connect.sh` script is:

```sh
#!/bin/bash

echo "Connecting to [database name]"
echo "Password is [actual password]"
psql -h <url> -U <db_name>
```

So, it connects you to the DB and tells you how to answer the password prompt.

## Why are we doing it this way?

Why not use a DB migration library? 

Because I haven't yet found one that's better than this process.

Sure, this is a pain in the ass, and requires some babysitting, but:

* Migrations are rare. This isn't something we do weekly or even monthy.
* We want them to be attended, because a botched migration is a big 
problem, and we want to be on it ASAP.
* AFAIK, no migration libraries offer idempotent DDL. You have to write
that by hand.
* We have only one target backend, Postgres, so multi-backend ORM 
libraries don't buy us anything. They would just be generating the same
DDL we're writing by hand.

If a good migration library does come along that satisfies all our needs,
we'll use it.

