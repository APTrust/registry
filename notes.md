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

After rewriting APTrust's [Preservation Services](https://github.com/APTrust/preservation-services) to increase throughput, load testing showed that Pharos always became unresponsive under loads we expect to encounter often in production. Rails memory and CPU usage was so high that we could not even SSH into our servers for days. The only way to ensure the stability of the system was to throttle the services' access to the Rails app, which defeats the point of rewriting those services in the first place. We have seen similar behavior from Pharos in production for more than two years, and we have dealt with it by 1) running expensive overpowered hardware (whose memory and CPU still gets maxed out) and 2) throttling our clients' access to the Rails app.

### Maintainability (Explicitness)

While Rails is an excellent platform for rapidly developing database-driven web applications and APIs, the cost of changing code and adding features is often much higher than in a codebase built on the principle of explicitness. Rails' "magic" initially hides tremendous complexity from the developer, but developers must often dig into that complexity to make seemingly simple changes. In eleven years of Rails work, I have consistently run into unexpected behaviors that require hours of research to understand what Rails is doing under the hood.

It's common to spend four hours or even an entire day wading through documentation and source code to figure out the unexpected consequences of changing a single line of code. Choosing Rails means committing your developers to years of difficult debugging.

These problems rarely occur in languages and frameworks that enforce explicitness. While explicit codebases require more up-front effort, their transparency and lack of assumptions make them far easier to maintain in the long run.

### Maintainability (Dependencies)

APTrust's projects that include large numbers of dependencies require considerable ongoing work to keep dependencies up to date. Even with Dependabot scanning for vulnerabilities and outdated packages, we still spend developer hours every week merging pull requests and ensuring tests and builds work across all platforms.

An empty Rails 6 project includes 18,688 files. Then yarn pulls in 770 node packages, giving us a total of 21,254 files before we even write a line of code.

We know we don't need many of these files. We know we don't need ANY of 770 node modules. We can spend a few days now weeding them all out, or we can leave them in place and update them every week when Dependabot nags us.

(In practice, removing files from a default Rails installation causes problems later, as new gems added later in the development cycle, depend on the presence of those default files. We'd likely have to add back many of the files we remove.)

We estimate that the using Go's [Gin Framework](https://github.com/gin-gonic/gin), we can implement all of our required featured in about 100-150 code files. That's 21,000 fewer files to maintain.

### Maintainability (Testing)

Our Rails applications have traditially included thousands of tests. Many of those simply test things a compiler will catch at build time, such as whether types and parameter lists are correct, whether values are (or even can be) nil, etc. Maintaining tests requires work. If we can remove a few hundred tests, we'll have that much less to maintain.

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

During development, we'll create a barebones web UI that simply renders the required data without polish. Later, we'll work toward [the UI that Simple Thread mocked up](./docs/APTrust_Wireframes.pdf). During development, however, we only need to ensure that our pages get the data that Simple Thread's UI requires.

Consider using some Tailwind themes, such as:

* [Semantic UI](https://semantic-ui.com/) may be the best of the lot. Very clean and simple. It's well documented and its [GitHub repo has 49k stars](https://github.com/semantic-org/semantic-ui). It's nearly as heavy as Bootstrap and includes it's own JavaScript and themes. Looks like maintenance ended in 2018. Take a look at [Formantic-UI](https://fomantic-ui.com/), which is the community fork. [Formantic on GitHub](https://github.com/fomantic/fomantic-ui).
* [UIKit](https://getuikit.com/docs/introduction) also looks top-notch. Very clean and feature-rich, though the class markup is a little heavier than Semantic UI. Development seems active on [GitHub](https://github.com/uikit/uikit) and it has a simple static distribution (i.e. does not rely on Node or SASS pre-processing).
* [Bulma](https://bulma.io/) is also simple and clean, and is easy to customize. It neither includes nor requires JavaScript. It has a number of [extensions](https://bulma.io/extensions/) and over [42k stars on GitHub](https://github.com/jgthms/bulma). The default colors are too pastel, but should be easy to change. Bulma's one drawback (compared to Bootstrap) is weaker accessibility support.
* [Boostrap](https://getbootstrap.com/) is useful, well documented, well-supported, and often a pain in the ass. It also requires JQuery, which is another pain in the ass.
* [Tailwind Admin](https://github.com/tailwindadmin/admin) - See the [live demo](https://tailwindadmin.netlify.app/index.html) for component samples, but note that it includes heavy class markup and 56k lines of CSS.
* [Zurb's stupid website](https://get.foundation/) isn't even working. That bodes ill.


## Accounts

### Login

* Email/password login
* Two-factor text/sms
* Two-factor Authy

To ensure users won't have to change their passwords when moving from the Rails app, implement the same password encryption scheme as Devise. The scheme is described [here](https://www.freecodecamp.org/news/how-does-devise-keep-your-passwords-safe-d367f6e816eb/), and the [Go bcrypt library](https://pkg.go.dev/golang.org/x/crypto/bcrypt) should be able to support it.

For two-factor auth, since we're already using Authy, try the [Go Client for Authy](https://github.com/dcu/go-authy).

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
* Combine filters

### Admin Features
* Edit item
* Requeue item

## Alerts

* list with sort, filter, and paging
* view
* mark read/unread

## Web UI Admin Features

The following features will be accessible only to APTrust admins.

### AWS Account Management

The admin panel will implement these features, which currently exist only in Ansible playbooks. This will vastly simplify common tasks associated with users and institutions.

* New organization setup will create required buckets.
* Admin can add and remove IAM users for each institution.
* Admin can create and deactivate IAM user keys.

### NSQ

We currently manage NSQ through its built-in web UI. We've configured it to allow access only from whitelisted IP addresses, and we must often change the IPs in the whitelist to maintain access. Accessing NSQ through the registry will ensure that 1) we can access it from anywhere and 2) only valid APTrust admins can access it.

* View current NSQ status for all topics, channels, and hosts.
* Pause any topic or channel.
* Unpause any topic or channel.
* Empty any topic or channel.

### Redis

Admin UI will show interim processing data in Redis. This helps us understand what's being ingested and whether there are problems in specific workers. All data in these views is read-only.

* View object list.
* View files list (list of all files belonging to an object)
* View file details (full report on individual file)

### Logs

Log viewing and searching will require an additional service in preservation services to tail and search logs. This feature gives admins access to worker logs and registry logs only, not to any other files.

* Admin can tail any log for any worker on any host.
* Admin can search specific log.
* Admin can search all logs across all hosts.
* Admin can dowload any worker or registry log from any host.

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
