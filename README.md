# Spot Log

**A spotlight to locate important logs**

Every log system provides log levels: set the minimum output level and throw out
the rest. The goal is to avoid wasting resources on necessary details. Those
details can be important though. Have you every thought: *"I wish I had the
debug logs for this error!"*

Spot Log stores logs below the level, then outputs all logs when a log above the
level is received.
