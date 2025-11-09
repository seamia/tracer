// Copyright 2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// (originated from github.com/seamia/tracer)

package tracer

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	instanceCounter int64
	instanceGuard   sync.Mutex
	instanceStore   = map[int64]Tracer{}
)

func getUniqueInstanceID() int64 {
	return atomic.AddInt64(&instanceCounter, 1)
}

func addTracer(instance int64, result Tracer) {
	instanceGuard.Lock()
	instanceStore[instance] = result
	instanceGuard.Unlock()
}

func removeTracer(instance int64) {
	instanceGuard.Lock()
	delete(instanceStore, instance)
	instanceGuard.Unlock()
}

func FindOverdueTracers(maxAge time.Duration) {
	instanceGuard.Lock()
	for id, tracer := range instanceStore {
		if tracer, found := tracer.(*memoryTracer); found && tracer != nil {
			if len(tracer.history) > 0 {
				age := time.Now().Sub(tracer.history[0].When)
				if age > maxAge {
					tracer.Message("this object seems to be stalled, it's age is %v, which is larger than max.allowed.age of %v", age, maxAge)
					tracer.publish(tracer.Render())
					delete(instanceStore, id)
				}
			}
		}
	}
	instanceGuard.Unlock()
}
