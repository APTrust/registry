[![Build Status](https://travis-ci.com/APTrust/registry.svg?branch=master)](https://travis-ci.org/APTrust/registry)
[![Maintainability](https://api.codeclimate.com/v1/badges/e4c7cfd351d6bae759e3/maintainability)](https://codeclimate.com/github/APTrust/registry/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/e4c7cfd351d6bae759e3/test_coverage)](https://codeclimate.com/github/APTrust/registry/test_coverage)

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

If you're running on Linux or OSX, you'll find the required Redis and NSQ binaries in this repo's .bin/linux and ./bin/osx directories. The registry script automatically starts them when it runs the server and tests.

## Requirements for Two Factor Authentication

As a developer, you generally won't need to send Authy or SMS messages for two-factor authentication. If you're running this in a demo or production environment, you will.

The AWS SNS library, which sends two-factor auth codes via text/SMS, requires the following config files:

* ~/.aws/credentials should contain the following:

```
[default]
aws_access_key_id=<valid access key id>
aws_secret_access_key=<valid secret key>
```

* ~/.aws/config should contain the following:

```
[default]
region=us-east-2
output=json
```


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
APT_ENV=dev registry test
```

Note that this will delete and rebuild your dev DB with some starter data. You probably don't want to run this after the first time you set up your dev environment.

If you want more data in your dev DB, we have a copy of the staging database in an undisclosed location.

[DBeaver](https://dbeaver.io/download/) is an excellent free GUI tool for interacting with the database.

# Running

`APT_ENV=dev registry serve`

You can change APT_ENV to test if you want to run against the test database, but note that the test DB is regenerated every time we run the test suite.

Or if you want to run in the test environment after running tests, use:

`registry serve`

To run the server with NSQ and Redis with a different env:

`APT_ENV=dev registry serve`

The run.sh script starts Redis and NSQ in addition to the registry. These services are required for some functionality, such as initiating restorations and deletions and requeueing WorkItems.

You'll have some minimal data available in the DB, including a number of user accounts. You can log in with any of the following:

| Email                | Password | Role                            |
| -------------------- | -------- | ------------------------------- |
| system@aptrust.org   | password | Sys Admin                       |
| admin@inst1.edu      | password | Institutional Admin at Inst 1   |
| user@inst1.edu       | password | Institutional User at Inst 1    |
| inactive@inst1.edu   | password | Deactivated Inst User at Inst 1 |
| admin@inst2.edu      | password | Institutional User at Inst 2    |


# Testing

`registry test`

This drops everything in the test DB, recreates the tables and views, loads in some fixtures, and runs the unit tests. Unless you say otherwise, the script assumes APT_ENV=test.

Note that Go does not rerun tests that passed on the prior run. If you want to force all tests to run, run this before the tests themselves:

`go clean -testcache`

This may be necessary if the tests passed on the prior run, but you want to force a reload of the schema or the fixtures.

# External Services

If you want to send two-factor OTP codes through SMS, or two-factor Authy push notifications, you'll need to enable these in the .env (or .env.test) file. Set the following:

```
ENABLE_TWO_FACTOR_SMS=true
ENABLE_TWO_FACTOR_AUTHY=true
```

If these services are causing problems on your dev machine, you can turn them off by changing the settings to `false`. You can still log in with two-factor authentication when these services are turned off locally. The registry will print the OTP to the terminal console during development. You can cut and paste the OTP code from there.

You'll have to manually set a test user's phone number and/or Authy ID to send OTP messages successfully.

## AWS Environment Variables

Set the following environment variables to send OTP codes via text message through Amazon's SNS:

```
AWS_ACCESS_KEY_ID=<your key>
AWS_SECRET_ACCESS_KEY=<your secret key>
AWS_REGION="us-east-1"
```

## Authy Environment Variables

To use Authy for OTP, set the `AUTHY_API_KEY` environment variable.
