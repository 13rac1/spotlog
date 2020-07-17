package spotlog_test

import (
	"bytes"
	"testing"

	"context"

	"github.com/13rac1/spotlog"
	"github.com/stretchr/testify/assert"
)

// The logrus.FieldLogger interface is implemented near exactly, but changed
// to return spotlog.Entry.
// var logruslogger logrus.FieldLogger = logger

func TestLogger(t *testing.T) {
	ctx := context.Background()
	ctx, logger := spotlog.Get(ctx)

	var stdout bytes.Buffer
	logger.Out = &stdout

	logger.Debug("debugmsg")
	assert.Empty(t, stdout.String())

	logger.Info("infomsg")
	assert.Empty(t, stdout.String())

	logger.Error("errormsg")
	assert.NotEmpty(t, stdout.String())
	assert.Contains(t, stdout.String(), "debugmsg")
	assert.Contains(t, stdout.String(), "infomsg")
	assert.Contains(t, stdout.String(), "errormsg")
}

func TestEntry(t *testing.T) {
	ctx := context.Background()
	ctx, logger := spotlog.Get(ctx)

	entry := logger.WithField("test", "value")

	var stdout bytes.Buffer
	logger.Out = &stdout

	entry.Debug("debugmsg")
	assert.Empty(t, stdout.String())

	entry.Info("infomsg")
	assert.Empty(t, stdout.String())

	entry.Error("errormsg")
	assert.NotEmpty(t, stdout.String())
	assert.Contains(t, stdout.String(), "msg=debugmsg test=value")
	assert.Contains(t, stdout.String(), "msg=infomsg test=value")
	assert.Contains(t, stdout.String(), "msg=errormsg test=value")
}
