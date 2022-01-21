[![Build Status](https://travis-ci.com/APTrust/registry.svg?branch=master)](https://travis-ci.com/APTrust/registry)
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

Note that in the test fixtures, the API secret key for all users is `password`.

# Testing

`registry test`

This drops everything in the test DB, recreates the tables and views, loads in some fixtures, and runs the unit tests. Unless you say otherwise, the script assumes APT_ENV=test.

Note that Go does not rerun tests that passed on the prior run. If you want to force all tests to run, run this before the tests themselves:

`go clean -testcache`

This may be necessary if the tests passed on the prior run, but you want to force a reload of the schema or the fixtures.

## Fixtures

The test script automatically loads fixtures when tests begin. Fixtures for Inst 1 and Inst 2 are stable. Tests don't add objects, files, etc. to those institutions. That means you can test for a certain number of items (e.g. the index page should return 4 objects, or 6 files) and you should get that number.

Items added dynamically during tests are added to the test.edu institution. Counts of these items may change between tests, so don't rely on them.

If you want to ensure a 100% known dataset before any test, call `db.ForceFixtureReload()` at the top of your test fuction. If your test is going to pollute the DB, call `defer db.ForceFixtureReload()` at the outset, so your test will clean up after itself.

The only reason we don't call `db.ForceFixtureReload()` for every test is because it will slow down the test suite.

## HTTP Tests

We use the httpexpect test library, which sometimes panics when calling Expect(). That seems to be a bug. It can also send repeat requests if you call Expect() and then later call Expect().Body(). The repeated requests mangle the URL, causing errors. That also seems to be a bug.

You can call `testutil.InitHTTPTests()` before any HTTP tests to ensure HTTP test clients are ready to use. This call is idempotent, and will only initialize clients that have not yet been initialized, so it's always safe to call.

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

# Notes on Routes

Registry routes deviate somewhat from standard REST patterns for two reasons.

1. When this project started, conflicts in the Gin router prevented the use of some standard REST patterns. See  [issue #1681](https://github.com/gin-gonic/gin/issues/1681), which has since been resolved.
2. Gin did not include the Rails `_method` hack, which told the backend to interpret a POST as a PUT, DELETE, or other method. While XHR and API clients support all HTTP methods, browsers only support GET and POST, so routes in the Web UI use those methods.
3. While a typical web app nests routes under a resource's parent ID, Registry nests sometimes nests routes under `institution_id`. This is because all of the Registry's permissions are based on a combination of user role and institution id.

For example, GenericFile is a child of IntellectualObject. A typical nested REST route to create a batch of GenericFiles would be `/files/create_batch/:object_id`. Instead, we use `/files/create_batch/:institution_id`. This helps the authorization middleware ensure that the current user has permission to perform the given action at the specified institution.

# Security Implementation

Users must be logged in to use the system. If a user is not logged in, the authentication middleware will redirect them to the login page.

The Registry includes three pieces of middleware that execute before the target handler for virtually all routes:

    * middleware/authenticate.go reads an encrypted cookie and sets the value of CurrentUser so that all subsequent code in the current request can access it.
    * middleware/authorize.go figures out the resource being requested, the institution that owns the resource, and the action the user wants to perform on that resource. It checks middleware/authorization_map.go to see if the user is allowed to perform that action on that resource. If so, the request proceeds. If not, the user gets a 403 error. Note that permissions are hard-coded in the authorization_map instead of being stored in the DB. This prevents hackers from elevating privileges by manipulating the DB. APTrust permissions are fixed and documented. There is no need to change them dynamically.
    * middleware/csrf.go provides CSRF protection for unsafe methods coming through the Web UI. Any method other than GET, HEAD, OPTIONS, and TRACE is considered unsafe.

With the exception of a few whitelisted pages, the middleware will not permit any requests to proceed unless an authorization check has been performed and the user has passed. "Missing authorization check" is a fatal error that immediately aborts the request.

Routes excepted from auth checks include:

* The login page
* The log off page
* The "forgot password" landing page, which comes from a link in the "forgot password" email
* The "complete password reset" page, which comes from a link in the "password reset" email
* The general error page
* Static files, such as scripts, stylesheets, images, favicon, etc.

The "forgot password" and "password reset" pages both use a secure token in the query string to identify the user.

The authentication middleware also hijacks all requests in cases where a user must complete registration activities. For example, if a user needs to reset their password, confirm an Authy account, or provide the second factor for multi-factor auth, the middleware will not let them past the reset/confirmation/token page until they've provided the required info.
