package graceful

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var goRoutineCounter sync.WaitGroup

var stopChan chan struct{}

func init() {
	stopChan = make(chan struct{})
}

func Go(f func()) {
	goRoutineCounter.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				hostname, _ := os.Hostname()
				content := hostname + " graceful Go 程序错误: " + fmt.Sprintf("%s", r)
				log.Println(hostname, content)
			}
		}()
		defer goRoutineCounter.Done()
		f()
	}()
}

func Wait() {
	goRoutineCounter.Wait()
}

func Add(delta int) {
	goRoutineCounter.Add(delta)
}

func Done() {
	goRoutineCounter.Done()
}

func Stop() {
	close(stopChan)
}

func TimeInterval(t time.Duration, callback func()) {
	ticker := time.NewTicker(t)
	go func() {
		for {
			select {
			case <-ticker.C:
				Go(callback)
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}
