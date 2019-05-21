package main

import (
	"time"

	"github.com/Allenxuxu/microservices/lib/token"
	"github.com/Allenxuxu/microservices/lib/tracer"
	"github.com/Allenxuxu/microservices/srv/user/db"
	"github.com/Allenxuxu/microservices/srv/user/handler"
	pb "github.com/Allenxuxu/microservices/srv/user/proto/user"

	"github.com/micro/go-grpc"
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

	service := grpc.NewService(
		micro.Name(name),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	service.Init()

	//从consul KV 获取 DB 配置
	db.Init("127.0.0.1:8500")
	token := &token.Token{}
	token.InitConfig("127.0.0.1:8500", "micro", "config", "jwt-key", "key")
	pb.RegisterUserServiceHandler(service.Server(), handler.New(token))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
