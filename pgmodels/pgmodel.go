package pgmodels

import (
	"context"
)

type PGModel interface {
	AfterDelete(context.Context) error
	AfterInsert(context.Context) error
	AfterScan(context.Context) error
	AfterSelect(context.Context) error
	AfterUpdate(context.Context) error
	BeforeDelete(context.Context) (context.Context, error)
	BeforeInsert(context.Context) (context.Context, error)
	BeforeScan(context.Context) error
	BeforeUpdate(context.Context) (context.Context, error)
}
