# Pharos Database Migrations

This folder contains legacy files used during the initial migration from 
Pharos to Registry. We ran these migrations extensively during development
and testing. We ran them on the production Pharos DB on November 2, 2022,
converting it to the structure required by the Registry.

**There is no need ever to run these files again.** They are included for 
archival purposes, so we know how the DB changed when we moved from Pharos
to Registry.

## schema.sql

This file contains the final Pharos schema as it existed before we applied
the registry migrations. Note the presence of Rails-specific tables like
`ar_internal_metadata` which ActiveRecord stored to track its own migration
info.

## migrations.sql

This file contains the structural changes applied to the final Pharos DB
to convert it to the Registry DB. Note that it also moves and updates some
data, particularly in the `generic_files` table.

## views.sql

This file defines the Registry DB's new views, as well as functions that
work on views. The contents of this file could have been included in the
migrations.sql file, but the views file changed so often during development,
it was easier to keep separate.
