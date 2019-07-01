package ezcache

import (
	"time"
)

type expiredItem struct {
	key  string
	time time.Time

	next *expiredItem
	prev *expiredItem
}

func (item *expiredItem) remove() {
	item_prev := item.prev
	if item_prev != nil {
		item_prev.next = item.next
	}

	item_next := item.next
	if item_next != nil {
		item_next.prev = item.prev
	}

	item.next = nil
	item.prev = nil

}

func (item *expiredItem) insertAfter(new *expiredItem) {
	item_next := item.next
	item.next, new.prev = new, item
	if item_next != nil {
		item.next.prev = new
		new.next = item_next
	}
}

func (item *expiredItem) insertBefore(new *expiredItem) {
	item_prev := item.prev
	item.prev, new.next = new, item
	if item_prev != nil {
		item_prev.next = new
		new.prev = item_prev
	}
}

func (c *Cache) SetExpiredSpace(t time.Duration) {
	c.expiredSpace = t
}

// 由于采用了标记的方式清楚expire, 所以可能会出现一种情况:
// 清理后, list中残留了一个键为[key]的expiredItem, 但后面只再次Set了[key]的值, 但没有Expire, 会导致expiredItem仍然生效
// 因此, 使用的时候需要注意这点
func (c *Cache) Expire(key string, duration time.Duration) {
	var (
		start *expiredItem // 起始节点
		right bool         // 是否往右, 即时间递增的方向
	)
	newTime := time.Now().Add(duration)
	newItem, existed := c.expiredSet[key]
	if existed {
		if newTime.Equal(newItem.time) {
			return
		}
		if newTime.After(newItem.time) {
			start = newItem.next
			right = true
		} else {
			start = newItem.prev
		}
		newItem.time = newTime
		if start == nil {
			// 不需要处理
			return
		}

		if newItem == c.expiredRoot {
			c.expiredRoot = newItem.next
		}
		newItem.remove()

	} else {
		newItem = &expiredItem{
			key:  key,
			time: newTime,
		}

		c.mu.Lock()
		c.expiredSet[key] = newItem
		c.mu.Unlock()
		// 第一个
		if c.expiredRoot == nil {
			c.expiredRoot = newItem
			return
		}

		start = c.expiredRoot
		right = true
	}

	if right {
		cur := start
		for {
			if newItem.time.Before(cur.time) {
				cur.insertBefore(newItem)
				if c.expiredRoot == cur {
					c.expiredRoot = newItem
				}
				break
			}

			if cur.next == nil {
				cur.insertAfter(newItem)
				break
			}

			cur = cur.next
		}
	} else {
		cur := start
		for {
			if newItem.time.After(cur.time) {
				cur.insertAfter(newItem)
				break
			}

			if cur.prev == nil {
				cur.insertBefore(newItem)
				if c.expiredRoot == cur {
					c.expiredRoot = newItem
				}
				break
			}

			cur = cur.prev
		}
	}
}

func (c *Cache) cleanExpireKey(key string) {

}

func (c *Cache) cleanExpireItem(item *expiredItem) {
	item.remove()
	c.mu.Lock()
	delete(c.expiredSet, item.key)
	c.mu.Unlock()
	c.Delete(item.key)
}

func (c *Cache) goCleanExpireList() {
	go func() {
		for {
			c.cleanExpireList()
			time.Sleep(c.expiredSpace)
		}
	}()
}

func (c *Cache) cleanExpireList() {
	if c.expiredRoot == nil {
		return
	}

	now := time.Now()
	cur := c.expiredRoot
	for {
		if now.Before(cur.time) {
			break
		}
		//
		c.expiredRoot = cur.next
		c.cleanExpireItem(cur)

		if cur.next == nil {
			break
		}

		cur = cur.next
	}
}
