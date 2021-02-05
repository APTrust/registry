# Implementation Notes

## Controllers

### Index

1. Construct where clause filters from query string. E.g. If User object has InstitutionID=3 and EnabledTwoFactor=true, we should filter down to users who match those criteria.

1. Forcibly apply filters if necessary. E.g. Institutional admin can view users only at their own institution, so force institution id filter in that case.

1. Run query and set result list in TemplateData.

1. Assemble other data and collections. E.g. Items that go into filter lists, such as the Institutions list on the Users page.

1. Set errors as you go.

1. Set all template data.

1. Render.

### Delete

...


### Show

...

### Create (API Only)

...

### Update (API Only)

...

## Controller Error Handling

* What is Gin's default behavior?
* What is our current behavior?

### Desired Behavior

* Add error to Gin context
* Log error
* Short-circuit if necessary: do not complete unauthorized tasks
* How is control interrupted when we want to abort?
* Return proper status code
* Display error page (web ui)
* Include specific error message in JSON response (API)

## Queries

* Building filters
* Idiomatic pg queries vs. home-grown query builder?
* Views vs tables

## Security

* Authentication is implemented in Auth middleware.
* Authorization is built into models.

How are authorization failures handled?
