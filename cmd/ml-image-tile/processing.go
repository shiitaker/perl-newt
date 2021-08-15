
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
func processImageBimg(logger log.Logger, filePath, srcDir, dstDir string, smallerTile bool, resize, width, height int) error {
	buffer, err := bimg.Read(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	img := bimg.NewImage(buffer)
	if err != nil {
		return fmt.Errorf("can't open image %s %v", filePath, err)
	}
	size, err := img.Size()
	if err != nil {
		return fmt.Errorf("can't read sizego image %s %v", filePath, err)
	}

	if resize > 1 {
		buffer, err = img.Resize(size.Width/resize, size.Height/resize)
		if err != nil {
			return fmt.Errorf("can't resize image %s %v", filePath, err)
		}
		img = bimg.NewImage(buffer)

		size, err = img.Size()
		if err != nil {
			return fmt.Errorf("can't read resized image %s %v", filePath, err)
		}

		level.Debug(logger).Log(
			"msg", "resizing image",
			"sizex", size.Width,
			"sizey", size.Height,
			"src_path", filePath,
		)
	}

	// generate tiles starting from the center
	if size.Width < width || size.Height < height {
		return fmt.Errorf("too small to be tilled %s", filePath)
	}

	count := 0

	// start at the top left
	var xpos, ypos int

	// find how many tiles we need
	// we want the tile a c if we have enough material (at least x/2)
	// we want the tile d e if we have enough material (at least y/2)
	// a | b | c
	// d | e | f

	// a line is the number of slice + the extra half overlap
	var modx, mody int

	// do we allow repetition in tiles ?
	if smallerTile {
		modx = size.Width % width
		mody = size.Height % height
	}

	needx := size.Width / width
	needy := size.Height / height

	if modx > width/2 {
		needx += 2
	}

	if mody > height/2 {
		needy += 2
	}