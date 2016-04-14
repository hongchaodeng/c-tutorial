package demo1

import "time"

type Cache interface {
	Get(url string) (string, bool)
	Put(url string, body string)
}

type item struct {
	body      string
	timestamp time.Time
}

func newLRUCache(size int) *lruCache {
	return &lruCache{
		size:      size,
		urlToItem: make(map[string]*item),
	}
}

type lruCache struct {
	size      int
	urlToItem map[string]*item
}

func (lc *lruCache) Get(url string) (string, bool) {
	if lc.urlToItem[url] != nil {
		lc.urlToItem[url].timestamp = time.Now()
		return lc.urlToItem[url].body, true
	}
	return "", false
}

func (lc *lruCache) Put(url string, body string) {
	if lc.Full() {
		lc.evictOne()
	}
	lc.set(url, body)
}

func (lc *lruCache) Full() {
	return len(lc.urlToItem) == lc.size
}

func (lc *lruCache) evictOne() {
	var toEvict string
	var evictItem *item
	// find the item with earliest timestamp, O(n)
	for url, item := range lc.urlToItem {
		if evictItem == nil || item.timestamp.Before(evictItem.timestamp) {
			toEvict = url
			evictItem = item
		}
	}
	if evictItem == nil {
		panic("It's not possible! REWRITE THE CODE!")
	}
	delete(lc.urlToItem, toEvict)
}

func (lc *lruCache) set(url string, body string) {
	lc.urlToItem[url] = &item{
		body:      body,
		timestamp: time.Now(),
	}
}

// --------------------
// Efficient algorithm
// --------------------

// type item struct {
// 	url  string
// 	body string
// }

// func newLRUCache(size int) *lruCache {
// 	return &lruCache{
// 		size:          size,
// 		urlToListElem: make(map[string]*list.Element),
// 		l:             list.New(),
// 	}
// }

// type lruCache struct {
// 	size          int
// 	urlToListElem map[string]*list.Element
// 	l             *list.List
// }

// func (lc *lruCache) Get(url string) (string, bool) {
// 	if lc.urlToListElem[url] != nil {
// 		lc.l.MoveToFront(lc.urlToListElem[url])
// 		return lc.urlToListElem[url].Value.(*item).body, true
// 	}
// 	return "", false
// }

// func (lc *lruCache) Put(url string, body string) {
// 	if len(lc.urlToListElem) == lc.size {
// 		lc.evictOne()
// 	}
// 	lc.set(url, body)
// }

// func (lc *lruCache) evictOne() {
// 	delete(lc.urlToListElem, lc.l.Back().Value.(*item).url)
// 	lc.l.Remove(lc.l.Back())
// }

// func (lc *lruCache) set(url string, body string) {
// 	lc.urlToListElem[url] = lc.l.PushFront(&item{
// 		url:  url,
// 		body: body,
// 	})
// }
