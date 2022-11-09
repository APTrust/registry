package common

import (
	"time"
)

var cronJobsInitialized = false

// initCronJobs initializes
func initCronJobs(ctx *APTContext) {
	if !cronJobsInitialized {
		ctx.Log.Info().Msg("Initializing cron jobs")
		updateSlowCounts(ctx)
		updateDepositStats(ctx)
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
func updateSlowCounts(ctx *APTContext) {
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

// updateDepositStats updates info about the quantity of depositor
// data in the system. This data appears on the dashboard after login,
// and in the "Reports" section. These queries take way too long to run,
// so we run them asynchronously once every hour.
//
// If we have multiple instance of Registry running in multiple containers,
// the DB ensures that this function runs no more than once per hour.
func updateDepositStats(ctx *APTContext) {
	if !cronJobsInitialized {
		go func() {
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
