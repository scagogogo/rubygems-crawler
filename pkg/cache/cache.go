// Package cache 提供了缓存接口和实现
// 用于存储和检索键值对数据，支持过期时间和自动清理
package cache

import (
	"sync"
	"time"
)

// Cache 定义了缓存的基本操作接口
// 该接口允许存储任意类型的值，并支持过期时间设置
type Cache interface {
	// Get 获取指定键的缓存值
	// 如果键存在且未过期，返回值和true
	// 如果键不存在或已过期，返回nil和false
	Get(key string) (interface{}, bool)

	// Set 设置缓存值
	// 使用默认的过期时间
	Set(key string, value interface{})

	// SetWithExpiration 设置缓存值并指定过期时间
	// 如果过期时间为0，则使用默认过期时间
	// 如果过期时间为负，则永不过期
	SetWithExpiration(key string, value interface{}, d time.Duration)

	// Delete 删除指定键的缓存
	Delete(key string)

	// Clear 清空所有缓存
	Clear()

	// Count 返回缓存中的项目数量
	Count() int

	// Close 关闭缓存，释放资源
	// 在不再使用缓存时应调用此方法
	Close()
}

// cacheItem 表示一个缓存项
type cacheItem struct {
	value      interface{} // 存储的值
	expiration time.Time   // 过期时间
	created    time.Time   // 创建时间
}

// MemoryCache 是Cache接口的内存实现
// 它将数据存储在内存中，并支持自动过期和定期清理
type MemoryCache struct {
	defaultExpiration time.Duration        // 默认过期时间
	cleanupInterval   time.Duration        // 清理周期
	items             map[string]cacheItem // 缓存项存储
	mu                sync.RWMutex         // 读写锁，保证并发安全
	stopCleanup       chan struct{}        // 停止清理的通道
	closed            bool                 // 缓存是否已关闭
}

// NewMemoryCache 创建一个新的内存缓存
// 参数:
//   - defaultExpiration: 默认的缓存项过期时间
//   - cleanupInterval: 自动清理过期项目的时间间隔
//
// 如果cleanupInterval为0，则不会自动清理过期项目
// 如果defaultExpiration为0，则使用1小时作为默认过期时间
func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) *MemoryCache {
	if defaultExpiration <= 0 {
		defaultExpiration = time.Hour
	}

	cache := &MemoryCache{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		items:             make(map[string]cacheItem),
		stopCleanup:       make(chan struct{}),
	}

	// 如果设置了清理间隔，启动自动清理
	if cleanupInterval > 0 {
		go cache.startCleanupTimer()
	}

	return cache
}

// Get 获取缓存值
// 如果键不存在或已过期，返回nil和false
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// 检查是否已过期
	if !item.expiration.IsZero() && item.expiration.Before(time.Now()) {
		return nil, false
	}

	return item.value, true
}

// Set 使用默认过期时间设置缓存值
func (c *MemoryCache) Set(key string, value interface{}) {
	c.SetWithExpiration(key, value, c.defaultExpiration)
}

// SetWithExpiration 设置缓存值并指定过期时间
// 如果d为0，使用默认过期时间
// 如果d为负数，则永不过期
func (c *MemoryCache) SetWithExpiration(key string, value interface{}, d time.Duration) {
	var expiration time.Time

	if d == 0 {
		d = c.defaultExpiration
	}

	// 如果持续时间为负，则永不过期
	if d > 0 {
		expiration = time.Now().Add(d)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: expiration,
		created:    time.Now(),
	}
}

// Delete 从缓存中删除指定键
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear 清空所有缓存项
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]cacheItem)
}

// Count 返回缓存中的项目数量
func (c *MemoryCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// Close 关闭缓存，停止自动清理
func (c *MemoryCache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		close(c.stopCleanup)
		c.closed = true
	}
}

// startCleanupTimer 启动定期清理过期项目的定时器
func (c *MemoryCache) startCleanupTimer() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// deleteExpired 删除所有过期的缓存项
func (c *MemoryCache) deleteExpired() {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, item := range c.items {
		// 如果设置了过期时间且已过期，则删除
		if !item.expiration.IsZero() && item.expiration.Before(now) {
			delete(c.items, k)
		}
	}
}
