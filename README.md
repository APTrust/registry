# Registry

APTrust's registry contains metadata about the millions of items we hold on behalf of our depositors in preservation storage. The registry provides web-based and API-based access to the metadata and to information about items currently undergoing ingest and restoration.

This will be the third-generation of our registry software, based on the [Gin Web Framework](https://github.com/gin-gonic/gin). It will replace our existing Rails application, [Pharos](https://github.com/APTrust/pharos), which suffers from performance and maintability problems, which are discussed in more detail in [these notes](notes.md).

# Requirements

To run the registry on your local dev machine, you will need the following for ALL operations:

* [Go](https://golang.org/dl/) 1.16 or higher
* [Postgres](https://www.postgresql.org/download/) 11 or higher

You will also need the following for some Admin operations:

* [Redis](https://redis.io/download) 5.07 or higher
* [NSQ](https://nsq.io/deployment/installing.html) version 1.20 or higher

If you're running on Linux or OSX, you'll find the required Redis and NSQ binaries in this repo's .bin/linux and ./bin/osx directories. The run.sh script automatically starts them when it runs the server and tests.

# Database Setup

Run the following commands in postgres to create the user and databases:

```sql
create user dev_user with password 'password';

create database apt_registry_development owner dev_user;
create database apt_registry_integration owner dev_user;
create database apt_registry_test owner dev_user;

-- This lets dev_user load data from csv files into tables.
-- Required for running tests.
grant pg_read_server_files to dev_user;
```

Now load the schema into the dev database:

```
psql -U dev_user -d apt_registry_development -a -f db/schema.sql
```

If you want more data in your dev DB, we have a copy of the staging database in an undisclosed location.

[DBeaver](https://dbeaver.io/download/) is an excellent free GUI tool for interacting with the database.

# Running

`APT_ENV=dev go run registry.go`

You can change APT_ENV to test if you want to run against the test database, but note that the test DB is regenerated every time we run the test suite.

Or if you want to run in the test environment after running tests, use:

`./run.sh server`

The run.sh script starts Redis and NSQ in addition to the registry. These services are required for some functionality, such as initiating restorations and deletions and requeueing WorkItems.

You'll have some minimal data available in the DB, including a number of user accounts. You can log in with any of the following:

| Email                | Password | Role                            |
| -------------------- | -------- | ------------------------------- |
| system@aptrust.org   | password | Sys Admin                       |
| admin@inst1.edu      | password | Institutional Admin at Inst 1   |
| user@inst1.edu       | password | Institutional User at Inst 1    |
| inactive@inst1.edu   | password | Deactivated Inst User at Inst 1 |


# Testing

`./run.sh tests`
