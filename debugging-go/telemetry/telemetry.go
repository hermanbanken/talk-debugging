package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.uber.org/zap"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/profiler"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	ttrace "go.opentelemetry.io/otel/trace"
)

var ServiceName = "service"

func init() {
	if strings.Contains(os.Args[0], "test") {
		return
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, _ := config.Build()
	zap.ReplaceGlobals(l)

	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(profiler.Config{
		ProjectID:      os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Service:        os.Getenv("OTEL_SERVICE_NAME"),
		ServiceVersion: "1.0.0",
	}); err != nil {
		log.Fatal(err)
	}

	// Set the network propagator format to W3C (https://www.w3.org/TR/trace-context/)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))

	exporterTraces, err := texporter.New(texporterOpts()...)
	if err != nil {
		log.Fatal("Failed to create google telemetry exporter", err)
	}

	// Trace provider
	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporterTraces)),
		trace.WithResource(resource.Environment()),
	))
}

func texporterOpts() (out []texporter.Option) {
	if proj, hasProject := os.LookupEnv("GOOGLE_CLOUD_PROJECT"); hasProject {
		out = append(out, texporter.WithProjectID(proj))
	}
	return
}

// WrapZap injects the current OpenTelemetry trace into the log statement, to link it between Cloud Trace and Cloud Logging.
func WrapZap(ctx context.Context, l *zap.Logger) *zap.Logger {
	if s := ttrace.SpanFromContext(ctx).SpanContext(); s.IsValid() {
		if l == nil {
			l = zap.L()
		}

		// Ignore errors getting the project ID, the app would have crashed before this code gets reached
		// It's also fine to have a slightly invalid trace ID metadata string than logging nothing...
		projectID, _ := GetProjectID()

		// Integrate Trace and Logging: https://cloud.google.com/trace/docs/trace-log-integration
		// Format: projects/[PROJECT_ID]/traces/[TRACE_ID]
		// Also: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
		// Special Fields docs: https://cloud.google.com/logging/docs/agent/configuration#special-fields
		l = l.With(
			zap.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", projectID, s.TraceID().String())),
			zap.String("logging.googleapis.com/spanId", s.SpanID().String()),
			zap.Bool("logging.googleapis.com/trace_sampled", s.IsSampled()), // required for Cloud Logs UI to allow clicking through to Cloud Trace
		)
	}
	return l
}

func GetProjectID() (string, error) {
	proj, err := metadata.ProjectID()
	if err != nil {
		return os.Getenv("GOOGLE_CLOUD_PROJECT"), nil
	}
	return proj, nil
}
