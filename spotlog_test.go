package spotlog_test

import (
	"fmt"
	"testing"

	"context"

	"github.com/13rac1/spotlog"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	ctx, logger := spotlog.Get(ctx)

	// The logrus.FieldLogger interface is implemented near exactly, but changed
	// to return spotlog.Entry.
	// var logruslogger logrus.FieldLogger = logger

	fmt.Println("debug")
	logger.Debug("debug")
	fmt.Println("info")
	logger.Info("info")
	fmt.Println("error")
	logger.Error("error")

	entry := logger.WithField("test", "value")
	fmt.Println("debug")
	entry.Debug("debug")
	fmt.Println("info")
	entry.Info("info")
	entry = entry.WithField("test3", "value3")
	fmt.Println("error")
	entry.Error("error")
}
