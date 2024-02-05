package OpenTelemetry

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type JaegerTracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func NewGlobalJaegerTracer(serviceName string) (*JaegerTracer, error) {
	// Initialize Jaeger
	// cfg := config.Configuration{
	// 	ServiceName: "my-tracing-app",
	// 	Sampler:     &config.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
	// 	Reporter:    &config.ReporterConfig{LogSpans: true},
	// }
	config := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := config.NewTracer()
	opentracing.SetGlobalTracer(tracer)

	if err != nil {
		return nil, err
	}

	return &JaegerTracer{
		tracer: tracer,
		closer: closer,
	}, nil
}
func (obj *JaegerTracer) StartConfig(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			// LocalAgentHostPort: jaeger_host_port,
		},
	}
	// tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	tracer, closer, err := cfg.NewTracer()

	return tracer, closer, err
}
func (obj *JaegerTracer) Close() {
	obj.closer.Close()
}
func (obj *JaegerTracer) StartGlobalSpan(operationtrachingNAme string) opentracing.Span {
	span := opentracing.StartSpan(operationtrachingNAme)
	return span
}
