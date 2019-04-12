package main

import (
	"sync"
)

type LockingMap struct {
	words map[string]int
	mux   *sync.Mutex
}

func (l *LockingMap) Listen() {
}

func (l *LockingMap) Stop() {
}

func (l *LockingMap) Reduce(functor ReduceFunc, accum_str string, accum_int int) (string, int) {
	l.mux.Lock()
	for k, v := range l.words {
		accum_str, accum_int = functor(accum_str, accum_int, k, v)
	}
	l.mux.Unlock()
	return accum_str, accum_int
}

func (l *LockingMap) AddWord(word string) {
	l.mux.Lock()
	l.words[word]++
	l.mux.Unlock()
}

func (l *LockingMap) GetCount(word string) (out int) {
	l.mux.Lock()
	out = l.words[word]
	l.mux.Unlock()
	return
}

func NewLockingMap() *LockingMap {
	lm := new(LockingMap)
	lm.words = make(map[string]int)
	lm.mux = &sync.Mutex{}
	return lm
}
