package prometrics

import "time"

var (
	// CrudOperationTotal counts the total number of CRUD operations, labeled by
	// object type and operation name (e.g. "person", "create").
	//
	// Metric type: CounterVec
	CrudOperationTotal = CreateCounter("crud_operations_total", "Total CRUD operations", []string{"object", "operation"})
	// CrudOperationDuration tracks the duration of CRUD operations in seconds,
	// labeled by object type and operation name.
	//
	// Metric type: HistogramVec
	CrudOperationDuration = CreateHistogram("object_operation_duration_seconds", "CRUD duration", []string{"object", "operation"}, nil)
	// CrudObjectCount reports the current number of objects of each type.
	//
	// Metric type: GaugeVec
	CrudObjectCount = CreateGauge("object_count", "Current number of objects", []string{"object"})
)

// TrackCRUD records metrics for a CRUD operation. It should be called
// immediately before and after performing an operation.
//
// Usage pattern:
//
//	done := metrics.TrackCRUD("person", "create")
//	defer done(time.Now())
//
// The returned function observes the operation duration and increments
// the total CRUD counter.
func TrackCRUD(object, operation string) func(start time.Time) {
	start := time.Now()
	return func(_ time.Time) {
		elapsed := time.Since(start).Seconds()
		CrudOperationTotal.WithLabelValues(object, operation).Inc()
		CrudOperationDuration.WithLabelValues(object, operation).Observe(elapsed)
	}
}

// SetObjectCount sets the gauge for the given object type to a specific value.
func SetObjectCount(object string, count float64) {
	CrudObjectCount.WithLabelValues(object).Set(count)
}

// IncObjectCount increments the gauge for the given object type by 1.
func IncObjectCount(object string) { CrudObjectCount.WithLabelValues(object).Inc() }

// DecObjectCount decrements the gauge for the given object type by 1.
func DecObjectCount(object string) { CrudObjectCount.WithLabelValues(object).Dec() }
