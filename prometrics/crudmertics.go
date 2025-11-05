package prometrics

import "time"

var (
	CrudOperationTotal    = CreateCounter("crud_operations_total", "Total CRUD operations", []string{"object", "operation"})
	CrudOperationDuration = CreateHistogram("object_operation_duration_seconds", "CRUD duration", []string{"object", "operation"}, nil)
	CrudObjectCount       = CreateGauge("object_count", "Current number of objects", []string{"object"})
)

func TrackCRUD(object, operation string) func(start time.Time) {
	start := time.Now()
	return func(_ time.Time) {
		elapsed := time.Since(start).Seconds()
		CrudOperationTotal.WithLabelValues(object, operation).Inc()
		CrudOperationDuration.WithLabelValues(object, operation).Observe(elapsed)
	}
}

func SetObjectCount(object string, count float64) {
	CrudObjectCount.WithLabelValues(object).Set(count)
}
func IncObjectCount(object string) { CrudObjectCount.WithLabelValues(object).Inc() }
func DecObjectCount(object string) { CrudObjectCount.WithLabelValues(object).Dec() }
