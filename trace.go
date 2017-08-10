// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"net/http"

	_ "github.com/zxfonline/expvar"
	"github.com/zxfonline/iptable"
	_ "github.com/zxfonline/pprof"
	"golang.org/x/net/trace"
)

func Init(enableTracing bool, checkip bool) {
	iptable.CHECK_IPTRUSTED = checkip
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		w := iptable.IsTrustedIP(iptable.GetRemoteAddrIP(req.RemoteAddr), false)
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

//ProxyTrace 跟踪
type ProxyTrace struct {
	tr trace.Trace
}

//TraceStart 开始跟踪
func TraceStart(family, title string) *ProxyTrace {
	if EnableTracing {
		pt := &ProxyTrace{tr: trace.New(family, title)}
		return pt
	}
	return nil
}

func TraceFinish(pt *ProxyTrace) {
	if pt != nil {
		if pt.tr != nil {
			pt.tr.Finish()
		}
	}
}

func TracePrintf(pt *ProxyTrace, format string, a ...interface{}) {
	if pt != nil {
		if pt.tr != nil {
			pt.tr.LazyPrintf(format, a...)
		}
	}
}

func TraceErrorf(pt *ProxyTrace, format string, a ...interface{}) {
	if pt != nil {
		if pt.tr != nil {
			pt.tr.LazyPrintf(format, a...)
			pt.tr.SetError()
		}
	}
}
