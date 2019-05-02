package main

import (
	"log"
	"net"
	"net/http"

	"microservice/lib/token"
	"microservice/lib/tracer"
	"microservice/lib/wrapper/auth"
	"microservice/lib/wrapper/breaker/hystrix"
	"microservice/lib/wrapper/metrics/prometheus"
	"microservice/lib/wrapper/tracer/opentracing/stdhttp"

	ph "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-plugins/micro/cors"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
	opentracing "github.com/opentracing/opentracing-go"
)

func init() {
	token := &token.Token{}
	token.InitConfig("127.0.0.1:8500", "micro", "config", "jwt-key", "key")

	plugin.Register(cors.NewPlugin())

	plugin.Register(plugin.NewPlugin(
		plugin.WithName("auth"),
		plugin.WithHandler(
			auth.JWTAuthWrapper(token),
		),
	))
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			stdhttp.TracerWrapper,
		),
	))
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("breaker"),
		plugin.WithHandler(
			hystrix.BreakerWrapper,
		),
	))
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("metrics"),
		plugin.WithHandler(
			prometheus.MetricsWrapper,
		),
	))
}
const name = "API gateway"

func main() {
	t, io, err := tracer.NewTracer(name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	hystrixStreamHandler := ph.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)

	cmd.Init()
}
