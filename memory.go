// Copyright 2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tracer

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type tracerEntry struct {
	When    time.Time `json:"when"`
	Message string    `json:"message"`
	Error   error     `json:"error,omitempty"`
	Data    []byte    `json:"data,omitempty"`
}

type memoryTracer struct {
	active  bool
	depth   int
	history []tracerEntry
	guard   sync.Mutex
	publish func([]byte)
}

func (t *memoryTracer) Stage(format string, args ...any) {
	if !t.active {
		entry := tracerEntry{
			When:    time.Now(),
			Message: fmt.Sprintf("STAGE: "+format, args...),
		}
		t.guard.Lock()
		t.active = true
		t.depth++
		t.history = append(t.history, entry)
		t.guard.Unlock()
	}
}

func (t *memoryTracer) Message(format string, args ...any) {
	if t.active {
		entry := tracerEntry{
			When:    time.Now(),
			Message: fmt.Sprintf("Msg: "+format, args...),
		}

		t.guard.Lock()
		t.history = append(t.history, entry)
		t.guard.Unlock()
	}
}

func (t *memoryTracer) Data(data any, format string, args ...any) {
	if t.active {
		entry := tracerEntry{
			When:    time.Now(),
			Message: fmt.Sprintf("DATA: "+format, args...),
		}

		switch actual := data.(type) {
		case []byte:
			entry.Data = actual
		case string:
			entry.Data = []byte(actual)
		default:
			if raw, err := json.Marshal(data); err != nil {
				entry.Data = []byte(fmt.Sprintf("error while marshalling: %v", err))
			} else {
				entry.Data = raw
			}
		}

		t.guard.Lock()
		t.history = append(t.history, entry)
		t.guard.Unlock()
	}
}

func (t *memoryTracer) Error(err error, format string, args ...any) {
	if t.active {
		entry := tracerEntry{
			When:    time.Now(),
			Message: fmt.Sprintf("ERROR: "+format, args...),
		}
		entry.Message += fmt.Sprintf(" (error: %v)", err)

		t.guard.Lock()
		t.history = append(t.history, entry)
		t.guard.Unlock()
	}
}

func (t *memoryTracer) Done(format string, args ...any) {
	entry := tracerEntry{
		When:    time.Now(),
		Message: fmt.Sprintf("DONE: "+format, args...),
	}
	t.guard.Lock()
	t.depth--
	t.history = append(t.history, entry)
	t.guard.Unlock()

	if t.depth == 0 {
		result := t.Render()
		t.publish(result)
	}
}

func (t *memoryTracer) Render() []byte {
	raw, err := json.Marshal(t.history)
	if err != nil {
		raw = []byte(fmt.Sprintf("error while marshalling: %v", err))
	}
	return raw
}

func Create() Tracer {
	return &memoryTracer{}
}
