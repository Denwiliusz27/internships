package main

import (
	"fmt"
	"sync"
)

const DebugOut = false

type LRU interface {
	GetValue(int) (int, error)
	SetValue(int, int)
}

type lru struct {
	sync.RWMutex
	values    map[int]int
	stack     []int
	stackSize int
	wg        sync.WaitGroup
}

func debug(format string, args ...any) {
	if DebugOut {
		fmt.Printf(format+"\n", args...)
	}
}

func newLru(size int) *lru {
	return &lru{
		RWMutex:   sync.RWMutex{},
		values:    make(map[int]int, size),
		stack:     make([]int, 0, size),
		stackSize: size,
		wg:        sync.WaitGroup{},
	}
}

func (l *lru) GetValue(key int) (int, error) {
	l.RLock()
	defer l.RUnlock()

	debug("-- START getValue -- %d \n", key)

	val, ok := l.values[key]
	if !ok {
		return 0, fmt.Errorf("key doesn't exist")
	}

	debug("-- END getValue -- %d \n", key)

	return val, nil
}

func (l *lru) SetValue(key int, value int) {
	l.Lock()
	defer l.Unlock()

	debug("-- START setKeyValue -- %d \n", key)

	l.values[key] = value
	l.addToStack(key)

	debug("map: %d \n", l.values)
	debug("stack: %d \n", l.stack)
	debug("-- END setKeyValue -- %d \n", key)
}

func (l *lru) addToStack(val int) {
	debug("-- START addToStack -- %d \n", val)

	present := false
	for i, el := range l.stack {
		if el == val {
			present = true
			debug("deleted %d from middle of stack \n", el)
			l.stack = append(l.stack[:i], l.stack[i+1:]...)
			break
		}
	}

	if len(l.stack) == l.stackSize {
		if !present {
			key := l.stack[0]
			l.stack = l.stack[1:l.stackSize]
			delete(l.values, key)
			debug("-- Deleted %d from map\n", key)
		}
	}

	l.stack = append(l.stack, val)

	debug("-- END addToStack -- %d \n", val)
}

func fib(lru LRU, n int) int {
	if n <= 1 {
		return 1
	}

	if lru != nil {
		a, err := lru.GetValue(n)
		if err == nil {
			return a
		}

		val1 := fib(lru, n-1)
		val2 := fib(lru, n-2)

		lru.SetValue(n, val1+val2)

		return val1 + val2
	}

	return fib(lru, n-1) + fib(lru, n-2)
}
