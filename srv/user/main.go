package main

import (
	"time"

	"github.com/Allenxuxu/microservices/lib/token"
	"github.com/Allenxuxu/microservices/lib/tracer"
	"github.com/Allenxuxu/microservices/srv/user/db"
	"github.com/Allenxuxu/microservices/srv/user/handler"
	pb "github.com/Allenxuxu/microservices/srv/user/proto/user"

	"github.com/micro/cli"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

const name = "go.micro.srv.user"

func main() {
	token := &token.Token{}
	var consulAddr string

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
		micro.Flags(cli.StringFlag{
			Name:   "consul_address",
			Usage:  "consul address for K/V",
			EnvVar: "CONSUL_ADDRESS",
			Value:  "127.0.0.1:8500",
		}),
		micro.Action(func(ctx *cli.Context) {
			consulAddr = ctx.String("consul_address")
			token.InitConfig(consulAddr, "micro", "config", "jwt-key", "key")
		}),
	)
	service.Init()

	//从consul KV 获取 DB 配置
	db.Init(consulAddr)
	pb.RegisterUserServiceHandler(service.Server(), handler.New(token))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
