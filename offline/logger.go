package offline

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Logger(inner httprouter.Handle, name string) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Printf(
			"Start to handle : %s\t%s\t%s\t",
			r.Method,
			r.RequestURI,
			name,
		)
		start := time.Now()
		//time.Sleep(time.Millisecond * 1000)

		inner(w, r, ps)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func LoggerNotFound(inner http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			"404 Not Found",
			time.Since(start),
		)
	})
}
