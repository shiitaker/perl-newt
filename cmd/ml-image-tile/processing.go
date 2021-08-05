
package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"time"

	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/h2non/bimg"
)

// processImage
//  on OSX CGO_CFLAGS_ALLOW="-Xpreprocessor" go get github.com/h2non/bimg