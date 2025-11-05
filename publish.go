// Copyright 2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tracer

import (
	"fmt"
	"os"
	"path"
	"sync/atomic"
	"time"
)

type PublishFunc func([]byte)

func FileBasedPublish(location string) PublishFunc {

	os.MkdirAll(location, 0755)

	return func(data []byte) {
		fileName := path.Join(location, constructUniqueFileName())
		fmt.Printf("\t=== %s\n", fileName)
		os.WriteFile(fileName, data, 0755)
	}
}

var (
	uniqueFileNameCounter int64
)

func constructUniqueFileName() string {
	return fmt.Sprintf("trace-%v-%v", time.Now().UnixNano(), atomic.AddInt64(&uniqueFileNameCounter, 1))
}
