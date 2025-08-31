package bybit

import (
	"context"
	"log"
	"os"
	"strings"
)

var runtimeCtx context.Context

var debugEnabled bool

func init() {
	// Enable extra logging if BYBIT_DEBUG is set to 1/true
	v := strings.ToLower(os.Getenv("BYBIT_DEBUG"))
	debugEnabled = v == "1" || v == "true" || v == "yes" || v == "on"
}

func SetRuntimeCtx(ctx context.Context) {
	runtimeCtx = ctx
}

func getRuntimeCtx() context.Context {
	return runtimeCtx
}

// dbg logs only when debug is enabled to reduce noisy output in production
func dbg(format string, args ...interface{}) {
	if debugEnabled {
		log.Printf(format, args...)
	}
}
