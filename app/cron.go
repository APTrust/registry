package app

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

var cronJobsInitialized = false

// initCronJobs initializes
func initCronJobs(ctx *common.APTContext) {
	// We set the maintenance mode flag in Parameter Store, and it comes in through
	// the environment. Changing this setting requires an application restart.
	// So if it's true now, it's going to remain true until the app restarts.
	// If it's true, we don't want to run these cron jobs against the DB, because
	// part of maintenance my be a DB migration.
	if ctx.Config.MaintenanceMode {
		ctx.Log.Warn().Msg("Registry is NOT initializing cron jobs because system is in maintenance mode.")
		return
	}
	if !cronJobsInitialized {
		ctx.Log.Info().Msg("Initializing cron jobs")
		updateSlowCounts(ctx)
		updateCurrentDepositStats(ctx)
		updateHistoricalDepositStats(ctx)
		populateEmptyDepositStats(ctx)
		initRestorationSpotTests(ctx)
		cronJobsInitialized = true
	}
}

// updateSlowCounts runs hourly, calling our custom postgres function
// update_counts(), which refreshes materialized views that hold count
// data for our largest tables.
//
// Count queries on the big tables (IntellectualObjects, GenericFiles,
// WorkItems, and PremisEvents) are some of our most frequently run queries,
// and because of postgres' MVCC, they require very slow table scans.
// Counting PremisEvents can take 15+ seconds.
//
// To combat this, we refresh the materialized views premis_event_counts,
// intellectual_object_counts, generic_file_counts and work_item_counts
// every hour in an async go routine that will not block requests.
//
// In all our use cases, hour-old counts are tolerable. For more on these
// views, see db/migrations/001_deposit_stats.sql.
//
// If we have multiple instance of Registry running in multiple containers,
// the DB ensures that this function runs no more than once per hour.
//
// Note that the SQL function also contains a guard against multiple
// instances of Registry running the stats update at the same time.
func updateSlowCounts(ctx *common.APTContext) {
	if !cronJobsInitialized {
		go func() {
			for {
				ctx.Log.Info().Msg("cron: starting update_counts() to refresh views")
				start := time.Now().UTC()
				_, err := ctx.DB.Exec("select update_counts()")
				end := time.Now().UTC()
				duration := end.Sub(start).Seconds()
				if err != nil {
					ctx.Log.Error().Msgf("cron: update_counts failed after %f seconds: %s", duration, err.Error())
				} else {
					ctx.Log.Info().Msgf("cron: update_counts completed after %f seconds.  (Less than one second indicates counts did not need to be updated.)", duration)
				}
				time.Sleep(1 * time.Hour)
			}
		}()
	}
}

// updateCurrentDepositStats updates info about the quantity of depositor
// data in the system. This data appears on the dashboard after login,
// and in the "Reports" section. These queries take way too long to run,
// so we run them asynchronously once every hour.
//
// If we have multiple instance of Registry running in multiple containers,
// the DB ensures that this function runs no more than once per hour.
func updateCurrentDepositStats(ctx *common.APTContext) {
	if !cronJobsInitialized {
		go func() {
			// Stagger this, so it doesn't overlap with update slow counts
			time.Sleep(12 * time.Minute)
			for {
				ctx.Log.Info().Msg("cron: starting update_current_deposit_stats()")
				start := time.Now().UTC()
				_, err := ctx.DB.Exec("select update_current_deposit_stats()")
				end := time.Now().UTC()
				duration := end.Sub(start).Seconds()
				if err != nil {
					ctx.Log.Error().Msgf("cron: update_current_deposit_stats failed after %f seconds: %s", duration, err.Error())
				} else {
					ctx.Log.Info().Msgf("cron: update_current_deposit_stats completed after %f seconds. (Less than one second indicates stats did not need to be updated.)", duration)
				}
				time.Sleep(1 * time.Hour)
			}
		}()
	}
}

// updateHistoricalDepositStats ensure that the historical_deposit_stats
// table contains a snapshot of deposit stats for every month from
// APTrust's inception in 2014 until the end of the prior month.
//
// The table should already be full, based on migration 001_deposit_stats.sql.
//
// On the first of each month, we want to update the table to include stats
// through the end of the prior month. The queries needed to fill this table
// are quite expensive, so the SQL function checks the historical_deposit_stats
// table and does not try to fill in data that already exists. It just adds stats
// for the prior month.
//
// Note that this job runs every 24 hours and does nothing at all unless it's
// the first of the month.
func updateHistoricalDepositStats(ctx *common.APTContext) {
	if !cronJobsInitialized {
		go func() {
			for {
				// We usually use UTC dates in Registry, but here, we'll check whether
				// it's the first of the month in the local timezone.
				if time.Now().Day() == 1 {
					ctx.Log.Info().Msg("cron: starting populate_all_historical_deposit_stats() because it's the first of the month")
					start := time.Now().UTC()
					_, err := ctx.DB.Exec("select populate_all_historical_deposit_stats()")
					end := time.Now().UTC()
					duration := end.Sub(start).Seconds()
					if err != nil {
						ctx.Log.Error().Msgf("cron: populate_all_historical_deposit_stats failed after %f seconds: %s", duration, err.Error())
					} else {
						ctx.Log.Info().Msgf("cron: populate_all_historical_deposit_stats completed after %f seconds. (Less than one second indicates stats did not need to be updated.)", duration)
					}
				} else {
					ctx.Log.Info().Msg("cron: no need to run populate_all_historical_deposit_stats() because it's not the first of the month")
				}
				time.Sleep(24 * time.Hour)
			}
		}()
	}
}

// This fills in stats for timeline reports where depositors had no
// data in the system in a given month. We only need to run this once.
// Afterwards, it will be called on the first of each month when we
// populate historical deposit stats.
func populateEmptyDepositStats(ctx *common.APTContext) {
	if !cronJobsInitialized {
		ctx.Log.Info().Msg("cron: starting populate_empty_deposit_stats()")
		start := time.Now().UTC()
		_, err := ctx.DB.Exec("select populate_empty_deposit_stats()")
		end := time.Now().UTC()
		duration := end.Sub(start).Seconds()
		if err != nil {
			ctx.Log.Error().Msgf("cron: populate_empty_deposit_stats failed after %f seconds: %s", duration, err.Error())
		} else {
			ctx.Log.Info().Msgf("cron: populate_empty_deposit_stats completed after %f seconds. (Less than one second indicates stats did not need to be updated.)", duration)
		}
	}
}

func initRestorationSpotTests(ctx *common.APTContext) {
	if !cronJobsInitialized {
		ctx.Log.Info().Msg("cron: initializing restoration spot tests. These will run every 24 hours.")
		go func() {
			for {
				// Run restoration spot tests once a day.
				// The function will run only if needed.
				runRestorationSpotTest(ctx)
				time.Sleep(24 * time.Hour)
			}
		}()
	}
}

func runRestorationSpotTest(ctx *common.APTContext) {
	// for each inst:
	// if inst needs spot test
	// find appropriate object
	// create restoration work item
	// link obj id and work item id to inst
	//
	// later, after restoration is complete, send restoration completed alert

	if !shouldRunSpotTest(ctx) {
		return
	}

	err := spotTestLock(ctx)
	if err != nil {
		ctx.Log.Error().Msgf("runRestorationSpotTest: error setting spot test lock: %v", err)
		return
	}
	defer spotTestUnlock(ctx)

	systemUser, err := pgmodels.UserByEmail(constants.SystemUser)
	if err != nil {
		ctx.Log.Error().Msgf("runRestorationSpotTest: error getting system user: %v", err)
		return
	}

	query := pgmodels.NewQuery().Limit(100).Offset(0)
	institutions, err := pgmodels.InstitutionSelect(query)
	if err != nil {
		ctx.Log.Error().Msgf("runRestorationSpotTest: error getting institutions list for restoration spot test: %v", err)
		return
	}
	for _, inst := range institutions {
		isDue, err := inst.DueForSpotRestore()
		if err != nil {
			ctx.Log.Error().Msgf("runRestorationSpotTest: error checking whether %s is due for spot test: %v", inst.Identifier, err)
			continue
		}
		if isDue {
			scheduleSpotRestoration(ctx, inst, systemUser)
		}
	}
}

func scheduleSpotRestoration(ctx *common.APTContext, inst *pgmodels.Institution, systemUser *pgmodels.User) error {
	ctx.Log.Info().Msgf("runRestorationSpotTest: %s is due for restoration spot test (%d days)", inst.Identifier, inst.SpotRestoreFrequency)
	objView, err := pgmodels.SmallestObjectNotRestoredInXDays(inst.ID, 10000, int(inst.SpotRestoreFrequency))
	if err != nil {
		ctx.Log.Error().Msgf("runRestorationSpotTest: error getting restoration object candidate for %s: %v", inst.Identifier, err)
		return err
	}
	obj, err := pgmodels.IntellectualObjectByID(objView.ID)
	if err != nil {
		ctx.Log.Error().Msgf("runRestorationSpotTest: error getting restoration object record for %s: %v", inst.Identifier, err)
		return err
	}
	ctx.Log.Info().Msgf("runRestorationSpotTest: object %d - %s chosen for restore for %s", objView.ID, objView.Identifier, inst.Identifier)
	workItem, err := pgmodels.NewRestorationItem(obj, nil, systemUser)
	if err != nil {
		ctx.Log.Error().Msgf("runRestorationSpotTest: error creating restoration work item for %s: %v", obj.Identifier, err)
		return err
	}

	// Queue the work item in NSQ
	topic := constants.TopicObjectRestore
	if obj.IsGlacierOnly() {
		topic = constants.TopicGlacierRestore
	}
	err = ctx.NSQClient.Enqueue(topic, workItem.ID)
	if err != nil {
		// Log this error and keep going.
		// Admin can manually queue the item, if necessary.
		// We want to record the LastSpotRestoreWorkItemID.
		ctx.Log.Error().Msgf("runRestorationSpotTest: error queuing work item %d in topic %s: %v", workItem.ID, topic, err)
	}

	inst.LastSpotRestoreWorkItemID = workItem.ID
	err = inst.Save()
	if err != nil {
		ctx.Log.Error().Msgf("runRestorationSpotTest: error updating last_spot_restore_work_item_id for %s: %v", inst.Identifier, err)
	}
	return err
}

func shouldRunSpotTest(ctx *common.APTContext) bool {
	metadata, err := pgmodels.InternalMetadataByKey(constants.MetaSpotTestsRunning)
	if err != nil {
		ctx.Log.Error().Msgf("shouldRunSpotTest: error getting metadata '%s': %v", constants.MetaSpotTestsRunning, err)
		return false
	}
	if metadata.Value == "true" {
		ctx.Log.Info().Msgf("shouldRunSpotTest: no, because spot test is currently running from another process")
		return false
	}

	metadata, err = pgmodels.InternalMetadataByKey(constants.MetaSpotTestsLastRun)
	if err != nil {
		ctx.Log.Error().Msgf("shouldRunSpotTest: error getting metadata '%s': %v", constants.MetaSpotTestsLastRun, err)
		return false
	}
	lastRunDate, err := time.Parse(time.RFC3339, metadata.Value)
	if err != nil {
		ctx.Log.Error().Msgf("shouldRunSpotTest: error parsing last run date '%s': %v", metadata.Value, err)
		return false
	}
	if lastRunDate.After(time.Now().AddDate(0, 0, -1)) {
		ctx.Log.Info().Msgf("shouldRunSpotTest: no, because spot test has run within the past 24 hours at %s", metadata.Value)
		return false
	}
	ctx.Log.Info().Msgf("shouldRunSpotTest: yes, because it's not currently running and last run was at %s", metadata.Value)
	return true
}

func spotTestLock(ctx *common.APTContext) error {
	metadata, err := pgmodels.InternalMetadataByKey(constants.MetaSpotTestsRunning)
	if err != nil {
		return err
	}
	metadata.Value = "true"
	return metadata.Save()
}

func spotTestUnlock(ctx *common.APTContext) {
	metadata, err := pgmodels.InternalMetadataByKey(constants.MetaSpotTestsRunning)
	if err != nil {
		ctx.Log.Error().Msgf("spotTestUnlock: error getting spot test metadata: %v", err)
		return
	}
	metadata.Value = "false"
	err = metadata.Save()
	if err != nil {
		ctx.Log.Error().Msgf("spotTestUnlock: error releasing spot test lock on %s: %v", constants.MetaSpotTestsRunning, err)
	}

	metadata, err = pgmodels.InternalMetadataByKey(constants.MetaSpotTestsLastRun)
	if err != nil {
		ctx.Log.Error().Msgf("spotTestUnlock: error getting metadata '%s': %v", constants.MetaSpotTestsLastRun, err)
		return
	}
	metadata.Value = time.Now().UTC().Format(time.RFC3339)
	err = metadata.Save()
	if err != nil {
		ctx.Log.Error().Msgf("spotTestUnlock: error releasing spot test lock on %s: %v", constants.MetaSpotTestsLastRun, err)
	}
}
