package concurrent

import (
	"fmt"
	"sync"
	"time"
)

var (
	greenLight = false             // 红绿灯状态
	mu         = sync.Mutex{}      // 互斥锁
	cond       = sync.NewCond(&mu) // 条件变量
)

func car(id int) {
	cond.L.Lock()
	for !greenLight {
		fmt.Printf("🚗 车辆 %d 等待红灯...\n", id)
		cond.Wait()
	}
	fmt.Printf("✅ 车辆 %d 通过绿灯！\n", id)
	cond.L.Unlock()
}

func Cond() {
	// 启动5辆车（5个等待的 goroutine）
	for i := 1; i <= 5; i++ {
		go car(i)
	}

	// 等待红灯 3 秒
	time.Sleep(3 * time.Second)

	// 改变红绿灯状态为绿灯
	cond.L.Lock()
	fmt.Println("🟢 绿灯了！通知所有车辆通行！")
	greenLight = true
	cond.Broadcast() // 广播通知所有等待的 goroutine
	cond.L.Unlock()

	// 等待所有 goroutine 打印完
	time.Sleep(2 * time.Second)
}

// 替代方式一：死循环自旋轮询变量
// 缺点：
// 多个协程都在不停轮询变量状态。

// 效率低，CPU 浪费。

// 即使加了 Sleep()，也不是实时响应。
func car1(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		mu.Lock()
		if greenLight {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(100 * time.Millisecond) // 避免吃满CPU
	}
	fmt.Printf("✅ 车辆 %d 通过绿灯（自旋检查）\n", id)
}

// 替代方式二：channel 广播（手动模拟）
// ⚠️ 注意：
// 只能广播一次，close(ch) 后所有接收者解除阻塞；

// 如果再想发第二次广播 → 不行，close 只能调用一次；

// 没法控制「只唤醒一个」或者「有条件唤醒」。

func car2(id int, ch <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ch // 阻塞等待消息
	fmt.Printf("✅ 车辆 %d 通过绿灯（channel）\n", id)
}

func main2() {
	ch := make(chan struct{}) // 无缓冲 channel
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go car(i, ch, &wg)
	}

	time.Sleep(2 * time.Second)
	close(ch) // 广播唤醒（所有接收者会 unblock）

	wg.Wait()
}

// | 能力         | channel     | sync.Cond    |
// | ---------- | ----------- | ------------ |
// | 点对点通信      | ✅           | ✅            |
// | 多个协程等待     | 🚫（要模拟）     | ✅            |
// | 条件控制逻辑     | ❌           | ✅ 支持复杂条件判断   |
// | 通知所有等待者    | ✅（通过 close） | ✅（Broadcast） |
// | 重复使用（多轮通知） | ❌           | ✅            |
