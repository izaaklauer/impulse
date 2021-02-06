package engine

import (
    "sync"
)

type PortAllocator struct {
    allocs []int
    mu     *sync.Mutex
}

func NewPortAllocator(start int, end int) *PortAllocator {
    allocs := make([]int, end-start + 1)
    for port := start; port <= end; port++ {
        allocs[port-start] = port
    }
    return &PortAllocator{
        allocs: allocs,
        mu:     &sync.Mutex{},
    }
}

func (pa *PortAllocator) GetPort() int {
    pa.mu.Lock()
    defer pa.mu.Unlock()

    port := pa.allocs[0]
    pa.allocs = pa.allocs[1:]
    return port
}

func (pa *PortAllocator) ReleasePort(port int) {
    pa.mu.Lock()
    defer pa.mu.Unlock()
    
    pa.allocs = append(pa.allocs, port)
    return
}
