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
- [Two Factor Authentication](#two-factor-authentication)
- [Reporting](#reporting)
- [Testing](#testing)
- [Database Changes](#database-changes)

# Tech Stack

The primary technology stack for the new application will consist of a web application and REST API built on the [Gin HTTP Framework](https://github.com/gin-gonic/gin) and the [pg Postgres client and ORM](https://github.com/go-pg/pg).

Like Pharos, the registry will use a Postgres database and will likely sit behind a Nginx or a similar reverse proxy.

We're also using the dead-simple and utra-flexible [govalidator](https://github.com/asaskevich/govalidator) instead of Gin's built-in [go-playground validator](https://github.com/go-playground/validator) for the following reasons:

1. The specific cross-field validation our app needs is much easier to do with govalidator than with the go-playground validator.
2. Custom error messages are easier to set with govalidator.
3. Gin's uses the go-playground validator only in the web handler context, when user data is bound to models. Govalidator lets us attach validators to the models themselves, so we can call them whenever we want, in any context.

We also evaluated [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) which is built on top of govalidator, but in the end, govalidator was the simplest, richest, and most flexible. It also led to the most readable and maintainable code.

Model validation takes place on the model itself before insert and update, triggered by go-pg's built-in hooks. Validation is explicit and clear, using no reflection or hard-to-debug abstraction, so we can easily follow the logic when debugging.

Clarity is a guiding principle for this entire code base. When choosing between writing five lines of explicit, unambiguous code and using a clever one-liner from a third-party reflection library, we write the five lines. That spares future developers hours of painful debugging.

# Web UI

During development, we'll create a barebones web UI that simply renders the required data without polish. Later, we'll work toward [the UI that Simple Thread mocked up](./docs/APTrust_Wireframes.pdf). During development, however, we only need to ensure that our pages get the data that Simple Thread's UI requires.

After two years of maintaining a Node.js app with it's whole tangled dependency hell, we want to keep JavaScript to a minimum. For simple tasks like XHR loading, modals, and the like, we'll write our own vanilla JavaScript.

While it take a little more time up front, avoiding huge dependency trees saves enormous time in the long run. Don't import thousand-line libraries to do something you can accomplish in ten lines of script.

For functionality we don't want to write ourselves, such as charts, look for small, proven, dependency-free libraries like chart.js.

Remember, depdenency hell and mountains of garbage code are only one npm package away.

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

# Two Factor Authentication

Current Pharos users who have enabled two-factor authentication receive one-time passwords through SMS or push notifications through Authy OneTouch. These methods were chosen after long discussion with depositors and we cannot change them without another long discussion. So for now, we're sticking with these two.

Notes on two-factor setup and workflow have grown large enoug to warrant their own document. See [Two Factor Notes](two_factor_notes.md).

# Testing

Testing should cover all major features of the Web UI, Member API and Admin API. That is, tests should make the same endpoint requests that users make and should ensure that results are complete and correct, and that side effects (e.g. generating an email alert) are complete and correct.

In general, high-level testing should be more useful than huge suites of low-level tests.

# Reporting

Depositors and APTrust should be able to get reports on demand describing:

* total current deposits
    * file count
    * object count
    * byte count
    * by institution
    * totals for member + subaccounts
* total deposits over time (end-of-month totals 2014 - present)
    * file count
    * object count
    * byte count
    * by institution
    * totals for member + subaccounts
* deposits by storage type (standard, glacier, glacier-deep - this report can be use to calculate billing)
* deposits by region and technology
* show deleted objects/files/bytes and (ideally) when those items were deleted, and by whom
* cost breakdown
    * bytes per storage option
    * minus 10TB (taken first from Standard, then from other options)
* deposits by mime type (??)
* fixity check counts
    * per institution
    * failed fixity checks
    * drilldown?
* work item summary (number of items in period)
    * action (ingest, restoration, deletion)
    * outcome
    * drilldown?
* stalled work items (??)

The Web UI should show data and charts. The Admin API should provide reporting for all institutions, for APTrust's billing and reporting needs. The Member API should show information for the member's own institution, but not for other institutions.

Count queries can be slow in Postgres. See the section on Index Only Scans and other options at https://www.citusdata.com/blog/2016/10/12/count-performance/.

# Database Changes

## Indexes

The existing Pharos DB includes a number of unnecessary indexes created in an attempt to optimize query performance. Many are not used and are likely just hurting write performance. See, for example, the generic_files table, which has 12  columns and 19 indexes. Other tables are affected as well.

## Columns

Remove unused columns, such as generic_files.ingest_state, as well as columns that can be shared using views, rather than being duplicated. See [Views](#views) below.

As of Feb. 2021, migrations delete the following columns:

- intellectual\_objects.ingest\_state - was never used
- generic\_files.ingest\_state - was never used
- work\_items.object\_identifier - replaced in work\_items\_view
- work\_items.generic\_file\_identifier - replaced in work\_items\_view
- premis\_events.object\_identifier - replaced in premis\_events\_view
- premis\_events.generic\_file\_identifier - replaced in premis\_events\_view

## Views

The existing Pharos DB duplicates a number of fields (such as premis_events.intellectual_object_identifier and premis_events.generic_file_identifier, which are duplicated from other tables. These columns exist to speed queries that return multiple records. We can delete the columns and create views instead that provide the same columns without duplicating data. This will likely reduce the size of the DB by several gigabytes and increase insert efficiency.

Other tables with duplicate columns include work_items (object_identifier and generic_file_identifier).

## Foreign Keys

Many tables do not include properly-defined foreign key constraints. The premis_events table is an example. Its intel obj id and generic file id columns should be proper foreign keys to the related tables.

As we add foreign keys, we need to change datatypes as well. All of the 32-bit integer serial IDs and the foreign keys that refer to them will need to change to 64-bit integers.

The Pharos DB already had the following formal foreign key definitions:

- checksums.generic\_file\_id -> generic\_files.id
- storage\_records.generic\_file\_id -> generic\_files.id
- users.institution\_id -> institutions.id

The migrations file will need to add the following, pointing to the obvious places. We also need to index these fields, if they're not already indexed, because Postgres doesn't automatically index foreign keys.

- checksums.generic\_file\_id
- generic\_files.institution\_id
- generic\_files.intellectual\_object\_id
- institutions.member\_institution\_id (to institution\_id)
- intellectual\_objects.institution\_id
- premis\_events.generic\_file\_id
- premis\_events.institution\_id
- premis\_events.intellectual\_object\_id
- storage\_records.generic\_file\_id
- work\_items.generic\_file\_id
- work\_items.institution\_id
- work\_items.intellectual\_object\_id

Note that as of Feb. 2021, we're skipping changes to the bulk\_delete tables, email tables, confirmation\_tokens and some others. We'll decide later whether we keep these tables or rearchitect them.

## Tables To Drop or Rearchitect

- ar\_internal\_metadata - No longer required once we stop using ActiveRecord.
- bulk\_delete\_jobs - Likely still required, but need review to ensure the structure still serves our needs.
- bulk\_delete\_jobs\_emails - Consolidate all email tables into one.
- bulk\_delete\_jobs\_generic\_files - Probably still required, since it lists which files should be deleted by a job.
- bulk\_delete\_jobs\_institutions - Huh?
- bulk\_delete\_jobs\_generic\_files - Probably still required, since it lists which objects should be deleted by a job.
- confirmation\_tokens - WTF is this?
- emails - I get why it's there, but the structure makes no sense.
- emails\_generic\_files - Why?
- emails\_intellectual\_objects - More why?
- emails\_premis\_events - Seeing a pattern here.
- emails\_work\_items - Sigh
- old\_passwords - Like used chewing gum, probably best to discard. This is probably used to prevent users from reusing an old password. Need to discuss password policy before deciding what to do about this. If we keep it, we have to guarantee the new registry code can use it. Who knows what kind of voodoo Rails used when hashing these values?
- schema\_migrations - Not necessary after we move away from Rails.
- snapshots - Contains summary info about total deposits as of the first of each month, broken down by each depositor. God knows how this is calculated. Needs review.
- usage_samples - No idea what this is. The table is empty in the demo DB. Check the production DB, and see if there's any code to populate it or display its contents.

## Other Changes

intellectual\_obects.bagit\_profile\_identifier should probably be an integer field, pointing to a lookup table containing profile identifiers. The actual identifiers are long URLs and we don't need to repeat them 100k times.

## Counts

We require count queries in a number of places such as for API results and paged web results. Count queries are notoriously slow in Postgres because the MVCC implementation needs to run a table scan to check row visibility.

Many developers and DBAs recommend using triggers on insert and delete to update a "counts" table. We'll look into how feasible that is in our case. The difficulty lies in count queries that include where clauses. We can't store counts for every possible where clause, but we can store counts by institution, which is often what we need.

---

# Working Notes

Delete these notes when implementation is complete.

## models.DataStore

Due to lack of Go generics, we have to implement list methods individually for each type. The hard part of doing this generically is that the underlying pg library wants an interface{} type to bind results to. That interface is actually a pointer to a slice of specific model types. After binding results, we need to check permissions on every object, which is done through the Authorize method of the Model interface. Since Go will not let us cast interface{} to []*<Type> to []<ModelInterface> as needed, for each conversion, we would have to copy and rebuild every element in each slice using reflection.  That process is slow, complex, cumbersome and generally unsafe. So we bite the bullet and implement our list method individually for each type.

This class may also return select lists for different types. For example, we need lists of institutions, event types, etc. so users can filter. All we need for these lists is an id and label, not the entire object. Items should be in alpha order.
