package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if element, ok := c.items[key]; ok {
		c.queue.MoveToFront(element)
		element.Value.(*cacheItem).value = value
		return true
	}

	newItem := &cacheItem{key: key, value: value}
	element := c.queue.PushFront(newItem)
	c.items[key] = element

	if c.queue.Len() > c.capacity {
		c.removeOldest()
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if c.items == nil {
		return nil, false
	}

	if element, ok := c.items[key]; ok {
		c.queue.MoveToFront(element)
		return element.Value.(*cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.items = nil
	c.queue = nil
}

func (c *lruCache) removeOldest() {
	if element := c.queue.Back(); element != nil {
		c.queue.Remove(element)
		item := element.Value.(*cacheItem)
		delete(c.items, item.key)
	}
}
