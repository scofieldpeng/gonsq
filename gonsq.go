package gonsq

import "github.com/vaughan0/go-ini"

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
