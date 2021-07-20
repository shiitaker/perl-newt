
package main

import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	fileCounter           uint64 = 0
	tileCounter           uint64 = 0
	errCounter            uint64 = 0
	rejectedBlurryCounter uint64 = 0

	//nolint deadcode
	fileCounterP = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "file_counter",
			Help: "Counts number of file processed",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&fileCounter))
		})
