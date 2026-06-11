package redisprometheus

import (
	metric "github.com/luxfi/metric"

	"github.com/hanzoai/kv-go/v9"
)

// StatGetter provides a method to get pool statistics.
type StatGetter interface {
	PoolStats() *redis.PoolStats
}

// Collector collects statistics from a redis client.
// It implements the metric.Collector interface.
type Collector struct {
	getter      StatGetter
	hitDesc     *metric.Desc
	missDesc    *metric.Desc
	timeoutDesc *metric.Desc
	totalDesc   *metric.Desc
	idleDesc    *metric.Desc
	staleDesc   *metric.Desc
}

var _ metric.Collector = (*Collector)(nil)

// NewCollector returns a new Collector based on the provided StatGetter.
// The given namespace and subsystem are used to build the fully qualified metric name,
// i.e. "{namespace}_{subsystem}_{metric}".
// The provided metrics are:
//   - pool_hit_total
//   - pool_miss_total
//   - pool_timeout_total
//   - pool_conn_total_current
//   - pool_conn_idle_current
//   - pool_conn_stale_total
func NewCollector(namespace, subsystem string, getter StatGetter) *Collector {
	return &Collector{
		getter: getter,
		hitDesc: metric.NewDesc(
			metric.BuildFQName(namespace, subsystem, "pool_hit_total"),
			"Number of times a connection was found in the pool",
			nil, nil,
		),
		missDesc: metric.NewDesc(
			metric.BuildFQName(namespace, subsystem, "pool_miss_total"),
			"Number of times a connection was not found in the pool",
			nil, nil,
		),
		timeoutDesc: metric.NewDesc(
			metric.BuildFQName(namespace, subsystem, "pool_timeout_total"),
			"Number of times a timeout occurred when looking for a connection in the pool",
			nil, nil,
		),
		totalDesc: metric.NewDesc(
			metric.BuildFQName(namespace, subsystem, "pool_conn_total_current"),
			"Current number of connections in the pool",
			nil, nil,
		),
		idleDesc: metric.NewDesc(
			metric.BuildFQName(namespace, subsystem, "pool_conn_idle_current"),
			"Current number of idle connections in the pool",
			nil, nil,
		),
		staleDesc: metric.NewDesc(
			metric.BuildFQName(namespace, subsystem, "pool_conn_stale_total"),
			"Number of times a connection was removed from the pool because it was stale",
			nil, nil,
		),
	}
}

// Describe implements the metric.Collector interface.
func (s *Collector) Describe(descs chan<- *metric.Desc) {
	descs <- s.hitDesc
	descs <- s.missDesc
	descs <- s.timeoutDesc
	descs <- s.totalDesc
	descs <- s.idleDesc
	descs <- s.staleDesc
}

// Collect implements the metric.Collector interface.
func (s *Collector) Collect(metrics chan<- metric.Metric) {
	stats := s.getter.PoolStats()
	metrics <- metric.MustNewConstMetric(
		s.hitDesc,
		metric.CounterValue,
		float64(stats.Hits),
	)
	metrics <- metric.MustNewConstMetric(
		s.missDesc,
		metric.CounterValue,
		float64(stats.Misses),
	)
	metrics <- metric.MustNewConstMetric(
		s.timeoutDesc,
		metric.CounterValue,
		float64(stats.Timeouts),
	)
	metrics <- metric.MustNewConstMetric(
		s.totalDesc,
		metric.GaugeValue,
		float64(stats.TotalConns),
	)
	metrics <- metric.MustNewConstMetric(
		s.idleDesc,
		metric.GaugeValue,
		float64(stats.IdleConns),
	)
	metrics <- metric.MustNewConstMetric(
		s.staleDesc,
		metric.CounterValue,
		float64(stats.StaleConns),
	)
}
