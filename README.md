# Spot Log

**A spotlight to locate important logs**

Every log system provides log levels: set the minimum output level and throw out
the rest. The goal is to avoid wasting resources on necessary details. Those
details can be important though. Have you every thought: *"I wish I had the
debug logs for this error!"*

Spot Log stores logs below the level, then outputs all logs when a log above the
level is received.

## Warning

`Logger` instances maintain the list of stored log entries and `Entry`
instances use their parent `Logger` instance's list. Therefore
there must be a separate `Logger` instance for each HTTP Request or goroutine.

TODO: Optionally store entries in the Context.
