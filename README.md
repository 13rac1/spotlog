# Spot Log

**A spotlight to locate important logs**

Every log system provides log levels: set the minimum output level and throw out
the rest. The goal is to avoid wasting resources on necessary details. Those
details can be important though. Have you every thought: _"I wish I had the
debug logs for this error!"_

Spot Log stores logs below the level, then outputs all logs when a log above the
level is received.

Spot Log wraps Logrus implementing 95% the same interface.

## Usage

The `SpotLogger` instance is stored within the application or request `context`.
Get it from the `context` and use like any other logger.

```go
func passingCalc(ctx context.Context, w http.ResponseWriter) {
	_, logger := spotlog.Get(ctx)
	logger.Error("passed calc")
	w.WriteHeader(http.StatusOK)
}

func failingCalc(ctx context.Context, w http.ResponseWriter) {
	_, logger := spotlog.Get(ctx)
	logger.Error("failed calc")
	w.WriteHeader(http.StatusInternalServerError)
}

func main() {
	passingHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx, logger := spotlog.Get(req.Context())
		logger.Debug("request received")
		passingCalc(ctx, w)
		// Output: nil
	}
	http.HandleFunc("/pass", passingHandler)

	failingHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx, logger := spotlog.Get(req.Context())
		logger.Debug("request received")
		failingCalc(ctx, w)
		// Output:
		// level=debug msg="request received"
		// level=error msg="failed calc"
	}
	http.HandleFunc("/fail", failingHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Ideas

* Logger data fields as global fields. Compare to the existing Entry fields
  which are attached to an Entry.
* Print an Entry, but not all stored entries. Probably best at the `Info` level.
