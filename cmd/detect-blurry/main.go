
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
