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
// views, see db/views.sql.
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
					ctx.Log.Info().Msgf("cron: update_counts completed after %f seconds", duration)
				}
				time.Sleep(1 * time.Hour)
			}
		}()
	}
}
