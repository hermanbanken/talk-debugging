package main

import (
	"log"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var ServiceName = "service"

func init() {
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
