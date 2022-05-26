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
- [ ] Add `active` classname to nav items when on current page
- [ ] If there are new alerts, show red dot on notifications icon, using `has-notifications` classname
- [ ] Confirm intended use of info tooltip on global search
- [ ] Make sure all form partials use Bulma, so far only the used ones have.
- [ ] Add `is-danger` classname to form fields on error
- [ ] Hook up filter chips, clearing filters functionality