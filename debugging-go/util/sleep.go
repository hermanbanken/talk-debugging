package util

import (
	"context"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
)

func ActBusy(ctx context.Context, d time.Duration) {
	_, span := otel.Tracer("sleeper").Start(ctx, "sleep")
	defer span.End()

	time.Sleep(d + time.Duration(float64(d/2)*(rand.Float64()-0.5))) // little bit of gauss
}
