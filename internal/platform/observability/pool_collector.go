package observability

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type DbPoolCollector struct {
	db      *sql.DB
	inUse   *prometheus.Desc
	idle    *prometheus.Desc
	maxOpen *prometheus.Desc
}

func NewDbPoolCollector(db *sql.DB) *DbPoolCollector {
	return &DbPoolCollector{
		db:      db,
		inUse:   prometheus.NewDesc("db_in_use_connections", "Number of connections currently in use", nil, nil),
		idle:    prometheus.NewDesc("db_idle_connections", "Number of idle connections", nil, nil),
		maxOpen: prometheus.NewDesc("db_max_open_connections", "Maximum number of open connections", nil, nil),
	}
}

func (c *DbPoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.inUse
	ch <- c.idle
	ch <- c.maxOpen
}

func (c *DbPoolCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()
	ch <- prometheus.MustNewConstMetric(c.inUse, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.idle, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.maxOpen, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
}
