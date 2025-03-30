package cache

import (
	"strconv"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	// 创建一个缓存，过期时间100ms，清理间隔200ms
	cache := NewMemoryCache(100*time.Millisecond, 200*time.Millisecond)
	defer cache.Close()

	// 测试Set和Get
	t.Run("Set and Get", func(t *testing.T) {
		cache.Set("key1", "value1")
		cache.Set("key2", 2)
		cache.Set("key3", struct{ Name string }{"test"})

		// 检查值是否正确
		if val, found := cache.Get("key1"); !found || val.(string) != "value1" {
			t.Errorf("Expected key1=value1, got %v, found=%v", val, found)
		}

		if val, found := cache.Get("key2"); !found || val.(int) != 2 {
			t.Errorf("Expected key2=2, got %v, found=%v", val, found)
		}

		if val, found := cache.Get("key3"); !found || val.(struct{ Name string }).Name != "test" {
			t.Errorf("Expected key3.Name=test, got %v, found=%v", val, found)
		}

		// 检查不存在的键
		if _, found := cache.Get("not_exists"); found {
			t.Error("Expected not_exists to not be found")
		}
	})

	// 测试Delete
	t.Run("Delete", func(t *testing.T) {
		cache.Set("key_to_delete", "value")
		if _, found := cache.Get("key_to_delete"); !found {
			t.Error("Expected key_to_delete to be found before deletion")
		}

		cache.Delete("key_to_delete")
		if _, found := cache.Get("key_to_delete"); found {
			t.Error("Expected key_to_delete to not be found after deletion")
		}
	})

	// 测试过期
	t.Run("Expiration", func(t *testing.T) {
		cache.SetWithExpiration("expire_key", "value", 50*time.Millisecond)
		if _, found := cache.Get("expire_key"); !found {
			t.Error("Expected expire_key to be found before expiration")
		}

		// 等待项过期
		time.Sleep(100 * time.Millisecond)
		if _, found := cache.Get("expire_key"); found {
			t.Error("Expected expire_key to not be found after expiration")
		}
	})

	// 测试计数
	t.Run("Count", func(t *testing.T) {
		cache.Clear()
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		if count := cache.Count(); count != 2 {
			t.Errorf("Expected count=2, got %d", count)
		}

		cache.Set("key3", "value3")
		if count := cache.Count(); count != 3 {
			t.Errorf("Expected count=3, got %d", count)
		}

		cache.Delete("key1")
		if count := cache.Count(); count != 2 {
			t.Errorf("Expected count=2 after deletion, got %d", count)
		}

		cache.Clear()
		if count := cache.Count(); count != 0 {
			t.Errorf("Expected count=0 after clear, got %d", count)
		}
	})

	// 测试自动清理
	t.Run("Auto Cleanup", func(t *testing.T) {
		cleanupCache := NewMemoryCache(50*time.Millisecond, 100*time.Millisecond)
		defer cleanupCache.Close()

		cleanupCache.Set("key1", "value1")
		cleanupCache.Set("key2", "value2")

		// 等待自动清理
		time.Sleep(200 * time.Millisecond)

		if _, found := cleanupCache.Get("key1"); found {
			t.Error("Expected key1 to be automatically cleaned up")
		}

		if _, found := cleanupCache.Get("key2"); found {
			t.Error("Expected key2 to be automatically cleaned up")
		}
	})

	// 测试永不过期的缓存项
	t.Run("Never Expire", func(t *testing.T) {
		cache.Clear()
		cache.SetWithExpiration("never_expire", "value", -1)

		// 等待正常过期时间
		time.Sleep(150 * time.Millisecond)

		// 验证项目仍然存在
		if val, found := cache.Get("never_expire"); !found || val.(string) != "value" {
			t.Errorf("Expected never_expire to still exist with value='value', got %v, found=%v", val, found)
		}
	})

	// 测试缓存覆盖
	t.Run("Cache Override", func(t *testing.T) {
		cache.Clear()
		cache.Set("override_key", "original")

		// 验证原始值
		if val, found := cache.Get("override_key"); !found || val.(string) != "original" {
			t.Errorf("Expected override_key=original, got %v, found=%v", val, found)
		}

		// 覆盖值
		cache.Set("override_key", "updated")

		// 验证更新的值
		if val, found := cache.Get("override_key"); !found || val.(string) != "updated" {
			t.Errorf("Expected override_key=updated, got %v, found=%v", val, found)
		}
	})
}

// 测试缓存创建时的默认值
func TestNewMemoryCache(t *testing.T) {
	// 测试默认过期时间
	t.Run("Default Expiration", func(t *testing.T) {
		cache := NewMemoryCache(0, 0)
		defer cache.Close()

		// 默认过期时间应为1小时
		cache.Set("key", "value")

		// 应能正常获取值
		if val, found := cache.Get("key"); !found || val.(string) != "value" {
			t.Errorf("Expected key=value with default expiration, got %v, found=%v", val, found)
		}
	})

	// 测试无清理间隔
	t.Run("No Cleanup Interval", func(t *testing.T) {
		cache := NewMemoryCache(50*time.Millisecond, 0)
		defer cache.Close()

		cache.Set("key", "value")

		// 等待项过期
		time.Sleep(100 * time.Millisecond)

		// 尽管项已过期，但没有自动清理，Get仍然会检查过期时间
		if _, found := cache.Get("key"); found {
			t.Error("Expected expired key to not be found even without cleanup")
		}
	})
}

// 测试多个并发清理
func TestMultipleCleanupRoutines(t *testing.T) {
	cache := NewMemoryCache(50*time.Millisecond, 20*time.Millisecond)

	// 添加一些项
	for i := 0; i < 5; i++ {
		cache.Set(strconv.Itoa(i), i)
	}

	// 等待一段时间，让清理程序运行多次
	time.Sleep(150 * time.Millisecond)

	// 关闭缓存
	cache.Close()

	// 在多次清理后所有项应该都过期并被删除
	if count := cache.Count(); count != 0 {
		t.Errorf("Expected all items to be cleaned up, but found %d items", count)
	}
}

// 测试关闭和重复关闭
func TestClose(t *testing.T) {
	cache := NewMemoryCache(100*time.Millisecond, 200*time.Millisecond)
	cache.Set("key", "value")

	// 正常关闭
	cache.Close()

	// 验证仍然可以使用，但清理协程已停止
	cache.Set("key2", "value2")
	val, found := cache.Get("key2")
	if !found || val.(string) != "value2" {
		t.Error("Cache should still be usable after close")
	}

	// 再次关闭不应导致问题
	cache.Close()
}

// 测试缓存过期
func TestExpiredItemRemoval(t *testing.T) {
	cache := NewMemoryCache(50*time.Millisecond, 0)
	defer cache.Close()

	// 添加一些很快过期的项
	cache.Set("expire1", "value1")
	cache.Set("expire2", "value2")
	cache.SetWithExpiration("never_expire", "value3", -1)

	// 等待一些项过期
	time.Sleep(100 * time.Millisecond)

	// 验证过期的项已不可访问
	if _, found := cache.Get("expire1"); found {
		t.Error("Expected expire1 to be expired")
	}

	if _, found := cache.Get("expire2"); found {
		t.Error("Expected expire2 to be expired")
	}

	// 验证永不过期的项仍然存在
	if val, found := cache.Get("never_expire"); !found || val.(string) != "value3" {
		t.Errorf("Expected never_expire to still exist, got %v, found=%v", val, found)
	}
}
