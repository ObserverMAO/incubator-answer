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
package user_external_login

import "testing"

func TestIsMixinBotAccount(t *testing.T) {
	cases := map[string]bool{
		// real payloads from production (trimmed); the first is an app-created network user
		`{"data":{"type":"user","user_id":"8a7fee01-d454-36d0-8291-738acdf05237","identity_number":"0","full_name":"qac2e2d7979b","app_id":"4e96fa4e-9331-4787-958b-84bab43b1df6"}}`: true,
		`{"data":{"type":"user","user_id":"b847a455-aa41-4f7d-8038-0aefbe40dcaa","identity_number":"31911","full_name":"crossle"}}`:                                                  false,
		``:            false,
		`not json`:    false,
		`{"data":{}}`: false,
	}
	for in, want := range cases {
		if got := isMixinBotAccount(in); got != want {
			t.Errorf("isMixinBotAccount(%q) = %v, want %v", in, got, want)
		}
	}
}
