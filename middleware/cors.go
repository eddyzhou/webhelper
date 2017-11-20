package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	defaultAllowOrigins = []string{"*"}
	defaultAllowHeaders = []string{
		"Mobile-Type",
		"Channel",
		"Version-Code",
		"Content-Type",
		"X-Requested-With",
		"User-Id",
		"Session-Id",
		"Peer-Id",
	}
	defaultAllowMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPatch,
	}
)

type Options struct {
	AllowOrigins     []string
	AllowCredentials bool
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           time.Duration
}

var DefaultOptions = Options{
	AllowOrigins: defaultAllowOrigins,
	AllowHeaders: defaultAllowHeaders,
	AllowMethods: defaultAllowMethods,
}

func Cors(options Options) func(next http.Handler) http.Handler {
	if options.AllowHeaders == nil {
		options.AllowHeaders = defaultAllowHeaders
	}
	if options.AllowMethods == nil {
		options.AllowMethods = defaultAllowMethods
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			method := r.Header.Get("Access-Control-Request-Method")
			headers := r.Header.Get("Access-Control-Request-Headers")

			if len(options.AllowOrigins) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", strings.Join(options.AllowOrigins, " "))
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			if options.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if len(options.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(options.ExposeHeaders, ","))
			}

			if r.Method == http.MethodOptions {
				if len(options.AllowMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowMethods, ","))
				} else if method != "" {
					w.Header().Set("Access-Control-Allow-Methods", method)
				}

				if len(options.AllowHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowHeaders, ","))
				} else if headers != "" {
					w.Header().Set("Access-Control-Allow-Headers", headers)
				}

				if options.MaxAge > time.Duration(0) {
					w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(options.MaxAge/time.Second), 10))
				}

				w.WriteHeader(http.StatusNoContent)
			} else {
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}
}
