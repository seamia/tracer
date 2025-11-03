// Copyright 2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tracer

import (
	"fmt"
	"os"
	"path"
	"time"
)

type PublishFunc func([]byte)

func FileBasedPublish(location string) PublishFunc {

	os.MkdirAll(location, 0755)

	return func(data []byte) {
		fileName := path.Join(location, fmt.Sprintf("trace-%v", time.Now().UnixNano()))
		os.WriteFile(fileName, data, 0755)
	}
}
