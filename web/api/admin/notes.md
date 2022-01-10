# Admin API Notes

## Alerts

[x] List - uses common API
[x] Show - uses common API

## Deletion Requests

[x] List - uses common API
[x] Show - uses common API

## Generic Files

[x] List - uses common API
[x] Show - uses common API
[x] Create
[x] Update
[ ] Batch Create
[ ] Delete

Create should support bulk insert.

## Institutions

[x] List - admin API only
[x] Show - admin API only

## Intellectual Objects

[x] List - uses common API
[x] Show - uses common API
[x] Create - admin API only
[x] Update - admin API only
[x] Delete - admin API only

Delete needs to ensure all files are already marked as deleted and must create a deletion premis event.

[ ] Ensure all files deleted
[ ] Create deletion premis event

## Premis Events

[x] List - uses common API
[x] Show - uses common API
[x] Create - admin API only

## Work Items

[x] List - uses common API
[x] Show - uses common API
[x] Create - admin API only
[x] Update - admin API only
