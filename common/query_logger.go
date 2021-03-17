package common

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
)

type QueryLogger struct {
	log zerolog.Logger
}

func NewQueryLogger(logger zerolog.Logger) *QueryLogger {
	return &QueryLogger{
		log: logger,
	}
}

func (l *QueryLogger) BeforeQuery(c context.Context, qe *pg.QueryEvent) (context.Context, error) {
	sql, _ := qe.FormattedQuery()
	l.log.Debug().Msgf("Starting SQL: %s", string(sql))
	return c, nil
}

func (l *QueryLogger) AfterQuery(c context.Context, qe *pg.QueryEvent) error {
	//sql, _ := qe.FormattedQuery()
	//l.log.Debug().Msgf("Finished SQL: %s", string(sql))
	return nil
}
