# Alert Templates

This directory contains text templates (not HTML templates) used to format alerts. Alerts appear in the dashboard and on the alerts page, and may also be emailed to users.

We use plain text for simplicity when sending alert emails, and these templates are NOT parsed by the gin router in registry.go. These are parsed and used separately for alerts only.

Also note that text templates don't and shouldn't be enclosed in a {{ define "dir/name.txt" }} block like the HTML templates. If they are enclosed in such a block, they will parse without error but will return an empty string when executed.
