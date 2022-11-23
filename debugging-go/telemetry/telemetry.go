package telemetry

import (
	"log"
	"os"
	"strings"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"

	"cloud.google.com/go/profiler"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var ServiceName = "service"

func init() {
	if strings.Contains(os.Args[0], "test") {
		return
	}

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
