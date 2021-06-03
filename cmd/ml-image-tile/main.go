
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