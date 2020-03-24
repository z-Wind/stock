package main

import (
	"sync"

	"github.com/z-Wind/stock/crawler"
)

type record struct {
	taskDone map[string]bool
	lock     sync.Mutex
}

// newRecord 建立 record
func newRecord() *record {
	var r record

	r.taskDone = make(map[string]bool)

	return &r
}

// isProcessedOrAdd 是否已處理，未處理就加入
// 不存在表示未處理，存在但 False 表示處理中，存在且 True 表示已處理
func (r *record) isProcessedOrAdd(req interface{}) bool {
	key := req.(crawler.Request).Item.(string)
	r.lock.Lock()
	_, ok := r.taskDone[key]
	if !ok {
		r.taskDone[key] = false
	}
	r.lock.Unlock()

	return ok
}

// isDone 是否完成
// 不存在表示未處理，存在但 False 表示處理中，存在且 True 表示已處理
func (r *record) isDone(req interface{}) bool {
	key := req.(crawler.Request).Item.(string)

	return r.taskDone[key]
}

// done 任務已完成
func (r *record) done(symbol string) {
	key := symbol

	r.lock.Lock()
	r.taskDone[key] = true
	r.lock.Unlock()
}
