package main

import (
	"time"

	"microservices/lib/tracer"
	"microservices/srv/hello/handler"

	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"

	example "microservices/srv/hello/proto/example"
)

func main() {
	// New Service
	t, io, err := tracer.NewTracer("go.micro.srv.hello", "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := grpc.NewService(
		micro.Name("go.micro.srv.hello"),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(t)),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		// micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// // Register Struct as Subscriber
	// micro.RegisterSubscriber("go.micro.srv.hello", service.Server(), new(subscriber.Example))

	// // Register Function as Subscriber
	// micro.RegisterSubscriber("go.micro.srv.hello", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
