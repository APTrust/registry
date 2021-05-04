# Workflows

In this document, "the system" means the registry. This document does not address actions that take place outside the system, such as the work done by ingest/restoration/deletion services.

These workflows assume that the user initiation the action has permission to do so. The registry checks permissions in middleware, and if that's done correctly, the user will never even reach the controller methods required to initiate these workflows.

When the system creates WorkItems in the steps outlined below, the WorkItem must always include the id or email address of the user who initiated the action and the id or email of the user who approved it.

## Object Deletion

### Step 1: Request

1. User clicks Delete button.
2. User confirms deletion in modal dialog.
3. System checks to see if this object or any of its files has pending WorkItem actions (ingest, restore, delete) or if item is already deleted (state = D)
    a. If so, system returns 409/Conflict and workflow ends.
4. If object can't be found, system returns 404.
5. System creates a deletion request in an email table. This request includes a secure confirmation token.
6. System delivers the email to institutional admins (or queues it for delivery).
7. System tells user that admins have been notified.

### Step 2: Approval

1. Institutional admin clicks the email confirmation link to review a deletion request.
2. System ensures the admin is logged in.
    a. If not logged in, redirect to login, then back to confirmation page.
3. System ensures the deletion confirmation token is valid.
    a. If token is invalid, show message and stop.
3. System checks that deletion has not already been confirmed. (Check for existing deletion WorkItem for this object that is newer than the most recent ingest.)
    a. If already confirmed, show message and stop.
4. System displays item or list of items to be deleted.
5. User clicks Confirm Delete button.
6. User confirms modal dialog message.
7. System creates confirmation email noting which user initiated and which confirmed the deletion, along with list of items to be deleted.
    a. System sends or queues email to initiator and institutional admins.
8. System creates a WorkItem to delete the object.
9. System queues the new WorkItem ID in NSQ.
10. System displays a message to the user that the item is queued for deletion.

## File Deletion

File deletion follows the same steps as object deletion.

### Bulk Deletion

In Pharos, bulk deletion required an institutional admin sending an APTrust admin a text file or spreadsheet listing all objects and files to be deleted. APTrust admin then triggered a bulk deletion endpoint with the list of identifiers. From there, the system followed the steps above, and then an APTrust admin received an email with a final confirmation link. Once the admin clicked the link, the WorkItems would be created and queued.

The registry should allow users to build a deletion list by selecting multiple objects and files. It should also include a page where users can choose to delete the items on the list. From there, the system would follow steps one and two above, but the confirmation emails would contain a list of all selected items instead of the system sending one email per item.

## Object Restoration

### Step 1: Request

1. User clicks Restore button on object detail page.
2. User confirms restoration in modal dialog.
3. If object does not exist or has been deleted, system returns error message.
4. System checks that there are no pending WorkItems (ingest, delete, restore) for the object or its files.
    a. If pending WorkItems exist, system returns an error message to user.
5. System creates WorkItem with action Restore or Glacier Restore, as necessary.
6. System queues the new WorkItem ID in NSQ.
7. System displays a confirmation message to the user.

### Step 2: Notify

1. System receives a message from external worker saying restoration is complete. (This message comes into the WorkItems endpoint via the admin API and includes the restoration bucket URL of the restored item.)
2. System creates an email (in email table) saying the restoration is complete. The message includes the restore bucket URL where the depositor can find the item.
3. System sends or queues email. It goes to original requestor (as recorded in WorkItem) and to institutional admins.

**Note:** Pharos ran restoration notifications as a cron job so it could batch multiple successful restorations into a single email instead of spamming depositors when fifty or a hundred restorations completed. Consider something similar here. The relevant Pharos code is [here ](https://github.com/APTrust/pharos/blob/master/app/controllers/work_items_controller.rb#L239-L267).

## File Restoration

File restoration follows the same steps as object restoration, except that the system checks the existence and deletion state on the file instead of the object. The system should reject the file restoration request if the file or its parent object has any pending ingest/restore/delete WorkItems.

## Spot Test Restorations

If an institution chooses to enable spot test restorations, the system will randomly restore one object each month. The system randomly chooses one intellectual object that meets the following criteria:

1. Less than 20GB total size (if possible).
2. Item has not been restored in at least six months.

The system creates a Restore or Glacier Restore WorkItem and queues the WorkItem ID in NSQ.

Note that spot tests are for object restoration only, not for file restoration. Part of the purpose of the spot tests is to ensure that all files are restored and that the BagIt package is complete and correct.

# Alerts

**TODO:** How is the current system finding these? Storing them? Sending them?

# Restructuring Database Tables

The Pharos database includes a number of tables to track emails and the items emails are related to. Some of these tables should be consolidated, and we need to preserve the data from legacy tables, possibly as a simple set of JSON records.

The existing Pharos tables include:

* bulk_delete_jobs - This tracks requestor and approver information for bulk delete jobs.
* bulk_delete_jobs_emails - Maps one bulk delete job to many recipient emails.
* bulk_delete_jobs_generic_files - Maps one bulk delete job to many files (files to be deleted).
* bulk_delete_jobs_institutions - Maps one bulk delete job to many institutions. _This table probably shouldn't exist because a bulk delete job can only belong to one institution._
* bulk_delete_jobs_intellectual_objects - Maps one bulk delete job to many objects (objects to be deleted).
* confirmation_tokens - Associates a unique confirmation token with a user, institution, object and/or generic file. These tokens are used to confirm deletions. Not sure why this is in its own table, or why it does not include a bulk deletion job id column.
* emails - Contains email body, type, recipient list and other info used by a cron job to construct and send emails. _This table needs review._
* emails_generic_files - Associates emails with generic files. **Why? Alerts?**
* emails_intellectual_objects - Associates emails with intellectual objects. **Why? Alerts?**
* emails_premis_events - Associates emails with generic files. **Why? Alerts?**
* emails_work_items - Associates emails with work items. **Why? Alerts?**


## Other Database Tables to Examine

* old_passwords - Tracks old passwords so users can't re-use them. Do we still need this? Are we still forcing user to reset passwords?
* schema_migrations - Rails-specific table to track DB migrations.
* snapshots - This is probably meant to track deposits by depositors over time. Mostly useless.
* usage_samples - No idea what this is.
