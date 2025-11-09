// Copyright 2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// (originated from github.com/seamia/tracer)

package tracer

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type (
	tracerEntry struct {
		When     time.Time `json:"when"`
		Relative int64     `json:"relative.ns"`
		Message  string    `json:"message"`
		Error    error     `json:"error,omitempty"`
		Data     any       `json:"data,omitempty"`
	}

	traceHistory []tracerEntry

	memoryTracer struct {
		instance int64
		depth    int
		history  traceHistory
		guard    sync.Mutex
		publish  func([]byte)
	}
)

const (
	prefixStage   = "STAGE: "
	prefixMessage = "Msg: "
	prefixData    = "DATA: "
	prefixError   = "ERROR: "
	prefixDone    = "DONE: "
)

func (t *memoryTracer) Stage(format string, args ...any) {
	entry := tracerEntry{
		When:    time.Now(),
		Message: fmt.Sprintf(prefixStage+format, args...),
	}
	t.guard.Lock()
	t.depth++
	t.history = append(t.history, entry)
	t.guard.Unlock()
}

func (t *memoryTracer) Message(format string, args ...any) {
	entry := tracerEntry{
		When:    time.Now(),
		Message: fmt.Sprintf(prefixMessage+format, args...),
	}

	t.guard.Lock()
	t.history = append(t.history, entry)
	t.guard.Unlock()
}

func (t *memoryTracer) Data(data any, format string, args ...any) {
	entry := tracerEntry{
		When:    time.Now(),
		Message: fmt.Sprintf(prefixData+format, args...),
	}

	switch actual := data.(type) {
	case []byte:
		entry.Data = actual // Q: is assignment good enough? or shall we copy the data?
	case string:
		entry.Data = actual
	default:
		if raw, err := json.Marshal(actual); err != nil {
			entry.Data = fmt.Sprintf("failed to save: %T", actual)
			// entry.Data = []byte(fmt.Sprintf("error while marshalling: %v", err))
		} else {
			entry.Data = string(raw)
		}
	}

	t.guard.Lock()
	t.history = append(t.history, entry)
	t.guard.Unlock()
}

func (t *memoryTracer) Error(err error, format string, args ...any) {
	entry := tracerEntry{
		When:    time.Now(),
		Message: fmt.Sprintf(prefixError+format, args...),
	}
	entry.Message += fmt.Sprintf(" (error: %v)", err)

	t.guard.Lock()
	t.history = append(t.history, entry)
	t.guard.Unlock()
}

func (t *memoryTracer) Done(format string, args ...any) {
	entry := tracerEntry{
		When:    time.Now(),
		Message: fmt.Sprintf(prefixDone+format, args...),
	}
	t.guard.Lock()
	t.depth--
	t.history = append(t.history, entry)
	t.guard.Unlock()

	if t.depth == 0 {
		result := t.Render()
		t.publish(result)
		removeTracer(t.instance)
	}
}

func (t *memoryTracer) Render() []byte {
	if len(t.history) > 0 {
		start := t.history[0].When
		for index, entry := range t.history {
			t.history[index].Relative = entry.When.Sub(start).Nanoseconds()
		}
	}

	raw, err := json.MarshalIndent(t.history, "", "\t")
	if err != nil {
		raw = []byte(fmt.Sprintf("error while marshalling: %v", err))
	}
	return raw
}

func Create(publish PublishFunc) Tracer {
	instance := getUniqueInstanceID()
	var result Tracer = &memoryTracer{
		instance: instance,
		publish:  publish,
	}

	addTracer(instance, result)

	return result
}
