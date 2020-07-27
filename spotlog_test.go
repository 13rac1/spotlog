package spotlog_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"context"

	"github.com/13rac1/spotlog"
	"github.com/sirupsen/logrus"
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

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	ctx, logger := spotlog.Get(r.Context())

	// Needed for testable example only - Sending to Stdout so the test runner
	// catches the output.
	logger.Out = os.Stdout
	// Needed for testable example only - Specify the formatter to disable the
	// timestamp for reproducible output.
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	logger.Info("request received")
	exampleCalculation(ctx, w, r)
}

func exampleCalculation(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, logger := spotlog.Get(ctx)
	logger.Error("failed calc")
	w.WriteHeader(http.StatusInternalServerError)
}

func ExampleLogger() {
	ts := httptest.NewServer(http.HandlerFunc(exampleHandler))
	defer ts.Close()

	_, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// level=info msg="request received"
	// level=error msg="failed calc"
}
