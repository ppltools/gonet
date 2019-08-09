package recovery

import (
	"context"
	"runtime/debug"
	"strconv"
)

// Tracer 定义日志接口，任何实现该接口的日志对象都可用
type Tracer interface {
	Fatalf(ctx context.Context, format string, args ...interface{})
}

// Recovery 函数是在 built-in 函数 recover 之上，将对 error 的日志做的一层封装。
// 参数:
//	     ctx: 请求上下文
//	     log: 实现Tracer接口的任何实例，用于输出日志，若为 nil 则不输出日志
//    module: 子模块名，每个 goroutine 定制一个子模块名
// 方法:
//   defer recovery.Recovery(ctx, logger, name)
func Recovery(ctx context.Context, log Tracer, module string) {
	if err := recover(); err != nil {
		if log != nil {
			log.Fatalf(ctx, "panic in %s goroutine||errmsg=%s||stack info=%s", module, err, strconv.Quote(string(debug.Stack())))
		}
	}
}
