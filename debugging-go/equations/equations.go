package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"math/rand"
	"net/http"
	"time"

	"dummy/atoi"
	"dummy/calculator/job"
	"dummy/telemetry"
	"dummy/util"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func init() {
	telemetry.ServiceName = "equations"
}

func main() {
	var h http.Handler = http.HandlerFunc(handler)
	h = otelhttp.NewHandler(h, "equations")
	http.ListenAndServe(util.Addr(), h)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}

		result, err := parseEquation(r.Context(), data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = gob.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func parseEquation(ctx context.Context, data []byte) (res job.Expr, err error) {
	ctx, span := otel.Tracer("equations").Start(ctx, "parseEquation")
	span.SetAttributes(attribute.String("eq", string(data)))
	defer span.End()
	defer func() {
		if err == nil {
			span.SetStatus(codes.Ok, "done")
		} else {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	telemetry.WrapZap(ctx, zap.L()).Sugar().Infof("parsing %q", string(data))

	// emulate random 0..50ms work
	time.Sleep(time.Duration(rand.Float32() * float32(50*time.Millisecond)))

	// TODO use the right algorithm (https://en.wikipedia.org/wiki/Shunting_yard_algorithm)

	// +/- lowest precedence
	if idx := bytes.IndexFunc(data, func(r rune) bool { return r == '+' || r == '-' }); idx >= 0 {
		res1, err := parseEquation(ctx, data[0:idx])
		if err != nil {
			return job.Job{}, err
		}
		res2, err := parseEquation(ctx, data[idx+1:])
		if err != nil {
			return job.Job{}, err
		}
		if data[idx] == '+' {
			return job.Job{A: res1, B: res2, Op: job.Add}, nil
		}
		return job.Job{A: res1, B: res2, Op: job.Remove}, nil
	}

	// * and / higher precedence
	if idx := bytes.IndexFunc(data, func(r rune) bool { return r == '*' || r == '/' }); idx >= 0 {
		res1, err := parseEquation(ctx, data[0:idx])
		if err != nil {
			return job.Job{}, err
		}
		res2, err := parseEquation(ctx, data[idx+1:])
		if err != nil {
			return job.Job{}, err
		}
		if data[idx] == '*' {
			return job.Job{A: res1, B: res2, Op: job.Multiply}, nil
		}
		return job.Job{A: res1, B: res2, Op: job.Divide}, nil
	}

	d := atoi.Atoi(string(data))
	return job.Value(d), nil
}
