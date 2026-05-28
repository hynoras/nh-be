package app

import (
	"context"
	"log/slog"
	"os"

	obs "nh-be/internal/platform/observability"
	"nh-be/pkg/env"

	"github.com/prometheus/client_golang/prometheus"
)

// SetupObservability registers Prometheus collectors for database pool metrics
// and initializes the OpenTelemetry tracer exporter. It returns a shutdown
// function that must be called during application teardown to flush any
// pending traces. The returned function is safe to pass to Service.Close().
func SetupObservability(service *Service) func(context.Context) error {
	prometheus.MustRegister(obs.NewDbPoolCollector(service.SQLDB))

	otelEndpoint := env.GetEnvOrDefault("OTEL_EXPORTER_ENDPOINT", "otel-collector:4317")
	tracerShutdown, err := obs.InitTracer(context.Background(), otelEndpoint)
	if err != nil {
		slog.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}

	return tracerShutdown
}
