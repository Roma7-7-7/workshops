package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type DynamicWP struct {
	// number of workers
	min, max           int
	currentWorkerCount *int32
}

// fill the correct arguments
func (w *DynamicWP) work(ctx context.Context, workerTasks chan func()) {
	atomic.AddInt32(w.currentWorkerCount, 1)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker is canceled")
			atomic.AddInt32(w.currentWorkerCount, -1)
			return
		case task := <-workerTasks:
			task()
		case <-time.NewTimer(time.Second * 2).C:
			if atomic.LoadInt32(w.currentWorkerCount) > int32(w.min) {
				fmt.Println("no tasks")
				atomic.AddInt32(w.currentWorkerCount, -1)
				return
			}
		}
	}
}

// Start starts dynamic worker pull logic
func (w *DynamicWP) Start(ctx context.Context, tasksCh chan func()) {
	workerTasks := make(chan func())

	for i := 0; i < w.min; i++ {
		go w.work(ctx, workerTasks)
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("terminated")
			return
		case task := <-tasksCh:
		LOOP:
			for {
				select {
				case <-ctx.Done():
					fmt.Println("terminated")
					return
				case workerTasks <- task:
					break LOOP
				case <-time.NewTimer(time.Millisecond * 50).C:
					if atomic.LoadInt32(w.currentWorkerCount) < int32(w.max) {
						go w.work(ctx, workerTasks)
						workerTasks <- task
					}
				}
			}
		}
	}
}

func NewDynamicWorkerPool(min, max int) *DynamicWP {
	return &DynamicWP{
		min:                min,
		max:                max,
		currentWorkerCount: new(int32),
	}
}
