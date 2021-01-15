# Database Changes

- [x] Add unique constraint on roles.name

## Multi-Step Involved Changes

### Move Role into Users Table

Currently a user has one institution (users.institution\_id) and can have multiple roles through roles\_users to the roles table. This makes no sense. Because of our business rules, a user can have only one role at one institution. To fix this:

1. Add column users.role.
2. Copy each user's role from user->roles_users->role.name to users.role.
3. Drop table users_roles.
4. Drop table roles.

We may have to apply this change back in the feature/storage-record branch of Pharos, or if that's too difficult due to the brittleness of the old Rails code, apply it stages here. If in stages, we would do items 1 and 2 above, and save items 3 and 4 (dropping those tables) until after the registry goes into production.
