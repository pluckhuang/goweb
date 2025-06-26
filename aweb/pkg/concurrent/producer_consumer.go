package concurrent

// 案例1：Goroutine与Channel实现生产者-消费者
// 场景：一个生产者生成数据，多个消费者并行处理。
// 说明：
// 生产者通过channel发送数据，消费者并行接收。

// 带缓冲通道（容量2）减少阻塞。

// close(ch)确保消费者知道数据结束。

import (
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i
		time.Sleep(time.Millisecond * 500)
	}
	close(ch) // 关闭通道，通知消费者数据发送完毕
}

func consumer(id int, ch <-chan int) {
	for num := range ch { // range自动检测通道关闭
		fmt.Printf("Consumer %d received: %d\n", id, num)
	}
}

func ProducerConsumer() {
	ch := make(chan int, 2) // 带缓冲通道
	go producer(ch)
	go consumer(1, ch)
	go consumer(2, ch)
	time.Sleep(time.Second * 3) // 等待goroutine完成
}
