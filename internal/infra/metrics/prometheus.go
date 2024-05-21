package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrderCanceled = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "order_canceled_total",
			Help: "Total number of canceled orders processed",
		},
	)
	OrderPending = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "order_pending_total",
			Help: "Total number of pending orders processed",
		},
	)
	OrderProcessed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "order_processed_total",
			Help: "Total number of processed orders processed",
		},
	)
	ProcessingErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "processing_errors_total",
			Help: "Total number of errors processing messages",
		},
	)
	ConcurrentOrdersProcessing = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "concurrent_orders_processing",
			Help: "Current number of orders being processed concurrently",
		},
	)
	OrdersEnqueued = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_enqueued_total",
			Help: "The total number of orders enqueued",
		},
	)
	OrdersDequeued = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_dequeued_total",
			Help: "The total number of orders dequeued",
		},
	)
	BufferFull = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "buffer_full_total",
			Help: "The total number of times the buffer was full",
		},
	)
	TotalValueProcessed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "total_value_processed",
			Help: "The total value of all processed orders",
		},
	)
)

func init() {
	prometheus.MustRegister(OrderCanceled)
	prometheus.MustRegister(OrderPending)
	prometheus.MustRegister(OrderProcessed)
	prometheus.MustRegister(ConcurrentOrdersProcessing)
	prometheus.MustRegister(ProcessingErrors)
}
