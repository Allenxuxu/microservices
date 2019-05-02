package handler

import (
	"errors"
	"microservice/lib/token"
	"microservice/lib/wrapper/tracer/opentracing/gin2micro"
	"net/http"

	// "time"

	helloS "microservice/srv/hello/proto/example"
	userS "microservice/srv/user/proto/user"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
)

// UserAPIService 服务
type UserAPIService struct {
	jwt    *token.Token
	helloC helloS.ExampleService
	userC  userS.UserService
}

// New UserAPIService
func New(client client.Client, token *token.Token) *UserAPIService {
	return &UserAPIService{
		jwt:    token,
		helloC: helloS.NewExampleService("", client),
		userC:  userS.NewUserService("", client),
	}
}

// Anything 测试demo，调用hello服务和user两个服务
func (s *UserAPIService) Anything(c *gin.Context) {
	log.Log("Received Say.Anything API request")

	ctx, ok := gin2micro.ContextWithSpan(c)
	if ok == false {
		log.Log("get context err")
	}

	res, err := s.helloC.Call(ctx, &helloS.Request{Name: "xuxu"})
	if err != nil {
		log.Log(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Log(res)

	userres, err := s.userC.Ping(ctx, &userS.Request{})
	if err != nil {
		log.Log(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Log(userres)

	c.JSON(http.StatusOK, map[string]string{
		"message": "Hi, this is the Greeter API",
	})
}

// Create 新建一个用户
// {
// 	"name":"徐旭",
// 	"email": "123.@qq.com",
// 	"tel":"tel1",
// 	"password":"d"
// }
func (s *UserAPIService) Create(c *gin.Context) {

	ctx, ok := gin2micro.ContextWithSpan(c)
	if ok == false {
		log.Log("get context err")
	}
	var user userS.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("JWT decode failed"))
		return
	}

	_, err := s.userC.Create(ctx, &user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
