package opentelemetry

// TODO: 8-02

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	// 准备1
	res, err := newResource("demo", "v0.0.1")
	require.NoError(t, err)

	// 准备2
	porp := newPropagator()
	otel.SetTextMapPropagator(porp)

	// 准备3
	tp, err := newTranceProvider(res)
	require.NoError(t, err)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	server := gin.Default()
	server.GET("/test", func(ginCtx *gin.Context) {
		// 这个 Tracer 的名字，最好设置为唯一的，比如说用所在包名
		tracer := otel.Tracer("opentelemtry")

		var ctx context.Context = ginCtx
		ctx, span := tracer.Start(ctx, "top-span")
		defer span.End()

		span.AddEvent("event-1")
		time.Sleep(time.Second)
		ctx, subSpan := tracer.Start(ctx, "sub-span")
		defer subSpan.End()
		time.Sleep(time.Millisecond * 300)
		subSpan.SetAttributes(attribute.String("key1", "value1"))
		ginCtx.String(http.StatusOK, "ok")
	})
	server.Run(":8082")

}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		)
}

func newTranceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}