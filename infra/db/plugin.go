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
	if err := db.Callback().Create().Before("gorm:create").Register("metrics:before_create", before); err != nil {
		return err
	}
	if err := db.Callback().Query().Before("gorm:query").Register("metrics:before_query", before); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register("metrics:before_update", before); err != nil {
		return err
	}
	if err := db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", before); err != nil {
		return err
	}

	// Register After callbacks
	if err := db.Callback().Create().After("gorm:create").Register("metrics:after_create", after("create")); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register("metrics:after_query", after("query")); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("metrics:after_update", after("update")); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("metrics:after_delete", after("delete")); err != nil {
		return err
	}

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
