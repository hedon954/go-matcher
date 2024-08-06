package safe

import (
	"runtime/debug"
	"sync"
)

var sg = safeGo{
	goCallbacks:   make([]PanicCallback, 0),
	callCallbacks: make([]PanicCallback, 0),
}

type safeGo struct {
	sync.WaitGroup
	goCallbacks   []PanicCallback
	callCallbacks []PanicCallback
}

type PanicCallback = func(err any, stack []byte)

func Callback(callbacks ...PanicCallback) {
	sg.goCallbacks = append(sg.goCallbacks, callbacks...)
	sg.callCallbacks = append(sg.callCallbacks, callbacks...)
}

func GoCallBack(callbacks ...PanicCallback) {
	sg.goCallbacks = append(sg.goCallbacks, callbacks...)
}

func CallCallBack(callbacks ...PanicCallback) {
	sg.callCallbacks = append(sg.callCallbacks, callbacks...)
}

func Call(f func(), panicCallback ...PanicCallback) {
	defer func() {
		debug.Stack()
		if err := recover(); err != nil {
			stack := debug.Stack()
			for _, callback := range sg.callCallbacks {
				callback(err, stack)
			}
			for _, callback := range panicCallback {
				callback(err, stack)
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
				stack := debug.Stack()
				for _, callback := range sg.goCallbacks {
					callback(err, stack)
				}
				for _, callback := range panicCallback {
					callback(err, stack)
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
