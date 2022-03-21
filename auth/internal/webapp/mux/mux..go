// Package mux manages the API handlers.
package mux

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/mroobert/go-tickets/auth/internal/foundation/web"
	"github.com/mroobert/go-tickets/auth/internal/webapp/mid"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	SignUpHandler web.Handler
	SignInHandler web.Handler
	Log           *zap.SugaredLogger
	Shutdown      chan os.Signal
}

// APIMux constructs a mux with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.AppMux {
	mux := web.NewAppMux(cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Panics(),
	)

	const group = "api"
	mux.Handle(http.MethodPost, group, "/signup", cfg.SignUpHandler)
	mux.Handle(http.MethodPost, group, "/signin", cfg.SignInHandler)

	return mux
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
