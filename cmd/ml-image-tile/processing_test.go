
package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
	"testing"

	log "github.com/go-kit/kit/log"
)

func Test_processImageBimg(t *testing.T) {
	testFile := "../../testdata/A/testimg.png"

	type args struct {
		logger log.Logger
		path   string
		srcDir string
		resize int
		width  int
		height int