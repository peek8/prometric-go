package prometrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// type MetricDefinition struct {
// 	Name    string    `yaml:"name"`
// 	Type    string    `yaml:"type"`
// 	Help    string    `yaml:"help"`
// 	Labels  []string  `yaml:"labels"`
// 	Buckets []float64 `yaml:"buckets,omitempty"`
// }

// type MetricConfig struct {
// 	Metrics []MetricDefinition `yaml:"metrics"`
// }

type MetricFactory struct {
	mu         sync.Mutex
	counters   map[string]*prometheus.CounterVec
	gauges     map[string]*prometheus.GaugeVec
	histograms map[string]*prometheus.HistogramVec
}

var factory = &MetricFactory{
	counters:   make(map[string]*prometheus.CounterVec),
	gauges:     make(map[string]*prometheus.GaugeVec),
	histograms: make(map[string]*prometheus.HistogramVec),
}

func CreateCounter(name, help string, labels []string) *prometheus.CounterVec {
	factory.mu.Lock()
	defer factory.mu.Unlock()
	if c, ok := factory.counters[name]; ok {
		return c
	}
	c := prometheus.NewCounterVec(prometheus.CounterOpts{Name: name, Help: help}, labels)
	prometheus.MustRegister(c)
	factory.counters[name] = c
	return c
}

func CreateGauge(name, help string, labels []string) *prometheus.GaugeVec {
	factory.mu.Lock()
	defer factory.mu.Unlock()
	if g, ok := factory.gauges[name]; ok {
		return g
	}
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labels)
	prometheus.MustRegister(g)
	factory.gauges[name] = g
	return g
}

func CreateHistogram(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec {
	factory.mu.Lock()
	defer factory.mu.Unlock()
	if h, ok := factory.histograms[name]; ok {
		return h
	}
	if len(buckets) == 0 {
		buckets = prometheus.DefBuckets
	}
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: name, Help: help, Buckets: buckets}, labels)
	prometheus.MustRegister(h)
	factory.histograms[name] = h
	return h
}

// func LoadMetricsFromYAML(path string) error {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return fmt.Errorf("read file: %w", err)
// 	}
// 	var cfg MetricConfig
// 	if err := yaml.Unmarshal(data, &cfg); err != nil {
// 		return fmt.Errorf("parse YAML: %w", err)
// 	}
// 	for _, m := range cfg.Metrics {
// 		switch m.Type {
// 		case "counter":
// 			CreateCounter(m.Name, m.Help, m.Labels)
// 		case "gauge":
// 			CreateGauge(m.Name, m.Help, m.Labels)
// 		case "histogram":
// 			CreateHistogram(m.Name, m.Help, m.Labels, m.Buckets)
// 		default:
// 			return fmt.Errorf("unknown metric type: %s", m.Type)
// 		}
// 	}
// 	return nil
// }
