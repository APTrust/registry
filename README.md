# Registry

APTrust's registry contains metadata about the millions of items we hold on behalf of our depositors in preservation storage. The registry provides web-based and API-based access to the metadata and to information about items currently undergoing ingest and restoration.

This will be the third-generation of our registry software, based on the [Gin Web Framework](https://github.com/gin-gonic/gin). It will replace our existing Rails application, [Pharos](https://github.com/APTrust/pharos), which suffers from performance and maintability problems, which are discussed in more detail in [these notes](notes.md).

# Requirements

To run the registry on your local dev machine, you will need the following for ALL operations:

* [Go](https://golang.org/dl/) 1.13 or higher
* [Postgres](https://www.postgresql.org/download/) 11 or higher

You will also need the following for some Admin operations:

* [Redis](https://redis.io/download) 5.07 or higher
* [NSQ](https://nsq.io/deployment/installing.html) version 1.20 or higher

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

[DBeaver](https://dbeaver.io/download/) is an excellent free GUI tool for interacting with the database.

__TODO__: Set up default admin account for dev environment, or create a setup script that will load the test fixtures into the default dev DB at startup.

# Testing

`APT_ENV=test go test ./...`
