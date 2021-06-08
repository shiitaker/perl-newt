
package main

import (
	"context"
	"fmt"
	"io/fs"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gocv.io/x/gocv"
	"golang.org/x/sync/errgroup"
)

var (
	version = "no version from LDFLAGS"

	source = flag.String("source", "", "Source directory for the images")
	dest   = flag.String("dest", "", "Destination directory for the images")
	width  = flag.Int("width", 400, "Size of the target tiles x")
	height = flag.Int("height", 400, "Size of the target tiles y")

	resize               = flag.Int("resize", 2, "Divide size tilling")
	smallerTile          = flag.Bool("smallerTile", false, "Allow tiling of remaining on the borders")
	workerCount          = flag.Int("workerCount", 8, "Parallel worker count")
	validationTileCount  = flag.Int("validationTileCount", 0, "Number of validation tiles")
	validationOnly       = flag.Bool("validationOnly", false, "Generate validation tiles only")
	rejectBlurry         = flag.Bool("rejectBlurry", false, "Reject blurry source image")
	rejectBlurryThresold = flag.Float64("rejectBlurryThresold", 6_000, "Thresold before rejecting blurry images")
	logLevel             = flag.String("logLevel", "INFO", "DEBUG|INFO|WARN|ERROR")
	httpMetricsPort      = flag.Int("httpMetricsPort", 34130, "http port")

	httpMetricsServer *http.Server
)

func main() {
	flag.Parse()

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "caller", log.Caller(5), "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "app", "ml-image-tile")
	logger = NewLevelFilterFromString(logger, *logLevel)

	stdlog.SetOutput(log.NewStdlibAdapter(logger))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// catch termination
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	// web server metrics
	g.Go(func() error {
		httpMetricsServer = &http.Server{
			Addr:         fmt.Sprintf(":%d", *httpMetricsPort),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		level.Info(logger).Log("msg", fmt.Sprintf("HTTP Metrics server listening at :%d", *httpMetricsPort))

		versionGauge.WithLabelValues(version).Add(1)
