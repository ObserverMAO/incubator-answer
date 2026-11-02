/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package checker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSupportedImageFile_ICOHeader(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]struct {
		body []byte
		ok   bool
	}{
		"real.ico":  {[]byte{0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x10, 0x10}, true},
		"probe.ico": {[]byte("PROBE7491_plaintext_ico"), false},
		"mz.ico":    {[]byte("MZ..---i3428---"), false},
		"empty.ico": {nil, false},
	}
	for name, c := range cases {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, c.body, 0o600); err != nil {
			t.Fatal(err)
		}
		if got := IsSupportedImageFile(p); got != c.ok {
			t.Errorf("%s: got %v, want %v", name, got, c.ok)
		}
	}
}
