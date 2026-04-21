package db

import (
	"context"
	"time"

	infra "nh-be/infra/observability"

	"gorm.io/gorm"
)

type dbMetricsKey struct{}

type DbMetricsPlugin struct{}

func (p *DbMetricsPlugin) Name() string { return "dbMetrics" }

func (p *DbMetricsPlugin) Initialize(db *gorm.DB) error {
	// Register Before callbacks
	db.Callback().Create().Before("gorm:create").Register("metrics:before_create", before)
	db.Callback().Query().Before("gorm:query").Register("metrics:before_query", before)
	db.Callback().Update().Before("gorm:update").Register("metrics:before_update", before)
	db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", before)

	// Register After callbacks
	db.Callback().Create().After("gorm:create").Register("metrics:after_create", after("create"))
	db.Callback().Query().After("gorm:query").Register("metrics:after_query", after("query"))
	db.Callback().Update().After("gorm:update").Register("metrics:after_update", after("update"))
	db.Callback().Delete().After("gorm:delete").Register("metrics:after_delete", after("delete"))

	return nil
}

func before(db *gorm.DB) {
	db.Statement.Context = context.WithValue(db.Statement.Context, dbMetricsKey{}, time.Now())
}

func after(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if start, ok := db.Statement.Context.Value(dbMetricsKey{}).(time.Time); ok {
			infra.DbQueryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
		}
	}
}
