package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	kafka "github.com/segmentio/kafka-go"
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

	fmt.Println("server arrancó; broker=", os.Getenv("KAFKA_BROKER"))

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_BROKER")},
		Topic:     "pruebas-k6",
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})

	http.Handle("/health",
		otelhttp.NewHandler(http.HandlerFunc(healthHandler), "Health"),
	)

	go http.ListenAndServe(":8080", nil)

	for {
		fmt.Println("intentando leer")

		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			fmt.Println("ERROR FETCH:", err)
			continue
		}
		fmt.Println("RECIBIDO:", string(msg.Value))

		_, span := tracer.Start(ctx, "consume")

		span.End()

	}

}
