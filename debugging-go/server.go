package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func serve() {
	ln, err := net.Listen("tcp", addr())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on http://%s", ln.Addr())

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		actBusy(ctx, 100*time.Millisecond) // act busy

		// Calculate atoi
		path := strings.TrimPrefix(r.URL.Path, "/")
		result := calculate(ctx, path)
		trace.SpanFromContext(ctx).AddEvent("calculated", trace.WithAttributes(attribute.Int("result", result)))
		w.Write([]byte(fmt.Sprintf("%d", result)))

		actBusy(ctx, 42*time.Millisecond) // act busy
	})
	handler = otelhttp.NewHandler(handler, "http", otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string { return r.URL.Path }))
	err = http.Serve(ln, handler)

	if err != nil {
		log.Fatal(err)
	}
}

func actBusy(ctx context.Context, d time.Duration) {
	_, span := otel.Tracer("myService").Start(ctx, "sleep")
	defer span.End()

	time.Sleep(d + time.Duration(float64(d/2)*(rand.Float64()-0.5))) // little bit of gauss
}

func calculate(ctx context.Context, str string) int {
	_, span := otel.Tracer("myService").Start(ctx, "calculate")
	defer span.End()

	return ft_atoi(str)
}

func addr() string {
	if port, hasPort := os.LookupEnv("PORT"); hasPort {
		return ":" + port
	}
	return ":8080"
}
