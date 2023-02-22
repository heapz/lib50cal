package logger

import (
	"fmt"

	"sync/atomic"
)

type Logger interface{ Output(int, string) error }

// crazy ass user may call NewLoggerConfig multiple times, we need law and order here
var debug = atomic.Pointer[Logger]{}
var print = atomic.Pointer[Logger]{}

var isDiscardPrint int32
var isDiscardDebug int32

func Printf(format string, v ...any) {
	if atomic.LoadInt32(&isDiscardPrint) != 0 {
		return
	}

	(*print.Load()).Output(2, fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...any) {
	if atomic.LoadInt32(&isDiscardDebug) != 0 {
		return
	}

	(*debug.Load()).Output(2, fmt.Sprintf(format, v...))
}

func RegisterDebug(out Logger) {
	var isDiscard int32
	if out == nil {
		isDiscard = 1
	}
	atomic.StoreInt32(&isDiscardDebug, isDiscard)
	debug.Swap(&out)
}

func RegisterPrint(out Logger) {
	var isDiscard int32
	if out == nil {
		isDiscard = 1
	}
	atomic.StoreInt32(&isDiscardPrint, isDiscard)
	print.Swap(&out)
}
