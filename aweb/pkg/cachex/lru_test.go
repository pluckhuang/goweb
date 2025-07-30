package cachex

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	// 创建一个容量为 2 的 LRU 缓存
	cache := NewLRUCache(2)

	// 插入两个键值对
	cache.Put(1, 10) // key=1, value=10
	cache.Put(2, 20) // key=2, value=20

	// 验证插入后可以正确获取值
	if val := cache.Get(1); val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	// 访问 key=1，更新其最近使用时间
	cache.Get(1)

	// 插入第三个键值对，触发淘汰
	cache.Put(3, 30) // key=3, value=30

	// 验证 key=2 被淘汰（最近最少使用）
	if val := cache.Get(2); val != -1 {
		t.Errorf("Expected -1, got %d", val)
	}

	// 验证 key=1 和 key=3 仍然存在
	if val := cache.Get(1); val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}
	if val := cache.Get(3); val != 30 {
		t.Errorf("Expected 30, got %d", val)
	}
}
