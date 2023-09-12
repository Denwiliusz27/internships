package main

import "container/list"

type LRU struct {
	values map[int]int
	stack  list.List
}
