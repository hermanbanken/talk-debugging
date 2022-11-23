package main

import (
	"context"
	"dummy/util"
	"encoding/gob"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"dummy/calculator/job"
	"dummy/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func init() {
	telemetry.ServiceName = "calculator"
}

func main() {
	var h http.Handler = http.HandlerFunc(handler)
	h = otelhttp.NewHandler(h, "calculator")
	http.ListenAndServe(util.Addr(), h)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		var job job.Job
		err := gob.NewDecoder(r.Body).Decode(&job)
		if err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}
		result, err := calc(r.Context(), job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = gob.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func calc(ctx context.Context, j job.Job) (res int, err error) {
	_, span := otel.Tracer("calc").Start(ctx, "calc")
	span.SetName("calculate " + j.String())
	defer span.End()
	defer func() {
		if err == nil {
			span.SetStatus(codes.Ok, "done")
		} else {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	// emulate random 0..50ms work
	time.Sleep(time.Duration(rand.Float32() * float32(50*time.Millisecond)))

	a, aIsValue := j.A.(job.Value)
	b, bIsValue := j.B.(job.Value)
	if !aIsValue || !bIsValue {
		return 0, errors.New("provide calculable jobs, not nested expressions")
	}

	switch j.Op {
	case job.Add:
		return int(a) + int(b), nil
	case job.Remove:
		return int(a) - int(b), nil
	case job.Multiply:
		return int(a) * int(b), nil
	case job.Divide:
		if int(b) == 0 {
			return 0, errors.New("division by zero")
		}
		return int(a) / int(b), nil
	}
	return 0, errors.New("unknown operation")
}
