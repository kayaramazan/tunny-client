// Package common provides shared utilities for logging and protocol handling.
package common

import "go.uber.org/zap"

// NewLogger creates a new zap logger with either development or production settings.
func NewLogger(dev bool) *zap.Logger {
	if dev {
		l, _ := zap.NewDevelopment()
		return l
	}
	l, _ := zap.NewProduction()
	return l
}
