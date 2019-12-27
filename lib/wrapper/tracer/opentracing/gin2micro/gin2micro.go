package gin2micro

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-log/log"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const contextTracerKey = "Tracer-context"

// sf sampling frequency
var sf = 100

func init() {
	rand.Seed(time.Now().Unix())
}

// SetSamplingFrequency 设置采样频率
// 0 <= n <= 100
func SetSamplingFrequency(n int) {
	sf = n
}

// TracerWrapper tracer 中间件
func TracerWrapper(c *gin.Context) {
	sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
	tracer := opentracing.GlobalTracer()
	md := make(map[string]string)
	nsf := sf
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	if err == nil {
		sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
		tracer = sp.Tracer()
		nsf = 100
	}
	defer sp.Finish()

	if err := tracer.Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(md)); err != nil {
		log.Log(err)
	}

	ctx := context.TODO()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	c.Set(contextTracerKey, ctx)

	c.Next()

	statusCode := c.Writer.Status()
	ext.HTTPStatusCode.Set(sp, uint16(statusCode))
	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	if statusCode >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	} else if rand.Intn(100) > nsf {
		ext.SamplingPriority.Set(sp, 0)
	}
}

// ContextWithSpan 返回context
func ContextWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(contextTracerKey)
	if exist == false {
		ok = false
		ctx = context.TODO()
		return
	}

	ctx, ok = v.(context.Context)
	return
}
