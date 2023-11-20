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
	otel.SetTextMapPropagator(prop)

	// 准备3
	tp, err := newTranceProvider(res)
	require.NoError(t, err)
	defer tp.Shutdown(context.Background())
	otel.SetTraceProcider(tp)

	server := gin.Default()
	server.GET("test", func(ginCtx *gin.Context) {

	})
	server.Run(":8082")

}

func newResource(serviceName, serviceVersion string) (*resource.Resouce, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			))
}

newPropagator

newTranceProvider