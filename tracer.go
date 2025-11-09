// Copyright 2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// (originated from github.com/seamia/tracer)

package tracer

type Tracer interface {
	Stage(string, ...any)

	Message(string, ...any)
	Data(any, string, ...any)
	Error(error, string, ...any)

	Done(string, ...any)
}
