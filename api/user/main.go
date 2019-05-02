package main

import (
	"microservice/api/user/handler"
	"microservice/lib/token"
	"microservice/lib/tracer"
	"microservice/lib/wrapper/tracer/opentracing/gin2micro"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	hystrixplugin "github.com/micro/go-plugins/wrapper/breaker/hystrix"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	web "github.com/micro/go-web"
	opentracing "github.com/opentracing/opentracing-go"
)

const name = "go.micro.api.user"

func main() {
	t, io, err := tracer.NewTracer(name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := web.NewService(
		web.Name(name),
		web.Version("lastest"),
	)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	sClient := service.Options().Service.Client()
	sClient.Init(
		client.Wrap(hystrixplugin.NewClientWrapper()),
		client.WrapCall(ocplugin.NewCallWrapper(t)),
	)

	token := &token.Token{}
	token.InitConfig("127.0.0.1:8500", "micro", "config", "jwt-key", "key")
	apiService := handler.New(sClient, token)
	router := gin.Default()
	r := router.Group("/user")
	r.Use(gin2micro.TracerWrapper)
	r.GET("/login", apiService.Anything)
	r.POST("/register", apiService.Create)

	service.Handle("/", router)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
