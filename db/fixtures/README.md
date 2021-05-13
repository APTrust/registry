# Registry Fixtures

The fixtures in this directory are used in unit and integration tests. They are loaded by the LoadFixtures function in db/test_util.go, which loads them only once, even if it's called repeatedly.

## Users

The password for all users in the user.csv file is `password`.

## DeletionRequests

For all deletion requests in deletion_requests.csv, the plaintext value of the encrypted confirmation token is `ConfirmationToken`. The plaintext value of the cancellation token is `CancelToken`.

| ID  | Description |
| --- | ----------- |
| 1   | Requested by user user@inst1.edu, awaiting approval or cancellation. |
| 2   | Requested by user user@inst1.edu, approved by admin@inst1.edu.       |
| 3   | Requested by user user@inst1.edu, cancelled by admin@inst1.edu.      |
