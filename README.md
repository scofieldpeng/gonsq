# gonsq

一个简单的基于官方 nsq 的 go client 的封装

## 使用

一. 安装

```go
go get github.com/scofieldpeng/gonsq
```

二. 配置

```ini
# app.ini
# 生产者配置
[producer]
nsqd=127.0.0.1:4151

# 消费者配置
[consumer]
# 直接连接上 nsqd,支持多个,用半角逗号隔开
nsqd=127.0.0.1:4150
# 连接上 nsqlookupd, 支持多个,用半角逗号隔开
nsqlookupd=127.0.0.1:4161,127.0.0.2:4161
# 最大并行数量
maxInFlight=100
# 并行数量,该值不能超过 maxInFlight
concurrent=100
```

三. quick start

```go
package main

import(
    "go.zhuzi.me/config"	
    "github.com/scofieldpeng/gonsq"
    "github.com/nsqio/go-nsq"
)

func main(){
    producerConfig := config.Data("nsq").Section("producer")
    consumerConfig := config.Data("nsq").Section("consumer")
    debug := false
    run := make(chan bool)

    // 初始化 producer
    if err := gonsq.Producer.Init(producerConfig,debug);err != nil {
	    panic(err)
    }
    // 初始化 consumer
    if err := gonsq.Consumer.Init(consumerConfig,debug);err != nil {
    	panic(err)
    }
    // 添加消费者处理函数
    gonsq.Consumer.AddHandler("testTopic",defaultHandler())
    gonsq.Consumer.AddHandler("testTopic2",defaultHandler())
    
    // run 起来
    if err := gonsq.RunAll();err != nil {
        panic(err)	
    }
    // 不要忘记退出时关闭
    defer gonsq.StopAll()
    
    <- run
}

// defaultHandler 编写消费者处理函数
func defaultHandler() (handler nsq.HandlerFunc) {
    return func (nm *nsq.Message) error {
        // 具体 consumer 逻辑	
    	return nil
    }
}
```

四. 单独使用producer或者consumer

**单独使用producer**

```go
// 单独使用producer
// 初始化
if err := gonsq.Producer.Init(config,debug);err != nil {
	panic(err)
}
// run起来
if err := gonsq.Producer.Run();err != nil {
	panic(err)
}
// 别忘了defer
defer gonsq.Producer.Stop()

// 调用
// msg支持字符型，数字型，或者直接struct
if err := gonsq.Producer.Publish(topic,msg);err != nil {
	log.Error(err)
}
// 批量发布
if err := gonsq.Producer.MultiPublish(topic,msgArr);err != nil {
	log.Error(err)
}
```

**单独使用consumer**

```go
// 初始化
gonsq.Consumer.Init(config,debug)
if err := gonsq.Consumer.Run();err != nil {
	panic(err)
}
// 添加handler
gonsq.Consumer.AddHandler(topic,handler)
gonsq.Consumer.AddHandler(topic2,handler2)

// run起来
if err := gonsq.Consumer.Run();err != nil {
	panic(err)
}
// 别忘了defer
defer gonsq.Consumer.Stop()

```