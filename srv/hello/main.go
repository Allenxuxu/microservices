package main

import (
	"context"
	"time"

	"github.com/Allenxuxu/microservices/lib/tracer"
	"github.com/Allenxuxu/microservices/srv/hello/handler"
	example "github.com/Allenxuxu/microservices/srv/hello/proto/example"
	"github.com/Allenxuxu/microservices/srv/hello/subscriber"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

func Handler(ctx context.Context, msg *example.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}

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

	// Register Function as Subscriber
	micro.RegisterSubscriber("/test", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
