package trace

import (
	"net/http"

	"github.com/zxfonline/iptable"
	_ "github.com/zxfonline/pprof"
	"golang.org/x/net/trace"
)

func Init(enableTracing bool) {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		w := iptable.IsTrustedIP(iptable.GetRemoteAddrIP(req.RemoteAddr))
		return w, w
	}
	EnableTracing = enableTracing

	//	if env, ok := os.LookupEnv("proj_env"); ok {
	//		switch env {
	//		case "development":
	//		case "production":
	//		}
	//	}
}

// EnableTracing controls whether to trace using the golang.org/x/net/trace package.
var EnableTracing = true
