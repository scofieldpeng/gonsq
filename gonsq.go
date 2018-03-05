package gonsq

import (
	"sync"

	"os"

	"os/signal"
	"syscall"

	"github.com/vaughan0/go-ini"
)

// InitAll 初始化producer 和 consumer
func InitAll(producerConfig, consumerConfig ini.Section, debug bool) (err error) {
	if err = Producer.Init(producerConfig, debug); err != nil {
		return
	}
	if err = Consumer.Init(consumerConfig, debug); err != nil {
		return
	}

	return
}

// RunAll 运行 nsq 服务
func RunAll() (err error) {
	if err = Producer.Run(); err != nil {
		return
	}
	err = Consumer.Run()
	return
}

// StopAll 停止nsq 服务
func StopAll() {
	Producer.Stop()
	Consumer.Stop()
}

// 优雅运行，相当于执行:
// RunAll()
// defer StopAll()
//
func GracefulRunAll(wg *sync.WaitGroup) (err error) {
	wg.Add(1)
	if err = RunAll(); err != nil {
		return
	}
	go func() {
		wg.Add(1)
		defer wg.Done()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
		<-stopChan
		StopAll()
	}()
}
