package spotlog_test

import (
	"fmt"
	"testing"

	"github.com/13rac1/spotlog"
	"github.com/docker/distribution/context"
)

func TestTest(t *testing.T) {
	ctx := context.Background()
	ctx, logger := spotlog.Get(ctx)

	fmt.Println("debug")
	logger.Debug("debug")
	fmt.Println("info")
	logger.Info("info")
	fmt.Println("error")
	logger.Error("error")
	fmt.Println("debug")
	logger.Debug("debug")
	fmt.Println("warn")
	logger.Warn("warn")
	fmt.Println("error")
	logger.Error("error")
	fmt.Println("warn")
	logger.Warn("warn")
}
