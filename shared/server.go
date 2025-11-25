package shared

import (
	"log"
	"net/http"
	"time"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware func(http.Handler) http.Handler
}

// statusWriter captures HTTP status codes for logging
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// global request logger middleware
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(ww, r)

		// Prefer X-Forwarded-For if present (behind proxies), else RemoteAddr
		remote := r.Header.Get("X-Forwarded-For")
		if remote == "" {
			remote = r.RemoteAddr
		}

		log.Printf("remote=%s method=%s path=%s status=%d duration=%s",
			remote, r.Method, r.URL.Path, ww.status, time.Since(start))
	})
}

// method guard for ServeMux (matches path only)
func methodGuard(method string, next http.Handler) http.Handler {
	if method == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RegisterRoute(mux *http.ServeMux, routes []Route) {

	for _, rt := range routes {
		var h http.Handler = rt.Handler
		// apply optional per-route middleware
		if rt.Middleware != nil {

			h = rt.Middleware(h)

		}
		// guard HTTP method
		h = methodGuard(rt.Method, h)
		// mount on path
		mux.Handle(rt.Path, h)

	}
}

func StartServer(routes []Route, addr string) {
	mux := http.NewServeMux()

	RegisterRoute(mux, routes)

	handler := requestLogger(mux)

	log.Println("Listening on :", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
