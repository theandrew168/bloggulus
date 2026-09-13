package otel

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/pgx-contrib/pgxotel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	traceAPI "go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/theandrew168/bloggulus"

func NewResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("bloggulus"),
		semconv.ServiceVersion("v0.8.0"),
	)
}

func NewLogger(ctx context.Context, res *resource.Resource) (*log.LoggerProvider, error) {
	victoriaLogsExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpointURL("http://localhost:9428/insert/opentelemetry/v1/logs"),
	)
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(log.NewBatchProcessor(victoriaLogsExporter)),
	)
	global.SetLoggerProvider(loggerProvider)

	otelSlogHandler := otelslog.NewHandler(
		"bloggulus",
		otelslog.WithLoggerProvider(loggerProvider),
	)
	slog.SetDefault(slog.New(otelSlogHandler))

	return loggerProvider, nil
}

func NewTracer(ctx context.Context, res *resource.Resource) (*trace.TracerProvider, error) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	victoriaTracesExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL("http://localhost:10428/insert/opentelemetry/v1/traces"),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(victoriaTracesExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}

func GetTracer() traceAPI.Tracer {
	return otel.Tracer(tracerName)
}

func NewQueryTracer(tracerProvider *trace.TracerProvider) pgx.QueryTracer {
	pgxTracer := pgxotel.QueryTracer{
		Name:     tracerName,
		Provider: tracerProvider,
	}
	return &pgxTracer
}
