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
package schema

import "testing"

func TestEmailCodeContent_IsSourceType(t *testing.T) {
	unsub := &EmailCodeContent{SourceType: UnsubscribeSourceType, Email: "admin@example.com", UserID: "2"}
	if unsub.IsSourceType(AccountActivationSourceType, BindingSourceType) {
		t.Error("unsubscribe code must not be accepted for account activation")
	}
	if unsub.IsSourceType(PasswordResetSourceType) {
		t.Error("unsubscribe code must not be accepted for password reset")
	}
	if !unsub.IsSourceType(UnsubscribeSourceType) {
		t.Error("unsubscribe code must be accepted for unsubscribe")
	}
	legacy := &EmailCodeContent{Email: "admin@example.com", UserID: "2"}
	if legacy.IsSourceType(AccountActivationSourceType, PasswordResetSourceType, ConfirmNewEmailSourceType, UnsubscribeSourceType, BindingSourceType) {
		t.Error("untyped code must be rejected everywhere")
	}
	if ConfirmNewEmailSourceType == PasswordResetSourceType {
		t.Error("confirm-new-email and password-reset codes must be distinct")
	}
}
