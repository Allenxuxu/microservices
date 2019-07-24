package main

import (
	"context"
	"time"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/Allenxuxu/microservices/api/user/handler"
	"github.com/Allenxuxu/microservices/lib/token"
	"github.com/Allenxuxu/microservices/lib/tracer"
	"github.com/Allenxuxu/microservices/lib/wrapper/tracer/opentracing/gin2micro"

	"github.com/gin-gonic/gin"
	"github.com/micro/cli"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	hystrixplugin "github.com/micro/go-plugins/wrapper/breaker/hystrix"
	web "github.com/micro/go-micro/web"
	opentracing "github.com/opentracing/opentracing-go"
)

const name = "go.micro.api.user"

func main() {
	token := &token.Token{}

	gin2micro.SetSamplingFrequency(50)
	t, io, err := tracer.NewTracer(name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := web.NewService(
		web.Name(name),
		web.Version("lastest"),
		web.RegisterTTL(time.Second*15),
		web.RegisterInterval(time.Second*10),
		web.MicroService(grpc.NewService()),
		web.Flags(cli.StringFlag{
			Name:   "consul_address",
			Usage:  "consul address for K/V",
			EnvVar: "CONSUL_ADDRESS",
			Value:  "127.0.0.1:8500",
		}),
		web.Action(func(ctx *cli.Context) {
			token.InitConfig(ctx.String("consul_address"), "micro", "config", "jwt-key", "key")
		}),
	)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	hystrix.DefaultTimeout = 5000

	sClient := hystrixplugin.NewClientWrapper()(service.Options().Service.Client())
	sClient.Init(
		// client.WrapCall(ocplugin.NewCallWrapper(t)),
		client.Retries(3),
		client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
			log.Log(req.Method(), retryCount, " client retry")
			return true, nil
		}),
	)

	pub := micro.NewPublisher("/test", sClient)

	apiService := handler.New(sClient, pub, token)
	router := gin.Default()
	r := router.Group("/user")
	r.Use(gin2micro.TracerWrapper)
	r.GET("/test", apiService.Anything)
	r.POST("/register", apiService.Create)

	service.Handle("/", router)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
