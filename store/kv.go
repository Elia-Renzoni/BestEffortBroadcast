package store

import (
	"sync"
)

type Storage interface {
	Set(key string, value []byte)
	Get(key string) []byte
	Delete(key string)
}

type InMemoryKV struct {
	dataStructure map[string][]byte
	lock *sync.Mutex
}

func NewInMemoryKV() *InMemoryKV {
	return &InMemoryKV{
		dataStructure: make(map[string][]byte),
	}
}

func (i *InMemoryKV) Set(key string, value []byte) {
	if key == "" {
		return
	}

	i.dataStructure[key] = value
}

func (i *InMemoryKV) Get(key string) []byte {
	if key == "" {
		return nil
	}
	v, ok := i.dataStructure[key]
	if !ok {
		return nil	
	}
	return v
}

func (i *InMemoryKV) Delete(key string) {
	if key == "" {
		return
	}

	_, ok := i.dataStructure[key]
	if !ok {
		return
	}
	delete(i.dataStructure, key)
}

