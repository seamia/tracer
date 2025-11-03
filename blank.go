// Copyright 2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tracer

type blankTracer struct{}

func (*blankTracer) Stage(string, ...any)        {}
func (*blankTracer) Message(string, ...any)      {}
func (*blankTracer) Data(any, string, ...any)    {}
func (*blankTracer) Error(error, string, ...any) {}
func (*blankTracer) Done(string, ...any)         {}

func CreateBlankTracer() Tracer {
	return &blankTracer{}
}
