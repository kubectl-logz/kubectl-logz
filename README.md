# `kubectl logz`

The aim of this tool is to provide a middle-ground between a Enterprise log facility (like Splunk) and tailing
container logs using the clunky and limited `kubectl logs` command.

Solve these problem:

- Find problems.
- Give context, not just the error.
- Do it fast.

## Log Structure

There are hundreds of different log formats. The parser must normalize log messages into a semantic format.

| field      | description                          |
| ---------- | ------------------------------------ |
| `time`     | RFC339 timestamp.                    |
| `level`    | `error`, `warn`, `info`, or `debug`. |
| `msg`      | Human readable text.                 |
| `threadId` | The thread or Coroutine ID.          |

This information can be contextualized, to find the exact thread that executed the code

| field  | description       |
| ------ | ----------------- |
| `host` | The process host. |
| `pid`  | The PID.          |

`(host,pid,threadId)` gives a GUID for work being done.

Requests may traverse processes,
recording [OpenTelemetry](https://github.com/opentracing/specification/blob/master/specification.md) fields:

| field     | description |
| --------- | ----------- |
| `traceId` | Trace ID.   |
| `spanId`  | Span ID.    |

It and does so by trying each of the following parsers:

- Structured:
  - [logfmt](https://brandur.org/logfmt), e.g. `time=2022-12-04T18:36:52Z level=warn msg="thing happened"`
  - JSON, e.g. `{"time":"2022-12-04T18:36:52Z","level":"warn","msg"="thing happened"}`
- Unstructured:
  - Space-separated fields, e.g. `2022-12-04T18:36:52Z [warn] thing happened"`.
  - Failover.

## Addendum

**HTTP access logs** include the following:

| field    | description               |
| -------- | ------------------------- |
| `ip`     | IP address of the client. |
| `user`   | Client user name.         |
| `method` | HTTP method.              |
| `path`   | HTTP path.                |
| `status` | HTTP status code.         |
