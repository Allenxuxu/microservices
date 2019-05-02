package main

import (
	"microservice/lib/token"
	"microservice/lib/tracer"
	_ "microservice/srv/user/db"
	"microservice/srv/user/handler"
	pb "microservice/srv/user/proto/user"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

const name = "go.micro.srv.user"

func main() {
	t, io, err := tracer.NewTracer(name, "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(name),
		micro.Version("latest"),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	service.Init()

	token := &token.Token{}
	token.InitConfig("127.0.0.1:8500", "micro", "config", "jwt-key", "key")
	pb.RegisterUserServiceHandler(service.Server(), handler.New(token))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
