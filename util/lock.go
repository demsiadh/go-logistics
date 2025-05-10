package util

import "sync"

// LockPool Key-based Lock Pool
var LockPool = struct {
	mu sync.Mutex
	m  map[string]*sync.Mutex
}{m: make(map[string]*sync.Mutex)}

func GetOrderLock(orderID string) *sync.Mutex {
	LockPool.mu.Lock()
	defer LockPool.mu.Unlock()

	if _, exists := LockPool.m[orderID]; !exists {
		LockPool.m[orderID] = &sync.Mutex{}
	}
	return LockPool.m[orderID]
}

func GetVehicleLock(plateNumber string) *sync.Mutex {
	LockPool.mu.Lock()
	defer LockPool.mu.Unlock()

	if _, exists := LockPool.m[plateNumber]; !exists {
		LockPool.m[plateNumber] = &sync.Mutex{}
	}
	return LockPool.m[plateNumber]
}
