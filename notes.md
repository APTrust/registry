# APTrust Registry

This repository will contain the source code for the new APTrust registry, replacing [Pharos](https://github.com/APTrust/pharos). We are replacing Pharos primarily because the code has become unmaintainable.

## Why Rewrite?

* The old code is such a horrible mess that rewriting is easier than incremental cleanup.
* It's unclear _why_ some elements of the old code work, much less _how_.
* The fact the key features do work is verified primarily by integration tests in the external [Exchange](https://github.com/APTrust/exchange) and [Preservation Services](https://github.com/APTrust/preservation-services) repositories.
* Changes to the old code result in unpredictable and potentially destructive consequences. In fact, the existing code is so brittle that virtually any change results in multiple cascading failures in the spec tests.
* Elements of the architecture and underlying database schema of some non-core features (such as alerts, emails, and the bulk delete workflow) need to be fully analyzed and possibly redesiged.
* We need to clarify what the code is doing by simplifying it and making everything explicit. This alone would require a rewrite.

## Rationale for Choosing Gin over Ruby on Rails

We had considered writing this replacement in Ruby on Rails, but decided to use Go's Gin framework for the following reasons:

### Performance

After rewriting APTrust's [Preservation Services](https://github.com/APTrust/preservation-services) to increase throughput, load testing showed that Pharos always became unresponsive under loads we expect to encounter often in production. Rails memory and CPU usage was so high that we could not even SSH into our servers for days. The only way to ensure the stability of the system was to throttle the services' access to the Rails app, which defeats the point of rewriting those services in the first place. We have seen similar behavior from Pharos in production for more than two years, and we have dealt with it by throttling the clients that access it.

### Maintainability (Explicitness)

While Rails is an excellent platform for rapidly developing database-driven web applications and APIs, the cost of changing code and adding features is often much higher than in a codebase built on the principle of explicitness over assumption. Rails' "magic" hides tremendous complexity from the developer, and developers must often dig into that complexity to make seemingly simple changes. In eleven years of Rails work, I have consistently run into unexpected behaviors that require hours of research to understand what Rails is doing under the hood.

These problems rarely occur in languages and frameworks that enforce explicitness. While explicit codebases require more up-front effort, their transparency and lack of assumptions make them far easier to maintain in the long run.

### Maintainability (Dependencies)

APTrust's projects that include large numbers of dependencies require considerable ongoing work to keep dependencies up to date. Even with Dependabot scanning for vulnerabilities and outdated packages, we still spend developer hours every week merging pull requests and ensuring tests and builds work across all platforms.

An empty Rails 6 project includes 18,688 files. Then yarn pulls in 770 node packages, giving us a total of 21,254 files before we even write a line of code.

We know we don't need many of these files. We know we don't need ANY of 770 node modules. We can spend a few days now weeding them all out, or we can leave them in place and update them every week when Dependabot nags us.

We estimate that the using Go's [Gin Framework](https://github.com/gin-gonic/gin), we can implement all of our required featured in about 100-150 code files. That takes about 21,000 items off our to-do list.

### Maintainability (Testing)

Our Rails applications have traditially included hundreds or thousands of tests. Many of those simply test things a compiler will catch at build time, such as whether types and parameter lists are correct, whether values are (or even can be) nil, etc.

# Features

Below is a barebones functional spec that the registry app will need to implement to match what Pharos did. With the exception of two-factor authentication, these features should be relatively easy to implement.

The new registry will also include a number of reports and administrative features that had never been implemented in Pharos.

The registry will include three sets of routes, leading to three sets of endpoints:

1. __Web UI__ - Used by both members and APTrust admins, the web UI supports almost exclusively read-only operations.
2. __Member API__ - The member API allows members to query objects, files, and events in the registry, as well as the items in progress. The Member API will be almost identical to the current API described in this [Swagger Doc](https://aptrust.github.io/pharos/), though it will clean up some inconsistencies in parameter names and will expand some search capabilities.
3. __Admin API__ - APTrust's [preservation services](https://github.com/APTrust/preservation-services) performs both read and write operations through the Admin API. The registry will implement all of the features in the Pharos admin API. Those features have never been formally documented. They will be outlined here.

# Contents

- [Tech Stack](#tech-stack)
- [Web UI](#web-ui)
- [Member API](#member-api)
- [Admin API](#admin-api)
- [Roles and Security](#roles-and-security)
- [Testing](#testing)
- [Database Changes](#database-changes)

# Tech Stack

The primary technology stack for the new application will consist of a web application and REST API built on the [Gin HTTP Framework](https://github.com/gin-gonic/gin) and the [pg Postgres client and ORM](https://github.com/go-pg/pg).

Like Pharos, the registry will use a Postgres database and will likely sit behind a Nginx or a similar reverse proxy.

# Web UI

## Accounts

### Login

* Email/password login
* Two-factor text/sms
* Two-factor Authy

### Edit

* edit details (phone, etc.)
* reset password

### Logout

* clear session auth cookie

## Institutions

* create, list, edit, enable/disable (admin only)

## Intellectual Objects

* list with sort, filter, and paging
* view
* request delete
* request bulk delete
* request restore

## Generic Files

* list with sort, filter, and paging
* view
* request delete
* request bulk delete
* request restore

## Premis Events
* list with sort, filter, and paging
* view

## User Management

* list with sort, filter, and paging
* view
* edit (admin and inst admin only)
* disable (admin and inst admin only)

## Checksums

* list with filter, sort, paging
* view (by identifier or id)
* create
* update
* delete

## WorkItems

###  Default Listing
* Filtered by user institution (or none for Admin)
* Ordered by date
* Paged, with 25 or so items per page

### Custom Listing
* Sort on any field
* Filter on any field
* Combile filters

## Alerts

* list with sort, filter, and paging
* view
* mark read/unread

# Member API

## Authentication

* via API token

## Intellectual Objects

* list with filter, sort, paging
* view (by identifier, possibly also by id)

## Generic Files

* list with filter, sort, paging
* view (by identifier, possibly also by id)

## Premis Events

* list with filter, sort, paging
* view (by identifier, possibly also by id)

## Checksums

* list with filter, sort, paging
* view (by identifier or id)
* create
* update
* delete

## Work Items

* list with filter, sort, paging
* view (by id)

## Alerts

* list with filter, sort, paging
* view (by id)

# Admin API

## Authentication

* by API token

## Institutions

* list with filter, sort, paging

## Intellectual Objects

* list with filter, sort, paging
* view (by identifier or id)
* create
* update
* soft delete

## Generic Files

* list with filter, sort, paging
* view (by identifier or id)
* create
* __bulk create__ (with checksums and events)
* update
* soft delete

## Premis Events

* list with filter, sort, paging
* view (by identifier or id)
* create
* __no delete allowed__
* __no updating allowed__

## Checksums

* list with filter, sort, paging
* view (by identifier or id)
* create
* update
* delete

## Work Items

* list with filter, sort, paging
* view (by identifier or id)
* create (??)
* update

## Alerts

* Preservation services workers likely will not create these.
* Preservation services workers will not consume there.
* Best to leave these in the Web UI and Member API only.

# Roles and Security

The term "items" below refers to Intellectual Objects, Generic Files, Checksums, Premis Events and Work Items.

## Institutional User

* can view items belonging to their institution
* can request deletion of objects and files belonging to their insitution (subject to approval of institutional admin)
* can request restoration of files and objects belonging to their institution
* can edit elements of their own user account, including phone number and password
* can generate API keys for themselves

## Institutional Admin

* has all institutional user privileges, plus:
* can add, deactivate, and reactivate users at their own institution
* can turn two-factor authentication on and off for their institution
* can approve or reject file and object deletions (including bulk deletions)

## System Admin

* has all the powers of institutional admin across all institutions, plus:
* can create, edit, enable and disable institutions
* can edit and requeue WorkItems
* can access admin features related to external services, such as
  * work queues
  * ingest and restoration services
  * interim processing data in redis
  * AWS (IAM, S3 and Glacier)
  * Wasabi

# Testing

Testing should cover all major features of the Web UI, Member API and Admin API. That is, tests should make the same endpoint requests that users make and should ensure that results are complete and correct, and that side effects (e.g. generating an email alert) are complete and correct.

In general, high-level testing should be more useful than huge suites of low-level tests.

# Database Changes

## Indexes

The existing Pharos DB includes a number of unnecessary indexes created in an attempt to optimize query performance. Many are not used and are likely just hurting write performance. See, for example, the generic_files table, which has 12  columns and 19 indexes. Other tables are affected as well.

## Columns

Remove unused columns, such as generic_files.ingest_state, as well as columns that can be shared using views, rather than being duplicated. See [Views](#views) below.

## Views

The existing Pharos DB duplicates a number of fields (such as premis_events.intellectual_object_identifier and premis_events.generic_file_identifier, which are duplicated from other tables. These columns exist to speed queries that return multiple records. We can delete the columns and create views instead that provide the same columns without duplicating data. This will likely reduce the size of the DB by several gigabytes and increase insert efficiency.

Other tables with duplicate columns include work_items (object_identifier and generic_file_identifier).

## Counts

We require count queries in a number of places such as for API results and paged web results. Count queries are notoriously slow in Postgres because the MVCC implementation needs to run a table scan to check row visibility.

Many developers and DBAs recommend using triggers on insert and delete to update a "counts" table. We'll look into how feasible that is in our case. The difficulty lies in count queries that include where clauses. We can't store counts for every possible where clause, but we can store counts by institution, which is often what we need.
