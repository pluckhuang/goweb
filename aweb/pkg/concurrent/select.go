package concurrent

// 场景：模拟任务处理，超时退出。

// 说明：
// select监听channel和超时事件。

// 如果process在1秒内未完成，触发超时。

import (
	"fmt"
	"time"
)

func process(ch chan string, data string) {
	time.Sleep(time.Second * 2) // 模拟耗时操作
	ch <- data
}

func Select() {
	ch := make(chan string)
	go process(ch, "result")

	select {
	case res := <-ch:
		fmt.Println("Received:", res)
	case <-time.After(time.Second * 1): // 1秒超时
		fmt.Println("Timeout")
	}
}
