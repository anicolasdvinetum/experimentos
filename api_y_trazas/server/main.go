package main

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func setupOTelSDK(ctx context.Context) (func(context.Context) error, error) {
	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)

	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.Shutdown, nil
}

var tracer = otel.Tracer("server")

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "hello")
	defer span.End()

	w.Write([]byte("hello world"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// responder 200
	w.WriteHeader(http.StatusOK)

}

func main() {
	ctx := context.Background()
	shutdown, err := setupOTelSDK(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	http.Handle("/hello",
		otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello"),
	)

	http.Handle("/health",
		otelhttp.NewHandler(http.HandlerFunc(healthHandler), "Health"),
	)

	http.ListenAndServe(":8080", nil)
}
