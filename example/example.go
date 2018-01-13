package main

import (
	"os"

	"errors"
	"github.com/nsqio/go-nsq"
	"github.com/scofieldpeng/gonsq"
	"github.com/vaughan0/go-ini"
	"go.zhuzi.me/log"
)

var (
	producerConfig ini.Section
	consumerConfig ini.Section
	debug          = true
	receiveChan    = make(chan string, 10)
	receiveNum     = 0
)

func main() {
	producerConfig = make(map[string]string)
	producerConfig["nsqd"] = "192.168.2.100:4150"

	consumerConfig = make(map[string]string)
	consumerConfig["nsqlookupd"] = "192.168.2.100:4161"
	consumerConfig["maxInFlight"] = "5"
	consumerConfig["concurrent"] = "3"
	consumerConfig["channelName"] = "chan1"
	consumerConfig["max_attempt"] = "2"
	log.SetDebug(debug)

	if err := gonsq.InitAll(producerConfig, consumerConfig, true); err != nil {
		log.Panic(err)
	}
	gonsq.Consumer.AddFailHandler("test", testFailHandler())
	gonsq.Consumer.AddHandler("test", testHandler())

	if err := gonsq.RunAll(); err != nil {
		log.Panic(err)
	}
	defer gonsq.StopAll()

	go func() {
		i := 0
		for {
			if i == 5 {
				break
			}
			if err := gonsq.Producer.Publish("test", "hello world"); err != nil {
				log.Error("produce error:", err.Error())
				continue
			}
			log.Info("producer success!")
			i++
		}
	}()

	for {
		select {
		case <-receiveChan:
			receiveNum++
			log.Info("receive,num:", receiveNum)
			if receiveNum == 5 {
				os.Exit(0)
			}
		}
	}
}

func testHandler() nsq.HandlerFunc {
	return func(nm *nsq.Message) error {
		return errors.New(string(nm.Body))
	}
}

func testFailHandler() gonsq.FailMessageFunc {
	return func(message gonsq.FailMessage) (err error) {
		log.Error("error msg trigger,msg:", string(message.Body), ",messageid:", message.MessageID)
		receiveChan <- string(message.Body)
		err = nil
		return
	}
}
