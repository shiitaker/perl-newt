package main

import (
	"strings"

	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewLevelFilterFromString filter the log level using the string "DEBUG|INFO|WARN|ERROR"
func NewLevelFilterFromString(next log.Logger, ls string) log.Logger {
	switch strings.ToLower(ls) {
	case "debug":
		return l