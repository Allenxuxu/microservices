package hystrix

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	status_code "github.com/Allenxuxu/microservices/lib/http"

	"github.com/afex/hystrix-go/hystrix"
)

// BreakerWrapper hystrix breaker
func BreakerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		log.Println(name)
		err := hystrix.Do(name, func() error {
			sct := &status_code.StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}
			h.ServeHTTP(sct.WrappedResponseWriter(), r)

			if sct.Status >= http.StatusBadRequest {
				str := fmt.Sprintf("status code %d", sct.Status)
				log.Println(str)
				return errors.New(str)
			}
			return nil
		}, nil)
		if err != nil {
			log.Println("hystrix breaker err: ", err)
			return
		}
	})
}
