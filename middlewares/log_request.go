package middlewares

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func MiddlewareLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Trace(r)
			log.Infoln(r.Method, r.RemoteAddr, r.URL)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		},
	)
}
