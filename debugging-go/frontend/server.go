package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/gob"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	jobt "dummy/calculator/job"
	"dummy/telemetry"
	"dummy/util"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	telemetry.ServiceName = "frontend"
}

//go:embed index.html
var indexTemplateStr string
var indexTemplate = template.Must(template.New("").Parse(indexTemplateStr))

type indexData struct {
	Equation string
	Result   *int
}

func serve() {
	ln, err := net.Listen("tcp", util.Addr())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on http://%s", ln.Addr())

	var handler http.Handler = http.HandlerFunc(handler)
	handler = otelhttp.NewHandler(handler, "http", otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string { return r.URL.Path }))
	err = http.Serve(ln, handler)

	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	util.DoBusyWork(10000) // act busy

	if r.Method == http.MethodPost {
		trace.SpanFromContext(ctx).SetName("submit form")
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := solveEquation(ctx, r.FormValue("equation"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		indexTemplate.Execute(w, indexData{
			Equation: r.FormValue("equation"),
			Result:   &result,
		})
		return
	}

	trace.SpanFromContext(ctx).SetName("render form")
	indexTemplate.Execute(w, indexData{
		Equation: "",
		Result:   nil,
	})
}

func solveEquation(ctx context.Context, equation string) (res int, err error) {
	_, span := otel.Tracer("frontend").Start(ctx, "solve")
	defer span.End()
	defer func() {
		if err == nil {
			span.SetStatus(codes.Ok, "done")
		} else {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	var job jobt.Job
	resp, err := otelhttp.Post(ctx, os.Getenv("URL_EQUATIONS"), "text/plain", strings.NewReader(equation))
	if err != nil {
		return 0, errors.Wrap(err, "equations")
	}
	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return 0, errors.New(string(data))
	}

	err = gob.NewDecoder(resp.Body).Decode(&job)
	if err != nil {
		return 0, errors.Wrap(err, "gob")
	}

	return solveRecursive(ctx, job)
}

func solveRecursive(ctx context.Context, job jobt.Job) (res int, err error) {
	var a, b int

	switch v := job.A.(type) {
	case jobt.Value:
		a = int(v)
	case jobt.Job:
		a, err = solveRecursive(ctx, v)
		if err != nil {
			return
		}
	}

	switch v := job.B.(type) {
	case jobt.Value:
		b = int(v)
	case jobt.Job:
		b, err = solveRecursive(ctx, v)
		if err != nil {
			return
		}
	}

	buf := bytes.NewBuffer(nil)
	err = errors.Wrap(gob.NewEncoder(buf).Encode(jobt.Job{
		Op: job.Op,
		A:  jobt.Value(a),
		B:  jobt.Value(b),
	}), "gob encode")
	if err != nil {
		return
	}

	resp, err := otelhttp.Post(ctx, os.Getenv("URL_CALCULATOR"), "text/plain", buf)
	if err != nil {
		return 0, errors.Wrap(err, "calculator")
	}
	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return 0, errors.New(string(data))
	}
	err = gob.NewDecoder(resp.Body).Decode(&res)
	return
}
