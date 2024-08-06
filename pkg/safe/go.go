package safe

import (
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
)

var sg = safeGo{
	logger: slog.Default(),
}

type safeGo struct {
	sync.WaitGroup
	logger ErrorLogger
}

type PanicCallback = func(err any)

type ErrorLogger interface {
	Error(msg string, v ...any)
}

func SetLogger(logger ErrorLogger) {
	sg.logger = logger
}

func Call(f func(), panicCallback ...PanicCallback) {
	defer func() {
		debug.Stack()
		if err := recover(); err != nil {
			sg.logger.Error(fmt.Sprintf("safe call occurs panic: %v \n", err), string(debug.Stack()))
			for _, callback := range panicCallback {
				callback(err)
			}
		}
	}()
	f()
}

func Go(f func(), panicCallback ...PanicCallback) {
	sg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				sg.logger.Error(fmt.Sprintf("safe go occurs panic: %v \n", err), string(debug.Stack()))
				for _, callback := range panicCallback {
					callback(err)
				}
			}
		}()
		defer sg.Done()
		f()
	}()
}

func Wait() {
	sg.Wait()
}
