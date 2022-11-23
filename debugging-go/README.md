# OpenTelemetry demo
https://opentelemetry.io/ is cross-platform (language, OS, etc.) cross-business (AWS, Google, DataDog, ZipKin, Splunk, Jaeger, etc.) tracing, logging & metric format. Each application exports telemetry to a common sink, which can then visualize traces across different applications.

Example configuration in `telemetry.go`, which picks up [some standard environment](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/sdk-environment-variables.md) variables:

```bash
OTEL_SERVICE_NAME=serviceA
OTEL_RESOURCE_ATTRIBUTES=g.co/gae/app/module=serviceA # Cloud Trace displays 'g.co/gae/app/module' in the UI as "Service"; there is no alternative key yet
OTEL_TRACES_SAMPLER=parentbased_always_on
```
