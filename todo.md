# TODO

Here's some general messiness to clean up, refactor, and/or document.

## Database Models

[ ] Use embedding and type promotion to get rid of redundant code.
[ ] Object Get methods should return nil if there are no matching rows.
[ ] Implement nested insert where appropriate. See DeletionRequest.Save()

## Testing

[ ] Document test fixtures.
[ ] Document fixtures vs. factory-generated models.
[ ] Move factory out of pgmodels namespace.
