# TODO

Here's some general messiness to clean up, refactor, and/or document.

## Database Models

- [x] Use embedding and type promotion to get rid of redundant code.
- [ ] Object Get methods should return nil if there are no matching rows.
  - [ ] Alert
  - [ ] AlertView
  - [ ] Checksum
  - [ ] DeletionRequest
  - [ ] DeletionRequestView
  - [x] GenericFile
  - [ ] GenericView
  - [ ] Institution
  - [ ] InstitutionView
  - [ ] IntellectualObject
  - [ ] IntellectualObjectView
  - [ ] PremisEvent
  - [ ] PremisEventView
  - [ ] StorageOption
  - [ ] StorageRecord
  - [ ] User
  - [ ] UserView
  - [ ] WorkItem
  - [ ] WorkItemView
- [ ] Implement nested insert where appropriate. See DeletionRequest.Save()
- [ ] Validate IntellectualObject identifier by regex, as in Pharos
- [ ] Validate GenericFile identifier by regex, as in Pharos

## Testing

- [ ] Document test fixtures.
- [ ] Document fixtures vs. factory-generated models.
- [ ] Move factory out of pgmodels namespace.

## Front-end

- [ ] Research `flash` in `shared/_header`.