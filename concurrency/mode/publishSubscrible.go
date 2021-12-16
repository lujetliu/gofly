package main

/* 发布/订阅模型 */

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

/*
 * TODO: 加深理解, 应用场景实践
 * 发布/订阅模型通常称为 pub/sub 模型, 在这个模型中, ./producerConsumer.go 中的
 * 消息生产者成为发布者(publisher), 而消息消费者成为订阅者(subscriber), 生产
 * 者和消费者是 M:N 的关; 在传统生产者/消费者模型中, 是将消息发送到一个队列中,
 * 而发布订/订阅模型是将消息发布给一个主题.
 *
 */

type (
	subscriber chan interface{}         // 订阅者为一个通道
	topicFunc  func(v interface{}) bool // 主题为一个过滤器
)

// 发布者对象
type Publisher struct {
	m           sync.RWMutex  // 读写锁
	buffer      int           // 订阅队列的缓存大小
	timeout     time.Duration // 发布超时时间
	subscribers map[subscriber]topicFunc
}

func NewePublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 添加新的订阅者, 订阅全部主题
func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

// 添加新的订阅者, 订阅过滤器筛选后的主题
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
}

// 退出订阅
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// 发布一个主题
func (p *Publisher) Publish(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

func (p *Publisher) sendTopic(sub subscriber, topic topicFunc,
	v interface{}, wg *sync.WaitGroup) {

	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

func main() {
	p := NewePublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscribe()                                  // 订阅了所有主题
	golang := p.SubscribeTopic(func(v interface{}) bool { // 订阅了 golang 主题
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello, world")
	p.Publish("hello, golang")

	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	time.Sleep(3 * time.Second)
}

// 在发布/订阅模型中, 每条消息都会传送给多个订阅者, 发布者通常不会知道,
// 也不关心哪一个订阅者正在接收主题消息; 订阅者和发布者可以在运行时动态
// 添加, 它们之间是一种松散的耦合关系,这使得系统的复杂性可以随时间的推
// 移而增长.
