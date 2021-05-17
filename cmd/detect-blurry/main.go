
package main

import (
	"fmt"
	"log"

	"github.com/namsral/flag"
	"gocv.io/x/gocv"
)

var source = flag.String("source", "", "Source image")

const thresold = 7_000

func main() {
	flag.Parse()

	img := gocv.IMRead(*source, gocv.IMReadAnyDepth)
	if img.Empty() {
		log.Fatalf("Error reading image from: %v\n", *source)
	}